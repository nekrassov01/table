// Package csv renders tabular data as delimiter-separated records written to
// an io.Writer.
//
// [NewTable] creates a [Table]. [Table.Render] accepts all body rows at once as
// [][]any. [NewStream] creates a [Stream]. [Stream.Render] accepts one []any
// body row at a time, and [Stream.Close] must be called after the final row.
// Both types are configured with options passed to their constructors.
//
// The default delimiter is a tab. [WithDelimiter] accepts the same delimiter
// runes as [encoding/csv.Writer]. Fields are quoted using the conventions of
// encoding/csv.Writer. Records end in LF by default. [WithCRLF] selects CRLF
// record endings and applies the same line break normalization as
// encoding/csv.Writer with UseCRLF enabled.
//
// A header contains exactly one record. Footer rows are emitted as ordinary
// records because delimiter-separated formats have no distinct footer
// section. Options can also add an index and transform values. The package
// does not apply visual styling or merge fields.
//
// [encoding/csv.Writer]: https://pkg.go.dev/encoding/csv#Writer
package csv
