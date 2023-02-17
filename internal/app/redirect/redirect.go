package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jvm986/url-shortener/internal/pkg/storage"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type redirectHandler struct {
	storage storage.Storage
	log     *zap.Logger
}

func NewRedirectHandler(dynamoDbStorage storage.Storage, log *zap.Logger) *redirectHandler {
	return &redirectHandler{
		storage: dynamoDbStorage,
		log:     log,
	}
}

func (h *redirectHandler) handleRedirect(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	p := request.PathParameters
	key, ok := p["key"]
	if !ok {
		e := "missing path parameter"
		h.log.Sugar().With("path", request.Path).Error(e)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       e,
		}, nil
	}

	url, err := h.storage.Get(ctx, key)
	var nfe *storage.NotFoundError
	if errors.As(err, &nfe) {
		h.log.Sugar().With("key", key).Error(err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       fmt.Sprintf("key [%s] not found in storage", key),
		}, nil
	}
	if err != nil {
		h.log.Sugar().With("key", key).Error(err.Error())
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
