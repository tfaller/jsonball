package operation

import (
	"context"
	"time"

	"github.com/tfaller/jsonball/internal/name"
	"github.com/tfaller/propchange"
)

// RequeueHandler requeue all listeners of the given handler. This is used if the handler code changed
// and the handler needs to process all data again.
func RequeueHandler(ctx context.Context, detector propchange.Detector, handler string) error {
	// To do this each listener listens for a requeue property.
	// We just have to change the property. We simply set the version to the current minutes since 1970 + 1.
	// This guarantees that all previous defined listener will trigger this.
	doc, err := detector.OpenDocument(ctx, name.CreateDocName(name.HandlerDocumentType, handler))
	if err != nil {
		return err
	}
	defer doc.Close()

	err = doc.SetProperty(name.HandlerRequeueProperty, uint64(time.Now().Unix()/60)+1)
	if err != nil {
		return err
	}

	return doc.Commit()
}
