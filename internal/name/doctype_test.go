package name

import (
	"errors"
	"testing"
)

func TestCheckDocTypeName(t *testing.T) {
	testCases := []struct {
		DocType string
		Error   error
	}{
		// simple basic name
		{"test", nil},
		// empty name
		{"", ErrDocTypeNameInvalid},
		// too long name
		{"1234567889012345678901", ErrDocTypeNameInvalid},
		// reserved name
		{"jsonball", ErrDocTypeNameReserved},
		// reserved name prefix
		{"jsonball2", ErrDocTypeNameReserved},
	}

	for idx, test := range testCases {
		err := CheckDocTypeName(test.DocType)
		if !errors.Is(err, test.Error) {
			t.Errorf("%v: Expected error %v but got %v", idx, test.Error, err)
		}
	}
}
