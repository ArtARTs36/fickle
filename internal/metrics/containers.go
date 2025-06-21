package metrics

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type Containers struct {
	running *CallbackGauge
	runs    *prometheus.CounterVec
	stops   *prometheus.CounterVec
}

func NewContainers(namespace string) *Containers {
	return &Containers{
		running: NewCallbackGauge(prometheus.GaugeOpts{
			Name:      "containers_running",
			Namespace: namespace,
			Help:      "Containers: count of running containers",
		}),
		runs: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:      "containers_runs_total",
			Namespace: namespace,
			Help:      "Containers: count of runs",
		}, []string{"host", "state"}),
		stops: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:      "containers_stops_total",
			Namespace: namespace,
			Help:      "Containers: count of stops",
		}, []string{"host", "state"}),
	}
}

func (c *Containers) Describe(ch chan<- *prometheus.Desc) {
	c.running.Describe(ch)
	c.runs.Describe(ch)
	c.stops.Describe(ch)
}

func (c *Containers) Collect(ch chan<- prometheus.Metric) {
	c.running.Collect(ch)
	c.runs.Collect(ch)
	c.stops.Collect(ch)
}

func (c *Containers) BindRunningCallback(callback func() float64) {
	c.running.Bind(callback)
}

func (c *Containers) IncRun(host string, state bool) {
	c.runs.WithLabelValues(host, strconv.FormatBool(state)).Inc()
}

func (c *Containers) IncStops(host string, state bool) {
	c.stops.WithLabelValues(host, strconv.FormatBool(state)).Inc()
}
