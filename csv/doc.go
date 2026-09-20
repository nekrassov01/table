// Package csv renders tabular data as delimiter-separated records written to
// an io.Writer.
//
// [NewTable] creates a [Table]. [Table.Render] accepts all body rows at once as
// [][]table.Value.
// [NewStream] creates a [Stream]. [Stream.Render] accepts one []table.Value
// body row at a time, and [Stream.Close] must be called after the final row.
// Both types are configured with options passed to their constructors.
//
// The default delimiter is a tab. [WithDelimiter] accepts the same delimiter
// runes as [encoding/csv.Writer]. Fields are quoted using the conventions of
// encoding/csv.Writer. Records end in LF by default. [WithCRLF] selects CRLF
// record endings and applies the same line break normalization as
// encoding/csv.Writer with UseCRLF enabled.
//
// The package does not neutralize spreadsheet formulas. Spreadsheet software
// may interpret fields beginning with =, +, -, @, tab, or carriage return even
// when CSV quoting is applied. Sanitize untrusted values before rendering when
// the output will be opened in a spreadsheet. A one-column record containing an
// empty field is written as a blank line, matching encoding/csv.Writer. Readers
// that skip blank lines do not preserve that record; use a non-empty placeholder
// or transformer when it must survive a round trip.
//
// A header contains exactly one record. Footer rows are emitted as ordinary
// records because delimiter-separated formats have no distinct footer
// section. Options can transform values. The package does not apply visual
// styling or merge fields.
//
// [encoding/csv.Writer]: https://pkg.go.dev/encoding/csv#Writer
package csv
