package cmap

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
)

// DefaultShardCount is the default number of shards to use when New() or NewFromJSON() are called.
const DefaultShardCount = 1 << 9 // 512

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
// shardCount must be a power of 2
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

func (cm *CMap) Foreach(fn func(key, val interface{}) error) error {
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

// DefaultKeyHasher returns a hash for the specific key for internal sharding.
// By default, those types are supported as keys: string, uint64, int64, uint32, int32, uint, int,
//  float64, float32 and KeyHasher.
func DefaultKeyHasher(key interface{}) uint32 {
	switch key := key.(type) {
	case string:
		return fnv32(key)
	case uint64:
		return uint32(key)
	case int64:
		return uint32(key)
	case float64:
		return uint32(math.Float64bits(key))
	case float32:
		return uint32(math.Float32bits(key))
	case int:
		return uint32(key)
	case uint:
		return uint32(key)
	case KeyHasher:
		return key.Hash()
	default:
		panic(fmt.Sprintf("unsupported type: %T (%v)", key, key))
	}
}
