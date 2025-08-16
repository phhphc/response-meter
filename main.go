package main

import (
	"context"
	"flag"
	"fmt"
	"os"

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

	ctx := context.TODO()

	f := collector.NewHTTPStatusCollectorFactory(*target)
	m := meter.New(f)

	if err := m.Measure(ctx, *concurrency); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
