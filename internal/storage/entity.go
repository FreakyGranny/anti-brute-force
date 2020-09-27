package storage

import (
	"context"
)

//go:generate mockgen -source=$GOFILE -destination=../mocks/storage_mock.go -package=mocks .

// WriteStorage storage provider.
type WriteStorage interface {
	AddToWhiteList(ctx context.Context, ip, mask string) error
	AddToBlackList(ctx context.Context, ip, mask string) error
	RemoveFromWhiteList(ctx context.Context, ip, mask string) error
	RemoveFromBlackList(ctx context.Context, ip, mask string) error
}

// ReadStorage storage provider only for reading.
type ReadStorage interface {
	GetBlackList(ctx context.Context) ([]*IPNet, error)
	GetWhiteList(ctx context.Context) ([]*IPNet, error)
}

// IPNet subnet.
type IPNet struct {
	IP   string `db:"ip"`
	Mask string `db:"mask"`
}
