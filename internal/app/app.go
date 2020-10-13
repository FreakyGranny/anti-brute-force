package app

import (
	"context"
	"errors"
	"net"

	"github.com/FreakyGranny/anti-brute-force/internal/storage"
)

// ErrInvalidArgument given string is not valid ip.
var ErrInvalidArgument = errors.New("invalid argument")

// Application business logic.
type Application interface {
	AddToWhiteList(ctx context.Context, ip, mask string) error
	AddToBlackList(ctx context.Context, ip, mask string) error
	RemoveFromWhiteList(ctx context.Context, ip, mask string) error
	RemoveFromBlackList(ctx context.Context, ip, mask string) error
	DropStat(ctx context.Context, login, ip string) error
	CheckRate(ctx context.Context, login, password, ip string) (bool, error)
}

// App business logic implementation.
type App struct {
	storage storage.WriteStorage
	limiter Limiter
	keeper  IPKeeper
}

// New returns application instance.
func New(storage storage.WriteStorage, keeper IPKeeper, limiter Limiter) *App {
	return &App{
		storage: storage,
		limiter: limiter,
		keeper:  keeper,
	}
}

// CheckRate returns true if request is allowed.
func (a *App) CheckRate(ctx context.Context, login, password, ip string) (bool, error) {
	if !IsValidIPFormat(ip) {
		return false, ErrInvalidArgument
	}
	if ipInSubnet(ip, a.keeper.GetWhitelist()) {
		return true, nil
	}
	if ipInSubnet(ip, a.keeper.GetBlacklist()) {
		return false, nil
	}

	return a.limiter.CheckLimits(ctx, login, password, ip)
}

func ipInSubnet(ip string, subnets []*net.IPNet) bool {
	for _, n := range subnets {
		if n.Contains(net.ParseIP(ip)) {
			return true
		}
	}

	return false
}

// DropStat drops all stats for given login, ip.
func (a *App) DropStat(ctx context.Context, login, ip string) error {
	if !IsValidIPFormat(ip) {
		return ErrInvalidArgument
	}

	return a.limiter.DropBuckets(ctx, login, ip)
}

// AddToBlackList adding ip and mask to blacklist.
func (a *App) AddToBlackList(ctx context.Context, ip, mask string) error {
	if !IsValidIPFormat(ip) || !IsValidIPFormat(mask) {
		return ErrInvalidArgument
	}

	return a.storage.AddToBlackList(ctx, ip, mask)
}

// AddToWhiteList adding ip and mask to whitelist.
func (a *App) AddToWhiteList(ctx context.Context, ip, mask string) error {
	if !IsValidIPFormat(ip) || !IsValidIPFormat(mask) {
		return ErrInvalidArgument
	}

	return a.storage.AddToWhiteList(ctx, ip, mask)
}

// RemoveFromWhiteList removes record from whitelist.
func (a *App) RemoveFromWhiteList(ctx context.Context, ip, mask string) error {
	if !IsValidIPFormat(ip) || !IsValidIPFormat(mask) {
		return ErrInvalidArgument
	}

	return a.storage.RemoveFromWhiteList(ctx, ip, mask)
}

// RemoveFromBlackList removes record from blacklist.
func (a *App) RemoveFromBlackList(ctx context.Context, ip, mask string) error {
	if !IsValidIPFormat(ip) || !IsValidIPFormat(mask) {
		return ErrInvalidArgument
	}

	return a.storage.RemoveFromBlackList(ctx, ip, mask)
}

// IsValidIPFormat check string is valid ip address or mask.
func IsValidIPFormat(ip string) bool {
	return net.ParseIP(ip) != nil
}
