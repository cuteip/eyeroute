package main

import (
	"fmt"
	"os"

	"github.com/cuteip/eyeroute/cmd/eyeroute/server"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "eyeroute",
		Short: "WIP: Looking Glass Server, etc ...",
	}

	rootCmd.AddCommand(server.RootCmd())

	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err)
		os.Exit(1)
	}
}
