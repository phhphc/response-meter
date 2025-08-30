package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/phhphc/response-meter/internal/collector"
	"github.com/phhphc/response-meter/internal/meter"
	"github.com/phhphc/response-meter/internal/reporter"
)

func main() {
	target := flag.String("t", "", "HTTP(S) URL to probe. (required)")
	interval := flag.Duration("i", 2*time.Second, "Interval between reports. (e.g., 500ms, 2s, 1m)")
	concurrency := flag.Int("c", 1, "Number of concurrent requests.")
	timeout := flag.Duration("d", 0, "Per-request timeout. (0 disables timeout)")
	flag.Parse()

	if *target == "" {
		fmt.Fprintf(os.Stderr, "error: target URL is required (use -t)\n")
		flag.Usage()
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	f := collector.HTTPStatusCollectorFactory{
		Target:  *target,
		Timeout: *timeout,
	}
	r := reporter.NewTUIReporter()
	m := meter.Meter{
		CollectorFactory: f,
		Reporter:         r,
		ReportInterval:   *interval,
	}

	err := m.Measure(ctx, *concurrency)
	if err != nil && !errors.Is(err, context.Canceled) {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}
}
