package gsync

import (
	"sync"
)

type Map[K any, V any] struct {
	m sync.Map
}

func (sm *Map[K, V]) Load(key K) (value V, ok bool) {
	v, ok := sm.m.Load(key)
	if !ok {
		return
	}

	value = v.(V)
	return
}

func (sm *Map[K, V]) Store(key K, value V) {
	sm.m.Store(key, value)
}

func (sm *Map[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	v, loaded := sm.m.LoadOrStore(key, value)
	if !loaded {
		actual = value
		return
	}

	actual = v.(V)

	return
}
