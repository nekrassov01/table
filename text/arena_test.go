package text

import (
	"bytes"
	"testing"

	"github.com/nekrassov01/table/internal/scope"
	"github.com/nekrassov01/table/internal/span"
	"github.com/nekrassov01/table/internal/testutil"
	"github.com/nekrassov01/table/internal/value"
)

func Test_arena_resetRows(t *testing.T) {
	type fields struct {
		strings  value.Store
		config   configState
		compiler compilerState
		solver   solverState
		painter  painterState
	}
	type want struct {
		stringMark int
		cells      int
		spanValues int
		rows       int
		layouts    int
		segments   int
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "clears per-row views",
			fields: fields{
				strings: func() value.Store {
					var strings value.Store
					strings.AppendString("value")
					return strings
				}(),
				compiler: compilerState{
					cells: []cell{
						{
							value: "value",
							attr: &Attr{
								Prefix: []byte("prefix"),
							},
						},
					},
					spanValues: []string{"value"},
					rows: []row{
						{
							cells: []cell{
								{
									value: "value",
								},
							},
						},
					},
				},
				painter: painterState{
					layouts: []layout{
						{
							value: "value",
						},
					},
					segments: []segment{
						{
							value: "value",
						},
					},
				},
			},
			want: want{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &arena{
				strings:  test.fields.strings,
				config:   test.fields.config,
				compiler: test.fields.compiler,
				solver:   test.fields.solver,
				painter:  test.fields.painter,
			}
			o.resetRows()
			got := want{
				stringMark: o.strings.Mark(),
				cells:      len(o.compiler.cells),
				spanValues: len(o.compiler.spanValues),
				rows:       len(o.compiler.rows),
				layouts:    len(o.painter.layouts),
				segments:   len(o.painter.segments),
			}
			testutil.AssertValue(t, got, test.want, "resetRows")
		})
	}
}

func Test_arena_newConfig(t *testing.T) {
	type fields struct {
		strings  value.Store
		config   configState
		compiler compilerState
		solver   solverState
		painter  painterState
	}
	type args struct {
		option      *option
		footer      [][]string
		bodyRows    int
		bodyColumns int
	}
	type want struct {
		bodyColumns int
		output      configResult
		ownsState   bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "binds config state",
			args: args{
				option: &option{
					placeholder: "-",
					header:      [][]string{{"header"}},
				},
				footer:      [][]string{{"total"}},
				bodyRows:    2,
				bodyColumns: 3,
			},
			want: want{
				output: configResult{
					option: &option{
						placeholder: "-",
						header:      [][]string{{"header"}},
					},
					header:   [][]string{{"header"}},
					footer:   [][]string{{"total"}},
					bodyRows: 2,
				},
				bodyColumns: 3,
				ownsState:   true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &arena{
				strings:  test.fields.strings,
				config:   test.fields.config,
				compiler: test.fields.compiler,
				solver:   test.fields.solver,
				painter:  test.fields.painter,
			}
			config := o.newConfig(test.args.option, test.args.footer, test.args.bodyRows, test.args.bodyColumns)
			got := want{
				bodyColumns: config.bodyColumns,
				output:      config.output,
				ownsState:   config.state == &o.config,
			}
			testutil.AssertValue(t, got, test.want, "newConfig")
		})
	}
}

func Test_arena_resumeConfig(t *testing.T) {
	type fields struct {
		strings  value.Store
		config   configState
		compiler compilerState
		solver   solverState
		painter  painterState
	}
	type args struct {
		option   *option
		footer   [][]string
		bodyRows int
	}
	type want struct {
		output    configResult
		ownsState bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "restores resolved columns",
			fields: fields{
				config: configState{
					columns: []columnConfig{
						{
							limit: 3,
						},
					},
				},
			},
			args: args{
				option: &option{
					placeholder: "-",
					header:      [][]string{{"header"}},
				},
				footer:   [][]string{{"total"}},
				bodyRows: 2,
			},
			want: want{
				output: configResult{
					option: &option{
						placeholder: "-",
						header:      [][]string{{"header"}},
					},
					header:   [][]string{{"header"}},
					footer:   [][]string{{"total"}},
					bodyRows: 2,
					columns: []columnConfig{
						{
							limit: 3,
						},
					},
					footerColumns: 1,
				},
				ownsState: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &arena{
				strings:  test.fields.strings,
				config:   test.fields.config,
				compiler: test.fields.compiler,
				solver:   test.fields.solver,
				painter:  test.fields.painter,
			}
			config := o.resumeConfig(test.args.option, test.args.footer, test.args.bodyRows)
			got := want{
				output:    config.output,
				ownsState: config.state == &o.config,
			}
			testutil.AssertValue(t, got, test.want, "resumeConfig")
		})
	}
}

