package text

import "strconv"

// Code identifies one ANSI Select Graphic Rendition parameter passed to
// [NewAttr].
type Code int

// SGR codes are defined in https://en.wikipedia.org/wiki/ANSI_escape_code#SGR.
const (
	CodeReset        Code = 0 // CodeReset resets all attributes.
	CodeBold         Code = 1 // CodeBold enables bold text.
	CodeFaint        Code = 2 // CodeFaint enables faint text.
	CodeItalic       Code = 3 // CodeItalic enables italic text.
	CodeUnderline    Code = 4 // CodeUnderline enables underlined text.
	CodeBlinkSlow    Code = 5 // CodeBlinkSlow enables slow blinking.
	CodeBlinkRapid   Code = 6 // CodeBlinkRapid enables rapid blinking.
	CodeReverseVideo Code = 7 // CodeReverseVideo swaps foreground and background.
	CodeConcealed    Code = 8 // CodeConcealed hides text.
	CodeCrossedOut   Code = 9 // CodeCrossedOut enables crossed-out text.

	CodeFgBlack   Code = 30 // CodeFgBlack selects a black foreground.
	CodeFgRed     Code = 31 // CodeFgRed selects a red foreground.
	CodeFgGreen   Code = 32 // CodeFgGreen selects a green foreground.
	CodeFgYellow  Code = 33 // CodeFgYellow selects a yellow foreground.
	CodeFgBlue    Code = 34 // CodeFgBlue selects a blue foreground.
	CodeFgMagenta Code = 35 // CodeFgMagenta selects a magenta foreground.
	CodeFgCyan    Code = 36 // CodeFgCyan selects a cyan foreground.
	CodeFgWhite   Code = 37 // CodeFgWhite selects a white foreground.

	CodeFgHiBlack   Code = 90 // CodeFgHiBlack selects a bright black foreground.
	CodeFgHiRed     Code = 91 // CodeFgHiRed selects a bright red foreground.
	CodeFgHiGreen   Code = 92 // CodeFgHiGreen selects a bright green foreground.
	CodeFgHiYellow  Code = 93 // CodeFgHiYellow selects a bright yellow foreground.
	CodeFgHiBlue    Code = 94 // CodeFgHiBlue selects a bright blue foreground.
	CodeFgHiMagenta Code = 95 // CodeFgHiMagenta selects a bright magenta foreground.
	CodeFgHiCyan    Code = 96 // CodeFgHiCyan selects a bright cyan foreground.
	CodeFgHiWhite   Code = 97 // CodeFgHiWhite selects a bright white foreground.

	CodeBgBlack   Code = 40 // CodeBgBlack selects a black background.
	CodeBgRed     Code = 41 // CodeBgRed selects a red background.
	CodeBgGreen   Code = 42 // CodeBgGreen selects a green background.
	CodeBgYellow  Code = 43 // CodeBgYellow selects a yellow background.
	CodeBgBlue    Code = 44 // CodeBgBlue selects a blue background.
	CodeBgMagenta Code = 45 // CodeBgMagenta selects a magenta background.
	CodeBgCyan    Code = 46 // CodeBgCyan selects a cyan background.
	CodeBgWhite   Code = 47 // CodeBgWhite selects a white background.

	CodeBgHiBlack   Code = 100 // CodeBgHiBlack selects a bright black background.
	CodeBgHiRed     Code = 101 // CodeBgHiRed selects a bright red background.
	CodeBgHiGreen   Code = 102 // CodeBgHiGreen selects a bright green background.
	CodeBgHiYellow  Code = 103 // CodeBgHiYellow selects a bright yellow background.
	CodeBgHiBlue    Code = 104 // CodeBgHiBlue selects a bright blue background.
	CodeBgHiMagenta Code = 105 // CodeBgHiMagenta selects a bright magenta background.
	CodeBgHiCyan    Code = 106 // CodeBgHiCyan selects a bright cyan background.
	CodeBgHiWhite   Code = 107 // CodeBgHiWhite selects a bright white background.
)

// Attr holds precomputed ANSI SGR sequences written around a cell value.
type Attr struct {
	Prefix []byte // SGR sequence written before a value.
	Suffix []byte // SGR sequence written after a value.
}

// Color presets provide common foreground and background attributes.
var (
	// ColorFgBlack sets a black foreground.
	ColorFgBlack = NewAttr(CodeFgBlack)

	// ColorFgRed sets a red foreground.
	ColorFgRed = NewAttr(CodeFgRed)

	// ColorFgGreen sets a green foreground.
	ColorFgGreen = NewAttr(CodeFgGreen)

	// ColorFgYellow sets a yellow foreground.
	ColorFgYellow = NewAttr(CodeFgYellow)

	// ColorFgBlue sets a blue foreground.
	ColorFgBlue = NewAttr(CodeFgBlue)

	// ColorFgMagenta sets a magenta foreground.
	ColorFgMagenta = NewAttr(CodeFgMagenta)

	// ColorFgCyan sets a cyan foreground.
	ColorFgCyan = NewAttr(CodeFgCyan)

	// ColorFgWhite sets a white foreground.
	ColorFgWhite = NewAttr(CodeFgWhite)

	// ColorBgBlack sets a black background.
	ColorBgBlack = NewAttr(CodeBgBlack)

	// ColorBgRed sets a red background.
	ColorBgRed = NewAttr(CodeBgRed)

	// ColorBgGreen sets a green background.
	ColorBgGreen = NewAttr(CodeBgGreen)

	// ColorBgYellow sets a yellow background.
	ColorBgYellow = NewAttr(CodeBgYellow)

	// ColorBgBlue sets a blue background.
	ColorBgBlue = NewAttr(CodeBgBlue)

	// ColorBgMagenta sets a magenta background.
	ColorBgMagenta = NewAttr(CodeBgMagenta)

	// ColorBgCyan sets a cyan background.
	ColorBgCyan = NewAttr(CodeBgCyan)

	// ColorBgWhite sets a white background.
	ColorBgWhite = NewAttr(CodeBgWhite)
)

// NewAttr returns an Attr for codes, or nil when no codes are supplied.
func NewAttr(codes ...Code) *Attr {
	if len(codes) == 0 {
		return nil
	}
	return &Attr{
		Prefix: buildSGR(codes),
		Suffix: buildSGR([]Code{CodeReset}),
	}
}

// isZero reports whether the receiver is nil or has no SGR prefix.
func (o *Attr) isZero() bool {
	return o == nil || len(o.Prefix) == 0
}

// size returns the byte length added when the attribute is written.
func (o *Attr) size() int {
	if o.isZero() {
		return 0
	}
	return len(o.Prefix) + len(o.Suffix)
}

// buildSGR serializes the codes as a CSI SGR sequence.
func buildSGR(codes []Code) []byte {
	if len(codes) == 0 {
		return nil
	}
	b := make([]byte, 0, 2+len(codes)*4+1)
	b = append(b, '\x1b', '[')
	for i, code := range codes {
		if i > 0 {
			b = append(b, ';')
		}
		b = strconv.AppendInt(b, int64(code), 10)
	}
	b = append(b, 'm')
	return b
}
