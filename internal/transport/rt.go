package transport

import (
	"net/http"
	"time"
)

type RetryableRoundTripper struct {
	cfg  RetryPolicy
	next http.RoundTripper
}

type RetryPolicy struct {
	Attempts int           `yaml:"attempts"`
	Backoff  time.Duration `yaml:"backoff"`
}

func Retryable(rp RetryPolicy, next http.RoundTripper) http.RoundTripper {
	return &RetryableRoundTripper{
		cfg:  rp,
		next: next,
	}
}

func (t *RetryableRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	backoff := t.cfg.Backoff

	for i := 0; i < t.cfg.Attempts; i++ {
		resp, err = t.next.RoundTrip(req)
		if err == nil {
			break
		}

		time.Sleep(backoff)
		backoff *= 2
	}

	return resp, err
}
