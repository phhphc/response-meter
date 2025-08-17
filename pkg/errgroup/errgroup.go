package errgroup

import (
	"sync"
)

// Group is a collection of goroutines working on subtasks that are part of
// the same overall task.
type Group struct {
	wg sync.WaitGroup

	errOnce sync.Once
	err     error
}

// Go calls the given function in a new goroutine.
// The first call to return a non-nil error cancels the group; its error will be
// returned by Wait.
func (g *Group) Go(f func() error) {
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()

		if err := f(); err != nil {
			g.errOnce.Do(func() {
				g.err = err
			})
		}
	}()
}

// Wait blocks until all function calls from the Go method have returned, then
// returns the first non-nil error (if any) from them.
func (g *Group) Wait() error {
	g.wg.Wait()
	return g.err
}
