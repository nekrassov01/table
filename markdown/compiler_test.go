package markdown

import (
	"errors"
	"testing"

	"github.com/nekrassov01/table"
	"github.com/nekrassov01/table/internal/span"
	"github.com/nekrassov01/table/internal/testutil"
	"github.com/nekrassov01/table/internal/value"
)

func Test_compiler_prepare(t *testing.T) {
	type fields struct {
		input configResult
		state compilerState
	}
	type want struct {
		rowsLen      int
		rowsCap      int
		cellsLen     int
		cellsCap     int
		escapesLen   int
		valuesLen    int
		valuesCap    int
		rowspans     uint64
		colspans     uint64
		previousSpan uint64
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "allocates storage and derives masks",
			fields: fields{
				input: configResult{
					option:   &option{},
					bodyRows: 2,
					columns: []columnConfig{
						{
							rowspan: true,
						},
						{
							colspan: true,
						},
					},
				},
				state: compilerState{
					rows:     []row{{}},
					cells:    []cell{{}},
					escapes:  []byte("old"),
					values:   []string{"old"},
					rowspans: 8,
					colspans: 8,
				},
			},
			want: want{
				rowsCap:   3,
				cellsCap:  6,
				valuesLen: 2,
				valuesCap: 2,
				rowspans:  1,
				colspans:  2,
			},
		},
		{
			name: "reuses storage",
			fields: fields{
				input: configResult{
					option:   &option{},
					bodyRows: 1,
					columns:  make([]columnConfig, 2),
				},
				state: compilerState{
					rows:   make([]row, 2, 5),
					cells:  make([]cell, 3, 8),
					values: make([]string, 3, 4),
				},
			},
			want: want{
				rowsCap:   5,
				cellsCap:  8,
				valuesLen: 2,
				valuesCap: 4,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			span.Rowspans(1, []string{"same"}, &state.previousBody)
			o := &compiler{
				input: test.fields.input,
				state: &state,
				output: compilerResult{
					configResult: test.fields.input,
				},
			}
			o.prepare()
			got := want{
				rowsLen:      len(state.rows),
				rowsCap:      cap(state.rows),
				cellsLen:     len(state.cells),
				cellsCap:     cap(state.cells),
				escapesLen:   len(state.escapes),
				valuesLen:    len(state.values),
				valuesCap:    cap(state.values),
				rowspans:     state.rowspans,
				colspans:     state.colspans,
				previousSpan: span.Rowspans(1, []string{"same"}, &state.previousBody),
			}
			testutil.AssertValue(t, got, test.want, "prepare")
		})
	}
}

func Test_compiler_compileHeader(t *testing.T) {
	type fields struct {
		input configResult
		state compilerState
	}
	type want struct {
		header row
		rows   []row
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "compiles index labels and markup",
			fields: fields{
				input: func() configResult {
					colored := columnConfig{}
					colored.transformer.colors.Set(ScopeHeader, ColorFgRed)
					colored.transformer.decorations.Set(ScopeHeader, DecorationBold)
					return configResult{
						option: &option{
							indexOffset: 1,
						},
						header:  []string{"A", ""},
						columns: []columnConfig{{}, colored, colored},
					}
				}(),
				state: compilerState{
					cells: func() []cell {
						cells := []cell{
							{
								color:      ColorFgBlue,
								decoration: DecorationItalic,
							},
							{},
							{
								color:      ColorFgBlue,
								decoration: DecorationItalic,
							},
						}
						return cells[:0]
					}(),
					rows:   make([]row, 0, 1),
					values: make([]string, 3),
				},
			},
			want: want{
				header: row{
					cells: []cell{
						{
							value: "#",
							width: 1,
							size:  1,
						},
						{
							value:      "A",
							width:      36,
							size:       36,
							color:      ColorFgRed,
							decoration: DecorationBold,
						},
						{},
					},
				},
				rows: []row{
					{
						cells: []cell{
							{
								value: "#",
								width: 1,
								size:  1,
							},
							{
								value:      "A",
								width:      36,
								size:       36,
								color:      ColorFgRed,
								decoration: DecorationBold,
							},
							{},
						},
					},
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
				output: compilerResult{
					configResult: test.fields.input,
				},
			}
			o.compileHeader()
			got := want{
				header: o.output.header,
				rows:   state.rows,
			}
			testutil.AssertValue(t, got, test.want, "compileHeader")
		})
	}
}

