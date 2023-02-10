package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jvm986/url-shortener/internal/pkg/shortener"
)

func main() {
	md5Config := shortener.MD5ShortenerConfig{
		PathLength: 10,
	}
	md5shortener := shortener.NewMD5Shortener(md5Config)
	shortenHandler := NewShortenHandler(md5shortener)

	lambda.Start(shortenHandler.handleShorten)
}
