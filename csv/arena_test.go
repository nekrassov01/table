package csv

import (
	"io"
	"testing"

	"github.com/nekrassov01/table/internal/testutil"
	"github.com/nekrassov01/table/internal/value"
)

func Test_arena_resetRows(t *testing.T) {
	type fields struct {
		strings  func() value.Store
		config   configState
		compiler compilerState
		painter  painterState
	}
	type want struct {
		cellLen      int
		rowLen       int
		quoteLen     int
		valueLen     int
		stringLen    int
		columnLen    int
		lineLen      int
		cellCleared  bool
		rowCleared   bool
		valueCleared bool
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "clears row storage and preserves configuration and painter",
			fields: fields{
				strings: func() value.Store {
					var strings value.Store
					strings.AppendString("value")
					return strings
				},
				config: configState{
					columns: []columnConfig{{}},
				},
				compiler: compilerState{
					cells: []cell{
						{
							value: "value",
						},
					},
					rows: []row{
						{
							cells: []cell{{value: "value"}},
						},
					},
					quotes: []byte("quoted"),
					values: []string{"value"},
				},
				painter: painterState{
					line: []byte("line"),
				},
			},
			want: want{
				columnLen:    1,
				lineLen:      4,
				cellCleared:  true,
				rowCleared:   true,
				valueCleared: true,
			},
		},
		{
			name: "accepts empty storage",
			fields: fields{
				strings: func() value.Store {
					return value.Store{}
				},
			},
			want: want{
				cellCleared:  true,
				rowCleared:   true,
				valueCleared: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			a := &arena{
				strings:  test.fields.strings(),
				config:   test.fields.config,
				compiler: test.fields.compiler,
				painter:  test.fields.painter,
			}
			cells := a.compiler.cells
			rows := a.compiler.rows
			values := a.compiler.values
			a.resetRows()
			got := want{
				cellLen:      len(a.compiler.cells),
				rowLen:       len(a.compiler.rows),
				quoteLen:     len(a.compiler.quotes),
				valueLen:     len(a.compiler.values),
				stringLen:    a.strings.Mark(),
				columnLen:    len(a.config.columns),
				lineLen:      len(a.painter.line),
				cellCleared:  len(cells) == 0 || cells[0] == (cell{}),
				rowCleared:   len(rows) == 0 || rows[0].cells == nil,
				valueCleared: len(values) == 0 || values[0] == "",
			}
			testutil.AssertValue(t, got, test.want, "resetRows")
		})
	}
}

func Test_arena_newConfig(t *testing.T) {
	type fields struct {
		state configState
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
		stateBound  bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "binds input and state",
			args: args{
				option: &option{
					header: []string{"header"},
				},
				footer:      [][]string{{"footer"}},
				bodyRows:    3,
				bodyColumns: 2,
			},
			want: want{
				output: configResult{
					option: &option{
						header: []string{"header"},
					},
					header:   []string{"header"},
					footer:   [][]string{{"footer"}},
					bodyRows: 3,
				},
				bodyColumns: 2,
				stateBound:  true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			a := &arena{
				config: test.fields.state,
			}
			configured := a.newConfig(test.args.option, test.args.footer, test.args.bodyRows, test.args.bodyColumns)
			got := want{
				bodyColumns: configured.bodyColumns,
				output:      configured.output,
				stateBound:  configured.state == &a.config,
			}
			testutil.AssertValue(t, got, test.want, "newConfig")
		})
	}
}

func Test_arena_resumeConfig(t *testing.T) {
	type fields struct {
		state configState
	}
	type args struct {
		option   *option
		footer   [][]string
		bodyRows int
	}
	type want struct {
		columnCount   int
		footerColumns int
		bodyRows      int
		stateBound    bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "pairs current input with retained columns",
			fields: fields{
				state: configState{
					columns: []columnConfig{{}, {}},
				},
			},
			args: args{
				option:   &option{},
				footer:   [][]string{{"a"}, {"a", "b", "c"}},
				bodyRows: 1,
			},
			want: want{
				columnCount:   2,
				footerColumns: 3,
				bodyRows:      1,
				stateBound:    true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			a := &arena{
				config: test.fields.state,
			}
			configured := a.resumeConfig(test.args.option, test.args.footer, test.args.bodyRows)
			got := want{
				columnCount:   len(configured.output.columns),
				footerColumns: configured.output.footerColumns,
				bodyRows:      configured.output.bodyRows,
				stateBound:    configured.state == &a.config,
			}
			testutil.AssertValue(t, got, test.want, "resumeConfig")
		})
	}
}

