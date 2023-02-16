package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"
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

	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic("failed to load default aws config")
	}

	dynamodbClient := dynamodb.NewFromConfig(awsConfig)

	dynamodbStorageConfig := storage.DynamoDbStorageConfig{
		StorageTable: storageTable,
	}

	dynamodbStorage := storage.NewDynamoDbStorage(dynamodbClient, dynamodbStorageConfig)

	shortenHandler := NewShortenHandler(md5shortener, dynamodbStorage, endpoint)

	lambda.Start(shortenHandler.handleShorten)
}
