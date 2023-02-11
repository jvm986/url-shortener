package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jvm986/url-shortener/internal/pkg/shortener"
)

type shortenHandler struct {
	shortener shortener.Shortener
	endpoint  string
}

type RequestBody struct {
	Url string `json:"url"`
}

type ResponseBody struct {
	Url         string `json:"url"`
	Key         string `json:"key"`
	RedirectUrl string `json:"redirect_url"`
}

func NewShortenHandler(shortener shortener.Shortener, endpoint string) *shortenHandler {
	return &shortenHandler{
		shortener: shortener,
		endpoint:  endpoint,
	}
}

func (h *shortenHandler) handleShorten(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	b := &RequestBody{}
	err := json.Unmarshal([]byte(request.Body), b)
	if err != nil {
		e := fmt.Sprintf("failed to unmarshal request body [%s]", request.Body)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       e,
		}, nil
	}

	if b.Url == "" {
		e := fmt.Sprintf("bad request body [%s]", request.Body)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       e,
		}, nil
	}

	key, value, err := h.shortener.Shorten(b.Url)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       fmt.Sprintf("failed to shorten url [%s]", b.Url),
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
