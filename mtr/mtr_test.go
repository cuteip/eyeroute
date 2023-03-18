package mtr

import (
	"io"
	"net/netip"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// テスト用
type ExecuterFake struct{}

func (e *ExecuterFake) Execute(host netip.Addr, count int) ([]byte, error) {
	f, err := os.Open("../testdata/mtr/stdout1.json")
	if err != nil {
		return nil, err
	}

	stdoutBytes, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return stdoutBytes, nil
}

func TestMtr_Run(t *testing.T) {
	tests := []struct {
		name    string
		want    Report
		wantErr bool
	}{
		{
			name: "test",
			want: Report{
				Mtr: ReportMtr{
					Src:        "fake-host",
					Dst:        "1.1.1.1",
					Tos:        0,
					Tests:      10,
					Psize:      "64",
					Bitpattern: "0x00",
				},
				Hubs: []ReportHub{
					{
						Count: 1,
						Host:  "192.168.0.1",
						Loss:  0.0,
						Snt:   10,
						Last:  0.206,
						Avg:   0.212,
						Best:  0.176,
						Wrst:  0.26,
						StDev: 0.026,
					},
					{
						Count: 2,
						Host:  "192.168.1.1",
						Loss:  0.0,
						Snt:   10,
						Last:  0.392,
						Avg:   0.356,
						Best:  0.324,
						Wrst:  0.394,
						StDev: 0.027,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ExecuterFake{}
			m := &Mtr{
				Executer: e,
			}
			targetIPAddrFake := netip.MustParseAddr("192.0.2.1")
			countFake := 10
			got, err := m.Run(targetIPAddrFake, countFake)
			if (err != nil) != tt.wantErr {
				t.Errorf("Mtr.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Mtr.Run() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
