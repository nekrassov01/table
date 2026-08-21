package html

import (
	"strings"

	"github.com/nekrassov01/table/internal/color"
)

// Color holds HTML markup that applies color to a cell value.
type Color = color.Color

// Color presets provide common foreground and background colors.
var (
	// ColorFgBlack sets a black foreground.
	ColorFgBlack = NewColor("black", "")

	// ColorFgRed sets a red foreground.
	ColorFgRed = NewColor("red", "")

	// ColorFgGreen sets a green foreground.
	ColorFgGreen = NewColor("green", "")

	// ColorFgYellow sets a yellow foreground.
	ColorFgYellow = NewColor("yellow", "")

	// ColorFgBlue sets a blue foreground.
	ColorFgBlue = NewColor("blue", "")

	// ColorFgMagenta sets a magenta foreground.
	ColorFgMagenta = NewColor("magenta", "")

	// ColorFgCyan sets a cyan foreground.
	ColorFgCyan = NewColor("cyan", "")

	// ColorFgWhite sets a white foreground.
	ColorFgWhite = NewColor("white", "")

	// ColorBgBlack sets a black background.
	ColorBgBlack = NewColor("", "black")

	// ColorBgRed sets a red background.
	ColorBgRed = NewColor("", "red")

	// ColorBgGreen sets a green background.
	ColorBgGreen = NewColor("", "green")

	// ColorBgYellow sets a yellow background.
	ColorBgYellow = NewColor("", "yellow")

	// ColorBgBlue sets a blue background.
	ColorBgBlue = NewColor("", "blue")

	// ColorBgMagenta sets a magenta background.
	ColorBgMagenta = NewColor("", "magenta")

	// ColorBgCyan sets a cyan background.
	ColorBgCyan = NewColor("", "cyan")

	// ColorBgWhite sets a white background.
	ColorBgWhite = NewColor("", "white")
)

// NewColor returns a Color that wraps values in a span using fg and bg as CSS
// color values. An empty value omits that property; if both are empty,
// NewColor returns nil. Values are normalized and escaped for an HTML
// attribute but are not otherwise validated as CSS.
func NewColor(fg, bg string) *Color {
	if fg == "" && bg == "" {
		return nil
	}
	fg = escapeAttr(fg)
	bg = escapeAttr(bg)
	size := len(`<span style="`) + len(`">`)
	if fg != "" {
		size += len("color:") + len(fg)
	}
	if bg != "" {
		size += len("background-color:") + len(bg)
	}
	if fg != "" && bg != "" {
		size += len(";")
	}
	var b strings.Builder
	b.Grow(size)
	b.WriteString(`<span style="`)
	if fg != "" {
		b.WriteString("color:")
		b.WriteString(fg)
	}
	if fg != "" && bg != "" {
		b.WriteString(";")
	}
	if bg != "" {
		b.WriteString("background-color:")
		b.WriteString(bg)
	}
	b.WriteString(`">`)
	return color.New(b.String(), "</span>")
}
