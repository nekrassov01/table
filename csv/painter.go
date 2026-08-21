package csv

import (
	"io"
	"unicode/utf8"
)

// painter writes compiled logical rows as delimiter-separated records.
type painter struct {
	input compilerResult // Compiled records to paint.
	state *painterState  // Reusable painting state.
	w     io.Writer      // Output destination.
	err   error          // Sticky output error.
}

// prepare sizes the reusable output buffer for the largest record in the
// current pass.
func (o *painter) prepare() {
	state := o.state
	if cap(state.line) > cap(state.lineBacking) {
		state.lineBacking = state.line[:0]
	}
	size := o.lineSize(&o.input.header)
	for index := range o.input.body {
		size = max(size, o.lineSize(&o.input.body[index]))
	}
	for index := range o.input.footer {
		size = max(size, o.lineSize(&o.input.footer[index]))
	}
	if cap(state.lineBacking) < size {
		state.lineBacking = make([]byte, size)
	}
	state.line = state.lineBacking[:0]
}

// lineSize returns the number of bytes required to paint one record.
func (o *painter) lineSize(r *row) int {
	if len(r.cells) == 0 {
		return 0
	}
	// Delimiters between adjacent fields and one trailing line terminator.
	size := (len(r.cells)-1)*utf8.RuneLen(o.input.option.delimiter) + 1
	if o.input.option.crlf {
		size++
	}
	for index := range r.cells {
		size += len(r.cells[index].value)
	}
	return size
}

// paintHeader writes the optional header record.
func (o *painter) paintHeader() {
	o.paintRow(o.input.header)
}

// paintBody writes compiled body records in input order.
func (o *painter) paintBody() {
	for index := range o.input.body {
		if o.err != nil {
			return
		}
		o.paintRow(o.input.body[index])
	}
}

// paintFooter writes compiled footer records in top-to-bottom order.
func (o *painter) paintFooter() {
	for index := range o.input.footer {
		if o.err != nil {
			return
		}
		o.paintRow(o.input.footer[index])
	}
}

// paintRow writes one logical row.
func (o *painter) paintRow(r row) {
	if len(r.cells) == 0 {
		return
	}
	o.resetLine()
	delimiter := o.input.option.delimiter
	for index := range r.cells {
		if index > 0 {
			o.state.line = utf8.AppendRune(o.state.line, delimiter)
		}
		o.state.line = append(o.state.line, r.cells[index].value...)
	}
	o.writeNewline()
	o.flushLine()
}

// resetLine resets the current output line.
func (o *painter) resetLine() {
	o.state.line = o.state.line[:0]
}

// flushLine writes the current output line.
func (o *painter) flushLine() {
	if o.err != nil || len(o.state.line) == 0 {
		return
	}
	_, err := o.w.Write(o.state.line)
	if err != nil {
		o.err = newWriteError(err)
	}
}

// writeNewline appends the configured line terminator.
func (o *painter) writeNewline() {
	if o.input.option.crlf {
		o.state.line = append(o.state.line, '\r')
	}
	o.state.line = append(o.state.line, '\n')
}
