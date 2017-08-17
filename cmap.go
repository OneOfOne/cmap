package cmap

import (
	"sync"
	"sync/atomic"
)

// DefaultShardCount is the default number of shards to use when New() or NewFromJSON() are called.
const DefaultShardCount = 1 << 8 // 256

// KeyHasher represents an interface to supply your own type of hashing for keys.
type KeyHasher interface {
	Hash() uint32
}

// CMap is a concurrent safe sharded map to scale on multiple cores.
type CMap struct {
	shards []lmap
	HashFn func(interface{}) uint32
	mod    uint32
}

// New is an alias for NewSize(DefaultShardCount)
func New() *CMap { return NewSize(DefaultShardCount) }

// NewSize returns a CMap with the specific shardSize, note that for performance reasons,
// shardCount must be a power of 2.
// Hash shardCount will improve concurrency but will consume much more memory.
func NewSize(shardCount int) *CMap {
	// must be a power of 2
	if shardCount < 1 {
		shardCount = DefaultShardCount
	} else if shardCount&(shardCount-1) != 0 {
		panic("shardCount must be a power of 2")
	}

	cm := &CMap{
		shards: make([]lmap, shardCount),
		mod:    uint32(shardCount) - 1,
		HashFn: DefaultKeyHasher,
	}

	for i := range cm.shards {
		cm.shards[i].m = make(map[interface{}]interface{})
	}

	return cm
}

func (cm *CMap) shard(key interface{}) *lmap {
	h := cm.HashFn(key)
	return &cm.shards[h&cm.mod]
}

func (cm *CMap) Get(key interface{}) (val interface{}) {
	return cm.shard(key).Get(key)
}

func (cm *CMap) GetOK(key interface{}) (val interface{}, ok bool) {
	return cm.shard(key).GetOK(key)
}

func (cm *CMap) Set(key, val interface{}) {
	cm.shard(key).Set(key, val)
}

func (cm *CMap) Has(key string) bool                 { return cm.shard(key).Has(key) }
func (cm *CMap) Delete(key string)                   { cm.shard(key).Delete(key) }
func (cm *CMap) DeleteAndGet(key string) interface{} { return cm.shard(key).DeleteAndGet(key) }

func (cm *CMap) Update(key string, fn func(oldVal interface{}) (newVal interface{})) {
	cm.shard(key).Update(key, fn)
}

func (cm *CMap) Swap(key string, val interface{}) interface{} {
	return cm.shard(key).Swap(key, val)
}

func (cm *CMap) Keys() []interface{} {
	out := make([]interface{}, 0, cm.Len())
	for i := range cm.shards {
		sh := &cm.shards[i]
		sh.l.RLock()
		for k := range sh.m {
			out = append(out, k)
		}
		sh.l.RUnlock()
	}
	return out
}

func (cm *CMap) ForEach(fn func(key, val interface{}) error) error {
	for i := range cm.shards {
		if err := cm.shards[i].ForEach(fn); err != nil {
			return err
		}
	}
	return nil
}

func (cm *CMap) ForEachParallel(fn func(key, val interface{}) error) error {
	var (
		wg   sync.WaitGroup
		errv atomic.Value
	)
	for i := range cm.shards {
		wg.Add(1)
		go func(i int) {
			cm.shards[i].ForEach(func(k, v interface{}) error {
				if err, _ := errv.Load().(error); err != nil {
					return err
				}

				if err := fn(k, v); err != nil {
					errv.Store(err)
					return err
				}
				return nil
			})
			wg.Done()
		}(i)
	}
	wg.Wait()

	err, _ := errv.Load().(error)
	return err
}

func (cm *CMap) Len() int {
	ln := 0
	for i := range cm.shards {
		ln += cm.shards[i].Len()
	}
	return ln
}
