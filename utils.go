package cmap

import "github.com/OneOfOne/xxhash/native"

func cmapHashString(s string) uint64 {
	return xxhash.ChecksumString64(s)
}
