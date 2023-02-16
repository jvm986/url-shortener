package storage

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/golang/mock/gomock"
	dynamodbclient "github.com/jvm986/url-shortener/internal/pkg/clients/dynamodbclient/mocks"
	"github.com/stretchr/testify/assert"
)

// TODO add more test cases
func TestGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := dynamodbclient.NewMockDynamoDbClient(ctrl)

	mockClient.
		EXPECT().
		GetItem(gomock.Any(), &dynamodb.GetItemInput{
			TableName: aws.String("table"),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "key"},
			},
		}, gomock.Any()).
		Return(&dynamodb.GetItemOutput{
			Item: map[string]types.AttributeValue{
				"id":    &types.AttributeValueMemberS{Value: "key"},
				"value": &types.AttributeValueMemberS{Value: "value"},
			},
		}, nil)

	s := DynamoDbStorage{
		client:       mockClient,
		storageTable: "table",
	}

	actual, err := s.Get(context.TODO(), "key")
	assert.NoError(t, err)
	assert.Equal(t, "value", actual)
}

// TODO add more test cases
func TestSet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := dynamodbclient.NewMockDynamoDbClient(ctrl)

	mockClient.
		EXPECT().
		PutItem(gomock.Any(), &dynamodb.PutItemInput{
			TableName: aws.String("table"),
			Item: map[string]types.AttributeValue{
				"id":    &types.AttributeValueMemberS{Value: "key"},
				"value": &types.AttributeValueMemberS{Value: "value"},
			},
		}, gomock.Any()).
		Return(nil, nil)

	s := DynamoDbStorage{
		client:       mockClient,
		storageTable: "table",
	}

	err := s.Set(context.TODO(), "key", "value")
	assert.NoError(t, err)
}
