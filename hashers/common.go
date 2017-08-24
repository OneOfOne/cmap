package hashers

import (
	"fmt"
	"math"
	"reflect"
)

// KeyHasher is a type that provides its own hash function.
type KeyHasher interface {
	Hash() uint64
}

// TypeHasher32 returns a hash for the specific key for internal sharding.
// By default, those types are supported as keys: KeyHasher, string, uint64, int64, uint32, int32, uint16, int16, uint8,
// int8, uint, int, float64, float32 and fmt.Stringer.
// Falls back to Fnv32(reflect.ValueOf(v).String()).
func TypeHasher32(v interface{}) uint32 {
	switch v := v.(type) {
	case KeyHasher:
		return Mix32(uint32(v.Hash()))
	case string:
		return Fnv32(v)
	case int:
		return Mix32(uint32(v))
	case uint:
		return Mix32(uint32(v))
	case uint64:
		return Mix32(uint32(v))
	case int64:
		return Mix32(uint32(v))
	case uint32:
		return Mix32(v)
	case int32:
		return Mix32(uint32(v))
	case uint16:
		return Mix32(uint32(v))
	case int16:
		return Mix32(uint32(v))
	case uint8:
		return Mix32(uint32(v))
	case int8:
		return Mix32(uint32(v))
	case float64:
		return Mix32(uint32(math.Float64bits(v)))
	case float32:
		return Mix32(math.Float32bits(v))
	case fmt.Stringer:
		return Fnv32(v.String())
	default:
		return Fnv32(reflect.ValueOf(v).String())
	}
}

// TypeHasher64 returns a hash for the specific key for internal sharding.
// By default, those types are supported as keys: KeyHasher, string, uint64, int64, uint32, int32, uint16, int16, uint8,
// int8, uint, int, float64, float32 and fmt.Stringer.
// Falls back to Fnv64(reflect.ValueOf(v).String()).
func TypeHasher64(v interface{}) uint64 {
	switch v := v.(type) {
	case KeyHasher:
		return Mix64(v.Hash())
	case string:
		return Fnv64(v)
	case int:
		return Mix64(uint64(v))
	case uint:
		return Mix64(uint64(v))
	case uint64:
		return Mix64(v)
	case int64:
		return Mix64(uint64(v))
	case uint32:
		return Mix64(uint64(v))
	case int32:
		return Mix64(uint64(v))
	case uint16:
		return Mix64(uint64(v))
	case int16:
		return Mix64(uint64(v))
	case uint8:
		return Mix64(uint64(v))
	case int8:
		return Mix64(uint64(v))
	case float64:
		return Mix64(math.Float64bits(v))
	case float32:
		return Mix64(uint64(math.Float32bits(v)))
	case fmt.Stringer:
		return Fnv64(v.String())
	default:
		return Fnv64(reflect.ValueOf(v).String())
	}
}

// Mix32 mixes the hash to make sure the bits are spread, borrowed from xxhash.
func Mix32(h uint32) uint32 {
	const prime32x2 = 2246822519
	const prime32x3 = 3266489917
	h ^= h >> 15
	h *= prime32x2
	h ^= h >> 13
	h *= prime32x3
	h ^= h >> 16
	return h
}

// Mix64to32 is a helper to mix 64bit number down to 32.
func Mix64to32(h uint64) uint32 {
	return uint32(Mix64(h))
}

// Mix64 mixes the hash to make sure the bits are spread, borrowed from xxhash.
func Mix64(h uint64) uint64 {
	const prime64x2 = 14029467366897019727
	const prime64x3 = 1609587929392839161

	h ^= h >> 33
	h *= prime64x2
	h ^= h >> 29
	h *= prime64x3
	h ^= h >> 32
	return h
}
