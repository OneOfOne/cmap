//go:generate go install ./cmd/cmap-gen
//go:generate cmap-gen -v -internal -n cmap -p ./
//go:generate cmap-gen -v -n stringcmap -kt string -hfn hashers.Fnv32
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
