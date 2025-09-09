package banji

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/AndrewChon/banji/bus"
)

const (
	stopPriority     int16 = -32768
	startPriority    int16 = -32768
	tickPrePriority  int16 = -32766
	tickPostPriority int16 = -32767
)

// The Engine brokers communication between decoupled components via Event and Receiver.
type Engine struct {
	options   *Options
	bus       Bus
	active    atomic.Bool
	accepting atomic.Bool
	ticker    *time.Ticker
	loopWg    sync.WaitGroup
	stopLoop  chan struct{}
}

func New(opts ...Option) *Engine {
	eng := &Engine{
		options:  NewOptions(opts...),
		stopLoop: make(chan struct{}, 1),
	}

	eng.bus = bus.NewBus[int16, Event, Receiver](
		bus.WithDemuxers(eng.options.Demuxers),
		bus.WithErrorBuilder(errorBuilder),
	)

	eng.ticker = time.NewTicker((1 * time.Second) / time.Duration(eng.options.TPS))

	for _, c := range eng.options.Components {
		rs, err := c.Bootstrap()

		if err != nil {
			panic(fmt.Sprintf("failed to load component %T: %v\n", c, err))
		}

		for _, r := range rs {
			eng.Subscribe(r)
		}
	}

	return eng
}

func (eng *Engine) Active() bool {
	return eng.active.Load()
}

// Start is a non-blocking operation that starts the engine. Components can listen for StartTopic to be notified when
// this function has been executed.
func (eng *Engine) Start() {
	if eng.active.Load() {
		return
	}

	eng.active.Store(true)
	eng.accepting.Store(true)

	eng.post(new(StartEvent), startPriority)
	go eng.runLoop()
}

// Stop is a blocking operation that gracefully shuts down the engine. Components can listen for StopTopic to be
// notified when this function has been executed.
func (eng *Engine) Stop() {
	if !eng.active.Load() {
		return
	}

	eng.post(new(StopEvent), stopPriority)
	eng.accepting.Store(false)

	eng.ticker.Stop()
	eng.stopLoop <- struct{}{}
	eng.loopWg.Wait()

	for eng.bus.Size() > 0 {
		eng.bus.Tick()
	}

	eng.active.Store(false)
}

// Subscribe registers a Receiver to its associated topic. A Receiver can only be subscribed once. Subsequent calls
// to Subscribe with the same Receiver will have no effect.
func (eng *Engine) Subscribe(r Receiver) {
	r.mark()
	eng.bus.Subscribe(r)
}

// Unsubscribe unregisters a Receiver from its associated topic.
func (eng *Engine) Unsubscribe(r Receiver) {
	eng.bus.Unsubscribe(r)
}

// Post posts an Event to the Engine, which will be handled on the next available tick.
func (eng *Engine) Post(event Event, priority uint8) {
	eng.post(event, int16(priority))
}

func (eng *Engine) post(event Event, priority int16) {
	if eng.accepting.Load() {
		event.mark()
		eng.bus.Post(event, priority)
	}
}

func (eng *Engine) runLoop() {
	for {
		select {
		case tick := <-eng.ticker.C:
			eng.loopWg.Add(1)

			eng.post(&PreTickEvent{
				tick: tick,
			}, tickPrePriority)

			eng.bus.Tick()

			eng.post(&PostTickEvent{
				tick: tick,
			}, tickPostPriority)

			eng.loopWg.Done()
		case <-eng.stopLoop:
			return
		}
	}
}

func errorBuilder(err error) bus.Emittable {
	return &ErrorEvent{
		err: err,
	}
}
