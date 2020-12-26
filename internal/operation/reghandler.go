package operation

import (
	"context"

	"github.com/tfaller/jsonball"
	"github.com/tfaller/jsonball/event"
	"github.com/tfaller/jsonball/internal/name"
	"github.com/tfaller/propchange"
)

// RegisterHandler registers a new handler
func RegisterHandler(ctx context.Context, registry jsonball.Registry, detector propchange.Detector, handler *event.RegisterHandler) error {
	err := name.CheckHandlerName(handler.Handler)
	if err != nil {
		return err
	}

	// create special internal handler document
	doc, err := detector.OpenDocument(ctx, name.CreateDocName(name.HandlerDocumentType, handler.Handler))
	if err != nil {
		return err
	}
	defer doc.Close()

	// set special requeue property. All listeners for this handler
	// will listen for this property. This property will be only updated
	// if the handler should requeue all listened documents.
	err = doc.SetProperty(name.HandlerRequeueProperty, 1)
	if err != nil {
		return err
	}

	err = doc.Commit()
	if err != nil {
		return err
	}

	// register the handler itself
	return registry.RegisterHandler(ctx, handler)
}
