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
	detector = startup.MustGetDetector()
)

func lambdaMain(ctx context.Context, adminCmd *event.AdminCmd) error {
	return operation.AdminCommands(ctx, registry, detector, adminCmd)
}

func main() {
	lambda.Start(lambdaMain)
}
