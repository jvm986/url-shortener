package main

import (
	"context"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang/mock/gomock"
	"github.com/jvm986/url-shortener/internal/pkg/domain"
	shortener "github.com/jvm986/url-shortener/internal/pkg/shortener/mocks"
	storage "github.com/jvm986/url-shortener/internal/pkg/storage/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestHandleShorten(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShortener := shortener.NewMockShortener(ctrl)
	mockShortener.
		EXPECT().
		Shorten("valid-url").
		Return("key", "valid-url", nil)

	mockStorage := storage.NewMockStorage(ctrl)
	mockStorage.
		EXPECT().
		Set(gomock.Any(), "key", "valid-url").
		Return(nil)

	h := shortenHandler{
		endpoint:  "base/",
		storage:   mockStorage,
		shortener: mockShortener,
		log:       zap.NewNop(),
	}

	actual, err := h.handleShorten(context.TODO(), events.APIGatewayProxyRequest{
		Body: `{"url": "valid-url"}`,
	})

	assert.NoError(t, err)
	assert.Equal(t, events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Headers:    domain.GetCorsHeaders(),
		Body:       `{"url":"valid-url","key":"key","redirect_url":"base/short/key"}`,
	}, actual)
}

func TestHandleShorten_InvalidJson(t *testing.T) {
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)

	h := shortenHandler{
		log: observedLogger,
	}

	_, err := h.handleShorten(context.TODO(), events.APIGatewayProxyRequest{
		Body: `invalid json`,
	})

	assert.NoError(t, err)
	assert.Equal(t, 1, observedLogs.Len())
	assert.Equal(t, "failed to unmarshal request body: invalid character 'i' looking for beginning of value", observedLogs.All()[0].Message)
}

func TestHandleShorten_NoUrl(t *testing.T) {
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)

	h := shortenHandler{
		log: observedLogger,
	}

	_, err := h.handleShorten(context.TODO(), events.APIGatewayProxyRequest{
		Body: `{"not-url": "value"}`,
	})

	assert.NoError(t, err)
	assert.Equal(t, 1, observedLogs.Len())
	assert.Equal(t, "no or empty url field in request body", observedLogs.All()[0].Message)
}
