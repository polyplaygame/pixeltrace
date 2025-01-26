package tools

import (
	"unsafe"
)

// SliceByteToString converts []byte to string without copy.
// DO NOT USE unless you know what you're doing.
func SliceByteToString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// StringToSliceByte converts string to []byte without copy.
// DO NOT USE unless you know what you're doing.
func StringToSliceByte(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
