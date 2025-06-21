package control

import (
	"github.com/artarts36/fickle/internal/cfg"
	"github.com/artarts36/fickle/internal/metricsscrapper"
	"github.com/artarts36/fickle/internal/proxy"
	"log/slog"
	"net/http"
)

type HostMetricsHandler struct {
	scrapper *metricsscrapper.Scrapper
	store    *metricsscrapper.Store
	config   *cfg.Config
	server   *proxy.Server
}

func NewHostMetricsHandler(
	scrapper *metricsscrapper.Scrapper,
	store *metricsscrapper.Store,
	config *cfg.Config,
	server *proxy.Server,
) *HostMetricsHandler {
	return &HostMetricsHandler{scrapper: scrapper, store: store, config: config, server: server}
}

func (h *HostMetricsHandler) HandleRequest(w http.ResponseWriter, req *http.Request) {
	name := req.PathValue("name")

	if h.server.Enabled(name) {
		slog.DebugContext(
			req.Context(),
			"[host-metrics] get metrics directly from container",
			slog.String("host", name),
		)

		pr := h.config.Proxy[name]

		h.scrapper.Scrape(name, pr.Metrics.Scrape.Address)
	}

	metrics, ok := h.store.Get(name)
	if !ok {
		http.NotFound(w, req)
		return
	}

	_, err := w.Write(metrics)
	if err != nil {
		slog.ErrorContext(req.Context(), "failed to write metrics response", slog.Any("err", err))
	}
}
