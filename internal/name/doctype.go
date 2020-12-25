package name

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const docTypePatternRaw = "^[a-z0-9\\-]{1,20}$"

var (
	docTypePattern = regexp.MustCompilePOSIX(docTypePatternRaw)

	// ErrDocTypeNameReserved indicates that a reserved name was used
	// as a document type name.
	ErrDocTypeNameReserved = errors.New("document type which starts with \"jsonball\" is reserved")

	// ErrDocTypeNameInvalid indicates that the document type name is not valid
	// (e.g. contains invalid characters).
	ErrDocTypeNameInvalid = fmt.Errorf("document type must match pattern %q", docTypePatternRaw)
)

// CheckDocTypeName checks if a docType name is valid
func CheckDocTypeName(docType string) error {
	if strings.HasPrefix(docType, "jsonball") {
		return ErrDocTypeNameReserved
	}

	if !docTypePattern.MatchString(docType) {
		return ErrDocTypeNameInvalid
	}

	return nil
}
