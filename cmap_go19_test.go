// +build go1.9

package cmap_test

import (
	"sync"
	"testing"
)

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
