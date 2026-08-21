// Package unsafe provides zero-copy conversion from byte slices to strings
// for append-only storage used within this module.
//
// A returned string shares its backing memory with the source slice. Callers
// must not modify the referenced bytes and must discard the string before the
// backing storage is reset or reused. Keeping this operation in one package
// centralizes that lifetime requirement and isolates direct use of the
// standard library's unsafe package.
package unsafe

import "unsafe"

// View returns b as a string without copying it. Callers must not modify b or
// reuse its backing storage while the returned string is in use.
func View(b []byte) string {
	//nolint:gosec // view into an append-only backing; appends never mutate written bytes
	return unsafe.String(unsafe.SliceData(b), len(b))
}
