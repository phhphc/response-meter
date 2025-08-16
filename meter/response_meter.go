package meter

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

type Meter struct {
	CollectorFactory CollectorFactory
}

type CollectorFactory interface {
	NewCollector() (Collector, error)
}

type Collector interface {
	Collect(ctx context.Context) (string, error)
}

func New(f CollectorFactory) *Meter {
	return &Meter{
		CollectorFactory: f,
	}
}

func (m Meter) Measure(ctx context.Context, concurrency int) error {
	g := errgroup.Group{}

	ch := make(chan string, 1000*concurrency)
	defer close(ch)

	for range concurrency {
		g.Go(func() error {
			c, err := m.CollectorFactory.NewCollector()
			if err != nil {
				return fmt.Errorf("failed to create collector: %w", err)
			}
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
					res, err := c.Collect(ctx)
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
		counts := make(map[string]int)
		for {
			select {
			case <-ctx.Done():
				return nil
			case res := <-ch:
				counts[res]++
			case <-t.C:
				dur := time.Since(start).Seconds()
				total := 0
				for _, c := range counts {
					total += c
				}
				avg := float64(total) / dur
				parcentages := make(map[string]float64)
				for res, c := range counts {
					parcentages[res] = (float64(c) / float64(total)) * 100
				}
				fmt.Printf("%0.3f resq/s | %v | %v\n", avg, parcentages, counts)
			}
		}
	})

	return g.Wait()
}
