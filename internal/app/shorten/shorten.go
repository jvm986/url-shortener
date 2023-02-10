package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jvm986/url-shortener/internal/pkg/shortener"
)

type shortenHandler struct {
	shortener shortener.Shortener
}

func NewShortenHandler(shortener shortener.Shortener) *shortenHandler {
	return &shortenHandler{
		shortener: shortener,
	}
}

func (h *shortenHandler) handleShorten(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	key, _, err := h.shortener.Shorten("example.com")
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       fmt.Sprintf("failed to shorten url [%s]", "example.com"),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       key,
	}, nil
}
