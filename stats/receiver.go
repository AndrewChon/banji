package stats

import (
	"banji"
)

type PostTickReceiver struct {
	banji.Receiver
	ticks *Ticks
}

func (r *PostTickReceiver) Topic() string {
	return banji.PostTickTopic
}

func (r *PostTickReceiver) Handle(em banji.Emittable) error {
	e, _ := em.(*banji.PostTickEvent)
	r.ticks.in(e.Tick())
	return nil
}
