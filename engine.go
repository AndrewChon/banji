package banji

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type Component interface {
	Bootstrap() []Listener
}

type Engine struct {
	cfg           *Configuration
	bus           *bus
	ticker        *time.Ticker
	tickWG        sync.WaitGroup
	context       context.Context
	contextCancel context.CancelFunc

	active atomic.Bool
}

func New(options ...Option) *Engine {
	eng := &Engine{cfg: NewConfiguration(options...)}
	eng.bootstrap()

	return eng
}

func (eng *Engine) Start() {
	eng.active.Store(true)
	eng.Post(newStartEvent(), important)

	go func() {
		for {
			select {
			case tick := <-eng.ticker.C:
				eng.tickWG.Add(1)
				eng.Post(newPreTickEvent(tick), important)
				eng.bus.Tick()
				eng.Post(newPostTickEvent(tick, time.Now()), important)
				eng.tickWG.Done()
			case <-eng.context.Done():
				return
			}
		}
	}()
}

func (eng *Engine) Stop() {
	if !eng.active.Load() {
		return
	}

	eng.Post(newStopEvent(), important)
	eng.active.Store(false)

	go func() {
		for eng.bus.Size() > 0 {
			time.Sleep(time.Second)
		}
	}()

	eng.ticker.Stop()
	eng.tickWG.Wait()
	eng.contextCancel()
}

func (eng *Engine) Size() int {
	return eng.bus.queue.Size() + eng.bus.bufQueue.Size()
}

func (eng *Engine) IsActive() bool {
	return eng.active.Load()
}

func (eng *Engine) AddReceiver(rs ...Listener) {
	for _, r := range rs {
		eng.bus.AddReceiver(r)
	}
}

func (eng *Engine) RemoveReceiver(r Listener) {
	eng.bus.RemoveReceiver(r)
}

func (eng *Engine) Post(em Emittable, priority uint8) {
	if !eng.active.Load() {
		return
	}

	eng.bus.Post(em, priority)
}

func (eng *Engine) bootstrap() {
	eng.bus = newBus(eng.cfg.Demuxers)
	eng.ticker = time.NewTicker(time.Second / time.Duration(eng.cfg.TPS))
	eng.context, eng.contextCancel = context.WithCancel(context.Background())

	for _, c := range eng.cfg.Components {
		eng.AddReceiver(c.Bootstrap()...)
	}
}
