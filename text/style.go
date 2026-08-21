package text

// Style defines a table's border and content appearance.
type Style struct {
	Border  BorderStyle  // The frame around and between the cells.
	Content ContentStyle // The text within them.
}

// BorderStyle defines border glyphs and attributes. A nil horizontal or
// vertical field omits the corresponding border or separator.
type BorderStyle struct {
	Top      *Horizontal // The top border.
	Header   *Horizontal // The separator below the header.
	Body     *Horizontal // The separator between data rows.
	Footer   *Horizontal // The separator above the footer.
	Bottom   *Horizontal // The bottom border.
	Vertical *Vertical   // The vertical borders and column separators.
	Attr     *Attr       // Attribute applied to every border glyph.
}

// maxGlyphLen returns the longest byte length any line in the border can
// write. A nil line contributes zero.
func (o *BorderStyle) maxGlyphLen() int {
	return max(
		o.Vertical.maxGlyphLen(),
		o.Top.maxGlyphLen(),
		o.Header.maxGlyphLen(),
		o.Body.maxGlyphLen(),
		o.Footer.maxGlyphLen(),
		o.Bottom.maxGlyphLen(),
	)
}

// ContentStyle defines text attributes for each table part. A nil attribute
// leaves that part unstyled.
type ContentStyle struct {
	Header  *Attr // The header values.
	Body    *Attr // The data values.
	Footer  *Attr // The footer values.
	Caption *Attr // The caption text.
}

// maxAttrLen returns the greatest byte length of non-caption attribute markup.
func (o *ContentStyle) maxAttrLen() int {
	n := 0
	for _, a := range [...]*Attr{o.Header, o.Body, o.Footer} {
		n = max(n, a.size())
	}
	return n
}

// resolve returns the attribute for a table part.
func (o *ContentStyle) resolve(sc Scope) *Attr {
	switch sc {
	case ScopeHeader:
		return o.Header
	case ScopeBody:
		return o.Body
	case ScopeFooter:
		return o.Footer
	}
	return nil
}

var (
	// StyleASCII uses only ASCII glyphs ("+", "-", "|") and is suitable where
	// Unicode line-drawing characters are unavailable.
	StyleASCII = Style{
		Border: BorderStyle{
			Top:      glyphsASCIITop,
			Header:   glyphsASCIIHeader,
			Body:     glyphsASCIIBody,
			Footer:   glyphsASCIIHeader,
			Bottom:   glyphsASCIIBottom,
			Vertical: verticalASCII,
		},
	}

	// StyleLight uses Unicode light box-drawing characters ("│", "─", "┌", ...).
	StyleLight = Style{
		Border: BorderStyle{
			Top:      glyphsLightTop,
			Header:   glyphsLightHeader,
			Body:     glyphsLightBody,
			Footer:   glyphsLightHeader,
			Bottom:   glyphsLightBottom,
			Vertical: verticalLight,
		},
	}

	// StyleColoredLight is StyleLight with a dim border and caption, and a bold
	// header and footer.
	StyleColoredLight = colored(StyleLight)

	// StyleRounded is StyleLight with rounded corner glyphs ("╭", "╯", ...).
	StyleRounded = Style{
		Border: BorderStyle{
			Top:      glyphsRoundedTop,
			Header:   glyphsRoundedHeader,
			Body:     glyphsRoundedBody,
			Footer:   glyphsRoundedHeader,
			Bottom:   glyphsRoundedBottom,
			Vertical: verticalLight,
		},
	}

	// StyleColoredRounded is StyleRounded with a dim border and caption, and a
	// bold header and footer.
	StyleColoredRounded = colored(StyleRounded)

	// StyleHeavy uses Unicode heavy box-drawing characters ("┃", "━", "┏", ...).
	StyleHeavy = Style{
		Border: BorderStyle{
			Top:      glyphsHeavyTop,
			Header:   glyphsHeavyHeader,
			Body:     glyphsHeavyBody,
			Footer:   glyphsHeavyHeader,
			Bottom:   glyphsHeavyBottom,
			Vertical: verticalHeavy,
		},
	}

	// StyleColoredHeavy is StyleHeavy with a dim border and caption, and a bold
	// header and footer.
	StyleColoredHeavy = colored(StyleHeavy)

	// StyleDouble uses Unicode double box-drawing characters ("║", "═", "╔", ...).
	StyleDouble = Style{
		Border: BorderStyle{
			Top:      glyphsDoubleTop,
			Header:   glyphsDoubleHeader,
			Body:     glyphsDoubleBody,
			Footer:   glyphsDoubleHeader,
			Bottom:   glyphsDoubleBottom,
			Vertical: verticalDouble,
		},
	}

	// StyleColoredDouble is StyleDouble with a dim border and caption, and a bold
	// header and footer.
	StyleColoredDouble = colored(StyleDouble)
)

