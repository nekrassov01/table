package html

import (
	"bytes"
	"testing"

	"github.com/nekrassov01/table/internal/span"
	"github.com/nekrassov01/table/internal/testutil"
	"github.com/nekrassov01/table/internal/value"
)

func Test_arena_resetRows(t *testing.T) {
	type fields struct {
		strings  value.Store
		compiler compilerState
	}
	type want struct {
		stringMark  int
		cells       int
		rows        int
		escapes     int
		values      int
		columnSizes []int
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
					rows:        []row{{}, {}, {}},
					escapes:     []byte("value"),
					values:      []string{"value"},
					columnSizes: []int{5},
				},
			},
			want: want{
				columnSizes: []int{5},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &arena{
				strings:  test.fields.strings,
				compiler: test.fields.compiler,
			}
			o.resetRows()
			got := want{
				stringMark:  o.strings.Mark(),
				cells:       len(o.compiler.cells),
				rows:        len(o.compiler.rows),
				escapes:     len(o.compiler.escapes),
				values:      len(o.compiler.values),
				columnSizes: o.compiler.columnSizes,
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
			name: "binds arena state",
			args: args{
				option: &option{
					placeholder: "-",
					header:      [][]string{{"header"}},
				},
				footer:      [][]string{{"footer"}},
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
					columns: []column{
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
					columns: []column{
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
		ownsState   bool
		ownsStrings bool
		bodyStart   int
		output      compilerResult
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
				ownsState:   true,
				ownsStrings: true,
				bodyStart:   -1,
				output: compilerResult{
					configResult: configResult{
						option: &option{},
					},
				},
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
				ownsState:   compiler.state == &o.compiler,
				ownsStrings: compiler.strings == &o.strings,
				bodyStart:   compiler.bodyStart,
				output:      compiler.output,
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
			name: "relays retained column sizes",
			fields: fields{
				compiler: compilerState{
					columnSizes: []int{3, 5},
				},
			},
			args: args{
				input: configResult{
					option: &option{},
				},
			},
			want: want{
				output: compilerResult{
					configResult: configResult{
						option: &option{},
					},
					columnSizes:     []int{3, 5},
					hasPreviousBody: true,
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
			name: "binds solver state and relays compiler output",
			args: args{
				input: compilerResult{
					caption: "caption",
				},
			},
			want: want{
				input: compilerResult{
					caption: "caption",
				},
				output: solverResult{
					compilerResult: compilerResult{
						caption: "caption",
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
			solver := o.newSolver(test.args.input)
			got := want{
				input:     solver.input,
				output:    solver.output,
				ownsState: solver.state == &o.solver,
			}
			testutil.AssertValue(t, got, test.want, "newSolver")
		})
	}
}

func Test_arena_newPainter(t *testing.T) {
	type fields struct {
		painter painterState
	}
	type args struct {
		input solverResult
	}
	type want struct {
		input     solverResult
		ownsState bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "binds painter state",
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
				ownsState: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &arena{
				painter: test.fields.painter,
			}
			w := &bytes.Buffer{}
			painter := o.newPainter(test.args.input, w)
			got := want{
				input:     painter.input,
				ownsState: painter.state == &o.painter,
			}
			testutil.AssertValue(t, got, test.want, "newPainter")
		})
	}
}

func Test_arena_release_nil(t *testing.T) {
	type fields struct {
		arena *arena
	}
	type want struct {
		panics bool
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "nil arena",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := want{}
			test.fields.arena.release()
			testutil.AssertValue(t, got, test.want, "release")
		})
	}
}

func Test_arena_release(t *testing.T) {
	type fields struct {
		config   configState
		compiler compilerState
		painter  painterState
	}
	type want struct {
		lineNil     bool
		lineBacking int
		columnZero  bool
		cellZero    bool
		bodyReset   bool
		bandReset   bool
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "severs views and retains grown line backing",
			fields: fields{
				config: configState{
					columns: []column{
						{
							transformer: transformer{
								fn: func(any) (string, *Color, *Decoration) {
									return "", nil, nil
								},
							},
						},
					},
				},
				compiler: compilerState{
					cells: []cell{
						{
							value: "value",
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
					lineBacking: make([]byte, 0, 2),
					line:        make([]byte, 0, 8),
				},
			},
			want: want{
				lineNil:     true,
				lineBacking: 8,
				columnZero:  true,
				cellZero:    true,
				bodyReset:   true,
				bandReset:   true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &arena{
				config:   test.fields.config,
				compiler: test.fields.compiler,
				painter:  test.fields.painter,
			}
			o.release()
			got := want{
				lineNil:     o.painter.line == nil,
				lineBacking: cap(o.painter.lineBacking),
				columnZero:  o.config.columns[0].transformer.fn == nil,
				cellZero:    o.compiler.cells[0].value == "",
				bodyReset:   span.Rowspans(1, []string{"value"}, &o.compiler.previousBody) == 0,
				bandReset:   span.Rowspans(1, []string{"value"}, &o.compiler.previousBand) == 0,
			}
			o.compiler.previousBody.Clear()
			o.compiler.previousBand.Clear()
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
