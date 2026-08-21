package text

import "github.com/nekrassov01/table/internal/caption"

// CaptionSide specifies the visual position of a table caption.
type CaptionSide = caption.Side

const (
	// CaptionDefault renders the caption below the table.
	CaptionDefault = caption.Default

	// CaptionTop renders the caption above the table.
	CaptionTop = caption.Top

	// CaptionBottom renders the caption below the table.
	CaptionBottom = caption.Bottom
)
