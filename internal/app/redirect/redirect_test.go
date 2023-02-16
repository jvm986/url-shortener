package main

import (
	"context"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang/mock/gomock"
	storage "github.com/jvm986/url-shortener/internal/pkg/storage/mocks"
	"github.com/stretchr/testify/assert"
)

func TestHandleRedirect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStorage := storage.NewMockStorage(ctrl)

	mockStorage.
		EXPECT().
		Get(gomock.Any(), "key").
		Return("redirect_url", nil)

	h := redirectHandler{
		storage: mockStorage,
	}

	actual, err := h.handleRedirect(context.TODO(), events.APIGatewayProxyRequest{
		PathParameters: map[string]string{
			"key": "key",
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, events.APIGatewayProxyResponse{
		StatusCode: http.StatusPermanentRedirect,
		Headers: map[string]string{
			"location": "redirect_url",
		},
	}, actual)
}
