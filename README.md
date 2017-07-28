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
	cm := cmap.New() // or cmap.NewString()
	// cm := cmap.NewSize(1 << 8) // the size must always be a power of 2
	cm.Set("key", "value")
	ok := cm.Has("key") == true
	if v, ok := cm.Get("key").(string); ok {
		// do something with v
	}
}
```

## Benchmark
```bash
➤ go test -bench=. -benchmem -tags streamrail -benchtime 2s -cpu 1,8,32 -short
# map[interface{}]interface{}
BenchmarkCMap/2048                      20000000               172 ns/op             0 B/op          0 allocs/op
BenchmarkCMap/4096                      20000000               166 ns/op             0 B/op          0 allocs/op
BenchmarkCMap/8192                      20000000               165 ns/op             0 B/op          0 allocs/op

BenchmarkCMap/4096-8                    50000000              42.8 ns/op             0 B/op          0 allocs/op
BenchmarkCMap/2048-8                    50000000              46.0 ns/op             0 B/op          0 allocs/op
BenchmarkCMap/8192-8                    100000000             41.8 ns/op             0 B/op          0 allocs/op

BenchmarkCMap/2048-32                   100000000             43.2 ns/op             0 B/op          0 allocs/op
BenchmarkCMap/4096-32                   100000000             41.3 ns/op             0 B/op          0 allocs/op
BenchmarkCMap/8192-32                   100000000             39.1 ns/op             0 B/op          0 allocs/op

# map[string]interface{}
BenchmarkCMapString/2048                20000000               185 ns/op            16 B/op          1 allocs/op
BenchmarkCMapString/4096                20000000               178 ns/op            16 B/op          1 allocs/op
BenchmarkCMapString/8192                20000000               171 ns/op            16 B/op          1 allocs/op

BenchmarkCMapString/2048-8              50000000              47.6 ns/op            16 B/op          1 allocs/op
BenchmarkCMapString/4096-8              100000000             44.7 ns/op            16 B/op          1 allocs/op
BenchmarkCMapString/8192-8              100000000             42.9 ns/op            16 B/op          1 allocs/op

BenchmarkCMapString/2048-32             50000000              51.1 ns/op            16 B/op          1 allocs/op
BenchmarkCMapString/4096-32             50000000              48.7 ns/op            16 B/op          1 allocs/op
BenchmarkCMapString/8192-32             100000000             46.0 ns/op            16 B/op          1 allocs/op

# map[interface{}]interface{} protected with sync.RWMutex
BenchmarkMutexMap                       20000000               136 ns/op             0 B/op          0 allocs/op
BenchmarkMutexMap-8                     20000000               177 ns/op             0 B/op          0 allocs/op
BenchmarkMutexMap-32                    20000000               178 ns/o              0 B/op          0 allocs/op

# sync.Map
BenchmarkSyncMap                        20000000               157 ns/op            16 B/op          1 allocs/op
BenchmarkSyncMap-8                      100000000             38.2 ns/op            16 B/op          1 allocs/op
BenchmarkSyncMap-32                     100000000             41.0 ns/op            16 B/op          1 allocs/op

BenchmarkStreamrail/2048                20000000               201 ns/op            16 B/op          1 allocs/op
BenchmarkStreamrail/4096                20000000               196 ns/op            16 B/op          1 allocs/op
BenchmarkStreamrail/8192                20000000               188 ns/op            16 B/op          1 allocs/op

BenchmarkStreamrail/2048-8              50000000              50.0 ns/op            16 B/op          1 allocs/op
BenchmarkStreamrail/4096-8              50000000              49.3 ns/op            16 B/op          1 allocs/op
BenchmarkStreamrail/8192-8              50000000              45.7 ns/op            16 B/op          1 allocs/op

BenchmarkStreamrail/2048-32             50000000              53.2 ns/op            16 B/op          1 allocs/op
BenchmarkStreamrail/4096-32             50000000              53.4 ns/op            16 B/op          1 allocs/op
BenchmarkStreamrail/8192-32             50000000              47.9 ns/op            16 B/op          1 allocs/op
PASS
ok      github.com/OneOfOne/cmap        118.483s
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
