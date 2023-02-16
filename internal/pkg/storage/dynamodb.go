package storage

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
)

type DynamoDbStorage struct {
	client       *dynamodb.Client
	storageTable string
}

type DynamoDbStorageConfig struct {
	StorageTable string
}

type Entry struct {
	Value string `dynamodbav:"value"`
}

func NewDynamoDbStorage(c *dynamodb.Client, cfg DynamoDbStorageConfig) Storage {
	return DynamoDbStorage{
		client:       c,
		storageTable: cfg.StorageTable,
	}
}

func (s DynamoDbStorage) Get(ctx context.Context, key string) (string, error) {
	response, err := s.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(s.storageTable),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: key},
		},
	})
	if err != nil {
		return "", errors.Wrapf(err, "failed to get item [%s] from dynamodb", key)
	}

	if response.Item == nil {
		return "", errors.Wrap(&NotFoundError{}, key)
	}

	entry := Entry{}
	err = attributevalue.UnmarshalMap(response.Item, &entry)
	if err != nil {
		return "", errors.Wrap(err, "failed to unmarshal entry")
	}

	return entry.Value, errors.Wrapf(err, "failed to get item [%s] from dynamodb", key)
}

func (s DynamoDbStorage) Set(ctx context.Context, key string, value string) error {
	_, err := s.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(s.storageTable),
		Item: map[string]types.AttributeValue{
			"id":    &types.AttributeValueMemberS{Value: key},
			"value": &types.AttributeValueMemberS{Value: value},
		},
	})

	return errors.Wrapf(err, "failed to put item [%s] to dynamodb", key)
}
