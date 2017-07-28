// +build streamrail

package cmap_test

import (
	"fmt"
	"strconv"
	"testing"

	CC "github.com/streamrail/concurrent-map"
)

func benchSR(b *testing.B, sz int) {
	keys := make([]string, 10000)
	for i := range keys {
		keys[i] = fmt.Sprintf("%010d", i)
	}
	CC.SHARD_COUNT = sz
	cm := CC.New()
	b.RunParallel(func(pb *testing.PB) {
		var i int
		for pb.Next() {
			key := keys[i%len(keys)]
			cm.Set(key, key)
			if v, ok := cm.Get(key); !ok || v.(string) != key {
				b.Fatalf("sz: %d, wanted %v, got %v", sz, key, v)
			}
			i++
		}
	})
}

func BenchmarkStreamrail(b *testing.B) {
	shardCounts := []int{2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192}
	if testing.Short() {
		shardCounts = shardCounts[len(shardCounts)-3:]
	}
	for _, sz := range shardCounts {
		b.Run(strconv.Itoa(sz), func(b *testing.B) {
			benchSR(b, sz)
		})
	}
}
