package name

import (
	"errors"
	"fmt"
	"testing"
)

func TestCheckDocTypeName(t *testing.T) {
	testCases := []struct {
		DocType  string
		Error    error
		Internal bool
	}{
		// simple basic name
		{"test", nil, false},
		// empty name
		{"", ErrDocTypeNameInvalid, false},
		// too long name
		{"1234567889012345678901", ErrDocTypeNameInvalid, false},
		// reserved name
		{InternalPrefix, ErrDocTypeNameReserved, false},
		// reserved name prefix
		{fmt.Sprintf("%v%v", InternalPrefix, "2"), ErrDocTypeNameReserved, false},
		// simple internal type
		{fmt.Sprintf("%v-test", InternalPrefix), nil, true},
		// wrong internal type
		{"test", ErrDocTypeNameInternal, true},
	}

	for idx, test := range testCases {
		err := CheckDocTypeName(test.DocType, test.Internal)
		if !errors.Is(err, test.Error) {
			t.Errorf("%v: Expected error %v but got %v", idx, test.Error, err)
		}
	}
}

func TestIsDocTypeInternal(t *testing.T) {
	if !IsDocTypeInternal(InternalPrefix + "test") {
		t.Error("Expected this to be internal")
	}
	if IsDocTypeInternal("test") {
		t.Error("Expected this be internal")
	}
}
