package table

import "iter"

// Tabular writes all body rows in one call.
type Tabular interface {
	// Render writes all body rows.
	Render(rows [][]any) error
}

// Streamer writes body rows one at a time and completes the output with Close.
type Streamer interface {
	// Render writes one body row.
	Render(row []any) error

	// Close writes any deferred output and completes the stream.
	Close() error
}

// TableOf converts values to rows by calling fn once for each value, in order.
func TableOf[T any](values []T, fn func(T) []any) [][]any {
	rows := make([][]any, len(values))
	for i, value := range values {
		rows[i] = fn(value)
	}
	return rows
}

// StreamOf converts values to rows by calling fn for each successful value.
// It stops after forwarding the first error from values.
func StreamOf[T any](values iter.Seq2[T, error], fn func(T) []any) iter.Seq2[[]any, error] {
	return func(yield func([]any, error) bool) {
		for value, err := range values {
			if err != nil {
				_ = yield(nil, err)
				return
			}
			if !yield(fn(value), nil) {
				return
			}
		}
	}
}
