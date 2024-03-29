package operation

import (
	"context"
	"errors"

	"github.com/tfaller/jsonball"
	"github.com/tfaller/jsonball/event"
	"github.com/tfaller/jsonball/internal/name"
	"github.com/tfaller/propchange"
)

// HandleChange handles an open change -> loads affected documents
func HandleChange(ctx context.Context, registry jsonball.Registry, change propchange.OnChange) (*event.Change, error) {
	lHandler, lName := name.ParseListenerName(change.Listener())

	changeEvent := &event.Change{
		Handler: lHandler,
		Name:    lName,
	}

	docs := change.Documents()
	changedDocs := make([]event.Document, 0, len(docs))

	// doc current document
	for _, doc := range docs {
		docType, docName, err := name.ParseDocName(doc)
		if err != nil {
			return nil, err
		}

		if name.IsDocTypeInternal(docType) {
			// Ignore this document. Just there for
			// internal handling of things.
			continue
		}

		d, err := GetDocumentContent(ctx, registry, docType, docName)
		if err != nil {
			if !errors.Is(err, jsonball.ErrDocumentNotExist) {
				return nil, err
			}
			d = event.Document{Type: docType, Name: docName}
		}
		changedDocs = append(changedDocs, d)
	}

	changeEvent.Documents = changedDocs
	return changeEvent, nil
}