func Test_arena_newCompiler(t *testing.T) {
	type fields struct {
		strings  value.Store
		config   configState
		compiler compilerState
		solver   solverState
		painter  painterState
	}
	type args struct {
		input configResult
	}
	type want struct {
		bodyStart   int
		output      compilerResult
		ownsState   bool
		ownsStrings bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "starts a fresh compilation",
			args: args{
				input: configResult{
					option: &option{},
					columns: []columnConfig{
						{},
					},
				},
			},
			want: want{
				bodyStart: -1,
				output: compilerResult{
					configResult: configResult{
						option: &option{},
						columns: []columnConfig{
							{},
						},
					},
					previousBars: allBars,
					lastBars:     allBars,
				},
				ownsState:   true,
				ownsStrings: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &arena{
				strings:  test.fields.strings,
				config:   test.fields.config,
				compiler: test.fields.compiler,
				solver:   test.fields.solver,
				painter:  test.fields.painter,
			}
			compiler := o.newCompiler(test.args.input)
			got := want{
				bodyStart:   compiler.bodyStart,
				output:      compiler.output,
				ownsState:   compiler.state == &o.compiler,
				ownsStrings: compiler.strings == &o.strings,
			}
			testutil.AssertValue(t, got, test.want, "newCompiler")
		})
	}
}

func Test_arena_resumeCompiler(t *testing.T) {
	type fields struct {
		strings  value.Store
		config   configState
		compiler compilerState
		solver   solverState
		painter  painterState
	}
	type args struct {
		input configResult
	}
	type want struct {
		output      compilerResult
		ownsState   bool
		ownsStrings bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "restores body span state",
			fields: fields{
				compiler: func() compilerState {
					var rowspans scope.Masks
					rowspans.Mark(ScopeBody, 1)
					return compilerState{
						rowspans: rowspans,
						lastBars: 0b101,
					}
				}(),
			},
			args: args{
				input: configResult{
					option: &option{},
					columns: []columnConfig{
						{},
						{},
					},
				},
			},
			want: want{
				output: compilerResult{
					configResult: configResult{
						option: &option{},
						columns: []columnConfig{
							{},
							{},
						},
					},
					rowspanMask:     0b10,
					previousBars:    0b101,
					lastBars:        0b101,
					hasPreviousBody: true,
				},
				ownsState:   true,
				ownsStrings: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &arena{
				strings:  test.fields.strings,
				config:   test.fields.config,
				compiler: test.fields.compiler,
				solver:   test.fields.solver,
				painter:  test.fields.painter,
			}
			compiler := o.resumeCompiler(test.args.input)
			got := want{
				output:      compiler.output,
				ownsState:   compiler.state == &o.compiler,
				ownsStrings: compiler.strings == &o.strings,
			}
			testutil.AssertValue(t, got, test.want, "resumeCompiler")
		})
	}
}

func Test_arena_newSolver(t *testing.T) {
	type fields struct {
		strings  value.Store
		config   configState
		compiler compilerState
		solver   solverState
		painter  painterState
	}
	type args struct {
		input compilerResult
	}
	type want struct {
		output     solverResult
		widthLimit int
		ownsState  bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "binds fresh geometry",
			fields: fields{
				solver: solverState{
					columnMetrics: []columnMetric{
						{
							limit: 3,
						},
					},
				},
			},
			args: args{
				input: compilerResult{
					configResult: configResult{
						option: &option{},
					},
				},
			},
			want: want{
				output: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: &option{},
						},
					},
					metrics: []columnMetric{
						{
							limit: 3,
						},
					},
				},
				ownsState: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &arena{
				strings:  test.fields.strings,
				config:   test.fields.config,
				compiler: test.fields.compiler,
				solver:   test.fields.solver,
				painter:  test.fields.painter,
			}
			w := &bytes.Buffer{}
			solver := o.newSolver(test.args.input, w)
			got := want{
				output:     solver.output,
				widthLimit: solver.widthLimit,
				ownsState:  solver.state == &o.solver,
			}
			testutil.AssertValue(t, got, test.want, "newSolver")
		})
	}
}

func Test_arena_resumeSolver(t *testing.T) {
	type fields struct {
		strings  value.Store
		config   configState
		compiler compilerState
		solver   solverState
		painter  painterState
	}
	type args struct {
		input compilerResult
	}
	type want struct {
		output    solverResult
		ownsState bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "restores solved geometry",
			fields: fields{
				solver: solverState{
					columnMetrics: []columnMetric{
						{
							limit: 3,
						},
					},
				},
			},
			args: args{
				input: compilerResult{
					configResult: configResult{
						option: &option{},
					},
				},
			},
			want: want{
				output: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: &option{},
						},
					},
					metrics: []columnMetric{
						{
							limit: 3,
						},
					},
				},
				ownsState: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &arena{
				strings:  test.fields.strings,
				config:   test.fields.config,
				compiler: test.fields.compiler,
				solver:   test.fields.solver,
				painter:  test.fields.painter,
			}
			solver := o.resumeSolver(test.args.input)
			got := want{
				output:    solver.output,
				ownsState: solver.state == &o.solver,
			}
			testutil.AssertValue(t, got, test.want, "resumeSolver")
		})
	}
}

