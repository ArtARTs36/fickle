package cfg

import (
	"time"

	"github.com/artarts36/specw"

	"github.com/artarts36/fickle/internal/transport"
)

type Config struct {
	Proxy   map[string]Proxy `yaml:"proxy"`
	Control struct {
		Address string `yaml:"address"`
	} `yaml:"control"`
	Log struct {
		Level specw.SlogLevel `yaml:"level"`
	} `yaml:"log"`
}

type Proxy struct {
	Host        string                `yaml:"-"`
	RetryPolicy transport.RetryPolicy `yaml:"retry_policy"`
	To          struct {
		NetworkName string `yaml:"network_name"`
		ServiceName string `yaml:"service_name"`
		Address     string `yaml:"address"`
	} `yaml:"to"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
	Metrics     struct {
		Scrape struct {
			Address string `yaml:"address"`
		} `yaml:"scrape"`
	} `yaml:"metrics"`
}
