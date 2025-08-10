package banji

import (
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

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
	e.id = uuid.New()
	e.postmark = time.Now()
	e.canceled.Store(false)
}
