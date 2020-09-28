package main

import (
	"errors"
	"net"
	"strconv"

	"github.com/FreakyGranny/anti-brute-force/internal/server"
	"google.golang.org/grpc"
)

var (
	errListNotChoosen = errors.New("one of (blacklist/whitelist) flag must be set")
	errIPInvalid      = errors.New("ip address is invalid")
	errMaskInvalid    = errors.New("mask is invalid")
)

func getGRPCClient() (server.ABruteforceClient, error) {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	conn, err := grpc.Dial(net.JoinHostPort(host, strconv.Itoa(port)), opts...)
	if err != nil {
		return nil, err
	}

	return server.NewABruteforceClient(conn), nil
}

func isValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}
