package cmap

import "sync"

type lockedMap struct {
	m map[string]interface{}
	l sync.RWMutex
}

func (ms *lockedMap) Set(key string, v interface{}) {
	ms.l.Lock()
	ms.m[key] = v
	ms.l.Unlock()
}

func (ms *lockedMap) Update(key string, fn UpdateFunc) {
	ms.l.Lock()

	if v := fn(ms.m[key]); v != DeleteValue {
		ms.m[key] = v
	} else {
		delete(ms.m, key)
	}
	ms.l.Unlock()
}

func (ms *lockedMap) Swap(key string, v interface{}) interface{} {
	ms.l.Lock()
	ov := ms.m[key]
	ms.m[key] = v
	ms.l.Unlock()
	return ov
}

func (ms *lockedMap) CompareAndSwap(key string, v interface{}, casFn CompareAndSwapFunc) bool {
	ms.l.Lock()
	ov := ms.m[key]
	ok := casFn(v, ov)
	if ok {
		ms.m[key] = v
	}
	ms.l.Unlock()
	return ok
}

func (ms *lockedMap) Get(key string) interface{} {
	ms.l.RLock()
	v := ms.m[key]
	ms.l.RUnlock()
	return v
}

func (ms *lockedMap) Has(key string) bool {
	ms.l.RLock()
	_, ok := ms.m[key]
	ms.l.RUnlock()
	return ok
}

func (ms *lockedMap) Delete(key string) {
	ms.l.Lock()
	delete(ms.m, key)
	ms.l.Unlock()
}

func (ms *lockedMap) DeleteAndGet(key string) interface{} {
	ms.l.Lock()
	v := ms.m[key]
	delete(ms.m, key)
	ms.l.Unlock()
	return v
}

func (ms *lockedMap) Len() int {
	ms.l.RLock()
	ln := len(ms.m)
	ms.l.RUnlock()
	return ln
}

func (ms *lockedMap) ForEach(fn ForEachFunc) {
	ms.l.RLock()
	for k, v := range ms.m {
		if fn(k, v) {
			break
		}
	}
	ms.l.RUnlock()
}

func (ms *lockedMap) iter(ch KeyValueChan, wg *sync.WaitGroup) {
	var kv KeyValue
	ms.l.RLock()
	for k, v := range ms.m {
		kv.Key, kv.Value = k, v
		if !ch.send(&kv) {
			break
		}
	}
	ms.l.RUnlock()
	wg.Done()
}
