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

type Bus[K cmp.Ordered, EM Emittable, SU Subscriber[EM]] struct {
	options *Options

	bufferQueue   PriorityQueue[K, EM]
	workingQueue  PriorityQueue[K, EM]
	bufferQueueMu sync.Mutex

	subscribeQueue     *pqueue.CircularBuffer[SU]
	unsubscribeQueue   *pqueue.CircularBuffer[SU]
	subscribers        gsync.Map[string, []SU]
	subscriberQueuesMu sync.Mutex

	wp *workerPool
}

func NewBus[K cmp.Ordered, EM Emittable, SU Subscriber[EM]](opts ...Option) *Bus[K, EM, SU] {
	options := NewOptions(opts...)
	b := &Bus[K, EM, SU]{
		options:          options,
		bufferQueue:      pqueue.NewPairing[K, EM](),
		workingQueue:     pqueue.NewPairing[K, EM](),
		subscribeQueue:   pqueue.NewCircularBuffer[SU](),
		unsubscribeQueue: pqueue.NewCircularBuffer[SU](),
		wp:               newWorkerPool(options.Demuxers),
	}

	return b
}

func (b *Bus[K, EM, SU]) Tick() {
	b.updateSubscribers()

	b.bufferQueueMu.Lock()
	b.workingQueue.(*pqueue.Pairing[K, EM]).Meld(b.bufferQueue.(*pqueue.Pairing[K, EM]))
	b.bufferQueueMu.Unlock()

	for em, ok := b.workingQueue.Pop(); ok; em, ok = b.workingQueue.Pop() {
		b.demux(em)
	}
}

func (b *Bus[K, EM, SU]) Subscribe(s SU) {
	topic := s.Topic()
	if topic == "" {
		return
	}

	b.subscriberQueuesMu.Lock()
	defer b.subscriberQueuesMu.Unlock()

	b.subscribeQueue.Push(s)
}

func (b *Bus[K, EM, SU]) Unsubscribe(s SU) {
	topic := s.Topic()
	if topic == "" {
		return
	}

	b.subscriberQueuesMu.Lock()
	defer b.subscriberQueuesMu.Unlock()

	b.unsubscribeQueue.Push(s)
}

func (b *Bus[K, EM, SU]) Post(em EM, priority K) {
	b.bufferQueueMu.Lock()
	defer b.bufferQueueMu.Unlock()

	b.bufferQueue.Push(em, priority)
}

func (b *Bus[K, EM, SU]) Size() int {
	return b.bufferQueue.Size() + b.workingQueue.Size()
}

func (b *Bus[K, EM, SU]) demux(em EM) {
	if em.Topic() == "" {
		return
	}

	subs, loaded := b.getSubscribers(em.Topic())
	if !loaded {
		return
	}

	for _, s := range subs {
		b.wp.post(func() { b.handlingAgent(em, s) })
	}
	b.wp.wait()
}

func (b *Bus[K, EM, SU]) handlingAgent(em EM, s SU) {
	err := s.Handle(em)
	if err == nil {
		return
	}

	errEm := b.options.ErrorBuilder(err)
	if errTyped, ok := errEm.(EM); ok {
		b.Post(errTyped, K(0))
	}
}

func (b *Bus[K, EM, SU]) updateSubscribers() {
	b.subscriberQueuesMu.Lock()
	defer b.subscriberQueuesMu.Unlock()

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

func (b *Bus[K, EM, SU]) findSubscriber(s SU) (index int, found bool) {
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

func (b *Bus[K, EM, SU]) getSubscribers(topic string) ([]SU, bool) {
	var foundSubscribers []SU

	// Global subscribers.
	if universalSubs, ok := b.subscribers.Load("*"); ok {
		foundSubscribers = append(foundSubscribers, universalSubs...)
	}

	// Retrieve subscribers listening to parent namespaces.
	for i := 0; i < len(topic); i++ {
		if topic[i] == '.' {
			key := topic[:i] + ".*"
			if subs, ok := b.subscribers.Load(key); ok {
				foundSubscribers = append(foundSubscribers, subs...)
			}
		}
	}

	// Finally, retrieve subscribers listening to the specific namespaced topic string.
	if subs, ok := b.subscribers.Load(topic); ok {
		foundSubscribers = append(foundSubscribers, subs...)
	}

	return foundSubscribers, true
}
