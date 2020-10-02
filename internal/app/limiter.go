package app

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/FreakyGranny/anti-brute-force/internal/cache"
	"github.com/jonboulle/clockwork"
)

const period = time.Minute

const (
	loginPrefix    = "LOGIN"
	passwordPrefix = "PASSWORD"
	ipPrefix       = "IP"
)

//go:generate mockgen -source=$GOFILE -destination=../mocks/limiter_mock.go -package=mocks Limiter

// Limiter ...
type Limiter interface {
	CheckLimits(ctx context.Context, login string, password string, ip string) (bool, error)
	DropBuckets(ctx context.Context, login string, ip string) error
}

// CacheLimiter limit checker based on cache.
type CacheLimiter struct {
	cache         cache.Cache
	clock         clockwork.Clock
	loginLimit    int
	passwordLimit int
	IPLimit       int
}

// NewCacheLimiter returns new cache limiter instance.
func NewCacheLimiter(cache cache.Cache, clock clockwork.Clock, loginLimit int, passwordLimit int, ipLimit int) *CacheLimiter {
	return &CacheLimiter{
		cache:         cache,
		clock:         clock,
		loginLimit:    loginLimit,
		passwordLimit: passwordLimit,
		IPLimit:       ipLimit,
	}
}

func buildKey(id string, val string, min int) string {
	return strings.Join([]string{id, val, strconv.Itoa(min)}, ":")
}

// CheckLimits returns true if request is allowed.
func (lim *CacheLimiter) CheckLimits(ctx context.Context, login string, password string, ip string) (bool, error) {
	p := lim.clock.Now().Minute()

	lKey := buildKey(loginPrefix, login, p)
	err := lim.cache.Incr(ctx, lKey, period)
	if err != nil {
		return false, err
	}
	v, err := lim.cache.Get(ctx, lKey)
	if err != nil {
		return false, err
	}
	if v > lim.loginLimit {
		return false, nil
	}

	pKey := buildKey(passwordPrefix, password, p)
	err = lim.cache.Incr(ctx, pKey, period)
	if err != nil {
		return false, err
	}
	v, err = lim.cache.Get(ctx, pKey)
	if err != nil {
		return false, err
	}
	if v > lim.passwordLimit {
		return false, nil
	}

	IPKey := buildKey(ipPrefix, ip, p)
	err = lim.cache.Incr(ctx, IPKey, period)
	if err != nil {
		return false, err
	}
	v, err = lim.cache.Get(ctx, IPKey)
	if err != nil {
		return false, err
	}
	if v > lim.IPLimit {
		return false, nil
	}

	return true, nil
}

// DropBuckets deletes buckets for given login, ip.
func (lim *CacheLimiter) DropBuckets(ctx context.Context, login string, ip string) error {
	p := lim.clock.Now().Minute()

	lKey := buildKey(loginPrefix, login, p)
	err := lim.cache.Del(ctx, lKey)
	if err != nil {
		return err
	}
	pKey := buildKey(ipPrefix, ip, p)

	return lim.cache.Del(ctx, pKey)
}
