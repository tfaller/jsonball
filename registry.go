package jsonball

import (
	"context"
	"errors"

	"github.com/tfaller/jsonball/event"
)

// Change represents a given change
type Change struct {
	Document string
}

// DocumentList is the result of a ListDocuments operation
type DocumentList struct {
	Documents []string
	NextToken string
}

// ErrDocumentNotExist indicates that an operation failed because
// a given document does not exist.
var ErrDocumentNotExist = errors.New("document does not exist")

// Registry provides backend features to store documents and handles.
type Registry interface {
	// OpenDocument opens a new or existing document to update it.
	// To simply read a document use GetDocument instead.
	OpenDocument(ctx context.Context, docType, name string) (DocOps, error)

	// GetDocument returns the current latest version
	// Note: Currently running "PutDocument" operation
	// will not affect this method. This method is NOT blocked
	// by a running "PutDocument" method. But it is not guaranteed, that
	// this document is the actual latest Version! It is possible that
	// the upstream system change the doc, but that change was not yet
	// published to and processed by the registry.
	GetDocument(ctx context.Context, docType, name string) (string, error)

	// GetNewDocHanders returns all handlers that should be triggered for a given
	// document type. This function is used to determin which listerns should be automatically
	// generated if a new document was created.
	GetNewDocHanders(ctx context.Context, docType string) ([]string, error)

	// GetHandlerQueueURL returns the SQS queue URL. This func is called
	// if a change for a given handler was found.
	GetHandlerQueueURL(ctx context.Context, handler string) (string, error)

	// Registers a document type in the registry
	RegisterDocumentType(ctx context.Context, docType string) error

	// RegisterHandler registers a new handler
	RegisterHandler(ctx context.Context, handler *event.RegisterHandler) error

	// ListDocuments lists documents of a given type. With the parameter startToken pagination is
	// possible. Either startToken is empty or the value of "NextToken" of a previous list operation.
	// Note: There is no guarantee that documents are listed which were created after the NextToken
	// was build. They might be listed or might not be.
	ListDocuments(ctx context.Context, docType, startToken string, maxDocs uint16) (*DocumentList, error)

	// HandlerNewDoc registers or de-registers that a handler should
	// trigger for a new document of given type.
	HandlerNewDoc(ctx context.Context, handler, docType string, register bool) error
}

// DocOps handles the operations possible on an opened document
type DocOps interface {
	// IsNew tells wether this doc was created because it got
	// opened now.
	IsNew() bool

	// Document returns the current json document.
	Document() string

	// Change updates the document and implicitely closes the document.
	Change(change Change) error

	// Refresh marks the document as recently re-checked for changes,
	// but no changes where found. Implicitely closes the document.
	Refresh() error

	// Close closes the document.
	Close() error
}
