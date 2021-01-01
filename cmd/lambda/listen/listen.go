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
	detector = startup.MustGetDetector()
)

func lambdaMain(ctx context.Context, listen event.ListenOnChange) error {
	return operation.Listen(ctx, detector, listen)
}

func main() {
	lambda.Start(lambdaMain)
}
