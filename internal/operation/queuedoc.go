package operation

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/tfaller/jsonball/event"
	"github.com/tfaller/jsonball/internal/name"
)

// ErrQueueIsNotFIFO indicates that the target queue is not a FIFO queue
var ErrQueueIsNotFIFO = errors.New("target queue should be a FIFO queue")

// QueueDocument queues a document which needs to be proccessed
func QueueDocument(ctx context.Context, sqsClient *sqs.Client, queueURL string, doc event.PostDocument) error {
	if !strings.HasPrefix(queueURL, ".fifo") {
		return ErrQueueIsNotFIFO
	}

	msg, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	msgStr := string(msg)
	// msgID is used so that all msg of the same document FIFO Queue in the right order
	// without blocking other unrelated documents
	msgID := name.CreateDocName(doc.Type, doc.Name)

	_, err = sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody:    &msgStr,
		MessageGroupId: &msgID,
		QueueUrl:       &queueURL,
	})

	return err
}
