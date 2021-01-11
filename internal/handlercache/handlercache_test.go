package handlercache

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/tfaller/jsonball/internal/mock/mock_jsonball"
)

func TestHandlerCache(t *testing.T) {
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	registry := mock_jsonball.NewMockRegistry(ctrl)
	registry.EXPECT().
		GetHandlerQueueURL(ctx, "test").
		Times(1).
		Return("url", nil)

	cache := NewHandlerCache(registry)

	for i := 0; i < 2; i++ {
		url, err := cache.GetHandlerQueueURL(ctx, "test")

		if err != nil {
			t.Errorf("Expected not error but got %v", err)
		}

		if url != "url" {
			t.Errorf("Expected url 'url' but got %v", url)
		}
	}
}
