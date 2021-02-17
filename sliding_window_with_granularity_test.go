package ratelimiter_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"ratelimiter"
	"time"
)

var _ = Describe("Sliding window with granularity", func() {
	Describe("Basic tests, sync, one instance", func() {
		var limiter ratelimiter.RateLimiter

		BeforeEach(func() {
			opts := ratelimiter.SlidingWindowOptions{
				Limit:    5,
				Interval: time.Millisecond * 5,
				Tick:     time.Millisecond,
			}
			limiter = ratelimiter.NewSlidingWindow(opts)
		})

		It("Get some, no tick advance", func() {
			Expect(limiter.GetTokens(2)).To(Equal(int32(2)))
			Expect(limiter.GetTokens(2)).To(Equal(int32(2)))
			Expect(limiter.GetTokens(2)).To(Equal(int32(1)))
			Expect(limiter.GetTokens(2)).To(Equal(int32(0)))
		})

		It("Get some, with tick advance", func() {
			for i := 0; i < 5; i++ {
				limiter.GetTokens(2)
				limiter.NextTick()
			}
			Expect(limiter.GetTokens(10)).To(Equal(int32(2)))
			limiter.NextTick()
			Expect(limiter.GetTokens(10)).To(Equal(int32(2)))
			limiter.NextTick()
			Expect(limiter.GetTokens(10)).To(Equal(int32(1)))
			limiter.NextTick()
			Expect(limiter.GetTokens(10)).To(Equal(int32(0)))
			limiter.NextTick()
			Expect(limiter.GetTokens(10)).To(Equal(int32(0)))
		})
	})
})
