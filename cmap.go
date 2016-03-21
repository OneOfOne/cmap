package cmap

import (
	"errors"
	"sync"
	"sync/atomic"
)

var Break = errors.New(":break:")

const DefaultShardCount = 1 << 8 // 256

// ForeachFunc is a function that gets passed to Foreach, returns true to break early
type ForEachFunc func(key string, val interface{}) (BreakEarly bool)

type MapShard struct {
	m map[string]interface{}
	l sync.RWMutex
}

func (ms *MapShard) Set(key string, v interface{}) {
	ms.l.Lock()
	ms.m[key] = v
	ms.l.Unlock()
}

func (ms *MapShard) Get(key string) interface{} {
	ms.l.RLock()
	v := ms.m[key]
	ms.l.RUnlock()
	return v
}

func (ms *MapShard) Has(key string) bool {
	ms.l.RLock()
	_, ok := ms.m[key]
	ms.l.RUnlock()
	return ok
}

func (ms *MapShard) Delete(key string) {
	ms.l.Lock()
	delete(ms.m, key)
	ms.l.Unlock()
}

func (ms *MapShard) DeleteAndGet(key string) interface{} {
	ms.l.Lock()
	v := ms.m[key]
	delete(ms.m, key)
	ms.l.Unlock()
	return v
}

func (ms *MapShard) Len() int {
	ms.l.RLock()
	ln := len(ms.m)
	ms.l.RUnlock()
	return ln
}

func (ms *MapShard) ForEach(fn ForEachFunc) {
	ms.l.RLock()
	for k, v := range ms.m {
		if !fn(k, v) {
			break
		}
	}
	ms.l.RUnlock()
}

func (ms *MapShard) iter(ch chan *KeyValue, wg *sync.WaitGroup) {
	var kv KeyValue
	ms.l.RLock()
	defer func() { recover(); ms.l.RUnlock(); wg.Done() }()
	for k, v := range ms.m {
		kv.Key, kv.Value = k, v
		ch <- &kv
	}
}

type CMap struct {
	shards []MapShard
	l      uint64
	HashFn func(s string) uint64 // HashFn returns a hash used to select which shard to map a key to.
}

// New is an alias for NewSize(DefaultShardCount)
func New() CMap { return NewSize(DefaultShardCount) }

func NewSize(shardCount int) CMap {
	// must be a power of 2
	if shardCount == 0 {
		shardCount = DefaultShardCount
	} else if shardCount&(shardCount-1) != 0 {
		panic("shardCount must be a power of 2")
	}
	cm := CMap{
		shards: make([]MapShard, shardCount),
		l:      uint64(shardCount) - 1,
		HashFn: FNV64aString,
	}
	for i := range cm.shards {
		cm.shards[i].m = make(map[string]interface{}, shardCount/2)
	}
	return cm
}

func (cm CMap) Shard(key string) *MapShard {
	h := cm.HashFn(key)
	return &cm.shards[h&cm.l]
}

func (cm CMap) Set(key string, val interface{})     { cm.Shard(key).Set(key, val) }
func (cm CMap) Get(key string) interface{}          { return cm.Shard(key).Get(key) }
func (cm CMap) Has(key string) bool                 { return cm.Shard(key).Has(key) }
func (cm CMap) Delete(key string)                   { cm.Shard(key).Delete(key) }
func (cm CMap) DeleteAndGet(key string) interface{} { return cm.Shard(key).DeleteAndGet(key) }

func (cm CMap) Foreach(fn ForEachFunc) {
	for i := range cm.shards {
		cm.shards[i].ForEach(fn)
	}
}

func (cm CMap) ForEachParallel(fn ForEachFunc) {
	var (
		wg   sync.WaitGroup
		exit uint32
	)
	wg.Add(len(cm.shards))
	for i := range cm.shards {
		go func(i int) {
			cm.shards[i].ForEach(func(k string, v interface{}) bool {
				if atomic.LoadUint32(&exit) == 1 {
					return true
				}
				b := fn(k, v)
				if b {
					atomic.StoreUint32(&exit, 1)
				}
				return b
			})
			wg.Done()
		}(i)
	}
	wg.Wait()
}

type KeyValue struct {
	Key   string
	Value interface{}
}

// Iter is an alias for IterBuffered(0)
func (cm CMap) Iter() (<-chan *KeyValue, func()) {
	return cm.IterBuffered(0)
}

// IterBuffered returns a buffered channal shardCount * sz, to return an unbuffered channel you can pass 0
// calling breakLoop will close the channel and consume any remaining values in it.
func (cm CMap) IterBuffered(sz int) (_ <-chan *KeyValue, breakLoop func()) {
	ch := make(chan *KeyValue, len(cm.shards)*sz)
	go func() {
		defer func() { recover() }()
		var wg sync.WaitGroup
		wg.Add(len(cm.shards))
		for i := range cm.shards {
			go cm.shards[i].iter(ch, &wg)
		}
		wg.Wait()
		close(ch)
		ch = nil
	}()
	return ch, func() {
		defer func() { recover() }()
		close(ch)
		for range ch {
		}
	}
}

func (cm CMap) Len() int {
	ln := 0
	for i := range cm.shards {
		ln += cm.shards[i].Len()
	}
	return ln
}
