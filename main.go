package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/phhphc/response-meter/collector"
	"github.com/phhphc/response-meter/meter"
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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	f := collector.NewHTTPStatusCollectorFactory(*target)
	m := meter.New(f)

	err := m.Measure(ctx, *concurrency)
	if err != nil && !errors.Is(err, context.Canceled) {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}
}
