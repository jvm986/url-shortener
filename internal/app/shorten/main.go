package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jvm986/url-shortener/internal/pkg/shortener"
)

func main() {
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

	shortenHandler := NewShortenHandler(md5shortener, endpoint)

	lambda.Start(shortenHandler.handleShorten)
}
