package main

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tfaller/jsonball/event"
	"github.com/tfaller/jsonball/internal/operation"
	"github.com/tfaller/jsonball/internal/startup"
)

var (
	detector = startup.MustGetDetector()
)

func lambdaMain(ctx context.Context, sqsEvent events.SQSEvent) error {
	if len(sqsEvent.Records) == 0 {
		return errors.New("invalid sqs event with no records")
	}
	for _, m := range sqsEvent.Records {
		err := handleMessage(ctx, m)
		if err != nil {
			return err
		}
	}
	return nil
}

func handleMessage(ctx context.Context, m events.SQSMessage) error {
	var listen event.ListenOnChange

	err := json.Unmarshal([]byte(m.Body), &listen)
	if err != nil {
		return err
	}

	return operation.Listen(ctx, detector, listen)
}

func main() {
	lambda.Start(lambdaMain)
}
