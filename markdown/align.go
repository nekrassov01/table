package markdown

import "github.com/nekrassov01/table/internal/align"

// AlignSide controls per-column horizontal text alignment.
type AlignSide = align.Side

const (
	// AlignDefault emits no alignment marker, so GFM uses its default left
	// alignment.
	AlignDefault = align.Default

	// AlignLeft left-aligns all values in the column.
	AlignLeft = align.Left

	// AlignRight right-aligns all values in the column.
	AlignRight = align.Right

	// AlignCenter center-aligns all values in the column.
	AlignCenter = align.Center
)
