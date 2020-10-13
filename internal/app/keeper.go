package app

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/FreakyGranny/anti-brute-force/internal/storage"
	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=$GOFILE -destination=../mocks/keeper_mock.go -package=mocks IPKeeper

// IPKeeper ip lists storage.
type IPKeeper interface {
	GetBlacklist() []*net.IPNet
	GetWhitelist() []*net.IPNet
	Refresh(ctx context.Context) error
	Watch(ctx context.Context, interval time.Duration)
}

// MemIPKeeper provides access to black and white lists.
type MemIPKeeper struct {
	storage   storage.ReadStorage
	bmu       sync.RWMutex
	wmu       sync.RWMutex
	blacklist []*net.IPNet
	whitelist []*net.IPNet
}

// NewMemIPKeeper returns keeper instance.
func NewMemIPKeeper(s storage.ReadStorage) *MemIPKeeper {
	return &MemIPKeeper{
		storage: s,
	}
}

// GetBlacklist returns blacklist from memory storage.
func (k *MemIPKeeper) GetBlacklist() []*net.IPNet {
	blacklist := make([]*net.IPNet, len(k.blacklist))
	copy(blacklist, k.blacklist)

	return blacklist
}

// GetWhitelist returns whitelist from memory storage.
func (k *MemIPKeeper) GetWhitelist() []*net.IPNet {
	whitelist := make([]*net.IPNet, len(k.whitelist))
	copy(whitelist, k.whitelist)

	return whitelist
}

func (k *MemIPKeeper) refreshBlackList(ctx context.Context) error {
	k.bmu.Lock()
	defer k.bmu.Unlock()

	values, err := k.storage.GetBlackList(ctx)
	if err != nil {
		return err
	}

	k.blacklist = make([]*net.IPNet, 0, len(values))
	for _, v := range values {
		k.blacklist = append(k.blacklist, parseNet(v))
	}

	return nil
}

func parseNet(subnet *storage.IPNet) *net.IPNet {
	byteMask := net.ParseIP(subnet.Mask).To4()
	return &net.IPNet{
		IP:   net.ParseIP(subnet.IP),
		Mask: net.IPv4Mask(byteMask[0], byteMask[1], byteMask[2], byteMask[3]),
	}
}

func (k *MemIPKeeper) refreshWhiteList(ctx context.Context) error {
	k.wmu.Lock()
	defer k.wmu.Unlock()

	values, err := k.storage.GetWhiteList(ctx)
	if err != nil {
		return err
	}

	k.whitelist = make([]*net.IPNet, 0, len(values))
	for _, v := range values {
		k.whitelist = append(k.whitelist, parseNet(v))
	}

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
