package mtr

import (
	"bytes"
	"encoding/json"
	"net/netip"
	"os/exec"
	"strconv"
	"strings"
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

// mtr 0.94 では 0 のように 10 進数の int で出力される
// mtr 0.93 では 0x0 のように 16 進数の string で出力される
// UnmarshalJSON() を細工するため、独自に型を定義する
type reportMtrTos int

func (t *reportMtrTos) UnmarshalJSON(b []byte) error {
	// まず、int としてパースしてみる
	var tosInt int
	err := json.Unmarshal(b, &tosInt)
	if err == nil {
		// パースに成功したので終了
		*t = reportMtrTos(tosInt)
		return nil
	}

	// 素直に int にパースできなかったので 16 進数として扱いパースを試みる
	tosString := strings.ReplaceAll(string(b), `"`, "")
	tosInt64, err := strconv.ParseInt(tosString, 0, 64)
	if err != nil {
		return err
	}

	*t = reportMtrTos(tosInt64)
	return nil
}

type ReportMtr struct {
	Src        string       `json:"src"`
	Dst        string       `json:"dst"`
	Tos        reportMtrTos `json:"tos"`
	Tests      int          `json:"tests"`
	Psize      string       `json:"psize"`
	Bitpattern string       `json:"bitpattern"`
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