func Test_arena_newPainter(t *testing.T) {
	type fields struct {
		strings  value.Store
		config   configState
		compiler compilerState
		solver   solverState
		painter  painterState
	}
	type args struct {
		input solverResult
	}
	type want struct {
		input       solverResult
		ownsState   bool
		ownsStrings bool
		ownsWriter  bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "binds painting resources",
			args: args{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: &option{},
						},
					},
				},
			},
			want: want{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: &option{},
						},
					},
				},
				ownsState:   true,
				ownsStrings: true,
				ownsWriter:  true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &arena{
				strings:  test.fields.strings,
				config:   test.fields.config,
				compiler: test.fields.compiler,
				solver:   test.fields.solver,
				painter:  test.fields.painter,
			}
			w := &bytes.Buffer{}
			painter := o.newPainter(test.args.input, w)
			got := want{
				input:       painter.input,
				ownsState:   painter.state == &o.painter,
				ownsStrings: painter.strings == &o.strings,
				ownsWriter:  painter.w == w,
			}
			testutil.AssertValue(t, got, test.want, "newPainter")
		})
	}
}

func Test_arena_release(t *testing.T) {
	type fields struct {
		strings  value.Store
		config   configState
		compiler compilerState
		solver   solverState
		painter  painterState
	}
	type want struct {
		nilArena   bool
		columnNil  bool
		cellNil    bool
		valueNil   bool
		rowNil     bool
		layoutNil  bool
		segmentNil bool
		lineNil    bool
		horizonNil bool
		bodyReset  bool
		bandReset  bool
	}
	tests := []struct {
		name   string
		fields *fields
		want   want
	}{
		{
			name: "nil arena",
			want: want{
				nilArena: true,
			},
		},
		{
			name: "drops retained views",
			fields: &fields{
				config: configState{
					columns: []columnConfig{
						{
							transformer: transformer{
								fn: func(any) (string, *Attr) {
									return "", nil
								},
							},
						},
					},
				},
				compiler: compilerState{
					cells: []cell{
						{
							value: "value",
							attr: &Attr{
								Prefix: []byte("prefix"),
							},
						},
					},
					spanValues: []string{"value"},
					rows: []row{
						{
							cells: []cell{
								{
									value: "value",
								},
							},
						},
					},
					previousBody: func() span.PreviousRow {
						var previous span.PreviousRow
						span.Rowspans(1, []string{"value"}, &previous)
						return previous
					}(),
					previousBand: func() span.PreviousRow {
						var previous span.PreviousRow
						span.Rowspans(1, []string{"value"}, &previous)
						return previous
					}(),
				},
				painter: painterState{
					layouts: []layout{
						{
							value: "value",
						},
					},
					segments: []segment{
						{
							value: "value",
						},
					},
					lineBacking: make([]byte, 0, 8),
					line:        []byte("line"),
					horizon:     []byte("horizon"),
				},
			},
			want: want{
				columnNil:  true,
				cellNil:    true,
				valueNil:   true,
				rowNil:     true,
				layoutNil:  true,
				segmentNil: true,
				lineNil:    true,
				horizonNil: true,
				bodyReset:  true,
				bandReset:  true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var o *arena
			if test.fields != nil {
				o = &arena{
					strings:  test.fields.strings,
					config:   test.fields.config,
					compiler: test.fields.compiler,
					solver:   test.fields.solver,
					painter:  test.fields.painter,
				}
			}
			o.release()
			got := want{
				nilArena: o == nil,
			}
			if o != nil {
				got.columnNil = o.config.columns[0].transformer.fn == nil
				got.cellNil = o.compiler.cells[0].value == "" && o.compiler.cells[0].attr == nil
				got.valueNil = o.compiler.spanValues[0] == ""
				got.rowNil = o.compiler.rows[0].cells == nil
				got.layoutNil = o.painter.layouts[0].value == ""
				got.segmentNil = o.painter.segments[0].value == ""
				got.lineNil = o.painter.line == nil
				got.horizonNil = o.painter.horizon == nil
				got.bodyReset = span.Rowspans(1, []string{"value"}, &o.compiler.previousBody) == 0
				got.bandReset = span.Rowspans(1, []string{"value"}, &o.compiler.previousBand) == 0
				o.compiler.previousBody.Clear()
				o.compiler.previousBand.Clear()
			}
			testutil.AssertValue(t, got, test.want, "release")
		})
	}
}

func Test_acquireArena(t *testing.T) {
	type want struct {
		nonNil     bool
		stringMark int
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "returns an empty string store",
			want: want{
				nonNil: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotArena := acquireArena()
			got := want{
				nonNil:     gotArena != nil,
				stringMark: gotArena.strings.Mark(),
			}
			gotArena.release()
			testutil.AssertValue(t, got, test.want, "acquireArena")
		})
	}
}
