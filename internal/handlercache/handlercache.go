package handlercache

import (
	"context"
	"sync"

	"github.com/tfaller/jsonball"
)

// HandlerCache caches handler data.
// If the wanted data is not already cached it will be loaded
type HandlerCache struct {
	m               sync.Mutex
	registry        jsonball.Registry
	handlerQueueURL map[string]string
}

// NewHandlerCache creates a new handler cache.
func NewHandlerCache(registry jsonball.Registry) *HandlerCache {
	return &HandlerCache{
		registry:        registry,
		handlerQueueURL: map[string]string{},
	}
}

// GetHandlerQueueURL gets the queue URL of a given handler.
// If the data is not already cached, the data will be loaded.
func (h *HandlerCache) GetHandlerQueueURL(ctx context.Context, handler string) (queueURL string, err error) {
	h.m.Lock()
	defer h.m.Unlock()

	queueURL = h.handlerQueueURL[handler]
	if queueURL == "" {
		queueURL, err = h.registry.GetHandlerQueueURL(ctx, handler)
		if err == nil {
			// cache the resolved URL
			h.handlerQueueURL[handler] = queueURL
		}
	}

	return
}
