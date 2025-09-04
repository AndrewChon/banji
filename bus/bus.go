// Package bus is an implementation of an event bus.
package bus

import (
	"cmp"
	"sync"

	"github.com/AndrewChon/gsync"
	"github.com/AndrewChon/pqueue"
	"github.com/google/uuid"
)

// A PriorityQueue is any data structure that can store and retrieve elements in order of priority.
type PriorityQueue[P cmp.Ordered, V any] interface {
	Push(elem V, priority P)
	Pop() (V, bool)
	Peek() V
	Size() int
	Clear()
}

// An Emittable is any type that can be emitted and routed by a Bus.
type Emittable interface {
	ID() uuid.UUID
	Topic() string
	Cancel()
	Canceled() bool
}

// A Subscriber is any type that can receive and handle Emittable types.
type Subscriber[EM Emittable] interface {
	ID() uuid.UUID
	Topic() string
	Handle(em EM) error
}

type Bus[EM Emittable, SU Subscriber[EM]] struct {
	options *Options

	bufferQueue  PriorityQueue[uint8, EM]
	workingQueue PriorityQueue[uint8, EM]

	wp *workerPool

	subscribeQueue   *pqueue.CircularBuffer[SU]
	unsubscribeQueue *pqueue.CircularBuffer[SU]
	subscribers      gsync.Map[string, []SU]

	bufferQueueMu      sync.Mutex
	subscribeQueueMu   sync.Mutex
	unsubscribeQueueMu sync.Mutex
}

func NewBus[EM Emittable, SU Subscriber[EM]](opts ...Option) *Bus[EM, SU] {
	options := NewOptions(opts...)
	b := &Bus[EM, SU]{
		options:          options,
		bufferQueue:      pqueue.NewPairing[uint8, EM](),
		workingQueue:     pqueue.NewPairing[uint8, EM](),
		subscribeQueue:   pqueue.NewCircularBuffer[SU](),
		unsubscribeQueue: pqueue.NewCircularBuffer[SU](),
		wp:               newWorkerPool(options.Demuxers),
	}

	return b
}

func (b *Bus[EM, SU]) Tick() {
	b.updateSubscribers()

	b.bufferQueueMu.Lock()
	b.workingQueue.(*pqueue.Pairing[uint8, EM]).Meld(b.bufferQueue.(*pqueue.Pairing[uint8, EM]))
	b.bufferQueueMu.Unlock()

	for em, ok := b.workingQueue.Pop(); ok; em, ok = b.workingQueue.Pop() {
		b.demux(em)
	}

	b.wp.wait()
}

func (b *Bus[EM, SU]) Subscribe(s SU) {
	topic := s.Topic()
	if topic == "" {
		return
	}

	b.subscribeQueueMu.Lock()
	defer b.subscribeQueueMu.Unlock()

	b.subscribeQueue.Push(s)
}

func (b *Bus[EM, SU]) Unsubscribe(s SU) {
	topic := s.Topic()
	if topic == "" {
		return
	}

	b.unsubscribeQueueMu.Lock()
	defer b.unsubscribeQueueMu.Unlock()

	b.unsubscribeQueue.Push(s)
}

func (b *Bus[EM, SU]) Post(em EM, priority uint8) {
	b.bufferQueueMu.Lock()
	b.bufferQueue.Push(em, priority)
	b.bufferQueueMu.Unlock()
}

func (b *Bus[EM, SU]) Size() int {
	return b.bufferQueue.Size() + b.workingQueue.Size()
}

func (b *Bus[EM, SU]) demux(em EM) {
	if em.Topic() == "" {
		return
	}

	subs, loaded := b.subscribers.Load(em.Topic())
	if !loaded {
		return
	}

	for _, s := range subs {
		b.wp.post(func() { b.handlingAgent(em, s) })
	}
}

func (b *Bus[EM, SU]) handlingAgent(em EM, s SU) {
	err := s.Handle(em)
	if err == nil {
		return
	}

	errEm := b.options.ErrorBuilder(err)
	if errTyped, ok := errEm.(EM); ok {
		b.Post(errTyped, 0)
	}
}

func (b *Bus[EM, SU]) updateSubscribers() {
	b.unsubscribeQueueMu.Lock()
	b.subscribeQueueMu.Lock()
	defer b.unsubscribeQueueMu.Unlock()
	defer b.subscribeQueueMu.Unlock()

	for s, ok := b.unsubscribeQueue.Pop(); ok; s, ok = b.unsubscribeQueue.Pop() {
		if i, found := b.findSubscriber(s); found {
			subs, _ := b.subscribers.Load(s.Topic())
			subs = append(subs[:i], subs[i+1:]...)
			b.subscribers.Store(s.Topic(), subs)
		}
	}

	for s, ok := b.subscribeQueue.Pop(); ok; s, ok = b.subscribeQueue.Pop() {
		actual, loaded := b.subscribers.Load(s.Topic())
		if !loaded {
			b.subscribers.Store(s.Topic(), []SU{s})
			continue
		}

		// Do not allow duplicate subscribers.
		if _, found := b.findSubscriber(s); found {
			continue
		}

		actual = append(actual, s)
		b.subscribers.Store(s.Topic(), actual)
	}
}

func (b *Bus[EM, SU]) findSubscriber(s SU) (index int, found bool) {
	subs, loaded := b.subscribers.Load(s.Topic())
	if !loaded {
		return 0, false
	}

	for i, x := range subs {
		if x.ID() == s.ID() {
			return i, true
		}
	}

	return 0, false
}
