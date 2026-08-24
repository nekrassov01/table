package html

import (
	"errors"
	"testing"

	"github.com/nekrassov01/table"
	"github.com/nekrassov01/table/internal/scope"
	spans "github.com/nekrassov01/table/internal/span"
	"github.com/nekrassov01/table/internal/testutil"
	"github.com/nekrassov01/table/internal/value"
)

func Test_compiler_prepare(t *testing.T) {
	type fields struct {
		input   configResult
		state   compilerState
		strings value.Store
	}
	type want struct {
		rowsCap       int
		cellsCap      int
		valuesLen     int
		valuesCap     int
		columnSizes   []int
		rowspanHeader uint64
		rowspanBody   uint64
		rowspanFooter uint64
		colspanHeader uint64
		colspanBody   uint64
		colspanFooter uint64
		caption       string
		previousSpan  uint64
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "resets storage and derives masks",
			fields: fields{
				input: configResult{
					option: &option{
						caption: "<caption>",
					},
					header:   [][]string{{"header"}},
					footer:   [][]string{{"footer"}},
					bodyRows: 1,
					columns: []columnConfig{
						{
							rowspan: ScopeHeader | ScopeBody,
							colspan: ScopeFooter,
						},
						{
							rowspan: ScopeBody | ScopeFooter,
							colspan: ScopeHeader | ScopeBody,
						},
					},
				},
				state: func() compilerState {
					var rowspans scope.Masks
					rowspans.Mark(ScopeHeader|ScopeBody|ScopeFooter, 3)
					var colspans scope.Masks
					colspans.Mark(ScopeHeader|ScopeBody|ScopeFooter, 3)
					return compilerState{
						cells: []cell{
							{
								value: "old",
							},
						},
						columnSizes: []int{9, 9},
						rowspans:    rowspans,
						colspans:    colspans,
					}
				}(),
			},
			want: want{
				rowsCap:       3,
				cellsCap:      6,
				valuesLen:     2,
				valuesCap:     2,
				columnSizes:   []int{0, 0},
				rowspanHeader: 0b01,
				rowspanBody:   0b11,
				rowspanFooter: 0b10,
				colspanHeader: 0b10,
				colspanBody:   0b10,
				colspanFooter: 0b01,
				caption:       "&lt;caption&gt;",
			},
		},
		{
			name: "reuses storage and allocates missing estimates",
			fields: fields{
				input: configResult{
					option:   &option{},
					bodyRows: 1,
					columns:  make([]columnConfig, 2),
				},
				state: compilerState{
					rows:  make([]row, 2, 5),
					cells: make([]cell, 3, 8),
					values: func() []string {
						values := make([]string, 3, 4)
						values[0] = "old"
						return values
					}(),
				},
			},
			want: want{
				rowsCap:     5,
				cellsCap:    8,
				valuesLen:   2,
				valuesCap:   4,
				columnSizes: []int{0, 0},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			spans.Rowspans(1, []string{"same"}, &state.previousBody)
			strings := test.fields.strings
			o := &compiler{
				input:   test.fields.input,
				state:   &state,
				strings: &strings,
				output: compilerResult{
					configResult: test.fields.input,
				},
			}
			o.prepare()
			got := want{
				rowsCap:       cap(state.rows),
				cellsCap:      cap(state.cells),
				valuesLen:     len(state.values),
				valuesCap:     cap(state.values),
				columnSizes:   state.columnSizes,
				rowspanHeader: state.rowspans.Resolve(ScopeHeader),
				rowspanBody:   state.rowspans.Resolve(ScopeBody),
				rowspanFooter: state.rowspans.Resolve(ScopeFooter),
				colspanHeader: state.colspans.Resolve(ScopeHeader),
				colspanBody:   state.colspans.Resolve(ScopeBody),
				colspanFooter: state.colspans.Resolve(ScopeFooter),
				caption:       o.output.caption,
				previousSpan:  spans.Rowspans(1, []string{"same"}, &state.previousBody),
			}
			testutil.AssertValue(t, got, test.want, "prepare")
		})
	}
}

