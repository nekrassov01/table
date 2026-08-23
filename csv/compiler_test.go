package csv

import (
	"errors"
	"testing"

	"github.com/nekrassov01/table"
	"github.com/nekrassov01/table/internal/testutil"
	"github.com/nekrassov01/table/internal/value"
)

func Test_compiler_prepare(t *testing.T) {
	type fields struct {
		input configResult
		state compilerState
	}
	type want struct {
		rowLen   int
		rowCap   int
		cellLen  int
		cellCap  int
		quoteLen int
		values   []string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "allocates for header body and footer",
			fields: fields{
				input: configResult{
					header:   []string{"A", "B"},
					footer:   [][]string{{"x", "y"}},
					bodyRows: 2,
					columns:  []column{{}, {}},
				},
			},
			want: want{
				rowCap:  4,
				cellCap: 8,
				values:  []string{"", ""},
			},
		},
		{
			name: "reuses and clears storage",
			fields: fields{
				input: configResult{
					bodyRows: 1,
					columns:  []column{{}, {}},
				},
				state: compilerState{
					rows:   []row{{}, {}},
					cells:  []cell{{value: "a"}, {value: "b"}},
					quotes: []byte("quoted"),
					values: []string{"a", "b", "c"},
				},
			},
			want: want{
				rowCap:  2,
				cellCap: 2,
				values:  []string{"", ""},
			},
		},
		{
			name: "prepares zero columns",
			fields: fields{
				input: configResult{},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			o := &compiler{
				input: test.fields.input,
				state: &state,
			}
			o.prepare()
			got := want{
				rowLen:   len(state.rows),
				rowCap:   cap(state.rows),
				cellLen:  len(state.cells),
				cellCap:  cap(state.cells),
				quoteLen: len(state.quotes),
				values:   state.values,
			}
			testutil.AssertValue(t, got, test.want, "prepare")
		})
	}
}

func Test_compiler_compileHeader(t *testing.T) {
	type fields struct {
		input configResult
		state compilerState
		err   error
	}
	type want struct {
		values []string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "compiles header and index label",
			fields: fields{
				input: configResult{
					option: &option{
						delimiter:   ',',
						indexOffset: 1,
					},
					header:  []string{"A,B", "C"},
					columns: []column{{}, {}, {}},
				},
				state: compilerState{
					cells:  make([]cell, 0, 3),
					values: make([]string, 3),
				},
			},
			want: want{
				values: []string{"#", `"A,B"`, "C"},
			},
		},
		{
			name: "skips absent header",
			fields: fields{
				input: configResult{
					option: &option{
						delimiter: '\t',
					},
				},
			},
		},
		{
			name: "skips after error",
			fields: fields{
				input: configResult{
					header: []string{"A"},
				},
				err: testutil.NewError(),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			var strings value.Store
			o := &compiler{
				input:   test.fields.input,
				state:   &state,
				strings: &strings,
				err:     test.fields.err,
			}
			o.compileHeader()
			var values []string
			if len(o.output.header.cells) > 0 {
				values = make([]string, len(o.output.header.cells))
				for index := range o.output.header.cells {
					values[index] = o.output.header.cells[index].value
				}
			}
			got := want{
				values: values,
			}
			testutil.AssertValue(t, got, test.want, "compileHeader")
		})
	}
}

func Test_compiler_compileBody(t *testing.T) {
	type fields struct {
		input configResult
		state compilerState
	}
	type args struct {
		sources [][]any
	}
	type want struct {
		values   [][]string
		hasError bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "compiles rows in order",
			fields: fields{
				input: configResult{
					option: &option{
						delimiter: '\t',
					},
					columns: []column{{}, {}},
				},
				state: compilerState{
					cells:  make([]cell, 0, 4),
					rows:   make([]row, 0, 2),
					values: make([]string, 2),
				},
			},
			args: args{
				sources: [][]any{{"a", 1}, {"b", 2}},
			},
			want: want{
				values: [][]string{{"a", "1"}, {"b", "2"}},
			},
		},
		{
			name: "stops after wider row",
			fields: fields{
				input: configResult{
					option: &option{
						delimiter: '\t',
					},
					columns: []column{{}},
				},
				state: compilerState{
					cells:  make([]cell, 0, 2),
					rows:   make([]row, 0, 2),
					values: make([]string, 1),
				},
			},
			args: args{
				sources: [][]any{{"a", "b"}, {"c"}},
			},
			want: want{
				hasError: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			var strings value.Store
			o := &compiler{
				input:     test.fields.input,
				state:     &state,
				strings:   &strings,
				bodyStart: -1,
			}
			o.compileBody(test.args.sources)
			var values [][]string
			if len(o.output.body) > 0 {
				values = make([][]string, len(o.output.body))
				for rowIndex := range o.output.body {
					values[rowIndex] = make([]string, len(o.output.body[rowIndex].cells))
					for columnIndex := range o.output.body[rowIndex].cells {
						values[rowIndex][columnIndex] = o.output.body[rowIndex].cells[columnIndex].value
					}
				}
			}
			got := want{
				values:   values,
				hasError: o.err != nil,
			}
			testutil.AssertValue(t, got, test.want, "compileBody")
		})
	}
}

