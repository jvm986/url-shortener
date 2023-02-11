package main

import (
	"context"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang/mock/gomock"
	shortener "github.com/jvm986/url-shortener/internal/pkg/shortener/mocks"
	"github.com/stretchr/testify/assert"
)

func TestHandleShorten(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShortener := shortener.NewMockShortener(ctrl)
	mockShortener.
		EXPECT().
		Shorten("valid-url").
		Return("key", "valid-url", nil)

	h := shortenHandler{
		shortener: mockShortener,
		endpoint:  "base/",
	}

	actual, err := h.handleShorten(context.TODO(), events.APIGatewayProxyRequest{
		Body: `{"url": "valid-url"}`,
	})

	assert.NoError(t, err)
	assert.Equal(t, events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       `{"url":"valid-url","key":"key","redirect_url":"base/short/key"}`,
	}, actual)
}
