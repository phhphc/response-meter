package collector

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/phhphc/response-meter/internal/meter"
)

type HTTPStatusCollectorFactory struct {
	Target  string
	Timeout time.Duration
}

func (f HTTPStatusCollectorFactory) NewCollector() (meter.Collector, error) {
	return &httpResponseCollector{
		target:  f.Target,
		timeout: f.Timeout,
		client: &http.Client{
			Transport: &http.Transport{},
		},
	}, nil
}

type httpResponseCollector struct {
	target  string
	timeout time.Duration
	client  *http.Client
}

func (h *httpResponseCollector) Collect(ctx context.Context) (string, error) {
	if h.timeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, h.timeout)
		defer cancel()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, h.target, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return "Timeout", nil
		}
		return "", fmt.Errorf("failed to do request: %w", err)
	}
	defer resp.Body.Close()

	return strconv.Itoa(resp.StatusCode), nil
}
