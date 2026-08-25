package text

import (
	"io"

	"github.com/nekrassov01/table/internal/repeat"
	"github.com/nekrassov01/table/internal/value"
	"github.com/nekrassov01/table/internal/width"
)

// ellipsis marks truncated content.
const ellipsis = "..."

// painter arranges solved logical rows into physical lines and writes bordered
// text.
type painter struct {
	input   solverResult  // Logical rows and solved column metrics to paint.
	state   *painterState // Reusable painting state.
	strings *value.Store  // Storage for truncated values.
	w       io.Writer     // Output destination.
	err     error         // Sticky output error.
}

// prepare sizes reusable painting buffers and resets cached horizon state.
func (o *painter) prepare() {
	state := o.state
	state.segments = state.segments[:0]
	state.rowspans = 0
	option := o.input.option
	metrics := o.input.metrics
	columns := o.input.columns
	if cap(state.layouts) < len(metrics) {
		state.layouts = make([]layout, 0, len(metrics))
	} else {
		state.layouts = state.layouts[:0]
	}
	if len(metrics) == 0 {
		const lineCap = 1
		if cap(state.lineBacking) < lineCap {
			state.lineBacking = make([]byte, lineCap)
		}
		state.line = state.lineBacking[:0:lineCap]
		state.horizon = nil
		return
	}
	glyphLen := max(1, option.style.Border.maxGlyphLen())
	columnCount := len(metrics)
	rowSize := glyphLen*(columnCount+1) + 1
	glyphCount := columnCount + 1
	for i := range metrics {
		metric := &metrics[i]
		boxWidth := metric.box.totalWidth()
		rowSize += boxWidth + metric.overhead
		glyphCount += boxWidth
	}
	borderAttrLen := option.style.Border.Attr.size()
	rowSize += (columnCount + 1) * borderAttrLen
	attrLen := max(option.style.Content.maxAttrLen(), int(o.input.attrLen))
	for i := range metrics {
		col := &columns[i]
		attrs := &col.transformer.attrs
		for _, attr := range [3]*Attr{
			attrs.Resolve(ScopeHeader),
			attrs.Resolve(ScopeBody),
			attrs.Resolve(ScopeFooter),
		} {
			attrLen = max(attrLen, attr.size())
		}
	}
	rowSize += columnCount * attrLen
	horizonSize := glyphCount*glyphLen + 1
	horizonSize += (columnCount*2 + 2) * borderAttrLen
	captionSize := 0
	if option.caption != "" {
		captionSize = len(option.caption) + 1
		captionSize += option.style.Content.Caption.size()
	}
	lineCap := max(rowSize, horizonSize, captionSize)
	rowspans := o.input.rowspanMask
	cacheHorizon := option.style.Border.Body != nil && (!option.compact || rowspans != 0)
	backingCap := lineCap
	if cacheHorizon {
		backingCap *= 2
	}
	if cap(state.lineBacking) < backingCap {
		state.lineBacking = make([]byte, backingCap)
	}
	state.line = state.lineBacking[:0:lineCap]
	state.horizon = nil
	if cacheHorizon {
		state.horizon = state.lineBacking[lineCap:lineCap:backingCap]
	}
}

// paintHeader paints the header and its surrounding horizons.
func (o *painter) paintHeader() {
	option := o.input.option
	header := o.input.header
	body := o.input.body
	footer := o.input.footer
	if option.captionSide == CaptionTop {
		o.paintCaption()
	}
	nextRows := body
	if len(nextRows) == 0 {
		nextRows = footer
	}
	downBars := allBars
	if len(nextRows) > 0 {
		downBars = nextRows[0].bars
	}
	if len(header) == 0 {
		o.paintHorizon(option.style.Border.Top, 0, noBars, downBars)
		return
	}
	o.paintHorizon(option.style.Border.Top, 0, noBars, header[0].bars)
	for i := range header {
		row := &header[i]
		o.paintRow(row, ScopeHeader)
		if i+1 < len(header) {
			h := option.style.Border.Body
			if h == nil {
				h = option.style.Border.Header
			}
			next := &header[i+1]
			o.paintHorizon(h, row.rowspans, row.bars, next.bars)
		}
	}
	last := &header[len(header)-1]
	o.paintHorizon(option.style.Border.Header, 0, last.bars, downBars)
}

// paintBody paints compiled body rows using solved geometry.
func (o *painter) paintBody() {
	body := o.input.body
	upBars := o.input.previousBars
	for i := range body {
		if o.err != nil {
			break
		}
		row := &body[i]
		if o.input.hasPreviousBody || i > 0 {
			o.paintRowHorizon(row.rowspans, upBars, row.bars)
		}
		o.paintRow(row, ScopeBody)
		upBars = row.bars
	}
}

