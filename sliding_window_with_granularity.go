package ratelimiter

import (
	"sync"
	"sync/atomic"
	"time"
)

type SlidingWindowOptions struct {
	//limit of requests
	Limit int32
	//whole window interval, limit sets per window
	Interval time.Duration
	//one tick of sliding window interval
	Tick time.Duration
}

type SlidingWindowWithGranularity struct {
	opts SlidingWindowOptions //do not want external changes, so copy

	//per tick requests storage
	requests []int32
	//current tick position
	position int
	//total requests allowed for current tick
	total int32

	requestsLock sync.RWMutex

	//already used requests for current tick
	used int32
}

func NewSlidingWindow(opts SlidingWindowOptions) *SlidingWindowWithGranularity {
	return &SlidingWindowWithGranularity{
		opts:     opts,
		total:    opts.Limit,
		requests: make([]int32, opts.Interval/opts.Tick),
	}
}

func (s *SlidingWindowWithGranularity) NextTick() {
	s.requestsLock.Lock()
	defer s.requestsLock.Unlock()

	used := atomic.LoadInt32(&s.used)

	//store this tick used requests
	s.requests[s.position] = used
	s.total -= used

	//advance to new tick, so no used requests yet
	s.used = 0

	s.position++
	if s.position >= len(s.requests) {
		s.position = 0
	}

	//restore used requests from same tick of previous whole window
	s.total += s.requests[s.position]
}

func (s *SlidingWindowWithGranularity) IsTokenExists() bool {
	s.requestsLock.RLock()
	defer s.requestsLock.RUnlock()

	return s.used < s.total
}

func (s *SlidingWindowWithGranularity) GetTokens(amount int32) (actual int32) {
	s.requestsLock.RLock()
	defer s.requestsLock.RUnlock()

	//return up to actual amount of ticks
	actual = s.total - atomic.LoadInt32(&s.used)
	if actual > amount {
		actual = amount
	}

	atomic.AddInt32(&s.used, actual)

	return actual
}

func (s *SlidingWindowWithGranularity) GetTickDuration() time.Duration {
	return s.opts.Tick
}
