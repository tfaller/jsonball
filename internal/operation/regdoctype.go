package operation

import (
	"context"

	"github.com/tfaller/jsonball"
	"github.com/tfaller/jsonball/internal/name"
)

// RegisterDocumentType registers a new document type
func RegisterDocumentType(ctx context.Context, registry jsonball.Registry, docType string) error {
	return registerDocumentType(ctx, registry, docType, false)
}

// RegisterDocumentTypeInternal registers a new internal document type.
// This function should only be used to register internal (reserved) document types.
// Don't use this for regular document types.
func RegisterDocumentTypeInternal(ctx context.Context, registry jsonball.Registry, docType string) error {
	return registerDocumentType(ctx, registry, docType, true)
}

func registerDocumentType(ctx context.Context, registry jsonball.Registry, docType string, internal bool) error {
	err := name.CheckDocTypeName(docType, internal)
	if err != nil {
		return err
	}
	return registry.RegisterDocumentType(ctx, docType)
}
