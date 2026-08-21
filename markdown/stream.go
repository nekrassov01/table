package markdown

import (
	"io"

	"github.com/nekrassov01/table"
)

var _ table.Streamer = (*Stream)(nil)

// Stream renders rows incrementally as a GFM table.
type Stream struct {
	option   option    // Options fixed at construction.
	w        io.Writer // Output destination.
	err      error     // Sticky output error.
	arena    *arena    // State retained while the stream is active.
	rendered int       // Body rows written so far; also the next index offset.
	closed   bool      // Whether Close has been called.
}

// NewStream creates a [Stream] that writes to w with the given options.
func NewStream(w io.Writer, opts ...Option) *Stream {
	s := &Stream{w: w}
	s.option.apply(opts...)
	return s
}

// Render writes one body row. The first call writes the required header and
// delimiter row.
func (o *Stream) Render(row []any) error {
	if o.err != nil {
		return o.err
	}
	if o.closed {
		return newClosedError()
	}
	if o.arena != nil {
		o.arena.resetRows()
		config := o.arena.resumeConfig(&o.option, 1)
		compiler := o.arena.resumeCompiler(config.output)
		compiler.compileRow(row, o.rendered)
		if compiler.err != nil {
			return compiler.err
		}
		solver := o.arena.resumeSolver(compiler.output)
		painter := o.arena.newPainter(solver.output, o.w)
		painter.prepare()
		painter.paintBody()
		o.err = painter.err
		o.rendered++
		return o.err
	}
	o.arena = acquireArena()
	config := o.arena.newConfig(&o.option, 1)
	config.prepare()
	if config.err != nil {
		err := config.err
		o.releaseArena()
		return err
	}
	compiler := o.arena.newCompiler(config.output)
	compiler.prepare()
	compiler.compileHeader()
	compiler.compileRow(row, 0)
	if compiler.err != nil {
		err := compiler.err
		o.releaseArena()
		return err
	}
	solver := o.arena.newSolver(compiler.output)
	solver.prepare()
	painter := o.arena.newPainter(solver.output, o.w)
	painter.prepare()
	painter.paintHeader()
	painter.paintBody()
	o.err = painter.err
	o.rendered = 1
	return o.err
}

// Close finalizes the stream. Repeated calls after a successful close return
// nil.
func (o *Stream) Close() error {
	if o.err != nil {
		o.releaseArena()
		return o.err
	}
	if o.closed {
		return nil
	}
	o.closed = true
	if o.arena != nil {
		o.releaseArena()
		return nil
	}
	o.arena = acquireArena()
	config := o.arena.newConfig(&o.option, 0)
	config.prepare()
	if config.err != nil {
		o.err = config.err
		o.releaseArena()
		return o.err
	}
	compiler := o.arena.newCompiler(config.output)
	compiler.prepare()
	compiler.compileHeader()
	solver := o.arena.newSolver(compiler.output)
	solver.prepare()
	painter := o.arena.newPainter(solver.output, o.w)
	painter.prepare()
	painter.paintHeader()
	o.err = painter.err
	o.releaseArena()
	return o.err
}

// releaseArena returns the stream's active arena to the pool.
func (o *Stream) releaseArena() {
	o.arena.release()
	o.arena = nil
}
