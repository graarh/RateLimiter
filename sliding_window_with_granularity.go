package ratelimiter

type SlidingWindowWithGranularity struct {
	Limiter

	//per tick requests storage
	requests []int32
	//current tick position
	position int
}

//check tests if you want to know difference with common sliding window algorithm
func NewSlidingWindowWithGranularity(opts LimiterOptions) *SlidingWindowWithGranularity {
	return &SlidingWindowWithGranularity{
		Limiter: Limiter{
			opts:  opts,
			total: opts.Limit,
		},
		requests: make([]int32, opts.Interval/opts.Tick),
	}
}

func (l *SlidingWindowWithGranularity) NextTick() {
	l.requestsLock.Lock()
	defer l.requestsLock.Unlock()

	//store this tick used requests
	l.requests[l.position] = l.used
	l.total -= l.used

	//advance to new tick, so no used requests yet
	l.used = 0

	l.position++
	if l.position >= len(l.requests) {
		l.position = 0
	}

	//restore used requests from same tick of previous whole window
	l.total += l.requests[l.position]
}