func Test_arena_newCompiler(t *testing.T) {
	type args struct {
		input configResult
	}
	type want struct {
		input       configResult
		stateBound  bool
		storeBound  bool
		bodyStart   int
		outputInput configResult
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "binds fresh compiler",
			args: args{
				input: configResult{
					columns: []columnConfig{{}},
				},
			},
			want: want{
				input: configResult{
					columns: []columnConfig{{}},
				},
				stateBound: true,
				storeBound: true,
				bodyStart:  -1,
				outputInput: configResult{
					columns: []columnConfig{{}},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			a := &arena{}
			compiler := a.newCompiler(test.args.input)
			got := want{
				input:       compiler.input,
				stateBound:  compiler.state == &a.compiler,
				storeBound:  compiler.strings == &a.strings,
				bodyStart:   compiler.bodyStart,
				outputInput: compiler.output.configResult,
			}
			testutil.AssertValue(t, got, test.want, "newCompiler")
		})
	}
}

func Test_arena_resumeCompiler(t *testing.T) {
	type fields struct {
		state compilerState
	}
	type args struct {
		input configResult
	}
	type want struct {
		stateBound bool
		storeBound bool
		bodyStart  int
		rowCap     int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "binds retained compiler state",
			fields: fields{
				state: compilerState{
					rows: make([]row, 0, 3),
				},
			},
			args: args{
				input: configResult{
					columns: []columnConfig{{}},
				},
			},
			want: want{
				stateBound: true,
				storeBound: true,
				bodyStart:  -1,
				rowCap:     3,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			a := &arena{
				compiler: test.fields.state,
			}
			compiler := a.resumeCompiler(test.args.input)
			got := want{
				stateBound: compiler.state == &a.compiler,
				storeBound: compiler.strings == &a.strings,
				bodyStart:  compiler.bodyStart,
				rowCap:     cap(compiler.state.rows),
			}
			testutil.AssertValue(t, got, test.want, "resumeCompiler")
		})
	}
}

func Test_arena_newPainter(t *testing.T) {
	type args struct {
		input compilerResult
		w     io.Writer
	}
	type want struct {
		input      compilerResult
		stateBound bool
		writer     io.Writer
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "binds compiled input and writer",
			args: args{
				input: compilerResult{
					header: row{
						cells: []cell{{value: "A"}},
					},
				},
				w: io.Discard,
			},
			want: want{
				input: compilerResult{
					header: row{
						cells: []cell{{value: "A"}},
					},
				},
				stateBound: true,
				writer:     io.Discard,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			a := &arena{}
			painter := a.newPainter(test.args.input, test.args.w)
			got := want{
				input:      painter.input,
				stateBound: painter.state == &a.painter,
				writer:     painter.w,
			}
			testutil.AssertValue(t, got, test.want, "newPainter")
		})
	}
}

func Test_arena_release(t *testing.T) {
	type fields struct {
		arena func() *arena
	}
	type want struct {
		lineNil        bool
		lineBackingCap int
		columnCleared  bool
		cellCleared    bool
		rowCleared     bool
		valueCleared   bool
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "nil receiver",
			fields: fields{
				arena: func() *arena {
					return nil
				},
			},
			want: want{
				lineNil:       true,
				columnCleared: true,
				cellCleared:   true,
				rowCleared:    true,
				valueCleared:  true,
			},
		},
		{
			name: "clears references and keeps largest line backing",
			fields: fields{
				arena: func() *arena {
					return &arena{
						config: configState{
							columns: []columnConfig{
								{
									transformer: func(any) string {
										return "value"
									},
								},
							},
						},
						compiler: compilerState{
							cells:  []cell{{value: "value"}},
							rows:   []row{{cells: []cell{{value: "value"}}}},
							values: []string{"value"},
						},
						painter: painterState{
							lineBacking: make([]byte, 0, 1),
							line:        make([]byte, 0, 8),
						},
					}
				},
			},
			want: want{
				lineNil:        true,
				lineBackingCap: 8,
				columnCleared:  true,
				cellCleared:    true,
				rowCleared:     true,
				valueCleared:   true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			a := test.fields.arena()
			a.release()
			got := want{
				lineNil:       a == nil || a.painter.line == nil,
				columnCleared: a == nil || len(a.config.columns) == 0 || a.config.columns[0].transformer == nil,
				cellCleared:   a == nil || len(a.compiler.cells) == 0 || a.compiler.cells[0] == (cell{}),
				rowCleared:    a == nil || len(a.compiler.rows) == 0 || a.compiler.rows[0].cells == nil,
				valueCleared:  a == nil || len(a.compiler.values) == 0 || a.compiler.values[0] == "",
			}
			if a != nil {
				got.lineBackingCap = cap(a.painter.lineBacking)
			}
			testutil.AssertValue(t, got, test.want, "release")
		})
	}
}

func Test_acquireArena(t *testing.T) {
	type want struct {
		stringLen int
		nonNil    bool
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "resets shared string store",
			want: want{
				nonNil: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			seed := acquireArena()
			seed.strings.AppendString("value")
			seed.release()
			a := acquireArena()
			t.Cleanup(a.release)
			got := want{
				stringLen: a.strings.Mark(),
				nonNil:    a != nil,
			}
			testutil.AssertValue(t, got, test.want, "acquireArena")
		})
	}
}
