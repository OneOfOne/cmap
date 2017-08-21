//go:generate go run ./cmd/cmap-gen/main.go -v -internal -n cmap -p ./
//go:generate go run ./cmd/cmap-gen/main.go -v -n stringcmap -kt string -hfn hashers.Fnv32
//go:generate gometalinter ./ ./stringcmap

package cmap

import (
	"github.com/OneOfOne/cmap/hashers"
)

// DefaultShardCount is the default number of shards to use when New() or NewFromJSON() are called.
// The default is 256
var DefaultShardCount = 1 << 8

// DefaultKeyHasher is an alias for hashers.TypeHasher32(key)
func DefaultKeyHasher(key interface{}) uint32 { return hashers.TypeHasher32(key) }
