package backlog

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
		compiler compilerState
		solver   solverState
	}
	type want struct {
		stringMark int
		cells      int
		rows       int
		escapes    int
		values     int
		rowspans   uint64
		colspans   uint64
		metrics    []columnMetric
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "clears row views and preserves stream state",
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
						},
					},
					rows:    []row{{}},
					escapes: []byte("value"),
					values:  []string{"value"},
					rowspans: func() scope.Masks {
						var masks scope.Masks
						masks.Mark(ScopeBody, 0)
						return masks
					}(),
					colspans: func() scope.Masks {
						var masks scope.Masks
						masks.Mark(ScopeBody, 1)
						return masks
					}(),
				},
				solver: solverState{
					columnMetrics: []columnMetric{
						{
							box: box{width: 5},
						},
					},
				},
			},
			want: want{
				rowspans: 1,
				colspans: 2,
				metrics: []columnMetric{
					{
						box: box{width: 5},
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &arena{
				strings:  test.fields.strings,
				compiler: test.fields.compiler,
				solver:   test.fields.solver,
			}
			o.resetRows()
			got := want{
				stringMark: o.strings.Mark(),
				cells:      len(o.compiler.cells),
				rows:       len(o.compiler.rows),
				escapes:    len(o.compiler.escapes),
				values:     len(o.compiler.values),
				rowspans:   o.compiler.rowspans.Resolve(ScopeBody),
				colspans:   o.compiler.colspans.Resolve(ScopeBody),
				metrics:    o.solver.columnMetrics,
			}
			testutil.AssertValue(t, got, test.want, "resetRows")
		})
	}
}

func Test_arena_newConfig(t *testing.T) {
	type fields struct {
		config configState
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
					header: [][]string{{"a"}},
				},
				footer:      [][]string{{"footer"}},
				bodyRows:    2,
				bodyColumns: 3,
			},
			want: want{
				output: configResult{
					option: &option{
						header: [][]string{{"a"}},
					},
					header:   [][]string{{"a"}},
					footer:   [][]string{{"footer"}},
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
				config: test.fields.config,
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
		config configState
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
							rowspan: ScopeBody,
						},
					},
				},
			},
			args: args{
				option: &option{
					placeholder: "-",
					header:      [][]string{{"header"}},
				},
				footer:   [][]string{{"sum", "total"}},
				bodyRows: 2,
			},
			want: want{
				output: configResult{
					option: &option{
						placeholder: "-",
						header:      [][]string{{"header"}},
					},
					header:   [][]string{{"header"}},
					footer:   [][]string{{"sum", "total"}},
					bodyRows: 2,
					columns: []columnConfig{
						{
							rowspan: ScopeBody,
						},
					},
					footerColumns: 2,
				},
				ownsState: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &arena{
				config: test.fields.config,
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
		compiler compilerState
	}
	type args struct {
		input configResult
	}
	type want struct {
		input       configResult
		output      compilerResult
		ownsState   bool
		ownsStrings bool
		bodyStart   int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "binds compiler state",
			args: args{
				input: configResult{
					option: &option{},
				},
			},
			want: want{
				input: configResult{
					option: &option{},
				},
				output: compilerResult{
					configResult: configResult{
						option: &option{},
					},
				},
				ownsState:   true,
				ownsStrings: true,
				bodyStart:   -1,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &arena{
				compiler: test.fields.compiler,
			}
			compiler := o.newCompiler(test.args.input)
			got := want{
				input:       compiler.input,
				output:      compiler.output,
				ownsState:   compiler.state == &o.compiler,
				ownsStrings: compiler.strings == &o.strings,
				bodyStart:   compiler.bodyStart,
			}
			testutil.AssertValue(t, got, test.want, "newCompiler")
		})
	}
}

func Test_arena_resumeCompiler(t *testing.T) {
	type fields struct {
		compiler compilerState
	}
	type args struct {
		input configResult
	}
	type want struct {
		input       configResult
		output      compilerResult
		ownsState   bool
		ownsStrings bool
		bodyStart   int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "binds retained compiler state",
			args: args{
				input: configResult{
					option: &option{},
				},
			},
			want: want{
				input: configResult{
					option: &option{},
				},
				output: compilerResult{
					configResult: configResult{
						option: &option{},
					},
				},
				ownsState:   true,
				ownsStrings: true,
				bodyStart:   -1,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &arena{
				compiler: test.fields.compiler,
			}
			compiler := o.resumeCompiler(test.args.input)
			got := want{
				input:       compiler.input,
				output:      compiler.output,
				ownsState:   compiler.state == &o.compiler,
				ownsStrings: compiler.strings == &o.strings,
				bodyStart:   compiler.bodyStart,
			}
			testutil.AssertValue(t, got, test.want, "resumeCompiler")
		})
	}
}

