package text

import (
	"slices"

	"github.com/nekrassov01/table/internal/width"
)

// solver measures logical rows and derives the column metrics required by the
// painter.
type solver struct {
	input      compilerResult // Logical table currently being solved.
	state      *solverState   // Reusable column-metric state.
	widthLimit int            // Terminal width limit for automatic fitting.
	output     solverResult   // Logical table paired with solved column metrics.
}

// prepare initializes column metrics from compiler output.
func (o *solver) prepare() {
	state := o.state
	state.spanRequirements = state.spanRequirements[:0]
	columns := o.input.columns
	columnCount := len(columns)
	if cap(state.columnMetrics) < columnCount {
		state.columnMetrics = make([]columnMetric, columnCount)
	} else {
		state.columnMetrics = state.columnMetrics[:columnCount]
	}
	metrics := state.columnMetrics
	for i := range metrics {
		col := &columns[i]
		metrics[i] = columnMetric{
			box: box{
				lPad: col.lPad,
				rPad: col.rPad,
			},
			limit: col.limit,
		}
	}
	o.output.metrics = metrics
}

// solve measures logical rows and resolves their column geometry.
func (o *solver) solve() {
	o.measureRows(o.input.header)
	o.measureRows(o.input.body)
	o.measureRows(o.input.footer)
	o.resolveWidths()
	o.fitColumns()
	o.offsetColumns()
}

// freeze turns natural widths into limits so later stream rows cannot change
// the established geometry. Empty columns reserve one display cell.
func (o *solver) freeze() {
	metrics := o.state.columnMetrics
	changed := false
	for i := range metrics {
		metric := &metrics[i]
		if metric.limit != 0 {
			continue
		}
		if metric.box.width == 0 {
			metric.box.width = 1
			changed = true
		}
		metric.limit = metric.box.width
	}
	if changed {
		o.offsetColumns()
	}
}

// measureRows accumulates column metrics from logical rows.
func (o *solver) measureRows(rows []row) {
	for index := range rows {
		o.measureRow(&rows[index])
	}
}

// measureRow accumulates column and span width requirements from a row.
func (o *solver) measureRow(r *row) {
	metrics := o.state.columnMetrics
	requirements := o.state.spanRequirements
	for i := range metrics {
		bit := uint64(1) << uint(i)
		if r.rowspans&bit != 0 || r.colspans&bit != 0 {
			continue
		}
		metric := &metrics[i]
		cell := &r.cells[i]
		value := cell.value
		cellWidth, hasBreak := measureLine(value)
		overhead := len(value) - cellWidth
		if hasBreak {
			overhead = 0
			for start := 0; ; {
				line, lineWidth, next, hasBreak := scanLine(value, start)
				cellWidth = max(cellWidth, lineWidth)
				overhead = max(overhead, len(line)-lineWidth)
				if !hasBreak {
					break
				}
				start = next
			}
		}
		cell.width = cellWidth
		cell.hasBreak = hasBreak
		metric.overhead = max(metric.overhead, overhead)
		if r.colspans&(bit<<1) == 0 {
			metric.box.width = max(metric.box.width, cellWidth)
			continue
		}
		end := i + 1
		for end < len(metrics) && r.colspans&(1<<uint(end)) != 0 {
			end++
		}
		insertAt := len(requirements)
		for j := range requirements {
			span := &requirements[j]
			if span.start == i && span.end == end {
				span.width = max(span.width, cellWidth)
				insertAt = -1
				break
			}
			if span.start > i || (span.start == i && span.end > end) {
				insertAt = j
				break
			}
		}
		if insertAt < 0 {
			continue
		}
		span := spanRequirement{start: i, end: end, width: cellWidth}
		requirements = slices.Insert(requirements, insertAt, span)
	}
	o.state.spanRequirements = requirements
}

