package cfg

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
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
	}

	return &cfg, nil
}
