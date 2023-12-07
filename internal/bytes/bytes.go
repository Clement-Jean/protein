package bytes

import "unsafe"

func ToString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

func FromString(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
