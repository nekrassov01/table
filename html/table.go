package html

import (
	"io"

	"github.com/nekrassov01/table"
	"github.com/nekrassov01/table/internal/column"
)

var _ table.Tabular = (*Table)(nil)

// Table renders a complete set of rows as an HTML table.
type Table struct {
	option option    // Options fixed at construction.
	w      io.Writer // Output destination.
}

// NewTable creates a [Table] that writes to w with the given options.
func NewTable(w io.Writer, opts ...Option) *Table {
	t := &Table{w: w}
	t.option.apply(opts...)
	return t
}

// Render writes rows as one complete HTML table.
func (o *Table) Render(rows [][]any) error {
	header := o.option.header
	body := rows
	if column.MaxColumns(header) == 0 {
		for len(body) > 0 && len(body[0]) == 0 {
			body = body[1:]
		}
	}
	bodyColumns := 0
	if len(body) > 0 {
		bodyColumns = len(body[0])
	}
	a := acquireArena()
	var footer [][]string
	if o.option.footer != nil {
		footer = o.option.footer()
	}
	config := a.newConfig(&o.option, footer, len(body), bodyColumns)
	config.prepare()
	configured := config.output
	if len(configured.columns) == 0 {
		a.release()
		return nil
	}
	compiler := a.newCompiler(configured)
	compiler.prepare()
	compiler.compileHeader()
	compiler.compileBody(body)
	compiler.compileFooter()
	if compiler.err != nil {
		a.release()
		return compiler.err
	}
	solver := a.newSolver(compiler.output)
	solver.solve()
	painter := a.newPainter(solver.output, o.w)
	painter.prepare()
	painter.paintHeader()
	painter.paintBody()
	painter.paintFooter()
	a.release()
	return painter.err
}