func Test_compiler_compileFooter(t *testing.T) {
	type fields struct {
		input configResult
		state compilerState
		err   error
	}
	type want struct {
		values    [][]string
		isColumns bool
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "compiles footer rows with empty index",
			fields: fields{
				input: configResult{
					option: &option{
						delimiter:   ',',
						indexOffset: 1,
					},
					footer:        [][]string{{"sum", "1,2"}, {"end"}},
					columns:       []column{{}, {}, {}},
					footerColumns: 2,
				},
				state: compilerState{
					cells:  make([]cell, 0, 6),
					rows:   make([]row, 0, 2),
					values: make([]string, 3),
				},
			},
			want: want{
				values: [][]string{{"", "sum", `"1,2"`}, {"", "end", ""}},
			},
		},
		{
			name: "rejects footer wider than columns",
			fields: fields{
				input: configResult{
					option:        &option{},
					footer:        [][]string{{"a", "b"}},
					columns:       []column{{}},
					footerColumns: 2,
				},
			},
			want: want{
				isColumns: true,
			},
		},
		{
			name: "skips empty footer",
			fields: fields{
				input: configResult{
					option: &option{},
				},
			},
		},
		{
			name: "skips after error",
			fields: fields{
				input: configResult{
					option: &option{},
					footer: [][]string{{"a"}},
				},
				err: testutil.NewError(),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			var strings value.Store
			o := &compiler{
				input:   test.fields.input,
				state:   &state,
				strings: &strings,
				err:     test.fields.err,
			}
			o.compileFooter()
			var values [][]string
			if len(o.output.footer) > 0 {
				values = make([][]string, len(o.output.footer))
				for rowIndex := range o.output.footer {
					values[rowIndex] = make([]string, len(o.output.footer[rowIndex].cells))
					for columnIndex := range o.output.footer[rowIndex].cells {
						values[rowIndex][columnIndex] = o.output.footer[rowIndex].cells[columnIndex].value
					}
				}
			}
			got := want{
				values:    values,
				isColumns: errors.Is(o.err, table.ErrColumnCount),
			}
			testutil.AssertValue(t, got, test.want, "compileFooter")
		})
	}
}

func Test_compiler_compileBand(t *testing.T) {
	type fields struct {
		input configResult
		state compilerState
	}
	type args struct {
		labels []string
		header bool
	}
	type want struct {
		values []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "header carries index marker and labels",
			fields: fields{
				input: configResult{
					option: &option{
						delimiter:   '\t',
						indexOffset: 1,
					},
					columns: []column{{}, {}, {}},
				},
				state: compilerState{
					cells:  make([]cell, 0, 3),
					values: make([]string, 3),
				},
			},
			args: args{
				labels: []string{"A", "B"},
				header: true,
			},
			want: want{
				values: []string{"#", "A", "B"},
			},
		},
		{
			name: "footer leaves index and missing tail empty",
			fields: fields{
				input: configResult{
					option: &option{
						delimiter:   '\t',
						indexOffset: 1,
					},
					columns: []column{{}, {}, {}},
				},
				state: compilerState{
					cells:  make([]cell, 0, 3),
					values: make([]string, 3),
				},
			},
			args: args{
				labels: []string{"sum"},
			},
			want: want{
				values: []string{"", "sum", ""},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			var strings value.Store
			o := &compiler{
				input:   test.fields.input,
				state:   &state,
				strings: &strings,
			}
			r := o.compileBand(test.args.labels, test.args.header)
			values := make([]string, len(r.cells))
			for index := range r.cells {
				values[index] = r.cells[index].value
			}
			got := want{
				values: values,
			}
			testutil.AssertValue(t, got, test.want, "compileBand")
		})
	}
}

