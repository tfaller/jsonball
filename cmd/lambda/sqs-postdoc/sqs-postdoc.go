package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/tfaller/jsonball/event"
	"github.com/tfaller/jsonball/internal/operation"
	"github.com/tfaller/jsonball/internal/startup"
)

var (
	sqsClient   = startup.MustGetSqsClient()
	docRegistry = startup.MustGetRegistry()
	detector    = startup.MustGetDetector()
	queueURL    = os.Getenv("QUEUE_URL")
)

func lambdaMain(ctx context.Context, sqsEvent events.SQSEvent) error {

	// one record after the other ... because FIFO queue
	for _, rec := range sqsEvent.Records {
		err := handleRecord(ctx, rec)
		if err != nil {
			// return instantly with an error
			// SQS will requeue the messages
			return err
		}

		// remove message from queue ... lambda would to this for us,
		// but if a following message causes an error, we would handle this
		// message as well again
		if queueURL != "" {
			_, err := sqsClient.DeleteMessage(ctx,
				&sqs.DeleteMessageInput{ReceiptHandle: &rec.ReceiptHandle, QueueUrl: &queueURL})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func handleRecord(ctx context.Context, rec events.SQSMessage) error {
	var putDoc event.PostDocument
	err := json.Unmarshal([]byte(rec.Body), &putDoc)
	if err != nil {
		return err
	}
	return operation.PutDocument(ctx, docRegistry, detector, putDoc.Type, putDoc.Name, putDoc.Document)
}

func main() {
	lambda.Start(lambdaMain)
}
