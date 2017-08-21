package cmap

import "sync"

// LMap is a simple sync.RWMutex locked map.
// Used by CMap internally for sharding.
type LMap struct {
	m map[KT]VT
	l *sync.RWMutex
}

// NewLMap returns a new LMap with the cap set to 0.
func NewLMap() *LMap {
	return NewLMapSize(0)
}

// NewLMapSize is the equivalent of `m := make(map[KT]VT, cap)`
func NewLMapSize(cap int) *LMap {
	return &LMap{
		m: make(map[KT]VT, cap),
		l: new(sync.RWMutex),
	}
}

// Set is the equivalent of `map[key] = val`.
func (lm *LMap) Set(key KT, v VT) {
	lm.l.Lock()
	lm.m[key] = v
	lm.l.Unlock()
}

// SetIfNotExists will only assign val to key if it wasn't already set.
// Use `Update` if you need more logic.
func (lm *LMap) SetIfNotExists(key KT, val VT) (set bool) {
	lm.l.Lock()
	if _, ok := lm.m[key]; !ok {
		lm.m[key], set = val, true
	}
	lm.l.Unlock()
	return
}

// Get is the equivalent of `val := map[key]`.
func (lm *LMap) Get(key KT) (v VT) {
	lm.l.RLock()
	v = lm.m[key]
	lm.l.RUnlock()
	return
}

// GetOK is the equivalent of `val, ok := map[key]`.
func (lm *LMap) GetOK(key KT) (v VT, ok bool) {
	lm.l.RLock()
	v, ok = lm.m[key]
	lm.l.RUnlock()
	return
}

// Has is the equivalent of `_, ok := map[key]`.
func (lm *LMap) Has(key KT) (ok bool) {
	lm.l.RLock()
	_, ok = lm.m[key]
	lm.l.RUnlock()
	return
}

// Delete is the equivalent of `delete(map, key)`.
func (lm *LMap) Delete(key KT) {
	lm.l.Lock()
	delete(lm.m, key)
	lm.l.Unlock()
}

// DeleteAndGet is the equivalent of `oldVal := map[key]; delete(map, key)`.
func (lm *LMap) DeleteAndGet(key KT) (v VT) {
	lm.l.Lock()
	v = lm.m[key]
	delete(lm.m, key)
	lm.l.Unlock()
	return v
}

// Update calls `fn` with the key's old value (or nil) and assigns the returned value to the key.
// The shard containing the key will be locked, it is NOT safe to call other cmap funcs inside `fn`.
func (lm *LMap) Update(key KT, fn func(oldVal VT) (newVal VT)) {
	lm.l.Lock()
	lm.m[key] = fn(lm.m[key])
	lm.l.Unlock()
}

// Swap is the equivalent of `oldVal, map[key] = map[key], newVal`.
func (lm *LMap) Swap(key KT, newV VT) (oldV VT) {
	lm.l.Lock()
	oldV = lm.m[key]
	lm.m[key] = newV
	lm.l.Unlock()
	return
}

// ForEach loops over all the key/values in the map.
// You can break early by returning an error .
// It **is** safe to modify the map while using this iterator, however it uses more memory and is slightly slower.
func (lm *LMap) ForEach(keys []KT, fn func(key KT, val VT) bool) bool {
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
		if !fn(key, val) {
			return false
		}
	}

	return true
}

// ForEachLocked loops over all the key/values in the map.
// You can break early by returning false
// It is **NOT* safe to modify the map while using this iterator.
func (lm *LMap) ForEachLocked(fn func(key KT, val VT) bool) bool {
	lm.l.RLock()
	defer lm.l.RUnlock()

	for key, val := range lm.m {
		if !fn(key, val) {
			return false
		}
	}

	return true
}

// Len returns the length of the map.
func (lm *LMap) Len() (ln int) {
	lm.l.RLock()
	ln = len(lm.m)
	lm.l.RUnlock()
	return
}

// Keys appends all the keys in the map to buf and returns buf.
// buf may be nil.
func (lm *LMap) Keys(buf []KT) []KT {
	lm.l.RLock()
	if cap(buf) == 0 {
		buf = make([]KT, 0, len(lm.m))
	}
	for k := range lm.m {
		buf = append(buf, k)
	}
	lm.l.RUnlock()
	return buf
}
