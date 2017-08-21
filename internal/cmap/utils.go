package cmap

import "github.com/OneOfOne/cmap/hashers"

type (
	KT interface{}
	VT interface{}
)

func DefaultKeyHasher(key KT) uint32 { return hashers.TypeHasher32(key) }
