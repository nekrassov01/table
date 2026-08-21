// Package repeat appends repeated bytes to caller-owned buffers.
package repeat

import "bytes"

const blockSize = 256

var (
	spaces = bytes.Repeat([]byte{' '}, blockSize)
	dashes = bytes.Repeat([]byte{'-'}, blockSize)
)

// AppendSpaces appends count spaces to dst.
func AppendSpaces(dst []byte, count int) []byte {
	return appendBlock(dst, spaces, count)
}

// AppendDashes appends count dashes to dst.
func AppendDashes(dst []byte, count int) []byte {
	return appendBlock(dst, dashes, count)
}

// appendBlock appends count bytes from a homogeneous reusable block.
func appendBlock(dst, block []byte, count int) []byte {
	if count <= 0 {
		return dst
	}
	for count > len(block) {
		dst = append(dst, block...)
		count -= len(block)
	}
	return append(dst, block[:count]...)
}
