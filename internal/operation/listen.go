package operation

import (
	"context"
	"fmt"
	"strings"

	"github.com/tfaller/jsonball/event"
	"github.com/tfaller/jsonball/internal/name"
	"github.com/tfaller/propchange"
)

// Listen adds a listener that listens to changes of properties
func Listen(ctx context.Context, detector propchange.Detector, listen event.ListenOnChange) error {
	if listen.Handler == "" {
		return fmt.Errorf("Handler name can't be empty")
	}

	docFilter := make([]propchange.ChangeFilter, len(listen.Documents))

	for i, d := range listen.Documents {
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
			Document:   name.CreateDocName(d.Type, d.Name),
			Properties: prop,
		}
	}

	return detector.AddListener(ctx, name.CreateListenerName(listen.Handler), docFilter)
}