// resolveWidths settles column widths, including placeholder and span requirements.
func (o *solver) resolveWidths() {
	option := o.input.option
	metrics := o.state.columnMetrics
	placeholderWidth := width.StringWidth(option.placeholder)
	placeholderOverhead := len(option.placeholder) - placeholderWidth
	for index := range metrics {
		metric := &metrics[index]
		if index < option.indexOffset {
			if option.indexWidth > metric.box.width {
				metric.box.width = option.indexWidth
			}
			continue
		}
		if placeholderOverhead > metric.overhead {
			metric.overhead = placeholderOverhead
		}
		if metric.limit > 0 {
			metric.box.width = metric.limit
			continue
		}
		if placeholderWidth > metric.box.width {
			metric.box.width = placeholderWidth
		}
	}
	for _, span := range o.state.spanRequirements {
		boxWidth := 0
		for i := span.start; i < span.end; i++ {
			boxWidth += metrics[i].box.totalWidth()
		}
		if vertical := option.style.Border.Vertical; vertical != nil {
			boxWidth += (span.end - span.start - 1) * width.StringWidth(vertical.Inner)
		}
		first := &metrics[span.start].box
		last := &metrics[span.end-1].box
		deficit := span.width - (boxWidth - first.lPad - last.rPad)
		if deficit <= 0 {
			continue
		}
		flexibleColumns := 0
		for i := span.start; i < span.end; i++ {
			if metrics[i].limit == 0 {
				flexibleColumns++
			}
		}
		if flexibleColumns == 0 {
			continue
		}
		// Walk twice to avoid allocating a temporary column list.
		share := deficit / flexibleColumns
		extra := deficit % flexibleColumns
		narrowestIndex := -1
		for i := span.start; i < span.end; i++ {
			if metrics[i].limit != 0 {
				continue
			}
			metrics[i].box.width += share
			if narrowestIndex < 0 || metrics[i].box.width < metrics[narrowestIndex].box.width {
				narrowestIndex = i
			}
		}
		metrics[narrowestIndex].box.width += extra
	}
}

// fitColumns fits automatic column widths within the terminal width limit.
func (o *solver) fitColumns() {
	option := o.input.option
	metrics := o.state.columnMetrics
	if !option.autoFit {
		return
	}
	columnCount := len(metrics)
	if o.widthLimit <= 0 || columnCount == 0 {
		return
	}
	frameWidth := 0
	if vertical := option.style.Border.Vertical; vertical != nil {
		frameWidth = 2 * width.StringWidth(vertical.Outer)
		frameWidth += (columnCount - 1) * width.StringWidth(vertical.Inner)
	}
	naturalWidth := 0
	for index := range metrics {
		metric := &metrics[index]
		frameWidth += metric.box.lPad + metric.box.rPad
		naturalWidth += max(metric.box.width, 1)
	}
	budget := max(o.widthLimit-frameWidth, columnCount)
	if naturalWidth <= budget {
		return
	}
	remaining := budget
	pendingColumns := columnCount
	if option.indexOffset != 0 && columnCount > 1 {
		metric := &metrics[0]
		columnWidth := max(metric.box.width, 1)
		metric.limit = columnWidth
		remaining -= columnWidth
		pendingColumns--
	}
	for {
		share := remaining / pendingColumns
		progressed := false
		for index := range metrics {
			metric := &metrics[index]
			columnWidth := max(metric.box.width, 1)
			if metric.limit != 0 || columnWidth > share {
				continue
			}
			metric.limit = columnWidth
			remaining -= columnWidth
			pendingColumns--
			progressed = true
		}
		if !progressed {
			break
		}
	}
	share := remaining / pendingColumns
	extra := remaining % pendingColumns
	for index := range metrics {
		metric := &metrics[index]
		if metric.limit != 0 {
			continue
		}
		limit := share
		if extra > 0 {
			limit++
			extra--
		}
		metric.limit = max(limit, 1)
	}
	for index := range metrics {
		metric := &metrics[index]
		metric.box.width = metric.limit
	}
}

// offsetColumns assigns horizontal offsets from the solved column widths.
func (o *solver) offsetColumns() {
	metrics := o.state.columnMetrics
	innerWidth := 0
	if vertical := o.input.option.style.Border.Vertical; vertical != nil {
		innerWidth = width.StringWidth(vertical.Inner)
	}
	offset := 0
	for i := range metrics {
		columnBox := &metrics[i].box
		columnBox.offset = offset
		offset += columnBox.totalWidth() + innerWidth
	}
}

// solverResult is the logical table paired with its solved column metrics.
type solverResult struct {
	compilerResult

	metrics []columnMetric // Solved metrics in logical column order.
}

// columnMetric holds solved geometry, width constraints, and byte overhead for
// one column.
type columnMetric struct {
	box      box // The geometry cells render against.
	limit    int // The display width to fit within; 0 allows any.
	overhead int // Greatest byte length beyond display width among cell lines.
}

// box describes a cell's solved horizontal geometry.
type box struct {
	offset int // The horizontal offset resolved by the solver.
	width  int // The display width the cell is padded to.
	lPad   int // The left padding width.
	rPad   int // The right padding width.
}

// totalWidth returns the display width of the cell box including padding.
func (o *box) totalWidth() int {
	return o.lPad + o.width + o.rPad
}

// spanRequirement records the minimum width required by a column span.
type spanRequirement struct {
	start int // First column in the span.
	end   int // Column after the span.
	width int // Required display width.
}
