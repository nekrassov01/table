package value

import "unsafe"

// packString keeps the input pointer and length together in a Value.
// The typed pointer keeps the immutable data reachable independently of Store.
func packString(s string) Value {
	// #nosec G103 -- The typed pointer retains the original string data.
	return Value{
		num: uint64(len(s)),
		any: tagString(unsafe.StringData(s)),
	}
}

// unpackString requires a Value created by packString.
func unpackString(v Value) string {
	// #nosec G103 -- The pointer and length describe the same retained input.
	return unsafe.String(v.any.(tagString), v.num)
}

// packBytes borrows the input and keeps its pointer and length together.
// The typed pointer preserves nilness and reachability independently of Store.
func packBytes(b []byte) Value {
	// #nosec G103 -- The typed pointer retains the original slice data.
	return Value{
		num: uint64(len(b)),
		any: tagBytes(unsafe.SliceData(b)),
	}
}

// unpackBytes requires a Value created by packBytes. The returned slice
// shares the input's storage and has capacity equal to its retained length.
func unpackBytes(v Value) []byte {
	// #nosec G103 -- The pointer and length describe the same retained input.
	return unsafe.Slice(v.any.(tagBytes), v.num)
}
