package cmap

import (
	"encoding/json"
	"strconv"
	"sync"
	"testing"
)

func benchCmapSetGet(b *testing.B, sz int) {
	b.Helper()
	cm := NewSize(sz)
	keys := make([]interface{}, 10000)
	for i := range keys {
		keys[i] = uint64(i)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var i int
		for pb.Next() {
			x := keys[i%len(keys)]
			cm.Set(x, x)
			if v, ok := cm.Get(x).(uint64); !ok || v != x {
				b.Fatalf("wanted %v, got %v", x, v)
			}
			i++
		}
	})

}

func BenchmarkCMap(b *testing.B) {
	for _, sz := range []int{128, 256, 512, 1024, 2048, 4096} {
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

func (mm *mutexMap) MarshalJSON() ([]byte, error) {
	mm.RLock()
	j, err := json.Marshal(mm.m)
	mm.RUnlock()
	return j, err
}

func (mm *mutexMap) UnmarshalJSON(j []byte) error {
	mm.Lock()
	err := json.Unmarshal(j, &mm.m)
	mm.Unlock()
	return err
}

func BenchmarkMutexMap(b *testing.B) {
	cm := mutexMap{m: make(map[interface{}]interface{}, DefaultShardCount)}
	keys := make([]interface{}, 10000)
	for i := range keys {
		keys[i] = uint64(i)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var i int
		for pb.Next() {
			x := keys[i%len(keys)]
			cm.Set(x, x)
			if v, ok := cm.Get(x).(uint64); !ok || v != x {
				b.Fatalf("wanted %v, got %v", x, v)
			}
			i++
		}
	})
	if testing.Verbose() {
		b.Logf("size: %v", len(cm.m))
	}
}
func BenchmarkSyncMap(b *testing.B) {
	var cm sync.Map
	keys := make([]interface{}, 10000)
	for i := range keys {
		keys[i] = uint64(i)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var i int
		for pb.Next() {
			x := keys[i%len(keys)]
			cm.Store(x, x)

			if v, ok := cm.Load(x); !ok || v.(uint64) != x {
				b.Fatalf("wanted %v, got %v", x, v)
			}
			i++
		}
	})
}
