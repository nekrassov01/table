package text

import "github.com/nekrassov01/table/internal/align"

// AlignSide specifies horizontal cell alignment.
type AlignSide = align.Side

const (
	// AlignDefault centers header values and left-aligns body and footer values.
	AlignDefault = align.Default

	// AlignLeft left-aligns cell content.
	AlignLeft = align.Left

	// AlignRight right-aligns cell content.
	AlignRight = align.Right

	// AlignCenter centers cell content.
	AlignCenter = align.Center
)
