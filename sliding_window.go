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

func (s *SlidingWindow) NextTick() {
	s.requestsLock.Lock()
	defer s.requestsLock.Unlock()

	s.position++

	//new window position
	if s.position >= s.ticksPerWindow {
		s.prevWindowRequests = s.used
		s.used = 0
		s.position = 0
	}

	//reserve some requests from limit based on prev window requests amount
	s.total = s.Limiter.opts.Limit - s.prevWindowRequests*
		int32(s.ticksPerWindow-s.position-1)/int32(s.ticksPerWindow)
}
