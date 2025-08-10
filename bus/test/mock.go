package test

import (
	"sync/atomic"

	"banji/bus"

	"github.com/google/uuid"
)

type MockEmittable struct {
	id       uuid.UUID
	canceled atomic.Bool
}

func (e *MockEmittable) ID() uuid.UUID {
	return e.id
}

func (e *MockEmittable) Topic() string {
	return "mock"
}

func (e *MockEmittable) Cancel() {
	e.canceled.Store(true)
}

func (e *MockEmittable) Canceled() bool {
	return e.canceled.Load()
}

type MockSubscriber[EM bus.Emittable] struct {
	id uuid.UUID
}

func (s *MockSubscriber[EM]) ID() uuid.UUID {
	return s.id
}

func (s *MockSubscriber[EM]) Topic() string {
	return "mock"
}

func (s *MockSubscriber[EM]) Handle(_ EM) error {
	return nil
}
