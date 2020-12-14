package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tfaller/jsonball/internal/operation"
	"github.com/tfaller/jsonball/internal/startup"
	"github.com/tfaller/propchange"
)

var (
	sqsClient = startup.MustGetSqsClient()
	detector  = startup.MustGetDetector()
	registry  = startup.MustGetRegistry()
)

func lambdaMain(ctx context.Context) error {
	for {
		change, err := detector.NextChange(ctx)
		if err != nil {
			if errors.Is(err, propchange.ErrNoMoreChanges) {
				// end early ... no more changes
				return nil
			}
			return err
		}

		err = processChange(ctx, change)
		if err != nil {
			change.Close()
			return err
		}
	}
}

func processChange(ctx context.Context, change propchange.OnChange) error {
	jsonball, err := operation.HandleChange(ctx, registry, change)
	if err != nil {
		return err
	}

	msgBody, err := json.Marshal(jsonball)
	if err != nil {
		return err
	}

	queueURL, err := registry.GetHandlerQueueURL(ctx, jsonball.Handler)
	if err != nil {
		return err
	}

	if !strings.HasSuffix(queueURL, ".fifo") {
		// should always be a fifo queue
		return fmt.Errorf("queue %q is not a fifo queue", queueURL)
	}

	_, err = sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    &queueURL,
		MessageBody: aws.String(string(msgBody)),
	})
	if err != nil {
		return err
	}

	return change.Commit()
}

func main() {
	lambda.Start(lambdaMain)
}
