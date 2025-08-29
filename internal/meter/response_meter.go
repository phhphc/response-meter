package meter

import (
	"context"
	"fmt"
	"time"

	"github.com/phhphc/response-meter/pkg/errgroup"
)

type Meter struct {
	CollectorFactory CollectorFactory
	Reporter         Reporter
}

type CollectorFactory interface {
	NewCollector() (Collector, error)
}

type Collector interface {
	Collect(ctx context.Context) (string, error)
}

type Reporter interface {
	Update(s Stats) error
}

type Stats struct {
	TotalCounts        map[string]int
	TotalDuration      time.Duration
	LastPeriodCounts   map[string]int
	LastPeriodDuration time.Duration
}

func New(f CollectorFactory, r Reporter) *Meter {
	return &Meter{
		CollectorFactory: f,
		Reporter:         r,
	}
}

func (m Meter) Measure(ctx context.Context, concurrency int) error {
	g, ctx := errgroup.WithContext(ctx)

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
		t := time.NewTicker(2 * time.Second)
		defer t.Stop()

		s := Stats{
			TotalCounts:        make(map[string]int),
			TotalDuration:      0 * time.Second,
			LastPeriodCounts:   make(map[string]int),
			LastPeriodDuration: 0 * time.Second,
		}
		start := time.Now()
		periodStart := start

		for {
			select {
			case <-ctx.Done():
				return nil
			case res := <-ch:
				s.LastPeriodCounts[res]++
			case <-t.C:
				now := time.Now()
				s.TotalDuration = now.Sub(start)
				s.LastPeriodDuration = now.Sub(periodStart)
				periodStart = now
				for k, v := range s.LastPeriodCounts {
					s.TotalCounts[k] += v
				}
				if err := m.Reporter.Update(s); err != nil {
					return fmt.Errorf("failed to update report: %w", err)
				}
				for k := range s.LastPeriodCounts {
					s.LastPeriodCounts[k] = 0
				}
			}
		}
	})

	return g.Wait()
}
