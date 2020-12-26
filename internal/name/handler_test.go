package name

import (
	"errors"
	"testing"
)

func TestCheckHandlerName(t *testing.T) {
	testCases := []struct {
		Name  string
		Error error
	}{
		// simple name
		{"test", nil},
		// invalid name
		{":", ErrHandlerNameInvalid},
	}

	for idx, test := range testCases {
		err := CheckHandlerName(test.Name)
		if !errors.Is(err, test.Error) {
			t.Errorf("%v: Expected error %v but got %v", idx, test.Error, err)
		}
	}
}
