package csv

import (
	"io"
	"sync"

	"github.com/nekrassov01/table/internal/column"
	"github.com/nekrassov01/table/internal/value"
)

// pool shares reusable pipeline storage across tables and streams.
var pool = sync.Pool{
	New: func() any {
		return new(arena)
	},
}

// arena owns reusable storage for one table render or active stream.
type arena struct {
	strings  value.Store   // Text backing shared across pipeline stages.
	config   configState   // Storage retained by the config.
	compiler compilerState // Storage retained by the compiler.
	painter  painterState  // Storage retained by the painter.
}

// resetRows clears data owned by the previous stream pass while preserving
// resolved columns and reusable capacity.
func (o *arena) resetRows() {
	compiler := &o.compiler
	clear(compiler.cells)
	clear(compiler.rows)
	clear(compiler.values)
	compiler.cells = compiler.cells[:0]
	compiler.rows = compiler.rows[:0]
	compiler.quotes = compiler.quotes[:0]
	compiler.values = compiler.values[:0]
	o.strings.Reset()
}

// newConfig starts the initial config stage for one rendering pass.
func (o *arena) newConfig(option *option, footer [][]string, bodyRows, bodyColumns int) config {
	return config{
		bodyColumns: bodyColumns,
		state:       &o.config,
		output: configResult{
			option:   option,
			header:   option.header,
			footer:   footer,
			bodyRows: bodyRows,
		},
	}
}

// resumeConfig pairs pass data with the columns retained by an active stream.
func (o *arena) resumeConfig(option *option, footer [][]string, bodyRows int) config {
	state := &o.config
	return config{
		state: state,
		output: configResult{
			option:        option,
			header:        option.header,
			footer:        footer,
			bodyRows:      bodyRows,
			columns:       state.columns,
			footerColumns: column.MaxColumns(footer),
		},
	}
}

// newCompiler starts the initial compiler stage for input.
func (o *arena) newCompiler(input configResult) compiler {
	return compiler{
		input:     input,
		state:     &o.compiler,
		strings:   &o.strings,
		bodyStart: -1,
		output: compilerResult{
			configResult: input,
		},
	}
}

// resumeCompiler continues compilation with storage retained from prior rows.
func (o *arena) resumeCompiler(input configResult) compiler {
	return compiler{
		input:     input,
		state:     &o.compiler,
		strings:   &o.strings,
		bodyStart: -1,
		output: compilerResult{
			configResult: input,
		},
	}
}

// newPainter starts a painter for compiled rows.
func (o *arena) newPainter(input compilerResult, w io.Writer) painter {
	return painter{
		input: input,
		state: &o.painter,
		w:     w,
	}
}

// release clears retained references and returns o to the pool while preserving
// reusable backing storage.
func (o *arena) release() {
	if o == nil {
		return
	}
	compiler := &o.compiler
	painter := &o.painter
	if cap(painter.line) > cap(painter.lineBacking) {
		painter.lineBacking = painter.line[:0]
	}
	clear(o.config.columns)
	clear(compiler.cells)
	clear(compiler.rows)
	clear(compiler.values)
	painter.line = nil
	pool.Put(o)
}

// configState retains resolved columns across stream passes.
type configState struct {
	columns []columnConfig // Resolved column settings in logical order.
}

// compilerState owns reusable compilation storage.
type compilerState struct {
	cells  []cell   // Backing for compiled logical cells.
	rows   []row    // Backing for compiled logical rows.
	quotes []byte   // Backing for quoted values.
	values []string // Resolved values for the current row.
}

// painterState owns the reusable output buffer.
type painterState struct {
	lineBacking []byte // Backing for the output line.
	line        []byte // Current output line.
}

// acquireArena takes an arena from the pool and resets its shared string store.
func acquireArena() *arena {
	a := pool.Get().(*arena)
	a.strings.Reset()
	return a
}
