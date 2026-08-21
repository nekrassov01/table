package html

import "github.com/nekrassov01/table/internal/scope"

// DefaultPlaceholder is the default placeholder for missing or empty body
// cells.
const DefaultPlaceholder = ""

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
		columns[index] = column{}
	}
	defaults := column{}
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
	attrs       TableAttr         // Classes and inline styles by element.
	placeholder string            // Text for a missing value.
	header      [][]string        // Header rows in top-to-bottom order.
	footer      func() [][]string // Generates footer rows for Render or Close.
	caption     string            // Caption text.
	columns     columnSet         // Input columns and their defaults.
	indexOffset int               // Synthetic leading column count: 0 or 1.
	captionSide CaptionSide       // Caption position.
}

// apply sets defaults and applies opts in order.
func (o *option) apply(opts ...Option) {
	o.placeholder = DefaultPlaceholder
	for _, opt := range opts {
		opt(o)
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
			defaults := column{}
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
		defaults := column{}
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
	transformer transformer             // Value transformation and inner markup.
	attrs       scope.Scopes[Attr]      // Cell attributes by table part.
	aligns      scope.Scopes[AlignSide] // Alignment by table part.
	rowspan     Scope                   // Parts that span equal values vertically.
	colspan     Scope                   // Parts that span equal values horizontally.
}

// transformer holds body-value transformation and inner markup for a column.
type transformer struct {
	colors      scope.Scopes[*Color]      // Color by table part.
	decorations scope.Scopes[*Decoration] // Decoration by table part.
	fn          func(any) (string, *Color, *Decoration)
}

// maxColumns returns the greatest column count among rows.
func maxColumns(rows [][]string) int {
	columnCount := 0
	for _, row := range rows {
		columnCount = max(columnCount, len(row))
	}
	return columnCount
}
