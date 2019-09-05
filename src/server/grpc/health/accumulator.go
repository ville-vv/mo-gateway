package health

import "time"

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

func NewAccumulator(metrics chan<- Metric) Accumulator {
	return &Collection{metrics: metrics}
}

func (c *Collection) AddStatus(server string, val map[string]interface{}, tm ...time.Time) {
	c.addFields(server, val, tm...)
}

func (c *Collection) AddError(server string, err error) {
}

// addFields
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
