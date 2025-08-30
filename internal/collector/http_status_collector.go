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
		target: f.Target,
		client: &http.Client{
			Transport: &http.Transport{},
			Timeout:   f.Timeout,
		},
	}, nil
}

type httpResponseCollector struct {
	target string
	client *http.Client
}

func (h *httpResponseCollector) Collect(ctx context.Context) (string, error) {
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
