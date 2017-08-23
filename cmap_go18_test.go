package cmap_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/OneOfOne/cmap"
)

var keys [1e5]interface{}

func init() {
	for i := range keys {
		keys[i] = fmt.Sprintf("%010d", i)
	}
}

func TestDistruption(t *testing.T) {
	cm := cmap.NewSize(32)
	for i := 0; i < 1e6; i++ {
		cm.Set(i, i)
	}

	t.Logf("%+v", cm.ShardDistribution())

}

func TestIter(t *testing.T) {
	cm := cmap.New()
	for i := 0; i < 100; i++ {
		cm.Set(i, i)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	i := 0
	for kv := range cm.IterLocked(ctx, 1) {
		t.Logf("%d: %+v", i, kv)
		if i++; i > 10 {
			cancel()
		}
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

func BenchmarkLMap(b *testing.B) {
	cm := cmap.NewLMapSize(cmap.DefaultShardCount)
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
		b.Logf("size: %v", cm.Len())
	}
}
