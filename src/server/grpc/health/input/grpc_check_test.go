package input

import (
	"fmt"
	"mo-gateway/src/server/grpc/health"
	"testing"
)

func TestNewClient(t *testing.T) {
	cli := NewClient(":10081", "health-check", 10)
	err := cli.Init()
	if err != nil {
		t.Error(err)
		return
	}
}

func TestClient_Gather(t *testing.T) {
	cli := NewClient(":10081", "health-check", 10)
	err := cli.Init()
	if err != nil {
		t.Error(err)
		return
	}

	inputM := make(chan health.Metric, 100)
	acc := health.NewAccumulator(inputM)
	err = cli.Gather(acc)
	if err != nil {
		t.Error(err)
		return
	}
	out := <-inputM
	fmt.Println("server:", out.name, " value:", out.value, "time:", out.tm)

}
