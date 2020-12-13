package name

import (
	"fmt"
	"strconv"
	"strings"
)

// CreateDocName creates a document name based of docType and docName
func CreateDocName(docType, name string) string {
	return fmt.Sprintf("%v:%v:%v", len(docType), docType, name)
}

// ParseDocName parses a DocName into docType and docName
func ParseDocName(name string) (docType, docName string, err error) {
	parts := strings.SplitN(name, ":", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("Invalid docName parts")
	}

	typeLen, err := strconv.ParseUint(parts[0], 10, 32)
	if err != nil {
		return "", "", fmt.Errorf("Invalid type name length: %w", err)
	}

	if typeLen > uint64(len(parts[1])) {
		return "", "", fmt.Errorf("TypeLen is too large")
	}

	if parts[1][typeLen] != ':' {
		return "", "", fmt.Errorf("TypeLen does not point to seperator")
	}

	return parts[1][:typeLen], parts[1][typeLen+1:], nil
}
