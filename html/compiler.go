package html

import (
	"slices"

	"github.com/nekrassov01/table/internal/param"
	"github.com/nekrassov01/table/internal/scope"
	"github.com/nekrassov01/table/internal/span"
	"github.com/nekrassov01/table/internal/value"
)

// compiler formats and escapes input values, then resolves them into the
// logical rows consumed by a solver.
type compiler struct {
	input     configResult   // Resolved columns and input being compiled.
	state     *compilerState // Reusable compilation state.
	strings   *value.Store   // Storage for formatted values.
	bodyStart int            // First body row in compilerState.rows; -1 before the body.
	err       error          // Structural input error from the current compilation.
	output    compilerResult // Compiled table accumulated by this compiler.
}

// prepare sizes reusable storage, initializes span state, and escapes the
// caption.
func (o *compiler) prepare() {
	state := o.state
	option := o.input.option
	columns := o.input.columns
	columnCount := len(columns)
	rowCount := o.input.bodyRows + len(o.input.header) + len(o.input.footer)
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
	state.escapes = state.escapes[:0]
	if cap(state.values) < columnCount {
		state.values = make([]string, columnCount)
	} else {
		state.values = state.values[:columnCount]
	}
	if cap(state.columnSizes) < columnCount {
		state.columnSizes = make([]int, columnCount)
	} else {
		state.columnSizes = state.columnSizes[:columnCount]
		clear(state.columnSizes)
	}
	state.rowspans = scope.Masks{}
	state.colspans = scope.Masks{}
	for index := 0; index < min(columnCount, param.SpanLimit); index++ {
		column := &o.input.columns[index]
		state.rowspans.Mark(column.rowspan, index)
		state.colspans.Mark(column.colspan, index)
	}
	state.previousBody.Reset()
	o.output.columnSizes = state.columnSizes
	o.output.caption, state.escapes = escapeValue(state.escapes, option.caption)
}

// compileHeader compiles header rows from top to bottom because an HTML
// rowspan belongs to its leading cell.
func (o *compiler) compileHeader() {
	state := o.state
	header := o.input.header
	if len(header) == 0 {
		return
	}
	state.previousBand.Reset()
	rows := o.reserveBand(len(header))
	for index := range header {
		rows[index] = o.compileBand(header[index], ScopeHeader)
	}
	o.output.header = rows
}

// compileBody compiles and retains body rows in input order.
func (o *compiler) compileBody(sources [][]any) {
	for index, source := range sources {
		o.compileRow(source, index)
		if o.err != nil {
			return
		}
	}
}

// compileFooter compiles footer rows from top to bottom because an HTML
// rowspan belongs to its leading cell.
func (o *compiler) compileFooter() {
	state := o.state
	footer := o.input.footer
	if o.err != nil || len(footer) == 0 {
		return
	}
	footerColumns := o.input.footerColumns + o.input.option.indexOffset
	if footerColumns > len(o.input.columns) {
		o.err = newColumnCountError(footerColumns, len(o.input.columns))
		return
	}
	state.previousBand.Reset()
	rows := o.reserveBand(len(footer))
	for index := range footer {
		rows[index] = o.compileBand(footer[index], ScopeFooter)
	}
	o.output.footer = rows
}

// compileBand compiles one header or footer row from labels.
func (o *compiler) compileBand(labels []string, sc Scope) row {
	config := &o.input
	state := o.state
	r := o.newRow()
	values := state.values[:len(config.columns)]
	for index := range config.columns {
		column := &config.columns[index]
		transformer := &column.transformer
		var text string
		if index < config.option.indexOffset && sc == ScopeHeader {
			text = param.IndexHeader
		}
		if index >= config.option.indexOffset {
			if source := index - config.option.indexOffset; source < len(labels) {
				text = labels[source]
			}
		}
		color := transformer.colors.Resolve(sc)
		decoration := transformer.decorations.Resolve(sc)
		if text == "" {
			color = nil
			decoration = nil
		}
		values[index] = text
		r.cells[index] = cell{
			color:      color,
			decoration: decoration,
			colspan:    1,
		}
	}
	state.values = values
	o.setSpans(&r, sc, &state.previousBand)
	o.compileCells(r)
	return r
}

