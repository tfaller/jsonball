package name

import (
	"fmt"
	"regexp"
)

const (
	handlerPatternRaw = "^[a-z0-9\\-]{1,30}$"

	// HandlerDocumentType is a reserved internal document type
	// use for special handler related documents
	HandlerDocumentType = InternalPrefix + "-handler"

	// HandlerRequeueProperty is used to signal that all listeners
	// of the given handler should trigger a change. This is used if
	// the handler code was changed and now all data must be processed again.
	HandlerRequeueProperty = "requeue"
)

var (
	handlerPattern = regexp.MustCompilePOSIX(handlerPatternRaw)

	// ErrHandlerNameInvalid indicates that the handler name is not valid
	// (e.g. contains invalid characters).
	ErrHandlerNameInvalid = fmt.Errorf("handler must match pattern %q", handlerPatternRaw)
)

// CheckHandlerName checks if the handler name is valid
func CheckHandlerName(handler string) error {
	if !handlerPattern.MatchString(handler) {
		return ErrHandlerNameInvalid
	}

	return nil
}
