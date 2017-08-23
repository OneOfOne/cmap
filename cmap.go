// +build genx

package cmap

import (
	"context"
	"sync"

	"github.com/OneOfOne/cmap/hashers"
)

type (
	KT interface{} // nolint
	VT interface{} // nolint
)

// DefaultShardCount is the default number of shards to use when New() or NewFromJSON() are called. The default is 256.
const DefaultShardCount = 1 << 8

// CMap is a concurrent safe sharded map to scale on multiple cores.
type CMap struct {
	shards   []*LMap
	keysPool sync.Pool
	// HashFn allows using a custom hash function that's used to determain the key's shard. Defaults to DefaultKeyHasher.
	HashFn func(KT) uint32
}

// New is an alias for NewSize(DefaultShardCount)
func New() *CMap { return NewSize(DefaultShardCount) }

// NewSize returns a CMap with the specific shardSize, note that for performance reasons,
// shardCount must be a power of 2.
// Higher shardCount will improve concurrency but will consume more memory.
func NewSize(shardCount int) *CMap {
	// must be a power of 2
	if shardCount < 1 {
		shardCount = DefaultShardCount
	} else if shardCount&(shardCount-1) != 0 {
		panic("shardCount must be a power of 2")
	}

	cm := &CMap{
		shards: make([]*LMap, shardCount),
		HashFn: DefaultKeyHasher,
	}

	cm.keysPool.New = func() interface{} {
		out := make([]KT, 0, DefaultShardCount) // good starting round

		return &out // return a ptr to avoid extra allocation on Get/Put
	}

	for i := range cm.shards {
		cm.shards[i] = NewLMapSize(shardCount)
	}

	return cm
}

// ShardForKey returns the LMap that may hold the specific key.
func (cm *CMap) ShardForKey(key KT) *LMap {
	h := cm.HashFn(key)
	return cm.shards[h&uint32(len(cm.shards)-1)]
}

// Set is the equivalent of `map[key] = val`.
func (cm *CMap) Set(key KT, val VT) {
	h := cm.HashFn(key)
	cm.shards[h&uint32(len(cm.shards)-1)].Set(key, val)
}

// SetIfNotExists will only assign val to key if it wasn't already set.
// Use `Update` if you need more logic.
func (cm *CMap) SetIfNotExists(key KT, val VT) (set bool) {
	h := cm.HashFn(key)
	return cm.shards[h&uint32(len(cm.shards)-1)].SetIfNotExists(key, val)
}

// Get is the equivalent of `val := map[key]`.
func (cm *CMap) Get(key KT) (val VT) {
	h := cm.HashFn(key)
	return cm.shards[h&uint32(len(cm.shards)-1)].Get(key)
}

// GetOK is the equivalent of `val, ok := map[key]`.
func (cm *CMap) GetOK(key KT) (val VT, ok bool) {
	h := cm.HashFn(key)
	return cm.shards[h&uint32(len(cm.shards)-1)].GetOK(key)
}

// Has is the equivalent of `_, ok := map[key]`.
func (cm *CMap) Has(key KT) bool {
	h := cm.HashFn(key)
	return cm.shards[h&uint32(len(cm.shards)-1)].Has(key)
}

// Delete is the equivalent of `delete(map, key)`.
func (cm *CMap) Delete(key KT) {
	h := cm.HashFn(key)
	cm.shards[h&uint32(len(cm.shards)-1)].Delete(key)
}

// DeleteAndGet is the equivalent of `oldVal := map[key]; delete(map, key)`.
func (cm *CMap) DeleteAndGet(key KT) VT {
	h := cm.HashFn(key)
	return cm.shards[h&uint32(len(cm.shards)-1)].DeleteAndGet(key)
}

// Update calls `fn` with the key's old value (or nil) and assign the returned value to the key.
// The shard containing the key will be locked, it is NOT safe to call other cmap funcs inside `fn`.
func (cm *CMap) Update(key KT, fn func(oldval VT) (newval VT)) {
	h := cm.HashFn(key)
	cm.shards[h&uint32(len(cm.shards)-1)].Update(key, fn)
}

// Swap is the equivalent of `oldVal, map[key] = map[key], newVal`.
func (cm *CMap) Swap(key KT, val VT) VT {
	h := cm.HashFn(key)
	return cm.shards[h&uint32(len(cm.shards)-1)].Swap(key, val)
}

// Keys returns a slice of all the keys of the map.
func (cm *CMap) Keys() []KT {
	out := make([]KT, 0, cm.Len())
	for _, sh := range cm.shards {
		out = sh.Keys(out)
	}
	return out
}

// ForEach loops over all the key/values in the map.
// You can break early by returning false.
// It **is** safe to modify the map while using this iterator, however it uses more memory and is slightly slower.
func (cm *CMap) ForEach(fn func(key KT, val VT) bool) bool {
	keysP := cm.keysPool.Get().(*[]KT)
	defer cm.keysPool.Put(keysP)

	for _, lm := range cm.shards {
		keys := (*keysP)[:0]
		if !lm.ForEach(keys, fn) {
			return false
		}
	}

	return false
}

// ForEachLocked loops over all the key/values in the map.
// You can break early by returning false.
// It is **NOT* safe to modify the map while using this iterator.
func (cm *CMap) ForEachLocked(fn func(key KT, val VT) bool) bool {
	for _, lm := range cm.shards {
		if !lm.ForEachLocked(fn) {
			return false
		}
	}

	return true
}

// Len returns the length of the map.
func (cm *CMap) Len() int {
	ln := 0
	for _, lm := range cm.shards {
		ln += lm.Len()
	}
	return ln
}

// KV holds the key/value returned when Iter is called.
type KV struct {
	Key   KT
	Value VT
}

// Iter returns a channel to be used in for range.
// Use `context.WithCancel` if you intend to break early or goroutines will leak.
// It **is** safe to modify the map while using this iterator, however it uses more memory and is slightly slower.
func (cm *CMap) Iter(ctx context.Context, buffer int) <-chan *KV {
	ch := make(chan *KV, buffer)
	go func() {
		cm.iterContext(ctx, ch, false)
		close(ch)
	}()
	return ch
}

// IterLocked returns a channel to be used in for range.
// Use `context.WithCancel` if you intend to break early or goroutines will leak and map access will deadlock.
// It is **NOT* safe to modify the map while using this iterator.
func (cm *CMap) IterLocked(ctx context.Context, buffer int) <-chan *KV {
	ch := make(chan *KV, buffer)
	go func() {
		cm.iterContext(ctx, ch, false)
		close(ch)
	}()
	return ch
}

// iterContext is used internally
func (cm *CMap) iterContext(ctx context.Context, ch chan<- *KV, locked bool) {
	fn := func(k KT, v VT) bool {
		select {
		case <-ctx.Done():
			return false
		case ch <- &KV{k, v}:
			return true
		}
	}

	if locked {
		_ = cm.ForEachLocked(fn)
	} else {
		_ = cm.ForEach(fn)
	}
}

// NumShards returns the number of shards in the map.
func (cm *CMap) NumShards() int { return len(cm.shards) }

// DefaultKeyHasher is an alias for hashers.TypeHasher32(key).
func DefaultKeyHasher(key KT) uint32 { return hashers.TypeHasher32(key) }
