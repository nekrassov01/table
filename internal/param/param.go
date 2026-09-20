// Package param defines format-independent constants shared by format packages.
package param

const (
	// SpanLimit is the number of positions represented by a uint64 span mask.
	// Span options affect only the first SpanLimit columns; later columns render
	// independently. Tables without spans have no column-count limit.
	SpanLimit = 64
)
