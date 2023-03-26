package mtr

import (
	"bytes"
	"encoding/json"
	"net/netip"
	"os/exec"
	"strconv"
)

type Mtr struct {
	// テスト時にモックするため
	Executer Executer
}

type Executer interface {
	// io.Writer を渡す、もしくは io.Reader を返すのが良いのかも
	Execute(host netip.Addr, count int) ([]byte, error)
}

type ExecuterReal struct{}

// --json オプションをつけたときに出力される JSON
type ReportJSON struct {
	Report Report `json:"report"`
}

type Report struct {
	Mtr  ReportMtr   `json:"mtr"`
	Hubs []ReportHub `json:"hubs"`
}

type ReportMtr struct {
	Src        string `json:"src"`
	Dst        string `json:"dst"`
	Tos        int    `json:"tos"`
	Tests      int    `json:"tests"`
	Psize      string `json:"psize"`
	Bitpattern string `json:"bitpattern"`
}

type ReportHub struct {
	Count int     `json:"count"`
	Host  string  `json:"host"`
	Loss  float64 `json:"Loss%"`
	Snt   int     `json:"Snt"`
	Last  float64 `json:"Last"`
	Avg   float64 `json:"Avg"`
	Best  float64 `json:"Best"`
	Wrst  float64 `json:"Wrst"`
	StDev float64 `json:"StDev"`
}

func New(executer Executer) *Mtr {
	return &Mtr{
		Executer: executer,
	}
}

func NewExecuter() Executer {
	e := &ExecuterReal{}
	return e
}

func (e *ExecuterReal) Execute(host netip.Addr, count int) ([]byte, error) {
	c := exec.Command("mtr", "--json", "-n", "-c", strconv.Itoa(count), host.String())

	var stdoutBuf bytes.Buffer
	c.Stdout = &stdoutBuf
	// stderr は無視する

	err := c.Run()
	return stdoutBuf.Bytes(), err
}

func (m *Mtr) Run(targetIPAddr netip.Addr, count int) (Report, error) {
	stdout, err := m.Executer.Execute(targetIPAddr, count)
	if err != nil {
		return Report{}, err
	}

	var ReportJSON ReportJSON
	err = json.Unmarshal(stdout, &ReportJSON)
	if err != nil {
		return Report{}, err
	}

	return ReportJSON.Report, nil
}
