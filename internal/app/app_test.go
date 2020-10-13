package app

import (
	"context"
	"net"
	"testing"

	"github.com/FreakyGranny/anti-brute-force/internal/mocks"
	"github.com/bxcodec/faker/v3"
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

func (s *AppSuite) TestCheckInvalidIP() {
	_, err := s.application.CheckRate(s.ctx, faker.Username(), faker.Password(), faker.Word())
	s.Require().Error(err)
	s.Require().Equal(ErrInvalidArgument, err)
}

func (s *AppSuite) TestIpNotAtBlacklist() {
	expect := []*net.IPNet{
		{
			IP:   net.IPv4(192, 168, 0, 0),
			Mask: net.IPv4Mask(255, 255, 225, 0),
		},
		{
			IP:   net.IPv4(10, 10, 0, 0),
			Mask: net.IPv4Mask(255, 255, 225, 0),
		},
	}
	login := faker.Username()
	pass := faker.Password()
	ip := "192.168.23.1"

	s.mockKeeper.EXPECT().GetWhitelist().Return(nil)
	s.mockKeeper.EXPECT().GetBlacklist().Return(expect)
	s.mockLimiter.EXPECT().CheckLimits(s.ctx, login, pass, ip).Return(true, nil)

	result, err := s.application.CheckRate(s.ctx, login, pass, ip)
	s.Require().NoError(err)
	s.Require().True(result)
}

func (s *AppSuite) TestIpAtBlacklist() {
	expect := []*net.IPNet{
		{
			IP:   net.IPv4(192, 168, 0, 0),
			Mask: net.IPv4Mask(255, 0, 0, 0),
		},
		{
			IP:   net.IPv4(10, 10, 0, 0),
			Mask: net.IPv4Mask(255, 255, 225, 0),
		},
	}
	s.mockKeeper.EXPECT().GetWhitelist().Return(nil)
	s.mockKeeper.EXPECT().GetBlacklist().Return(expect)

	result, err := s.application.CheckRate(s.ctx, faker.Username(), faker.Password(), "192.168.23.1")
	s.Require().NoError(err)
	s.Require().False(result)
}

func (s *AppSuite) TestIpNotAtWhitelist() {
	expect := []*net.IPNet{
		{
			IP:   net.IPv4(17, 18, 0, 0),
			Mask: net.IPv4Mask(255, 255, 225, 0),
		},
		{
			IP:   net.IPv4(123, 15, 89, 0),
			Mask: net.IPv4Mask(255, 255, 225, 0),
		},
	}
	login := faker.Username()
	pass := faker.Password()
	ip := "127.0.0.1"
	s.mockKeeper.EXPECT().GetWhitelist().Return(expect)
	s.mockKeeper.EXPECT().GetBlacklist().Return(nil)
	s.mockLimiter.EXPECT().CheckLimits(s.ctx, login, pass, ip).Return(true, nil)

	result, err := s.application.CheckRate(s.ctx, login, pass, ip)
	s.Require().NoError(err)
	s.Require().True(result)
}

func (s *AppSuite) TestIpAtWhitelist() {
	expect := []*net.IPNet{
		{
			IP:   net.IPv4(192, 168, 0, 0),
			Mask: net.IPv4Mask(255, 0, 0, 0),
		},
		{
			IP:   net.IPv4(10, 10, 0, 0),
			Mask: net.IPv4Mask(255, 255, 225, 0),
		},
	}
	s.mockKeeper.EXPECT().GetWhitelist().Return(expect)

	result, err := s.application.CheckRate(s.ctx, faker.Username(), faker.Password(), "192.168.23.1")
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

func (s *AppSuite) TestAddBLInvalidIP() {
	err := s.application.AddToBlackList(s.ctx, faker.Word(), faker.IPv4())
	s.Require().Error(err)
	s.Require().Equal(ErrInvalidArgument, err)
}

func (s *AppSuite) TestAddBLInvalidMask() {
	err := s.application.AddToBlackList(s.ctx, faker.IPv4(), faker.Word())
	s.Require().Error(err)
	s.Require().Equal(ErrInvalidArgument, err)
}

func (s *AppSuite) TestAddWL() {
	ip := "192.168.1.0"
	mask := "255.255.255.0"
	s.mockStorage.EXPECT().AddToWhiteList(s.ctx, ip, mask).Return(nil)

	err := s.application.AddToWhiteList(s.ctx, ip, mask)
	s.Require().NoError(err)
}

func (s *AppSuite) TestAddWLInvalidIP() {
	err := s.application.AddToWhiteList(s.ctx, faker.Word(), faker.IPv4())
	s.Require().Error(err)
	s.Require().Equal(ErrInvalidArgument, err)
}

func (s *AppSuite) TestAddWLInvalidMask() {
	err := s.application.AddToWhiteList(s.ctx, faker.IPv4(), faker.Word())
	s.Require().Error(err)
	s.Require().Equal(ErrInvalidArgument, err)
}

func (s *AppSuite) TestRemoveBL() {
	ip := faker.IPv4()
	mask := faker.IPv4()
	s.mockStorage.EXPECT().RemoveFromBlackList(s.ctx, ip, mask).Return(nil)
	err := s.application.RemoveFromBlackList(s.ctx, ip, mask)
	s.Require().NoError(err)
}

func (s *AppSuite) TestRemoveBLInvalidIP() {
	err := s.application.RemoveFromBlackList(s.ctx, faker.Word(), faker.IPv4())
	s.Require().Error(err)
	s.Require().Equal(ErrInvalidArgument, err)
}

func (s *AppSuite) TestRemoveBLInvalidMask() {
	err := s.application.RemoveFromBlackList(s.ctx, faker.IPv4(), faker.Word())
	s.Require().Error(err)
	s.Require().Equal(ErrInvalidArgument, err)
}

func (s *AppSuite) TestRemoveWL() {
	ip := faker.IPv4()
	mask := faker.IPv4()
	s.mockStorage.EXPECT().RemoveFromWhiteList(s.ctx, ip, mask).Return(nil)
	err := s.application.RemoveFromWhiteList(s.ctx, ip, mask)
	s.Require().NoError(err)
}

func (s *AppSuite) TestRemoveWLInvalidIP() {
	err := s.application.RemoveFromWhiteList(s.ctx, faker.Word(), faker.IPv4())
	s.Require().Error(err)
	s.Require().Equal(ErrInvalidArgument, err)
}

func (s *AppSuite) TestRemoveWLInvalidMask() {
	err := s.application.RemoveFromWhiteList(s.ctx, faker.IPv4(), faker.Word())
	s.Require().Error(err)
	s.Require().Equal(ErrInvalidArgument, err)
}

func (s *AppSuite) TestDropStat() {
	login := faker.Username()
	ip := faker.IPv4()
	s.mockLimiter.EXPECT().DropBuckets(s.ctx, login, ip).Return(nil)

	err := s.application.DropStat(s.ctx, login, ip)
	s.Require().NoError(err)
}

func (s *AppSuite) TestDropStatIpInvalid() {
	err := s.application.DropStat(s.ctx, faker.Username(), faker.Word())
	s.Require().Error(err)
	s.Require().Equal(ErrInvalidArgument, err)
}

func TestAppSuite(t *testing.T) {
	suite.Run(t, new(AppSuite))
}
