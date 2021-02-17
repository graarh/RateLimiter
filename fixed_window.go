package ratelimiter

type FixedWindow struct {
	Limiter

	//current tick position
	position int
	//ticks per window
	ticksPerWindow int
}

func NewFixedWindow(opts LimiterOptions) *FixedWindow {
	return &FixedWindow{
		Limiter: Limiter{
			opts:  opts,
			total: opts.Limit,
		},
		ticksPerWindow: int(opts.Interval / opts.Tick),
	}
}

func (f *FixedWindow) NextTick() {
	f.requestsLock.Lock()
	defer f.requestsLock.Unlock()

	f.position++

	//new window position
	if f.position >= f.ticksPerWindow {
		f.total = f.opts.Limit
		f.used = 0
		f.position = 0
	}
}
