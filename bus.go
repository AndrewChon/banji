package banji

import (
	"banji/gmap"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

const (
	Lowest    uint8 = 255
	Low       uint8 = 204
	Medium    uint8 = 153
	High      uint8 = 102
	Highest   uint8 = 51
	important uint8 = 0
)

type Emittable interface {
	ID() uuid.UUID
	Postmark() time.Time
	Topic() string
	Canceled() bool
	Cancel()

	mark()
}

type Listener interface {
	ID() uuid.UUID
	Topic() string
	Postmark() time.Time
	Handle(em Emittable) error

	mark()
}

type Event struct {
	id       uuid.UUID
	postmark time.Time
	canceled atomic.Bool
}

func (e *Event) ID() uuid.UUID {
	return e.id
}

func (e *Event) Postmark() time.Time {
	return e.postmark
}

func (e *Event) Canceled() bool {
	return e.canceled.Load()
}

func (e *Event) Cancel() {
	e.canceled.Store(true)
}

func (e *Event) mark() {
	e.id = uuid.New()
	e.postmark = time.Now()
	e.canceled.Store(false)
}

type Receiver struct {
	id       uuid.UUID
	postmark time.Time
}

func (r *Receiver) ID() uuid.UUID {
	return r.id
}

func (r *Receiver) Postmark() time.Time {
	return r.postmark
}

func (r *Receiver) mark() {
	r.id = uuid.New()
	r.postmark = time.Now()
}

type bus struct {
	queue    *Queue[uint8, Emittable]
	bufQueue *Queue[uint8, Emittable]

	rmap gmap.SyncMap[string, []Listener]

	bufQueueCond   *sync.Cond
	bufQueueLocked bool

	demuxWG  sync.WaitGroup
	demuxSem chan struct{}
}

func newBus(demuxers int) *bus {
	return &bus{
		queue:        new(Queue[uint8, Emittable]),
		bufQueue:     new(Queue[uint8, Emittable]),
		bufQueueCond: sync.NewCond(&sync.Mutex{}),

		demuxSem: make(chan struct{}, demuxers),
	}
}

func (b *bus) Tick() {
	b.drainBuffer()
	for e, ok := b.queue.Pop(); ok; e, ok = b.queue.Pop() {
		b.demux(e)
	}
}

func (b *bus) AddReceiver(r Listener) {
	topic := r.Topic()
	if topic == "" {
		return
	}

	r.mark()

	actual, loaded := b.rmap.LoadOrStore(topic, []Listener{r})
	if loaded {
		actual = append(actual, r)
		b.rmap.Store(topic, actual)
	}
}

func (b *bus) RemoveReceiver(r Listener) {
	topic := r.Topic()
	if topic == "" {
		return
	}

	val, ok := b.rmap.Load(topic)
	if !ok {
		return
	}

	for i, x := range val {
		if x == r {
			val = append(val[:i], val[i+1:]...)
			b.rmap.Store(topic, val)
			return
		}
	}
}

func (b *bus) Post(em Emittable, priority uint8) {
	b.bufQueueCond.L.Lock()
	for b.bufQueueLocked {
		b.bufQueueCond.Wait()
	}
	b.bufQueueCond.L.Unlock()

	em.mark()

	b.queue.Push(em, priority)
}

func (b *bus) Size() int {
	return b.queue.Size() + b.bufQueue.Size()
}

func (b *bus) demux(em Emittable) {
	if em.Topic() == "" {
		return
	}

	receivers, ok := b.rmap.Load(em.Topic())
	if !ok {
		return
	}

	for _, r := range receivers {
		b.demuxWG.Add(1)
		go func() {
			b.demuxSem <- struct{}{}
			err := r.Handle(em)
			if err != nil {
				b.Post(NewErrorEvent(err), Lowest)
			}
			<-b.demuxSem
			b.demuxWG.Done()
		}()
	}
	b.demuxWG.Wait()
}

func (b *bus) drainBuffer() {
	b.bufQueueCond.L.Lock()
	b.bufQueueLocked = true
	b.bufQueueCond.L.Unlock()

	b.queue.Meld(b.bufQueue)

	b.bufQueueCond.L.Lock()
	b.bufQueueLocked = false
	b.bufQueueCond.Broadcast()
	b.bufQueueCond.L.Unlock()
}
