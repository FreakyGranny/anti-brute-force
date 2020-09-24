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
	mockStorage    *mocks.MockStorage
}

func (s *AppSuite) SetupTest() {
	s.mockStorageCtl = gomock.NewController(s.T())
	s.mockStorage = mocks.NewMockStorage(s.mockStorageCtl)
}

func (s *AppSuite) TearDownTest() {
	s.mockStorageCtl.Finish()
}

func (s *AppSuite) TestNotInBL() {
	ctx := context.Background()
	a := App{
		storage: s.mockStorage,
		limiter: nil,
	}
	expect := []*storage.IPNet{
		{
			ID:     1,
			Subnet: "192.168.0.0/24",
		},
		{
			ID:     2,
			Subnet: "10.10.0.0/24",
		},
	}
	s.mockStorage.EXPECT().GetBlackList(ctx).Return(expect, nil)

	c, err := a.CheckBlackList(ctx, "127.0.0.1")
	s.Require().NoError(err)
	s.Require().False(c)
}

func (s *AppSuite) TestInBL() {
	ctx := context.Background()
	a := App{
		storage: s.mockStorage,
		limiter: nil,
	}
	expect := []*storage.IPNet{
		{
			ID:     1,
			Subnet: "192.168.0.0/8",
		},
		{
			ID:     2,
			Subnet: "10.10.0.0/24",
		},
	}
	s.mockStorage.EXPECT().GetBlackList(ctx).Return(expect, nil)

	c, err := a.CheckBlackList(ctx, "192.168.23.1")
	s.Require().NoError(err)
	s.Require().True(c)
}

func (s *AppSuite) TestNotInWL() {
	ctx := context.Background()
	a := App{
		storage: s.mockStorage,
		limiter: nil,
	}
	expect := []*storage.IPNet{
		{
			ID:     1,
			Subnet: "192.168.0.0/24",
		},
		{
			ID:     2,
			Subnet: "10.10.0.0/24",
		},
	}
	s.mockStorage.EXPECT().GetWhiteList(ctx).Return(expect, nil)

	c, err := a.CheckWhiteList(ctx, "127.0.0.1")
	s.Require().NoError(err)
	s.Require().False(c)
}

func (s *AppSuite) TestInWL() {
	ctx := context.Background()
	a := App{
		storage: s.mockStorage,
		limiter: nil,
	}
	expect := []*storage.IPNet{
		{
			ID:     1,
			Subnet: "192.168.0.0/8",
		},
		{
			ID:     2,
			Subnet: "10.10.0.0/24",
		},
	}
	s.mockStorage.EXPECT().GetWhiteList(ctx).Return(expect, nil)

	c, err := a.CheckWhiteList(ctx, "192.168.23.1")
	s.Require().NoError(err)
	s.Require().True(c)
}

func (s *AppSuite) TestAddBL() {
	ctx := context.Background()
	a := App{
		storage: s.mockStorage,
		limiter: nil,
	}
	rec := &storage.IPNet{Subnet: "192.168.1.0/24"}
	s.mockStorage.EXPECT().AddToBlackList(ctx, rec).Return(nil)

	err := a.AddToBlackList(ctx, "192.168.1.0", "255.255.255.0")
	s.Require().NoError(err)
}

func (s *AppSuite) TestAddWL() {
	ctx := context.Background()
	a := App{
		storage: s.mockStorage,
		limiter: nil,
	}
	rec := &storage.IPNet{Subnet: "192.168.1.0/24"}
	s.mockStorage.EXPECT().AddToWhiteList(ctx, rec).Return(nil)

	err := a.AddToWhiteList(ctx, "192.168.1.0", "255.255.255.0")
	s.Require().NoError(err)
}

func (s *AppSuite) TestRemoveBL() {
	ctx := context.Background()
	a := App{
		storage: s.mockStorage,
		limiter: nil,
	}
	s.mockStorage.EXPECT().RemoveFromBlackList(ctx, 1).Return(nil)

	err := a.RemoveFromBlackList(ctx, 1)
	s.Require().NoError(err)
}

func (s *AppSuite) TestRemoveWL() {
	ctx := context.Background()
	a := App{
		storage: s.mockStorage,
		limiter: nil,
	}
	s.mockStorage.EXPECT().RemoveFromWhiteList(ctx, 1).Return(nil)

	err := a.RemoveFromWhiteList(ctx, 1)
	s.Require().NoError(err)
}

func TestAppSuite(t *testing.T) {
	suite.Run(t, new(AppSuite))
}
