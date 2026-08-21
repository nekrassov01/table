package markdown

// minDashes is the shortest readable delimiter run. GFM accepts one dash,
// but a lone dash looks accidental.
const minDashes = 3

// solver measures compiled rows and derives the column metrics required by
// the painter.
type solver struct {
	input  compilerResult // Compiled table currently being solved.
	state  *solverState   // Reusable column-metric state.
	output solverResult   // Compiled table paired with column geometry.
}

// prepare initializes column metrics.
func (o *solver) prepare() {
	state := o.state
	columns := o.input.columns
	columnCount := len(columns)
	if cap(state.columnMetrics) < columnCount {
		state.columnMetrics = make([]columnMetric, columnCount)
	} else {
		state.columnMetrics = state.columnMetrics[:columnCount]
	}
	for index := range state.columnMetrics {
		align := columns[index].align
		if index < o.input.option.indexOffset {
			align = AlignRight
		}
		state.columnMetrics[index] = columnMetric{
			box:       box{align: align},
			separator: resolveSeparator(align),
		}
	}
	o.output.metrics = state.columnMetrics
}

// solve measures compiled rows and resolves final column geometry from their
// content and the delimiter-row minimums. Only Table calls it because a Stream
// cannot look ahead to future rows.
func (o *solver) solve() {
	o.measureRow(&o.input.header)
	o.measureRows(o.input.body)
	for index := range o.state.columnMetrics {
		metric := &o.state.columnMetrics[index]
		metric.box.width = max(metric.box.width, metric.separator.width)
	}
}

// measureRows accumulates metrics from rows in input order.
func (o *solver) measureRows(rows []row) {
	for index := range rows {
		o.measureRow(&rows[index])
	}
}

// measureRow accumulates the maximum cell width in every column.
func (o *solver) measureRow(r *row) {
	metrics := o.state.columnMetrics
	for index := range metrics {
		cell := &r.cells[index]
		metric := &metrics[index]
		metric.box.width = max(metric.box.width, cell.width)
	}
}

// solverResult pairs compiled rows with their column geometry.
type solverResult struct {
	compilerResult

	metrics []columnMetric // Metrics in logical column order.
}

// columnMetric holds the geometry of one logical column.
type columnMetric struct {
	box       box       // Geometry used to paint values.
	separator separator // Geometry used to paint the delimiter cell.
}

// box is the geometry needed to paint one column.
type box struct {
	width int       // Display width cells pad to; zero for a Stream.
	align AlignSide // Column alignment encoded by the delimiter row.
}

// separator is the resolved geometry of one delimiter cell.
type separator struct {
	width int  // Minimum display width, including alignment marks.
	lead  bool // Whether the dash run has a leading colon.
	trail bool // Whether the dash run has a trailing colon.
}

// resolveSeparator derives delimiter geometry from a column alignment.
func resolveSeparator(side AlignSide) separator {
	resolved := separator{width: minDashes}
	switch side {
	case AlignLeft:
		resolved.lead = true
	case AlignRight:
		resolved.trail = true
	case AlignCenter:
		resolved.lead = true
		resolved.trail = true
	}
	if resolved.lead {
		resolved.width++
	}
	if resolved.trail {
		resolved.width++
	}
	return resolved
}
