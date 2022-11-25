package xstrings

import (
	"unsafe"
)

// SameData reports whether s1 and s2 have identical underlying data pointers.
//
// This could be used in a fast string comparison path when you don't
// need the precise string equality answer.
//
// For example, when SameData(s1, s2) evaluates to true, you'll never need
// to check the data bytes; the strings are known to be identical.
func SameData(s1, s2 string) bool {
	return stringptr(s1) == stringptr(s2)
}

type stringHeader struct {
	data *byte
	len  int
}

func stringptr(s string) *byte {
	return (*stringHeader)(unsafe.Pointer(&s)).data
}
