package text

import (
	"slices"

	"github.com/nekrassov01/table/internal/param"
	"github.com/nekrassov01/table/internal/scope"
	"github.com/nekrassov01/table/internal/span"
	"github.com/nekrassov01/table/internal/value"
)

// noBars marks the absence of vertical bars.
const noBars = uint64(0)

// allBars marks vertical bars at every joint.
const allBars = ^uint64(0)

// compiler formats input values and resolves spans and visible boundaries into
// the logical rows consumed by a solver.
type compiler struct {
	input     configResult   // Resolved columns and input being compiled.
	state     *compilerState // Reusable compilation state.
	strings   *value.Store   // Storage for formatted and derived values.
	bodyStart int            // First body row in compilerState.rows; -1 before the body.
	err       error          // Structural input error from the current compilation.
	output    compilerResult // Logical table accumulated by this compiler.
}

// prepare sizes reusable storage and initializes span state from the resolved
// columns.
func (o *compiler) prepare() {
	state := o.state
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
	if cap(state.spanValues) < columnCount {
		state.spanValues = make([]string, columnCount)
	} else {
		state.spanValues = state.spanValues[:columnCount]
	}
	state.rowspans = scope.Masks{}
	state.colspans = scope.Masks{}
	for index := 0; index < min(columnCount, param.SpanLimit); index++ {
		col := &columns[index]
		state.rowspans.Mark(col.rowspan, index)
		state.colspans.Mark(col.colspan, index)
	}
	o.output.rowspanMask = state.rowspans.Resolve(ScopeBody)
	state.previousBody.Reset()
	state.lastBars = allBars
}

// compileHeader compiles rows from bottom to top so a vertical header span
// keeps its label on the row nearest the body.
func (o *compiler) compileHeader() {
	state := o.state
	header := o.input.header
	if len(header) == 0 {
		return
	}
	state.previousBand.Reset()
	rows := o.reserveBand(len(header))
	for i := len(header) - 1; i >= 0; i-- {
		rows[i] = o.compileBand(header[i], ScopeHeader)
	}
	bars := allBars
	for i := len(rows) - 1; i >= 0; i-- {
		o.setBars(&rows[i], bars, ScopeHeader)
		bars = rows[i].bars
	}
	o.output.header = rows
}

// compileBody compiles and retains body rows in input order.
func (o *compiler) compileBody(sources [][]any) {
	for i, source := range sources {
		o.compileRow(source, i)
		if o.err != nil {
			return
		}
	}
}

// compileFooter compiles rows from top to bottom so a vertical footer span
// keeps its label on the row nearest the body.
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
	for i := range footer {
		rows[i] = o.compileBand(footer[i], ScopeFooter)
	}
	bars := allBars
	for i := range rows {
		o.setBars(&rows[i], bars, ScopeFooter)
		bars = rows[i].bars
	}
	o.output.footer = rows
}

// compileBand compiles one header or footer row from labels.
func (o *compiler) compileBand(labels []string, sc Scope) row {
	config := &o.input
	state := o.state
	r := o.newRow()
	for index := range config.columns {
		col := &config.columns[index]
		var text string
		if index < config.option.indexOffset && sc == ScopeHeader {
			text = param.IndexHeader
		}
		if index >= config.option.indexOffset {
			if raw := index - config.option.indexOffset; raw < len(labels) {
				text = labels[raw]
			}
		}
		attr := col.transformer.attrs.Resolve(sc)
		if text == "" {
			// Empty labels remain empty and carry no attribute.
			attr = nil
		}
		r.cells[index] = cell{
			value: text,
			attr:  attr,
		}
	}
	o.setSpans(&r, sc, &state.previousBand)
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
	o.compileCells(r, source, rowIndex)
	o.setSpans(&r, ScopeBody, &state.previousBody)
	o.setBars(&r, o.output.lastBars, ScopeBody)
	o.output.lastBars = r.bars
	state.lastBars = r.bars
	state.rows = append(state.rows, r)
	o.output.body = state.rows[o.bodyStart:]
}

