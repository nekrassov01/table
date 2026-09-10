package text

import (
	"testing"

	"github.com/nekrassov01/table/internal/scope"
	"github.com/nekrassov01/table/internal/span"
	"github.com/nekrassov01/table/internal/testutil"
	"github.com/nekrassov01/table/internal/value"
)

func Test_compiler_prepare(t *testing.T) {
	type fields struct {
		input     configResult
		state     compilerState
		strings   value.Store
		bodyStart int
		err       error
		output    compilerResult
		previous  []string
	}
	type want struct {
		rowsLen        int
		rowsCap        int
		cellsLen       int
		cellsCap       int
		spanValuesLen  int
		spanValuesCap  int
		rowspanHeader  uint64
		rowspanBody    uint64
		rowspanFooter  uint64
		colspanHeader  uint64
		colspanBody    uint64
		colspanFooter  uint64
		outputRowspans uint64
		placeholder    string
		previousSpan   uint64
		lastBars       uint64
		stringMark     int
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "allocates storage and span masks",
			fields: fields{
				input: configResult{
					option: &option{
						placeholder: "\t",
					},
					header:   [][]string{{"header"}},
					footer:   [][]string{{"footer"}},
					bodyRows: 2,
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
				bodyStart: -1,
				previous:  []string{"same"},
			},
			want: want{
				rowsCap:        4,
				cellsCap:       8,
				spanValuesLen:  2,
				spanValuesCap:  2,
				rowspanHeader:  0b01,
				rowspanBody:    0b11,
				rowspanFooter:  0b10,
				colspanHeader:  0b10,
				colspanBody:    0b10,
				colspanFooter:  0b01,
				outputRowspans: 0b11,
				placeholder:    "    ",
				lastBars:       allBars,
				stringMark:     4,
			},
		},
		{
			name: "reuses storage and clears old masks",
			fields: fields{
				input: configResult{
					option:   &option{},
					bodyRows: 1,
					columns:  []columnConfig{{}},
				},
				state: func() compilerState {
					var rowspans scope.Masks
					rowspans.Mark(ScopeHeader|ScopeBody|ScopeFooter, 3)
					var colspans scope.Masks
					colspans.Mark(ScopeHeader|ScopeBody|ScopeFooter, 3)
					return compilerState{
						rows:       make([]row, 2, 6),
						cells:      make([]cell, 3, 12),
						spanValues: make([]string, 3, 12),
						rowspans:   rowspans,
						colspans:   colspans,
					}
				}(),
				bodyStart: -1,
				output: compilerResult{
					rowspanMask: 9,
				},
				previous: []string{"same"},
			},
			want: want{
				rowsCap:       6,
				cellsCap:      12,
				spanValuesLen: 1,
				spanValuesCap: 12,
				lastBars:      allBars,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			if test.fields.previous != nil {
				span.Rowspans(1, test.fields.previous, &state.previousBody)
			}
			strings := test.fields.strings
			o := &compiler{
				input:     test.fields.input,
				state:     &state,
				strings:   &strings,
				bodyStart: test.fields.bodyStart,
				err:       test.fields.err,
				output:    test.fields.output,
			}
			o.prepare()
			previousSpan := uint64(0)
			if test.fields.previous != nil {
				previousSpan = span.Rowspans(1, test.fields.previous, &state.previousBody)
			}
			got := want{
				rowsLen:        len(state.rows),
				rowsCap:        cap(state.rows),
				cellsLen:       len(state.cells),
				cellsCap:       cap(state.cells),
				spanValuesLen:  len(state.spanValues),
				spanValuesCap:  cap(state.spanValues),
				rowspanHeader:  state.rowspans.Resolve(ScopeHeader),
				rowspanBody:    state.rowspans.Resolve(ScopeBody),
				rowspanFooter:  state.rowspans.Resolve(ScopeFooter),
				colspanHeader:  state.colspans.Resolve(ScopeHeader),
				colspanBody:    state.colspans.Resolve(ScopeBody),
				colspanFooter:  state.colspans.Resolve(ScopeFooter),
				outputRowspans: o.output.rowspanMask,
				placeholder:    o.output.placeholder,
				previousSpan:   previousSpan,
				lastBars:       state.lastBars,
				stringMark:     strings.Mark(),
			}
			testutil.AssertValue(t, got, test.want, "prepare")
		})
	}
}

