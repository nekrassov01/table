package markdown

import (
	"github.com/nekrassov01/table/internal/param"
	"github.com/nekrassov01/table/internal/span"
	"github.com/nekrassov01/table/internal/value"
	"github.com/nekrassov01/table/internal/width"
)

// compiler formats and escapes input values, then resolves them into the
// logical rows consumed by a solver.
type compiler struct {
	input     configResult   // Resolved columns and input being compiled.
	state     *compilerState // Reusable compilation state.
	strings   *value.Store   // Storage for formatted and derived values.
	bodyStart int            // First body row in compilerState.rows; -1 before the body.
	err       error          // Structural input error from the current compilation.
	output    compilerResult // Compiled table accumulated by this compiler.
}

// prepare sizes reusable storage and initializes body span state from the
// resolved columns.
func (o *compiler) prepare() {
	state := o.state
	columns := o.input.columns
	columnCount := len(columns)
	rowCount := o.input.bodyRows + 1
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
	state.rowspans = 0
	state.colspans = 0
	for index := 0; index < min(columnCount, param.SpanLimit); index++ {
		column := &columns[index]
		if column.rowspan {
			state.rowspans |= 1 << uint(index)
		}
		if column.colspan {
			state.colspans |= 1 << uint(index)
		}
	}
	state.previousBody.Reset()
}

// compileHeader compiles the required header row.
func (o *compiler) compileHeader() {
	config := &o.input
	state := o.state
	r := o.newRow()
	values := state.values[:len(config.columns)]
	for index := range config.columns {
		column := &config.columns[index]
		transformer := &column.transformer
		text := param.IndexHeader
		if index >= config.option.indexOffset {
			text = config.header[index-config.option.indexOffset]
		}
		color := transformer.colors.Resolve(ScopeHeader)
		decoration := transformer.decorations.Resolve(ScopeHeader)
		if text == "" {
			color = nil
			decoration = nil
		}
		values[index] = text
		r.cells[index] = cell{
			color:      color,
			decoration: decoration,
		}
	}
	state.values = values
	o.compileCells(r)
	state.rows = append(state.rows, r)
	o.output.header = r
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

// compileRow compiles one body row and updates span continuation state.
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
			r.cells[index] = cell{}
			continue
		}
		sourceIndex := index - config.option.indexOffset
		if sourceIndex >= len(source) {
			values[index] = config.option.placeholder
			r.cells[index] = cell{}
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
		}
	}
	state.values = values
	o.setSpans(&r, &state.previousBody)
	o.compileCells(r)
	state.rows = append(state.rows, r)
	o.output.body = state.rows[o.bodyStart:]
}

// compileCells escapes the current row values and applies its inner markup.
func (o *compiler) compileCells(r row) {
	state := o.state
	values := state.values[:len(r.cells)]
	for index := range r.cells {
		bit := uint64(1) << uint(index)
		if r.rowspans&bit != 0 || r.colspans&bit != 0 {
			r.cells[index] = cell{}
			continue
		}
		compiled := &r.cells[index]
		value := values[index]
		color := compiled.color
		decoration := compiled.decoration
		*compiled = cell{}
		markup := 0
		if ticks := resolveTicks(decoration, value); ticks > 0 {
			compiled.ticks = ticks
			markup = ticks * 2
			value, state.escapes = escapeCode(state.escapes, value)
		} else {
			value, state.escapes = escapeValue(state.escapes, value)
			if !decoration.IsZero() {
				compiled.decoration = decoration
				markup = len(decoration.Prefix) + len(decoration.Suffix)
			}
		}
		compiled.value = value
		compiled.width = width.StringWidth(value) + markup
		compiled.size = len(value) + markup
		if !color.IsZero() {
			markup = len(color.Prefix) + len(color.Suffix)
			compiled.color = color
			compiled.width += markup
			compiled.size += markup
		}
	}
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

// setSpans resolves vertical continuations and horizontal absorptions for one
// body row.
func (o *compiler) setSpans(r *row, previous *span.PreviousRow) {
	state := o.state
	if state.rowspans|state.colspans == 0 {
		return
	}
	values := state.values[:len(r.cells)]
	if state.rowspans != 0 {
		r.rowspans = span.Rowspans(state.rowspans, values, previous)
	}
	if state.colspans != 0 {
		r.colspans = span.Colspans(state.colspans, values, r.rowspans)
	}
}

// compilerResult pairs compiled rows with their resolved configuration.
type compilerResult struct {
	configResult

	header row   // Required header row.
	body   []row // Body rows compiled in input order.
}

// row holds compiled cells and span facts for one logical row.
type row struct {
	cells    []cell // Compiled cells; a view into compilerState.cells.
	rowspans uint64 // The bitset of cells omitted by rowspans.
	colspans uint64 // The bitset of cells absorbed into the cell on their left.
}

// cell holds an escaped value and the markup written around it.
type cell struct {
	value      string      // Escaped display value.
	width      int         // Display width including visible markup bytes.
	size       int         // Bytes written for the value and its markup.
	color      *Color      // Optional HTML color span around value or code markup.
	decoration *Decoration // Optional non-code Markdown decoration.
	ticks      int         // Code-fence length; zero outside code spans.
}
