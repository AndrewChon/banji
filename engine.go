package banji

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/AndrewChon/banji/bus"
)

// The Engine brokers communication between decoupled components via Event and Receiver.
type Engine struct {
	options  *Options
	bus      Bus
	active   atomic.Bool
	ticker   *time.Ticker
	loopWg   sync.WaitGroup
	stopLoop chan struct{}
}

func New(opts ...Option) *Engine {
	eng := &Engine{
		options:  NewOptions(opts...),
		stopLoop: make(chan struct{}, 1),
	}

	eng.bus = bus.NewBus[Event, Receiver](
		bus.WithDemuxers(8),
		bus.WithErrorBuilder(errorBuilder),
	)

	for _, component := range eng.options.Components {
		for _, r := range component.Bootstrap() {
			eng.bus.Subscribe(r)
		}
	}

	eng.ticker = time.NewTicker((1 * time.Second) / time.Duration(eng.options.TPS))

	return eng
}

func (eng *Engine) Active() bool {
	return eng.active.Load()
}

// Start starts the engine. Start is a non-blocking operation.
//
// Receivers can subscribe to StartTopic to listen for this operation.
func (eng *Engine) Start() {
	if eng.active.Load() {
		return
	}

	eng.active.Store(true)
	eng.Post(new(StartEvent), 0)

	go func() {
		for {
			select {
			case tick := <-eng.ticker.C:
				eng.loopWg.Add(1)

				eng.Post(&PreTickEvent{
					tick: tick,
				}, 0)
				eng.bus.Tick()
				eng.Post(&PostTickEvent{
					tick: tick,
				}, 0)

				eng.loopWg.Done()
			case <-eng.stopLoop:
				return
			}
		}
	}()
}

// Stop gracefully shuts down the engine, and thus, is a blocking operation.
//
// Receivers can subscribe to StopTopic to listen for this operation.
func (eng *Engine) Stop() {
	if !eng.active.Load() {
		return
	}

	eng.Post(new(StopEvent), 0)
	eng.active.Store(false) // Prevent any new incoming events.

	eng.ticker.Stop()
	eng.stopLoop <- struct{}{}
	eng.loopWg.Wait()

	for eng.bus.Size() > 0 {
		eng.bus.Tick()
	}
}

// Subscribe registers a Receiver.
func (eng *Engine) Subscribe(r Receiver) {
	r.mark()
	eng.bus.Subscribe(r)
}

// Unsubscribe unregisters a Receiver. Note that this can be an expensive operation, and thus it should be used
// sparingly.
func (eng *Engine) Unsubscribe(r Receiver) {
	eng.bus.Unsubscribe(r)
}

// Post posts an Event to the engine. The Event will be handled on the next tick.
func (eng *Engine) Post(event Event, priority uint8) {
	if eng.active.Load() {
		event.mark()
		eng.bus.Post(event, priority)
	}
}

func errorBuilder(err error) bus.Emittable {
	return &ErrorEvent{
		err: err,
	}
}