func Test_compiler_compileHeader(t *testing.T) {
	type fields struct {
		input     configResult
		state     compilerState
		strings   value.Store
		bodyStart int
		err       error
		output    compilerResult
	}
	type want struct {
		header []row
		rows   []row
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "no header",
			fields: fields{
				input: configResult{
					option: &option{},
				},
				output: compilerResult{
					header: []row{
						{
							cells: []cell{
								{
									value: "existing",
								},
							},
						},
					},
				},
			},
			want: want{
				header: []row{
					{
						cells: []cell{
							{
								value: "existing",
							},
						},
					},
				},
			},
		},
		{
			name: "retains header order",
			fields: fields{
				input: func() configResult {
					attr := NewAttr(CodeBold)
					configured := defaultColumn()
					configured.transformer.attrs.Set(ScopeHeader, attr)
					return configResult{
						option: &option{},
						header: [][]string{
							{"top"},
							{"bottom"},
						},
						columns: []columnConfig{configured},
					}
				}(),
				state: compilerState{
					cells:      make([]cell, 0, 2),
					spanValues: make([]string, 1),
				},
			},
			want: want{
				header: []row{
					{
						cells: []cell{
							{
								value: "top",
								attr:  NewAttr(CodeBold),
								width: 3,
							},
						},
						bars: allBars,
					},
					{
						cells: []cell{
							{
								value: "bottom",
								attr:  NewAttr(CodeBold),
								width: 6,
							},
						},
						bars: allBars,
					},
				},
				rows: []row{
					{
						cells: []cell{
							{
								value: "top",
								attr:  NewAttr(CodeBold),
								width: 3,
							},
						},
						bars: allBars,
					},
					{
						cells: []cell{
							{
								value: "bottom",
								attr:  NewAttr(CodeBold),
								width: 6,
							},
						},
						bars: allBars,
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			strings := test.fields.strings
			o := &compiler{
				input:     test.fields.input,
				state:     &state,
				strings:   &strings,
				bodyStart: test.fields.bodyStart,
				err:       test.fields.err,
				output:    test.fields.output,
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
		input     configResult
		state     compilerState
		strings   value.Store
		bodyStart int
		err       error
		output    compilerResult
	}
	type args struct {
		sources [][]any
	}
	type want struct {
		values   []string
		rowspans []uint64
		rows     int
		err      string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "compiles all rows",
			fields: fields{
				input: configResult{
					option: &option{
						placeholder: "-",
					},
					columns: []columnConfig{{}},
				},
				state: compilerState{
					cells:      make([]cell, 0, 2),
					spanValues: make([]string, 1),
					rows:       make([]row, 0, 2),
				},
				bodyStart: -1,
				output: compilerResult{
					lastBars: allBars,
				},
			},
			args: args{
				sources: [][]any{
					{"first"},
					{"second"},
				},
			},
			want: want{
				values:   []string{"first", "second"},
				rowspans: []uint64{0, 0},
				rows:     2,
			},
		},
		{
			name: "compares expanded values for spans",
			fields: fields{
				input: configResult{
					option:  &option{},
					columns: []columnConfig{{}},
				},
				state: func() compilerState {
					var rowspans scope.Masks
					rowspans.Mark(ScopeBody, 0)
					return compilerState{
						cells:      make([]cell, 0, 2),
						spanValues: make([]string, 1),
						rows:       make([]row, 0, 2),
						rowspans:   rowspans,
					}
				}(),
				bodyStart: -1,
				output: compilerResult{
					lastBars: allBars,
				},
			},
			args: args{
				sources: [][]any{
					{"a\tb"},
					{"a    b"},
				},
			},
			want: want{
				values:   []string{"a    b", "a    b"},
				rowspans: []uint64{0, 1},
				rows:     2,
			},
		},
		{
			name: "stops after structural error",
			fields: fields{
				input: configResult{
					option: &option{
						placeholder: "-",
					},
					columns: []columnConfig{{}},
				},
				state: compilerState{
					cells:      make([]cell, 0, 3),
					spanValues: make([]string, 1),
					rows:       make([]row, 0, 3),
				},
				bodyStart: -1,
				output: compilerResult{
					lastBars: allBars,
				},
			},
			args: args{
				sources: [][]any{
					{"first"},
					{"too", "wide"},
					{"third"},
				},
			},
			want: want{
				values:   []string{"first"},
				rowspans: []uint64{0},
				rows:     1,
				err:      "text: column count exceeded: got 2, want 1",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			strings := test.fields.strings
			o := &compiler{
				input:     test.fields.input,
				state:     &state,
				strings:   &strings,
				bodyStart: test.fields.bodyStart,
				err:       test.fields.err,
				output:    test.fields.output,
			}
			o.compileBody(test.args.sources)
			got := want{
				rows: len(o.output.body),
			}
			for _, row := range o.output.body {
				got.rowspans = append(got.rowspans, row.rowspans)
				for _, cell := range row.cells {
					got.values = append(got.values, cell.value)
				}
			}
			if o.err != nil {
				got.err = o.err.Error()
			}
			testutil.AssertValue(t, got, test.want, "compileBody")
		})
	}
}

func Test_compiler_compileFooter(t *testing.T) {
	type fields struct {
		input     configResult
		state     compilerState
		strings   value.Store
		bodyStart int
		err       error
		output    compilerResult
	}
	type want struct {
		footer []row
		rows   []row
		err    bool
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "no footer",
			fields: fields{
				input: configResult{
					option: &option{},
				},
				output: compilerResult{
					footer: []row{
						{
							cells: []cell{
								{
									value: "existing",
								},
							},
						},
					},
				},
			},
			want: want{
				footer: []row{
					{
						cells: []cell{
							{
								value: "existing",
							},
						},
					},
				},
			},
		},
		{
			name: "current error skips footer",
			fields: fields{
				input: configResult{
					option: &option{},
					footer: [][]string{{"footer"}},
				},
				err: testutil.NewError(),
				output: compilerResult{
					footer: []row{
						{
							cells: []cell{
								{
									value: "existing",
								},
							},
						},
					},
				},
			},
			want: want{
				footer: []row{
					{
						cells: []cell{
							{
								value: "existing",
							},
						},
					},
				},
				err: true,
			},
		},
		{
			name: "retains footer order",
			fields: fields{
				input: func() configResult {
					attr := NewAttr(CodeUnderline)
					configured := defaultColumn()
					configured.transformer.attrs.Set(ScopeFooter, attr)
					return configResult{
						option: &option{},
						footer: [][]string{
							{"top"},
							{"bottom"},
						},
						columns:       []columnConfig{configured},
						footerColumns: 1,
					}
				}(),
				state: compilerState{
					cells:      make([]cell, 0, 2),
					spanValues: make([]string, 1),
				},
			},
			want: want{
				footer: []row{
					{
						cells: []cell{
							{
								value: "top",
								attr:  NewAttr(CodeUnderline),
								width: 3,
							},
						},
						bars: allBars,
					},
					{
						cells: []cell{
							{
								value: "bottom",
								attr:  NewAttr(CodeUnderline),
								width: 6,
							},
						},
						bars: allBars,
					},
				},
				rows: []row{
					{
						cells: []cell{
							{
								value: "top",
								attr:  NewAttr(CodeUnderline),
								width: 3,
							},
						},
						bars: allBars,
					},
					{
						cells: []cell{
							{
								value: "bottom",
								attr:  NewAttr(CodeUnderline),
								width: 6,
							},
						},
						bars: allBars,
					},
				},
			},
		},
		{
			name: "rejects footer wider than header",
			fields: fields{
				input: configResult{
					option:        &option{},
					footer:        [][]string{{"a", "b"}},
					columns:       []columnConfig{defaultColumn()},
					footerColumns: 2,
				},
			},
			want: want{
				err: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			strings := test.fields.strings
			o := &compiler{
				input:     test.fields.input,
				state:     &state,
				strings:   &strings,
				bodyStart: test.fields.bodyStart,
				err:       test.fields.err,
				output:    test.fields.output,
			}
			o.compileFooter()
			got := want{
				footer: o.output.footer,
				rows:   state.rows,
				err:    o.err != nil,
			}
			testutil.AssertValue(t, got, test.want, "compileFooter")
		})
	}
}

func Test_compiler_compileBand(t *testing.T) {
	type fields struct {
		input     configResult
		state     compilerState
		strings   value.Store
		bodyStart int
		err       error
		output    compilerResult
	}
	type args struct {
		labels []string
		scope  Scope
	}
	type want struct {
		row        row
		stringMark int
	}
	attr := NewAttr(CodeBold)
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "header index and labels",
			fields: fields{
				input: func() configResult {
					columns := make([]columnConfig, 3)
					for i := range columns {
						columns[i].transformer.attrs.Set(ScopeHeader|ScopeFooter, attr)
					}
					return configResult{
						option: &option{
							indexOffset: 1,
						},
						columns: columns,
					}
				}(),
				state: compilerState{
					cells:      make([]cell, 0, 3),
					spanValues: make([]string, 3),
				},
			},
			args: args{
				labels: []string{"la\tbel"},
				scope:  ScopeHeader,
			},
			want: want{
				row: row{
					cells: []cell{
						{
							value: "#",
							attr:  attr,
							width: 1,
						},
						{
							value: "la    bel",
							attr:  attr,
							width: 9,
						},
						{},
					},
				},
				stringMark: 9,
			},
		},
		{
			name: "footer labels omit index",
			fields: fields{
				input: func() configResult {
					columns := make([]columnConfig, 3)
					for i := range columns {
						columns[i].transformer.attrs.Set(ScopeHeader|ScopeFooter, attr)
					}
					return configResult{
						option: &option{
							indexOffset: 1,
						},
						columns: columns,
					}
				}(),
				state: compilerState{
					cells:      make([]cell, 0, 3),
					spanValues: make([]string, 3),
				},
			},
			args: args{
				labels: []string{"to\ttal", ""},
				scope:  ScopeFooter,
			},
			want: want{
				row: row{
					cells: []cell{
						{},
						{
							value: "to    tal",
							attr:  attr,
							width: 9,
						},
						{},
					},
				},
				stringMark: 9,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			strings := test.fields.strings
			o := &compiler{
				input:     test.fields.input,
				state:     &state,
				strings:   &strings,
				bodyStart: test.fields.bodyStart,
				err:       test.fields.err,
				output:    test.fields.output,
			}
			got := want{
				row:        o.compileBand(test.args.labels, test.args.scope),
				stringMark: strings.Mark(),
			}
			testutil.AssertValue(t, got, test.want, "compileBand")
		})
	}
}

