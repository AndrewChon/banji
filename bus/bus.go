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
	subscribers  gsync.Map[string, []SU]

	bufferQueueCond   *sync.Cond
	bufferQueueLocked bool

	demuxSema chan struct{}
}

func NewBus[EM Emittable, SU Subscriber[EM]](opts ...Option) *Bus[EM, SU] {
	options := NewOptions(opts...)
	b := &Bus[EM, SU]{
		options:         options,
		bufferQueue:     pqueue.NewBinary[uint8, EM](),
		workingQueue:    pqueue.NewBinary[uint8, EM](),
		bufferQueueCond: sync.NewCond(&sync.Mutex{}),
		demuxSema:       make(chan struct{}, options.Demuxers),
	}

	return b
}

func (b *Bus[EM, SU]) Tick() {
	b.lockBufferQueue()
	defer b.unlockBufferQueue()

	// TODO: Honestly, I just got lazy. It is inevitable that I will begrudgingly return to this later.
	b.workingQueue.(*pqueue.Binary[uint8, EM]).Meld(b.bufferQueue.(*pqueue.Binary[uint8, EM]))
	for em, ok := b.workingQueue.Pop(); ok; em, ok = b.workingQueue.Pop() {
		b.demux(em)
	}
}

func (b *Bus[EM, SU]) Subscribe(s SU) {
	topic := s.Topic()
	if topic == "" {
		return
	}

	actual, loaded := b.subscribers.LoadOrStore(topic, []SU{s})
	if !loaded {
		return
	}

	actual = append(actual, s)
	b.subscribers.Store(topic, actual)
}

func (b *Bus[EM, SU]) Unsubscribe(s SU) {
	topic := s.Topic()
	if topic == "" {
		return
	}

	subs, loaded := b.subscribers.Load(topic)
	if !loaded {
		return
	}

	for i, x := range subs {
		if x.ID() != s.ID() {
			continue
		}

		subs = append(subs[:i], subs[i+1:]...)
		b.subscribers.Store(topic, subs)
		return
	}
}

func (b *Bus[EM, SU]) Post(em EM, priority uint8) {
	b.bufferQueueCond.L.Lock()
	defer b.bufferQueueCond.L.Unlock()

	for b.bufferQueueLocked {
		b.bufferQueueCond.Wait()
	}

	b.bufferQueue.Push(em, priority)
}

func (b *Bus[EM, SU]) Size() int {
	return b.bufferQueue.Size() + b.workingQueue.Size()
}

func (b *Bus[EM, SU]) lockBufferQueue() {
	b.bufferQueueCond.L.Lock()
	defer b.bufferQueueCond.L.Unlock()

	b.bufferQueueLocked = true
}

func (b *Bus[EM, SU]) unlockBufferQueue() {
	b.bufferQueueCond.L.Lock()
	defer func() {
		b.bufferQueueCond.L.Unlock()
		b.bufferQueueCond.Broadcast()
	}()

	b.bufferQueueLocked = false
}

func (b *Bus[EM, SU]) demux(em EM) {
	if em.Topic() == "" {
		return
	}

	subs, ok := b.subscribers.Load(em.Topic())
	if !ok {
		return
	}

	for _, sub := range subs {
		b.demuxSema <- struct{}{}
		go func(sub Subscriber[EM]) {
			defer func() { <-b.demuxSema }()
			err := sub.Handle(em)
			if err == nil {
				return
			}

			errEm := b.options.ErrorBuilder(err)
			if errTyped, ok := errEm.(EM); ok {
				b.Post(errTyped, 0)
			}
		}(sub)
	}
}
