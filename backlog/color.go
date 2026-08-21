package backlog

import (
	"strings"

	"github.com/nekrassov01/table/internal/color"
)

// unsafeChars cannot be represented safely in Backlog color arguments.
const unsafeChars = "),{}|\r\n"

// Color holds Backlog markup that applies color to a cell value.
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

	// ColorBgWhite sets a white background.
	ColorBgWhite = NewColor("", "white")
)

// NewColor returns a Color that wraps values in Backlog color notation using
// fg and bg. An empty argument omits that color; if both are empty, NewColor
// returns nil. It also returns nil when either argument contains syntax that
// cannot be escaped safely.
func NewColor(fg, bg string) *Color {
	if fg == "" && bg == "" {
		return nil
	}
	if strings.ContainsAny(fg, unsafeChars) || strings.ContainsAny(bg, unsafeChars) {
		return nil
	}
	size := len("&color(") + len("){")
	if fg != "" {
		size += len(fg)
	} else {
		size += len(`""`)
	}
	if bg != "" {
		size += len(", ") + len(bg)
	}
	var b strings.Builder
	b.Grow(size)
	b.WriteString("&color(")
	if fg != "" {
		b.WriteString(fg)
	} else {
		b.WriteString(`""`)
	}
	if bg != "" {
		b.WriteString(", ")
		b.WriteString(bg)
	}
	b.WriteString("){")
	return color.New(b.String(), "}")
}
