package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/tfaller/jsonball/internal/handlercache"
	"github.com/tfaller/jsonball/internal/operation"
	"github.com/tfaller/jsonball/internal/startup"
	"github.com/tfaller/propchange"
	"golang.org/x/sync/errgroup"
)

var (
	sqsClient = startup.MustGetSqsClient()
	detector  = startup.MustGetDetector()
	registry  = startup.MustGetRegistry()

	messageID       = "jsonball"
	handlerQueueURL = handlercache.NewHandlerCache(registry)
)

func lambdaMain(ctx context.Context) error {
	throttle := make(chan struct{}, 10)
	errGroup, errCtx := errgroup.WithContext(ctx)

loop:
	for {
		select {
		case throttle <- struct{}{}:
		case <-errCtx.Done():
			break loop
		}

		// check again whether no error exists
		// it could be that only throttle was checked
		if errCtx.Err() != nil {
			break
		}

		change, err := detector.NextChange(ctx)
		if err != nil {
			if !errors.Is(err, propchange.ErrNoMoreChanges) {
				errGroup.Go(func() error { return err })
			}
			// stop because of no more changes or error
			break
		}

		// process changes concurrently
		errGroup.Go(func() error {
			err := processChange(ctx, change)
			if err != nil {
				// because we don't return all errors, log
				// concurent errors
				log.Println(err)
				change.Close()
			}

			<-throttle
			return err
		})
	}

	return errGroup.Wait()
}

func processChange(ctx context.Context, change propchange.OnChange) error {
	jsonball, err := operation.HandleChange(ctx, registry, change)
	if err != nil {
		return fmt.Errorf("can't build change event: %w", err)
	}

	msgBody, err := json.Marshal(jsonball)
	if err != nil {
		return fmt.Errorf("can't marshal change event: %w", err)
	}

	queueURL, err := handlerQueueURL.GetHandlerQueueURL(ctx, jsonball.Handler)
	if err != nil {
		return fmt.Errorf("can't resolve event handler queue: %w", err)
	}

	if !strings.HasSuffix(queueURL, ".fifo") {
		// should always be a fifo queue
		return fmt.Errorf("queue %q is not a fifo queue", queueURL)
	}

	_, err = sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    &queueURL,
		MessageBody: aws.String(string(msgBody)),
		// For now we don't have a safe way to use different message ids.
		// By using the same for all, messages are processed in a strict FIFO manner.
		MessageGroupId:         &messageID,
		MessageDeduplicationId: aws.String(uuid.New().String()),
	})
	if err != nil {
		return fmt.Errorf("can't send message to handler queue: %w", err)
	}

	if err = change.Commit(); err != nil {
		return fmt.Errorf("can't commit change as handled: %w", err)
	}

	return nil
}

func main() {
	lambda.Start(lambdaMain)
}
