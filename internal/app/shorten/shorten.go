package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jvm986/url-shortener/internal/pkg/shortener"
	"github.com/jvm986/url-shortener/internal/pkg/storage"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type shortenHandler struct {
	shortener shortener.Shortener
	storage   storage.Storage
	endpoint  string
	log       *zap.Logger
}

type RequestBody struct {
	Url string `json:"url"`
}

type ResponseBody struct {
	Url         string `json:"url"`
	Key         string `json:"key"`
	RedirectUrl string `json:"redirect_url"`
}

func NewShortenHandler(
	shortener shortener.Shortener,
	storage storage.Storage,
	endpoint string,
	log *zap.Logger,
) *shortenHandler {
	return &shortenHandler{
		shortener: shortener,
		storage:   storage,
		endpoint:  endpoint,
		log:       log,
	}
}

func (h *shortenHandler) handleShorten(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	b := &RequestBody{}
	err := json.Unmarshal([]byte(request.Body), b)
	if err != nil {
		e := "failed to unmarshal request body"
		h.log.Sugar().With("requestBody", request.Body).Error(errors.Wrap(err, e))
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       e,
		}, nil
	}

	if b.Url == "" {
		e := "no or empty url field in request body"
		h.log.Sugar().With("requestBody", request.Body).Error(e)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       e,
		}, nil
	}

	key, value, err := h.shortener.Shorten(b.Url)
	if err != nil {
		h.log.Sugar().With("url", b.Url).Error(err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       fmt.Sprintf("failed to shorten url [%s]", b.Url),
		}, nil
	}

	err = h.storage.Set(ctx, key, value)
	if err != nil {
		h.log.Sugar().With("key", key, "value", value).Error(err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("failed to set storage for [%s]", key),
		}, nil
	}

	responseBody := ResponseBody{
		Url:         value,
		Key:         key,
		RedirectUrl: h.endpoint + "short/" + key,
	}

	responseBytes, err := json.Marshal(responseBody)
	if err != nil {
		e := "failed to marshal response body"
		h.log.Sugar().With("responseBody", responseBody).Error(err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       e,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       string(responseBytes),
	}, nil
}
