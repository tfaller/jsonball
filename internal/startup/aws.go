package startup

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

var (
	awsConfig    aws.Config
	awsSqsClient *sqs.Client
)

func getAwsConfig() (aws.Config, error) {
	if awsConfig.Region != "" {
		return awsConfig, nil
	}

	var err error
	awsConfig, err = config.LoadDefaultConfig()
	if err != nil {
		return awsConfig, fmt.Errorf("can't load aws config: %w", err)
	}
	return awsConfig, nil
}

// GetSqsClient gets an sqs client
func GetSqsClient() (*sqs.Client, error) {
	if awsSqsClient != nil {
		return awsSqsClient, nil
	}

	awsConfig, err := getAwsConfig()
	if err != nil {
		return nil, err
	}

	awsSqsClient = sqs.NewFromConfig(awsConfig)
	return awsSqsClient, nil
}

// MustGetSqsClient like GetSqsClient but panics
// if an error happened
func MustGetSqsClient() *sqs.Client {
	client, err := GetSqsClient()
	if err != nil {
		panic(err)
	}
	return client
}
