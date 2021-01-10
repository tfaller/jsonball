package startup

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

var (
	awsConfig    aws.Config
	awsSqsClient *sqs.Client
)

func getAwsConfig(ctx context.Context) (aws.Config, error) {
	if awsConfig.Region != "" {
		return awsConfig, nil
	}

	var err error
	awsConfig, err = config.LoadDefaultConfig(ctx)
	if err != nil {
		return awsConfig, fmt.Errorf("can't load aws config: %w", err)
	}
	return awsConfig, nil
}

// GetSqsClient gets an sqs client
func GetSqsClient(ctx context.Context) (*sqs.Client, error) {
	if awsSqsClient != nil {
		return awsSqsClient, nil
	}

	awsConfig, err := getAwsConfig(ctx)
	if err != nil {
		return nil, err
	}

	awsSqsClient = sqs.NewFromConfig(awsConfig)
	return awsSqsClient, nil
}

// MustGetSqsClient like GetSqsClient but panics
// if an error happened
func MustGetSqsClient() *sqs.Client {
	client, err := GetSqsClient(context.Background())
	if err != nil {
		panic(err)
	}
	return client
}
