package banji

import (
	"time"

	"github.com/google/uuid"
)

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
	r.id = uuid.New()
	r.postmark = time.Now()
}