func Test_compiler_compileHeader(t *testing.T) {
	type fields struct {
		input  configResult
		state  compilerState
		output compilerResult
	}
	type want struct {
		values []string
		rows   int
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "retains header order",
			fields: fields{
				input: configResult{
					option:  &option{},
					header:  [][]string{{"top"}, {"bottom"}},
					columns: []columnConfig{{}},
				},
				state: compilerState{
					cells:       make([]cell, 0, 2),
					values:      make([]string, 1),
					columnSizes: make([]int, 1),
				},
			},
			want: want{
				values: []string{"top", "bottom"},
				rows:   2,
			},
		},
		{
			name: "keeps existing output without header",
			fields: fields{
				input: configResult{
					option: &option{},
				},
				output: compilerResult{
					header: []row{
						{
							cells: []cell{{value: "existing"}},
						},
					},
				},
			},
			want: want{
				values: []string{"existing"},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			strings := value.Store{}
			o := &compiler{
				input:   test.fields.input,
				state:   &state,
				strings: &strings,
				output:  test.fields.output,
			}
			o.compileHeader()
			got := want{
				rows: len(state.rows),
			}
			for _, r := range o.output.header {
				for _, compiled := range r.cells {
					got.values = append(got.values, compiled.value)
				}
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
		values []string
		err    bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "retains body order",
			fields: fields{
				input: configResult{
					option:  &option{},
					columns: []columnConfig{{}},
				},
				state: compilerState{
					cells:       make([]cell, 0, 2),
					values:      make([]string, 1),
					columnSizes: make([]int, 1),
				},
			},
			args: args{
				sources: [][]any{{"first"}, {"second"}},
			},
			want: want{
				values: []string{"first", "second"},
			},
		},
		{
			name: "stops after a structural error",
			fields: fields{
				input: configResult{
					option:  &option{},
					columns: []columnConfig{{}},
				},
				state: compilerState{
					columnSizes: make([]int, 1),
				},
			},
			args: args{
				sources: [][]any{{"a", "b"}, {"unreached"}},
			},
			want: want{
				err: true,
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
			o.compileBody(test.args.sources)
			got := want{
				err: errors.Is(o.err, table.ErrColumnCount),
			}
			for _, r := range o.output.body {
				for _, compiled := range r.cells {
					got.values = append(got.values, compiled.value)
				}
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
		values []string
		err    bool
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "retains footer order",
			fields: fields{
				input: configResult{
					option:        &option{},
					footer:        [][]string{{"subtotal"}, {"total"}},
					columns:       []columnConfig{{}},
					footerColumns: 1,
				},
				state: compilerState{
					cells:       make([]cell, 0, 2),
					values:      make([]string, 1),
					columnSizes: make([]int, 1),
				},
			},
			want: want{
				values: []string{"subtotal", "total"},
			},
		},
		{
			name: "rejects footer wider than header",
			fields: fields{
				input: configResult{
					option:        &option{},
					footer:        [][]string{{"a", "b"}},
					columns:       []columnConfig{{}},
					footerColumns: 2,
				},
			},
			want: want{
				err: true,
			},
		},
		{
			name: "keeps existing error",
			fields: fields{
				input: configResult{
					option: &option{},
					footer: [][]string{{"total"}},
				},
				err: testutil.NewError(),
			},
			want: want{
				err: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			strings := value.Store{}
			o := &compiler{
				input:   test.fields.input,
				state:   &state,
				strings: &strings,
				err:     test.fields.err,
			}
			o.compileFooter()
			got := want{
				err: o.err != nil,
			}
			for _, r := range o.output.footer {
				for _, compiled := range r.cells {
					got.values = append(got.values, compiled.value)
				}
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
		scope  Scope
	}
	type want struct {
		values []string
	}
	columns := make([]columnConfig, 3)
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "header index and labels",
			fields: fields{
				input: configResult{
					option: &option{
						indexOffset: 1,
					},
					columns: columns,
				},
				state: compilerState{
					cells:       make([]cell, 0, 3),
					values:      make([]string, 3),
					columnSizes: make([]int, 3),
				},
			},
			args: args{
				labels: []string{"label"},
				scope:  ScopeHeader,
			},
			want: want{
				values: []string{"#", "label", ""},
			},
		},
		{
			name: "footer labels omit index",
			fields: fields{
				input: configResult{
					option: &option{
						indexOffset: 1,
					},
					columns: columns,
				},
				state: compilerState{
					cells:       make([]cell, 0, 3),
					values:      make([]string, 3),
					columnSizes: make([]int, 3),
				},
			},
			args: args{
				labels: []string{"total", ""},
				scope:  ScopeFooter,
			},
			want: want{
				values: []string{"", "total", ""},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			strings := value.Store{}
			o := &compiler{
				input:   test.fields.input,
				state:   &state,
				strings: &strings,
			}
			compiled := o.compileBand(test.args.labels, test.args.scope)
			got := want{}
			for _, cell := range compiled.cells {
				got.values = append(got.values, cell.value)
			}
			testutil.AssertValue(t, got, test.want, "compileBand")
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
		values []string
		err    bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "index and placeholder",
			fields: fields{
				input: configResult{
					option: &option{
						placeholder: "-",
						indexOffset: 1,
					},
					columns: make([]columnConfig, 2),
				},
				state: compilerState{
					cells:       make([]cell, 0, 2),
					values:      make([]string, 2),
					columnSizes: make([]int, 2),
				},
			},
			args: args{
				rowIndex: 2,
			},
			want: want{
				values: []string{"3", "-"},
			},
		},
		{
			name: "rejects row wider than header",
			fields: fields{
				input: configResult{
					option:  &option{},
					columns: []columnConfig{{}},
				},
				state: compilerState{
					columnSizes: make([]int, 1),
				},
			},
			args: args{
				source: []any{"a", "b"},
			},
			want: want{
				err: true,
			},
		},
		{
			name: "applies transformer and formatting fallbacks",
			fields: fields{
				input: configResult{
					option: &option{
						placeholder: "-",
					},
					columns: func() []columnConfig {
						columns := make([]columnConfig, 3)
						columns[0].transformer.fn = func(any) (string, *Color, *Decoration) {
							return "new", ColorFgRed, DecorationBold
						}
						columns[1].transformer.colors.Set(ScopeBody, ColorFgRed)
						columns[1].transformer.decorations.Set(ScopeBody, DecorationBold)
						columns[1].transformer.fn = func(any) (string, *Color, *Decoration) {
							return "", nil, nil
						}
						return columns
					}(),
				},
				state: compilerState{
					cells:       make([]cell, 0, 3),
					values:      make([]string, 3),
					columnSizes: make([]int, 3),
				},
			},
			args: args{
				source: []any{testutil.PanicStringer{}, "", 42},
			},
			want: want{
				values: []string{"new", "-", "42"},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			strings := value.Store{}
			o := &compiler{
				input:   test.fields.input,
				state:   &state,
				strings: &strings,
				output: compilerResult{
					configResult: test.fields.input,
				},
			}
			o.compileRow(test.args.source, test.args.rowIndex)
			values := []string(nil)
			if len(o.output.body) > 0 {
				for _, compiled := range o.output.body[0].cells {
					values = append(values, compiled.value)
				}
			}
			got := want{
				values: values,
				err:    errors.Is(o.err, table.ErrColumnCount),
			}
			testutil.AssertValue(t, got, test.want, "compileRow")
		})
	}
}

func Test_compiler_reserveBand(t *testing.T) {
	type fields struct {
		input configResult
		state compilerState
	}
	type args struct {
		rows int
	}
	type want struct {
		stateRows int
		bandRows  int
		cellRoom  bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "reserves rows and cell capacity",
			fields: fields{
				input: configResult{
					columns: make([]columnConfig, 2),
				},
				state: compilerState{
					cells: make([]cell, 1),
					rows:  make([]row, 1),
				},
			},
			args: args{
				rows: 2,
			},
			want: want{
				stateRows: 3,
				bandRows:  2,
				cellRoom:  true,
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
			band := o.reserveBand(test.args.rows)
			got := want{
				stateRows: len(state.rows),
				bandRows:  len(band),
				cellRoom:  cap(state.cells)-len(state.cells) >= test.args.rows*len(test.fields.input.columns),
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
		cells       int
		rowCells    int
		firstValue  string
		secondValue string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "takes a row from cell storage",
			fields: fields{
				input: configResult{
					columns: make([]columnConfig, 2),
				},
				state: compilerState{
					cells: func() []cell {
						cells := make([]cell, 1, 3)
						cells[0].value = "existing"
						return cells
					}(),
				},
			},
			want: want{
				cells:      3,
				rowCells:   2,
				firstValue: "first",
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
			row := o.newRow()
			row.cells[0].value = "first"
			got := want{
				cells:       len(state.cells),
				rowCells:    len(row.cells),
				firstValue:  state.cells[1].value,
				secondValue: state.cells[2].value,
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
		scope    Scope
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
				state: func() compilerState {
					var rowspans scope.Masks
					rowspans.Mark(ScopeBody, 0)
					var colspans scope.Masks
					colspans.Mark(ScopeBody, 0)
					colspans.Mark(ScopeBody, 1)
					return compilerState{
						values:   []string{"same", "same"},
						rowspans: rowspans,
						colspans: colspans,
					}
				}(),
			},
			args: args{
				row: row{
					cells: make([]cell, 2),
				},
				scope:    ScopeBody,
				previous: []string{"same", "old"},
			},
			want: want{
				rowspans: 0b01,
			},
		},
		{
			name: "detects horizontal candidates",
			fields: fields{
				state: func() compilerState {
					var colspans scope.Masks
					colspans.Mark(ScopeHeader, 0)
					colspans.Mark(ScopeHeader, 1)
					return compilerState{
						values:   []string{"same", "same"},
						colspans: colspans,
					}
				}(),
			},
			args: args{
				row: row{
					cells: make([]cell, 2),
				},
				scope: ScopeHeader,
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
				scope: ScopeFooter,
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
			spans.Rowspans(state.rowspans.Resolve(test.args.scope), test.args.previous, &state.previousBody)
			o := &compiler{
				state: &state,
			}
			r := test.args.row
			o.setSpans(&r, test.args.scope, &state.previousBody)
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
		values      []string
		sizes       []int
		colors      []*Color
		decorations []*Decoration
		columnSizes []int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "escapes values and omits vertical continuations",
			fields: fields{
				state: compilerState{
					values:      []string{"same", "<x>"},
					columnSizes: make([]int, 2),
				},
			},
			args: args{
				row: row{
					cells: []cell{
						{
							colspan: 1,
						},
						{
							color:      &Color{},
							decoration: &Decoration{},
							colspan:    1,
						},
					},
					rowspans: 0b01,
				},
			},
			want: want{
				values:      []string{"", "&lt;x&gt;"},
				sizes:       []int{0, 9},
				colors:      []*Color{nil, nil},
				decorations: []*Decoration{nil, nil},
				columnSizes: []int{0, 9},
			},
		},
		{
			name: "stores value and markup size",
			fields: fields{
				state: compilerState{
					values:      []string{"x"},
					columnSizes: make([]int, 1),
				},
			},
			args: args{
				row: row{
					cells: []cell{
						{
							color:      ColorFgRed,
							decoration: DecorationBold,
							colspan:    1,
						},
					},
				},
			},
			want: want{
				values:      []string{"x"},
				sizes:       []int{49},
				colors:      []*Color{ColorFgRed},
				decorations: []*Decoration{DecorationBold},
				columnSizes: []int{49},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			o := &compiler{
				state: &state,
			}
			o.compileCells(test.args.row)
			got := want{
				columnSizes: state.columnSizes,
			}
			for _, compiled := range test.args.row.cells {
				got.values = append(got.values, compiled.value)
				got.sizes = append(got.sizes, compiled.size)
				got.colors = append(got.colors, compiled.color)
				got.decorations = append(got.decorations, compiled.decoration)
			}
			testutil.AssertValue(t, got, test.want, "compileCells")
		})
	}
}
