package markdown

import (
	"github.com/nekrassov01/table/internal/column"
	"github.com/nekrassov01/table/internal/scope"
)

// placeholder is the default text for missing or empty body cells.
const placeholder = " "

// config resolves construction options and the required header into logical
// columns.
type config struct {
	state  *configState // Reusable config state.
	err    error        // Invalid Markdown table configuration.
	output configResult // Pass data paired with resolved columns.
}

// prepare validates the header and resolves logical columns.
func (o *config) prepare() {
	result := &o.output
	option := result.option
	headerColumns := len(result.header)
	if headerColumns == 0 {
		o.err = newHeaderError()
		return
	}
	columnCount := headerColumns + option.indexOffset
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
}

// configResult carries pass data and resolved logical columns.
type configResult struct {
	option   *option        // Options fixed at construction.
	header   []string       // The required header row.
	bodyRows int            // Number of body rows in this pass.
	columns  []columnConfig // Resolved column settings in logical order.
}

// option holds settings fixed when a Table or Stream is constructed.
type option struct {
	placeholder string    // Text for a missing value.
	header      []string  // The required header row.
	columns     columnSet // Input columns and their defaults.
	indexOffset int       // Synthetic leading column count: 0 or 1.
}

// apply sets defaults and applies opts in order.
func (o *option) apply(opts ...Option) {
	o.placeholder = placeholder
	for _, opt := range opts {
		opt(o)
	}
}

// columnSet applies shared column selection to Markdown column settings.
type columnSet column.Set[columnConfig]

// apply updates every input column selected by selector.
func (o *columnSet) apply(selector ColumnSelector, fn func(*columnConfig)) {
	(*column.Set[columnConfig])(o).Apply(selector.selector, nil, fn)
}

// columnConfig holds Markdown settings for one logical column.
type columnConfig struct {
	transformer transformer // Value transformation and inner markup.
	align       AlignSide   // Alignment encoded by the delimiter row.
	rowspan     bool        // Whether equal body values span vertically.
	colspan     bool        // Whether equal body values span horizontally.
}

// transformer holds body-value transformation and inner markup for a column.
type transformer struct {
	colors      scope.Scopes[*Color]      // Color by table part.
	decorations scope.Scopes[*Decoration] // Decoration by table part.
	fn          func(any) (string, *Color, *Decoration)
}
