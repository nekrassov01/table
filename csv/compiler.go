package csv

import (
	"slices"

	"github.com/nekrassov01/table/internal/param"
	"github.com/nekrassov01/table/internal/value"
)

// compiler formats and quotes input values into the logical rows consumed by
// a painter.
type compiler struct {
	input     configResult   // Resolved columns and input being compiled.
	state     *compilerState // Reusable compilation state.
	strings   *value.Store   // Storage for formatted values.
	bodyStart int            // First body row in compilerState.rows; -1 before the body.
	err       error          // Structural input error from the current compilation.
	output    compilerResult // Compiled records accumulated by this compiler.
}

// prepare sizes reusable storage for the current pass.
func (o *compiler) prepare() {
	state := o.state
	columnCount := len(o.input.columns)
	rowCount := o.input.bodyRows + len(o.input.footer)
	if len(o.input.header) > 0 {
		rowCount++
	}
	cellCount := columnCount * rowCount
	if cap(state.rows) < rowCount {
		state.rows = make([]row, 0, rowCount)
	} else {
		state.rows = state.rows[:0]
	}
	if cap(state.cells) < cellCount {
		state.cells = make([]cell, 0, cellCount)
	} else {
		state.cells = state.cells[:0]
	}
	state.quotes = state.quotes[:0]
	if cap(state.values) < columnCount {
		state.values = make([]string, columnCount)
	} else {
		state.values = state.values[:columnCount]
		clear(state.values)
	}
}

// compileHeader compiles the optional header record.
func (o *compiler) compileHeader() {
	if o.err != nil || len(o.input.header) == 0 {
		return
	}
	o.output.header = o.compileBand(o.input.header, true)
}

// compileBody compiles and retains body records in input order.
func (o *compiler) compileBody(sources [][]any) {
	for index, source := range sources {
		o.compileRow(source, index)
		if o.err != nil {
			return
		}
	}
}

// compileFooter compiles footer records from top to bottom.
func (o *compiler) compileFooter() {
	footer := o.input.footer
	if o.err != nil || len(footer) == 0 {
		return
	}
	footerColumns := o.input.footerColumns + o.input.option.indexOffset
	if footerColumns > len(o.input.columns) {
		o.err = newColumnCountError(footerColumns, len(o.input.columns))
		return
	}
	rows := o.reserveBand(len(footer))
	for index := range footer {
		rows[index] = o.compileBand(footer[index], false)
	}
	o.output.footer = rows
}

// compileBand compiles one header or footer record from labels.
func (o *compiler) compileBand(labels []string, header bool) row {
	config := &o.input
	state := o.state
	r := o.newRow()
	values := state.values[:len(config.columns)]
	for index := range config.columns {
		text := ""
		if index < config.option.indexOffset && header {
			text = param.IndexHeader
		}
		if index >= config.option.indexOffset {
			if source := index - config.option.indexOffset; source < len(labels) {
				text = labels[source]
			}
		}
		values[index] = text
	}
	state.values = values
	o.compileCells(r)
	return r
}

// compileRow validates and compiles one body record.
func (o *compiler) compileRow(source []any, rowIndex int) {
	config := &o.input
	state := o.state
	rowColumns := len(source) + config.option.indexOffset
	columnCount := len(config.columns)
	if rowColumns > columnCount {
		o.err = newColumnCountError(rowColumns, columnCount)
		return
	}
	if o.bodyStart < 0 {
		o.bodyStart = len(state.rows)
	}
	r := o.newRow()
	values := state.values[:columnCount]
	for index := range config.columns {
		if index < config.option.indexOffset {
			values[index] = value.Number(o.strings, int64(rowIndex)+1)
			continue
		}
		sourceIndex := index - config.option.indexOffset
		if sourceIndex >= len(source) {
			values[index] = config.option.placeholder
			continue
		}
		rawValue := source[sourceIndex]
		text := ""
		if transformer := config.columns[index].transformer; transformer != nil {
			text = transformer(rawValue)
		}
		if text == "" {
			text = value.Format(o.strings, rawValue)
		}
		if text == "" {
			text = config.option.placeholder
		}
		values[index] = text
	}
	state.values = values
	o.compileCells(r)
	state.rows = append(state.rows, r)
	o.output.body = state.rows[o.bodyStart:]
}

// compileCells applies delimiter-aware quoting to the current row values.
func (o *compiler) compileCells(r row) {
	values := o.state.values[:len(r.cells)]
	for index := range r.cells {
		r.cells[index] = cell{
			value: o.quoteValue(values[index]),
		}
	}
}

// reserveBand reserves storage for n footer rows and returns it. A stream
// resolves its footer after initial storage has been prepared.
func (o *compiler) reserveBand(n int) []row {
	state := o.state
	state.cells = slices.Grow(state.cells, n*len(o.input.columns))
	start := len(state.rows)
	state.rows = slices.Grow(state.rows, n)[:start+n]
	return state.rows[start:]
}

// newRow reserves one logical row in reusable cell storage.
func (o *compiler) newRow() row {
	state := o.state
	start := len(state.cells)
	end := start + len(o.input.columns)
	state.cells = state.cells[:end]
	return row{
		cells: state.cells[start:end],
	}
}

// compilerResult pairs compiled records with their resolved configuration.
type compilerResult struct {
	configResult

	header row   // Optional compiled header record.
	body   []row // Compiled body records in input order.
	footer []row // Compiled footer records in top-to-bottom order.
}

// row holds compiled cells for one logical record.
type row struct {
	cells []cell // A view into compilerState.cells.
}

// cell holds one quoted display value.
type cell struct {
	value string // Display value with delimiter-aware quoting applied.
}
