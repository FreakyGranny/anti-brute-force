package app

import (
	"context"
	"testing"

	"github.com/FreakyGranny/anti-brute-force/internal/mocks"
	"github.com/FreakyGranny/anti-brute-force/internal/storage"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type MemIPKeeperSuite struct {
	suite.Suite
	mockStorageCtl *gomock.Controller
	mockStorage    *mocks.MockReadStorage
	keeper         *MemIPKeeper
	ctx            context.Context
}

func (s *MemIPKeeperSuite) SetupTest() {
	s.mockStorageCtl = gomock.NewController(s.T())
	s.mockStorage = mocks.NewMockReadStorage(s.mockStorageCtl)
	s.keeper = NewMemIPKeeper(s.mockStorage)
	s.ctx = context.Background()
}

func (s *MemIPKeeperSuite) TearDownTest() {
	s.mockStorageCtl.Finish()
}

func (s *MemIPKeeperSuite) TestWithoutRefresh() {
	res := s.keeper.GetBlacklist()
	s.Require().Empty(res)
	res = s.keeper.GetWhitelist()
	s.Require().Empty(res)
}

func (s *MemIPKeeperSuite) TestEmptyLists() {
	s.mockStorage.EXPECT().GetBlackList(s.ctx).Return(nil, nil)
	s.mockStorage.EXPECT().GetWhiteList(s.ctx).Return(nil, nil)

	err := s.keeper.Refresh(s.ctx)
	s.Require().NoError(err)
	res := s.keeper.GetBlacklist()
	s.Require().Empty(res)
	res = s.keeper.GetWhitelist()
	s.Require().Empty(res)
}

func (s *MemIPKeeperSuite) TestWhiteList() {
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

	s.mockStorage.EXPECT().GetBlackList(s.ctx).Return(nil, nil)
	s.mockStorage.EXPECT().GetWhiteList(s.ctx).Return(expect, nil)

	err := s.keeper.Refresh(s.ctx)
	s.Require().NoError(err)
	res := s.keeper.GetBlacklist()
	s.Require().Empty(res)
	res = s.keeper.GetWhitelist()
	s.Require().Equal(expect, res)
}

func (s *MemIPKeeperSuite) TestBlackList() {
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

	s.mockStorage.EXPECT().GetBlackList(s.ctx).Return(expect, nil)
	s.mockStorage.EXPECT().GetWhiteList(s.ctx).Return(nil, nil)

	err := s.keeper.Refresh(s.ctx)
	s.Require().NoError(err)
	res := s.keeper.GetBlacklist()
	s.Require().Equal(expect, res)
	res = s.keeper.GetWhitelist()
	s.Require().Empty(res)
}

func TestMemIpKeeperSuite(t *testing.T) {
	suite.Run(t, new(MemIPKeeperSuite))
}
