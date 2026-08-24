package markdown

import "github.com/nekrassov01/table/internal/column"

// Option configures a [Table] or [Stream] during construction. Column indexes
// refer to positions in input rows; a generated index column does not change
// them.
//
// Options that set values replace earlier values for the same setting, column,
// and scope. Enable-only options never clear previously enabled targets:
// repeated calls are idempotent, and calls for different columns or scopes
// retain all selected targets.
type Option func(*option)

// WithHeader sets the required header row. GFM tables have exactly one header
// row followed immediately by a delimiter row.
func WithHeader(header []string) Option {
	return func(o *option) {
		o.header = header
	}
}

// WithIndex prepends a column that numbers body rows from 1. Column indexes in
// other options continue to refer to input positions.
func WithIndex() Option {
	return func(o *option) {
		o.indexOffset = 1
	}
}

// WithPlaceholder sets the value displayed for a missing or empty body cell.
// Empty header labels remain empty.
func WithPlaceholder(s string) Option {
	return func(o *option) {
		o.placeholder = s
	}
}

// WithAlign sets the alignment of the selected columns. GFM stores one
// alignment in the delimiter row, so it applies to both header and body cells.
func WithAlign(columns ColumnSelector, align AlignSide) Option {
	return func(o *option) {
		o.columns.apply(columns, func(c *columnConfig) {
			c.align = align
		})
	}
}

// WithRowspan blanks consecutive equal body values vertically in the selected
// columns. A change in a selected column also ends spans in selected columns
// to its right. Markdown has no merged-cell representation. Spanning is
// limited to the first 64 rendered columns.
func WithRowspan(columns ColumnSelector) Option {
	return func(o *option) {
		o.columns.apply(columns, func(c *columnConfig) {
			c.rowspan = true
		})
	}
}

// WithColspan blanks trailing equal body values horizontally in the selected
// columns. Markdown has no merged-cell representation. Spanning is limited to
// the first 64 rendered columns.
func WithColspan(columns ColumnSelector) Option {
	return func(o *option) {
		o.columns.apply(columns, func(c *columnConfig) {
			c.colspan = true
		})
	}
}

// WithColor applies color to values in the selected columns and table parts.
func WithColor(scopes Scope, columns ColumnSelector, color *Color) Option {
	return func(o *option) {
		o.columns.apply(columns, func(c *columnConfig) {
			c.transformer.colors.Set(scopes, color)
		})
	}
}

// WithDecoration applies decoration to values in the selected columns and
// table parts. DecorationCode uses a content-sized backtick fence, and
// DecorationPreformatted preserves whitespace through raw HTML.
func WithDecoration(scopes Scope, columns ColumnSelector, decoration *Decoration) Option {
	return func(o *option) {
		o.columns.apply(columns, func(c *columnConfig) {
			c.transformer.decorations.Set(scopes, decoration)
		})
	}
}

// WithTransformer sets a transformer for body cells in the selected columns.
// The function receives the raw value and may return a replacement display
// string, color, and decoration. A non-empty string skips formatting the raw
// value; an empty string uses the formatted raw value. Nil color or decoration
// values keep the corresponding column settings.
func WithTransformer(columns ColumnSelector, fn func(any) (string, *Color, *Decoration)) Option {
	return func(o *option) {
		o.columns.apply(columns, func(c *columnConfig) {
			c.transformer.fn = fn
		})
	}
}

// ColumnSelector identifies input columns for a column option.
type ColumnSelector struct {
	selector column.Selector
}

// Columns selects the input columns at indexes. Negative indexes are ignored.
func Columns(indexes ...int) ColumnSelector {
	return ColumnSelector{
		selector: column.NewSelector(indexes...),
	}
}

// AllColumns selects every input column, including columns discovered after
// options are applied. A generated index column is excluded.
func AllColumns() ColumnSelector {
	return ColumnSelector{
		selector: column.All(),
	}
}
