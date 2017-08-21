//go:generate env IS_MAIN_PACKAGE=1 /bin/sh ./scripts/mapSpec.sh cmap "interface{}" "interface{}" "" ./
//go:generate /bin/sh ./scripts/mapSpec.sh stringcmap "string" "interface{}" "hashers.Fnv32" ./stringcmap

package cmap

import (
	"errors"

	"github.com/OneOfOne/cmap/hashers"
)

// DefaultShardCount is the default number of shards to use when New() or NewFromJSON() are called.
// The default is 256
var DefaultShardCount = 1 << 8

// Break is returned to break early from ForEach without returning an error.
var Break = errors.New("break")

func DefaultKeyHasher(key interface{}) uint32 { return hashers.TypeHasher32(key) }
