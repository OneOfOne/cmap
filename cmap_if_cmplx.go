// +build genx
// +build genx_kt_complex64 genx_kt_complex128

package cmap

import "github.com/OneOfOne/cmap/hashers"

func hasher(key KT) uint32 { return hashers.Mix64to32(uint64(real(key) + imag(key))) }
