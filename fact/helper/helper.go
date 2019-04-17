package helper

import "unsafe"

// BytePtrToString converts byte pointer to a Go string.
// borrowed from github.com/shirou/gopsutil/internal/common
func BytePtrToString(p *uint8) string {
	a := (*[10000]uint8)(unsafe.Pointer(p))
	i := 0
	for a[i] != 0 {
		i++
	}
	return string(a[:i])
}
