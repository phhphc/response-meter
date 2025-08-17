package collector

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/phhphc/response-meter/internal/meter"
)

type HTTPStatusCollectorFactory struct {
	target            string
	DisableKeepAlives bool
}

func NewHTTPStatusCollectorFactory(target string) meter.CollectorFactory {
	return HTTPStatusCollectorFactory{
		target: target,
	}
}

func (f HTTPStatusCollectorFactory) NewCollector() (meter.Collector, error) {
	return &httpResponseCollector{
		target: f.target,
		client: &http.Client{
			Transport: &http.Transport{},
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
		return "", fmt.Errorf("failed to do request: %w", err)
	}
	defer resp.Body.Close()

	return strconv.Itoa(resp.StatusCode), nil
}
