package cfg

import (
	"github.com/artarts36/fickle/internal/transport"
	"time"

	"github.com/artarts36/specw"
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
	Host        string `yaml:"-"`
	ServiceName string `yaml:"service_name"`

	Forward struct {
		Address     string                `yaml:"address"`
		RetryPolicy transport.RetryPolicy `yaml:"retry_policy"`
	} `yaml:"forward"`

	IdleTimeout time.Duration `yaml:"idle_timeout"`
	Metrics     struct {
		Scrape struct {
			Address string `yaml:"address"`
		} `yaml:"scrape"`
	} `yaml:"metrics"`
}
