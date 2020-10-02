package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/FreakyGranny/anti-brute-force/internal/app"
	"github.com/FreakyGranny/anti-brute-force/internal/server"
	"github.com/spf13/cobra"
)

// NewDropCmd returns cmd for dropping stat fro given login and ip.
func NewDropCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "drop login ip",
		Short: "Drop stats for login,ip",
		Long:  "Allows drop stats for given login and ip",
		Args:  cobra.ExactArgs(2),
		RunE:  Drop,
	}
}

// Drop drops stats for given login, ip.
func Drop(cmd *cobra.Command, args []string) error {
	client, err := getGRPCClient()
	if err != nil {
		fmt.Println("Connection refused")
		os.Exit(504)
	}

	login := args[0]
	ip := args[1]
	if !app.IsValidIPFormat(ip) {
		return errIPInvalid
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = client.DropStat(ctx, &server.DropStatRequest{Login: login, Ip: ip})
	cancel()

	if err != nil {
		fmt.Printf("error while dropping stat for [login: %s ip: %s]: %s\n", login, ip, err.Error())
		os.Exit(13)
	}
	fmt.Printf("stats for [login: %s ip: %s] successfully dropped\n", login, ip)

	return nil
}