// compileRow validates and compiles one body row.
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
	values := state.values[:len(config.columns)]
	for index := range config.columns {
		column := &config.columns[index]
		transformer := &column.transformer
		if index < config.option.indexOffset {
			values[index] = value.Number(o.strings, int64(rowIndex)+1)
			r.cells[index] = cell{
				colspan: 1,
			}
			continue
		}
		sourceIndex := index - config.option.indexOffset
		if sourceIndex >= len(source) {
			values[index] = config.option.placeholder
			r.cells[index] = cell{
				colspan: 1,
			}
			continue
		}
		rawValue := source[sourceIndex]
		text := ""
		color := transformer.colors.Resolve(ScopeBody)
		decoration := transformer.decorations.Resolve(ScopeBody)
		if transformer.fn != nil {
			transformed, transformedColor, transformedDecoration := transformer.fn(rawValue)
			if transformed != "" {
				text = transformed
			}
			if transformedColor != nil {
				color = transformedColor
			}
			if transformedDecoration != nil {
				decoration = transformedDecoration
			}
		}
		if text == "" {
			text = value.Format(o.strings, rawValue)
		}
		if text == "" {
			text = config.option.placeholder
			color = nil
			decoration = nil
		}
		values[index] = text
		r.cells[index] = cell{
			color:      color,
			decoration: decoration,
			colspan:    1,
		}
	}
	state.values = values
	o.setSpans(&r, ScopeBody, &state.previousBody)
	o.compileCells(r)
	state.rows = append(state.rows, r)
	o.output.body = state.rows[o.bodyStart:]
}

// compileCells escapes the current row values and applies its inner markup.
func (o *compiler) compileCells(r row) {
	state := o.state
	values := state.values[:len(r.cells)]
	for index := range r.cells {
		compiled := &r.cells[index]
		if r.rowspans&(1<<uint(index)) != 0 {
			*compiled = cell{
				colspan: 1,
			}
		} else {
			compiled.value, state.escapes = escapeValue(state.escapes, values[index])
			if compiled.color.IsZero() {
				compiled.color = nil
			}
			if compiled.decoration.IsZero() {
				compiled.decoration = nil
			}
		}
		compiled.size = len(compiled.value)
		if compiled.color != nil {
			compiled.size += len(compiled.color.Prefix) + len(compiled.color.Suffix)
		}
		if compiled.decoration != nil {
			compiled.size += len(compiled.decoration.Prefix) + len(compiled.decoration.Suffix)
		}
		if compiled.size > state.columnSizes[index] {
			state.columnSizes[index] = compiled.size
		}
	}
}

// reserveBand reserves storage for n header or footer rows and returns them.
// A stream resolves its footer after initial storage has been prepared.
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

// setSpans resolves vertical continuations and horizontal merge candidates for
// one logical row.
func (o *compiler) setSpans(r *row, sc Scope, previous *span.PreviousRow) {
	state := o.state
	rowspan := state.rowspans.Resolve(sc)
	colspan := state.colspans.Resolve(sc)
	if rowspan|colspan == 0 {
		return
	}
	values := state.values[:len(r.cells)]
	if rowspan != 0 {
		r.rowspans = span.Rowspans(rowspan, values, previous)
	}
	if colspan != 0 {
		r.colspans = span.Colspans(colspan, values, r.rowspans)
	}
}

// compilerResult is a logical HTML table with escaped values and span state.
type compilerResult struct {
	configResult

	header          []row  // Compiled header rows in top-to-bottom order.
	body            []row  // Compiled body rows in top-to-bottom order.
	footer          []row  // Compiled footer rows in top-to-bottom order.
	columnSizes     []int  // Greatest value and inner-markup byte size by column.
	caption         string // Escaped caption text.
	hasPreviousBody bool   // Whether a body row precedes this result.
}

// row holds compiled cells and unresolved span facts for one logical row.
type row struct {
	cells    []cell // A view into compilerState.cells.
	rowspans uint64 // Cells that continue a vertical span from the row above.
	colspans uint64 // Cells that are candidates for absorption into the cell on their left.
}

// cell holds an escaped display value and its span state.
type cell struct {
	value      string      // Escaped display value.
	size       int         // Bytes written for the value and its inner markup.
	color      *Color      // Optional color markup inside decoration markup.
	decoration *Decoration // Optional outer decoration markup.
	rowspan    int         // Resolved rowspan count; zero or one needs no attribute.
	colspan    int         // Resolved colspan count; zero when another cell emits this cell.
}
