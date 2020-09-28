package main

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	host        string
	port        int
	isBlacklist bool
	isWhitelist bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:  "cli",
		Long: "CLI interface for anti-brute-force service",
	}
	rootCmd.PersistentFlags().StringVar(&host, "host", "localhost", "gRPC server host (default 'localhost')")
	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", 50051, "gRPC server port (default '50051')")

	rootCmd.AddCommand(NewAddCmd(), NewRemoveCmd(), NewDropCmd())
	if rootCmd.Execute() != nil {
		os.Exit(1)
	}
}
