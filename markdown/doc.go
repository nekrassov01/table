// Package markdown renders tabular data as GitHub Flavored Markdown tables
// written to an io.Writer.
//
// [NewTable] creates a [Table]. [Table.Render] accepts all body rows at once as
// [][]any. [NewStream] creates a [Stream]. [Stream.Render] accepts one []any
// body row at a time, and [Stream.Close] must be called after the final row.
// Both types are configured with options passed to their constructors.
//
// GitHub Flavored Markdown requires exactly one header row followed
// immediately by a delimiter row, so callers must configure a header and
// cannot configure a footer. Options control alignment markers, indexing,
// value transformation, repeated-value suppression, colors, and text
// decorations.
//
// Column-targeting options accept a [ColumnSelector]. [Columns] selects input
// positions explicitly, while [AllColumns] includes columns discovered after
// options are applied. A generated index column is never selected by either.
//
// Markdown has no merged-cell representation. Row- and column-spanning
// options therefore render absorbed body cells as blank cells, and do not
// alter header cells. Ordinary cell content is escaped so Markdown and HTML
// delimiters render as literal text, though GFM may still autolink URL-shaped
// text. Its whitespace follows GFM and HTML trimming and collapsing rules; use
// [DecorationPreformatted] when spacing must be preserved. Colors are emitted
// as inline HTML spans. A decoration normally wraps the color span; for inline
// code, the color span wraps the code-span fence because HTML inside the fence
// would render as code.
package markdown
