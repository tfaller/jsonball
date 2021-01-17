package main

import (
	"context"
	"errors"

	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tfaller/jsonball"
	"github.com/tfaller/jsonball/event"
	"github.com/tfaller/jsonball/internal/operation"
	"github.com/tfaller/jsonball/internal/startup"
)

var (
	registry = startup.MustGetRegistry()
)

func lambdaMain(ctx context.Context, e event.GetDocument) (doc event.Document, err error) {
	doc, err = operation.GetDocumentContent(ctx, registry, e.Type, e.Name)
	if errors.Is(err, jsonball.ErrDocumentNotExist) {
		// return empty document
		return event.Document{Type: e.Type, Name: e.Name}, nil
	}
	return
}

func main() {
	lambda.Start(lambdaMain)
}
