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

// NewAddCmd returns cmd for adding subnets to black/white list.
func NewAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add 127.0.0.1 255.255.255.0",
		Short: "Add to black/white lists",
		Long:  "Allows add subnets to black/white lists",
		Args:  cobra.ExactArgs(2),
		RunE:  Add,
	}
	cmd.Flags().BoolVarP(&isBlacklist, "blacklist", "b", false, "adds to blacklist")
	cmd.Flags().BoolVarP(&isWhitelist, "whitelist", "w", false, "adds to whitelist")

	return cmd
}

// Add adds subnet to black/white list.
func Add(cmd *cobra.Command, args []string) error {
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

	request := &server.AddSubnetRequest{Ip: ip, Mask: mask}
	response := &server.AddSubnetResponse{}
	if isBlacklist {
		response, err = client.AddToBlackList(ctx, request)
	}
	if isWhitelist {
		response, err = client.AddToWhiteList(ctx, request)
	}
	cancel()

	if err != nil {
		fmt.Printf("error while adding subnet [ip: %s mask: %s]: %s\n", ip, mask, err.Error())
		os.Exit(13)
	}
	fmt.Printf("subnet [ip: %s mask: %s] successfully added\n", response.GetIp(), response.GetMask())

	return nil
}
