package value

import (
	"strconv"

	"github.com/nekrassov01/table/internal/unsafe"
)

// Store builds formatted values in shared backing storage and returns
// zero-copy string views. Callers must discard the views before Reset because
// subsequent appends may overwrite their backing bytes.
type Store struct {
	buf []byte
}

// Mark returns the current length for a later call to Since.
func (o *Store) Mark() int {
	return len(o.buf)
}

// Reset empties the store. Existing views must not be used after Reset.
func (o *Store) Reset() {
	o.buf = o.buf[:0]
}

// Since returns an unsafe string view of the bytes appended since mark.
func (o *Store) Since(mark int) string {
	return unsafe.View(o.buf[mark:])
}

// AppendString appends s.
func (o *Store) AppendString(s string) {
	o.grow()
	o.buf = append(o.buf, s...)
}

// AppendBytes appends b.
func (o *Store) AppendBytes(b []byte) {
	o.grow()
	o.buf = append(o.buf, b...)
}

// AppendInt appends the decimal text of x.
func (o *Store) AppendInt(x int64) {
	o.grow()
	o.buf = strconv.AppendInt(o.buf, x, 10)
}

// AppendUint appends the decimal text of x.
func (o *Store) AppendUint(x uint64) {
	o.grow()
	o.buf = strconv.AppendUint(o.buf, x, 10)
}

// AppendFloat appends the shortest decimal text of x.
func (o *Store) AppendFloat(x float64, bitSize int) {
	o.grow()
	o.buf = strconv.AppendFloat(o.buf, x, 'g', -1, bitSize)
}

// grow gives a new store enough initial capacity to avoid repeated growth for
// small tables.
func (o *Store) grow() {
	if cap(o.buf) == 0 {
		o.buf = make([]byte, 0, 128)
	}
}
