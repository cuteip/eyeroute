package main

import (
	"fmt"
	"os"

	"github.com/cuteip/eyeroute/cmd/eyeroute-cli/mtr"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "eyeroute-cli",
		Short: "WIP: Looking Glass Server, etc ...",
	}

	rootCmd.AddCommand(mtr.RootCmd())

	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err)
		os.Exit(1)
	}
}
