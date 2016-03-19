// +build streamrail

package cmap

import (
	"strconv"
	"sync/atomic"
	"testing"

	"github.com/streamrail/concurrent-map"
)

func benchSRCmapSetGet(b *testing.B, sz int) {
	cmap.SHARD_COUNT = sz
	cm := cmap.New()
	var i uint64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			x := strconv.FormatUint(atomic.AddUint64(&i, 1), 10)
			cm.Set(x, x)
			if v, ok := cm.Get(x); !ok || v.(string) != x {
				b.Fatalf("wanted %v, got %v", x, v)
			}
		}
	})
	if testing.Verbose() {
		b.Logf("size: %v", cm.Count())
	}
}

func BenchmarkSRCMap8Shards(b *testing.B)   { benchSRCmapSetGet(b, 8) }
func BenchmarkSRCMap16Shards(b *testing.B)  { benchSRCmapSetGet(b, 16) }
func BenchmarkSRCMap32Shards(b *testing.B)  { benchSRCmapSetGet(b, 32) }
func BenchmarkSRCMap64Shards(b *testing.B)  { benchSRCmapSetGet(b, 64) }
func BenchmarkSRCMap128Shards(b *testing.B) { benchSRCmapSetGet(b, 128) }
func BenchmarkSRCMap256Shards(b *testing.B) { benchSRCmapSetGet(b, 256) }
func BenchmarkSRCMap512Shards(b *testing.B) { benchSRCmapSetGet(b, 512) }
