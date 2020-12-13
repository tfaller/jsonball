package name

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// CreateListenerName creates a unique listener name for a given handler.
func CreateListenerName(handler string) string {
	return fmt.Sprintf("%v:%v", handler, uuid.New().String())
}

// ParseListenerName returns the handler of a listener.
func ParseListenerName(listener string) string {
	parts := strings.SplitN(listener, ":", 2)
	if len(parts) < 2 {
		return ""
	}
	return parts[0]
}
