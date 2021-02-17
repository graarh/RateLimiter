package ratelimiter

type SlidingWindow struct {
	Limiter

	//previous window requests
	prevWindowRequests int32
	//current tick position
	position int
	//ticks per window
	ticksPerWindow int
}

func NewSlidingWindow(opts LimiterOptions) *SlidingWindow {
	return &SlidingWindow{
		Limiter: Limiter{
			opts:  opts,
			total: opts.Limit,
		},
		ticksPerWindow: int(opts.Interval / opts.Tick),
	}
}

func (f *SlidingWindow) NextTick() {
	f.requestsLock.Lock()
	defer f.requestsLock.Unlock()

	f.position++

	//new window position
	if f.position >= f.ticksPerWindow {
		f.prevWindowRequests = f.used
		f.used = 0
		f.position = 0
	}

	//reserve some requests from limit based on prev window requests amount
	f.total = f.Limiter.opts.Limit - f.prevWindowRequests*
		int32(f.ticksPerWindow-f.position-1)/int32(f.ticksPerWindow)
}
