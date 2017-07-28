package cmap

import (
	"sync"
	"sync/atomic"
)

type CMapString struct {
	shards []lmapString
	mod    uint32
}

// New is an alias for NewSize(DefaultShardCount)
func NewString() *CMapString { return NewSizeString(DefaultShardCount) }

// NewSize returns a CMap with the specific shardSize, note that for performance reasons,
// shardCount must be a power of 2
func NewSizeString(shardCount int) *CMapString {
	// must be a power of 2
	if shardCount == 0 {
		shardCount = DefaultShardCount
	} else if shardCount&(shardCount-1) != 0 {
		panic("shardCount must be a power of 2")
	}

	cm := &CMapString{
		shards: make([]lmapString, shardCount),
		mod:    uint32(shardCount) - 1,
	}

	for i := range cm.shards {
		cm.shards[i].m = make(map[string]interface{})
	}

	return cm
}

func (cm *CMapString) shard(key string) *lmapString {
	return &cm.shards[fnv32(key)&cm.mod]
}

func (cm *CMapString) Get(key string) (val interface{}) {
	return cm.shard(key).Get(key)
}

func (cm *CMapString) Set(key string, val interface{}) {
	cm.shard(key).Set(key, val)
}

func (cm *CMapString) Has(key string) bool                 { return cm.shard(key).Has(key) }
func (cm *CMapString) Delete(key string)                   { cm.shard(key).Delete(key) }
func (cm *CMapString) DeleteAndGet(key string) interface{} { return cm.shard(key).DeleteAndGet(key) }

func (cm *CMapString) Update(key string, fn func(oldVal interface{}) (newVal interface{})) {
	cm.shard(key).Update(key, fn)
}

func (cm *CMapString) Swap(key string, val interface{}) interface{} {
	return cm.shard(key).Swap(key, val)
}

func (cm *CMapString) Foreach(fn func(key string, val interface{}) error) error {
	for i := range cm.shards {
		if err := cm.shards[i].ForEach(fn); err != nil {
			return err
		}
	}
	return nil
}

func (cm *CMapString) ForEachParallel(fn func(key, val interface{}) error) error {
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

func (cm *CMapString) Len() int {
	ln := 0
	for i := range cm.shards {
		ln += cm.shards[i].Len()
	}
	return ln
}
