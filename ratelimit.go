package ratelimit

import (
	"context"
	"time"

	"golang.org/x/sync/errgroup"
)

// ErrGroup implements the errgroup.Group interface, but limits goroutines to a
// execute at a max throughput.
type ErrGroup struct {
	ticker *time.Ticker
	eg     *errgroup.Group
}

// WithContext initializes and returns a new RateLimiter
func WithContext(ctx context.Context, rps int) (*ErrGroup, context.Context) {
	eg, ctx := errgroup.WithContext(ctx)
	r := &ErrGroup{
		ticker: time.NewTicker(time.Second / time.Duration(rps)),
		eg:     eg,
	}
	return r, ctx
}

// Go runs a new job as a goroutine. It will wait until the next available
// time so that the ratelimit is not exceeded.
func (e *ErrGroup) Go(fn func() error) {
	<-e.ticker.C
	go e.eg.Go(fn)
}

// Wait will wait until all jobs are processed. Once Wait() is called, no more jobs
// can be added.
func (e *ErrGroup) Wait() error {
	defer e.ticker.Stop()
	return e.eg.Wait()
}
