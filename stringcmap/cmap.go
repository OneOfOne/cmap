package stringcmap

import (
	"sync"
	"sync/atomic"

	"github.com/OneOfOne/cmap"
)

// DefaultShardCount is the default number of shards to use when New() or NewFromJSON() are called.
const DefaultShardCount = cmap.DefaultShardCount

// CMap is a concurrent safe sharded map to scale on multiple cores.
type CMap struct {
	shards []lmap
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
	}

	for i := range cm.shards {
		cm.shards[i].m = make(map[string]interface{})
	}

	return cm
}

func (cm *CMap) shard(key string) *lmap {
	h := fnv32(key)
	return &cm.shards[h&cm.mod]
}

func (cm *CMap) Get(key string) (val interface{}) {
	return cm.shard(key).Get(key)
}

func (cm *CMap) Set(key string, val interface{}) {
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

func (cm *CMap) Foreach(fn func(key string, val interface{}) error) error {
	for i := range cm.shards {
		if err := cm.shards[i].ForEach(fn); err != nil {
			return err
		}
	}
	return nil
}

func (cm *CMap) ForEachParallel(fn func(key string, val interface{}) error) error {
	var (
		wg   sync.WaitGroup
		errv atomic.Value
	)
	for i := range cm.shards {
		wg.Add(1)
		go func(i int) {
			cm.shards[i].ForEach(func(k string, v interface{}) error {
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

func fnv32(key string) uint32 {
	const prime32 = uint32(16777619)
	hash := uint32(2166136261)
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}