// paintFooter paints the footer and its surrounding horizons.
func (o *painter) paintFooter() {
	option := o.input.option
	footer := o.input.footer
	upBars := o.input.lastBars
	hasBody := o.input.hasPreviousBody || len(o.input.body) > 0
	for i := range footer {
		row := &footer[i]
		h := option.style.Border.Footer
		if i == 0 && !hasBody {
			h = nil
		}
		if i > 0 && option.style.Border.Body != nil {
			h = option.style.Border.Body
		}
		o.paintHorizon(h, row.rowspans, upBars, row.bars)
		o.paintRow(row, ScopeFooter)
		upBars = row.bars
	}
	o.paintHorizon(option.style.Border.Bottom, 0, upBars, noBars)
	if option.captionSide != CaptionTop {
		o.paintCaption()
	}
}

// paintCaption paints the caption.
func (o *painter) paintCaption() {
	option := o.input.option
	if option.caption == "" {
		return
	}
	o.resetLine()
	o.writeValue(option.caption, option.style.Content.Caption)
	o.writeNewline()
	o.flushLine(o.state.line)
}

// paintRow paints a logical row as one or more physical lines.
func (o *painter) paintRow(r *row, sc Scope) {
	option := o.input.option
	vertical := option.style.Border.Vertical
	attr := option.style.Content.resolve(sc)
	height := o.layoutRow(r, sc)
	layouts := o.state.layouts
	for line := range height {
		o.resetLine()
		if vertical != nil {
			o.writeGlyph(vertical.Outer)
		}
		for cellIndex := range layouts {
			if cellIndex > 0 && vertical != nil {
				o.writeGlyph(vertical.Inner)
			}
			o.paintCell(&layouts[cellIndex], line, attr)
		}
		if vertical != nil {
			o.writeGlyph(vertical.Outer)
		}
		o.writeNewline()
		o.flushLine(o.state.line)
	}
}

// paintCell paints one physical line of a laid-out cell.
func (o *painter) paintCell(cell *layout, line int, attr *Attr) {
	var seg segment
	switch {
	case line < len(cell.segments):
		seg = cell.segments[line]
	case len(cell.segments) == 0 && line == 0:
		seg.value = cell.value
		seg.width = cell.width
	}
	if cell.attr != nil {
		attr = cell.attr
	}
	o.paintSegment(&cell.box, seg, cell.align, attr)
}

// paintSegment paints a cell segment with padding and alignment.
func (o *painter) paintSegment(cellBox *box, seg segment, align AlignSide, attr *Attr) {
	pad := max(cellBox.width-seg.width, 0)
	leftPad, rightPad := 0, pad
	switch align {
	case AlignRight:
		leftPad, rightPad = pad, 0
	case AlignCenter:
		leftPad = pad / 2
		rightPad = pad - leftPad
	}
	o.writeSpaces(cellBox.lPad + leftPad)
	o.writeValue(seg.value, attr)
	o.writeSpaces(rightPad + cellBox.rPad)
}

// paintRowHorizon paints the horizon between two body rows.
func (o *painter) paintRowHorizon(rowspans, upBars, downBars uint64) {
	option := o.input.option
	state := o.state
	h := option.style.Border.Body
	sameBars := upBars == downBars
	if sameBars && h == nil {
		return
	}
	enabledRowspans := o.input.rowspanMask
	if sameBars && option.compact && rowspans&enabledRowspans == enabledRowspans {
		return
	}
	if h == nil {
		h = option.style.Border.Header
	}
	if h == nil {
		return
	}
	cacheable := upBars == allBars && downBars == allBars
	if cacheable && len(state.horizon) != 0 && state.rowspans == rowspans {
		o.flushLine(state.horizon)
		return
	}
	o.paintHorizon(h, rowspans, upBars, downBars)
	if !cacheable || len(state.horizon) != 0 {
		return
	}
	state.rowspans = rowspans
	state.horizon = append(state.horizon, state.line...)
}

