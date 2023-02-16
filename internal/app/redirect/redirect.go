package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jvm986/url-shortener/internal/pkg/storage"
	"github.com/pkg/errors"
)

type redirectHandler struct {
	storage storage.Storage
}

func NewRedirectHandler(dynamoDbStorage storage.Storage) *redirectHandler {
	return &redirectHandler{
		storage: dynamoDbStorage,
	}
}

func (h *redirectHandler) handleRedirect(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	p := request.PathParameters
	key, ok := p["key"]
	if !ok {
		e := fmt.Sprintf("missing path parameter on [%s]", request.Path)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       e,
		}, nil
	}

	url, err := h.storage.Get(ctx, key)
	var nfe *storage.NotFoundError
	if errors.As(err, &nfe) {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       fmt.Sprintf("key [%s] not found in storage", key),
		}, nil
	}
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "failed to get from storage",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusPermanentRedirect,
		Headers: map[string]string{
			"location": url,
		},
	}, nil
}
