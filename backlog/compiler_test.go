package backlog

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
		rowsLen       int
		rowsCap       int
		cellsLen      int
		cellsCap      int
		valuesLen     int
		valuesCap     int
		rowspanHeader uint64
		rowspanBody   uint64
		rowspanFooter uint64
		colspanHeader uint64
		colspanBody   uint64
		colspanFooter uint64
		previousSpan  uint64
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
					header:   [][]string{{"header"}},
					footer:   [][]string{{"footer"}},
					bodyRows: 1,
					columns: []column{
						{rowspan: ScopeHeader | ScopeBody, colspan: ScopeFooter},
						{rowspan: ScopeBody | ScopeFooter, colspan: ScopeHeader | ScopeBody},
					},
				},
			},
			want: want{
				rowsCap:       3,
				cellsCap:      6,
				valuesLen:     2,
				valuesCap:     2,
				rowspanHeader: 0b01,
				rowspanBody:   0b11,
				rowspanFooter: 0b10,
				colspanHeader: 0b10,
				colspanBody:   0b10,
				colspanFooter: 0b01,
			},
		},
		{
			name: "reuses storage",
			fields: fields{
				input: configResult{
					bodyRows: 1,
					columns:  make([]column, 2),
				},
				state: compilerState{
					rows:  make([]row, 2, 5),
					cells: make([]cell, 3, 8),
					values: func() []string {
						values := make([]string, 3, 4)
						values[0] = "old"
						return values
					}(),
					escapes: []byte("old"),
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
			spans.Rowspans(1, []string{"same"}, &state.previousBody)
			strings := test.fields.strings
			o := &compiler{
				input:   test.fields.input,
				state:   &state,
				strings: &strings,
			}
			o.prepare()
			previousSpan := spans.Rowspans(1, []string{"same"}, &state.previousBody)
			got := want{
				rowsLen:       len(state.rows),
				rowsCap:       cap(state.rows),
				cellsLen:      len(state.cells),
				cellsCap:      cap(state.cells),
				valuesLen:     len(state.values),
				valuesCap:     cap(state.values),
				rowspanHeader: state.rowspans.Resolve(ScopeHeader),
				rowspanBody:   state.rowspans.Resolve(ScopeBody),
				rowspanFooter: state.rowspans.Resolve(ScopeFooter),
				colspanHeader: state.colspans.Resolve(ScopeHeader),
				colspanBody:   state.colspans.Resolve(ScopeBody),
				colspanFooter: state.colspans.Resolve(ScopeFooter),
				previousSpan:  previousSpan,
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
		values   [][]string
		rowspans []uint64
		colspans []uint64
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "empty",
			fields: fields{
				input: configResult{
					option: &option{},
				},
			},
		},
		{
			name: "compiles stacked header in output order",
			fields: fields{
				input: configResult{
					option: &option{},
					header: [][]string{
						{"same", "Top", "Top"},
						{"same", "A", "B"},
					},
					columns: []column{
						{rowspan: ScopeHeader},
						{colspan: ScopeHeader},
						{colspan: ScopeHeader},
					},
				},
			},
			want: want{
				values: [][]string{
					{"", "Top", ""},
					{"same", "A", "B"},
				},
				rowspans: []uint64{0b001, 0},
				colspans: []uint64{0b100, 0},
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
			o.compileHeader()
			var got want
			for index := range o.output.header {
				values := make([]string, len(o.output.header[index].cells))
				for cellIndex := range o.output.header[index].cells {
					values[cellIndex] = o.output.header[index].cells[cellIndex].value
				}
				got.values = append(got.values, values)
				got.rowspans = append(got.rowspans, o.output.header[index].rowspans)
				got.colspans = append(got.colspans, o.output.header[index].colspans)
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
		rows      int
		first     string
		hasError  bool
		isColumns bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "compiles rows",
			fields: fields{
				input: configResult{
					option: &option{}, bodyRows: 2,
					columns: []column{{}},
				},
			},
			args: args{
				sources: [][]any{{"a"}, {"b"}},
			},
			want: want{
				rows:  2,
				first: "a",
			},
		},
		{
			name: "stops on wider row",
			fields: fields{
				input: configResult{
					option: &option{}, bodyRows: 2,
					columns: []column{{}},
				},
			},
			args: args{
				sources: [][]any{{"a"}, {"b", "c"}},
			},
			want: want{
				rows:      1,
				first:     "a",
				hasError:  true,
				isColumns: true,
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
			o.prepare()
			o.compileBody(test.args.sources)
			got := want{
				rows:      len(o.output.body),
				hasError:  o.err != nil,
				isColumns: errors.Is(o.err, table.ErrColumnCount),
			}
			if len(o.output.body) > 0 {
				got.first = o.output.body[0].cells[0].value
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
		rows      int
		value     string
		last      string
		hasError  bool
		isColumns bool
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "empty",
			fields: fields{
				input: configResult{option: &option{}},
			},
		},
		{
			name: "keeps prior error",
			fields: fields{
				input: configResult{
					option: &option{}, footer: [][]string{{"f"}},
					columns: []column{{}},
				},
				err: testutil.NewError(),
			},
			want: want{
				hasError: true,
			},
		},
		{
			name: "rejects wider footer",
			fields: fields{
				input: configResult{
					option: &option{}, footer: [][]string{{"a", "b"}},
					columns:       []column{{}},
					footerColumns: 2,
				},
			},
			want: want{
				hasError:  true,
				isColumns: true,
			},
		},
		{
			name: "compiles rows",
			fields: fields{
				input: configResult{
					option:        &option{},
					footer:        [][]string{{"f"}, {"f"}},
					columns:       []column{{rowspan: ScopeFooter}},
					footerColumns: 1,
				},
			},
			want: want{
				rows:  2,
				value: "f",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			o := &compiler{
				input: test.fields.input,
				state: &state,
				err:   test.fields.err,
			}
			o.prepare()
			o.compileFooter()
			got := want{
				rows:      len(o.output.footer),
				hasError:  o.err != nil,
				isColumns: errors.Is(o.err, table.ErrColumnCount),
			}
			if len(o.output.footer) > 0 {
				got.value = o.output.footer[0].cells[0].value
				got.last = o.output.footer[len(o.output.footer)-1].cells[0].value
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
		values      []string
		colors      []*Color
		decorations []*Decoration
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "header index and empty label",
			fields: fields{
				input: configResult{
					option: &option{indexOffset: 1},
					header: [][]string{{"", "b"}},
					columns: []column{
						{},
						{transformer: transformer{colors: func() scope.Scopes[*Color] {
							var colors scope.Scopes[*Color]
							colors.Set(ScopeHeader, ColorFgRed)
							return colors
						}()}},
						{transformer: transformer{decorations: func() scope.Scopes[*Decoration] {
							var decorations scope.Scopes[*Decoration]
							decorations.Set(ScopeHeader, DecorationBold)
							return decorations
						}()}},
					},
				},
			},
			args: args{
				labels: []string{"", "b"},
				scope:  ScopeHeader,
			},
			want: want{
				values:      []string{"#", "", "b"},
				colors:      []*Color{nil, nil, nil},
				decorations: []*Decoration{nil, nil, DecorationBold},
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
			r := o.compileBand(test.args.labels, test.args.scope)
			got := want{
				values:      make([]string, len(r.cells)),
				colors:      make([]*Color, len(r.cells)),
				decorations: make([]*Decoration, len(r.cells)),
			}
			for index := range r.cells {
				got.values[index] = r.cells[index].value
				got.colors[index] = r.cells[index].color
				got.decorations[index] = r.cells[index].decoration
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
		values      []string
		colors      []*Color
		decorations []*Decoration
		hasError    bool
		isColumns   bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "index transformer placeholder and missing value",
			fields: fields{
				input: configResult{
					option: &option{
						placeholder: "-",
						indexOffset: 1,
					},
					bodyRows: 1,
					columns: []column{
						{},
						{},
						{transformer: transformer{fn: func(any) (string, *Color, *Decoration) {
							return "two", ColorFgRed, DecorationBold
						}}},
						{},
					},
				},
			},
			args: args{
				source:   []any{"", testutil.PanicStringer{}},
				rowIndex: 2,
			},
			want: want{
				values:      []string{"3", "-", "two", "-"},
				colors:      []*Color{nil, nil, ColorFgRed, nil},
				decorations: []*Decoration{nil, nil, DecorationBold, nil},
			},
		},
		{
			name: "wider row",
			fields: fields{
				input: configResult{
					option: &option{}, bodyRows: 1,
					columns: []column{{}},
				},
			},
			args: args{
				source: []any{"a", "b"},
			},
			want: want{
				hasError:  true,
				isColumns: true,
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
			o.prepare()
			o.compileRow(test.args.source, test.args.rowIndex)
			got := want{
				hasError:  o.err != nil,
				isColumns: errors.Is(o.err, table.ErrColumnCount),
			}
			if len(o.output.body) > 0 {
				got.values = make([]string, len(o.output.body[0].cells))
				got.colors = make([]*Color, len(o.output.body[0].cells))
				got.decorations = make([]*Decoration, len(o.output.body[0].cells))
				for index := range o.output.body[0].cells {
					got.values[index] = o.output.body[0].cells[index].value
					got.colors[index] = o.output.body[0].cells[index].color
					got.decorations[index] = o.output.body[0].cells[index].decoration
				}
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
					columns: make([]column, 2),
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
					columns: make([]column, 2),
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
			r := o.newRow()
			r.cells[0].value = "first"
			got := want{
				cells:       len(state.cells),
				rowCells:    len(r.cells),
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
		cells []cell
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "escapes applies markup and clears spans",
			fields: fields{
				state: compilerState{
					values: []string{"a|b&br;", "bold", "code", "rowspan", "colspan", "color"},
				},
			},
			args: args{
				row: row{
					cells: []cell{
						{},
						{decoration: DecorationBold},
						{color: ColorFgRed, decoration: DecorationCode},
						{color: ColorFgBlue},
						{decoration: DecorationItalic},
						{color: ColorFgRed},
					},
					rowspans: 0b001000,
					colspans: 0b010000,
				},
			},
			want: want{
				cells: []cell{
					{value: `a\\|b\\&br;`, width: 11, size: 11},
					{
						value:      "bold",
						width:      4 + len(DecorationBold.Prefix) + len(DecorationBold.Suffix),
						size:       4 + len(DecorationBold.Prefix) + len(DecorationBold.Suffix),
						decoration: DecorationBold,
					},
					{
						value:      "code",
						width:      4 + len(DecorationCode.Prefix) + len(DecorationCode.Suffix),
						size:       4 + len(DecorationCode.Prefix) + len(DecorationCode.Suffix),
						decoration: DecorationCode,
					},
					{},
					{},
					{
						value: "color",
						width: 5 + len(ColorFgRed.Prefix) + len(ColorFgRed.Suffix),
						size:  5 + len(ColorFgRed.Prefix) + len(ColorFgRed.Suffix),
						color: ColorFgRed,
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
				cells: r.cells,
			}
			testutil.AssertValue(t, got, test.want, "compileCells")
		})
	}
}
