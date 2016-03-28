package cmap

import (
	"runtime"
	"sync/atomic"
)

var nilValue *KeyValue

type aValue struct {
	v      *KeyValue
	lk     uint32
	hasVal bool
}

func (a *aValue) lock() {
	for !atomic.CompareAndSwapUint32(&a.lk, 0, 1) {
		runtime.Gosched()
	}
}

func (a *aValue) unlock() { atomic.StoreUint32(&a.lk, 0) }

func (a *aValue) Store(v *KeyValue) {
	a.lock()
	a.v = v
	a.unlock()
}

func (a *aValue) Load() *KeyValue {
	a.lock()
	v := a.v
	a.unlock()
	return v
}

func (a *aValue) CompareAndSwapIfNil(newVal *KeyValue) bool {
	var b bool
	a.lock()
	if b = !a.hasVal; b {
		a.v, a.hasVal = newVal, true
	}
	a.unlock()
	return b
}
func (a *aValue) SwapWithNil() (*KeyValue, bool) {
	var (
		v  *KeyValue
		ok bool
	)
	a.lock()
	v, a.v, ok, a.hasVal = a.v, nilValue, a.hasVal, false
	a.unlock()
	return v, ok
}
