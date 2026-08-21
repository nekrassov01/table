package html

import "github.com/nekrassov01/table/internal/caption"

// CaptionSide specifies the visual position of a table caption.
type CaptionSide = caption.Side

const (
	// CaptionDefault emits no caption-side declaration and leaves placement to
	// the user agent.
	CaptionDefault = caption.Default

	// CaptionTop emits caption-side:top.
	CaptionTop = caption.Top

	// CaptionBottom emits caption-side:bottom.
	CaptionBottom = caption.Bottom
)

// resolveCaptionSide returns the CSS declaration for side, or an empty string
// for [CaptionDefault].
func resolveCaptionSide(side CaptionSide) string {
	switch side {
	case CaptionTop:
		return "caption-side:top"
	case CaptionBottom:
		return "caption-side:bottom"
	}
	return ""
}
