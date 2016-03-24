package cmap

const (
	offset32 uint32 = 2166136261
	offset64 uint64 = 14695981039346656037
	prime32  uint32 = 16777619
	prime64  uint64 = 1099511628211
)

// FNV64aString is a fnv64a hash without copying the value to a []byte
func FNV64aString(s string) uint64 {
	h := offset64
	for i := 0; i < len(s); i++ {
		h ^= uint64(s[i])
		h *= prime64
	}
	return h
}

// FNV32aString is a fnv32a hash without copying the value to a []byte
func FNV32aString(s string) uint32 {
	h := offset32
	for i := 0; i < len(s); i++ {
		h ^= uint32(s[i])
		h *= prime32
	}
	return h
}
