package name

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestListener(t *testing.T) {
	assert := assert.New(t)

	a1 := CreateListenerName("handlerA", "")
	a2 := CreateListenerName("handlerA", "")

	// name should be auto-generated -> no simlar names
	assert.NotEqual(a1, a2)

	h1, n1 := ParseListenerName(a1)
	assert.Equal("handlerA", h1)
	assert.NotPanics(func() { uuid.MustParse(n1) })

	h2, n2 := ParseListenerName(a2)
	assert.Equal("handlerA", h2)
	assert.NotPanics(func() { uuid.MustParse(n2) })

	// invalid name
	h, _ := ParseListenerName("abc123")
	assert.Empty(h)

	b1 := CreateListenerName("handlerB", "bName")
	b2 := CreateListenerName("handlerB", "bName")

	assert.Equal(b1, b2)

	h3, n3 := ParseListenerName(b1)
	assert.Equal("handlerB", h3)
	assert.Equal("bName", n3)
}
