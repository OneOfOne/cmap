//go:generate go install github.com/OneOfOne/genx/...
//go:generate genx -pkg github.com/OneOfOne/cmap -v -t KT=interface{},VT=interface{} -o ./cmap_iface_iface.go
//go:generate genx -pkg github.com/OneOfOne/cmap -v -n stringcmap -t KT=string,VT=interface{} -o ./stringcmap/cmap_string_iface.go
//go:generate genx -pkg github.com/OneOfOne/cmap -v -n u64cmap -t KT=uint64,VT=interface{} -o ./u64cmap/cmap_u64_iface.go
//go:generate gometalinter --aggregate --cyclo-over=17 ./...

package cmap
