package cmap

import "sync"

type lmapString struct {
	m map[string]interface{}
	l sync.RWMutex
}

func (ms *lmapString) Set(key string, v interface{}) {
	ms.l.Lock()
	ms.m[key] = v
	ms.l.Unlock()
}

func (ms *lmapString) Update(key string, fn func(oldVal interface{}) (newVal interface{})) {
	ms.l.Lock()
	ms.m[key] = fn(ms.m[key])
	ms.l.Unlock()
}

func (ms *lmapString) Swap(key string, newV interface{}) (oldV interface{}) {
	ms.l.Lock()
	oldV = ms.m[key]
	ms.m[key] = newV
	ms.l.Unlock()
	return
}

func (ms *lmapString) Get(key string) (v interface{}) {
	ms.l.RLock()
	v = ms.m[key]
	ms.l.RUnlock()
	return
}
func (ms *lmapString) GetOK(key string) (v interface{}, ok bool) {
	ms.l.RLock()
	v, ok = ms.m[key]
	ms.l.RUnlock()
	return
}

func (ms *lmapString) Has(key string) (ok bool) {
	ms.l.RLock()
	_, ok = ms.m[key]
	ms.l.RUnlock()
	return
}

func (ms *lmapString) Delete(key string) {
	ms.l.Lock()
	delete(ms.m, key)
	ms.l.Unlock()
}

func (ms *lmapString) DeleteAndGet(key string) (v interface{}) {
	ms.l.Lock()
	v = ms.m[key]
	delete(ms.m, key)
	ms.l.Unlock()
	return v
}

func (ms *lmapString) Len() (ln int) {
	ms.l.RLock()
	ln = len(ms.m)
	ms.l.RUnlock()
	return
}

func (ms *lmapString) ForEach(fn func(key string, val interface{}) error) (err error) {
	ms.l.RLock()
	keys := make([]string, 0, len(ms.m))
	for key := range ms.m {
		keys = append(keys, key)
	}
	ms.l.RUnlock()

	for _, key := range keys {
		ms.l.RLock()
		val, ok := ms.m[key]
		ms.l.RUnlock()

		if !ok {
			continue
		}

		if err = fn(key, val); err != nil {
			return
		}
	}

	return
}