func Test_compiler_compileRow(t *testing.T) {
	type fields struct {
		input     configResult
		state     compilerState
		strings   value.Store
		bodyStart int
		err       error
		output    compilerResult
	}
	type args struct {
		source   []any
		rowIndex int
	}
	type want struct {
		values        []string
		attrs         []*Attr
		rowspans      uint64
		colspans      uint64
		bars          uint64
		bodyStart     int
		rows          int
		bodyRows      int
		attrLen       uint32
		lastBars      uint64
		stateLastBars uint64
		stringMark    int
		err           string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "index transformed and missing values",
			fields: fields{
				input: func() configResult {
					configuredAttr := NewAttr(CodeBold)
					dynamicAttr := NewAttr(CodeFgRed)
					dynamicColumn := defaultColumn()
					dynamicColumn.transformer.fn = func(any) (string, *Attr) {
						return "answer", dynamicAttr
					}
					emptyColumn := defaultColumn()
					emptyColumn.transformer.attrs.Set(ScopeBody, configuredAttr)
					return configResult{
						option: &option{
							placeholder: "-",
							indexOffset: 1,
						},
						columns: []columnConfig{
							defaultColumn(),
							dynamicColumn,
							emptyColumn,
							defaultColumn(),
						},
					}
				}(),
				state: compilerState{
					rows: []row{
						{
							cells: []cell{
								{
									value: "existing",
								},
							},
						},
					},
					cells:      make([]cell, 0, 4),
					spanValues: make([]string, 4),
				},
				bodyStart: -1,
				output: compilerResult{
					lastBars: allBars,
				},
			},
			args: args{
				source:   []any{testutil.PanicStringer{}, ""},
				rowIndex: 4,
			},
			want: want{
				values:        []string{"5", "answer", "-", "-"},
				attrs:         []*Attr{nil, NewAttr(CodeFgRed), nil, nil},
				bars:          allBars,
				bodyStart:     1,
				rows:          2,
				bodyRows:      1,
				attrLen:       9,
				lastBars:      allBars,
				stateLastBars: allBars,
				stringMark:    1,
			},
		},
		{
			name: "row wider than header",
			fields: fields{
				input: configResult{
					option:  &option{},
					columns: []columnConfig{{}},
				},
				bodyStart: -1,
			},
			args: args{
				source: []any{"first", "second"},
			},
			want: want{
				bodyStart: -1,
				err:       "text: column count exceeded: got 2, want 1",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			strings := test.fields.strings
			o := &compiler{
				input:     test.fields.input,
				state:     &state,
				strings:   &strings,
				bodyStart: test.fields.bodyStart,
				err:       test.fields.err,
				output:    test.fields.output,
			}
			o.compileRow(test.args.source, test.args.rowIndex)
			got := want{
				bodyStart:     o.bodyStart,
				rows:          len(state.rows),
				bodyRows:      len(o.output.body),
				attrLen:       o.output.attrLen,
				lastBars:      o.output.lastBars,
				stateLastBars: state.lastBars,
				stringMark:    strings.Mark(),
			}
			if len(o.output.body) > 0 {
				r := o.output.body[len(o.output.body)-1]
				for _, cell := range r.cells {
					got.values = append(got.values, cell.value)
					got.attrs = append(got.attrs, cell.attr)
				}
				got.rowspans = r.rowspans
				got.colspans = r.colspans
				got.bars = r.bars
			}
			if o.err != nil {
				got.err = o.err.Error()
			}
			testutil.AssertValue(t, got, test.want, "compileRow")
		})
	}
}

