package operation

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/tfaller/jsonball/event"
	"github.com/tfaller/jsonball/internal/mock/mock_jsonball"
	"github.com/tfaller/jsonball/internal/mock/mock_propchange"
	"github.com/tfaller/jsonball/internal/name"
)

func TestRegisterHandler(t *testing.T) {
	ctx := context.TODO()
	ctrl := gomock.NewController(t)

	registry := mock_jsonball.NewMockRegistry(ctrl)
	detector := mock_propchange.NewMockDetector(ctrl)

	testCases := []struct {
		Handler *event.RegisterHandler
		Error   error
	}{
		// simple test case
		{&event.RegisterHandler{Handler: "handler", QueueURL: "queueurl"}, nil},
		// too long name
		{&event.RegisterHandler{Handler: "veryveryverylongname", QueueURL: "queueurl"}, name.ErrHandlerNameInvalid},
	}

	for idx, test := range testCases {
		registry.EXPECT().RegisterHandler(ctx, test.Handler)

		// requeue property must be set
		docOps := mock_propchange.NewMockDocumentOps(ctrl)
		docOps.EXPECT().SetProperty(name.HandlerRequeueProperty, gomock.Any())
		docOps.EXPECT().Commit()
		docOps.EXPECT().Close()

		detector.EXPECT().OpenDocument(ctx, name.CreateDocName(name.HandlerDocumentType, test.Handler.Handler)).Return(
			docOps, nil,
		)

		err := RegisterHandler(ctx, registry, detector, test.Handler)

		if !errors.Is(err, test.Error) {
			t.Errorf("%v: Expected error %v but got %v", idx, test.Error, err)
		}
	}
}
