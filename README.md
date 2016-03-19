# cmap [![GoDoc](http://godoc.org/github.com/OneOfOne/cmap?status.svg)](http://godoc.org/github.com/OneOfOne/cmap) [![Build Status](https://travis-ci.org/OneOfOne/cmap.svg?branch=master)](https://travis-ci.org/OneOfOne/cmap)
--

CMap (concurrent-map) is a sharded map implementation to support fast concurrent access.

## Install

	go get github.com/OneOfOne/cmap

## Usage

```go
import (
	"github.com/OneOfOne/cmap"
)

func main() {
	cm := cmap.New() // or
	// cm := cmap.NewSize(1 << 8) // the size must always be a power of 2
	cm.Set("key", "value")
	ok := cm.Has("key") == true
	if v, ok := cm.Get("key").(string); ok {
		// do something with v
	}
	for kv := range cmap.Iter() {
		if kv.Key == "key" {
			kv.Break = true
			break
		}
	}
	NewSize(DefaultShardCount)
}
```

## Benchmark
```bash
➜ go test -bench=. -benchmem -tags streamrail -benchtime 3s
testing: warning: no tests to run
# this package
BenchmarkCMap8Shards-8          10000000               545 ns/op             153 B/op          2 allocs/op
BenchmarkCMap16Shards-8         10000000               462 ns/op             153 B/op          2 allocs/op
BenchmarkCMap32Shards-8         10000000               428 ns/op             153 B/op          2 allocs/op
BenchmarkCMap64Shards-8         10000000               485 ns/op             153 B/op          2 allocs/op
BenchmarkCMap128Shards-8        10000000               415 ns/op             154 B/op          2 allocs/op
BenchmarkCMap256Shards-8        10000000               384 ns/op             153 B/op          2 allocs/op
BenchmarkCMap512Shards-8        20000000               453 ns/op             154 B/op          2 allocs/op

# simple RWMutex-guarded map
BenchmarkMutexMap-8              5000000               786 ns/op             153 B/op          2 allocs/op

# https://github.com/streamrail/concurrent-map
BenchmarkSRCMap8Shards-8         5000000               809 ns/op             153 B/op          2 allocs/op
BenchmarkSRCMap16Shards-8       10000000               665 ns/op             153 B/op          2 allocs/op
BenchmarkSRCMap32Shards-8       10000000               627 ns/op             153 B/op          2 allocs/op
BenchmarkSRCMap64Shards-8       10000000               491 ns/op             154 B/op          2 allocs/op
BenchmarkSRCMap128Shards-8      10000000               404 ns/op             154 B/op          2 allocs/op
BenchmarkSRCMap256Shards-8      10000000               400 ns/op             154 B/op          2 allocs/op
BenchmarkSRCMap512Shards-8      20000000               387 ns/op             154 B/op          2 allocs/op
PASS
ok      github.com/OneOfOne/cmap        88.383s
```

## License

Apache v2.0 (see [LICENSE](https://github.com/OneOfOne/cmap/blob/master/LICENSE file).

Copyright 2016-2016 Ahmed <[OneOfOne](https://github.com/OneOfOne/)> W.

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
