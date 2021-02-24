package ratelimiter_test

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"ratelimiter"
	"time"
)

type FakeLimiter struct {
	ratelimiter.FixedWindow

	NextTickCalled bool
}

func NewFakeLimiter(opts ratelimiter.LimiterOptions) *FakeLimiter {
	return &FakeLimiter{
		FixedWindow: *ratelimiter.NewFixedWindow(opts),
	}
}

func (f *FakeLimiter) NextTick() {
	f.NextTickCalled = true
}

var _ = Describe("Clock controller", func() {

	opts := ratelimiter.LimiterOptions{
		Limit:    10,
		Interval: 10 * time.Second,
		Tick:     10 * time.Millisecond,
	}

	It("Ticker controller should call NextTick every tick", func() {
		l1 := NewFakeLimiter(opts)
		l2 := NewFakeLimiter(opts)

		ctx, cancel := context.WithCancel(context.Background())
		ticker := make(chan struct{})
		tickerOpts := ratelimiter.TickerControllerOptions{
			Ctx:          ctx,
			Ticker:       ticker,
			TickDuration: 10 * time.Millisecond,
		}
		controller := ratelimiter.NewTickerController(tickerOpts)

		err := controller.AddLimiter(l1)
		Expect(err).To(BeNil())

		err = controller.AddLimiter(l2)
		Expect(err).To(BeNil())

		<-ticker
		<-ticker

		Expect(l1.NextTickCalled).To(BeTrue())
		Expect(l2.NextTickCalled).To(BeTrue())

		cancel()
	})
})
