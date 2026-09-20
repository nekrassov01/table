package value

import (
	"fmt"
	"strconv"
)

// Store builds formatted values in shared backing storage and returns
// zero-copy string views. Append methods return the text they appended; callers
// may ignore the result when building a larger value. Discard views before
// Reset because subsequent appends may overwrite their backing bytes.
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

// Since returns a zero-copy string view of the bytes appended since mark.
func (o *Store) Since(mark int) string {
	return View(o.buf[mark:])
}

// AppendString appends s.
func (o *Store) AppendString(s string) string {
	mark := o.Mark()
	o.grow()
	o.buf = append(o.buf, s...)
	return o.Since(mark)
}

// AppendBytes appends b.
func (o *Store) AppendBytes(b []byte) string {
	mark := o.Mark()
	o.grow()
	o.buf = append(o.buf, b...)
	return o.Since(mark)
}

// AppendInt appends the decimal text of x.
func (o *Store) AppendInt(x int64) string {
	mark := o.Mark()
	o.grow()
	o.buf = strconv.AppendInt(o.buf, x, 10)
	return o.Since(mark)
}

// AppendUint appends the decimal text of x.
func (o *Store) AppendUint(x uint64) string {
	mark := o.Mark()
	o.grow()
	o.buf = strconv.AppendUint(o.buf, x, 10)
	return o.Since(mark)
}

// AppendFloat appends the shortest decimal text of x.
func (o *Store) AppendFloat(x float64, bitSize int) string {
	mark := o.Mark()
	o.grow()
	o.buf = strconv.AppendFloat(o.buf, x, 'g', -1, bitSize)
	return o.Since(mark)
}

// AppendDefault appends the fmt.Sprint representation of v.
func (o *Store) AppendDefault(v any) string {
	mark := o.Mark()
	o.grow()
	o.buf = fmt.Append(o.buf, v)
	return o.Since(mark)
}

// grow gives a new store enough initial capacity to avoid repeated growth for
// small tables.
func (o *Store) grow() {
	if cap(o.buf) == 0 {
		o.buf = make([]byte, 0, 128)
	}
}