// compileCells formats the values and resolves the attributes of one body row.
func (o *compiler) compileCells(r row, source []any, rowIndex int) {
	config := &o.input
	for index := range config.columns {
		if index < config.option.indexOffset {
			text := value.Number(o.strings, int64(rowIndex)+1)
			r.cells[index] = cell{
				value: text,
			}
			continue
		}
		sourceIndex := index - config.option.indexOffset
		if sourceIndex >= len(source) {
			text := config.option.placeholder
			r.cells[index] = cell{
				value: text,
			}
			continue
		}
		transformer := &config.columns[index].transformer
		rawValue := source[sourceIndex]
		text := ""
		attr := transformer.attrs.Resolve(ScopeBody)
		if transformer.fn != nil {
			transformed, transformedAttr := transformer.fn(rawValue)
			if transformed != "" {
				text = transformed
			}
			if transformedAttr != nil && !config.option.plain {
				attr = transformedAttr
				// #nosec G115 -- Attr.size returns a non-negative byte count.
				attrLen := min(uint64(attr.size()), uint64(^uint32(0)))
				// #nosec G115 -- attrLen is clamped to uint32.
				o.output.attrLen = max(o.output.attrLen, uint32(attrLen))
			}
		}
		if text == "" {
			text = value.Format(o.strings, rawValue)
		}
		if text == "" {
			text = config.option.placeholder
			attr = nil
		}
		r.cells[index] = cell{
			value: text,
			attr:  attr,
		}
	}
}

// reserveBand reserves storage for n header or footer rows and returns them.
// A stream resolves its footer after initial storage has been prepared.
func (o *compiler) reserveBand(n int) []row {
	state := o.state
	cells := n * len(o.input.columns)
	state.cells = slices.Grow(state.cells, cells)
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

// setSpans marks vertical continuations and horizontally absorbed cells.
func (o *compiler) setSpans(r *row, sc Scope, previous *span.PreviousRow) {
	state := o.state
	rowspan := state.rowspans.Resolve(sc)
	colspan := state.colspans.Resolve(sc)
	if rowspan|colspan == 0 {
		return
	}
	values := state.spanValues[:len(r.cells)]
	for i := range r.cells {
		values[i] = r.cells[i].value
	}
	if rowspan != 0 {
		r.rowspans = span.Rowspans(rowspan, values, previous)
	}
	if colspan != 0 {
		r.colspans = span.Colspans(colspan, values, r.rowspans)
	}
}

// setBars resolves current colspans and inherits hidden boundaries only
// between adjacent rowspan continuations.
func (o *compiler) setBars(r *row, previousBars uint64, sc Scope) {
	r.bars = allBars
	if o.state.colspans.Resolve(sc) == 0 {
		return
	}
	for i := 1; i < len(r.cells); i++ {
		bit := uint64(1) << uint(i)
		if r.rowspans&bit != 0 {
			if previousBars&bit != 0 || r.rowspans&(bit>>1) == 0 {
				continue
			}
			r.bars &^= bit
			continue
		}
		if r.colspans&bit == 0 {
			continue
		}
		r.bars &^= bit
	}
}

// compilerResult is a logical table with resolved spans and boundaries.
type compilerResult struct {
	configResult

	header          []row  // Compiled header rows in top-to-bottom order.
	body            []row  // Compiled body rows in top-to-bottom order.
	footer          []row  // Compiled footer rows in top-to-bottom order.
	rowspanMask     uint64 // Columns configured for body rowspans.
	previousBars    uint64 // Boundaries inherited from the preceding body row.
	lastBars        uint64 // Final body row boundaries, or inherited boundaries without a body.
	attrLen         uint32 // Greatest dynamic attribute byte length used for capacity estimation.
	hasPreviousBody bool   // Whether a body row precedes this result.
}

// row holds cells and span-derived boundary masks for one logical row.
type row struct {
	cells    []cell // Resolved logical cells; a view into compilerState.cells.
	rowspans uint64 // The bitset of cells omitted by rowspans.
	colspans uint64 // The bitset of cells absorbed into the cell on their left.
	bars     uint64 // The bitset of visible vertical boundaries.
}

// cell holds the resolved value, attribute, and solved width of one logical
// cell.
type cell struct {
	value string // Resolved display value.
	attr  *Attr  // Optional display attribute.
	width int    // Solved display width for a non-empty single-line value.
}
