package name

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// CreateListenerName creates a unique listener name for a given handler.
func CreateListenerName(handler, name string) string {
	if name == "" {
		name = uuid.NewString()
	}
	return fmt.Sprintf("%v:%v", handler, name)
}

// ParseListenerName returns the handler of a listener.
func ParseListenerName(listener string) (string, string) {
	parts := strings.SplitN(listener, ":", 2)
	if len(parts) < 2 {
		return "", ""
	}
	return parts[0], parts[1]
}
