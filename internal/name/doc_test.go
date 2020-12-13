package name

import "testing"

func TestDoc(t *testing.T) {
	name := CreateDocName("typeA", "A1234567")
	docType, docName, err := ParseDocName(name)

	if err != nil {
		t.Errorf("Expected no error but got: %v", err)
	}

	if docType != "typeA" {
		t.Errorf("Expected typeA but got %v", docType)
	}

	if docName != "A1234567" {
		t.Errorf("Expected A1234567 but got %v", docName)
	}
}
