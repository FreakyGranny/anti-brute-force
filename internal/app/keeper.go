package app

import (
	"context"
	"sync"
	"time"

	"github.com/FreakyGranny/anti-brute-force/internal/storage"
	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=$GOFILE -destination=../mocks/keeper_mock.go -package=mocks IPKeeper

// IPKeeper ip lists storage.
type IPKeeper interface {
	GetBlacklist() []*storage.IPNet
	GetWhitelist() []*storage.IPNet
	Refresh(ctx context.Context) error
	Watch(ctx context.Context, interval time.Duration)
}

// MemIPKeeper provides access to black and white lists.
type MemIPKeeper struct {
	storage   storage.ReadStorage
	bmu       sync.RWMutex
	wmu       sync.RWMutex
	blacklist []*storage.IPNet
	whitelist []*storage.IPNet
}

// NewMemIPKeeper returns keeper instance.
func NewMemIPKeeper(s storage.ReadStorage) *MemIPKeeper {
	return &MemIPKeeper{
		storage:   s,
		blacklist: make([]*storage.IPNet, 0),
		whitelist: make([]*storage.IPNet, 0),
	}
}

// GetBlacklist returns blacklist from memory storage.
func (k *MemIPKeeper) GetBlacklist() []*storage.IPNet {
	k.bmu.RLock()
	v := k.blacklist
	k.bmu.RUnlock()

	return v
}

// GetWhitelist returns whitelist from memory storage.
func (k *MemIPKeeper) GetWhitelist() []*storage.IPNet {
	k.wmu.RLock()
	v := k.whitelist
	k.wmu.RUnlock()

	return v
}

func (k *MemIPKeeper) refreshBlackList(ctx context.Context) error {
	k.bmu.Lock()
	defer k.bmu.Unlock()

	values, err := k.storage.GetBlackList(ctx)
	if err != nil {
		return err
	}

	k.blacklist = values

	return nil
}

func (k *MemIPKeeper) refreshWhiteList(ctx context.Context) error {
	k.wmu.Lock()
	defer k.wmu.Unlock()

	values, err := k.storage.GetWhiteList(ctx)
	if err != nil {
		return err
	}

	k.whitelist = values

	return nil
}

// Refresh refresh black and white lists.
func (k *MemIPKeeper) Refresh(ctx context.Context) error {
	if err := k.refreshBlackList(ctx); err != nil {
		return err
	}
	return k.refreshWhiteList(ctx)
}

// Watch periodically refreshing black and white lists.
func (k *MemIPKeeper) Watch(ctx context.Context, interval time.Duration) {
	log.Info().Msgf("refreshing every %s", interval)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := k.refreshBlackList(ctx); err != nil {
				log.Error().Err(err).Msg("error while refreshing blacklist")
			}
			if err := k.refreshWhiteList(ctx); err != nil {
				log.Error().Err(err).Msg("error while refreshing whitelist")
			}
		}
	}
}
