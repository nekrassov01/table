package backlog

// solver measures compiled rows and derives the column metrics required by
// the painter.
type solver struct {
	state  *solverState   // Reusable column-metric state.
	input  compilerResult // Compiled table currently being solved.
	output solverResult   // Compiled table paired with column geometry.
}

// prepare initializes column metrics.
func (o *solver) prepare() {
	state := o.state
	columnCount := len(o.input.columns)
	if cap(state.columnMetrics) < columnCount {
		state.columnMetrics = make([]columnMetric, columnCount)
	} else {
		state.columnMetrics = state.columnMetrics[:columnCount]
		clear(state.columnMetrics)
	}
	o.output.metrics = state.columnMetrics
}

// solve measures compiled rows and resolves final column geometry. Only Table
// calls it because a Stream cannot look ahead to future rows.
func (o *solver) solve() {
	o.measureRows(o.input.header)
	o.measureRows(o.input.body)
	o.measureRows(o.input.footer)
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
		cellWidth := r.cells[index].width
		metric := &metrics[index]
		metric.box.width = max(metric.box.width, cellWidth)
	}
}

// solverResult pairs compiled rows with their column geometry.
type solverResult struct {
	compilerResult

	metrics []columnMetric // Metrics in logical column order.
}

// columnMetric holds the geometry of one logical column.
type columnMetric struct {
	box box // Geometry used to paint values.
}

// box is the geometry needed to paint one column.
type box struct {
	width int // Display width cells pad to; zero for a Stream.
}
