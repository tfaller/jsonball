package mysql

import (
	"context"
	"crypto/cipher"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/tfaller/go-sqlprepare"
	"github.com/tfaller/jsonball"
	"github.com/tfaller/jsonball/event"
	"golang.org/x/crypto/chacha20poly1305"
)

// Registry is a registry with a mysql backend
type Registry struct {
	m  sync.Mutex
	db *sql.DB

	stmtOpenDoc       *sql.Stmt
	stmtNewDoc        *sql.Stmt
	stmtGetDoc        *sql.Stmt
	stmtDocType       *sql.Stmt
	stmtDocTypeReg    *sql.Stmt
	stmtDocRefresh    *sql.Stmt
	stmtDocUpdate     *sql.Stmt
	stmtHandlerQueue  *sql.Stmt
	stmtHandlerNewDoc *sql.Stmt
	stmtHandlerReg    *sql.Stmt

	docType map[string]uint64
	docAead cipher.AEAD
}

// Document is an open document
type Document struct {
	tx  *sql.Tx
	reg *Registry

	isNew        bool
	ID           uint64
	Type         uint64
	Name         string
	document     string
	RegisteredAt *time.Time
	RefreshedAt  *time.Time
}

// NewRegistry opens a registry based on a connection string
func NewRegistry(cs string, key []byte) (*Registry, error) {
	if !strings.Contains(cs, "parseTime=true") {
		return nil, fmt.Errorf("Required option \"parseTime=true\" is missing in connection string")
	}

	m, err := &Registry{docType: map[string]uint64{}}, error(nil)

	if m.db, err = sql.Open("mysql", cs); err != nil {
		return nil, err
	}

	if err = m.prepare(); err != nil {
		return nil, err
	}

	m.docAead, err = chacha20poly1305.NewX(key)
	if err != nil {
		return nil, fmt.Errorf("can't init document encryption: %w", err)
	}

	return m, nil
}

func (r *Registry) prepare() error {
	return sqlprepare.Prepare(r.db, []sqlprepare.ToPrepare{
		{Name: "open-doc", Target: &r.stmtOpenDoc,
			Query: "SELECT * FROM document WHERE type = ? AND name = ? FOR UPDATE"},

		{Name: "get-doc", Target: &r.stmtGetDoc,
			Query: "SELECT id, document, encrypted FROM document WHERE type = ? AND name = ? FOR SHARE"},

		{Name: "new-doc", Target: &r.stmtNewDoc,
			Query: "INSERT INTO document (type, name, document) VALUES (?, ?, ?)"},

		{Name: "refresh-doc", Target: &r.stmtDocRefresh,
			Query: "UPDATE document SET refreshedat = NOW() WHERE id = ?"},

		{Name: "update-doc", Target: &r.stmtDocUpdate,
			Query: "UPDATE document SET document = ?, refreshedat = NOW(), encrypted = 1 WHERE id = ?"},

		{Name: "doc-type", Target: &r.stmtDocType,
			Query: "SELECT id FROM document_type WHERE name = ?"},

		{Name: "doc-type-reg", Target: &r.stmtDocTypeReg,
			Query: "INSERT INTO document_type (name) VALUES (?)"},

		{Name: "handler-queueurl", Target: &r.stmtHandlerQueue,
			Query: "SELECT queueurl FROM handler WHERE name = ?"},

		{Name: "handler-register", Target: &r.stmtHandlerReg,
			Query: "INSERT INTO handler (name, queueurl) VALUES (?, ?)"},

		{Name: "handler-newdoc", Target: &r.stmtHandlerNewDoc,
			Query: "SELECT h.name FROM document_type dt JOIN document_type_handler dth ON dth.doctype = dt.id JOIN handler h ON h.id = dth.handler WHERE dt.name = ?"},
	}...)
}

func (r *Registry) getTypeID(docType string) (uint64, error) {
	r.m.Lock()
	defer r.m.Unlock()

	tID, isSet := r.docType[docType]
	if isSet {
		return tID, nil
	}

	// type is not cached ... read from DB
	row := r.stmtDocType.QueryRow(docType)
	err := row.Scan(&tID)
	if err != nil {
		return 0, fmt.Errorf("Can't find doctype %v: %w", docType, err)
	}

	// cache it
	r.docType[docType] = tID

	return tID, nil
}

// GetHandlerQueueURL get the queue URL of an given handler
func (r *Registry) GetHandlerQueueURL(ctx context.Context, name string) (url string, err error) {
	err = r.stmtHandlerQueue.QueryRowContext(ctx, name).Scan(&url)
	return
}

// OpenDocument opens a new or existing document.
func (r *Registry) OpenDocument(ctx context.Context, docType, name string) (jsonball.DocOps, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("can't open transaction for PutDocument: %w", err)
	}

	ops, err := r.doOpenDocument(docType, name, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return ops, err
}

