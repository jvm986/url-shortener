package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/jvm986/url-shortener/internal/pkg/shortener"
	"github.com/jvm986/url-shortener/internal/pkg/storage"
)

func main() {
	ctx := context.Background()

	pathLengthString, ok := os.LookupEnv("PATH_LENGTH")
	if !ok {
		panic("failed to lookup PATH_LENGTH from environment variables")
	}

	pathLength, err := strconv.Atoi(pathLengthString)
	if err != nil {
		panic(fmt.Sprintf("PATH_LENGTH [%s] is not an integer", pathLengthString))
	}

	md5Config := shortener.MD5ShortenerConfig{
		PathLength: pathLength,
	}
	md5shortener := shortener.NewMD5Shortener(md5Config)

	endpoint, ok := os.LookupEnv("ENDPOINT")
	if !ok {
		panic("failed to lookup ENDPOINT from environment variables")
	}

	storageTable, ok := os.LookupEnv("STORAGE_TABLE_NAME")
	if !ok {
		panic("failed to lookup STORAGE_TABLE_NAME from environment variables")
	}

	env, ok := os.LookupEnv("ENV")
	if !ok {
		panic("failed to lookup ENV from environment variables")
	}
	dynamodbEndpoint, ok := os.LookupEnv("DDB_ENDPOINT")
	if !ok {
		panic("failed to lookup DDB_ENDPOINT from environment variables")
	}

	awsRegion, ok := os.LookupEnv("REGION")
	if !ok {
		panic("failed to lookup REGION from environment variables")
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
		panic("failed to load aws config")
	}

	dynamodbClient := dynamodb.NewFromConfig(awsConfig)

	dynamodbStorageConfig := storage.DynamoDbStorageConfig{
		StorageTable: storageTable,
	}

	dynamodbStorage := storage.NewDynamoDbStorage(dynamodbClient, dynamodbStorageConfig)

	shortenHandler := NewShortenHandler(md5shortener, dynamodbStorage, endpoint)

	lambda.Start(shortenHandler.handleShorten)
}
