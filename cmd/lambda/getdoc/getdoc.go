package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tfaller/jsonball/event"
	"github.com/tfaller/jsonball/internal/operation"
	"github.com/tfaller/jsonball/internal/startup"
)

var (
	registry = startup.MustGetRegistry()
)

func lambdaMain(ctx context.Context, e event.GetDocument) (event.Document, error) {
	return operation.GetDocumentContent(ctx, registry, e.Type, e.Name)
}

func main() {
	lambda.Start(lambdaMain)
}
