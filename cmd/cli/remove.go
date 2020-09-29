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

// NewRemoveCmd returns cmd for removing subnets from black/white list.
func NewRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove 127.0.0.1 255.255.255.0",
		Short: "Remove from black/white lists",
		Long:  "Allows remove subnets from black/white lists",
		Args:  cobra.ExactArgs(2),
		RunE:  Remove,
	}
	cmd.Flags().BoolVarP(&isBlacklist, "blacklist", "b", false, "adds to blacklist")
	cmd.Flags().BoolVarP(&isWhitelist, "whitelist", "w", false, "adds to whitelist")

	return cmd
}

// Remove removes subnet from black/white list.
func Remove(cmd *cobra.Command, args []string) error {
	if !isBlacklist && !isWhitelist {
		return errListNotChoosen
	}
	ip := args[0]
	mask := args[1]
	if !app.IsValidIPFormat(ip) {
		return errIPInvalid
	}
	if !app.IsValidIPFormat(mask) {
		return errMaskInvalid
	}
	client, err := getGRPCClient()
	if err != nil {
		fmt.Println("Connection refused")
		os.Exit(504)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	request := &server.RemoveSubnetRequest{Ip: ip, Mask: mask}
	if isBlacklist {
		_, err = client.RemoveFromBlackList(ctx, request)
	}
	if isWhitelist {
		_, err = client.RemoveFromWhiteList(ctx, request)
	}
	cancel()

	if err != nil {
		fmt.Printf("error while removing subnet [ip: %s mask: %s]: %s\n", ip, mask, err.Error())
		os.Exit(13)
	}
	fmt.Printf("subnet [ip: %s mask: %s] successfully removed\n", ip, mask)

	return nil
}
