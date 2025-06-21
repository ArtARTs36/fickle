package proxy

import (
	"github.com/artarts36/fickle/internal/metricsscrapper"
	"log/slog"
	"net/http"

	"github.com/artarts36/fickle/internal/cfg"
	"github.com/docker/docker/client"
)

type Server struct {
	containers map[string]*ContainerProxy
}

func NewServer(
	config *cfg.Config,
	dockerClient *client.Client,
	metricsScrapper *metricsscrapper.Scrapper,
) *Server {
	s := &Server{
		containers: map[string]*ContainerProxy{},
	}

	for _, c := range config.Proxy {
		s.containers[c.Host] = NewContainerProxy(c, dockerClient, metricsScrapper)
	}

	return s
}

func (s *Server) Run() error {
	slog.Info("running proxy server", slog.String("address", ":80"))

	http.HandleFunc("/", s.handleRequest)

	return http.ListenAndServe(":80", nil)
}

func (s *Server) Enabled(host string) bool {
	cont, ok := s.containers[host]
	return ok && cont.containerRan
}

func (s *Server) handleRequest(w http.ResponseWriter, req *http.Request) {
	slog.DebugContext(req.Context(), "[proxy-server] handling request", slog.String("host", req.Host))

	cont, ok := s.containers[req.Host]
	if !ok {
		slog.WarnContext(req.Context(), "[proxy-server] proxy for host not found", slog.String("host", req.Host))
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	cont.HandleRequest(w, req)
}