func (r *Registry) doOpenDocument(docType, name string, tx *sql.Tx) (jsonball.DocOps, error) {
	// we need the doctype ... not the name
	docTypeID, err := r.getTypeID(docType)
	if err != nil {
		return nil, fmt.Errorf("can't find docType: %w", err)
	}

	doc, err := r.tryOpenDoc(name, docTypeID, tx)
	if errors.Is(err, sql.ErrNoRows) {
		// create a new document
		doc, err = r.doNewDoc(name, docTypeID, tx)
	}
	if err != nil {
		return nil, fmt.Errorf("Can't open document: %w", err)
	}
	return doc, nil
}

func (r *Registry) tryOpenDoc(name string, typeID uint64, tx *sql.Tx) (*Document, error) {
	stmt := tx.Stmt(r.stmtOpenDoc)
	defer stmt.Close()

	row := stmt.QueryRow(typeID, name)
	doc := &Document{
		reg: r,
		tx:  tx,
	}

	var encrypted bool
	err := row.Scan(
		&doc.ID,
		&doc.Type,
		&doc.Name,
		&doc.document,
		&doc.RegisteredAt,
		&doc.RefreshedAt,
		&encrypted,
	)

	if encrypted {
		plainDoc, err := decryptDoc(r.docAead, doc.ID, doc.Type, doc.Name, []byte(doc.document))
		if err != nil {
			return nil, err
		}
		doc.document = string(plainDoc)
	}

	return doc, err
}

func (r *Registry) doNewDoc(name string, typeID uint64, tx *sql.Tx) (*Document, error) {
	stmt := tx.Stmt(r.stmtNewDoc)
	defer stmt.Close()

	res, err := stmt.Exec(typeID, name, "null")
	if err != nil {
		return nil, fmt.Errorf("Can't insert new doc: %w", err)
	}

	docID, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("Can't get ID of new doc: %w", err)
	}

	return &Document{
		ID:           uint64(docID),
		Name:         name,
		Type:         typeID,
		document:     "null",
		RefreshedAt:  nil,
		RegisteredAt: nil,
		isNew:        true,
		reg:          r,
		tx:           tx,
	}, nil
}

// GetDocument gets the current document content
func (r *Registry) GetDocument(ctx context.Context, docType, name string) (string, error) {
	// resolve the doc type first
	docTypeID, err := r.getTypeID(docType)
	if err != nil {
		return "", err
	}

	encrypted := false
	docID := uint64(0)
	doc := ""
	row := r.stmtGetDoc.QueryRow(docTypeID, name)

	if err = row.Scan(&docID, &doc, &encrypted); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", jsonball.ErrDocumentNotExist
		}
		return "", err
	}

	if encrypted {
		plainDoc, err := decryptDoc(r.docAead, docID, docTypeID, name, []byte(doc))
		if err != nil {
			return "", err
		}
		doc = string(plainDoc)
	}

	return doc, nil
}

// GetNewDocHanders gets handlers that want to handle new document of a given type
func (r *Registry) GetNewDocHanders(ctx context.Context, docType string) ([]string, error) {
	rows, err := r.stmtHandlerNewDoc.QueryContext(ctx, docType)
	if err != nil {
		return nil, err
	}

	var handlers []string
	for rows.Next() {
		var handler string
		if err = rows.Scan(&handler); err != nil {
			return nil, err
		}
		handlers = append(handlers, handler)
	}

	return handlers, nil
}

// RegisterDocumentType registers a new document type
func (r *Registry) RegisterDocumentType(ctx context.Context, docType string) error {
	r.m.Lock()
	defer r.m.Unlock()

	res, err := r.stmtDocTypeReg.ExecContext(ctx, docType)
	if err != nil {
		return err
	}

	// cache document type
	docTypeID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	r.docType[docType] = uint64(docTypeID)

	return nil
}

// RegisterHandler registers a new handler
func (r *Registry) RegisterHandler(ctx context.Context, handler *event.RegisterHandler) error {
	_, err := r.stmtHandlerReg.ExecContext(ctx, handler.Handler, handler.QueueURL)
	return err
}

// Change updates a document
func (d *Document) Change(change jsonball.Change) error {
	// the document updated
	stmt := d.tx.Stmt(d.reg.stmtDocUpdate)
	defer stmt.Close()

	doc, err := encryptDoc(d.reg.docAead, d.ID, d.Type, d.Name, []byte(change.Document))
	if err != nil {
		return err
	}

	_, err = stmt.Exec(doc, d.ID)
	if err != nil {
		return fmt.Errorf("Failed to update doc: %w", err)
	}

	return d.tx.Commit()
}

// Refresh marks the document as up-to-date even no ne changed was happened
func (d *Document) Refresh() error {
	stmt := d.tx.Stmt(d.reg.stmtDocRefresh)
	defer stmt.Close()

	_, err := stmt.Exec(d.ID)
	if err != nil {
		return fmt.Errorf("could not refresh document: %w", err)
	}

	return d.tx.Commit()
}

// Document gets the current document
func (d *Document) Document() string {
	return d.document
}

// Close coloses the current document.
// No further operation are possible
func (d *Document) Close() error {
	return d.tx.Rollback()
}

// IsNew tells whether the document is a new one
func (d *Document) IsNew() bool {
	return d.isNew
}

var _ jsonball.Registry = &Registry{}
