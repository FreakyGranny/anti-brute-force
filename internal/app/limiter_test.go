package app

import (
	"context"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/FreakyGranny/anti-brute-force/internal/mocks"
	"github.com/golang/mock/gomock"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/suite"
)

type CacheLimiterSuite struct {
	suite.Suite
	mockCacheCtl *gomock.Controller
	mockCache    *mocks.MockCache
	ctx          context.Context
}

func (s *CacheLimiterSuite) SetupTest() {
	s.mockCacheCtl = gomock.NewController(s.T())
	s.mockCache = mocks.NewMockCache(s.mockCacheCtl)
	s.ctx = context.Background()
}

func (s *CacheLimiterSuite) TearDownTest() {
	s.mockCacheCtl.Finish()
}

func (s *CacheLimiterSuite) setLoginExpect(login string, minute, loginReturn int) {
	strMinute := strconv.Itoa(minute)
	loginKey := strings.Join([]string{"LOGIN", login, strMinute}, ":")
	s.mockCache.EXPECT().Incr(s.ctx, loginKey, time.Minute).Return(nil)
	s.mockCache.EXPECT().Get(s.ctx, loginKey).Return(loginReturn, nil)
}

func (s *CacheLimiterSuite) setPassExpect(password string, minute, passReturn int) {
	strMinute := strconv.Itoa(minute)
	passKey := strings.Join([]string{"PASSWORD", password, strMinute}, ":")
	s.mockCache.EXPECT().Incr(s.ctx, passKey, time.Minute).Return(nil)
	s.mockCache.EXPECT().Get(s.ctx, passKey).Return(passReturn, nil)
}

func (s *CacheLimiterSuite) setIPExpect(ip string, minute, ipReturn int) {
	strMinute := strconv.Itoa(minute)
	ipKey := strings.Join([]string{"IP", ip, strMinute}, ":")
	s.mockCache.EXPECT().Incr(s.ctx, ipKey, time.Minute).Return(nil)
	s.mockCache.EXPECT().Get(s.ctx, ipKey).Return(ipReturn, nil)
}

func (s *CacheLimiterSuite) TestSimple() {
	fakeTime := clockwork.NewFakeClockAt(time.Date(2020, time.September, 8, 11, 6, 0, 0, time.UTC))
	lim := NewCacheLimiter(s.mockCache, fakeTime, 2, 4, 6)

	s.setLoginExpect("test", 6, 1)
	s.setPassExpect("xpass", 6, 1)
	s.setIPExpect("127.0.0.1", 6, 1)

	c, err := lim.CheckLimits(s.ctx, "test", "xpass", "127.0.0.1")

	s.Require().NoError(err)
	s.Require().True(c)
}

func (s *CacheLimiterSuite) TestLoginLimit() {
	fakeTime := clockwork.NewFakeClockAt(time.Date(2020, time.September, 8, 11, 6, 0, 0, time.UTC))
	lim := NewCacheLimiter(s.mockCache, fakeTime, 2, 4, 6)

	s.setLoginExpect("test", 6, 3)

	c, err := lim.CheckLimits(s.ctx, "test", "xpass", "127.0.0.1")

	s.Require().NoError(err)
	s.Require().False(c)
}

func (s *CacheLimiterSuite) TestPassLimit() {
	fakeTime := clockwork.NewFakeClockAt(time.Date(2020, time.September, 8, 11, 6, 0, 0, time.UTC))
	lim := NewCacheLimiter(s.mockCache, fakeTime, 2, 4, 6)

	s.setLoginExpect("test", 6, 1)
	s.setPassExpect("xpass", 6, 5)

	c, err := lim.CheckLimits(s.ctx, "test", "xpass", "127.0.0.1")

	s.Require().NoError(err)
	s.Require().False(c)
}

func (s *CacheLimiterSuite) TestIPLimit() {
	fakeTime := clockwork.NewFakeClockAt(time.Date(2020, time.September, 24, 9, 39, 0, 0, time.UTC))
	lim := NewCacheLimiter(s.mockCache, fakeTime, 2, 4, 6)

	s.setLoginExpect("login", 39, 2)
	s.setPassExpect("supersecretpassword", 39, 3)
	s.setIPExpect("10.10.99.101", 39, 7)

	c, err := lim.CheckLimits(s.ctx, "login", "supersecretpassword", "10.10.99.101")

	s.Require().NoError(err)
	s.Require().False(c)
}

func (s *CacheLimiterSuite) TestDrop() {
	fakeTime := clockwork.NewFakeClockAt(time.Date(2020, time.September, 8, 15, 31, 0, 0, time.UTC))
	lim := NewCacheLimiter(s.mockCache, fakeTime, 2, 4, 6)

	s.mockCache.EXPECT().Del(s.ctx, "LOGIN:test:31").Return(nil)
	s.mockCache.EXPECT().Del(s.ctx, "IP:10.10.34.5:31").Return(nil)

	err := lim.DropBuckets(s.ctx, "test", "10.10.34.5")
	s.Require().NoError(err)
}

func TestCacheLimiterSuite(t *testing.T) {
	suite.Run(t, new(CacheLimiterSuite))
}
