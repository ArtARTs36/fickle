package main

import (
	"context"
	"fmt"
	"github.com/artarts36/fickle/internal/cfg"
	"github.com/artarts36/fickle/internal/control"
	"github.com/artarts36/fickle/internal/metricsscrapper"
	"github.com/artarts36/fickle/internal/proxy"
	"github.com/docker/docker/client"
	"log/slog"
	"os"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	ctx := context.Background()

	if err := run(ctx); err != nil {
		slog.ErrorContext(ctx, "run application failed", slog.Any("err", err))
	}
}

func run(ctx context.Context) error {
	slog.InfoContext(ctx, "read config")

	config, err := cfg.Load("fickle.yaml")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	slog.InfoContext(ctx, fmt.Sprintf("log level: %s", config.Log.Level.String()))

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: config.Log.Level.Value,
	})))

	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("create docker client: %w", err)
	}

	metricsStore := metricsscrapper.NewStore()
	metricsScrapper := metricsscrapper.NewScrapper(metricsStore)

	prox := proxy.NewServer(config, dockerClient, metricsScrapper)

	go func() {
		if config.Control.Address == "" {
			return
		}

		slog.Info("running control server", slog.String("address", config.Control.Address))

		controlServer := control.NewServer(
			metricsScrapper,
			metricsStore,
			config,
			prox,
		)

		if err = controlServer.Run(); err != nil {
			slog.Error("failed to run control server", slog.Any("err", err))
		}
	}()

	return prox.Run()
}