func Test_compiler_compileCells(t *testing.T) {
	type fields struct {
		input   configResult
		strings value.Store
		output  compilerResult
	}
	type args struct {
		row      row
		source   []any
		rowIndex int
	}
	type want struct {
		cells      []cell
		attrLen    uint32
		stringMark int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "transformed value and attribute",
			fields: fields{
				input: func() configResult {
					attr := NewAttr(CodeFgRed)
					configured := defaultColumn()
					configured.transformer.fn = func(any) (string, *Attr) {
						return "ans\twer", attr
					}
					return configResult{
						option:  &option{},
						columns: []columnConfig{configured},
					}
				}(),
			},
			args: args{
				row: row{
					cells: make([]cell, 1),
				},
				source: []any{testutil.PanicStringer{}},
			},
			want: want{
				cells: []cell{
					{
						value: "ans    wer",
						attr:  NewAttr(CodeFgRed),
						width: 10,
					},
				},
				attrLen:    9,
				stringMark: 10,
			},
		},
		{
			name: "plain ignores transformed attribute",
			fields: fields{
				input: func() configResult {
					configured := defaultColumn()
					configured.transformer.attrs.Set(ScopeBody, NewAttr(CodeBold))
					configured.transformer.fn = func(any) (string, *Attr) {
						return "transformed", NewAttr(CodeFgRed)
					}
					return configResult{
						option: &option{
							plain: true,
						},
						columns: []columnConfig{configured},
					}
				}(),
			},
			args: args{
				row: row{
					cells: make([]cell, 1),
				},
				source: []any{"value"},
			},
			want: want{
				cells: []cell{
					{
						value: "transformed",
						attr:  NewAttr(CodeBold),
						width: 11,
					},
				},
			},
		},
		{
			name: "empty transformation retains formatted value",
			fields: fields{
				input: func() configResult {
					configured := defaultColumn()
					configured.transformer.attrs.Set(ScopeBody, NewAttr(CodeBold))
					configured.transformer.fn = func(any) (string, *Attr) {
						return "", nil
					}
					return configResult{
						option:  &option{},
						columns: []columnConfig{configured},
					}
				}(),
			},
			args: args{
				row: row{
					cells: make([]cell, 1),
				},
				source: []any{12},
			},
			want: want{
				cells: []cell{
					{
						value: "12",
						attr:  NewAttr(CodeBold),
						width: 2,
					},
				},
				stringMark: 2,
			},
		},
		{
			name: "empty value clears dynamic attribute",
			fields: fields{
				input: func() configResult {
					attr := NewAttr(CodeFgRed)
					configured := defaultColumn()
					configured.transformer.fn = func(any) (string, *Attr) {
						return "", attr
					}
					return configResult{
						option: &option{
							placeholder: "\t",
						},
						columns: []columnConfig{configured},
					}
				}(),
			},
			args: args{
				row: row{
					cells: make([]cell, 1),
				},
				source: []any{""},
			},
			want: want{
				cells: []cell{
					{
						value: "    ",
						width: 4,
					},
				},
				attrLen:    9,
				stringMark: 4,
			},
		},
		{
			name: "index and missing value",
			fields: fields{
				input: configResult{
					option: &option{
						placeholder: "\t",
						indexOffset: 1,
					},
					columns: []columnConfig{
						defaultColumn(),
						defaultColumn(),
						defaultColumn(),
					},
				},
			},
			args: args{
				row: row{
					cells: make([]cell, 3),
				},
				source:   []any{"va\tlue"},
				rowIndex: 4,
			},
			want: want{
				cells: []cell{
					{value: "5", width: 1},
					{value: "va    lue", width: 9},
					{value: "    ", width: 4},
				},
				stringMark: 14,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			strings := test.fields.strings
			o := &compiler{
				input:   test.fields.input,
				strings: &strings,
				output:  test.fields.output,
			}
			r := test.args.row
			o.compileCells(r, test.args.source, test.args.rowIndex)
			got := want{
				cells:      r.cells,
				attrLen:    o.output.attrLen,
				stringMark: strings.Mark(),
			}
			testutil.AssertValue(t, got, test.want, "compileCells")
		})
	}
}

