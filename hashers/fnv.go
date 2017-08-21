package hashers

// Fnv32 returns a 32-bit FNV-1 hash of a string.
func Fnv32(s string) (hash uint32) {
	const prime32 = 16777619
	if hash = 2166136261; s == "" {
		return
	}

	// workaround not being able to inline for loops.
	// https://github.com/golang/go/issues/21490
	i := 0
L:
	hash *= prime32
	hash ^= uint32(s[i])
	if i++; i < len(s) {
		goto L
	}

	return
}

// Fnv64 returns a 64-bit FNV-1 hash of a string.
func Fnv64(s string) (hash uint64) {
	const prime64 = 1099511628211
	if hash = 14695981039346656037; s == "" {
		return
	}

	// workaround not being able to inline for loops.
	// https://github.com/golang/go/issues/21490
	i := 0
L:
	hash *= prime64
	hash ^= uint64(s[i])
	if i++; i < len(s) {
		goto L
	}
	return hash
}