func Test_compiler_compileBody(t *testing.T) {
	type fields struct {
		input configResult
	}
	type args struct {
		sources [][]any
	}
	type want struct {
		body []row
		err  bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "retains rows and span masks",
			fields: fields{
				input: configResult{
					option: &option{
						placeholder: "-",
					},
					bodyRows: 2,
					columns: []columnConfig{
						{
							rowspan: true,
							colspan: true,
						},
						{
							rowspan: true,
							colspan: true,
						},
					},
				},
			},
			args: args{
				sources: [][]any{
					{"x", "x"},
					{"x", "x"},
				},
			},
			want: want{
				body: []row{
					{
						cells: []cell{
							{
								value: "x",
								width: 1,
								size:  1,
							},
							{},
						},
						colspans: 0b10,
					},
					{
						cells:    []cell{{}, {}},
						rowspans: 0b11,
					},
				},
			},
		},
		{
			name: "stops at an overflowing row",
			fields: fields{
				input: configResult{
					option:   &option{},
					bodyRows: 2,
					columns:  []columnConfig{{}},
				},
			},
			args: args{
				sources: [][]any{
					{"a", "b"},
					{"c"},
				},
			},
			want: want{
				err: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := compilerState{}
			strings := value.Store{}
			o := &compiler{
				input:     test.fields.input,
				state:     &state,
				strings:   &strings,
				bodyStart: -1,
				output: compilerResult{
					configResult: test.fields.input,
				},
			}
			o.prepare()
			o.compileBody(test.args.sources)
			got := want{
				body: o.output.body,
				err:  errors.Is(o.err, table.ErrColumnCount),
			}
			testutil.AssertValue(t, got, test.want, "compileBody")
		})
	}
}

func Test_compiler_compileRow(t *testing.T) {
	type fields struct {
		input configResult
		state compilerState
	}
	type args struct {
		source   []any
		rowIndex int
	}
	type want struct {
		bodyStart int
		body      []row
		err       bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "resolves index transformed and missing cells",
			fields: fields{
				input: configResult{
					option: &option{
						placeholder: "-",
						indexOffset: 1,
					},
					columns: []columnConfig{
						{},
						{
							transformer: transformer{
								fn: func(any) (string, *Color, *Decoration) {
									return "x|", ColorFgRed, DecorationBold
								},
							},
						},
						{},
					},
				},
				state: compilerState{
					cells: func() []cell {
						cells := []cell{
							{
								color:      ColorFgBlue,
								decoration: DecorationItalic,
							},
							{},
							{
								color:      ColorFgBlue,
								decoration: DecorationItalic,
							},
						}
						return cells[:0]
					}(),
					rows:   make([]row, 0, 1),
					values: make([]string, 3),
				},
			},
			args: args{
				source:   []any{testutil.PanicStringer{}},
				rowIndex: 2,
			},
			want: want{
				bodyStart: 0,
				body: []row{
					{
						cells: []cell{
							{
								value: "3",
								width: 1,
								size:  1,
							},
							{
								value:      `x\|`,
								width:      38,
								size:       38,
								color:      ColorFgRed,
								decoration: DecorationBold,
							},
							{
								value: "-",
								width: 1,
								size:  1,
							},
						},
					},
				},
			},
		},
		{
			name: "rejects an overflowing row",
			fields: fields{
				input: configResult{
					option:  &option{},
					columns: []columnConfig{{}},
				},
				state: compilerState{
					cells:  make([]cell, 0, 1),
					rows:   make([]row, 0, 1),
					values: make([]string, 1),
				},
			},
			args: args{
				source: []any{"a", "b"},
			},
			want: want{
				bodyStart: -1,
				err:       true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			strings := value.Store{}
			o := &compiler{
				input:     test.fields.input,
				state:     &state,
				strings:   &strings,
				bodyStart: -1,
				output: compilerResult{
					configResult: test.fields.input,
				},
			}
			o.compileRow(test.args.source, test.args.rowIndex)
			got := want{
				bodyStart: o.bodyStart,
				body:      o.output.body,
				err:       errors.Is(o.err, table.ErrColumnCount),
			}
			testutil.AssertValue(t, got, test.want, "compileRow")
		})
	}
}

