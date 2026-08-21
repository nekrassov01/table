// Package html formats tabular data as semantic HTML table markup and writes it
// to an io.Writer.
//
// [NewTable] creates a [Table]. [Table.Render] accepts all body rows at once as
// [][]any. [NewStream] creates a [Stream]. [Stream.Render] accepts one []any
// body row at a time, and [Stream.Close] must be called after the final row.
// Both types are configured with options passed to their constructors.
//
// Output uses table, caption, thead, tbody, and tfoot elements as configured.
// Options control alignment, multi-row headers and footers, row and column
// spans, colors, text decorations, and the class and inline style attached to
// each element. Text and attribute values are escaped for their respective
// HTML contexts. When color and decoration both apply, decoration wraps the
// color span so block decorations such as pre remain valid HTML.
//
// Table can inspect the complete input and emits rowspan attributes for
// repeated values. Stream cannot know a row span's final length while writing,
// so it emits an empty cell for each continuation instead. [Table.Render]
// resolves a dynamic footer before processing body rows; [Stream.Close]
// resolves it after the last body row. Row- and column-spanning options affect
// only the first 64 rendered columns. Table does not split vertical runs at
// HTML's 65,534-row rowspan limit, so a longer run may emit a non-conforming
// rowspan value.
package html
