package backlog

import (
	"io"

	"github.com/nekrassov01/table/internal/repeat"
)

// bandMarker marks header and footer values in Backlog notation.
const bandMarker = "~"

// painter writes solved logical rows as a Backlog table.
type painter struct {
	input solverResult  // Compiled rows and column geometry.
	state *painterState // Reusable painting state.
	w     io.Writer     // Output destination.
	err   error         // Sticky output error.
}

// prepare sizes the reusable output buffer for the maximum line size in the
// current pass.
func (o *painter) prepare() {
	state := o.state
	if cap(state.line) > cap(state.lineBacking) {
		state.lineBacking = state.line[:0]
	}
	size := 2
	for index := range o.input.header {
		size = max(size, o.lineSize(&o.input.header[index], true))
	}
	for index := range o.input.body {
		size = max(size, o.lineSize(&o.input.body[index], false))
	}
	for index := range o.input.footer {
		size = max(size, o.lineSize(&o.input.footer[index], true))
	}
	if cap(state.lineBacking) < size {
		state.lineBacking = make([]byte, size)
	}
	state.line = state.lineBacking[:0]
}

// lineSize returns the number of bytes required to paint a row.
func (o *painter) lineSize(r *row, band bool) int {
	// Leading pipe and trailing newline.
	size := 2
	for index := range r.cells {
		cell := &r.cells[index]
		width := cell.width
		cellSize := cell.size
		if band {
			width += len(bandMarker)
			cellSize += len(bandMarker)
		}
		padding := max(o.input.metrics[index].box.width-width, 0)
		// Two spaces and the trailing pipe frame each cell. Band rows keep both
		// spaces after the value so the marker stays adjacent to the leading pipe.
		size += cellSize + padding + 3
	}
	return size
}

// paintHeader writes compiled header rows in top-to-bottom order.
func (o *painter) paintHeader() {
	o.paintBand(o.input.header)
}

// paintBody writes compiled body rows in input order.
func (o *painter) paintBody() {
	for index := range o.input.body {
		if o.err != nil {
			return
		}
		o.paintRow(o.input.body[index], false)
	}
}

// paintFooter writes compiled footer rows in top-to-bottom order.
func (o *painter) paintFooter() {
	o.paintBand(o.input.footer)
}

// paintBand writes header or footer rows with a marker before every cell.
func (o *painter) paintBand(rows []row) {
	for index := range rows {
		if o.err != nil {
			return
		}
		o.paintRow(rows[index], true)
	}
}

// paintRow writes one logical row.
func (o *painter) paintRow(r row, band bool) {
	o.resetLine()
	o.state.line = append(o.state.line, '|')
	for index := range r.cells {
		if !band {
			o.state.line = append(o.state.line, ' ')
		}
		o.paintCell(&r.cells[index], &o.input.metrics[index].box, band)
		o.state.line = append(o.state.line, ' ')
		if band {
			o.state.line = append(o.state.line, ' ')
		}
		o.state.line = append(o.state.line, '|')
	}
	o.writeNewline()
	o.flushLine()
}

// paintCell writes one value and its right padding.
func (o *painter) paintCell(cell *cell, box *box, band bool) {
	width := cell.width
	if band {
		o.state.line = append(o.state.line, bandMarker...)
		width += len(bandMarker)
	}
	o.writeValue(cell)
	o.writeSpaces(box.width - width)
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

// writeValue writes a cell's color, decoration, and escaped value.
func (o *painter) writeValue(cell *cell) {
	if cell.value == "" {
		return
	}
	if cell.color != nil {
		o.state.line = append(o.state.line, cell.color.Prefix...)
	}
	if cell.decoration != nil {
		o.state.line = append(o.state.line, cell.decoration.Prefix...)
	}
	o.state.line = append(o.state.line, cell.value...)
	if cell.decoration != nil {
		o.state.line = append(o.state.line, cell.decoration.Suffix...)
	}
	if cell.color != nil {
		o.state.line = append(o.state.line, cell.color.Suffix...)
	}
}

// writeSpaces appends n padding bytes.
func (o *painter) writeSpaces(n int) {
	o.state.line = repeat.AppendSpaces(o.state.line, n)
}

// writeNewline appends a line terminator.
func (o *painter) writeNewline() {
	o.state.line = append(o.state.line, '\n')
}
