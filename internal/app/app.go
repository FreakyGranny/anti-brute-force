package app

import (
	"context"
	"net"

	"github.com/FreakyGranny/anti-brute-force/internal/cache"
	"github.com/FreakyGranny/anti-brute-force/internal/storage"

	"github.com/jonboulle/clockwork"
	"github.com/rs/zerolog/log"
)

// Application business logic.
type Application interface {
	AddToWhiteList(ctx context.Context, ip string, mask string) error
	AddToBlackList(ctx context.Context, ip string, mask string) error
	RemoveFromWhiteList(ctx context.Context, id int) error
	RemoveFromBlackList(ctx context.Context, id int) error
	DropStat(ctx context.Context, login string, password string) error
	CheckRate(ctx context.Context, login string, password string, ip string) (bool, error)
}

// App ...
type App struct {
	storage storage.Storage
	limiter *Limiter
}

// New returns application instance.
func New(storage storage.Storage, cache cache.Cache, clock clockwork.Clock, loginLimit int, passwordLimit int, IPLimit int) *App {
	return &App{
		storage: storage,
		limiter: NewLimiter(
			cache,
			clock,
			loginLimit,
			passwordLimit,
			IPLimit,
		),
	}
}

func ipInSubnet(ip string, subnets []*storage.IPNet) (bool, error) {
	ipCheck := func(ip string, subnet string) (bool, error) {
		_, ipv4Net, err := net.ParseCIDR(subnet)
		if err != nil {
			return false, err
		}

		return ipv4Net.Contains(net.ParseIP(ip)), nil
	}
	for _, n := range subnets {
		contains, err := ipCheck(ip, n.Subnet)
		if err != nil {
			return false, err
		}
		if contains {
			log.Info().Msg(n.Subnet)
			return true, nil
		}
	}

	return false, nil
}

// CheckBlackList returns true if request in black list.
func (a *App) CheckBlackList(ctx context.Context, ip string) (bool, error) {
	sns, err := a.storage.GetBlackList(ctx)
	if err != nil {
		return false, err
	}

	return ipInSubnet(ip, sns)
}

// CheckWhiteList returns true if request in white list.
func (a *App) CheckWhiteList(ctx context.Context, ip string) (bool, error) {
	sns, err := a.storage.GetWhiteList(ctx)
	if err != nil {
		return false, err
	}

	return ipInSubnet(ip, sns)
}

// CheckRate returns true if request is allowed.
func (a *App) CheckRate(ctx context.Context, login string, password string, ip string) (bool, error) {
	in, err := a.CheckWhiteList(ctx, ip)
	if err != nil {
		return false, err
	}
	if in {
		return true, nil
	}
	in, err = a.CheckBlackList(ctx, ip)
	if err != nil {
		return false, err
	}
	if in {
		return false, nil
	}

	return a.limiter.CheckLimits(ctx, login, password, ip)
}

// DropStat drops all stats for given login, password
func (a *App) DropStat(ctx context.Context, login string, password string) error {
	return a.limiter.DropBuckets(ctx, login, password)
}

// AddToBlackList adding ip and mask to blacklist
func (a *App) AddToBlackList(ctx context.Context, ip string, mask string) error {
	return a.storage.AddToBlackList(ctx, createRecord(ip, mask))
}

// AddToWhiteList adding ip and mask to whitelist
func (a *App) AddToWhiteList(ctx context.Context, ip string, mask string) error {
	return a.storage.AddToWhiteList(ctx, createRecord(ip, mask))
}

func createRecord(ip string, mask string) *storage.IPNet {
	byteMask := net.ParseIP(mask).To4()
	CIDR := net.IPNet{
		IP:   net.ParseIP(ip),
		Mask: net.IPv4Mask(byteMask[0], byteMask[1], byteMask[2], byteMask[3]),
	}

	return &storage.IPNet{Subnet: CIDR.String()}
}

// RemoveFromWhiteList removes record from whitelist
func (a *App) RemoveFromWhiteList(ctx context.Context, id int) error {
	return a.storage.RemoveFromWhiteList(ctx, id)
}

// RemoveFromBlackList removes record from blacklist
func (a *App) RemoveFromBlackList(ctx context.Context, id int) error {
	return a.storage.RemoveFromBlackList(ctx, id)
}
