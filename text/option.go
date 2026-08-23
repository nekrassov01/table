package text

import (
	"slices"
)

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
// caller. For a Table, fn runs before body transformers. For a Stream, column
// geometry is already fixed, so wider footer values wrap or truncate within
// the established columns. A footer row cannot contain more cells than the
// column count established by a header or earlier Stream output.
func WithFooter(fn func() [][]string) Option {
	return func(o *option) {
		o.footer = fn
	}
}

// WithCaption sets the caption text and its position. [CaptionDefault] places
// the caption below the table.
func WithCaption(s string, side CaptionSide) Option {
	return func(o *option) {
		o.caption = s
		o.captionSide = side
	}
}

// WithStyle sets the border and content style.
func WithStyle(style Style) Option {
	return func(o *option) {
		o.style = style
	}
}

// WithCompact omits horizontal separators between body rows. Row spans retain
// the separator segments at their boundaries.
func WithCompact() Option {
	return func(o *option) {
		o.compact = true
	}
}

// WithIndex prepends a column that numbers body rows from 1. Column indexes in
// other options continue to refer to input positions.
func WithIndex() Option {
	return func(o *option) {
		o.indexOffset = 1
	}
}

// WithIndexWidth prepends an index column as [WithIndex] does and gives it a
// minimum width of n digits. Stream reserves three digits by default; use this
// option when the expected row count needs more. A non-positive n enables the
// index without changing a minimum set by an earlier option.
func WithIndexWidth(n int) Option {
	return func(o *option) {
		o.indexOffset = 1
		if n > 0 {
			o.indexWidth = n
		}
	}
}

// WithAutoFit shrinks columns to the writer's terminal width and wraps content
// within the reduced widths. It has no effect for a non-terminal writer or
// when any column uses [WithWidth] or [WithTruncate].
func WithAutoFit() Option {
	return func(o *option) {
		o.autoFit = true
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
		o.columns.apply(columns, func(c *column) {
			c.aligns.Set(scopes, align)
		})
	}
}

// WithWidth sets the display-width boundary of the selected columns. Narrower
// content is padded; wider content wraps, or truncates where [WithTruncate]
// applies. Wrapping preserves grapheme clusters, so a cluster wider than the
// boundary remains intact. A non-positive width clears the setting.
func WithWidth(columns ColumnSelector, width int) Option {
	// Do not mutate values captured by a reusable option.
	limit := max(width, 0)
	return func(o *option) {
		o.columns.apply(columns, func(c *column) {
			c.limit = limit
		})
	}
}

// WithTruncate replaces content that exceeds the selected column width with an
// ellipsis instead of wrapping it. Without [WithWidth], it affects only later
// Stream values that exceed the geometry fixed when output began.
func WithTruncate(columns ColumnSelector) Option {
	return func(o *option) {
		o.columns.apply(columns, func(c *column) {
			c.truncate = true
		})
	}
}

// WithPadding sets the left and right padding of the selected columns.
// Negative values are treated as zero.
func WithPadding(columns ColumnSelector, left, right int) Option {
	// Do not mutate values captured by a reusable option.
	lPad, rPad := max(left, 0), max(right, 0)
	return func(o *option) {
		o.columns.apply(columns, func(c *column) {
			c.lPad = lPad
			c.rPad = rPad
		})
	}
}

// WithRowspan merges consecutive equal values vertically in the selected
// columns and table parts. A change in a selected column also ends spans in
// selected columns to its right. Spans never cross a part boundary and are
// limited to the first 64 rendered columns.
func WithRowspan(scopes Scope, columns ColumnSelector) Option {
	return func(o *option) {
		o.columns.apply(columns, func(c *column) {
			c.rowspan |= scopes
		})
	}
}

// WithColspan merges adjacent equal values horizontally in the selected
// columns and table parts. Spanning is limited to the first 64 rendered
// columns.
func WithColspan(scopes Scope, columns ColumnSelector) Option {
	return func(o *option) {
		o.columns.apply(columns, func(c *column) {
			c.colspan |= scopes
		})
	}
}

// WithAttr applies attr to values in the selected columns and table parts. It
// overrides the part attribute from [Style]; an attribute returned by
// [WithTransformer] overrides both.
func WithAttr(scopes Scope, columns ColumnSelector, attr *Attr) Option {
	return func(o *option) {
		o.columns.apply(columns, func(c *column) {
			c.transformer.attrs.Set(scopes, attr)
		})
	}
}

// WithTransformer sets a transformer for body cells in the selected columns.
// The function receives the raw value and may return a replacement display
// string and attribute. A non-empty string skips formatting the raw value; an
// empty string uses the formatted raw value. A nil attribute keeps the
// corresponding column setting.
func WithTransformer(columns ColumnSelector, fn func(any) (string, *Attr)) Option {
	return func(o *option) {
		o.columns.apply(columns, func(c *column) {
			c.transformer.fn = fn
		})
	}
}

// ColumnSelector identifies input columns for a column option.
type ColumnSelector struct {
	indexes []int
	all     bool
}

// Columns selects the input columns at indexes. Negative indexes are ignored.
func Columns(indexes ...int) ColumnSelector {
	return ColumnSelector{
		indexes: slices.Clone(indexes),
	}
}

// AllColumns selects every input column, including columns discovered after
// options are applied. A generated index column is excluded.
func AllColumns() ColumnSelector {
	return ColumnSelector{
		all: true,
	}
}
