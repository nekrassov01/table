package csv

import (
	"unicode/utf8"

	"github.com/nekrassov01/table/internal/column"
)

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
	footerColumns := column.MaxColumns(result.footer)
	columnCount := headerColumns
	if columnCount == 0 {
		columnCount = max(o.bodyColumns, footerColumns)
	}
	if columnCount > 0 {
		columnCount += option.indexOffset
	}
	columns := o.state.columns
	if cap(columns) < columnCount {
		columns = make([]columnConfig, columnCount)
	} else {
		columns = columns[:columnCount]
	}
	for index := 0; index < min(option.indexOffset, columnCount); index++ {
		columns[index] = columnConfig{}
	}
	defaults := columnConfig{}
	if option.columns.Defaults != nil {
		defaults = *option.columns.Defaults
	}
	for index := option.indexOffset; index < columnCount; index++ {
		columns[index] = defaults
	}
	if option.indexOffset < columnCount {
		copy(columns[option.indexOffset:], option.columns.Values)
	}
	o.state.columns = columns
	result.columns = columns
	result.footerColumns = footerColumns
}

// configResult carries pass data and resolved logical columns.
type configResult struct {
	option        *option        // Options fixed at construction.
	header        []string       // Header fields in input order.
	footer        [][]string     // Footer rows in top-to-bottom order.
	bodyRows      int            // Number of body rows in this pass.
	columns       []columnConfig // Resolved column settings in logical order.
	footerColumns int            // Footer width resolved for the current config pass.
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

// columnSet applies shared column selection to CSV column settings.
type columnSet column.Set[columnConfig]

// apply updates every input column selected by selector.
func (o *columnSet) apply(selector ColumnSelector, fn func(*columnConfig)) {
	(*column.Set[columnConfig])(o).Apply(selector.selector, nil, fn)
}

// columnConfig holds CSV settings for one logical column.
type columnConfig struct {
	transformer func(any) string // Optional body-value transformation.
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