func Test_compiler_reserveBand(t *testing.T) {
	type fields struct {
		input     configResult
		state     compilerState
		strings   value.Store
		bodyStart int
		err       error
		output    compilerResult
	}
	type args struct {
		rows int
	}
	type want struct {
		cellsLen    int
		rowsLen     int
		reservedLen int
		cellsCap    bool
		rowsCap     bool
		aliases     bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "grows reusable storage",
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
				cellsLen:    1,
				rowsLen:     3,
				reservedLen: 2,
				cellsCap:    true,
				rowsCap:     true,
				aliases:     true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			strings := test.fields.strings
			o := &compiler{
				input:     test.fields.input,
				state:     &state,
				strings:   &strings,
				bodyStart: test.fields.bodyStart,
				err:       test.fields.err,
				output:    test.fields.output,
			}
			rows := o.reserveBand(test.args.rows)
			rows[0].bars = 7
			got := want{
				cellsLen:    len(state.cells),
				rowsLen:     len(state.rows),
				reservedLen: len(rows),
				cellsCap:    cap(state.cells) >= 5,
				rowsCap:     cap(state.rows) >= 3,
				aliases:     state.rows[1].bars == 7,
			}
			testutil.AssertValue(t, got, test.want, "reserveBand")
		})
	}
}

