// +build go1.9

package cmap_test

import (
	"sync"
	"testing"
)

func BenchmarkSyncMap(b *testing.B) {
	var cm sync.Map
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var i int
		for pb.Next() {
			x := keys[i%len(keys)]
			cm.Store(x, x)

			if v, ok := cm.Load(x); !ok || v.(string) != x {
				b.Fatalf("wanted %v, got %v", x, v)
			}
			i++
		}
	})
}