// paintHorizon paints a horizontal border.
func (o *painter) paintHorizon(h *Horizontal, rowspans, upBars, downBars uint64) {
	metrics := o.input.metrics
	if h == nil || len(metrics) == 0 {
		return
	}
	o.resetLine()
	fillWidth := width.StringWidth(h.Fill)
	last := len(metrics) - 1
	up, down := upBars&1 != 0, downBars&1 != 0
	o.writeGlyph(h.Outer.resolve(up, down, false, !hasRowspan(rowspans, 0)))
	for i := range metrics {
		if i > 0 {
			bit := uint64(1) << uint(i)
			o.writeGlyph(h.Inner.resolve(
				bit == 0 || upBars&bit != 0,
				bit == 0 || downBars&bit != 0,
				!hasRowspan(rowspans, i-1),
				!hasRowspan(rowspans, i),
			))
		}
		width := metrics[i].box.totalWidth()
		if hasRowspan(rowspans, i) {
			o.writeSpaces(width)
		} else {
			o.writeFill(h.Fill, fillWidth, width)
		}
	}
	o.writeGlyph(h.Outer.resolve(up, down, !hasRowspan(rowspans, last), false))
	o.writeNewline()
	o.flushLine(o.state.line)
}

// layoutRow arranges one logical row into reusable per-cell layouts.
func (o *painter) layoutRow(r *row, sc Scope) int {
	state := o.state
	state.layouts = state.layouts[:0]
	state.segments = state.segments[:0]
	input := &o.input
	height := 1
	columnIndex := 0
	for columnIndex < len(input.metrics) {
		end := columnIndex + 1
		for end < len(input.metrics) {
			bit := uint64(1) << uint(end)
			if bit == 0 || r.bars&bit != 0 {
				break
			}
			end++
		}
		metric := &input.metrics[columnIndex]
		col := &input.columns[columnIndex]
		limit := metric.limit
		truncate := col.truncate
		cellBox := metric.box
		if end-columnIndex >= 2 {
			firstBox := &metric.box
			lastBox := &input.metrics[end-1].box
			totalWidth := lastBox.offset + lastBox.totalWidth() - firstBox.offset
			cellBox = box{
				offset: firstBox.offset,
				width:  totalWidth - firstBox.lPad - lastBox.rPad,
				lPad:   firstBox.lPad,
				rPad:   lastBox.rPad,
			}
			limit = cellBox.width
			truncate = false
		}
		compiled := r.cells[columnIndex]
		bit := uint64(1) << uint(columnIndex)
		if bit != 0 && (r.rowspans&bit != 0 || r.colspans&bit != 0) {
			compiled = cell{}
		}
		layout := o.layoutCell(compiled, limit, truncate)
		layout.attr = compiled.attr
		layout.box = cellBox
		layout.align = col.aligns.Resolve(sc)
		if layout.align == AlignDefault && sc == ScopeHeader {
			layout.align = AlignCenter
		}
		if columnIndex < input.option.indexOffset && sc != ScopeHeader {
			layout.align = AlignRight
		}
		state.layouts = append(state.layouts, layout)
		height = max(height, max(1, len(layout.segments)))
		columnIndex = end
	}
	return height
}

// layoutCell splits one logical cell into the physical lines required by its
// limit.
func (o *painter) layoutCell(cell cell, limit int, truncate bool) layout {
	value := cell.value
	cellWidth := cell.width
	hasBreak := false
	if cellWidth == 0 && value != "" {
		cellWidth, hasBreak = measureLine(value)
	}
	if !hasBreak && truncate && limit > 0 && cellWidth > limit {
		value, cellWidth = o.truncateLine(value, cellWidth, limit)
	}
	if !hasBreak && (limit <= 0 || cellWidth <= limit) {
		return layout{
			value: value,
			width: cellWidth,
		}
	}
	segmentStart := len(o.state.segments)
	maxWidth := 0
	start := 0
	for {
		line, lineWidth, next, hasBreak := scanLine(value, start)
		if truncate && limit > 0 && lineWidth > limit {
			line, lineWidth = o.truncateLine(line, lineWidth, limit)
		}
		var lineMaxWidth int
		if !truncate && limit > 0 && lineWidth > limit {
			lineMaxWidth = o.wrapLine(line, limit)
		} else {
			o.state.segments = append(o.state.segments, segment{
				value: line,
				width: lineWidth,
			})
			lineMaxWidth = lineWidth
		}
		maxWidth = max(maxWidth, lineMaxWidth)
		if !hasBreak {
			break
		}
		start = next
	}
	return layout{
		value:    value,
		segments: o.state.segments[segmentStart:],
		width:    maxWidth,
	}
}

