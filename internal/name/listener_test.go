package name

import "testing"

func TestListener(t *testing.T) {
	a1 := CreateListenerName("handlerA")
	a2 := CreateListenerName("handlerA")

	if a1 == a2 {
		t.Fatalf("handler names must be unique")
	}

	if h1 := ParseListenerName(a1); h1 != "handlerA" {
		t.Fatalf("%v != handlerA", h1)
	}

	if h2 := ParseListenerName(a2); h2 != "handlerA" {
		t.Fatalf("%v != handlerA", h2)
	}

	if ParseListenerName("abc123") != "" {
		t.Fatalf("Parsed an invalid listener")
	}
}
