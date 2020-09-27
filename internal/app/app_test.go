package app

import (
	"context"
	"testing"

	"github.com/FreakyGranny/anti-brute-force/internal/mocks"
	"github.com/FreakyGranny/anti-brute-force/internal/storage"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type AppSuite struct {
	suite.Suite
	mockStorageCtl *gomock.Controller
	mockStorage    *mocks.MockWriteStorage
	mockKeeperCtl  *gomock.Controller
	mockKeeper     *mocks.MockIPKeeper
	mockLimiterCtl *gomock.Controller
	mockLimiter    *mocks.MockLimiter
	application    *App
	ctx            context.Context
}

func (s *AppSuite) SetupTest() {
	s.mockStorageCtl = gomock.NewController(s.T())
	s.mockStorage = mocks.NewMockWriteStorage(s.mockStorageCtl)
	s.mockKeeperCtl = gomock.NewController(s.T())
	s.mockKeeper = mocks.NewMockIPKeeper(s.mockKeeperCtl)
	s.mockLimiterCtl = gomock.NewController(s.T())
	s.mockLimiter = mocks.NewMockLimiter(s.mockLimiterCtl)
	s.application = New(s.mockStorage, s.mockKeeper, s.mockLimiter)
	s.ctx = context.Background()
}

func (s *AppSuite) TearDownTest() {
	s.mockStorageCtl.Finish()
	s.mockKeeperCtl.Finish()
}

func (s *AppSuite) TestIpNotAtBlacklist() {
	expect := []*storage.IPNet{
		{
			IP:   "192.168.0.0/24",
			Mask: "255.255.225.0",
		},
		{
			IP:   "10.10.0.0/24",
			Mask: "255.255.225.0",
		},
	}
	login := "superuser"
	pass := "password"
	ip := "192.168.23.1"

	s.mockKeeper.EXPECT().GetWhitelist().Return(nil)
	s.mockKeeper.EXPECT().GetBlacklist().Return(expect)
	s.mockLimiter.EXPECT().CheckLimits(s.ctx, login, pass, ip).Return(true, nil)

	result, err := s.application.CheckRate(s.ctx, login, pass, ip)
	s.Require().NoError(err)
	s.Require().True(result)
}

func (s *AppSuite) TestIpAtBlacklist() {
	expect := []*storage.IPNet{
		{
			IP:   "192.168.0.0",
			Mask: "255.0.0.0",
		},
		{
			IP:   "10.10.0.0",
			Mask: "255.255.225.0",
		},
	}
	s.mockKeeper.EXPECT().GetWhitelist().Return(nil)
	s.mockKeeper.EXPECT().GetBlacklist().Return(expect)

	result, err := s.application.CheckRate(s.ctx, "user", "pass", "192.168.23.1")
	s.Require().NoError(err)
	s.Require().False(result)
}

func (s *AppSuite) TestIpNotAtWhitelist() {
	expect := []*storage.IPNet{
		{
			IP:   "192.168.0.0",
			Mask: "255.255.225.0",
		},
		{
			IP:   "10.10.0.0",
			Mask: "255.255.225.0",
		},
	}
	login := "USER"
	pass := "PASS"
	ip := "127.0.0.1"
	s.mockKeeper.EXPECT().GetWhitelist().Return(expect)
	s.mockKeeper.EXPECT().GetBlacklist().Return(nil)
	s.mockLimiter.EXPECT().CheckLimits(s.ctx, login, pass, ip).Return(true, nil)

	result, err := s.application.CheckRate(s.ctx, login, pass, ip)
	s.Require().NoError(err)
	s.Require().True(result)
}

func (s *AppSuite) TestIpAtWhitelist() {
	expect := []*storage.IPNet{
		{
			IP:   "192.168.0.0",
			Mask: "255.0.0.0",
		},
		{
			IP:   "10.10.0.0",
			Mask: "255.255.225.0",
		},
	}
	s.mockKeeper.EXPECT().GetWhitelist().Return(expect)

	result, err := s.application.CheckRate(s.ctx, "what-ever", "dont-care", "192.168.23.1")
	s.Require().NoError(err)
	s.Require().True(result)
}

func (s *AppSuite) TestAddBL() {
	ip := "192.168.1.0"
	mask := "255.255.255.0"
	s.mockStorage.EXPECT().AddToBlackList(s.ctx, ip, mask).Return(nil)

	err := s.application.AddToBlackList(s.ctx, ip, mask)
	s.Require().NoError(err)
}

func (s *AppSuite) TestAddWL() {
	ip := "192.168.1.0"
	mask := "255.255.255.0"
	s.mockStorage.EXPECT().AddToWhiteList(s.ctx, ip, mask).Return(nil)

	err := s.application.AddToWhiteList(s.ctx, ip, mask)
	s.Require().NoError(err)
}

func (s *AppSuite) TestRemoveBL() {
	s.mockStorage.EXPECT().RemoveFromBlackList(s.ctx, "ip", "mask").Return(nil)
	err := s.application.RemoveFromBlackList(s.ctx, "ip", "mask")
	s.Require().NoError(err)
}

func (s *AppSuite) TestRemoveWL() {
	s.mockStorage.EXPECT().RemoveFromWhiteList(s.ctx, "ip", "mask").Return(nil)
	err := s.application.RemoveFromWhiteList(s.ctx, "ip", "mask")
	s.Require().NoError(err)
}

func (s *AppSuite) TestDropStat() {
	login := "user"
	pass := "pass"
	s.mockLimiter.EXPECT().DropBuckets(s.ctx, login, pass).Return(nil)

	err := s.application.DropStat(s.ctx, login, pass)
	s.Require().NoError(err)
}

func TestAppSuite(t *testing.T) {
	suite.Run(t, new(AppSuite))
}
