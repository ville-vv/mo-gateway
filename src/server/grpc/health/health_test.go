package health

import (
	"context"
	"fmt"
	"mo-gateway/src/server/grpc/health/input"
	"testing"
	"time"
)

type inputMock struct {
}

func (*inputMock) Init() error {
	fmt.Println("插件初始化 input")
	return nil
}

func (*inputMock) Gather(acc Accumulator) error {
	acc.AddStatus("test", map[string]interface{}{"a": "b"})
	return nil
}

func (*inputMock) Stop() {
	fmt.Println("server input 插件关闭")
}

func (*inputMock) Start(acc Accumulator) error {
	fmt.Println("插件 server input 启动")
	return nil
}

func (*inputMock) Label() string {
	return ""
}

type outputMock struct {
}

func (*outputMock) Init() error {
	fmt.Println("插件初始化 output")
	return nil
}

func (*outputMock) Connect() error {
	return nil
}

func (*outputMock) Close() {
	return
}

func (*outputMock) Report(m Metric) {
	fmt.Println(m.value)
}

func TestServerCheckor_Run(t *testing.T) {
	check := NewServerCheckor([]Input{&inputMock{}}, []Output{&outputMock{}}, 2*time.Second)
	if err := check.Run(context.Background()); err != nil {
		t.Error(err)
	}
	time.Sleep(time.Second * 5)
	check.Stop()
	time.Sleep(time.Second * 5)
}

func TestServerCheckor_Run2(t *testing.T) {

	inputs := []Input{
		&inputMock{},
		input.NewClient(":10081", "health-check-rpc", 10),
	}

	check := NewServerCheckor(inputs, []Output{&outputMock{}}, 2*time.Second)
	if err := check.Run(context.Background()); err != nil {
		t.Error(err)
	}
	time.Sleep(time.Second * 30)
	check.Stop()
	time.Sleep(time.Second * 5)
}
