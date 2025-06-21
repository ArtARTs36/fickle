package proxy

import (
	"context"
	"errors"
	"fmt"
	"github.com/artarts36/fickle/internal/metrics"
	"github.com/docker/docker/api/types/filters"

	"github.com/artarts36/fickle/internal/metricsscrapper"
	"github.com/artarts36/fickle/internal/transport"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"

	"github.com/artarts36/fickle/internal/cfg"
)

type ContainerProxy struct {
	config       cfg.Proxy
	dockerClient *client.Client

	containerRanLock sync.RWMutex
	containerRan     bool

	proxy  *httputil.ReverseProxy
	target *url.URL

	lastRequestStartedLock sync.RWMutex
	lastRequestStarted     time.Time

	metricsScrapper *metricsscrapper.Scrapper
	metricsGroup    *metrics.Group
}

func NewContainerProxy(
	config cfg.Proxy,
	dockerClient *client.Client,
	metricsScrapper *metricsscrapper.Scrapper,
	metricsGroup *metrics.Group,
) *ContainerProxy {
	p := &ContainerProxy{
		target: &url.URL{
			Scheme: "http",
			Host:   config.To.Address,
		},
		config:          config,
		dockerClient:    dockerClient,
		metricsScrapper: metricsScrapper,
		metricsGroup:    metricsGroup,
	}

	p.proxy = httputil.NewSingleHostReverseProxy(p.target)
	p.proxy.Transport = transport.Retryable(config.RetryPolicy, http.DefaultTransport)

	go func() {
		p.recycle()
	}()

	return p
}

func (p *ContainerProxy) HandleRequest(w http.ResponseWriter, req *http.Request) {
	req.Host = p.target.Host

	p.lastRequestStartedLock.Lock()
	p.lastRequestStarted = time.Now()
	p.lastRequestStartedLock.Unlock()

	p.containerRanLock.Lock()
	if p.containerRan {
		p.containerRanLock.Unlock()
	} else {
		err := p.startContainer(req.Context())
		if err != nil {
			p.containerRanLock.Unlock()

			slog.ErrorContext(
				req.Context(),
				"[container-proxy] failed to start container",
				slog.Any("err", err),
				slog.String("service", p.config.To.ServiceName),
			)
			p.error(w)
			return
		}

		p.containerRan = true
		p.containerRanLock.Unlock()
	}

	p.proxy.ServeHTTP(w, req)
}

func (p *ContainerProxy) startContainer(ctx context.Context) error {
	cont, err := p.findContainer(ctx)
	if err != nil {
		return fmt.Errorf("find: %w", err)
	}

	if cont.State == container.StateRunning || cont.State == container.StateRestarting {
		slog.InfoContext(ctx, "[container-proxy] container already running", slog.String("container_id", cont.ID))

		return nil
	}

	slog.InfoContext(ctx, "[container-proxy] running container", slog.String("container_id", cont.ID))

	err = p.dockerClient.ContainerStart(ctx, cont.ID, container.StartOptions{})
	p.metricsGroup.Containers.IncRun(p.config.Host, err == nil)

	return err
}

func (p *ContainerProxy) findContainer(ctx context.Context) (*container.Summary, error) {
	listOpts := container.ListOptions{
		All: true,
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "label",
			Value: fmt.Sprintf("fickle.service.name=%s", p.config.To.ServiceName),
		}),
	}

	containers, err := p.dockerClient.ContainerList(ctx, listOpts)
	if err != nil {
		return nil, fmt.Errorf("list containers: %w", err)
	}

	for _, cont := range containers {
		return &cont, nil
	}

	return nil, errors.New("container not found")
}

func (p *ContainerProxy) stopContainer(ctx context.Context) error {
	cont, err := p.findContainer(ctx)
	if err != nil {
		return fmt.Errorf("find container: %w", err)
	}

	slog.InfoContext(ctx, "[container-proxy] stopping container", slog.String("container_id", cont.ID))

	if p.config.Metrics.Scrape.Address != "" {
		p.metricsScrapper.Scrape(p.config.Host, p.config.Metrics.Scrape.Address)
	}

	err = p.dockerClient.ContainerStop(ctx, cont.ID, container.StopOptions{})
	p.metricsGroup.Containers.IncStops(p.config.Host, err == nil)

	return nil
}

func (p *ContainerProxy) recycle() {
	tick := time.NewTicker(5 * time.Second)

	stop := func(t time.Time) error {
		p.containerRanLock.Lock()
		defer p.containerRanLock.Unlock()
		if !p.containerRan {
			return nil
		}

		if t.Before(p.lastRequestStarted.Add(p.config.IdleTimeout)) {
			return nil
		}

		p.containerRan = false

		return p.stopContainer(context.Background())
	}

	for t := range tick.C {
		err := stop(t)
		if err != nil {
			slog.InfoContext(context.Background(), "[container-proxy] failed to stop container", slog.Any("err", err))
		}
	}
}

func (p *ContainerProxy) error(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadGateway)
}
