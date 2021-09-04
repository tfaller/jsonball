package operation

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/tfaller/jsonball/event"
	"github.com/tfaller/jsonball/internal/name"
	"github.com/tfaller/propchange"
)

// Listen adds a listener that listens to changes of properties
func Listen(ctx context.Context, detector propchange.Detector, listen event.ListenOnChange) error {
	if listen.Handler == "" {
		return fmt.Errorf("Handler name can't be empty")
	}

	docFilter := make([]propchange.ChangeFilter, len(listen.Documents)+1)

	for i, d := range listen.Documents {
		// make sure that the document type is valid
		if err := name.CheckDocTypeName(d.Type, false); err != nil {
			return err
		}

		prop := map[string]uint64{}
		for _, p := range d.Properties {
			if p == "" {
				// json path of the whole document.
				// use special whole document property
				p = "p"
			} else {
				if !strings.HasPrefix(p, "/") {
					p = "/" + p
				}
				p = "p" + p
			}
			prop[p] = uint64(d.Version)
		}
		docFilter[i] = propchange.ChangeFilter{
			Document:    name.CreateDocName(d.Type, d.Name),
			NewDocument: d.NewDocument,
			Properties:  prop,
		}
	}

	// add special handler requeue property
	docFilter[len(docFilter)-1] = propchange.ChangeFilter{
		Document:   name.CreateDocName(name.HandlerDocumentType, listen.Handler),
		Properties: map[string]uint64{name.HandlerRequeueProperty: uint64(time.Now().Unix() / 60)},
	}

	return detector.AddListener(ctx, name.CreateListenerName(listen.Handler, listen.Name), docFilter)
}
