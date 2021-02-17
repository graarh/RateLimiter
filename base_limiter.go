package ratelimiter

import (
	"sync"
	"sync/atomic"
	"time"
)

type LimiterOptions struct {
	//limit of requests
	Limit int32
	//whole window interval, limit sets per window
	Interval time.Duration
	//one tick of sliding window interval
	Tick time.Duration
}

type Limiter struct {
	opts LimiterOptions //do not want external changes, so copy

	requestsLock sync.RWMutex

	//total requests allowed for current tick
	total int32

	//already used requests for current tick
	used int32
}

func (l *Limiter) IsTokenExists() bool {
	l.requestsLock.RLock()
	defer l.requestsLock.RUnlock()

	return atomic.LoadInt32(&l.used) < l.total
}

func (l *Limiter) GetTokens(amount int32) (actual int32) {
	l.requestsLock.RLock()
	defer l.requestsLock.RUnlock()

	//return up to actual amount of ticks
	actual = l.total - atomic.LoadInt32(&l.used)
	if actual < 0 {
		actual = 0
	}
	if actual > amount {
		actual = amount
	}

	atomic.AddInt32(&l.used, actual)

	return actual
}

func (l *Limiter) GetTickDuration() time.Duration {
	return l.opts.Tick
}
