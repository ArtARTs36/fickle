package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"log/slog"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	visits := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "helloworld_visits",
	})
	prometheus.MustRegister(visits)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("Hello, world"))
		visits.Inc()
	})

	slog.Info("listening on 8000")

	if err := http.ListenAndServe(":8000", http.DefaultServeMux); err != nil {
		slog.Error("failed to listen", slog.Any("err", err))
	}
}
