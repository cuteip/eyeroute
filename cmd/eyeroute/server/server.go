package server

import (
	"github.com/cuteip/eyeroute/cmd/eyeroute/server/run"
	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "server",
		Short: "server",
	}

	rootCmd.AddCommand(run.RootCmd())
	return rootCmd
}