func Test_compiler_newRow(t *testing.T) {
	type fields struct {
		input     configResult
		state     compilerState
		strings   value.Store
		bodyStart int
		err       error
		output    compilerResult
	}
	type want struct {
		cellsLen int
		rowLen   int
		aliases  bool
	}
	attr := NewAttr(CodeItalic)
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "slices the next row from cell storage",
			fields: fields{
				input: configResult{
					columns: make([]columnConfig, 3),
				},
				state: compilerState{
					cells: make([]cell, 1, 4),
				},
			},
			want: want{
				cellsLen: 4,
				rowLen:   3,
				aliases:  true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			strings := test.fields.strings
			o := &compiler{
				input:     test.fields.input,
				state:     &state,
				strings:   &strings,
				bodyStart: test.fields.bodyStart,
				err:       test.fields.err,
				output:    test.fields.output,
			}
			r := o.newRow()
			r.cells[0] = cell{
				value: "value",
				attr:  attr,
			}
			got := want{
				cellsLen: len(state.cells),
				rowLen:   len(r.cells),
				aliases:  state.cells[1].value == "value" && state.cells[1].attr == attr,
			}
			testutil.AssertValue(t, got, test.want, "newRow")
		})
	}
}

func Test_compiler_setSpans(t *testing.T) {
	type fields struct {
		input     configResult
		state     compilerState
		strings   value.Store
		bodyStart int
		err       error
		output    compilerResult
		previous  [][]string
	}
	type args struct {
		row   row
		scope Scope
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
			name: "rowspans take cells before colspans",
			fields: fields{
				state: func() compilerState {
					var spans scope.Masks
					spans.Mark(ScopeBody, 0)
					spans.Mark(ScopeBody, 1)
					return compilerState{
						rowspans: spans,
						colspans: spans,
					}
				}(),
				previous: [][]string{{"same", "same"}},
			},
			args: args{
				row: row{
					cells: []cell{
						{
							value: "same",
						},
						{
							value: "same",
						},
					},
				},
				scope: ScopeBody,
			},
			want: want{
				rowspans: 0b11,
			},
		},
		{
			name: "equal adjacent values colspan",
			fields: fields{
				state: func() compilerState {
					var colspans scope.Masks
					colspans.Mark(ScopeBody, 0)
					colspans.Mark(ScopeBody, 1)
					return compilerState{
						colspans: colspans,
					}
				}(),
			},
			args: args{
				row: row{
					cells: []cell{
						{
							value: "same",
						},
						{
							value: "same",
						},
					},
				},
				scope: ScopeBody,
			},
			want: want{
				colspans: 0b10,
			},
		},
		{
			name: "scope without configured spans",
			fields: fields{
				state: func() compilerState {
					var spans scope.Masks
					spans.Mark(ScopeBody, 0)
					spans.Mark(ScopeBody, 1)
					return compilerState{
						rowspans: spans,
						colspans: spans,
					}
				}(),
			},
			args: args{
				row: row{
					cells: []cell{
						{
							value: "same",
						},
						{
							value: "same",
						},
					},
				},
				scope: ScopeHeader,
			},
			want: want{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			state.spanValues = make([]string, len(test.args.row.cells))
			strings := test.fields.strings
			o := &compiler{
				input:     test.fields.input,
				state:     &state,
				strings:   &strings,
				bodyStart: test.fields.bodyStart,
				err:       test.fields.err,
				output:    test.fields.output,
			}
			var previous span.PreviousRow
			for _, values := range test.fields.previous {
				span.Rowspans(state.rowspans.Resolve(test.args.scope), values, &previous)
			}
			r := test.args.row
			o.setSpans(&r, test.args.scope, &previous)
			got := want{
				rowspans: r.rowspans,
				colspans: r.colspans,
			}
			testutil.AssertValue(t, got, test.want, "setSpans")
		})
	}
}

