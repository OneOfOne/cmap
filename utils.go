package cmap

import (
	"fmt"
	"math"
)

// DefaultKeyHasher returns a hash for the specific key for internal sharding.
// By default, those types are supported as keys: string, uint64, int64, uint32, int32, uint, int,
//  float64, float32 and KeyHasher.
func DefaultKeyHasher(key interface{}) uint32 {
	switch key := key.(type) {
	case string:
		return fnv32(key)
	case uint64:
		return rehash32(uint32(key))
	case int64:
		return rehash32(uint32(key))
	case float64:
		return rehash32(uint32(math.Float64bits(key)))
	case float32:
		return rehash32(uint32(math.Float32bits(key)))
	case int:
		return rehash32(uint32(key))
	case uint:
		return rehash32(uint32(key))
	case KeyHasher:
		return rehash32(key.Hash())
	default:
		panic(fmt.Sprintf("unsupported type: %T (%v)", key, key))
	}
}

func fnv64(key string) uint64 {
	const prime64 = 1099511628211
	hash := uint64(14695981039346656037)
	for _, r := range key {
		hash *= prime64
		hash ^= uint64(r)
	}
	return hash
}

func fnv32(key string) uint32 {
	const prime32 = uint32(16777619)
	hash := uint32(2166136261)
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}

// bastardized fnv64 for numeric keys
func rehash64(key uint64) uint64 {
	const prime64 = 1099511628211
	hash := uint64(14695981039346656037)
	for ; key > 0; key = key >> 2 {
		hash *= prime64
		hash ^= uint64(key)
	}
	return hash
}

func rehash32(key uint32) uint32 {
	const prime32 = uint32(16777619)
	hash := uint32(2166136261)
	for ; key > 0; key = key >> 2 {
		hash *= prime32
		hash ^= uint32(key)
	}
	return hash
}
