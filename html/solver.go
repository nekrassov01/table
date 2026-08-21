package html

// solver resolves row and column span geometry in compiled HTML rows.
type solver struct {
	input  compilerResult // Compiled table being solved.
	state  *solverState   // Reusable span-geometry state.
	output solverResult   // Compiled table with span geometry resolved in place.
}

// solve resolves span geometry for rows available in the current pass. A
// stream body pass contains one row, so its future rowspan extent remains
// unknown.
func (o *solver) solve() {
	o.resolveRows(o.input.header)
	o.resolveRows(o.input.body)
	o.resolveRows(o.input.footer)
}

// resolveRows settles vertical and horizontal spans in one table part.
func (o *solver) resolveRows(rows []row) {
	var rowspans uint64
	var colspans uint64
	for index := range rows {
		rowspans |= rows[index].rowspans
		colspans |= rows[index].colspans
	}
	if rowspans != 0 {
		o.resolveRowspans(rows)
	}
	if colspans != 0 {
		o.resolveColspans(rows)
	}
}

// resolveRowspans turns vertical continuation markers into rowspan counts.
func (o *solver) resolveRowspans(rows []row) {
	counts := o.state.rowspanCounts[:]
	clear(counts)
	for rowIndex := len(rows) - 1; rowIndex >= 0; rowIndex-- {
		r := &rows[rowIndex]
		for columnIndex := range min(len(r.cells), len(counts)) {
			compiled := &r.cells[columnIndex]
			if r.rowspans&(1<<uint(columnIndex)) != 0 {
				counts[columnIndex]++
				continue
			}
			if counts[columnIndex] == 0 {
				continue
			}
			compiled.rowspan = counts[columnIndex] + 1
			for offset := 1; offset <= counts[columnIndex]; offset++ {
				rows[rowIndex+offset].cells[columnIndex].colspan = 0
			}
			counts[columnIndex] = 0
		}
	}
}

// resolveColspans absorbs adjacent equal cells into their leading cell.
func (o *solver) resolveColspans(rows []row) {
	for rowIndex := range rows {
		r := &rows[rowIndex]
		cells := r.cells
		for columnIndex := len(cells) - 1; columnIndex > 0; columnIndex-- {
			if r.colspans&(1<<uint(columnIndex)) == 0 {
				continue
			}
			if cells[columnIndex].rowspan != cells[columnIndex-1].rowspan {
				continue
			}
			cells[columnIndex-1].colspan += cells[columnIndex].colspan
			cells[columnIndex].colspan = 0
		}
	}
}

// solverResult is a logical HTML table with resolved span geometry.
type solverResult struct {
	compilerResult
}
