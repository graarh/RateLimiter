package ratelimiter

import (
	"context"
	"errors"
	"sync"
	"time"
)

type TickerControllerOptions struct {
	Ctx context.Context
	// syncing, handy for tests, set to nil if do not use
	Ticker       chan<- struct{}
	TickDuration time.Duration
}

type TickerController struct {
	// i do not want external changes, so i store copy
	opts TickerControllerOptions

	// limiters under control
	limiters     []RateLimiter
	limitersLock sync.Mutex
}

func (t *TickerController) controller() {
	// Not protected from system hanging
	// can skip some ticks in this case
	ticker := time.NewTicker(t.opts.TickDuration)
	for {
		select {
		case <-t.opts.Ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			if t.opts.Ticker != nil {
				t.opts.Ticker <- struct{}{}
			}
			t.limitersLock.Lock()
			for _, limiter := range t.limiters {
				limiter.NextTick()
			}
			t.limitersLock.Unlock()
		}
	}
}

func (t *TickerController) AddLimiter(limiter RateLimiter) error {
	if limiter.GetTickDuration() != t.opts.TickDuration {
		return errors.New("limiter tick duration differs from ticker controller")
	}

	t.limitersLock.Lock()
	defer t.limitersLock.Unlock()

	t.limiters = append(t.limiters, limiter)

	return nil
}

func NewTickerController(opts TickerControllerOptions) *TickerController {
	tc := TickerController{
		opts:     opts,
		limiters: make([]RateLimiter, 0, 1),
	}
	go tc.controller()
	return &tc
}
