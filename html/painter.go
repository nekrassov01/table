package html

import (
	"io"
	"strconv"

	"github.com/nekrassov01/table/internal/repeat"
)

// painter writes solved logical rows as HTML elements.
type painter struct {
	input solverResult  // Solved table to paint.
	state *painterState // Reusable painting state.
	w     io.Writer     // Output destination.
	err   error         // Sticky output error.
}

// prepare sizes the reusable output buffer for the largest block written at
// once.
func (o *painter) prepare() {
	state := o.state
	if cap(state.line) > cap(state.lineBacking) {
		state.lineBacking = state.line[:0]
	}
	option := o.input.option
	attrs := &option.attrs
	sectionAttrLen := maxAttrLen(attrs.Header.Section, attrs.Body.Section, attrs.Footer.Section)
	rowAttrLen := maxAttrLen(attrs.Header.Row, attrs.Body.Row, attrs.Footer.Row)
	cellAttrLen := maxAttrLen(attrs.Header.Cell, attrs.Body.Cell, attrs.Footer.Cell)
	// Surrounding table, caption, section, and row markup needs at most 86
	// fixed bytes for indentation, tags, attribute syntax, caption-side, and
	// newlines.
	size := 86 + len(o.input.caption) +
		attrs.Table.size() + attrs.Caption.size() + max(sectionAttrLen, rowAttrLen)
	for index := range o.input.columns {
		column := &o.input.columns[index]
		columnAttrLen := maxAttrLen(
			column.attrs.Resolve(ScopeHeader),
			column.attrs.Resolve(ScopeBody),
			column.attrs.Resolve(ScopeFooter),
		)
		// Each cell needs at most 86 fixed bytes for indentation, tags,
		// attributes, alignment, both span counts, and its newline.
		cellSize := 86 + o.input.columnSizes[index] +
			columnAttrLen + cellAttrLen
		size += cellSize
	}
	if cap(state.lineBacking) < size {
		state.lineBacking = make([]byte, size)
	}
	state.line = state.lineBacking[:0]
}

// paintHeader paints the table opening, caption, header, and body opening.
func (o *painter) paintHeader() {
	o.resetLine()
	o.writeOpenTag(element{
		tag:  "table",
		attr: o.input.option.attrs.Table,
	}, "")
	o.writeNewline()
	o.paintCaption()
	o.flushLine()
	o.paintBand(o.input.header, ScopeHeader)
	if len(o.input.body) > 0 {
		o.paintOpenSection(o.resolveSection(ScopeBody))
	}
}

// paintBody paints body rows without opening or closing the body section.
func (o *painter) paintBody() {
	section := o.resolveSection(ScopeBody)
	for index := range o.input.body {
		if o.err != nil {
			return
		}
		o.paintRow(o.input.body[index], section)
	}
}

// paintFooter paints the body closing, footer, and table closing.
func (o *painter) paintFooter() {
	if o.input.hasPreviousBody || len(o.input.body) > 0 {
		o.paintCloseSection(o.resolveSection(ScopeBody))
	}
	o.paintBand(o.input.footer, ScopeFooter)
	o.resetLine()
	o.writeCloseTag("table")
	o.writeNewline()
	o.flushLine()
}

// paintCaption paints the caption inside the table.
func (o *painter) paintCaption() {
	option := o.input.option
	if option.caption == "" {
		return
	}
	o.writeIndent(1)
	o.writeOpenTag(element{
		tag:  "caption",
		attr: option.attrs.Caption,
	}, resolveCaptionSide(option.captionSide))
	o.state.line = append(o.state.line, o.input.caption...)
	o.writeCloseTag("caption")
	o.writeNewline()
}

// paintBand paints one header or footer band within its section.
func (o *painter) paintBand(rows []row, sc Scope) {
	if len(rows) == 0 {
		return
	}
	section := o.resolveSection(sc)
	o.paintOpenSection(section)
	for index := range rows {
		if o.err != nil {
			return
		}
		o.paintRow(rows[index], section)
	}
	o.paintCloseSection(section)
}

// paintOpenSection paints a table-section opening tag.
func (o *painter) paintOpenSection(section section) {
	o.resetLine()
	o.writeIndent(1)
	o.writeOpenTag(section.section, "")
	o.writeNewline()
	o.flushLine()
}

// paintCloseSection paints a table-section closing tag.
func (o *painter) paintCloseSection(section section) {
	o.resetLine()
	o.writeIndent(1)
	o.writeCloseTag(section.section.tag)
	o.writeNewline()
	o.flushLine()
}

// paintRow paints one compiled row as a single output block.
func (o *painter) paintRow(r row, section section) {
	o.resetLine()
	o.writeIndent(2)
	o.writeOpenTag(section.row, "")
	o.writeNewline()
	for index := range r.cells {
		compiled := &r.cells[index]
		if compiled.colspan == 0 {
			continue
		}
		column := &o.input.columns[index]
		align := column.aligns.Resolve(section.scope)
		if index < o.input.option.indexOffset && section.scope != ScopeHeader {
			align = AlignRight
		}
		o.writeIndent(3)
		o.writeOpenCell(
			section.cell,
			column.attrs.Resolve(section.scope),
			align,
			compiled,
		)
		o.writeCellValue(compiled)
		o.writeCloseCell(section.cell)
		o.writeNewline()
	}
	o.writeIndent(2)
	o.writeCloseTag(section.row.tag)
	o.writeNewline()
	o.flushLine()
}

