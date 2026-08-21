// Package backlog renders tabular data using [Backlog table notation] and
// writes the result to an io.Writer.
//
// [NewTable] creates a [Table]. [Table.Render] accepts all body rows at once as
// [][]any. [NewStream] creates a [Stream]. [Stream.Render] accepts one []any
// body row at a time, and [Stream.Close] must be called after the final row.
// Both types are configured with options passed to their constructors.
//
// Backlog notation marks a header cell by placing '~' immediately after its
// opening '|'. The package uses those cells for configured header and footer
// bands, so either band may contain multiple rows without a separate delimiter.
// A footer is an API-level convention; its output uses ordinary Backlog
// header-cell notation. Options can also add an index, transform values,
// suppress repeated values across rows or columns, and apply Backlog color and
// decoration markup.
//
// Values are escaped where they would otherwise change the table notation.
// Color markup wraps ordinary decoration markup. When code decoration and
// color both apply, the code decoration is emitted and the color is omitted
// because Backlog notation has no representation for that combination.
//
// [Backlog table notation]: https://help-center.backlog.com/%E3%83%86%E3%82%AD%E3%82%B9%E3%83%88%E6%95%B4%E5%BD%A2%E3%83%AB%E3%83%BC%E3%83%AB%EF%BC%88Backlog%E8%A8%98%E6%B3%95%EF%BC%89-6a1d4d7f3abb3ada78c5658b#index7
package backlog
