package text

import (
	"io"

	"github.com/nekrassov01/table/internal/column"
	"github.com/nekrassov01/table/internal/scope"
)

// placeholder is the default text for missing or empty body cells.
const placeholder = " "

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
	columns := option.columns.resolve(o.state.columns, columnCount, option.indexOffset)
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
	o.placeholder = placeholder
	for _, opt := range opts {
		opt(o)
	}
	if o.indexOffset != 0 {
		o.indexWidth = max(o.indexWidth, minIndexWidth)
	}
	columns := (*column.Set[columnConfig])(&o.columns)
	o.plain = !isTerminal(w)
	if o.plain {
		o.style.Border.Attr = nil
		o.style.Content = ContentStyle{}
		if defaults := columns.Default(); defaults != nil {
			defaults.transformer.attrs = scope.Scopes[*Attr]{}
		}
		columns.Range(func(column *columnConfig) {
			column.transformer.attrs = scope.Scopes[*Attr]{}
		})
	}
	if !o.autoFit {
		return
	}
	if defaults := columns.Default(); defaults != nil && (defaults.limit > 0 || defaults.truncate) {
		o.autoFit = false
		return
	}
	if columns.Any(func(column *columnConfig) bool {
		return column.limit > 0 || column.truncate
	}) {
		o.autoFit = false
	}
}

// columnSet applies shared column selection to text column settings.
type columnSet column.Set[columnConfig]

// apply updates every input column selected by selector.
func (o *columnSet) apply(selector ColumnSelector, fn func(*columnConfig)) {
	(*column.Set[columnConfig])(o).Apply(selector.selector, defaultColumn, fn)
}

// resolve applies input settings to logical columns.
func (o *columnSet) resolve(columns []columnConfig, columnCount, indexOffset int) []columnConfig {
	return (*column.Set[columnConfig])(o).Resolve(columns, columnCount, indexOffset, defaultColumn())
}

// columnConfig holds text settings for one logical column.
type columnConfig struct {
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

// defaultColumn returns an unconfigured input column.
func defaultColumn() columnConfig {
	return columnConfig{
		lPad: 1,
		rPad: 1,
	}
}
