package text

import (
	"io"

	"github.com/nekrassov01/table/internal/scope"
)

// DefaultPlaceholder is the default placeholder for missing or empty body
// cells.
const DefaultPlaceholder = " "

// config resolves construction options and pass data into logical columns.
type config struct {
	bodyColumns int          // Body width used to resolve logical columns.
	state       *configState // Reusable config state.
	output      configResult // Pass data paired with resolved columns.
}

// prepare resolves logical columns. A non-empty header fixes their count;
// otherwise the current body and footer determine it.
func (o *config) prepare() {
	result := &o.output
	option := result.option
	headerColumns := maxColumns(result.header)
	footerColumns := maxColumns(result.footer)
	columnCount := headerColumns
	if columnCount == 0 {
		columnCount = max(o.bodyColumns, footerColumns)
	}
	if columnCount > 0 {
		columnCount += option.indexOffset
	}
	columns := o.state.columns
	if cap(columns) < columnCount {
		columns = make([]column, columnCount)
	} else {
		columns = columns[:columnCount]
	}
	for index := 0; index < min(option.indexOffset, columnCount); index++ {
		columns[index] = defaultColumn()
	}
	defaults := defaultColumn()
	if option.columns.defaults != nil {
		defaults = *option.columns.defaults
	}
	for index := option.indexOffset; index < columnCount; index++ {
		columns[index] = defaults
	}
	if option.indexOffset < columnCount {
		copy(columns[option.indexOffset:], option.columns.values)
	}
	o.state.columns = columns
	result.columns = columns
	result.footerColumns = footerColumns
}

// configResult carries pass data and resolved logical columns.
type configResult struct {
	option        *option    // Options fixed at construction.
	header        [][]string // Header rows in top-to-bottom order.
	footer        [][]string // Footer rows in top-to-bottom order.
	bodyRows      int        // Number of body rows in this pass.
	columns       []column   // Resolved column settings in logical order.
	footerColumns int        // Footer width resolved for the current config pass.
}

// option holds settings fixed when a Table or Stream is constructed.
type option struct {
	style       Style             // Border and content style.
	placeholder string            // Text for a missing value.
	header      [][]string        // Header rows in top-to-bottom order.
	footer      func() [][]string // Generates footer rows for Render or Close.
	caption     string            // Caption text.
	columns     columnSet         // Input columns and their defaults.
	indexOffset int               // Synthetic leading column count: 0 or 1.
	indexWidth  int               // Minimum index width; zero adds no explicit minimum.
	captionSide CaptionSide       // Caption position.
	compact     bool              // Whether body separators are omitted.
	autoFit     bool              // Whether columns are fitted to the terminal width.
	plain       bool              // Whether ANSI attributes are disabled.
}

// apply sets defaults, applies opts in order, and resolves writer-dependent
// behavior.
func (o *option) apply(w io.Writer, minIndexWidth int, opts ...Option) {
	o.style = StyleLight
	o.placeholder = DefaultPlaceholder
	for _, opt := range opts {
		opt(o)
	}
	if o.indexOffset != 0 {
		o.indexWidth = max(o.indexWidth, minIndexWidth)
	}
	columns := &o.columns
	o.plain = !isTerminal(w)
	if o.plain {
		o.style.Border.Attr = nil
		o.style.Content = ContentStyle{}
		if columns.defaults != nil {
			columns.defaults.transformer.attrs = scope.Scopes[*Attr]{}
		}
		for i := range columns.values {
			columns.values[i].transformer.attrs = scope.Scopes[*Attr]{}
		}
	}
	if !o.autoFit {
		return
	}
	if defaults := columns.defaults; defaults != nil && (defaults.limit > 0 || defaults.truncate) {
		o.autoFit = false
		return
	}
	for i := range columns.values {
		col := &columns.values[i]
		if col.limit > 0 || col.truncate {
			o.autoFit = false
			return
		}
	}
}

// columnSet holds explicit input columns and the default used by AllColumns.
type columnSet struct {
	values   []column // Explicit columns in input order.
	defaults *column  // Column inherited by every input position.
}

// apply updates every input column selected by selector.
func (o *columnSet) apply(selector ColumnSelector, fn func(*column)) {
	if selector.all {
		if o.defaults == nil {
			defaults := defaultColumn()
			o.defaults = &defaults
		}
		fn(o.defaults)
		for i := range o.values {
			fn(&o.values[i])
		}
		return
	}
	columnCount := 0
	for _, index := range selector.indexes {
		columnCount = max(columnCount, index+1)
	}
	existing := len(o.values)
	if columnCount > existing {
		o.values = append(o.values, make([]column, columnCount-existing)...)
		defaults := defaultColumn()
		if o.defaults != nil {
			defaults = *o.defaults
		}
		for i := existing; i < columnCount; i++ {
			o.values[i] = defaults
		}
	}
	for _, index := range selector.indexes {
		if index >= 0 {
			fn(&o.values[index])
		}
	}
}

// column defines the behavior of one resolved logical column.
type column struct {
	transformer transformer             // Value transformation and attributes.
	aligns      scope.Scopes[AlignSide] // Alignment by table part.
	limit       int                     // Configured display width; zero is unconstrained.
	lPad        int                     // Left padding width.
	rPad        int                     // Right padding width.
	rowspan     Scope                   // Parts that span equal values vertically.
	colspan     Scope                   // Parts that span equal values horizontally.
	truncate    bool                    // Whether overflow is truncated instead of wrapped.
}

// transformer holds value transformation and attributes for a column.
type transformer struct {
	attrs scope.Scopes[*Attr]       // Attributes applied in each table part.
	fn    func(any) (string, *Attr) // Per-cell transformer for body values.
}

// maxColumns returns the greatest column count among rows.
func maxColumns(rows [][]string) int {
	columnCount := 0
	for _, row := range rows {
		columnCount = max(columnCount, len(row))
	}
	return columnCount
}

// defaultColumn returns an unconfigured input column.
func defaultColumn() column {
	return column{
		lPad: 1,
		rPad: 1,
	}
}
