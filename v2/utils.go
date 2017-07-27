package cmap

import (
	"hash/fnv"
	"reflect"
	"unsafe"
)

func ClosestPowerOfTwo(n uint64, gt bool) (v uint64) {
	for i, p := uint64(1<<1), uint64(2); i <= n; i, p = i*2, p+1 {
		v = i
	}
	if gt && v < n {
		v *= 2
	}
	return
}

func fnvHashString(s string) uint64 {
	var (
		ss  = (*reflect.StringHeader)(unsafe.Pointer(&s))
		fnv = fnv.New64()
	)
	fnv.Write((*[0x7fffffff]byte)(unsafe.Pointer(ss.Data))[:len(s):len(s)])
	return fnv.Sum64()

}
