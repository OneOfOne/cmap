// +build genx
// +build genx_kt_string

package cmap

import "github.com/OneOfOne/cmap/hashers"

func hasher(key KT) uint32 { return hashers.TypeHasher32(key) }
