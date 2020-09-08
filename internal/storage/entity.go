package storage

import (
	"context"
)

//go:generate mockgen -source=$GOFILE -destination=../mocks/storage_mock.go -package=mocks .

// Storage storage provider.
type Storage interface {
	AddToWhiteList(ctx context.Context, e *IPNet) error
	AddToBlackList(ctx context.Context, e *IPNet) error
	RemoveFromWhiteList(ctx context.Context, id int) error
	RemoveFromBlackList(ctx context.Context, id int) error
	GetBlackList(ctx context.Context) ([]*IPNet, error)
	GetWhiteList(ctx context.Context) ([]*IPNet, error)
	Close() error
}

// IPNet subnet.
type IPNet struct {
	ID     int    `db:"id"`
	Subnet string `db:"subnet"`
}
