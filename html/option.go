package html

import "github.com/nekrassov01/table/internal/column"

// Option configures a [Table] or [Stream] during construction. Column indexes
// refer to positions in the input rows; a generated index column does not
// change them.
//
// Options that set values replace earlier values for the same setting, column,
// and scope. Enable-only options never clear previously enabled targets:
// repeated calls are idempotent, and calls for different columns or scopes
// retain all selected targets.
type Option func(*option)

// WithHeader sets header rows in top-to-bottom order. [WithRowspan] and
// [WithColspan] also apply to header labels when their scope includes
// [ScopeHeader].
func WithHeader(rows ...[]string) Option {
	return func(o *option) {
		o.header = rows
	}
}

// WithFooter sets a function that returns footer rows in top-to-bottom order.
// The constructor retains fn without calling it. [Table.Render] calls it once
// before processing body rows, and [Stream.Close] calls it once after the last
// body row. A nil function or nil result produces no footer.
//
// The deferred call allows the footer to read aggregates collected by the
// caller. For a Table, fn runs before body transformers. A footer row cannot
// contain more cells than the column count established by a header or earlier
// Stream output.
func WithFooter(fn func() [][]string) Option {
	return func(o *option) {
		o.footer = fn
	}
}

// WithCaption sets the caption text and its position. [CaptionDefault] emits no
// caption-side declaration. The caption element remains the table's first
// child regardless of its visual position.
func WithCaption(s string, side CaptionSide) Option {
	return func(o *option) {
		o.caption = s
		o.captionSide = side
	}
}

// WithTableAttr sets classes and inline styles for the table and its child
// elements.
func WithTableAttr(attr TableAttr) Option {
	// Escape the captured value once so applying a reusable option is read-only.
	attr = attr.escape()
	return func(o *option) {
		o.attrs = attr
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
// Empty header and footer labels remain empty.
func WithPlaceholder(s string) Option {
	return func(o *option) {
		o.placeholder = s
	}
}

// WithAlign sets the alignment of the selected columns in the table parts
// selected by scopes. [AlignDefault] restores the default for each part.
func WithAlign(scopes Scope, columns ColumnSelector, align AlignSide) Option {
	return func(o *option) {
		o.columns.apply(columns, func(c *columnConfig) {
			c.aligns.Set(scopes, align)
		})
	}
}

// WithRowspan merges consecutive equal values vertically in the selected
// columns and table parts. A change in a selected column also ends spans in
// selected columns to its right. Spans never cross a part boundary and are
// limited to the first 64 rendered columns. Table emits a rowspan attribute;
// Stream emits an empty cell for each continuation.
func WithRowspan(scopes Scope, columns ColumnSelector) Option {
	return func(o *option) {
		o.columns.apply(columns, func(c *columnConfig) {
			c.rowspan |= scopes
		})
	}
}

// WithColspan merges adjacent equal values horizontally in the selected
// columns and table parts by emitting a colspan attribute. Spanning is limited
// to the first 64 rendered columns.
func WithColspan(scopes Scope, columns ColumnSelector) Option {
	return func(o *option) {
		o.columns.apply(columns, func(c *columnConfig) {
			c.colspan |= scopes
		})
	}
}

// WithColor applies color to values in the selected columns and table parts.
// The generated span is nested inside the cell and composes with
// [WithCellAttr].
func WithColor(scopes Scope, columns ColumnSelector, color *Color) Option {
	return func(o *option) {
		o.columns.apply(columns, func(c *columnConfig) {
			c.transformer.colors.Set(scopes, color)
		})
	}
}

// WithDecoration applies decoration to values in the selected columns and
// table parts.
func WithDecoration(scopes Scope, columns ColumnSelector, decoration *Decoration) Option {
	return func(o *option) {
		o.columns.apply(columns, func(c *columnConfig) {
			c.transformer.decorations.Set(scopes, decoration)
		})
	}
}

// WithCellAttr sets cell classes and inline styles in the selected columns and
// table parts. The class is joined with the corresponding [TableAttr] cell
// class. The style is joined with its cell style and alignment declaration.
func WithCellAttr(scopes Scope, columns ColumnSelector, attr Attr) Option {
	// Escape the captured value once so applying a reusable option is read-only.
	attr = attr.escape()
	return func(o *option) {
		o.columns.apply(columns, func(c *columnConfig) {
			c.attrs.Set(scopes, attr)
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
