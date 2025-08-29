package errgroup

import (
	"context"
	"sync"
)

func WithContext(ctx context.Context) (*Group, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &Group{cancel: cancel}, ctx
}

// Group is a collection of goroutines working on subtasks of the same overall task.
// A Group must not be copied after first use (because it contains sync types).
type Group struct {
	wg      sync.WaitGroup
	errOnce sync.Once
	err     error
	cancel  context.CancelFunc
}

// Go starts f in a new goroutine. On the first non-nil error, the group's context is canceled.
func (g *Group) Go(f func() error) {
	g.wg.Go(func() {
		if err := f(); err != nil {
			g.errOnce.Do(func() {
				g.err = err
				if g.cancel != nil {
					g.cancel()
				}
			})
		}
	})
}

// Wait blocks for all goroutines and returns the first error, if any.
// It also cancels the context to release resources.
func (g *Group) Wait() error {
	g.wg.Wait()
	if g.cancel != nil {
		g.cancel()
	}
	return g.err
}
