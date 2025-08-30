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
	target := flag.String("t", "", "target url")
	concurrency := flag.Int("c", 1, "number of concurrent requests")
	interval := flag.Duration("i", 2*time.Second, "report interval (default: 2s)")
	timeout := flag.Duration("d", 0, "timeout duration (default: 0 - no timeout)")
	flag.Parse()

	if *target == "" {
		fmt.Fprintf(os.Stderr, "error: target url is required (use -t flag)\n")
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
