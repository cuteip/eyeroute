package mtr

import (
	"context"
	"fmt"
	"net/http"
	"net/netip"
	"os"
	"strconv"

	"github.com/bufbuild/connect-go"
	mtrv1alpha1 "github.com/cuteip/eyeroute/gen/eyeroute/mtr/v1alpha1"
	"github.com/cuteip/eyeroute/gen/eyeroute/mtr/v1alpha1/mtrv1alpha1connect"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "mtr [IP Address]",
		Short: "execute mtr",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd, args)
		},
	}

	rootCmd.Flags().String("server", "", "eyeroute server (ex: https://eyeroute.example.com)")
	rootCmd.MarkFlagRequired("server")

	return rootCmd
}

func run(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("ip address is required")
	}

	ipAddr, err := netip.ParseAddr(args[0])
	if err != nil {
		return errors.Wrap(err, "failed to parse ip address")
	}

	server, err := cmd.Flags().GetString("server")
	if err != nil {
		return err
	}

	client := mtrv1alpha1connect.NewMtrServiceClient(
		http.DefaultClient,
		server,
	)

	req := connect.NewRequest(&mtrv1alpha1.ExecuteMtrRequest{
		IpAddress:    ipAddr.String(),
		ReportCycles: 10,
	})

	res, err := client.ExecuteMtr(context.Background(), req)
	if err != nil {
		return err
	}

	tableList := hubsToTableList(res.Msg)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Hop", "Host", "Sent", "Last", "Avg", "Best", "Worst", "StDev"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.AppendBulk(tableList)
	table.Render()

	return nil
}

func hubsToTableList(hubs *mtrv1alpha1.ExecuteMtrResponse) [][]string {
	table := [][]string{}
	for index, hub := range hubs.Hubs {
		line := []string{strconv.Itoa(index + 1), hub.Host, strconv.Itoa(int(hub.Sent)), fmt.Sprintf("%.1f", hub.Last), fmt.Sprintf("%.1f", hub.Avg), fmt.Sprintf("%.1f", hub.Best), fmt.Sprintf("%.1f", hub.Worst), fmt.Sprintf("%.1f", hub.Stdev)}
		table = append(table, line)
	}

	return table
}