var (
	// jointsASCII crosses an ASCII horizontal with an ASCII vertical:
	// every joint uses the same glyph.
	jointsASCII = Joints{
		"+", "+", "+", "|",
		"+", "+", "+", "|",
		"+", "+", "+", "|",
		"-", "-", "-", " ",
	}

	// jointsLight crosses a light horizontal with a light vertical.
	jointsLight = Joints{
		"┼", "┤", "├", "│",
		"┬", "┐", "┌", "╷",
		"┴", "┘", "└", "╵",
		"─", "╴", "╶", " ",
	}

	// jointsRounded is jointsLight with rounded corners.
	jointsRounded = Joints{
		"┼", "┤", "├", "│",
		"┬", "╮", "╭", "╷",
		"┴", "╯", "╰", "╵",
		"─", "╴", "╶", " ",
	}

	// jointsDouble crosses a double horizontal with a light vertical.
	// A double horizontal has no one-sided form, so a lone arm keeps
	// the fill.
	jointsDouble = Joints{
		"╪", "╡", "╞", "│",
		"╤", "╕", "╒", "╷",
		"╧", "╛", "╘", "╵",
		"═", "═", "═", " ",
	}

	// jointsHeavy crosses a heavy horizontal with a light vertical.
	jointsHeavy = Joints{
		"┿", "┥", "┝", "│",
		"┯", "┑", "┍", "╷",
		"┷", "┙", "┕", "╵",
		"━", "╸", "╺", " ",
	}

	// jointsHeavyOnHeavy crosses a heavy horizontal with a heavy vertical.
	jointsHeavyOnHeavy = Joints{
		"╋", "┫", "┣", "┃",
		"┳", "┓", "┏", "╻",
		"┻", "┛", "┗", "╹",
		"━", "╸", "╺", " ",
	}

	// jointsLightOnHeavy crosses a light horizontal with a heavy vertical.
	jointsLightOnHeavy = Joints{
		"╂", "┨", "┠", "┃",
		"┰", "┒", "┎", "╻",
		"┸", "┚", "┖", "╹",
		"─", "╴", "╶", " ",
	}

	// jointsDoubleOnDouble crosses a double horizontal with a double vertical.
	jointsDoubleOnDouble = Joints{
		"╬", "╣", "╠", "║",
		"╦", "╗", "╔", "║",
		"╩", "╝", "╚", "║",
		"═", "═", "═", " ",
	}

	// jointsLightOnDouble crosses a light horizontal with a double vertical.
	jointsLightOnDouble = Joints{
		"╫", "╢", "╟", "║",
		"╥", "╖", "╓", "║",
		"╨", "╜", "╙", "║",
		"─", "╴", "╶", " ",
	}
)

