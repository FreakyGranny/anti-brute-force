package app

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/suite"

	"github.com/FreakyGranny/anti-brute-force/internal/mocks"
)

type LimiterSuite struct {
	suite.Suite
	mockCacheCtl *gomock.Controller
	mockCache    *mocks.MockCache
	ctx          context.Context
}

func (s *LimiterSuite) SetupTest() {
	s.mockCacheCtl = gomock.NewController(s.T())
	s.mockCache = mocks.NewMockCache(s.mockCacheCtl)
	s.ctx = context.Background()
}

func (s *LimiterSuite) TearDownTest() {
	s.mockCacheCtl.Finish()
}

func (s *LimiterSuite) TestSimple() {
	fakeTime := clockwork.NewFakeClockAt(time.Date(2020, time.September, 8, 11, 6, 0, 0, time.UTC))
	lim := NewLimiter(s.mockCache, fakeTime, 2, 4, 6)

	s.mockCache.EXPECT().Incr(s.ctx, "LOGIN:test:6", time.Minute).Return(nil)
	s.mockCache.EXPECT().Incr(s.ctx, "PASS:xpass:6", time.Minute).Return(nil)
	s.mockCache.EXPECT().Incr(s.ctx, "IP:127.0.0.1:6", time.Minute).Return(nil)
	s.mockCache.EXPECT().Get(s.ctx, "LOGIN:test:6").Return(1, nil)
	s.mockCache.EXPECT().Get(s.ctx, "PASS:xpass:6").Return(1, nil)
	s.mockCache.EXPECT().Get(s.ctx, "IP:127.0.0.1:6").Return(1, nil)

	c, err := lim.CheckLimits(s.ctx, "test", "xpass", "127.0.0.1")

	s.Require().NoError(err)
	s.Require().True(c)
}

func (s *LimiterSuite) TestLoginLimit() {
	fakeTime := clockwork.NewFakeClockAt(time.Date(2020, time.September, 8, 11, 6, 0, 0, time.UTC))
	lim := NewLimiter(s.mockCache, fakeTime, 2, 4, 6)

	s.mockCache.EXPECT().Incr(s.ctx, "LOGIN:test:6", time.Minute).Return(nil)
	s.mockCache.EXPECT().Get(s.ctx, "LOGIN:test:6").Return(3, nil)

	c, err := lim.CheckLimits(s.ctx, "test", "xpass", "127.0.0.1")

	s.Require().NoError(err)
	s.Require().False(c)
}

func (s *LimiterSuite) TestPassLimit() {
	fakeTime := clockwork.NewFakeClockAt(time.Date(2020, time.September, 8, 11, 6, 0, 0, time.UTC))
	lim := NewLimiter(s.mockCache, fakeTime, 2, 4, 6)

	s.mockCache.EXPECT().Incr(s.ctx, "LOGIN:test:6", time.Minute).Return(nil)
	s.mockCache.EXPECT().Incr(s.ctx, "PASS:xpass:6", time.Minute).Return(nil)
	s.mockCache.EXPECT().Get(s.ctx, "LOGIN:test:6").Return(1, nil)
	s.mockCache.EXPECT().Get(s.ctx, "PASS:xpass:6").Return(5, nil)

	c, err := lim.CheckLimits(s.ctx, "test", "xpass", "127.0.0.1")

	s.Require().NoError(err)
	s.Require().False(c)
}

func (s *LimiterSuite) TestIPLimit() {
	fakeTime := clockwork.NewFakeClockAt(time.Date(2020, time.September, 8, 11, 6, 0, 0, time.UTC))
	lim := NewLimiter(s.mockCache, fakeTime, 2, 4, 6)

	s.mockCache.EXPECT().Incr(s.ctx, "LOGIN:test:6", time.Minute).Return(nil)
	s.mockCache.EXPECT().Incr(s.ctx, "PASS:xpass:6", time.Minute).Return(nil)
	s.mockCache.EXPECT().Incr(s.ctx, "IP:127.0.0.1:6", time.Minute).Return(nil)
	s.mockCache.EXPECT().Get(s.ctx, "LOGIN:test:6").Return(2, nil)
	s.mockCache.EXPECT().Get(s.ctx, "PASS:xpass:6").Return(3, nil)
	s.mockCache.EXPECT().Get(s.ctx, "IP:127.0.0.1:6").Return(7, nil)

	c, err := lim.CheckLimits(s.ctx, "test", "xpass", "127.0.0.1")

	s.Require().NoError(err)
	s.Require().False(c)
}

func (s *LimiterSuite) TestDrop() {
	fakeTime := clockwork.NewFakeClockAt(time.Date(2020, time.September, 8, 15, 31, 0, 0, time.UTC))
	lim := NewLimiter(s.mockCache, fakeTime, 2, 4, 6)

	s.mockCache.EXPECT().Del(s.ctx, "LOGIN:test:31").Return(nil)
	s.mockCache.EXPECT().Del(s.ctx, "PASS:xpass:31").Return(nil)

	err := lim.DropBuckets(s.ctx, "test", "xpass")
	s.Require().NoError(err)
}

func TestLimiterSuite(t *testing.T) {
	suite.Run(t, new(LimiterSuite))
}
