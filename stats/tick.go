package stats

import "time"

type Ticks struct {
	wm *WelfordMean

	lastTick    time.Time
	lastTickGap time.Duration
	lastReset   time.Time
	lastUpdate  time.Time
}

func NewTicks() *Ticks {
	now := time.Now()
	return &Ticks{
		wm:         new(WelfordMean),
		lastReset:  now,
		lastUpdate: now,
	}
}

func (t *Ticks) PerSecond() float64 {
	if t.wm.Mean != 0 {
		return 1.0 / time.Duration(t.wm.Mean).Seconds()
	}
	return 0
}

func (t *Ticks) in(tick time.Time) {
	if !t.lastTick.IsZero() {
		t.lastTickGap = tick.Sub(t.lastTick)
	}
	t.lastTick = tick

	shouldReset := tick.Sub(t.lastReset) > time.Second

	if shouldReset {
		t.wm.Reset()
		t.lastReset = tick
	}

	if tick.Sub(t.lastUpdate) >= time.Nanosecond || shouldReset {
		t.wm.In(float64(t.lastTickGap))
		t.lastUpdate = tick
	}
}
