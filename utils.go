package cmap

import (
	"fmt"
	"math"
)

// DefaultKeyHasher returns a hash for the specific key for internal sharding.
// By default, those types are supported as keys: KeyHasher, string, uint64, int64, uint32, int32, uint16, int16, uint8, int8, uint, int,
//  float64, float32 and fmt.Stringer.
func DefaultKeyHasher(key KT) uint32 {
	switch key := key.(type) {
	case KeyHasher:
		return rehash32(key.Hash())
	case string:
		return Fnv32(key)
	case int:
		return rehash32(uint32(key))
	case uint:
		return rehash32(uint32(key))
	case uint64:
		return rehash32(uint32(key))
	case int64:
		return rehash32(uint32(key))
	case uint32:
		return rehash32(uint32(key))
	case int32:
		return rehash32(uint32(key))
	case uint16:
		return rehash32(uint32(key))
	case int16:
		return rehash32(uint32(key))
	case uint8:
		return rehash32(uint32(key))
	case int8:
		return rehash32(uint32(key))
	case float64:
		return rehash32(uint32(math.Float64bits(key)))
	case float32:
		return rehash32(uint32(math.Float32bits(key)))
	case fmt.Stringer:
		return Fnv32(key.String())
	default:
		panic(fmt.Sprintf("unsupported type: %T (%v)", key, key))
	}
}

// Fnv32 the default hash func we use for strings.
func Fnv32(key string) uint32 {
	const prime32 = uint32(16777619)
	hash := uint32(2166136261)
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}

func rehash32(h uint32) uint32 {
	// We apply this secondary hashing discovered by Doug Lea to defend
	// against bad hashes.
	h += h ^ (h << 9)
	h ^= h >> 14
	h += h << 4
	h ^= h >> 10
	return h
}