var (
	// glyphsASCIITop is the top border of StyleASCII.
	glyphsASCIITop = &Horizontal{
		Inner: jointsASCII,
		Outer: jointsASCII,
		Fill:  "-",
	}

	// glyphsASCIIHeader is the header and footer separator of StyleASCII.
	glyphsASCIIHeader = &Horizontal{
		Inner: jointsASCII,
		Outer: jointsASCII,
		Fill:  "-",
	}

	// glyphsASCIIBody is the separator between the data rows of StyleASCII.
	glyphsASCIIBody = &Horizontal{
		Inner: jointsASCII,
		Outer: jointsASCII,
		Fill:  "-",
	}

	// glyphsASCIIBottom is the bottom border of StyleASCII.
	glyphsASCIIBottom = &Horizontal{
		Inner: jointsASCII,
		Outer: jointsASCII,
		Fill:  "-",
	}

	// glyphsLightTop is the top border of StyleLight.
	glyphsLightTop = &Horizontal{
		Inner: jointsLight,
		Outer: jointsLight,
		Fill:  "─",
	}

	// glyphsLightHeader is the header and footer separator of StyleLight,
	// drawn double to distinguish it from row separators.
	glyphsLightHeader = &Horizontal{
		Inner: jointsDouble,
		Outer: jointsDouble,
		Fill:  "═",
	}

	// glyphsLightBody is the separator between the data rows of StyleLight.
	glyphsLightBody = &Horizontal{
		Inner: jointsLight,
		Outer: jointsLight,
		Fill:  "─",
	}

	// glyphsLightBottom is the bottom border of StyleLight.
	glyphsLightBottom = &Horizontal{
		Inner: jointsLight,
		Outer: jointsLight,
		Fill:  "─",
	}

	// glyphsRoundedTop is the top border of StyleRounded.
	glyphsRoundedTop = &Horizontal{
		Inner: jointsLight,
		Outer: jointsRounded,
		Fill:  "─",
	}

	// glyphsRoundedHeader is the header and footer separator of StyleRounded.
	glyphsRoundedHeader = &Horizontal{
		Inner: jointsDouble,
		Outer: jointsDouble,
		Fill:  "═",
	}

	// glyphsRoundedBody is the separator between the data rows of StyleRounded.
	glyphsRoundedBody = &Horizontal{
		Inner: jointsLight,
		Outer: jointsRounded,
		Fill:  "─",
	}

	// glyphsRoundedBottom is the bottom border of StyleRounded.
	glyphsRoundedBottom = &Horizontal{
		Inner: jointsLight,
		Outer: jointsRounded,
		Fill:  "─",
	}

	// glyphsHeavyTop is the top border of StyleHeavy.
	glyphsHeavyTop = &Horizontal{
		Inner: jointsHeavy,
		Outer: jointsHeavyOnHeavy,
		Fill:  "━",
	}

	// glyphsHeavyHeader is the header and footer separator of StyleHeavy.
	glyphsHeavyHeader = &Horizontal{
		Inner: jointsHeavy,
		Outer: jointsHeavyOnHeavy,
		Fill:  "━",
	}

	// glyphsHeavyBody is the light separator between data rows inside the heavy frame.
	glyphsHeavyBody = &Horizontal{
		Inner: jointsLight,
		Outer: jointsLightOnHeavy,
		Fill:  "─",
	}

	// glyphsHeavyBottom is the bottom border of StyleHeavy.
	glyphsHeavyBottom = &Horizontal{
		Inner: jointsHeavy,
		Outer: jointsHeavyOnHeavy,
		Fill:  "━",
	}

	// glyphsDoubleTop is the top border of StyleDouble.
	glyphsDoubleTop = &Horizontal{
		Inner: jointsDouble,
		Outer: jointsDoubleOnDouble,
		Fill:  "═",
	}

	// glyphsDoubleHeader is the header and footer separator of StyleDouble.
	glyphsDoubleHeader = &Horizontal{
		Inner: jointsDouble,
		Outer: jointsDoubleOnDouble,
		Fill:  "═",
	}

	// glyphsDoubleBody is the light separator between data rows inside the double frame.
	glyphsDoubleBody = &Horizontal{
		Inner: jointsLight,
		Outer: jointsLightOnDouble,
		Fill:  "─",
	}

	// glyphsDoubleBottom is the bottom border of StyleDouble.
	glyphsDoubleBottom = &Horizontal{
		Inner: jointsDouble,
		Outer: jointsDoubleOnDouble,
		Fill:  "═",
	}
)

var (
	// verticalASCII is the vertical set for StyleASCII.
	verticalASCII = &Vertical{
		Outer: "|",
		Inner: "|",
	}

	// verticalLight is the vertical set for StyleLight and StyleRounded.
	verticalLight = &Vertical{
		Outer: "│",
		Inner: "│",
	}

	// verticalHeavy is the vertical set for StyleHeavy: a heavy frame
	// around light separators.
	verticalHeavy = &Vertical{
		Outer: "┃",
		Inner: "│",
	}

	// verticalDouble is the vertical set for StyleDouble.
	verticalDouble = &Vertical{
		Outer: "║",
		Inner: "│",
	}
)

// colored applies the default attributes to a style.
func colored(s Style) Style {
	s.Border.Attr = NewAttr(CodeFaint)
	s.Content = ContentStyle{
		Header:  NewAttr(CodeBold),
		Footer:  NewAttr(CodeBold),
		Caption: NewAttr(CodeFaint),
	}
	return s
}
