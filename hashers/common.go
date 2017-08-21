package hashers

import (
	"fmt"
	"math"
)

// KeyHasher is a type that provides its own hash function.
type KeyHasher interface {
	Hash() uint64
}

// TypeHasher32 returns a hash for the specific key for internal sharding.
// By default, those types are supported as keys: KeyHasher, string, uint64, int64, uint32, int32, uint16, int16, uint8,
// int8, uint, int, float64, float32 and fmt.Stringer.
func TypeHasher32(v interface{}) uint32 {
	switch v := v.(type) {
	case KeyHasher:
		return MixHash32(uint32(v.Hash()))
	case string:
		return Fnv32(v)
	case int:
		return MixHash32(uint32(v))
	case uint:
		return MixHash32(uint32(v))
	case uint64:
		return MixHash32(uint32(v))
	case int64:
		return MixHash32(uint32(v))
	case uint32:
		return MixHash32(v)
	case int32:
		return MixHash32(uint32(v))
	case uint16:
		return MixHash32(uint32(v))
	case int16:
		return MixHash32(uint32(v))
	case uint8:
		return MixHash32(uint32(v))
	case int8:
		return MixHash32(uint32(v))
	case float64:
		return MixHash32(uint32(math.Float64bits(v)))
	case float32:
		return MixHash32(math.Float32bits(v))
	case fmt.Stringer:
		return Fnv32(v.String())
	default:
		panic(fmt.Sprintf("unsupported type: %T (%v)", v, v))
	}
}

// TypeHasher64 returns a hash for the specific key for internal sharding.
// By default, those types are supported as keys: KeyHasher, string, uint64, int64, uint32, int32, uint16, int16, uint8,
// int8, uint, int, float64, float32 and fmt.Stringer.
func TypeHasher64(v interface{}) uint64 {
	switch v := v.(type) {
	case KeyHasher:
		return MixHash64(v.Hash())
	case string:
		return Fnv64(v)
	case int:
		return MixHash64(uint64(v))
	case uint:
		return MixHash64(uint64(v))
	case uint64:
		return MixHash64(v)
	case int64:
		return MixHash64(uint64(v))
	case uint32:
		return MixHash64(uint64(v))
	case int32:
		return MixHash64(uint64(v))
	case uint16:
		return MixHash64(uint64(v))
	case int16:
		return MixHash64(uint64(v))
	case uint8:
		return MixHash64(uint64(v))
	case int8:
		return MixHash64(uint64(v))
	case float64:
		return MixHash64(math.Float64bits(v))
	case float32:
		return MixHash64(uint64(math.Float32bits(v)))
	case fmt.Stringer:
		return Fnv64(v.String())
	default:
		panic(fmt.Sprintf("unsupported type: %T (%v)", v, v))
	}
}

// MixHash32 mixes the hash to make sure the bits are spread, borrowed from xxhash.
func MixHash32(h uint32) uint32 {
	const prime32x2 = 2246822519
	const prime32x3 = 3266489917
	h ^= h >> 15
	h *= prime32x2
	h ^= h >> 13
	h *= prime32x3
	h ^= h >> 16
	return h
}

// MixHash64 mixes the hash to make sure the bits are spread, borrowed from xxhash.
func MixHash64(h uint64) uint64 {
	const prime64x2 = 14029467366897019727
	const prime64x3 = 1609587929392839161

	h ^= h >> 33
	h *= prime64x2
	h ^= h >> 29
	h *= prime64x3
	h ^= h >> 32
	return h
}
