package connecthandler

import (
	"context"
	"net/netip"

	"github.com/pkg/errors"

	"github.com/bufbuild/connect-go"
	mtrv1alpha1 "github.com/cuteip/eyeroute/gen/eyeroute/mtr/v1alpha1"
	"github.com/cuteip/eyeroute/mtr"
)

type MtrServer struct {
	Mtr mtr.Mtr
}

func NewMtrServer(m mtr.Mtr) *MtrServer {
	return &MtrServer{Mtr: m}
}

func (m MtrServer) ExecuteMtr(
	ctx context.Context,
	req *connect.Request[mtrv1alpha1.ExecuteMtrRequest],
) (*connect.Response[mtrv1alpha1.ExecuteMtrResponse], error) {
	targetIPAddr, err := netip.ParseAddr(req.Msg.IpAddress)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	const reportCyclesLimit = 10
	if req.Msg.ReportCycles > reportCyclesLimit {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.Errorf("report cycles must be <= %d", reportCyclesLimit))
	}

	report, err := m.Mtr.Run(targetIPAddr, int(req.Msg.ReportCycles))
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	var hubs []*mtrv1alpha1.ReportHub
	for _, reportHub := range report.Hubs {
		hubs = append(hubs, convertReportHubToProto(reportHub))
	}

	res := connect.NewResponse(&mtrv1alpha1.ExecuteMtrResponse{
		Hubs: hubs,
	})
	return res, nil
}

func convertReportHubToProto(hub mtr.ReportHub) *mtrv1alpha1.ReportHub {
	return &mtrv1alpha1.ReportHub{
		Count: int32(hub.Count),
		Host:  hub.Host,
		Loss:  float32(hub.Loss),
		Sent:  int32(hub.Snt),
		Last:  float32(hub.Last),
		Avg:   float32(hub.Avg),
		Best:  float32(hub.Best),
		Worst: float32(hub.Wrst),
		Stdev: float32(hub.StDev),
	}
}
