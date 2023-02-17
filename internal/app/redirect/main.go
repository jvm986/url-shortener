package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/jvm986/url-shortener/internal/pkg/storage"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	log, err := zap.NewProduction()
	if err != nil {
		panic("failed to init logger")
	}

	storageTable, ok := os.LookupEnv("STORAGE_TABLE_NAME")
	if !ok {
		log.Sugar().Fatal("failed to lookup STORAGE_TABLE_NAME from environment variables")
	}

	env, ok := os.LookupEnv("ENV")
	if !ok {
		log.Sugar().Fatal("failed to lookup ENV from environment variables")
	}
	dynamodbEndpoint, ok := os.LookupEnv("DDB_ENDPOINT")
	if !ok {
		log.Sugar().Fatal("failed to lookup DDB_ENDPOINT from environment variables")
	}

	awsRegion, ok := os.LookupEnv("REGION")
	if !ok {
		log.Sugar().Fatal("failed to lookup REGION from environment variables")
	}

	dynamodbEndpointResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == dynamodb.ServiceID && env == "development" {
			return aws.Endpoint{
				URL:           dynamodbEndpoint,
				SigningRegion: awsRegion,
			}, nil
		}
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	awsConfig, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(awsRegion),
		config.WithEndpointResolverWithOptions(dynamodbEndpointResolver),
	)
	if err != nil {
		log.Sugar().Fatal("failed to load aws config")
	}

	dynamodbClient := dynamodb.NewFromConfig(awsConfig)

	dynamodbStorageConfig := storage.DynamoDbStorageConfig{
		StorageTable: storageTable,
	}

	dynamodbStorage := storage.NewDynamoDbStorage(dynamodbClient, dynamodbStorageConfig)

	redirectHandler := NewRedirectHandler(dynamodbStorage, log)

	lambda.Start(redirectHandler.handleRedirect)
}
