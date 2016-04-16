//go:generate go run "$GOPATH/src/github.com/OneOfOne/lfchan/gen.go" "*KeyValue" .

package cmap

import (
	"bytes"
	"encoding/json"
	"io"
	"sync"
	"sync/atomic"
)

var (
	// IgnoreValue can be returned from the func called to NewFromJSON to ignore setting the value
	IgnoreValue = &struct{ bool }{true}
	// DeleteValue can be returns from Update to delete the value
	DeleteValue = &struct{ bool }{true}
)

// DefaultShardCount is the default number of shards to use when New() or NewFromJSON() are called.
const DefaultShardCount = 1 << 4 // 16

// ForeachFunc is a function that gets passed to Foreach, returns true to break early
type ForEachFunc func(key string, val interface{}) (BreakEarly bool)

type CompareAndSwapFunc func(a, b interface{}) bool
type UpdateFunc func(oldVal interface{}) (newVal interface{})

// CMap is a sharded thread-safe concurrent map.
type CMap struct {
	shards []lockedMap
	l      uint64
}

// New is an alias for NewSize(DefaultShardCount)
func New() CMap { return NewSize(DefaultShardCount) }

// NewSize returns a CMap with the specific shardSize, note that for performance reasons,
// shardCount must be a power of 2
func NewSize(shardCount int) CMap {
	// must be a power of 2
	if shardCount == 0 {
		shardCount = DefaultShardCount
	} else if shardCount&(shardCount-1) != 0 {
		panic("shardCount must be a power of 2")
	}
	cm := CMap{
		shards: make([]lockedMap, shardCount),
		l:      uint64(shardCount) - 1,
	}
	for i := range cm.shards {
		cm.shards[i].m = make(map[string]interface{}, shardCount/2)
	}
	return cm
}

// NewFromJSON is an alias for NewSizeFromJSON(DefaultShardCount, r, fn)
func NewFromJSON(r io.Reader, fn func(v interface{}) interface{}) (CMap, error) {
	return NewSizeFromJSON(DefaultShardCount, r, fn)
}

// NewFromJSON returns a cmap constructed from json, fn will return the "proper" value, for example:
// json by default reads all numbers as float64, so fn(v) where v is supposed to be an int should look like:
// 	func(v interface{}) interface{} { n, _ := v.(json.Number).Int64(); return int(n) }
// note that by default all numbers will be json.Number
func NewSizeFromJSON(shardCount int, r io.Reader, fn func(v interface{}) interface{}) (CMap, error) {
	//TODO use json.RawMessage
	cm := NewSize(shardCount)
	dec := json.NewDecoder(r)
	dec.UseNumber()
	var key string
	for dec.More() {
		t, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return cm, err
		}
		if t, ok := t.(string); ok && key == "" {
			key = t
			continue
		}

		if key != "" {
			if v := fn(t); v != IgnoreValue {
				cm.shard(key).m[key] = v // no need to use locks for this
			}
			key = ""
		}
	}

	return cm, nil
}

func (cm CMap) shard(key string) *lockedMap {
	h := FNV64aString(key)
	return &cm.shards[h&cm.l]
}

func (cm CMap) Set(key string, val interface{}) { cm.shard(key).Set(key, val) }

func (cm CMap) Get(key string) interface{}          { return cm.shard(key).Get(key) }
func (cm CMap) Has(key string) bool                 { return cm.shard(key).Has(key) }
func (cm CMap) Delete(key string)                   { cm.shard(key).Delete(key) }
func (cm CMap) DeleteAndGet(key string) interface{} { return cm.shard(key).DeleteAndGet(key) }

func (cm CMap) Update(key string, fn func(ov interface{}) interface{}) {
	cm.shard(key).Update(key, fn)
}
func (cm CMap) Swap(key string, val interface{}) interface{} {
	return cm.shard(key).Swap(key, val)
}

func (cm CMap) CompareAndSwap(key string, val interface{}, eqFn CompareAndSwapFunc) bool {
	return cm.shard(key).CompareAndSwap(key, val, eqFn)
}

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
func (cm CMap) Iter() KeyValueChan { return cm.IterBuffered(0) }

// IterBuffered returns a buffered channel sz, to return an unbuffered channel you can pass 0
// ch.Break() on the returned channel can allow breaking early.
// note that calling ch.Break() will show as a race on the race detector but it's more or less a "safe" race,
// and it is the only clean way to break out of a channel.
func (cm CMap) IterBuffered(sz int) KeyValueChan {
	ch := make(KeyValueChan, sz)
	go func() {
		var wg sync.WaitGroup
		wg.Add(len(cm.shards))
		for i := range cm.shards {
			go cm.shards[i].iter(ch, &wg)
		}
		wg.Wait()
		ch.Break()
	}()
	return ch
}

func (cm CMap) Len() int {
	ln := 0
	for i := range cm.shards {
		ln += cm.shards[i].Len()
	}
	return ln
}

func (cm CMap) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('{')
	enc := json.NewEncoder(&buf)
	cm.Foreach(func(k string, v interface{}) bool {
		buf.WriteString(`"` + k + `":`)
		enc.Encode(v)
		buf.Bytes()[buf.Len()-1] = ','
		return false
	})
	if buf.Bytes()[buf.Len()-1] == ',' {
		buf.Truncate(buf.Len() - 1)
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}
