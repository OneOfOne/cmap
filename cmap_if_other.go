// +build genx
// +build !genx_kt_int
// +build !genx_kt_uint
// +build !genx_kt_int32
// +build !genx_kt_uint32
// +build !genx_kt_int64
// +build !genx_kt_uint64
// +build !genx_kt_float64
// +build !genx_kt_float32
// +build !genx_kt_string

package cmap

import "github.com/OneOfOne/cmap/hashers"

func hasher(key KT) uint32 { return hashers.TypeHasher32(key) }
