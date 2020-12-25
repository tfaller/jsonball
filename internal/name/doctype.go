package name

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	docTypePatternRaw = "^[a-z0-9\\-]{1,20}$"

	// InternalPrefix is the prefix that each internal
	// document type must have.
	InternalPrefix = "jb"
)

var (
	docTypePattern = regexp.MustCompilePOSIX(docTypePatternRaw)

	// ErrDocTypeNameReserved indicates that a reserved name was used
	// as a document type name.
	ErrDocTypeNameReserved = fmt.Errorf("document type which starts with %q is reserved", InternalPrefix)

	// ErrDocTypeNameInternal indicates that a document type that is
	ErrDocTypeNameInternal = fmt.Errorf("internal document type must has %q as a prefix", InternalPrefix)

	// ErrDocTypeNameInvalid indicates that the document type name is not valid
	// (e.g. contains invalid characters).
	ErrDocTypeNameInvalid = fmt.Errorf("document type must match pattern %q", docTypePatternRaw)
)

// CheckDocTypeName checks if a docType name is valid
func CheckDocTypeName(docType string, internal bool) error {

	// prefix checks
	if internal && !strings.HasPrefix(docType, InternalPrefix) {
		return ErrDocTypeNameInternal
	}
	if !internal && strings.HasPrefix(docType, InternalPrefix) {
		return ErrDocTypeNameReserved
	}

	// general check
	if !docTypePattern.MatchString(docType) {
		return ErrDocTypeNameInvalid
	}

	return nil
}
