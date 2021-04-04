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

	errCommit := errors.New("commit failed")
	errSetProp := errors.New("set property failed")
	errOpenDoc := errors.New("open doc failed")

	testCases := []struct {
		Handler    *event.RegisterHandler
		Error      error
		ErrCommit  error
		ErrSetProp error
		ErrOpenDoc error
	}{
		// simple test case
		{&event.RegisterHandler{Handler: "handler", QueueURL: "queueurl"}, nil, nil, nil, nil},
		// too long name
		{&event.RegisterHandler{Handler: "veryveryverylongnameveryveryverylongname", QueueURL: "queueurl"}, name.ErrHandlerNameInvalid, nil, nil, nil},
		// commit failed
		{&event.RegisterHandler{Handler: "handler", QueueURL: "queueurl"}, errCommit, errCommit, nil, nil},
		// set requeue property failed
		{&event.RegisterHandler{Handler: "handler", QueueURL: "queueurl"}, errSetProp, nil, errSetProp, nil},
		// open document failed
		{&event.RegisterHandler{Handler: "handler", QueueURL: "queueurl"}, errOpenDoc, nil, nil, errOpenDoc},
	}

	for idx, test := range testCases {
		registry.EXPECT().RegisterHandler(ctx, test.Handler)

		// requeue property must be set
		docOps := mock_propchange.NewMockDocumentOps(ctrl)
		docOps.EXPECT().SetProperty(name.HandlerRequeueProperty, gomock.Any()).Return(test.ErrSetProp)
		docOps.EXPECT().Commit().Return(test.ErrCommit)
		docOps.EXPECT().Close()

		detector.EXPECT().OpenDocument(ctx, name.CreateDocName(name.HandlerDocumentType, test.Handler.Handler)).Return(
			docOps, test.ErrOpenDoc,
		)

		err := RegisterHandler(ctx, registry, detector, test.Handler)

		if !errors.Is(err, test.Error) {
			t.Errorf("%v: Expected error %v but got %v", idx, test.Error, err)
		}
	}
}
