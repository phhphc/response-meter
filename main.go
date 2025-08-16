package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	target := flag.String("t", "", "target url")
	concurrency := flag.Int("c", 10, "number of concurrent requests")
	flag.Parse()

	if *target == "" {
		fmt.Fprintf(os.Stderr, "error: target url is required (use -t flag)\n")
		flag.Usage()
		os.Exit(1)
	}

	ctx := context.TODO()
	if err := do(ctx, *target, *concurrency); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func do(ctx context.Context, target string, concurrency int) error {
	g := errgroup.Group{}

	ch := make(chan string, 1000*concurrency)
	defer close(ch)

	for range concurrency {
		g.Go(func() error {
			c := NewCollector()

			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
					res, err := c.Collect(ctx, target)
					if err != nil {
						return fmt.Errorf("failed to collect: %w", err)
					}
					ch <- res
				}
			}

		})
	}

	g.Go(func() error {
		start := time.Now()
		t := time.NewTicker(2 * time.Second)
		m := make(map[string]int)
		for {
			select {
			case <-ctx.Done():
				return nil
			case res := <-ch:
				m[res]++
			case <-t.C:
				dur := time.Since(start).Seconds()
				total := 0
				for _, c := range m {
					total += c
				}
				avg := float64(total) / dur
				fmt.Printf("%0.3f resq/s, %v\n", avg, m)
			}
		}
	})

	return g.Wait()
}

type Collector struct {
	client *http.Client
}

func NewCollector() *Collector {
	return &Collector{
		client: &http.Client{
			Transport: &http.Transport{},
		},
	}
}

func (h *Collector) Collect(ctx context.Context, target string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
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
