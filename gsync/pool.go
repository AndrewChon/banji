package gsync

import (
	"sync"
)

type Pool[V any] struct {
	p sync.Pool
}

func (sp *Pool[V]) Put(x V) {
	sp.p.Put(x)
}

func (sp *Pool[V]) Get() (x V, ok bool) {
	v := sp.p.Get()
	if v == nil {
		return x, false
	}

	return v.(V), true
}
