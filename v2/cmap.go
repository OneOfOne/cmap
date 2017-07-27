package cmap

import (
	"fmt"
	"math"
	"reflect"
	"sync"
	"sync/atomic"

	"github.com/OneOfOne/xxhash"
)

// DefaultShardCount is the default number of shards to use when New() or NewFromJSON() are called.
const DefaultShardCount = 1 << 4 // 16

var (
	hashers    = map[reflect.Type]func(key interface{}) uint64{}
	hashersMux sync.RWMutex
)

func RegisterHasher(t reflect.Type, fn func(key interface{}) uint64) {
	hashersMux.Lock()
	hashers[t] = fn
	hashersMux.Unlock()
}

type CMap struct {
	shards []smap
	l      uint64
	HashFn func(interface{}) uint64
}

// New is an alias for NewSize(DefaultShardCount)
func New() *CMap { return NewSize(DefaultShardCount) }

// NewSize returns a CMap with the specific shardSize, note that for performance reasons,
// shardCount must be a power of 2
func NewSize(shardCount int) *CMap {
	// must be a power of 2
	if shardCount == 0 {
		shardCount = DefaultShardCount
	} else if shardCount&(shardCount-1) != 0 {
		panic("shardCount must be a power of 2")
	}
	cm := CMap{
		shards: make([]smap, shardCount),
		l:      uint64(shardCount) - 1,
	}
	for i := range cm.shards {
		cm.shards[i].m = make(map[interface{}]interface{})
	}
	cm.HashFn = DefaultHasher
	return &cm
}

func (cm *CMap) shard(key interface{}) *smap {
	h := cm.HashFn(key)
	return &cm.shards[h&cm.l]
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

func DefaultHasher(key interface{}) uint64 {
	switch key := key.(type) {
	case string:
		return xxhash.ChecksumString64(key)
	case uint64:
		return key
	case int64:
		return uint64(key)
	case float64:
		return math.Float64bits(key)
	case float32:
		return uint64(math.Float32bits(key))
	case int:
		return uint64(key)
	case uint:
		return uint64(key)
	default:
		hashersMux.RLock()
		fn, ok := hashers[reflect.TypeOf(key)]
		hashersMux.RUnlock()
		if ok {
			return fn(key)
		}
		panic(fmt.Sprintf("unsupported type: %T (%v)", key, key))
	}
}
