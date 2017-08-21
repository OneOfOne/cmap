package cmap

import "sync"

type LockedMap struct {
	m map[KT]VT
	l *sync.RWMutex
}

func NewLockedMap() *LockedMap {
	return NewLockedMapSize(0)
}

func NewLockedMapSize(cap int) *LockedMap {
	return &LockedMap{
		m: make(map[KT]VT, cap),
		l: new(sync.RWMutex),
	}
}

func (lm LockedMap) Set(key KT, v VT) {
	lm.l.Lock()
	lm.m[key] = v
	lm.l.Unlock()
}

func (lm LockedMap) Update(key KT, fn func(oldVal VT) (newVal VT)) {
	lm.l.Lock()
	lm.m[key] = fn(lm.m[key])
	lm.l.Unlock()
}

func (lm LockedMap) Swap(key KT, newV VT) (oldV VT) {
	lm.l.Lock()
	oldV = lm.m[key]
	lm.m[key] = newV
	lm.l.Unlock()
	return
}

func (lm LockedMap) Get(key KT) (v VT) {
	lm.l.RLock()
	v = lm.m[key]
	lm.l.RUnlock()
	return
}
func (lm LockedMap) GetOK(key KT) (v VT, ok bool) {
	lm.l.RLock()
	v, ok = lm.m[key]
	lm.l.RUnlock()
	return
}

func (lm LockedMap) Has(key KT) (ok bool) {
	lm.l.RLock()
	_, ok = lm.m[key]
	lm.l.RUnlock()
	return
}

func (lm LockedMap) Delete(key KT) {
	lm.l.Lock()
	delete(lm.m, key)
	lm.l.Unlock()
}

func (lm LockedMap) DeleteAndGet(key KT) (v VT) {
	lm.l.Lock()
	v = lm.m[key]
	delete(lm.m, key)
	lm.l.Unlock()
	return v
}

func (lm LockedMap) Len() (ln int) {
	lm.l.RLock()
	ln = len(lm.m)
	lm.l.RUnlock()
	return
}

func (lm LockedMap) ForEach(keys []KT, fn func(key KT, val VT) error) (err error) {
	lm.l.RLock()
	for key := range lm.m {
		keys = append(keys, key)
	}
	lm.l.RUnlock()

	for _, key := range keys {
		lm.l.RLock()
		val, ok := lm.m[key]
		lm.l.RUnlock()
		if !ok {
			continue
		}
		if err = fn(key, val); err != nil {
			return
		}
	}

	return
}

func (lm LockedMap) ForEachLocked(fn func(key KT, val VT) error) (err error) {
	lm.l.RLock()
	defer lm.l.RUnlock()

	for key, val := range lm.m {
		if err = fn(key, val); err != nil {
			return
		}
	}

	return
}
