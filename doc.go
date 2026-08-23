// Package table provides common interfaces, row adapters, and errors for the
// format-specific table packages in this module.
//
// Rendering is implemented by the text, html, markdown, backlog, and csv
// subpackages. Applications choose one of those packages according to the
// output format they need. This package does not select a format or write
// output itself.
//
// Each format package provides a Table for batch input and a Stream for
// incremental input. A Table implements [Tabular], whose Render method accepts
// all body rows at once as [][]any. A Stream implements [Streamer], whose
// Render method accepts one []any body row at a time and whose Close method
// must be called after the final row.
//
// [TableOf] converts a typed slice to [][]any. [StreamOf] adapts an iterator
// to rows without first collecting the sequence in memory.
//
// Body values are converted consistently across format packages. Scalars,
// fmt.Stringer values, errors, and byte slices use their text form. Other
// values use the representation produced by fmt.Sprint. Nil values and empty
// strings, slices, and arrays use the configured placeholder. Transformer
// options can replace the default display text without evaluating that
// representation.
//
// Failures from format packages are returned as [Error] values that identify
// the package and unwrap the underlying cause. The sentinel errors in this
// package describe shared conditions and can be matched with errors.Is.
package table
