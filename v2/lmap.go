package cmap

import "sync"


type smap struct {
	m map[interface{}]interface{}
	l sync.RWMutex
}

func (ms *smap) Set(key interface{}, v interface{}) {
	ms.l.Lock()
	ms.m[key] = v
	ms.l.Unlock()
}

func (ms *smap) Update(key interface{}, fn func(oldVal interface{}) (newVal interface{})) {
	ms.l.Lock()
	ms.m[key] = fn(ms.m[key])
	ms.l.Unlock()
}

func (ms *smap) Swap(key interface{}, newV interface{}) (oldV interface{}) {
	ms.l.Lock()
	oldV = ms.m[key]
	ms.m[key] = newV
	ms.l.Unlock()
	return
}

func (ms *smap) Get(key interface{}) (v interface{}) {
	ms.l.RLock()
	v = ms.m[key]
	ms.l.RUnlock()
	return
}
func (ms *smap) GetOK(key interface{}) (v interface{}, ok bool) {
	ms.l.RLock()
	v, ok = ms.m[key]
	ms.l.RUnlock()
	return
}

func (ms *smap) Has(key interface{}) (ok bool) {
	ms.l.RLock()
	_, ok = ms.m[key]
	ms.l.RUnlock()
	return
}

func (ms *smap) Delete(key interface{}) {
	ms.l.Lock()
	delete(ms.m, key)
	ms.l.Unlock()
}

func (ms *smap) DeleteAndGet(key interface{}) (v interface{}) {
	ms.l.Lock()
	v = ms.m[key]
	delete(ms.m, key)
	ms.l.Unlock()
	return v
}

func (ms *smap) Len() (ln int) {
	ms.l.RLock()
	ln = len(ms.m)
	ms.l.RUnlock()
	return
}

func (ms *smap) ForEach(fn func(key, val interface{}) error) (err error) {
	ms.l.RLock()
	keys := make([]interface{}, 0, len(ms.m))
	for key := range ms.m {
		keys = append(keys, key)
	}
	ms.l.RUnlock()

	for _, key := range keys {
		ms.l.RLock()
		val := ms.m[key]
		ms.l.RUnlock()
		if err = fn(key, val); err != nil {
			return
		}
	}

	return
}
