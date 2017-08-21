# cmap [![GoDoc](https://godoc.org/github.com/OneOfOne/cmap?status.svg)](https://godoc.org/github.com/OneOfOne/cmap) [![Build Status](https://travis-ci.org/OneOfOne/cmap.svg?branch=master)](https://travis-ci.org/OneOfOne/cmap) [![Coverage](https://gocover.io/_badge/github.com/OneOfOne/cmap)](https://gocover.io/github.com/OneOfOne/cmap)
--

CMap (concurrent-map) is a sharded map implementation to support fast concurrent access.

## Install

	go get github.com/OneOfOne/cmap

## Features

* Full concurrent access (except for Update).
* Supports `Get`, `Set`, `SetIfNotExists`, `Swap`, `Update`, `Delete`, `DeleteAndGet` (Pop).
* `ForEach` / `Iter` supports modifing the map during the iteration like `map` and `sync.Map`.
* `stringcmap.CMap` gives a specialized version to support map[string]interface{}.
* `stringcmap.MapWithJSON` implements json.Unmarshaler with a custom value unmarshaler.

## FAQ

### Why?
* A simple sync.RWMutex wrapped map is much slower as the concurrency increase.
* Provides several helper functions, Swap(), Update, DeleteAndGet.

### Why not `sync.Map`?
* `sync.Map` is great, I absolute love it if all you need is pure Load/Store, however you can't safely update values in it.

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
	cm.Update("key", func(old interface{}) interface{} {
		v, _ := old.(uint64)
		return v + 1
	})
}
```

## Benchmark
```bash
➤ go version; go test -tags streamrail -short -bench=. -benchmem -count 5 ./ ./stringcmap/ | benchstat /dev/stdin
go version devel +ff90f4af66 2017-08-19 12:56:24 +0000 linux/amd64

name               time/op
# pkg:github.com/OneOfOne/cmap goos:linux goarch:amd64
CMap/2048-8        85.3ns ± 2%
CMap/4096-8        86.5ns ± 1%
CMap/8192-8        95.0ns ±16%

# simple map[interface{}]interface{} wrapped with a sync.RWMutex
MutexMap-8          486ns ± 9%

# sync.Map
SyncMap-8           511ns ±28%

# pkg:github.com/OneOfOne/cmap/stringcmap goos:linux goarch:amd64
StringCMap/2048-8  38.3ns ± 3%
StringCMap/4096-8  37.9ns ± 5%
StringCMap/8192-8  38.5ns ±17%

Streamrail/2048-8  47.2ns ± 1%
Streamrail/4096-8  46.6ns ± 1%
Streamrail/8192-8  46.7ns ± 2%

name               alloc/op
# pkg:github.com/OneOfOne/cmap goos:linux goarch:amd64
CMap/2048-8         48.0B ± 0%
CMap/4096-8         48.0B ± 0%
CMap/8192-8         48.0B ± 0%

MutexMap-8          35.0B ± 0%

SyncMap-8           63.4B ± 7%

# pkg:github.com/OneOfOne/cmap/stringcmap goos:linux goarch:amd64

# specialized version of CMap, using map[string]interface{} internally
StringCMap/2048-8   16.0B ± 0%
StringCMap/4096-8   16.0B ± 0%
StringCMap/8192-8   16.0B ± 0%

# github.com/streamrail/concurrent-map
Streamrail/2048-8   16.0B ± 0%
Streamrail/4096-8   16.0B ± 0%
Streamrail/8192-8   16.0B ± 0%

name               allocs/op
# pkg:github.com/OneOfOne/cmap goos:linux goarch:amd64
CMap/2048-8          3.00 ± 0%
CMap/4096-8          3.00 ± 0%
CMap/8192-8          3.00 ± 0%

MutexMap-8           2.00 ± 0%

SyncMap-8            3.00 ± 0%

# pkg:github.com/OneOfOne/cmap/stringcmap goos:linux goarch:amd64
StringCMap/2048-8    1.00 ± 0%
StringCMap/4096-8    1.00 ± 0%
StringCMap/8192-8    1.00 ± 0%

Streamrail/2048-8    1.00 ± 0%
Streamrail/4096-8    1.00 ± 0%
Streamrail/8192-8    1.00 ± 0%
```

## License

Apache v2.0 (see [LICENSE](https://github.com/OneOfOne/cmap/blob/master/LICENSE) file).

Copyright 2016-2017 Ahmed <[OneOfOne](https://github.com/OneOfOne/)> W.

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