// wrapLine splits a line into segments within the configured width.
func (o *painter) wrapLine(line string, limit int) int {
	scanner := width.NewScanner(line)
	maxWidth := 0
	segmentStart := 0
	segmentWidth := 0
	for start, end, displayWidth, ok := scanner.Next(); ok; start, end, displayWidth, ok = scanner.Next() {
		if segmentWidth+displayWidth > limit && segmentWidth > 0 {
			o.state.segments = append(o.state.segments, segment{
				value: line[segmentStart:start],
				width: segmentWidth,
			})
			maxWidth = max(maxWidth, segmentWidth)
			segmentStart = start
			segmentWidth = 0
		}
		segmentWidth += displayWidth
		if end == len(line) {
			o.state.segments = append(o.state.segments, segment{
				value: line[segmentStart:end],
				width: segmentWidth,
			})
			maxWidth = max(maxWidth, segmentWidth)
		}
	}
	return maxWidth
}

// truncateLine truncates a line to the configured display width.
func (o *painter) truncateLine(line string, lineWidth, limit int) (string, int) {
	if lineWidth <= limit {
		return line, lineWidth
	}
	contentLimit := limit - len(ellipsis)
	if contentLimit <= 0 {
		return ellipsis[:limit], limit
	}
	scanner := width.NewScanner(line)
	end := 0
	keptWidth := 0
	for start, next, displayWidth, ok := scanner.Next(); ok; start, next, displayWidth, ok = scanner.Next() {
		if keptWidth+displayWidth > contentLimit {
			end = start
			break
		}
		keptWidth += displayWidth
		end = next
	}
	mark := o.strings.Mark()
	o.strings.AppendString(line[:end])
	o.strings.AppendString(ellipsis)
	return o.strings.Since(mark), keptWidth + len(ellipsis)
}

// resetLine resets the line buffer before writing the next physical line.
func (o *painter) resetLine() {
	o.state.line = o.state.line[:0]
}

// flushLine writes one rendered physical line to the writer.
func (o *painter) flushLine(line []byte) {
	if o.err != nil || len(line) == 0 {
		return
	}
	_, err := o.w.Write(line)
	if err != nil {
		o.err = newWriteError(err)
	}
}

// writeGlyph appends a single glyph wrapped in the border attribute.
func (o *painter) writeGlyph(s string) {
	o.writeValue(s, o.input.option.style.Border.Attr)
}

// writeFill fills n display columns and wraps the result once in the border attribute.
func (o *painter) writeFill(s string, fillWidth, n int) {
	if s == "" || n <= 0 {
		return
	}
	fillCount := 0
	remainder := n
	if fillWidth > 0 {
		fillCount = n / fillWidth
		remainder = n % fillWidth
	}
	attr := o.input.option.style.Border.Attr
	hasAttr := !attr.isZero()
	if hasAttr {
		o.state.line = append(o.state.line, attr.Prefix...)
	}
	if fillCount > 0 {
		if len(s) == 1 {
			c := s[0]
			start := len(o.state.line)
			o.state.line = o.state.line[:start+fillCount]
			for i := range o.state.line[start:] {
				o.state.line[start+i] = c
			}
		} else {
			length := len(s)
			start := len(o.state.line)
			total := fillCount * length
			o.state.line = o.state.line[:start+total]
			copy(o.state.line[start:], s)
			for written := length; written < total; written *= 2 {
				copy(o.state.line[start+written:], o.state.line[start:start+written])
			}
		}
	}
	o.writeSpaces(remainder)
	if hasAttr {
		o.state.line = append(o.state.line, attr.Suffix...)
	}
}

// writeSpaces appends n spaces to line.
func (o *painter) writeSpaces(n int) {
	o.state.line = repeat.AppendSpaces(o.state.line, n)
}

// writeNewline appends a newline to line.
func (o *painter) writeNewline() {
	o.state.line = append(o.state.line, '\n')
}

// writeValue appends s to line, optionally wrapped in an attribute.
func (o *painter) writeValue(s string, attr *Attr) {
	if s == "" {
		return
	}
	if attr.isZero() {
		o.state.line = append(o.state.line, s...)
		return
	}
	o.state.line = append(o.state.line, attr.Prefix...)
	o.state.line = append(o.state.line, s...)
	o.state.line = append(o.state.line, attr.Suffix...)
}

// layout holds one logical cell arranged for immediate painting.
type layout struct {
	value    string    // Single-line value painted directly when segments is empty.
	attr     *Attr     // The resolved cell attribute.
	segments []segment // Multiline content; a view into painterState.segments.
	box      box       // The geometry this cell renders against.
	align    AlignSide // Alignment within box.
	width    int       // The widest line display width.
}

// segment holds one physical line of an arranged cell.
type segment struct {
	value string // The physical line content.
	width int    // The display width in terminal columns.
}

// hasRowspan reports whether a rowspan covers column i.
func hasRowspan(rowspans uint64, i int) bool {
	return rowspans&(uint64(1)<<uint(i)) != 0
}
