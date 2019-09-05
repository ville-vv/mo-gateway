package output

import (
	"mo-gateway/src/server/grpc/health"
)

type PrometheusPlg struct {
}

func NewPrometheusPlg() *PrometheusPlg {
	return &PrometheusPlg{}
}

func (p *PrometheusPlg) Init() error {
	return nil
}

func (p *PrometheusPlg) Connect() error {
	return nil
}

func (p *PrometheusPlg) Close() {
	return
}

func (p *PrometheusPlg) Report(health.Metric) {
	return
}
