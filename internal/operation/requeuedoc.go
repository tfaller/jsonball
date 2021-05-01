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

func RequeueDoc(ctx context.Context, registry jsonball.Registry, detector propchange.Detector, docType, docName string) error {
	openDoc, err := registry.OpenDocument(ctx, docType, docName)
	if err != nil {
		return fmt.Errorf("can't open document: %w", err)
	}
	defer openDoc.Close()

	if openDoc.IsNew() {
		return fmt.Errorf("document does not exist")
	}

	doc := &jsonera.EraDocument{}
	err = json.Unmarshal([]byte(openDoc.Document()), &doc)
	if err != nil {
		return fmt.Errorf("can't decode document: %w", err)
	}

	doc.DocEra++
	era := uint64(doc.DocEra)

	detectDoc, err := detector.OpenDocument(ctx, name.CreateDocName(docType, docName))
	if err != nil {
		return err
	}
	defer detectDoc.Close()

	// update special document prop
	if err = detectDoc.SetProperty("d", era); err != nil {
		return fmt.Errorf("can't update document version: %w", err)
	}

	// update all regular properties
	jsonvisitor.Visit(doc.Doc, func(path []string, value interface{}) bool {
		if err != nil {
			return false
		}
		err = detectDoc.SetProperty("p"+jsonp.Format(path), era)
		return true
	})

	if err != nil {
		return fmt.Errorf("can't update property version: %w", err)
	}

	// apply changes
	if err = detectDoc.Commit(); err != nil {
		return fmt.Errorf("can't save property doc: %w", err)
	}

	newDoc, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("can't marshal updated doc: %w", err)
	}

	return openDoc.Change(jsonball.Change{
		Document: string(newDoc),
	})
}