func Test_compiler_setBars(t *testing.T) {
	type fields struct {
		input     configResult
		state     compilerState
		strings   value.Store
		bodyStart int
		err       error
		output    compilerResult
	}
	type args struct {
		row          row
		previousBars uint64
		scope        Scope
	}
	type want struct {
		bars uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "no configured colspan keeps all bars",
			args: args{
				row: row{
					cells:    make([]cell, 3),
					colspans: 0b110,
				},
				scope: ScopeBody,
			},
			want: want{
				bars: allBars,
			},
		},
		{
			name: "spans clear interior bars",
			fields: fields{
				state: func() compilerState {
					var colspans scope.Masks
					colspans.Mark(ScopeBody, 2)
					colspans.Mark(ScopeBody, 3)
					return compilerState{
						colspans: colspans,
					}
				}(),
			},
			args: args{
				row: row{
					cells:    make([]cell, 5),
					rowspans: 0b00110,
					colspans: 0b01000,
				},
				previousBars: 0b00010,
				scope:        ScopeBody,
			},
			want: want{
				bars: allBars &^ 0b01100,
			},
		},
		{
			name: "rowspan continuation beside a normal cell keeps the bar",
			fields: fields{
				state: func() compilerState {
					var colspans scope.Masks
					colspans.Mark(ScopeBody, 0)
					colspans.Mark(ScopeBody, 1)
					return compilerState{
						colspans: colspans,
					}
				}(),
			},
			args: args{
				row: row{
					cells:    make([]cell, 2),
					rowspans: 0b10,
				},
				previousBars: allBars &^ 0b10,
				scope:        ScopeBody,
			},
			want: want{
				bars: allBars,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			strings := test.fields.strings
			o := &compiler{
				input:     test.fields.input,
				state:     &state,
				strings:   &strings,
				bodyStart: test.fields.bodyStart,
				err:       test.fields.err,
				output:    test.fields.output,
			}
			r := test.args.row
			o.setBars(&r, test.args.previousBars, test.args.scope)
			got := want{
				bars: r.bars,
			}
			testutil.AssertValue(t, got, test.want, "setBars")
		})
	}
}

