package sandbox

import "unsafe"

// faster with no memory copy
func strToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func byteArrayToString(buf []byte) string {
	return *(*string)(unsafe.Pointer(&buf))
}
