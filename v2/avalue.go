package cmap

import (
	"sync"
)

type AtomicValue struct {
	m sync.RWMutex
	v interface{}
}

func (av *AtomicValue) Load() (v interface{}) {
	av.m.RLock()
	v = av.v
	av.m.RUnlock()
	return
}

func (av *AtomicValue) Store(v interface{}) {
	av.m.Lock()
	av.v = v
	av.m.Unlock()
}

func (av *AtomicValue) Swap(newV interface{}) (oldV interface{}) {
	av.m.Lock()
	oldV, av.v = av.v, newV
	av.m.Unlock()
	return
}

func (av *AtomicValue) Update(fn func(oldVal interface{}) interface{}) {
	av.m.Lock()
	av.v = fn(av.v)
	av.m.Unlock()
}

func (av *AtomicValue) CompareAndSwap(fn func(oldV interface{}) (newV interface{}, ok bool)) (ok bool) {
	var newV interface{}
	av.m.Lock()
	if newV, ok = fn(av.v); ok {
		av.v = newV
	}
	av.m.Unlock()
	return
}
