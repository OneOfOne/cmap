package stringcmap_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/OneOfOne/cmap/stringcmap"
)

var keys [1e5]string

func init() {
	for i := range keys {
		keys[i] = fmt.Sprintf("%010d", i)
	}
}

func benchCmapStringSetGet(b *testing.B, sz int) {
	cm := stringcmap.NewSize(sz)
	keys := make([]string, 10000)
	for i := range keys {
		keys[i] = fmt.Sprintf("%010d", i)
	}
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

func BenchmarkCMapString(b *testing.B) {
	shardCounts := []int{2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192}
	if testing.Short() {
		shardCounts = shardCounts[len(shardCounts)-3:]
	}
	for _, sz := range shardCounts {
		b.Run(strconv.Itoa(sz), func(b *testing.B) {
			benchCmapStringSetGet(b, sz)
		})
	}
}
