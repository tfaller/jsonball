package operation

import (
	"context"
	"encoding/json"

	"github.com/tfaller/jsonball"
	"github.com/tfaller/jsonball/event"
)

// GetDocumentContent gets the actual content of a document
func GetDocumentContent(ctx context.Context, registry jsonball.Registry, docType, docName string) (d event.Document, err error) {
	doc, err := registry.GetDocument(ctx, docType, docName)
	if err != nil {
		return
	}

	var rawEraDoc = struct {
		Doc    json.RawMessage `json:"doc"`
		DocEra uint32          `json:"docEra"`
	}{}

	err = json.Unmarshal([]byte(doc), &rawEraDoc)
	if err != nil {
		return
	}

	return event.Document{
		Type:     docType,
		Name:     docName,
		Version:  rawEraDoc.DocEra,
		Document: rawEraDoc.Doc,
	}, nil
}
