package operation

import (
	"context"

	"github.com/tfaller/jsonball"
	"github.com/tfaller/jsonball/event"
	"github.com/tfaller/jsonball/internal/name"
	"github.com/tfaller/propchange"
)

// HandleChange handles an open change -> loads affected documents
func HandleChange(ctx context.Context, registry jsonball.Registry, change propchange.OnChange) (*event.Change, error) {
	changeEvent := &event.Change{
		Handler: name.ParseListenerName(change.Listener()),
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
			return nil, err
		}
		changedDocs = append(changedDocs, d)
	}

	changeEvent.Documents = changedDocs
	return changeEvent, nil
}
