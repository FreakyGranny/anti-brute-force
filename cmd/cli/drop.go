package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/FreakyGranny/anti-brute-force/internal/server"
	"github.com/spf13/cobra"
)

// NewDropCmd returns cmd for dropping stat fro given login and password.
func NewDropCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "drop login password",
		Short: "Drop stats for login/password",
		Long:  "Allows drop stats for given login and password",
		Args:  cobra.ExactArgs(2),
		RunE:  Drop,
	}
}

// Drop drops stats for given login/password.
func Drop(cmd *cobra.Command, args []string) error {
	client, err := getGRPCClient()
	if err != nil {
		fmt.Println("Connection refused")
		os.Exit(504)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	login := args[0]
	password := args[1]

	_, err = client.DropStat(ctx, &server.DropStatRequest{Login: login, Password: password})
	cancel()

	if err != nil {
		fmt.Printf("error while dropping stat for [login: %s password: %s]: %s\n", login, password, err.Error())
		os.Exit(13)
	}
	fmt.Printf("stats for [login: %s passrod: %s] successfully dropped\n", login, password)

	return nil
}
