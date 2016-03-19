package cmap

import (
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
)

func benchCmapSetGet(b *testing.B, sz int) {
	cm := NewSize(sz)
	var i uint64
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			x := strconv.FormatUint(atomic.AddUint64(&i, 1), 10)
			cm.Set(x, x)
			if v, ok := cm.Get(x).(string); !ok || v != x {
				b.Fatalf("wanted %v, got %v", x, v)
			}
		}
	})
	if testing.Verbose() {
		shardCounts := make([]int, sz)
		for i := range cm.shards {
			shardCounts[i] = cm.shards[i].Len()
		}
		b.Logf("size: %v: %v", cm.Len(), shardCounts)
	}
}

func BenchmarkCMap8Shards(b *testing.B)   { benchCmapSetGet(b, 8) }
func BenchmarkCMap16Shards(b *testing.B)  { benchCmapSetGet(b, 16) }
func BenchmarkCMap32Shards(b *testing.B)  { benchCmapSetGet(b, 32) }
func BenchmarkCMap64Shards(b *testing.B)  { benchCmapSetGet(b, 64) }
func BenchmarkCMap128Shards(b *testing.B) { benchCmapSetGet(b, 128) }
func BenchmarkCMap256Shards(b *testing.B) { benchCmapSetGet(b, 256) }

// func BenchmarkCMap512Shards(b *testing.B) { benchCmapSetGet(b, 512) }

type mutexMap struct {
	sync.RWMutex
	m map[string]interface{}
}

func (mm *mutexMap) Set(k string, v interface{}) {
	mm.Lock()
	mm.m[k] = v
	mm.Unlock()
}

func (mm *mutexMap) Get(k string) interface{} {
	mm.RLock()
	v := mm.m[k]
	mm.RUnlock()
	return v
}

func BenchmarkMutexMap(b *testing.B) {
	cm := mutexMap{m: make(map[string]interface{}, DefaultShardCount)}
	var i uint64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			x := strconv.FormatUint(atomic.AddUint64(&i, 1), 10)
			cm.Set(x, x)
			if v, ok := cm.Get(x).(string); !ok || v != x {
				b.Fatalf("wanted %v, got %v", x, v)
			}
		}
	})
	if testing.Verbose() {
		b.Logf("size: %v", len(cm.m))
	}
}
