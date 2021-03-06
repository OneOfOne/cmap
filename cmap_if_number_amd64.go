// +build genx
// +build genx_kt_int genx_kt_uint genx_kt_int32 genx_kt_uint32 genx_kt_int64 genx_kt_uint64 genx_kt_float64 genx_kt_float32
// +build amd64

package cmap

import (
	"github.com/OneOfOne/cmap/hashers"
)

// hasher uses Mix64to32 because it's faster on 64bit
func hasher(key KT) uint32 {
	return hashers.Mix64to32(uint64(key)) // nolint:unconvert
}
