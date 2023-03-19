package run

import (
	_ "embed"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/netip"
	"net/url"

	"github.com/cuteip/eyeroute/front"
	"github.com/cuteip/eyeroute/gen/eyeroute/mtr/v1alpha1/mtrv1alpha1connect"
	"github.com/cuteip/eyeroute/interfaces/connecthandler"
	"github.com/cuteip/eyeroute/mtr"
	"github.com/spf13/cobra"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

//go:embed mtr_stdout_dummy.json
var dummyJSONBytes []byte

func RootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "run",
		Short: "Run server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd, args)
		},
	}

	rootCmd.Flags().Bool("dummy", false, "enable dummy response")
	rootCmd.Flags().String("dev-upstream-front-url", "", "for develop: upstream url (default: serve static files in Go embed). (ex: http://127.0.0.1:3000)")

	return rootCmd
}

type mtrExecuterDummy struct{}

func (e mtrExecuterDummy) Execute(host netip.Addr, count int) ([]byte, error) {
	return dummyJSONBytes, nil
}

func run(cmd *cobra.Command, args []string) error {
	isDummy, err := cmd.Flags().GetBool("dummy")
	if err != nil {
		return err
	}

	var mtrServer *connecthandler.MtrServer
	if isDummy {
		mtrServer = connecthandler.NewMtrServer(*mtr.New(&mtrExecuterDummy{}))
	} else {
		mtrServer = connecthandler.NewMtrServer(*mtr.New(mtr.NewExecuter()))
	}

	api := http.NewServeMux()
	api.Handle(mtrv1alpha1connect.NewMtrServiceHandler(mtrServer))

	mux := http.NewServeMux()
	mux.Handle(fmt.Sprintf("/%s/", mtrv1alpha1connect.MtrServiceName), api)

	// front/ を編集するときはそっちにプロキシしたほうがやりやすいので、そのための機能
	upstreamFrontURLStr, err := cmd.Flags().GetString("dev-upstream-front-url")
	if err != nil {
		return err
	}
	if upstreamFrontURLStr != "" {
		upstreamFrontURL, err := url.Parse(upstreamFrontURLStr)
		if err != nil {
			return err
		}

		mux.Handle("/", httputil.NewSingleHostReverseProxy(upstreamFrontURL))
	} else {
		buildFS, err := front.GetBuildFS()
		if err != nil {
			return err
		}
		mux.Handle("/", http.FileServer(buildFS))
	}

	http.ListenAndServe(
		"127.0.0.1:8080",
		h2c.NewHandler(mux, &http2.Server{}),
	)

	return nil
}
