package markdown

import (
	"io"

	"github.com/nekrassov01/table/internal/repeat"
)

// painter writes solved logical rows as a GFM table.
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
	// Leading pipe and trailing newline.
	size := 2
	for index := range o.input.metrics {
		metric := &o.input.metrics[index]
		// One space on each side and the trailing pipe surround each cell.
		size += max(metric.box.width, metric.separator.width) + 3
	}
	size = max(size, o.lineSize(&o.input.header))
	for index := range o.input.body {
		size = max(size, o.lineSize(&o.input.body[index]))
	}
	if cap(state.lineBacking) < size {
		state.lineBacking = make([]byte, size)
	}
	state.line = state.lineBacking[:0]
}

// lineSize returns the number of bytes required to paint a row.
func (o *painter) lineSize(r *row) int {
	// Leading pipe and trailing newline.
	size := 2
	for index := range r.cells {
		cell := &r.cells[index]
		box := &o.input.metrics[index].box
		padding := max(box.width-cell.width, 0)
		// One space on each side and the trailing pipe surround each cell.
		size += cell.size + padding + 3
	}
	return size
}

// paintHeader writes the required header row and its delimiter row.
func (o *painter) paintHeader() {
	o.paintRow(o.input.header)
	o.paintSeparator()
}

// paintBody writes compiled body rows in input order.
func (o *painter) paintBody() {
	for index := range o.input.body {
		if o.err != nil {
			return
		}
		o.paintRow(o.input.body[index])
	}
}

// paintSeparator writes the delimiter row between the header and body.
func (o *painter) paintSeparator() {
	o.resetLine()
	o.state.line = append(o.state.line, '|')
	for index := range o.input.metrics {
		metric := &o.input.metrics[index]
		o.state.line = append(o.state.line, ' ')
		o.paintSeparatorCell(metric)
		o.state.line = append(o.state.line, ' ', '|')
	}
	o.writeNewline()
	o.flushLine()
}

// paintSeparatorCell writes one delimiter cell.
func (o *painter) paintSeparatorCell(metric *columnMetric) {
	separator := &metric.separator
	n := max(metric.box.width, separator.width)
	if separator.lead {
		o.state.line = append(o.state.line, ':')
		n--
	}
	if separator.trail {
		n--
	}
	o.writeDashes(n)
	if separator.trail {
		o.state.line = append(o.state.line, ':')
	}
}

// paintRow writes one logical row.
func (o *painter) paintRow(r row) {
	o.resetLine()
	o.state.line = append(o.state.line, '|')
	for index := range r.cells {
		box := &o.input.metrics[index].box
		o.state.line = append(o.state.line, ' ')
		o.paintCell(&r.cells[index], box)
		o.state.line = append(o.state.line, ' ', '|')
	}
	o.writeNewline()
	o.flushLine()
}

// paintCell writes one value with its resolved alignment and padding.
func (o *painter) paintCell(cell *cell, box *box) {
	gap := box.width - cell.width
	switch box.align {
	case AlignRight:
		o.writeSpaces(gap)
		o.writeValue(cell)
	case AlignCenter:
		left := gap / 2
		o.writeSpaces(left)
		o.writeValue(cell)
		o.writeSpaces(gap - left)
	default:
		o.writeValue(cell)
		o.writeSpaces(gap)
	}
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
	// A color span must wrap a code fence; inside it would render as code.
	if cell.ticks > 0 {
		if cell.color != nil {
			o.state.line = append(o.state.line, cell.color.Prefix...)
		}
		o.writeBackticks(cell.ticks)
		o.state.line = append(o.state.line, cell.value...)
		o.writeBackticks(cell.ticks)
		if cell.color != nil {
			o.state.line = append(o.state.line, cell.color.Suffix...)
		}
		return
	}
	if cell.decoration != nil {
		o.state.line = append(o.state.line, cell.decoration.Prefix...)
	}
	if cell.color != nil {
		o.state.line = append(o.state.line, cell.color.Prefix...)
	}
	o.state.line = append(o.state.line, cell.value...)
	if cell.color != nil {
		o.state.line = append(o.state.line, cell.color.Suffix...)
	}
	if cell.decoration != nil {
		o.state.line = append(o.state.line, cell.decoration.Suffix...)
	}
}

// writeBackticks appends n code-fence bytes.
func (o *painter) writeBackticks(n int) {
	for range n {
		o.state.line = append(o.state.line, '`')
	}
}

// writeDashes appends n delimiter bytes.
func (o *painter) writeDashes(n int) {
	o.state.line = repeat.AppendDashes(o.state.line, n)
}

// writeSpaces appends n padding bytes.
func (o *painter) writeSpaces(n int) {
	o.state.line = repeat.AppendSpaces(o.state.line, n)
}

// writeNewline appends a line terminator.
func (o *painter) writeNewline() {
	o.state.line = append(o.state.line, '\n')
}
