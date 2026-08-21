// Package text formats tabular data as bordered text tables and writes them to
// an io.Writer.
//
// [NewTable] creates a [Table]. [Table.Render] accepts all body rows at once as
// [][]any. [NewStream] creates a [Stream]. [Stream.Render] accepts one []any
// body row at a time, and [Stream.Close] must be called after the final row.
// Both types are configured with options passed to their constructors.
//
// The package supports Unicode and ASCII border styles, captions, multi-row
// headers and footers, row and column spans, alignment, padding, wrapping,
// truncation, and fitting to a terminal width. [Attr] represents both ANSI
// colors and text decorations because SGR encodes them in the same parameter
// sequence. ANSI attributes are omitted when the destination is not a terminal.
//
// Table measures the complete input before writing, so every row and a
// dynamically generated footer can affect column geometry. Stream writes
// incrementally and fixes its geometry when output begins, normally from the
// configured header and first body row. Later values wrap or truncate within
// that geometry. A stream resolves its footer only when [Stream.Close] is
// called. If no body row was written, the header and footer determine the
// geometry together.
package text
