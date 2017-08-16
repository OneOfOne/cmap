package stringcmap

import "sync"

type lmap struct {
	m map[string]interface{}
	l sync.RWMutex
}

func (ms *lmap) Set(key string, v interface{}) {
	ms.l.Lock()
	ms.m[key] = v
	ms.l.Unlock()
}

func (ms *lmap) Update(key string, fn func(oldVal interface{}) (newVal interface{})) {
	ms.l.Lock()
	ms.m[key] = fn(ms.m[key])
	ms.l.Unlock()
}

func (ms *lmap) Swap(key string, newV interface{}) (oldV interface{}) {
	ms.l.Lock()
	oldV = ms.m[key]
	ms.m[key] = newV
	ms.l.Unlock()
	return
}

func (ms *lmap) Get(key string) (v interface{}) {
	ms.l.RLock()
	v = ms.m[key]
	ms.l.RUnlock()
	return
}
func (ms *lmap) GetOK(key string) (v interface{}, ok bool) {
	ms.l.RLock()
	v, ok = ms.m[key]
	ms.l.RUnlock()
	return
}

func (ms *lmap) Has(key string) (ok bool) {
	ms.l.RLock()
	_, ok = ms.m[key]
	ms.l.RUnlock()
	return
}

func (ms *lmap) Delete(key string) {
	ms.l.Lock()
	delete(ms.m, key)
	ms.l.Unlock()
}

func (ms *lmap) DeleteAndGet(key string) (v interface{}) {
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

func (ms *lmap) Keys() []string {
	ms.l.RLock()
	keys := make([]string, 0, len(ms.m))
	for key := range ms.m {
		keys = append(keys, key)
	}
	ms.l.RUnlock()
	return keys
}

func (ms *lmap) ForEach(fn func(key string, val interface{}) error) (err error) {
	for _, key := range ms.Keys() {
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
