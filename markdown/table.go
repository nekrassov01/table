package markdown

import (
	"io"

	"github.com/nekrassov01/table"
)

var _ table.Tabular = (*Table)(nil)

// Table renders a complete set of rows as a GFM table.
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

// Render resolves column geometry from rows and writes one complete table.
func (o *Table) Render(rows [][]any) error {
	a := acquireArena()
	config := a.newConfig(&o.option, len(rows))
	config.prepare()
	if config.err != nil {
		a.release()
		return config.err
	}
	compiler := a.newCompiler(config.output)
	compiler.prepare()
	compiler.compileHeader()
	compiler.compileBody(rows)
	if compiler.err != nil {
		a.release()
		return compiler.err
	}
	solver := a.newSolver(compiler.output)
	solver.prepare()
	solver.solve()
	painter := a.newPainter(solver.output, o.w)
	painter.prepare()
	painter.paintHeader()
	painter.paintBody()
	a.release()
	return painter.err
}
