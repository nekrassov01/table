package value

import "unsafe"

// View returns b as a string without copying it. Callers must not modify b or
// reuse its backing storage while the returned string is in use.
func View(b []byte) string {
	// #nosec G103 -- view into an append-only backing; appends never mutate written bytes
	return unsafe.String(unsafe.SliceData(b), len(b))
}
