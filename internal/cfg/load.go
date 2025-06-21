package cfg

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
	"time"
)

const (
	defaultIdleTimeout           = 1 * time.Minute
	defaultRetryRequestsAttempts = 3
	defaultRetryBackoff          = 100 * time.Millisecond
)

func Load(path string) (*Config, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	var cfg Config

	err = yaml.Unmarshal(content, &cfg)
	if err != nil {
		return nil, fmt.Errorf("unmarhsal yaml: %w", err)
	}

	for host, proxy := range cfg.Proxy {
		if !strings.HasPrefix("http", proxy.Metrics.Scrape.Address) {
			proxy.Host = host
			proxy.Metrics.Scrape.Address = "http://" + proxy.Metrics.Scrape.Address
			cfg.Proxy[host] = proxy
		}

		if proxy.IdleTimeout <= 0 {
			proxy.IdleTimeout = defaultIdleTimeout
		}

		if proxy.Forward.RetryPolicy.Attempts <= 0 {
			proxy.Forward.RetryPolicy.Attempts = defaultRetryRequestsAttempts
		}

		if proxy.Forward.RetryPolicy.Backoff <= 0 {
			proxy.Forward.RetryPolicy.Backoff = defaultRetryBackoff
		}
	}

	return &cfg, nil
}
