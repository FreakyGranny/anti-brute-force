// +build integration

package integration_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/FreakyGranny/anti-brute-force/internal/server"
	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

type IntegrationSuite struct {
	suite.Suite
	client server.ABruteforceClient
}

func (s *IntegrationSuite) SetupTest() {
	grpcURL := os.Getenv("GRPC_URL")
	if grpcURL == "" {
		grpcURL = "127.0.0.1:50051"
	}
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	conn, err := grpc.Dial(grpcURL, opts...)
	if err != nil {
		s.T().Fail()
	}
	s.client = server.NewABruteforceClient(conn)

	if time.Now().Second() >= 50 {
		// Waiting for a new minute
		time.Sleep(10)
	}
}

func (s *IntegrationSuite) TestAuthIPInvalid() {
	_, err := s.client.Auth(context.Background(), &server.AuthRequest{
		Login:    faker.Username(),
		Password: faker.Password(),
		Ip:       faker.Word(),
	})
	s.Require().Error(err)
}

func (s *IntegrationSuite) TestIpAtBlacklist() {
	login := faker.Username()
	pass := faker.Password()
	ip := "10.0.0.1"
	mask := "255.255.255.0"
	s.AddToBlacklist(ip, mask)
	time.Sleep(2 * time.Second) // Wait for app updates blacklist
	ok := s.CheckAuth(login, pass, "10.0.0.29")
	s.Require().False(ok)

	s.RemoveFromBlacklist(ip, mask)
	time.Sleep(2 * time.Second) // Wait for app updates blacklist

	ok = s.CheckAuth(login, pass, "10.0.0.29")
	s.Require().True(ok)
}

func (s *IntegrationSuite) CheckAuth(login, password, ip string) bool {
	s.T().Helper()
	response, err := s.client.Auth(context.Background(), &server.AuthRequest{
		Login:    login,
		Password: password,
		Ip:       ip,
	})
	s.Require().NoError(err)
	s.Require().NotNil(response)

	return response.GetOk()
}

func (s *IntegrationSuite) AddToBlacklist(ip, mask string) {
	s.T().Helper()
	response, err := s.client.AddToBlackList(context.Background(), &server.AddSubnetRequest{
		Ip:   ip,
		Mask: mask,
	})
	s.Require().NoError(err)
	s.Require().NotNil(response)
	s.Require().Equal(ip, response.GetIp())
	s.Require().Equal(mask, response.GetMask())
}

func (s *IntegrationSuite) RemoveFromBlacklist(ip, mask string) {
	s.T().Helper()
	_, err := s.client.RemoveFromBlackList(context.Background(), &server.RemoveSubnetRequest{
		Ip:   ip,
		Mask: mask,
	})
	s.Require().NoError(err)
}

func (s *IntegrationSuite) TestAddToBlacklistInvalid() {
	_, err := s.client.AddToBlackList(context.Background(), &server.AddSubnetRequest{
		Ip:   faker.Word(),
		Mask: faker.IPv4(),
	})
	s.Require().Error(err)
}

func (s *IntegrationSuite) TestRemoveFromBlacklistInvalid() {
	_, err := s.client.RemoveFromBlackList(context.Background(), &server.RemoveSubnetRequest{
		Ip:   faker.Word(),
		Mask: faker.IPv4(),
	})
	s.Require().Error(err)
}

func (s *IntegrationSuite) TestIpAtWhitelist() {
	login := faker.Username()
	pass := faker.Password()
	ip := faker.IPv4()
	ok := false
	for i := 0; i < 5; i++ {
		ok = s.CheckAuth(login, pass, ip)
		s.Require().True(ok)
	}
	ok = s.CheckAuth(login, pass, ip)
	s.Require().False(ok)

	ip = "192.168.0.1"
	mask := "255.255.0.0"
	s.AddToWhitelist(ip, mask)
	time.Sleep(2 * time.Second) // Wait for app updates whitelist

	ok = s.CheckAuth(login, pass, "192.168.9.29")
	s.Require().True(ok)

	s.RemoveFromWhitelist(ip, mask)
	time.Sleep(2 * time.Second) // Wait for app updates whitelist
}

func (s *IntegrationSuite) AddToWhitelist(ip, mask string) {
	s.T().Helper()
	response, err := s.client.AddToWhiteList(context.Background(), &server.AddSubnetRequest{
		Ip:   ip,
		Mask: mask,
	})
	s.Require().NoError(err)
	s.Require().NotNil(response)
	s.Require().Equal(ip, response.GetIp())
	s.Require().Equal(mask, response.GetMask())
}

func (s *IntegrationSuite) RemoveFromWhitelist(ip, mask string) {
	s.T().Helper()
	_, err := s.client.RemoveFromWhiteList(context.Background(), &server.RemoveSubnetRequest{
		Ip:   ip,
		Mask: mask,
	})
	s.Require().NoError(err)
}

func (s *IntegrationSuite) TestAddToWhitelistInvalid() {
	_, err := s.client.AddToWhiteList(context.Background(), &server.AddSubnetRequest{
		Ip:   faker.Word(),
		Mask: faker.IPv4(),
	})
	s.Require().Error(err)
}

func (s *IntegrationSuite) TestRemoveFromWhitelistInvalid() {
	_, err := s.client.RemoveFromWhiteList(context.Background(), &server.RemoveSubnetRequest{
		Ip:   faker.Word(),
		Mask: faker.IPv4(),
	})
	s.Require().Error(err)
}

func (s *IntegrationSuite) TestUserLimit() {
	login := faker.Username()
	pass := faker.Password()
	ip := faker.IPv4()
	ok := false
	for i := 0; i < 5; i++ {
		ok = s.CheckAuth(login, pass, ip)
		s.Require().True(ok)
	}
	ok = s.CheckAuth(login, pass, ip)
	s.Require().False(ok)
}

func (s *IntegrationSuite) TestPasswordLimit() {
	pass := faker.Password()
	ip := faker.IPv4()
	ok := false
	for i := 0; i < 10; i++ {
		ok = s.CheckAuth(faker.Username(), pass, ip)
		s.Require().True(ok)
	}
	ok = s.CheckAuth(faker.Username(), pass, ip)
	s.Require().False(ok)
}

func (s *IntegrationSuite) TestIpLimit() {
	ip := faker.IPv4()
	ok := false
	for i := 0; i < 20; i++ {
		ok = s.CheckAuth(faker.Username(), faker.Password(), ip)
		s.Require().True(ok)
	}
	ok = s.CheckAuth(faker.Username(), faker.Password(), ip)
	s.Require().False(ok)
}

func (s *IntegrationSuite) TestDropStat() {
	login := faker.Username()
	pass := faker.Password()
	ip := faker.IPv4()
	ok := false
	for i := 0; i < 5; i++ {
		ok = s.CheckAuth(login, pass, ip)
		s.Require().True(ok)
	}
	ok = s.CheckAuth(login, pass, ip)
	s.Require().False(ok)

	_, err := s.client.DropStat(context.Background(), &server.DropStatRequest{
		Login: login,
		Ip:    ip,
	})
	s.Require().NoError(err)

	ok = s.CheckAuth(login, pass, ip)
	s.Require().True(ok)
}

func TestIntegration(t *testing.T) {
	suite.Run(t, new(IntegrationSuite))
}
