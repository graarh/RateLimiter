package ratelimiter

import "time"

type UnixNanoTime int64

type RateLimiter interface {
	//request amount of tokens for actions
	GetTokens(amount int32) (actual int32)
	//check that tokens exists
	IsTokenExists() bool

	//advance to the next tick of limiter
	NextTick()
	//get single tick duration
	GetTickDuration() time.Duration
}
