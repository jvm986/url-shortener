package storage

import (
	"context"
)

type NotFoundError struct{}

func (e *NotFoundError) Error() string {
	return "key not found in storage"

}

//go:generate mockgen -destination=mocks/mock_storage.go -package=storagemocks . Storage
type Storage interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string) error
}
