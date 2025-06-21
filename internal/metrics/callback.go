package metrics

import "github.com/prometheus/client_golang/prometheus"

type CallbackGauge struct {
	callback func() float64
	gauge    prometheus.Gauge
}

func NewCallbackGauge(opts prometheus.GaugeOpts) *CallbackGauge {
	return &CallbackGauge{
		callback: func() float64 {
			return 0
		},
		gauge: prometheus.NewGauge(opts),
	}
}

func (c *CallbackGauge) Bind(callback func() float64) {
	c.callback = callback
}

func (c *CallbackGauge) Describe(ch chan<- *prometheus.Desc) {
	c.gauge.Describe(ch)
}

func (c *CallbackGauge) Collect(ch chan<- prometheus.Metric) {
	result := c.callback()

	c.gauge.Set(result)

	c.gauge.Collect(ch)
}
