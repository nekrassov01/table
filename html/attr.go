package html

// Attr holds the optional class and inline style of one HTML element. Options
// escape both fields for use in a double-quoted attribute.
type Attr struct {
	Class string // CSS class names, space-separated.
	Style string // Inline CSS declarations, semicolon-separated.
}

// escape returns o with both fields escaped for an HTML attribute.
func (o Attr) escape() Attr {
	o.Class = escapeAttr(o.Class)
	o.Style = escapeAttr(o.Style)
	return o
}

// size returns the byte length of the attribute values.
func (o Attr) size() int {
	return len(o.Class) + len(o.Style)
}

// TableAttr defines attributes for the table and each table part.
type TableAttr struct {
	Table   Attr        // The table element.
	Caption Attr        // The caption element.
	Header  SectionAttr // The header section.
	Body    SectionAttr // The body section.
	Footer  SectionAttr // The footer section.
}

// escape returns o with every class and style escaped for an HTML attribute.
func (o TableAttr) escape() TableAttr {
	o.Table = o.Table.escape()
	o.Caption = o.Caption.escape()
	o.Header = o.Header.escape()
	o.Body = o.Body.escape()
	o.Footer = o.Footer.escape()
	return o
}

// resolve returns the section attributes for sc.
func (o *TableAttr) resolve(sc Scope) SectionAttr {
	switch sc {
	case ScopeHeader:
		return o.Header
	case ScopeBody:
		return o.Body
	case ScopeFooter:
		return o.Footer
	}
	return SectionAttr{}
}

// SectionAttr defines attributes for the elements within one table part.
type SectionAttr struct {
	Section Attr // The section element: thead, tbody or tfoot.
	Row     Attr // Every row of the section.
	Cell    Attr // Every cell of the section.
}

// escape returns o with every class and style escaped for an HTML attribute.
func (o SectionAttr) escape() SectionAttr {
	o.Section = o.Section.escape()
	o.Row = o.Row.escape()
	o.Cell = o.Cell.escape()
	return o
}
