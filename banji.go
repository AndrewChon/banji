// Package banji provides a concurrent proactor and a framework for implementing EDA.
package banji

import (
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
	Handle(em Event) error

	mark()
}

// A Bus is an entity that can receive and route Event types to Receiver types.
type Bus interface {
	Tick()
	Subscribe(r Receiver)
	Unsubscribe(r Receiver)
	Post(e Event, priority uint8)
	Size() int
}
