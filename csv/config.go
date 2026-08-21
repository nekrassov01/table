package csv

import "unicode/utf8"

// DefaultPlaceholder is the default placeholder for missing or empty body
// cells.
const DefaultPlaceholder = ""

// config validates construction options and resolves pass data into logical
// columns.
type config struct {
	bodyColumns int          // Body width used to resolve logical columns.
	state       *configState // Reusable config state.
	err         error        // Invalid delimiter configuration.
	output      configResult // Pass data paired with resolved columns.
}

// prepare validates the delimiter and resolves logical columns. A non-empty
// header fixes their count; otherwise the current body and footer determine it.
func (o *config) prepare() {
	result := &o.output
	option := result.option
	if !validDelimiter(option.delimiter) {
		o.err = newDelimiterError(option.delimiter)
		return
	}
	headerColumns := len(result.header)
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
	header        []string   // Header fields in input order.
	footer        [][]string // Footer rows in top-to-bottom order.
	bodyRows      int        // Number of body rows in this pass.
	columns       []column   // Resolved column settings in logical order.
	footerColumns int        // Footer width resolved for the current config pass.
}

// option holds settings fixed when a Table or Stream is constructed.
type option struct {
	placeholder string            // Text for a missing value.
	header      []string          // The optional header row.
	footer      func() [][]string // Generates footer rows for Render or Close.
	columns     columnSet         // Input columns and their defaults.
	delimiter   rune              // Field delimiter.
	crlf        bool              // Whether records use CRLF line endings.
	indexOffset int               // Synthetic leading column count: 0 or 1.
}

// apply sets defaults and applies opts in order.
func (o *option) apply(opts ...Option) {
	o.delimiter = '\t'
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
		for index := range o.values {
			fn(&o.values[index])
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
		for index := existing; index < columnCount; index++ {
			o.values[index] = defaults
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
	transformer func(any) string // Optional body-value transformation.
}

// maxColumns returns the greatest column count among rows.
func maxColumns(rows [][]string) int {
	columnCount := 0
	for _, row := range rows {
		columnCount = max(columnCount, len(row))
	}
	return columnCount
}

// validDelimiter reports whether delimiter can separate fields.
func validDelimiter(delimiter rune) bool {
	return delimiter != 0 &&
		delimiter != '"' &&
		delimiter != '\r' &&
		delimiter != '\n' &&
		utf8.ValidRune(delimiter) &&
		delimiter != utf8.RuneError
}
