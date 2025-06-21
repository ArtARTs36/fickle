package metrics

import "github.com/prometheus/client_golang/prometheus"

type Group struct {
	Containers *Containers
}

func NewGroup(namespace string) *Group {
	return &Group{
		Containers: NewContainers(namespace),
	}
}

func (g *Group) Describe(ch chan<- *prometheus.Desc) {
	g.Containers.Describe(ch)
}

func (g *Group) Collect(ch chan<- prometheus.Metric) {
	g.Containers.Collect(ch)
}
