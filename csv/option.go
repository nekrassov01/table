package csv

import "slices"

// Option configures a [Table] or [Stream] during construction. Column indexes
// refer to positions in the input rows; a generated index column does not
// change them.
//
// Options that set values replace earlier values for the same setting or
// column. WithIndex and WithCRLF only enable their features; repeated calls are
// idempotent, and no option disables them.
type Option func(*option)

// WithHeader sets the optional header row. Delimiter-separated formats cannot
// distinguish additional header rows from body records.
func WithHeader(header []string) Option {
	return func(o *option) {
		o.header = header
	}
}

// WithFooter sets a function that returns footer rows in top-to-bottom order.
// The constructor retains fn without calling it. [Table.Render] calls it once
// before processing body rows, and [Stream.Close] calls it once after the last
// body row. A nil function or nil result produces no footer.
//
// Delimiter-separated formats emit footer rows as ordinary records. The
// deferred call allows those records to contain aggregates collected by the
// caller. For a Table, fn runs before body transformers. A footer row cannot
// contain more fields than the column count established by a header or earlier
// Stream output.
func WithFooter(fn func() [][]string) Option {
	return func(o *option) {
		o.footer = fn
	}
}

// WithDelimiter sets the field delimiter. It defaults to a tab. NUL, double
// quote, CR, LF, invalid runes, and U+FFFD are rejected.
func WithDelimiter(delimiter rune) Option {
	return func(o *option) {
		o.delimiter = delimiter
	}
}

// WithCRLF terminates records with CRLF. Within quoted fields, line feeds are
// converted to CRLF and carriage returns are removed, matching
// encoding/csv.Writer with UseCRLF enabled.
func WithCRLF() Option {
	return func(o *option) {
		o.crlf = true
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

// WithTransformer sets a transformer for body cells in the selected columns.
// The function receives the raw value and may return a replacement display
// string. A non-empty string skips formatting the raw value; an empty string
// uses the formatted raw value.
func WithTransformer(columns ColumnSelector, fn func(any) string) Option {
	return func(o *option) {
		o.columns.apply(columns, func(c *column) {
			c.transformer = fn
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
