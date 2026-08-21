package markdown

import "github.com/nekrassov01/table/internal/scope"

// DefaultPlaceholder is the default placeholder for missing or empty body
// cells.
const DefaultPlaceholder = " "

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
}

// configResult carries pass data and resolved logical columns.
type configResult struct {
	option   *option  // Options fixed at construction.
	header   []string // The required header row.
	bodyRows int      // Number of body rows in this pass.
	columns  []column // Resolved column settings in logical order.
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
