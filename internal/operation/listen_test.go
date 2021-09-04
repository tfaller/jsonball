package operation

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/tfaller/jsonball/event"
	"github.com/tfaller/jsonball/internal/mock/mock_propchange"
)

func TestListenHandlerName(t *testing.T) {
	ctx := context.TODO()
	ctrl := gomock.NewController(t)

	detector := mock_propchange.NewMockDetector(ctrl)

	detector.EXPECT().AddListener(ctx, "handlerName:listenerName", gomock.Any())

	err := Listen(ctx, detector, event.ListenOnChange{
		Name: "listenerName", Handler: "handlerName",
		Documents: []event.ListenOnChangeDocument{{NewDocument: true, Type: "test", Name: "a"}}},
	)

	assert.NoError(t, err)
}
