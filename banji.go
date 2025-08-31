// Package banji provides a concurrent proactor and a framework for implementing EDA.
package banji

import (
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

// A Component is a component that provides a collection of receivers.
type Component interface {
	Bootstrap() ([]Receiver, error)
}

// An Event is any type that can be emitted and routed by the engine.
type Event interface {
	ID() uuid.UUID
	Postmark() time.Time
	Topic() string
	Cancel()
	Canceled() bool

	mark()
}

// A Receiver is any type that can receive and handle Event types.
type Receiver interface {
	ID() uuid.UUID
	Postmark() time.Time
	Topic() string
	Handle(event Event) error

	mark()
}

// A Bus is an entity that can receive and route Event types to Receiver types.
type Bus interface {
	Tick()
	Subscribe(r Receiver)
	Unsubscribe(r Receiver)
	Post(event Event, priority uint8)
	Size() int
}

// EventEmbed contains internal methods required to implement the Event interface.
type EventEmbed struct {
	id       uuid.UUID
	postmark time.Time
	canceled atomic.Bool
}

func (e *EventEmbed) ID() uuid.UUID {
	return e.id
}

func (e *EventEmbed) Postmark() time.Time {
	return e.postmark
}

func (e *EventEmbed) Cancel() {
	e.canceled.Store(true)
}

func (e *EventEmbed) Canceled() bool {
	return e.canceled.Load()
}

func (e *EventEmbed) mark() {
	e.postmark = time.Now()
	e.id = uuid.New()
	e.canceled.Store(false)
}

// ReceiverEmbed contains internal methods required to implement the Receiver interface.
type ReceiverEmbed struct {
	id       uuid.UUID
	postmark time.Time
}

func (r *ReceiverEmbed) ID() uuid.UUID {
	return r.id
}

func (r *ReceiverEmbed) Postmark() time.Time {
	return r.postmark
}

func (r *ReceiverEmbed) mark() {
	r.postmark = time.Now()
	r.id = uuid.New()
}