// resolveSection returns the element names and attributes for one table part.
func (o *painter) resolveSection(sc Scope) section {
	attrs := o.input.option.attrs.resolve(sc)
	sectionTag := "tbody"
	cellTag := "td"
	switch sc {
	case ScopeHeader:
		sectionTag = "thead"
		cellTag = "th"
	case ScopeFooter:
		sectionTag = "tfoot"
	}
	return section{
		section: element{
			tag:  sectionTag,
			attr: attrs.Section,
		},
		row: element{
			tag:  "tr",
			attr: attrs.Row,
		},
		cell: element{
			tag:  cellTag,
			attr: attrs.Cell,
		},
		scope: sc,
	}
}

// resetLine resets the current output block.
func (o *painter) resetLine() {
	o.state.line = o.state.line[:0]
}

// flushLine writes the current output block.
func (o *painter) flushLine() {
	if o.err != nil || len(o.state.line) == 0 {
		return
	}
	_, err := o.w.Write(o.state.line)
	if err != nil {
		o.err = newWriteError(err)
	}
}

// writeOpenTag appends an opening tag.
func (o *painter) writeOpenTag(element element, extraStyle string) {
	o.state.line = append(o.state.line, '<')
	o.state.line = append(o.state.line, element.tag...)
	o.writeClass(element.attr.Class, "")
	o.writeStyle(element.attr.Style, extraStyle, AlignDefault)
	o.state.line = append(o.state.line, '>')
}

// writeCloseTag appends a closing tag.
func (o *painter) writeCloseTag(tag string) {
	o.state.line = append(o.state.line, '<', '/')
	o.state.line = append(o.state.line, tag...)
	o.state.line = append(o.state.line, '>')
}

// writeOpenCell appends a cell opening tag and its attributes.
func (o *painter) writeOpenCell(element element, columnAttr Attr, align AlignSide, compiled *cell) {
	o.state.line = append(o.state.line, '<')
	o.state.line = append(o.state.line, element.tag...)
	o.writeClass(element.attr.Class, columnAttr.Class)
	o.writeStyle(element.attr.Style, columnAttr.Style, align)
	if compiled.rowspan > 1 {
		o.writeSpan(` rowspan="`, compiled.rowspan)
	}
	if compiled.colspan > 1 {
		o.writeSpan(` colspan="`, compiled.colspan)
	}
	o.state.line = append(o.state.line, '>')
}

// writeCloseCell appends a cell closing tag.
func (o *painter) writeCloseCell(element element) {
	o.writeCloseTag(element.tag)
}

// writeClass appends a joined class attribute.
func (o *painter) writeClass(class, columnClass string) {
	if class == "" && columnClass == "" {
		return
	}
	o.state.line = append(o.state.line, ` class="`...)
	if class != "" {
		o.state.line = append(o.state.line, class...)
		if columnClass != "" {
			o.state.line = append(o.state.line, ' ')
		}
	}
	o.state.line = append(o.state.line, columnClass...)
	o.state.line = append(o.state.line, '"')
}

// writeStyle appends a joined style attribute.
func (o *painter) writeStyle(style, extra string, align AlignSide) {
	if style == "" && extra == "" && align == AlignDefault {
		return
	}
	o.state.line = append(o.state.line, ` style="`...)
	first := true
	if style != "" {
		o.state.line = append(o.state.line, style...)
		first = false
	}
	if extra != "" {
		if !first {
			o.state.line = append(o.state.line, ';')
		}
		o.state.line = append(o.state.line, extra...)
		first = false
	}
	if align != AlignDefault {
		if !first {
			o.state.line = append(o.state.line, ';')
		}
		o.state.line = append(o.state.line, "text-align:"...)
		o.state.line = append(o.state.line, align.String()...)
	}
	o.state.line = append(o.state.line, '"')
}

// writeSpan appends a rowspan or colspan attribute.
func (o *painter) writeSpan(name string, count int) {
	o.state.line = append(o.state.line, name...)
	o.state.line = strconv.AppendInt(o.state.line, int64(count), 10)
	o.state.line = append(o.state.line, '"')
}

// writeCellValue appends a cell value wrapped in its inner markup.
func (o *painter) writeCellValue(compiled *cell) {
	if compiled.value == "" {
		return
	}
	if compiled.decoration != nil {
		o.state.line = append(o.state.line, compiled.decoration.Prefix...)
	}
	if compiled.color != nil {
		o.state.line = append(o.state.line, compiled.color.Prefix...)
	}
	o.state.line = append(o.state.line, compiled.value...)
	if compiled.color != nil {
		o.state.line = append(o.state.line, compiled.color.Suffix...)
	}
	if compiled.decoration != nil {
		o.state.line = append(o.state.line, compiled.decoration.Suffix...)
	}
}

// writeIndent appends indentation levels.
func (o *painter) writeIndent(level int) {
	o.state.line = repeat.AppendSpaces(o.state.line, level*2)
}

// writeNewline appends a newline.
func (o *painter) writeNewline() {
	o.state.line = append(o.state.line, '\n')
}

// section describes a table part by its section, row, and cell elements.
type section struct {
	section element
	row     element
	cell    element
	scope   Scope
}

// element is one HTML element with its attributes.
type element struct {
	tag  string
	attr Attr
}

// maxAttrLen returns the greatest attribute-value byte length.
func maxAttrLen(a, b, c Attr) int {
	return max(a.size(), b.size(), c.size())
}
