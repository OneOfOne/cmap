package cmap_test

import (
	"fmt"
	"strconv"
	"sync"
	"testing"

	"github.com/OneOfOne/cmap"
)

var keys [1e5]string

func init() {
	for i := range keys {
		keys[i] = fmt.Sprintf("%010d", i)
	}
}

func benchCmapSetGet(b *testing.B, sz int) {
	cm := cmap.NewSize(sz)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var i int
		for pb.Next() {
			x := keys[i%len(keys)]
			cm.Set(x, x)
			if v, ok := cm.Get(x).(string); !ok || v != x {
				b.Fatalf("sz: %d, wanted %v, got %v", sz, x, v)
			}
			i++
		}
	})

}

func BenchmarkCMap(b *testing.B) {
	shardCounts := []int{32, 64, 128, 256, 512, 1024, 2048, 4096, 8192}
	if testing.Short() {
		shardCounts = shardCounts[len(shardCounts)-3:]
	}
	for _, sz := range shardCounts {
		b.Run(strconv.Itoa(sz), func(b *testing.B) {
			benchCmapSetGet(b, sz)
		})
	}
}

type mutexMap struct {
	sync.RWMutex
	m map[interface{}]interface{}
}

func (mm *mutexMap) Set(k interface{}, v interface{}) {
	mm.Lock()
	mm.m[k] = v
	mm.Unlock()
}

func (mm *mutexMap) Get(k interface{}) interface{} {
	mm.RLock()
	v := mm.m[k]
	mm.RUnlock()
	return v
}

func BenchmarkMutexMap(b *testing.B) {
	cm := mutexMap{m: make(map[interface{}]interface{}, cmap.DefaultShardCount)}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var i int
		for pb.Next() {
			x := keys[i%len(keys)]
			cm.Set(x, x)
			if v, ok := cm.Get(x).(string); !ok || v != x {
				b.Fatalf("wanted %v, got %v", x, v)
			}
			i++
		}
	})
	if testing.Verbose() {
		b.Logf("size: %v", len(cm.m))
	}
}
