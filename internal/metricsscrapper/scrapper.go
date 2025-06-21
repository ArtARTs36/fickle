package metricsscrapper

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type Scrapper struct {
	metrics *Store
}

func NewScrapper(metrics *Store) *Scrapper {
	return &Scrapper{metrics: metrics}
}

func (m *Scrapper) Scrape(name, address string) {
	err := m.scrape(name, address)
	if err != nil {
		slog.Error("failed to scrape metrics", slog.Any("err", err))
	}
}

func (m *Scrapper) scrape(name, address string) error {
	resp, err := http.Get(address)
	if err != nil {
		return fmt.Errorf("get metrics: %w", err)
	}

	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	m.metrics.Put(name, content)

	return nil
}
