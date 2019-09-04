package health

import (
	"time"
)

type Metric struct {
	name  string
	value map[string]interface{}
	tm    time.Time
}

func NewMetric(server string, val map[string]interface{}, tm time.Time) Metric {
	m := Metric{
		name:  server,
		value: nil,
		tm:    tm,
	}
	m.value = make(map[string]interface{})
	for k, v := range val {
		m.value[k] = v
	}
	return m
}

type Collection struct {
	metrics chan<- Metric
}

func NewCollection(metrics chan<- Metric) *Collection {
	return &Collection{metrics: metrics}
}

func (c *Collection) AddStatus(server string, val map[string]interface{}, tm ...time.Time) {
	c.addFields(server, val, tm...)
}

func (c *Collection) AddError(server string, err error) {
}

func (c *Collection) addFields(server string, val map[string]interface{}, tm ...time.Time) {
	c.metrics <- NewMetric(server, val, c.getTime(tm))
}

func (c *Collection) getTime(t []time.Time) time.Time {
	var timestamp time.Time
	if len(t) > 0 {
		timestamp = t[0]
	} else {
		timestamp = time.Now()
	}
	return timestamp
}

// HealthCheckor
type HealthCheckor struct {
	inputs  []Input
	outputs []Output
}

// NewHealthCheckor
func NewHealthCheckor(inputs []Input, outputs []Output) *HealthCheckor {
	return &HealthCheckor{inputs: inputs, outputs: outputs}
}
