package control

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/artarts36/fickle/internal/cfg"
	"github.com/artarts36/fickle/internal/metricsscrapper"
	"github.com/artarts36/fickle/internal/proxy"
)

type Server struct {
	mux    *http.ServeMux
	config *cfg.Config
}

func NewServer(
	metricsScrapper *metricsscrapper.Scrapper,
	metricsStore *metricsscrapper.Store,
	config *cfg.Config,
	proxyServer *proxy.Server,
) *Server {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/metrics/{name}", NewHostMetricsHandler(metricsScrapper, metricsStore, config, proxyServer).HandleRequest)

	return &Server{mux: mux, config: config}
}

func (s *Server) Run() error {
	return http.ListenAndServe(s.config.Control.Address, s.mux)
}