func Test_compiler_compileRow(t *testing.T) {
	type fields struct {
		input     configResult
		state     compilerState
		bodyStart int
	}
	type args struct {
		source   []any
		rowIndex int
	}
	type want struct {
		values      []string
		bodyStart   int
		isColumns   bool
		stateValues []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "resolves index transformed fallback missing and quoted values",
			fields: fields{
				input: configResult{
					option: &option{
						delimiter:   ',',
						placeholder: "-",
						indexOffset: 1,
					},
					columns: []column{
						{},
						{
							transformer: func(any) string {
								return "new,value"
							},
						},
						{
							transformer: func(any) string {
								return ""
							},
						},
						{},
					},
				},
				state: compilerState{
					cells:  make([]cell, 0, 4),
					rows:   make([]row, 0, 1),
					values: make([]string, 4),
				},
				bodyStart: -1,
			},
			args: args{
				source:   []any{testutil.PanicStringer{}, ""},
				rowIndex: 2,
			},
			want: want{
				values:      []string{"3", `"new,value"`, "-", "-"},
				bodyStart:   0,
				stateValues: []string{"3", "new,value", "-", "-"},
			},
		},
		{
			name: "rejects wider row",
			fields: fields{
				input: configResult{
					option:  &option{},
					columns: []column{{}},
				},
				bodyStart: -1,
			},
			args: args{
				source: []any{"a", "b"},
			},
			want: want{
				bodyStart: -1,
				isColumns: true,
			},
		},
		{
			name: "retains existing body start",
			fields: fields{
				input: configResult{
					option: &option{
						delimiter: '\t',
					},
					columns: []column{{}},
				},
				state: compilerState{
					cells:  make([]cell, 0, 1),
					rows:   make([]row, 0, 1),
					values: make([]string, 1),
				},
				bodyStart: 0,
			},
			args: args{
				source: []any{"a"},
			},
			want: want{
				values:      []string{"a"},
				bodyStart:   0,
				stateValues: []string{"a"},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			var strings value.Store
			o := &compiler{
				input:     test.fields.input,
				state:     &state,
				strings:   &strings,
				bodyStart: test.fields.bodyStart,
			}
			o.compileRow(test.args.source, test.args.rowIndex)
			values := []string(nil)
			if len(o.output.body) > 0 {
				values = make([]string, len(o.output.body[0].cells))
				for index := range o.output.body[0].cells {
					values[index] = o.output.body[0].cells[index].value
				}
			}
			got := want{
				values:      values,
				bodyStart:   o.bodyStart,
				isColumns:   errors.Is(o.err, table.ErrColumnCount),
				stateValues: state.values,
			}
			testutil.AssertValue(t, got, test.want, "compileRow")
		})
	}
}

func Test_compiler_compileCells(t *testing.T) {
	type fields struct {
		input configResult
		state compilerState
	}
	type args struct {
		row row
	}
	type want struct {
		values []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "quotes current values",
			fields: fields{
				input: configResult{
					option: &option{
						delimiter: ',',
					},
				},
				state: compilerState{
					values: []string{"plain", "a,b", `a"b`},
				},
			},
			args: args{
				row: row{
					cells: make([]cell, 3),
				},
			},
			want: want{
				values: []string{"plain", `"a,b"`, `"a""b"`},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			o := &compiler{
				input: test.fields.input,
				state: &state,
			}
			o.compileCells(test.args.row)
			values := make([]string, len(test.args.row.cells))
			for index := range test.args.row.cells {
				values[index] = test.args.row.cells[index].value
			}
			got := want{
				values: values,
			}
			testutil.AssertValue(t, got, test.want, "compileCells")
		})
	}
}

func Test_compiler_reserveBand(t *testing.T) {
	type fields struct {
		input configResult
		state compilerState
	}
	type args struct {
		count int
	}
	type want struct {
		reservedLen int
		rowLen      int
		cellCap     int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "grows rows and cell capacity",
			fields: fields{
				input: configResult{
					columns: []column{{}, {}},
				},
				state: compilerState{
					rows:  []row{{}},
					cells: make([]cell, 0, 1),
				},
			},
			args: args{
				count: 2,
			},
			want: want{
				reservedLen: 2,
				rowLen:      3,
				cellCap:     4,
			},
		},
		{
			name: "reserves zero rows",
			fields: fields{
				input: configResult{
					columns: []column{{}},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			o := &compiler{
				input: test.fields.input,
				state: &state,
			}
			reserved := o.reserveBand(test.args.count)
			got := want{
				reservedLen: len(reserved),
				rowLen:      len(state.rows),
				cellCap:     cap(state.cells),
			}
			testutil.AssertValue(t, got, test.want, "reserveBand")
		})
	}
}

func Test_compiler_newRow(t *testing.T) {
	type fields struct {
		input configResult
		state compilerState
	}
	type want struct {
		rowLen  int
		cellLen int
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "reserves one row view",
			fields: fields{
				input: configResult{
					columns: []column{{}, {}},
				},
				state: compilerState{
					cells: make([]cell, 0, 2),
				},
			},
			want: want{
				rowLen:  2,
				cellLen: 2,
			},
		},
		{
			name: "reserves zero-column row",
			fields: fields{
				state: compilerState{
					cells: make([]cell, 0),
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			o := &compiler{
				input: test.fields.input,
				state: &state,
			}
			r := o.newRow()
			got := want{
				rowLen:  len(r.cells),
				cellLen: len(state.cells),
			}
			testutil.AssertValue(t, got, test.want, "newRow")
		})
	}
}
