// Package param defines format-independent constants shared by format packages.
package param

const (
	// IndexHeader labels the synthetic index column.
	IndexHeader = "#"

	// SpanLimit is the number of positions represented by a uint64 span mask.
	// Span options affect only the first SpanLimit columns; later columns render
	// independently. Tables without spans have no column-count limit.
	SpanLimit = 64
)
