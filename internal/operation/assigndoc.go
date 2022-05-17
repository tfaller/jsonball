package operation

import (
	"context"
	"fmt"

	"github.com/tfaller/jsonball"
	"github.com/tfaller/jsonball/event"
	"github.com/tfaller/propchange"
)

func assignDocs(ctx context.Context, registry jsonball.Registry, detector propchange.Detector, handler string, docs []event.GetDocument) error {

	// simple check whether the handler exists
	_, err := registry.GetHandlerQueueURL(ctx, handler)
	if err != nil {
		return fmt.Errorf("handler issue: %w", err)
	}

	for _, doc := range docs {
		err = Listen(ctx, detector, event.ListenOnChange{
			Handler: handler,
			Documents: []event.ListenOnChangeDocument{
				{
					Type:        doc.Type,
					Name:        doc.Name,
					NewDocument: true,
				},
			},
		})
		if err != nil {
			return fmt.Errorf("can't add listener to the doc: %w", err)
		}
	}

	return nil
}

func assignDoc(ctx context.Context, registry jsonball.Registry, detector propchange.Detector, handler, docType, docName string) error {
	return assignDocs(ctx, registry, detector, handler, []event.GetDocument{{Type: docType, Name: docName}})
}
