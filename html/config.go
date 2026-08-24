package html

import (
	"github.com/nekrassov01/table/internal/column"
	"github.com/nekrassov01/table/internal/scope"
)

// placeholder is the default text for missing or empty body cells.
const placeholder = ""

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
	headerColumns := column.MaxColumns(result.header)
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
	header        [][]string     // Header rows in top-to-bottom order.
	footer        [][]string     // Footer rows in top-to-bottom order.
	bodyRows      int            // Number of body rows in this pass.
	columns       []columnConfig // Resolved column settings in logical order.
	footerColumns int            // Footer width resolved for the current config pass.
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
	o.placeholder = placeholder
	for _, opt := range opts {
		opt(o)
	}
}

// columnSet applies shared column selection to HTML column settings.
type columnSet column.Set[columnConfig]

// apply updates every input column selected by selector.
func (o *columnSet) apply(selector ColumnSelector, fn func(*columnConfig)) {
	(*column.Set[columnConfig])(o).Apply(selector.selector, nil, fn)
}

// columnConfig holds HTML settings for one logical column.
type columnConfig struct {
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
