package stats

import (
	"banji"
)

type Stats struct {
	ticks      *Ticks
	ptReceiver *PostTickReceiver
}

func New() *Stats {
	ticks := NewTicks()
	return &Stats{
		ticks: ticks,
		ptReceiver: &PostTickReceiver{
			ticks: ticks,
		},
	}
}

func (s *Stats) Ticks() *Ticks {
	return s.ticks
}

func (s *Stats) Bootstrap() []banji.Listener {
	return []banji.Listener{s.ptReceiver}
}