func Test_arena_newSolver(t *testing.T) {
	type fields struct {
		solver solverState
	}
	type args struct {
		input compilerResult
	}
	type want struct {
		input     compilerResult
		ownsState bool
		output    solverResult
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "binds solver state",
			args: args{
				input: compilerResult{
					configResult: configResult{
						option: &option{},
					},
				},
			},
			want: want{
				input: compilerResult{
					configResult: configResult{
						option: &option{},
					},
				},
				ownsState: true,
				output: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: &option{},
						},
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &arena{
				solver: test.fields.solver,
			}
			solver := o.newSolver(test.args.input)
			got := want{
				input:     solver.input,
				ownsState: solver.state == &o.solver,
				output:    solver.output,
			}
			testutil.AssertValue(t, got, test.want, "newSolver")
		})
	}
}

func Test_arena_resumeSolver(t *testing.T) {
	type fields struct {
		solver solverState
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
			name: "restores measured metrics",
			fields: fields{
				solver: solverState{
					columnMetrics: []columnMetric{
						{
							box: box{width: 5},
						},
					},
				},
			},
			args: args{
				input: compilerResult{
					body: []row{
						{
							cells: []cell{{value: "value"}},
						},
					},
				},
			},
			want: want{
				output: solverResult{
					compilerResult: compilerResult{
						body: []row{
							{
								cells: []cell{{value: "value"}},
							},
						},
					},
					metrics: []columnMetric{
						{
							box: box{width: 5},
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
				solver: test.fields.solver,
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
		painter painterState
	}
	type args struct {
		input solverResult
		w     *bytes.Buffer
	}
	type want struct {
		input      solverResult
		ownsState  bool
		ownsWriter bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "binds painter state and writer",
			args: args{
				input: solverResult{
					metrics: []columnMetric{{}},
				},
				w: &bytes.Buffer{},
			},
			want: want{
				input: solverResult{
					metrics: []columnMetric{{}},
				},
				ownsState:  true,
				ownsWriter: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &arena{
				painter: test.fields.painter,
			}
			painter := o.newPainter(test.args.input, test.args.w)
			got := want{
				input:      painter.input,
				ownsState:  painter.state == &o.painter,
				ownsWriter: painter.w == test.args.w,
			}
			testutil.AssertValue(t, got, test.want, "newPainter")
		})
	}
}

func Test_arena_release(t *testing.T) {
	type fields struct {
		arena *arena
	}
	type want struct {
		lineIsNil bool
		lineCap   int
		columns   []columnConfig
		cells     []cell
		values    []string
		bodyReset bool
		bandReset bool
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "nil arena",
			want: want{
				lineIsNil: true,
			},
		},
		{
			name: "clears references and adopts grown line",
			fields: fields{
				arena: &arena{
					config: configState{
						columns: []columnConfig{{}},
					},
					compiler: compilerState{
						cells: []cell{
							{
								value: "value",
							},
						},
						values: []string{"value"},
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
						lineBacking: make([]byte, 0, 1),
						line:        make([]byte, 0, 8),
					},
				},
			},
			want: want{
				lineIsNil: true,
				lineCap:   8,
				columns:   []columnConfig{{}},
				cells:     []cell{{}},
				values:    []string{""},
				bodyReset: true,
				bandReset: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := test.fields.arena
			o.release()
			got := want{
				lineIsNil: true,
			}
			if o != nil {
				got.lineIsNil = o.painter.line == nil
				got.lineCap = cap(o.painter.lineBacking)
				got.columns = o.config.columns
				got.cells = o.compiler.cells
				got.values = o.compiler.values
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
			name: "returns reset arena",
			want: want{
				nonNil:     true,
				stringMark: 0,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := acquireArena()
			got := want{
				nonNil:     o != nil,
				stringMark: o.strings.Mark(),
			}
			o.release()
			testutil.AssertValue(t, got, test.want, "acquireArena")
		})
	}
}
