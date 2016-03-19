package cmap

func cmapHashString(s string) uint64 {
	//return xxhash.ChecksumString64(s)
	return JSWStringHash(s)
}

// based on the implementation at http://eternallyconfuzzled.com/tuts/algorithms/jsw_tut_hashing.aspx
func JSWStringHash(s string) uint64 {
	h := uint64(16777551)
	for i := range s {
		h = (h<<1 | h>>31) ^ uint64(s[i])
	}
	return h
}

// based on the implementation at http://eternallyconfuzzled.com/tuts/algorithms/jsw_tut_hashing.aspx
func JSWHash(p []byte) uint64 {
	h := uint64(16777551)
	for i := range p {
		h = (h<<1 | h>>31) ^ uint64(p[i])
	}
	return h
}
