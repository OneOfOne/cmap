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
	it := cmap.Iter()
	for v := it.Recv(); v != nil; v = it.Recv() {
		if kv.Key == "key" {
			it.Break()
		}
	}
}
```

## Benchmark
```bash
➜ go test -bench='256|Mutex' -benchmem -tags streamrail -benchtime 2s -cpu 1,8,32
# https://github.com/streamrail/concurrent-map
BenchmarkSRCMap256Shards       	 2000000	      1210 ns/op	     184 B/op	       2 allocs/op
BenchmarkSRCMap256Shards-8     	10000000	       365 ns/op	     154 B/op	       2 allocs/op
BenchmarkSRCMap256Shards-32    	10000000	       287 ns/op	     154 B/op	       2 allocs/op

BenchmarkCMap256Shards         	 3000000	       743 ns/op	     135 B/op	       2 allocs/op
BenchmarkCMap256Shards-8       	10000000	       307 ns/op	     153 B/op	       2 allocs/op
BenchmarkCMap256Shards-32      	10000000	       235 ns/op	     153 B/op	       2 allocs/op

BenchmarkMutexMap              	 2000000	      1275 ns/op	     183 B/op	       2 allocs/op
BenchmarkMutexMap-8            	 3000000	       739 ns/op	     135 B/op	       2 allocs/op
BenchmarkMutexMap-32           	 3000000	       716 ns/op	     135 B/op	       2 allocs/op
PASS
ok  	github.com/OneOfOne/cmap	30.817s

➜ go test -bench=. -benchmem -tags streamrail -benchtime 5s
testing: warning: no tests to run
# https://github.com/streamrail/concurrent-map
BenchmarkSRCMap8Shards-8        10000000               772 ns/op             153 B/op          2 allocs/op
BenchmarkSRCMap16Shards-8       10000000               710 ns/op             153 B/op          2 allocs/op
BenchmarkSRCMap32Shards-8       20000000               632 ns/op             153 B/op          2 allocs/op
BenchmarkSRCMap64Shards-8       20000000               495 ns/op             153 B/op          2 allocs/op
BenchmarkSRCMap128Shards-8      20000000               420 ns/op             154 B/op          2 allocs/op
BenchmarkSRCMap256Shards-8      20000000               388 ns/op             154 B/op          2 allocs/op

# this package
BenchmarkCMap8Shards-8          10000000               591 ns/op             153 B/op          2 allocs/op
BenchmarkCMap16Shards-8         20000000               492 ns/op             153 B/op          2 allocs/op
BenchmarkCMap32Shards-8         20000000               446 ns/op             153 B/op          2 allocs/op
BenchmarkCMap64Shards-8         20000000               362 ns/op             153 B/op          2 allocs/op
BenchmarkCMap128Shards-8        20000000               368 ns/op             153 B/op          2 allocs/op
BenchmarkCMap256Shards-8        20000000               320 ns/op             153 B/op          2 allocs/op

# simple RWMutex-guarded map
BenchmarkMutexMap-8             10000000               901 ns/op             153 B/op          2 allocs/op
PASS
ok      github.com/OneOfOne/cmap        88.383s
```

## License

Apache v2.0 (see [LICENSE](https://github.com/OneOfOne/cmap/blob/master/LICENSE) file).

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