func Test_compiler_compileValue(t *testing.T) {
	type fields struct {
		prefix string
	}
	type args struct {
		text      string
		fromStore bool
	}
	type want struct {
		text       string
		width      int
		stringMark int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "unchanged",
			args: args{
				text: "plain",
			},
			want: want{
				text:  "plain",
				width: 5,
			},
		},
		{
			name: "consecutive tabs",
			args: args{
				text: "a\t\tb",
			},
			want: want{
				text:       "a        b",
				width:      10,
				stringMark: 10,
			},
		},
		{
			name: "store-backed input",
			fields: fields{
				prefix: "a\tb",
			},
			args: args{
				fromStore: true,
			},
			want: want{
				text:       "a    b",
				width:      6,
				stringMark: 9,
			},
		},
		{
			name: "non-ASCII with tab",
			args: args{
				text: "あ\tb",
			},
			want: want{
				text:       "あ    b",
				stringMark: 8,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var strings value.Store
			strings.AppendString(test.fields.prefix)
			text := test.args.text
			if test.args.fromStore {
				text = strings.Since(0)
			}
			o := &compiler{
				strings: &strings,
			}
			gotText, gotWidth := o.compileValue(text)
			got := want{
				text:       gotText,
				width:      gotWidth,
				stringMark: strings.Mark(),
			}
			testutil.AssertValue(t, got, test.want, "compileValue")
		})
	}
}
