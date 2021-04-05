package operation

import (
	"context"
	"fmt"

	"github.com/tfaller/jsonball"
	"github.com/tfaller/jsonball/event"
	"github.com/tfaller/propchange"
)

func HandlerNewDoc(ctx context.Context, registry jsonball.Registry, detector propchange.Detector, handerNewDoc *event.HandlerNewDoc) error {

	err := registry.HandlerNewDoc(ctx, handerNewDoc.Handler, handerNewDoc.Type, true)
	if err != nil {
		return fmt.Errorf("can't register handler for new doc: %w", err)
	}

	if handerNewDoc.Existing {
		// We must trigger for already existing docs.
		// To do this queue these documents as a change.
		nextToken := ""

		for {
			docList, err := registry.ListDocuments(ctx, handerNewDoc.Type, nextToken, 100)
			if err != nil {
				return fmt.Errorf("can't list existing docs: %w", err)
			}

			if len(docList.Documents) == 0 {
				// we iterated over all existing documents
				break
			}

			for _, doc := range docList.Documents {
				err = Listen(ctx, detector, event.ListenOnChange{
					Handler: handerNewDoc.Handler,
					Documents: []event.ListenOnChangeDocument{{
						Type:        handerNewDoc.Type,
						Name:        doc,
						NewDocument: true,
					}}})

				if err != nil {
					return fmt.Errorf("can't add queue-listener for existing doc: %w", err)
				}
			}

			nextToken = docList.NextToken
		}

	}

	return nil
}
