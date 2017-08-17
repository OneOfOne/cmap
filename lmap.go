package cmap

import "sync"

type lmap struct {
	m map[KT]VT
	l sync.RWMutex
}

func (ms *lmap) Set(key KT, v VT) {
	ms.l.Lock()
	ms.m[key] = v
	ms.l.Unlock()
}

func (ms *lmap) Update(key KT, fn func(oldVal VT) (newVal VT)) {
	ms.l.Lock()
	ms.m[key] = fn(ms.m[key])
	ms.l.Unlock()
}

func (ms *lmap) Swap(key KT, newV VT) (oldV VT) {
	ms.l.Lock()
	oldV = ms.m[key]
	ms.m[key] = newV
	ms.l.Unlock()
	return
}

func (ms *lmap) Get(key KT) (v VT) {
	ms.l.RLock()
	v = ms.m[key]
	ms.l.RUnlock()
	return
}
func (ms *lmap) GetOK(key KT) (v VT, ok bool) {
	ms.l.RLock()
	v, ok = ms.m[key]
	ms.l.RUnlock()
	return
}

func (ms *lmap) Has(key KT) (ok bool) {
	ms.l.RLock()
	_, ok = ms.m[key]
	ms.l.RUnlock()
	return
}

func (ms *lmap) Delete(key KT) {
	ms.l.Lock()
	delete(ms.m, key)
	ms.l.Unlock()
}

func (ms *lmap) DeleteAndGet(key KT) (v VT) {
	ms.l.Lock()
	v = ms.m[key]
	delete(ms.m, key)
	ms.l.Unlock()
	return v
}

func (ms *lmap) Len() (ln int) {
	ms.l.RLock()
	ln = len(ms.m)
	ms.l.RUnlock()
	return
}

func (ms *lmap) ForEach(fn func(key KT, val VT) error) (err error) {
	ms.l.RLock()
	keys := make([]KT, 0, len(ms.m))
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

func (ms *lmap) ForEachLocked(fn func(key KT, val VT) error) (err error) {
	ms.l.RLock()
	defer ms.l.RUnlock()

	for key, val := range ms.m {
		if err = fn(key, val); err != nil {
			return
		}
	}

	return
}
