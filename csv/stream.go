package csv

import (
	"io"

	"github.com/nekrassov01/table"
)

var _ table.Streamer = (*Stream)(nil)

// Stream renders rows incrementally as delimiter-separated records.
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

// Render writes one body record. The first row that establishes columns writes
// the configured header before the body.
func (o *Stream) Render(row []any) error {
	if o.err != nil {
		return o.err
	}
	if o.closed {
		return newClosedError()
	}
	if o.arena != nil {
		o.arena.resetRows()
		config := o.arena.resumeConfig(&o.option, nil, 1)
		compiler := o.arena.resumeCompiler(config.output)
		compiler.compileRow(row, o.rendered)
		if compiler.err != nil {
			return compiler.err
		}
		painter := o.arena.newPainter(compiler.output, o.w)
		painter.prepare()
		painter.paintBody()
		o.err = painter.err
		o.rendered++
		return o.err
	}
	o.arena = acquireArena()
	config := o.arena.newConfig(&o.option, nil, 1, len(row))
	config.prepare()
	if config.err != nil {
		err := config.err
		o.releaseArena()
		return err
	}
	configured := config.output
	if len(configured.columns) == 0 {
		o.releaseArena()
		return nil
	}
	compiler := o.arena.newCompiler(configured)
	compiler.prepare()
	compiler.compileHeader()
	compiler.compileRow(row, 0)
	if compiler.err != nil {
		err := compiler.err
		o.releaseArena()
		return err
	}
	painter := o.arena.newPainter(compiler.output, o.w)
	painter.prepare()
	painter.paintHeader()
	painter.paintBody()
	o.err = painter.err
	o.rendered = 1
	return o.err
}

// Close writes the footer and completes the output. Repeated calls after a
// successful close return nil.
func (o *Stream) Close() error {
	if o.err != nil {
		o.releaseArena()
		return o.err
	}
	if o.closed {
		return nil
	}
	o.closed = true
	var footer [][]string
	if o.option.footer != nil {
		footer = o.option.footer()
	}
	if o.arena != nil {
		o.arena.resetRows()
		config := o.arena.resumeConfig(&o.option, footer, 0)
		compiler := o.arena.resumeCompiler(config.output)
		compiler.compileFooter()
		if compiler.err != nil {
			o.err = compiler.err
			o.releaseArena()
			return o.err
		}
		painter := o.arena.newPainter(compiler.output, o.w)
		painter.prepare()
		painter.paintFooter()
		o.err = painter.err
		o.releaseArena()
		return o.err
	}
	if len(o.option.header) == 0 && len(footer) == 0 {
		return nil
	}
	o.arena = acquireArena()
	config := o.arena.newConfig(&o.option, footer, 0, 0)
	config.prepare()
	if config.err != nil {
		o.err = config.err
		o.releaseArena()
		return o.err
	}
	configured := config.output
	if len(configured.columns) == 0 {
		o.releaseArena()
		return nil
	}
	compiler := o.arena.newCompiler(configured)
	compiler.prepare()
	compiler.compileHeader()
	compiler.compileFooter()
	if compiler.err != nil {
		o.err = compiler.err
		o.releaseArena()
		return o.err
	}
	painter := o.arena.newPainter(compiler.output, o.w)
	painter.prepare()
	painter.paintHeader()
	painter.paintFooter()
	o.err = painter.err
	o.releaseArena()
	return o.err
}

// releaseArena returns the stream's active arena to the pool.
func (o *Stream) releaseArena() {
	o.arena.release()
	o.arena = nil
}