func Test_compiler_newRow(t *testing.T) {
	type fields struct {
		input configResult
		state compilerState
	}
	type args struct {
		value string
	}
	type want struct {
		row   row
		cells []cell
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "reserves column cells",
			fields: fields{
				input: configResult{
					columns: make([]columnConfig, 2),
				},
				state: compilerState{
					cells: func() []cell {
						cells := make([]cell, 3)
						cells[0] = cell{
							value: "kept",
						}
						cells[1] = cell{
							value: "old",
						}
						return cells[:1]
					}(),
				},
			},
			args: args{
				value: "new",
			},
			want: want{
				row: row{
					cells: []cell{
						{
							value: "new",
						},
						{},
					},
				},
				cells: []cell{
					{
						value: "kept",
					},
					{
						value: "new",
					},
					{},
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
			r.cells[0].value = test.args.value
			got := want{
				row:   r,
				cells: state.cells,
			}
			testutil.AssertValue(t, got, test.want, "newRow")
		})
	}
}

func Test_compiler_setSpans(t *testing.T) {
	type fields struct {
		state compilerState
	}
	type args struct {
		row      row
		previous []string
	}
	type want struct {
		rowspans uint64
		colspans uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "detects vertical continuations",
			fields: fields{
				state: compilerState{
					values:   []string{"same", "same"},
					rowspans: 0b01,
					colspans: 0b11,
				},
			},
			args: args{
				row: row{
					cells: make([]cell, 2),
				},
				previous: []string{"same", "old"},
			},
			want: want{
				rowspans: 0b01,
			},
		},
		{
			name: "detects horizontal absorptions",
			fields: fields{
				state: compilerState{
					values:   []string{"same", "same"},
					colspans: 0b11,
				},
			},
			args: args{
				row: row{
					cells: make([]cell, 2),
				},
			},
			want: want{
				colspans: 0b10,
			},
		},
		{
			name: "leaves row without configured spans",
			args: args{
				row: row{
					cells:    make([]cell, 2),
					rowspans: 0b01,
					colspans: 0b10,
				},
			},
			want: want{
				rowspans: 0b01,
				colspans: 0b10,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			span.Rowspans(state.rowspans, test.args.previous, &state.previousBody)
			o := &compiler{
				state: &state,
			}
			r := test.args.row
			o.setSpans(&r, &state.previousBody)
			got := want{
				rowspans: r.rowspans,
				colspans: r.colspans,
			}
			testutil.AssertValue(t, got, test.want, "setSpans")
		})
	}
}

func Test_compiler_compileCells(t *testing.T) {
	type fields struct {
		state compilerState
	}
	type args struct {
		row row
	}
	type want struct {
		row row
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "omits horizontal absorption",
			fields: fields{
				state: compilerState{
					values: []string{"x", "x"},
				},
			},
			args: args{
				row: row{
					cells:    []cell{{}, {}},
					colspans: 0b10,
				},
			},
			want: want{
				row: row{
					cells: []cell{
						{
							value: "x",
							width: 1,
							size:  1,
						},
						{},
					},
					colspans: 0b10,
				},
			},
		},
		{
			name: "escapes decorated colored value",
			fields: fields{
				state: compilerState{
					values: []string{"x|"},
				},
			},
			args: args{
				row: row{
					cells: []cell{
						{
							color:      ColorFgRed,
							decoration: DecorationBold,
						},
					},
				},
			},
			want: want{
				row: row{
					cells: []cell{
						{
							value:      `x\|`,
							width:      38,
							size:       38,
							color:      ColorFgRed,
							decoration: DecorationBold,
						},
					},
				},
			},
		},
		{
			name: "grows code fence",
			fields: fields{
				state: compilerState{
					values: []string{"a`b"},
				},
			},
			args: args{
				row: row{
					cells: []cell{
						{
							decoration: DecorationCode,
						},
					},
				},
			},
			want: want{
				row: row{
					cells: []cell{
						{
							value: "a`b",
							width: 7,
							size:  7,
							ticks: 2,
						},
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			o := &compiler{
				state: &state,
			}
			r := test.args.row
			o.compileCells(r)
			got := want{
				row: r,
			}
			testutil.AssertValue(t, got, test.want, "compileCells")
		})
	}
}
