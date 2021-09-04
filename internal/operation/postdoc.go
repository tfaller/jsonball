package operation

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tfaller/go-jsonera"
	"github.com/tfaller/go-jsonera/pkg/jsonp"
	"github.com/tfaller/go-jsonvisitor"
	"github.com/tfaller/jsonball"
	"github.com/tfaller/jsonball/internal/name"
	"github.com/tfaller/propchange"
)

// PutDocument identifies changes of document and updates tha properties revision to later find changes
func PutDocument(ctx context.Context, docRegistry jsonball.Registry, detector propchange.Detector, docType, docName string, doc interface{}) error {
	registryDoc, err := docRegistry.OpenDocument(ctx, docType, docName)
	if err != nil {
		return err
	}

	var eraDoc *jsonera.EraDocument
	if !registryDoc.IsNew() {
		// update existing document
		if err = json.Unmarshal([]byte(registryDoc.Document()), &eraDoc); err != nil {
			return fmt.Errorf("Parsing of stored era document failed: %w", err)
		}
	} else {
		// new document
		eraDoc = jsonera.NewEraDocument(jsonvisitor.Undefined)
	}

	changes := eraDoc.UpdateDoc(doc)

	if len(changes) == 0 {
		// we found no changes.
		// only refresh the document.
		// no more operations are neeede.
		return registryDoc.Refresh()
	}

	// open propchange document to apply changes
	propDocName := name.CreateDocName(docType, docName)
	propDoc, err := detector.OpenDocument(ctx, propDocName)
	if err != nil {
		return err
	}

	for _, c := range changes {
		name := "p" + jsonp.Format(c.Path)
		if c.Mode == jsonera.ChangeDelete {
			err = propDoc.DelProperty(name)
		} else {
			err = propDoc.SetProperty(name, uint64(eraDoc.DocEra))
		}
		if err != nil {
			return err
		}
	}

	// update special whole document property
	err = propDoc.SetProperty("d", uint64(eraDoc.DocEra))
	if err != nil {
		return err
	}

	err = propDoc.Commit()
	if err != nil {
		return err
	}

	if registryDoc.IsNew() {
		handlers, err := docRegistry.GetNewDocHanders(ctx, docType)
		if err != nil {
			return err
		}

		// add "new doc" listerens to the new docs
		for _, handler := range handlers {
			err = detector.AddListener(ctx, name.CreateListenerName(handler, ""), []propchange.ChangeFilter{
				{
					Document:   propDocName,
					Properties: map[string]uint64{"d": 0},
				},
			})
			if err != nil {
				return fmt.Errorf("can't register listener: %w", err)
			}
		}
	}

	eraDocNew, err := json.Marshal(eraDoc)
	if err != nil {
		return err
	}

	return registryDoc.Change(jsonball.Change{
		Document: string(eraDocNew),
	})
}
