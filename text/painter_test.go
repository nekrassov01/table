package text

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/nekrassov01/table/internal/testutil"
	"github.com/nekrassov01/table/internal/value"
)

func Test_painter_prepare(t *testing.T) {
	type fields struct {
		input   solverResult
		state   painterState
		strings value.Store
		w       io.Writer
		err     error
	}
	type want struct {
		layoutsLen  int
		layoutsCap  int
		segmentsLen int
		backingLen  int
		backingCap  int
		lineLen     int
		lineCap     int
		horizonNil  bool
		horizonLen  int
		horizonCap  int
		rowspans    uint64
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "empty geometry reserves one byte",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: &option{},
						},
					},
				},
				state: painterState{
					layouts: []layout{
						{},
					},
					segments: []segment{
						{
							value: "old",
						},
					},
					rowspans: 7,
				},
			},
			want: want{
				layoutsCap: 1,
				backingLen: 1,
				backingCap: 1,
				lineCap:    1,
				horizonNil: true,
			},
		},
		{
			name: "reserves row and cached horizon",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: &option{
								style: Style{
									Border: BorderStyle{
										Body: &Horizontal{
											Fill: "-",
										},
									},
								},
							},
							columns: make([]column, 2),
						},
					},
					metrics: []columnMetric{
						{
							box: box{
								width: 2,
								lPad:  1,
							},
							overhead: 1,
						},
						{
							box: box{
								width: 2,
							},
						},
					},
				},
			},
			want: want{
				layoutsCap: 2,
				backingLen: 20,
				backingCap: 20,
				lineCap:    10,
				horizonLen: 0,
				horizonCap: 10,
			},
		},
		{
			name: "reuses backing with caption capacity",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: &option{
								caption: "caption",
							},
							columns: []column{{}},
						},
					},
					metrics: []columnMetric{
						{
							box: box{
								width: 1,
							},
						},
					},
				},
				state: painterState{
					layouts:     make([]layout, 1, 2),
					lineBacking: make([]byte, 50),
					horizon:     []byte("old"),
				},
			},
			want: want{
				layoutsCap: 2,
				backingLen: 50,
				backingCap: 50,
				lineCap:    8,
				horizonNil: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			store := test.fields.strings
			o := &painter{
				input:   test.fields.input,
				state:   &state,
				strings: &store,
				w:       test.fields.w,
				err:     test.fields.err,
			}
			o.prepare()
			got := want{
				layoutsLen:  len(state.layouts),
				layoutsCap:  cap(state.layouts),
				segmentsLen: len(state.segments),
				backingLen:  len(state.lineBacking),
				backingCap:  cap(state.lineBacking),
				lineLen:     len(state.line),
				lineCap:     cap(state.line),
				horizonNil:  state.horizon == nil,
				horizonLen:  len(state.horizon),
				horizonCap:  cap(state.horizon),
				rowspans:    state.rowspans,
			}
			testutil.AssertValue(t, got, test.want, "prepare")
		})
	}
}

func Test_painter_paintHeader(t *testing.T) {
	type fields struct {
		input   solverResult
		state   painterState
		strings value.Store
		w       io.Writer
		err     error
	}
	type want struct {
		output string
		err    string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "caption precedes empty header horizon",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: &option{
								style:       StyleASCII,
								caption:     "cap",
								captionSide: CaptionTop,
							},
							columns: []column{{}},
						},
						body: []row{
							{
								bars: allBars,
							},
						},
					},
					metrics: []columnMetric{
						{
							box: box{
								width: 1,
							},
						},
					},
				},
			},
			want: want{
				output: "cap\n+-+\n",
			},
		},
		{
			name: "empty header follows first body bars",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: &option{
								style: StyleLight,
							},
							columns: []column{{}, {}},
						},
						body: []row{
							{
								bars: allBars &^ 0b10,
							},
						},
					},
					metrics: []columnMetric{
						{
							box: box{
								width: 1,
							},
						},
						{
							box: box{
								width: 1,
							},
						},
					},
				},
			},
			want: want{
				output: "┌───┐\n",
			},
		},
		{
			name: "empty header follows first footer bars",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: &option{
								style: StyleLight,
							},
							columns: []column{{}, {}},
						},
						footer: []row{
							{
								bars: allBars &^ 0b10,
							},
						},
					},
					metrics: []columnMetric{
						{
							box: box{
								width: 1,
							},
						},
						{
							box: box{
								width: 1,
							},
						},
					},
				},
			},
			want: want{
				output: "┌───┐\n",
			},
		},
		{
			name: "header separator follows footer bars without body",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: &option{
								style: StyleLight,
							},
							columns: []column{{}, {}},
						},
						header: []row{
							{
								cells: []cell{
									{
										value: "A",
									},
									{
										value: "B",
									},
								},
								bars: allBars,
							},
						},
						footer: []row{
							{
								bars: allBars &^ 0b10,
							},
						},
					},
					metrics: []columnMetric{
						{
							box: box{
								width: 1,
							},
						},
						{
							box: box{
								width: 1,
							},
						},
					},
				},
			},
			want: want{
				output: "┌─┬─┐\n│A│B│\n╞═╧═╡\n",
			},
		},
		{
			name: "multiple header rows use header fallback",
			fields: fields{
				input: func() solverResult {
					style := StyleASCII
					style.Border.Body = nil
					return solverResult{
						compilerResult: compilerResult{
							configResult: configResult{
								option: &option{
									style: style,
								},
								columns: []column{{}},
							},
							header: []row{
								{
									cells: []cell{
										{
											value: "A",
										},
									},
									bars: allBars,
								},
								{
									cells: []cell{
										{
											value: "B",
										},
									},
									bars: allBars,
								},
							},
							body: []row{
								{
									bars: noBars,
								},
							},
						},
						metrics: []columnMetric{
							{
								box: box{
									width: 1,
								},
							},
						},
					}
				}(),
			},
			want: want{
				output: "+-+\n|A|\n+-+\n|B|\n+-+\n",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			store := test.fields.strings
			var output bytes.Buffer
			writer := test.fields.w
			if writer == nil {
				writer = &output
			}
			o := &painter{
				input:   test.fields.input,
				state:   &state,
				strings: &store,
				w:       writer,
				err:     test.fields.err,
			}
			o.paintHeader()
			got := want{
				output: output.String(),
			}
			if o.err != nil {
				got.err = o.err.Error()
			}
			testutil.AssertValue(t, got, test.want, "paintHeader")
		})
	}
}

func Test_painter_paintBody(t *testing.T) {
	type fields struct {
		input   solverResult
		state   painterState
		strings value.Store
		w       io.Writer
		err     error
	}
	type want struct {
		output string
		err    string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "rows are separated after the first",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: &option{
								style: StyleASCII,
							},
							columns: []column{{}},
						},
						body: []row{
							{
								cells: []cell{
									{
										value: "A",
									},
								},
								bars: allBars,
							},
							{
								cells: []cell{
									{
										value: "B",
									},
								},
								bars: allBars,
							},
						},
						previousBars: allBars,
					},
					metrics: []columnMetric{
						{
							box: box{
								width: 1,
							},
						},
					},
				},
			},
			want: want{
				output: "|A|\n+-+\n|B|\n",
			},
		},
		{
			name: "continued body starts with horizon",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: &option{
								style: StyleASCII,
							},
							columns: []column{{}},
						},
						body: []row{
							{
								cells: []cell{
									{
										value: "A",
									},
								},
								bars: allBars,
							},
						},
						previousBars:    allBars,
						hasPreviousBody: true,
					},
					metrics: []columnMetric{
						{
							box: box{
								width: 1,
							},
						},
					},
				},
			},
			want: want{
				output: "+-+\n|A|\n",
			},
		},
		{
			name: "sticky error stops before rows",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: &option{
								style: StyleASCII,
							},
							columns: []column{{}},
						},
						body: []row{
							{
								cells: []cell{
									{
										value: "A",
									},
								},
								bars: allBars,
							},
							{
								cells: []cell{
									{
										value: "B",
									},
								},
								bars: allBars,
							},
						},
						previousBars: allBars,
					},
					metrics: []columnMetric{
						{
							box: box{
								width: 1,
							},
						},
					},
				},
				err: testutil.NewError(),
			},
			want: want{
				err: "test",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			store := test.fields.strings
			var output bytes.Buffer
			writer := test.fields.w
			if writer == nil {
				writer = &output
			}
			o := &painter{
				input:   test.fields.input,
				state:   &state,
				strings: &store,
				w:       writer,
				err:     test.fields.err,
			}
			o.paintBody()
			got := want{
				output: output.String(),
			}
			if o.err != nil {
				got.err = o.err.Error()
			}
			testutil.AssertValue(t, got, test.want, "paintBody")
		})
	}
}

func Test_painter_paintFooter(t *testing.T) {
	type fields struct {
		input   solverResult
		state   painterState
		strings value.Store
		w       io.Writer
		err     error
	}
	type want struct {
		output string
		err    string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "footer without body omits first separator",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: &option{
								style:   StyleASCII,
								caption: "cap",
							},
							columns: []column{{}},
						},
						footer: []row{
							{
								cells: []cell{
									{
										value: "A",
									},
								},
								bars: allBars,
							},
							{
								cells: []cell{
									{
										value: "B",
									},
								},
								bars: allBars,
							},
						},
						lastBars: allBars,
					},
					metrics: []columnMetric{
						{
							box: box{
								width: 1,
							},
						},
					},
				},
			},
			want: want{
				output: "|A|\n+-+\n|B|\n+-+\ncap\n",
			},
		},
		{
			name: "footer after body starts with separator",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: &option{
								style:       StyleASCII,
								caption:     "cap",
								captionSide: CaptionTop,
							},
							columns: []column{{}},
						},
						footer: []row{
							{
								cells: []cell{
									{
										value: "A",
									},
								},
								bars: allBars,
							},
						},
						lastBars:        allBars,
						hasPreviousBody: true,
					},
					metrics: []columnMetric{
						{
							box: box{
								width: 1,
							},
						},
					},
				},
			},
			want: want{
				output: "+-+\n|A|\n+-+\n",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			store := test.fields.strings
			var output bytes.Buffer
			writer := test.fields.w
			if writer == nil {
				writer = &output
			}
			o := &painter{
				input:   test.fields.input,
				state:   &state,
				strings: &store,
				w:       writer,
				err:     test.fields.err,
			}
			o.paintFooter()
			got := want{
				output: output.String(),
			}
			if o.err != nil {
				got.err = o.err.Error()
			}
			testutil.AssertValue(t, got, test.want, "paintFooter")
		})
	}
}

func Test_painter_paintCaption(t *testing.T) {
	type fields struct {
		input   solverResult
		state   painterState
		strings value.Store
		w       io.Writer
		err     error
	}
	type want struct {
		output string
		err    string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "empty caption",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: &option{},
						},
					},
				},
			},
			want: want{},
		},
		{
			name: "styled caption",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: &option{
								caption: "caption",
								style: Style{
									Content: ContentStyle{
										Caption: NewAttr(CodeBold),
									},
								},
							},
						},
					},
				},
			},
			want: want{
				output: func() string {
					attr := NewAttr(CodeBold)
					return string(attr.Prefix) + "caption" + string(attr.Suffix) + "\n"
				}(),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			store := test.fields.strings
			var output bytes.Buffer
			writer := test.fields.w
			if writer == nil {
				writer = &output
			}
			o := &painter{
				input:   test.fields.input,
				state:   &state,
				strings: &store,
				w:       writer,
				err:     test.fields.err,
			}
			o.paintCaption()
			got := want{
				output: output.String(),
			}
			if o.err != nil {
				got.err = o.err.Error()
			}
			testutil.AssertValue(t, got, test.want, "paintCaption")
		})
	}
}

func Test_painter_paintRow(t *testing.T) {
	type fields struct {
		input   solverResult
		state   painterState
		strings value.Store
		w       io.Writer
		err     error
	}
	type args struct {
		row   row
		scope Scope
	}
	type want struct {
		output string
		err    string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "multiline row paints every physical line",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: &option{
								style: StyleASCII,
							},
							columns: []column{{}, {}},
						},
					},
					metrics: []columnMetric{
						{
							box: box{
								width: 1,
							},
						},
						{
							box: box{
								width: 1,
							},
						},
					},
				},
			},
			args: args{
				row: row{
					cells: []cell{
						{
							value: "A\nB",
						},
						{
							value: "X",
							attr:  NewAttr(CodeUnderline),
						},
					},
					bars: allBars,
				},
				scope: ScopeBody,
			},
			want: want{
				output: func() string {
					attr := NewAttr(CodeUnderline)
					return "|A|" + string(attr.Prefix) + "X" + string(attr.Suffix) + "|\n|B| |\n"
				}(),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			store := test.fields.strings
			var output bytes.Buffer
			writer := test.fields.w
			if writer == nil {
				writer = &output
			}
			o := &painter{
				input:   test.fields.input,
				state:   &state,
				strings: &store,
				w:       writer,
				err:     test.fields.err,
			}
			r := test.args.row
			o.paintRow(&r, test.args.scope)
			got := want{
				output: output.String(),
			}
			if o.err != nil {
				got.err = o.err.Error()
			}
			testutil.AssertValue(t, got, test.want, "paintRow")
		})
	}
}

func Test_painter_paintCell(t *testing.T) {
	type fields struct {
		state painterState
	}
	type args struct {
		cell layout
		line int
		attr *Attr
	}
	type want struct {
		line string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "selected segment",
			args: args{
				cell: layout{
					segments: []segment{
						{
							value: "a",
							width: 1,
						},
						{
							value: "b",
							width: 1,
						},
					},
					box: box{
						width: 2,
					},
				},
				line: 1,
			},
			want: want{
				line: "b ",
			},
		},
		{
			name: "single line with cell attribute",
			args: args{
				cell: layout{
					value: "x",
					width: 1,
					box: box{
						width: 1,
					},
					attr: NewAttr(CodeUnderline),
				},
				attr: NewAttr(CodeBold),
			},
			want: want{
				line: func() string {
					attr := NewAttr(CodeUnderline)
					return string(attr.Prefix) + "x" + string(attr.Suffix)
				}(),
			},
		},
		{
			name: "missing segment paints padding",
			args: args{
				cell: layout{
					segments: []segment{
						{
							value: "x",
							width: 1,
						},
					},
					box: box{
						width: 2,
					},
				},
				line: 1,
			},
			want: want{
				line: "  ",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			o := &painter{
				state: &state,
			}
			cell := test.args.cell
			o.paintCell(&cell, test.args.line, test.args.attr)
			got := want{
				line: string(state.line),
			}
			testutil.AssertValue(t, got, test.want, "paintCell")
		})
	}
}

func Test_painter_paintSegment(t *testing.T) {
	type fields struct {
		input   solverResult
		state   painterState
		strings value.Store
		w       io.Writer
		err     error
	}
	type args struct {
		box     box
		segment segment
		align   AlignSide
		attr    *Attr
	}
	type want struct {
		line string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "left alignment",
			args: args{
				box: box{
					width: 5,
					lPad:  1,
					rPad:  1,
				},
				segment: segment{
					value: "x",
					width: 1,
				},
				align: AlignLeft,
			},
			want: want{
				line: " x     ",
			},
		},
		{
			name: "right alignment",
			args: args{
				box: box{
					width: 5,
					lPad:  1,
					rPad:  1,
				},
				segment: segment{
					value: "x",
					width: 1,
				},
				align: AlignRight,
			},
			want: want{
				line: "     x ",
			},
		},
		{
			name: "center alignment with attribute",
			args: args{
				box: box{
					width: 5,
					lPad:  1,
					rPad:  1,
				},
				segment: segment{
					value: "x",
					width: 1,
				},
				align: AlignCenter,
				attr:  NewAttr(CodeUnderline),
			},
			want: want{
				line: func() string {
					attr := NewAttr(CodeUnderline)
					return "   " + string(attr.Prefix) + "x" + string(attr.Suffix) + "   "
				}(),
			},
		},
		{
			name: "overflow has no alignment padding",
			args: args{
				box: box{
					width: 3,
					lPad:  1,
					rPad:  1,
				},
				segment: segment{
					value: "value",
					width: 5,
				},
				align: AlignCenter,
			},
			want: want{
				line: " value ",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			store := test.fields.strings
			o := &painter{
				input:   test.fields.input,
				state:   &state,
				strings: &store,
				w:       test.fields.w,
				err:     test.fields.err,
			}
			cellBox := test.args.box
			o.paintSegment(&cellBox, test.args.segment, test.args.align, test.args.attr)
			got := want{
				line: string(state.line),
			}
			testutil.AssertValue(t, got, test.want, "paintSegment")
		})
	}
}

func Test_painter_paintRowHorizon(t *testing.T) {
	type fields struct {
		input   solverResult
		state   painterState
		strings value.Store
		w       io.Writer
		err     error
	}
	type args struct {
		rowspans uint64
		upBars   uint64
		downBars uint64
	}
	type want struct {
		output   string
		horizon  []byte
		rowspans uint64
		err      string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "same bars without body border",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: &option{},
						},
					},
				},
			},
			args: args{
				upBars:   allBars,
				downBars: allBars,
			},
			want: want{},
		},
		{
			name: "different bars without fallback border",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: &option{},
						},
					},
				},
			},
			args: args{
				upBars:   allBars,
				downBars: noBars,
			},
			want: want{},
		},
		{
			name: "compact full rowspan omits horizon",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: &option{
								style:   StyleASCII,
								compact: true,
							},
						},
						rowspanMask: 1,
					},
				},
			},
			args: args{
				rowspans: 1,
				upBars:   allBars,
				downBars: allBars,
			},
			want: want{},
		},
		{
			name: "header border is fallback",
			fields: fields{
				input: func() solverResult {
					style := StyleASCII
					style.Border.Body = nil
					return solverResult{
						compilerResult: compilerResult{
							configResult: configResult{
								option: &option{
									style: style,
								},
							},
						},
						metrics: []columnMetric{
							{
								box: box{
									width: 1,
								},
							},
						},
					}
				}(),
			},
			args: args{
				upBars:   allBars,
				downBars: noBars,
			},
			want: want{
				output: "+-+\n",
			},
		},
		{
			name: "cached horizon is reused",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: &option{
								style: StyleASCII,
							},
						},
					},
					metrics: []columnMetric{
						{
							box: box{
								width: 1,
							},
						},
					},
				},
				state: painterState{
					horizon:  []byte("cached\n"),
					rowspans: 1,
				},
			},
			args: args{
				rowspans: 1,
				upBars:   allBars,
				downBars: allBars,
			},
			want: want{
				output:   "cached\n",
				horizon:  []byte("cached\n"),
				rowspans: 1,
			},
		},
		{
			name: "new cache stores painted horizon",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: &option{
								style: StyleASCII,
							},
						},
					},
					metrics: []columnMetric{
						{
							box: box{
								width: 1,
							},
						},
					},
				},
				state: painterState{
					horizon: make([]byte, 0, 8),
				},
			},
			args: args{
				upBars:   allBars,
				downBars: allBars,
			},
			want: want{
				output:  "+-+\n",
				horizon: []byte("+-+\n"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			store := test.fields.strings
			var output bytes.Buffer
			writer := test.fields.w
			if writer == nil {
				writer = &output
			}
			o := &painter{
				input:   test.fields.input,
				state:   &state,
				strings: &store,
				w:       writer,
				err:     test.fields.err,
			}
			o.paintRowHorizon(test.args.rowspans, test.args.upBars, test.args.downBars)
			got := want{
				output:   output.String(),
				horizon:  state.horizon,
				rowspans: state.rowspans,
			}
			if o.err != nil {
				got.err = o.err.Error()
			}
			testutil.AssertValue(t, got, test.want, "paintRowHorizon")
		})
	}
}

func Test_painter_paintHorizon(t *testing.T) {
	type fields struct {
		input   solverResult
		state   painterState
		strings value.Store
		w       io.Writer
		err     error
	}
	type args struct {
		horizontal *Horizontal
		rowspans   uint64
		upBars     uint64
		downBars   uint64
	}
	type want struct {
		output string
		err    string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "nil horizontal",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: &option{},
						},
					},
				},
			},
			want: want{},
		},
		{
			name: "rowspan leaves horizontal gap",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: &option{},
						},
					},
					metrics: []columnMetric{
						{
							box: box{
								width: 2,
							},
						},
						{
							box: box{
								width: 2,
							},
						},
					},
				},
			},
			args: args{
				horizontal: StyleASCII.Border.Top,
				rowspans:   0b10,
				upBars:     allBars,
				downBars:   allBars,
			},
			want: want{
				output: "+--+  |\n",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			store := test.fields.strings
			var output bytes.Buffer
			writer := test.fields.w
			if writer == nil {
				writer = &output
			}
			o := &painter{
				input:   test.fields.input,
				state:   &state,
				strings: &store,
				w:       writer,
				err:     test.fields.err,
			}
			o.paintHorizon(
				test.args.horizontal,
				test.args.rowspans,
				test.args.upBars,
				test.args.downBars,
			)
			got := want{
				output: output.String(),
			}
			if o.err != nil {
				got.err = o.err.Error()
			}
			testutil.AssertValue(t, got, test.want, "paintHorizon")
		})
	}
}

func Test_painter_layoutRow(t *testing.T) {
	type fields struct {
		input   solverResult
		state   painterState
		strings value.Store
		w       io.Writer
		err     error
	}
	type args struct {
		row   row
		scope Scope
	}
	type want struct {
		height   int
		layouts  []layout
		segments []segment
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "body index and configured alignment",
			fields: fields{
				input: func() solverResult {
					columns := make([]column, 2)
					columns[1].aligns.Set(ScopeBody, AlignCenter)
					return solverResult{
						compilerResult: compilerResult{
							configResult: configResult{
								option: &option{
									indexOffset: 1,
								},
								columns: columns,
							},
						},
						metrics: []columnMetric{
							{
								box: box{
									width: 1,
								},
							},
							{
								box: box{
									width: 5,
								},
							},
						},
					}
				}(),
			},
			args: args{
				row: row{
					cells: []cell{
						{
							value: "1",
						},
						{
							value: "x",
						},
					},
					bars: allBars,
				},
				scope: ScopeBody,
			},
			want: want{
				height: 1,
				layouts: []layout{
					{
						value: "1",
						box: box{
							width: 1,
						},
						align: AlignRight,
						width: 1,
					},
					{
						value: "x",
						box: box{
							width: 5,
						},
						align: AlignCenter,
						width: 1,
					},
				},
			},
		},
		{
			name: "header defaults to center",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option:  &option{},
							columns: []column{{}},
						},
					},
					metrics: []columnMetric{
						{
							box: box{
								width: 3,
							},
						},
					},
				},
			},
			args: args{
				row: row{
					cells: []cell{
						{
							value: "H",
						},
					},
					bars: allBars,
				},
				scope: ScopeHeader,
			},
			want: want{
				height: 1,
				layouts: []layout{
					{
						value: "H",
						box: box{
							width: 3,
						},
						align: AlignCenter,
						width: 1,
					},
				},
			},
		},
		{
			name: "colspan combines adjacent boxes",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option:  &option{},
							columns: make([]column, 3),
						},
					},
					metrics: []columnMetric{
						{
							box: box{
								width: 2,
								lPad:  1,
								rPad:  1,
							},
							limit: 2,
						},
						{
							box: box{
								offset: 5,
								width:  2,
								lPad:   1,
								rPad:   1,
							},
							limit: 2,
						},
						{
							box: box{
								offset: 10,
								width:  2,
								lPad:   1,
								rPad:   1,
							},
							limit: 2,
						},
					},
				},
			},
			args: args{
				row: row{
					cells: []cell{
						{
							value: "wide",
							attr:  NewAttr(CodeBold),
						},
						{},
						{
							value: "z",
						},
					},
					colspans: 0b010,
					bars:     allBars &^ 0b010,
				},
				scope: ScopeBody,
			},
			want: want{
				height: 1,
				layouts: []layout{
					{
						value: "wide",
						attr:  NewAttr(CodeBold),
						box: box{
							width: 7,
							lPad:  1,
							rPad:  1,
						},
						width: 4,
					},
					{
						value: "z",
						box: box{
							offset: 10,
							width:  2,
							lPad:   1,
							rPad:   1,
						},
						width: 1,
					},
				},
			},
		},
		{
			name: "rowspan clears value and attribute",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option:  &option{},
							columns: []column{{}},
						},
					},
					metrics: []columnMetric{
						{
							box: box{
								width: 3,
							},
						},
					},
				},
			},
			args: args{
				row: row{
					cells: []cell{
						{
							value: "value",
							attr:  NewAttr(CodeBold),
						},
					},
					rowspans: 1,
					bars:     allBars,
				},
				scope: ScopeBody,
			},
			want: want{
				height: 1,
				layouts: []layout{
					{
						box: box{
							width: 3,
						},
					},
				},
			},
		},
		{
			name: "multiline cell determines row height",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option:  &option{},
							columns: []column{{}},
						},
					},
					metrics: []columnMetric{
						{
							box: box{
								width: 1,
							},
						},
					},
				},
			},
			args: args{
				row: row{
					cells: []cell{
						{
							value: "a\nb",
						},
					},
					bars: allBars,
				},
				scope: ScopeBody,
			},
			want: want{
				height: 2,
				layouts: []layout{
					{
						value: "a\nb",
						segments: []segment{
							{
								value: "a",
								width: 1,
							},
							{
								value: "b",
								width: 1,
							},
						},
						box: box{
							width: 1,
						},
						width: 1,
					},
				},
				segments: []segment{
					{
						value: "a",
						width: 1,
					},
					{
						value: "b",
						width: 1,
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			store := test.fields.strings
			o := &painter{
				input:   test.fields.input,
				state:   &state,
				strings: &store,
				w:       test.fields.w,
				err:     test.fields.err,
			}
			r := test.args.row
			height := o.layoutRow(&r, test.args.scope)
			got := want{
				height:   height,
				layouts:  state.layouts,
				segments: state.segments,
			}
			testutil.AssertValue(t, got, test.want, "layoutRow")
		})
	}
}

func Test_painter_layoutCell(t *testing.T) {
	type fields struct {
		input   solverResult
		state   painterState
		strings value.Store
		w       io.Writer
		err     error
	}
	type args struct {
		value    string
		limit    int
		truncate bool
	}
	type want struct {
		layout     layout
		segments   []segment
		stringMark int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "single line without limit",
			args: args{
				value: "界a",
			},
			want: want{
				layout: layout{
					value: "界a",
					width: 3,
				},
			},
		},
		{
			name: "single line truncates before layout",
			args: args{
				value:    "abcdef",
				limit:    5,
				truncate: true,
			},
			want: want{
				layout: layout{
					value: "ab...",
					width: 5,
				},
				stringMark: 5,
			},
		},
		{
			name: "multiline appends physical segments",
			fields: fields{
				state: painterState{
					segments: []segment{
						{
							value: "old",
							width: 3,
						},
					},
				},
			},
			args: args{
				value: "a\n界",
			},
			want: want{
				layout: layout{
					value: "a\n界",
					segments: []segment{
						{
							value: "a",
							width: 1,
						},
						{
							value: "界",
							width: 2,
						},
					},
					width: 2,
				},
				segments: []segment{
					{
						value: "old",
						width: 3,
					},
					{
						value: "a",
						width: 1,
					},
					{
						value: "界",
						width: 2,
					},
				},
			},
		},
		{
			name: "long line wraps by display width",
			args: args{
				value: "a界b",
				limit: 2,
			},
			want: want{
				layout: layout{
					value: "a界b",
					segments: []segment{
						{
							value: "a",
							width: 1,
						},
						{
							value: "界",
							width: 2,
						},
						{
							value: "b",
							width: 1,
						},
					},
					width: 2,
				},
				segments: []segment{
					{
						value: "a",
						width: 1,
					},
					{
						value: "界",
						width: 2,
					},
					{
						value: "b",
						width: 1,
					},
				},
			},
		},
		{
			name: "multiline truncates each physical line",
			args: args{
				value:    "abcdef\nxy",
				limit:    4,
				truncate: true,
			},
			want: want{
				layout: layout{
					value: "abcdef\nxy",
					segments: []segment{
						{
							value: "a...",
							width: 4,
						},
						{
							value: "xy",
							width: 2,
						},
					},
					width: 4,
				},
				segments: []segment{
					{
						value: "a...",
						width: 4,
					},
					{
						value: "xy",
						width: 2,
					},
				},
				stringMark: 4,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			store := test.fields.strings
			o := &painter{
				input:   test.fields.input,
				state:   &state,
				strings: &store,
				w:       test.fields.w,
				err:     test.fields.err,
			}
			got := want{
				layout:     o.layoutCell(test.args.value, test.args.limit, test.args.truncate),
				segments:   state.segments,
				stringMark: store.Mark(),
			}
			testutil.AssertValue(t, got, test.want, "layoutCell")
		})
	}
}

func Test_painter_wrapLine(t *testing.T) {
	type fields struct {
		input   solverResult
		state   painterState
		strings value.Store
		w       io.Writer
		err     error
	}
	type args struct {
		line  string
		limit int
	}
	type want struct {
		width    int
		segments []segment
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "splits ASCII and wide runes",
			args: args{
				line:  "a界b",
				limit: 2,
			},
			want: want{
				width: 2,
				segments: []segment{
					{
						value: "a",
						width: 1,
					},
					{
						value: "界",
						width: 2,
					},
					{
						value: "b",
						width: 1,
					},
				},
			},
		},
		{
			name: "rune wider than limit stays intact",
			args: args{
				line:  "界",
				limit: 1,
			},
			want: want{
				width: 2,
				segments: []segment{
					{
						value: "界",
						width: 2,
					},
				},
			},
		},
		{
			name: "keeps grapheme cluster intact",
			args: args{
				line:  "👩‍💻x",
				limit: 2,
			},
			want: want{
				width: 2,
				segments: []segment{
					{
						value: "👩‍💻",
						width: 2,
					},
					{
						value: "x",
						width: 1,
					},
				},
			},
		},
		{
			name: "empty line",
			args: args{
				limit: 3,
			},
			want: want{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			store := test.fields.strings
			o := &painter{
				input:   test.fields.input,
				state:   &state,
				strings: &store,
				w:       test.fields.w,
				err:     test.fields.err,
			}
			got := want{
				width:    o.wrapLine(test.args.line, test.args.limit),
				segments: state.segments,
			}
			testutil.AssertValue(t, got, test.want, "wrapLine")
		})
	}
}

func Test_painter_truncateLine(t *testing.T) {
	type fields struct {
		input   solverResult
		state   painterState
		strings value.Store
		w       io.Writer
		err     error
	}
	type args struct {
		line      string
		lineWidth int
		limit     int
	}
	type want struct {
		line       string
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
			name: "line within limit",
			args: args{
				line:      "abc",
				lineWidth: 3,
				limit:     3,
			},
			want: want{
				line:  "abc",
				width: 3,
			},
		},
		{
			name: "limit shorter than ellipsis",
			args: args{
				line:      "abcdef",
				lineWidth: 6,
				limit:     2,
			},
			want: want{
				line:  "..",
				width: 2,
			},
		},
		{
			name: "keeps content before ellipsis",
			args: args{
				line:      "abcdef",
				lineWidth: 6,
				limit:     5,
			},
			want: want{
				line:       "ab...",
				width:      5,
				stringMark: 5,
			},
		},
		{
			name: "does not split a wide rune",
			args: args{
				line:      "界abc",
				lineWidth: 5,
				limit:     4,
			},
			want: want{
				line:       "...",
				width:      3,
				stringMark: 3,
			},
		},
		{
			name: "does not split a grapheme cluster",
			args: args{
				line:      "a👩‍💻bcdef",
				lineWidth: 8,
				limit:     6,
			},
			want: want{
				line:       "a👩‍💻...",
				width:      6,
				stringMark: 15,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			store := test.fields.strings
			o := &painter{
				input:   test.fields.input,
				state:   &state,
				strings: &store,
				w:       test.fields.w,
				err:     test.fields.err,
			}
			line, lineWidth := o.truncateLine(
				test.args.line,
				test.args.lineWidth,
				test.args.limit,
			)
			got := want{
				line:       line,
				width:      lineWidth,
				stringMark: store.Mark(),
			}
			testutil.AssertValue(t, got, test.want, "truncateLine")
		})
	}
}

func Test_painter_resetLine(t *testing.T) {
	type fields struct {
		input   solverResult
		state   painterState
		strings value.Store
		w       io.Writer
		err     error
	}
	type want struct {
		lineLen int
		lineCap int
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "retains line capacity",
			fields: fields{
				state: painterState{
					line: make([]byte, 4, 12),
				},
			},
			want: want{
				lineCap: 12,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			store := test.fields.strings
			o := &painter{
				input:   test.fields.input,
				state:   &state,
				strings: &store,
				w:       test.fields.w,
				err:     test.fields.err,
			}
			o.resetLine()
			got := want{
				lineLen: len(state.line),
				lineCap: cap(state.line),
			}
			testutil.AssertValue(t, got, test.want, "resetLine")
		})
	}
}

func Test_painter_flushLine(t *testing.T) {
	type fields struct {
		input   solverResult
		state   painterState
		strings value.Store
		w       io.Writer
		err     error
	}
	type args struct {
		line []byte
	}
	type want struct {
		output string
		err    string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "sticky error skips write",
			fields: fields{
				err: testutil.NewError(),
			},
			args: args{
				line: []byte("value"),
			},
			want: want{
				err: "test",
			},
		},
		{
			name: "empty line skips write",
			want: want{},
		},
		{
			name: "writes complete line",
			args: args{
				line: []byte("value\n"),
			},
			want: want{
				output: "value\n",
			},
		},
		{
			name: "wraps writer error",
			fields: fields{
				w: &testutil.ErrorWriter{
					Err: testutil.NewError(),
				},
			},
			args: args{
				line: []byte("value"),
			},
			want: want{
				err: "text: write failed: test",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			store := test.fields.strings
			var output bytes.Buffer
			writer := test.fields.w
			if writer == nil {
				writer = &output
			}
			o := &painter{
				input:   test.fields.input,
				state:   &state,
				strings: &store,
				w:       writer,
				err:     test.fields.err,
			}
			o.flushLine(test.args.line)
			got := want{
				output: output.String(),
			}
			if o.err != nil {
				got.err = o.err.Error()
			}
			testutil.AssertValue(t, got, test.want, "flushLine")
		})
	}
}

func Test_painter_writeGlyph(t *testing.T) {
	type fields struct {
		input   solverResult
		state   painterState
		strings value.Store
		w       io.Writer
		err     error
	}
	type args struct {
		glyph string
	}
	type want struct {
		line string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "plain glyph",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: &option{},
						},
					},
				},
				state: painterState{
					line: []byte(">"),
				},
			},
			args: args{
				glyph: "|",
			},
			want: want{
				line: ">|",
			},
		},
		{
			name: "styled glyph",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: &option{
								style: Style{
									Border: BorderStyle{
										Attr: NewAttr(CodeFaint),
									},
								},
							},
						},
					},
				},
			},
			args: args{
				glyph: "|",
			},
			want: want{
				line: func() string {
					attr := NewAttr(CodeFaint)
					return string(attr.Prefix) + "|" + string(attr.Suffix)
				}(),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			store := test.fields.strings
			o := &painter{
				input:   test.fields.input,
				state:   &state,
				strings: &store,
				w:       test.fields.w,
				err:     test.fields.err,
			}
			o.writeGlyph(test.args.glyph)
			got := want{
				line: string(state.line),
			}
			testutil.AssertValue(t, got, test.want, "writeGlyph")
		})
	}
}

func Test_painter_writeFill(t *testing.T) {
	type fields struct {
		input   solverResult
		state   painterState
		strings value.Store
		w       io.Writer
		err     error
	}
	type args struct {
		fill      string
		fillWidth int
		count     int
	}
	type want struct {
		line string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "empty fill",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: &option{},
						},
					},
				},
				state: painterState{
					line: []byte(">"),
				},
			},
			args: args{
				count: 3,
			},
			want: want{
				line: ">",
			},
		},
		{
			name: "nonpositive count",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: &option{},
						},
					},
				},
				state: painterState{
					line: []byte(">"),
				},
			},
			args: args{
				fill:      "-",
				fillWidth: 1,
			},
			want: want{
				line: ">",
			},
		},
		{
			name: "single byte fill",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: &option{},
						},
					},
				},
				state: painterState{
					line: []byte(">"),
				},
			},
			args: args{
				fill:      "-",
				fillWidth: 1,
				count:     4,
			},
			want: want{
				line: ">----",
			},
		},
		{
			name: "multi-column fill",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: &option{},
						},
					},
				},
			},
			args: args{
				fill:      "ab",
				fillWidth: 2,
				count:     3,
			},
			want: want{
				line: "ab ",
			},
		},
		{
			name: "zero-width fill",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: &option{},
						},
					},
				},
			},
			args: args{
				fill:  "\x00",
				count: 3,
			},
			want: want{
				line: "   ",
			},
		},
		{
			name: "attribute wraps entire fill",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: &option{
								style: Style{
									Border: BorderStyle{
										Attr: NewAttr(CodeFaint),
									},
								},
							},
						},
					},
				},
			},
			args: args{
				fill:      "-",
				fillWidth: 1,
				count:     3,
			},
			want: want{
				line: func() string {
					attr := NewAttr(CodeFaint)
					return string(attr.Prefix) + "---" + string(attr.Suffix)
				}(),
			},
		},
		{
			name: "attribute wraps multi-column fill",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: &option{
								style: Style{
									Border: BorderStyle{
										Attr: NewAttr(CodeFaint),
									},
								},
							},
						},
					},
				},
			},
			args: args{
				fill:      "ab",
				fillWidth: 2,
				count:     3,
			},
			want: want{
				line: func() string {
					attr := NewAttr(CodeFaint)
					return string(attr.Prefix) + "ab " + string(attr.Suffix)
				}(),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			line := make([]byte, len(state.line), 64)
			copy(line, state.line)
			state.line = line
			store := test.fields.strings
			o := &painter{
				input:   test.fields.input,
				state:   &state,
				strings: &store,
				w:       test.fields.w,
				err:     test.fields.err,
			}
			o.writeFill(test.args.fill, test.args.fillWidth, test.args.count)
			got := want{
				line: string(state.line),
			}
			testutil.AssertValue(t, got, test.want, "writeFill")
		})
	}
}

func Test_painter_writeSpaces(t *testing.T) {
	type fields struct {
		input   solverResult
		state   painterState
		strings value.Store
		w       io.Writer
		err     error
	}
	type args struct {
		count int
	}
	type want struct {
		line string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "nonpositive count",
			fields: fields{
				state: painterState{
					line: []byte(">"),
				},
			},
			want: want{
				line: ">",
			},
		},
		{
			name: "small count",
			fields: fields{
				state: painterState{
					line: []byte(">"),
				},
			},
			args: args{
				count: 3,
			},
			want: want{
				line: ">   ",
			},
		},
		{
			name: "count larger than shared block",
			fields: fields{
				state: painterState{
					line: []byte(">"),
				},
			},
			args: args{
				count: 258,
			},
			want: want{
				line: ">" + strings.Repeat(" ", 258),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			store := test.fields.strings
			o := &painter{
				input:   test.fields.input,
				state:   &state,
				strings: &store,
				w:       test.fields.w,
				err:     test.fields.err,
			}
			o.writeSpaces(test.args.count)
			got := want{
				line: string(state.line),
			}
			testutil.AssertValue(t, got, test.want, "writeSpaces")
		})
	}
}

func Test_painter_writeNewline(t *testing.T) {
	type fields struct {
		input   solverResult
		state   painterState
		strings value.Store
		w       io.Writer
		err     error
	}
	type want struct {
		line string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "appends newline",
			fields: fields{
				state: painterState{
					line: []byte("value"),
				},
			},
			want: want{
				line: "value\n",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			store := test.fields.strings
			o := &painter{
				input:   test.fields.input,
				state:   &state,
				strings: &store,
				w:       test.fields.w,
				err:     test.fields.err,
			}
			o.writeNewline()
			got := want{
				line: string(state.line),
			}
			testutil.AssertValue(t, got, test.want, "writeNewline")
		})
	}
}

func Test_painter_writeValue(t *testing.T) {
	type fields struct {
		input   solverResult
		state   painterState
		strings value.Store
		w       io.Writer
		err     error
	}
	type args struct {
		value string
		attr  *Attr
	}
	type want struct {
		line string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "empty value",
			fields: fields{
				state: painterState{
					line: []byte(">"),
				},
			},
			args: args{
				attr: NewAttr(CodeItalic),
			},
			want: want{
				line: ">",
			},
		},
		{
			name: "nil attribute",
			fields: fields{
				state: painterState{
					line: []byte(">"),
				},
			},
			args: args{
				value: "value",
			},
			want: want{
				line: ">value",
			},
		},
		{
			name: "zero attribute",
			fields: fields{
				state: painterState{
					line: []byte(">"),
				},
			},
			args: args{
				value: "value",
				attr:  &Attr{},
			},
			want: want{
				line: ">value",
			},
		},
		{
			name: "styled value",
			fields: fields{
				state: painterState{
					line: []byte(">"),
				},
			},
			args: args{
				value: "value",
				attr:  NewAttr(CodeItalic),
			},
			want: want{
				line: func() string {
					attr := NewAttr(CodeItalic)
					return ">" + string(attr.Prefix) + "value" + string(attr.Suffix)
				}(),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			store := test.fields.strings
			o := &painter{
				input:   test.fields.input,
				state:   &state,
				strings: &store,
				w:       test.fields.w,
				err:     test.fields.err,
			}
			o.writeValue(test.args.value, test.args.attr)
			got := want{
				line: string(state.line),
			}
			testutil.AssertValue(t, got, test.want, "writeValue")
		})
	}
}

func Test_hasRowspan(t *testing.T) {
	type args struct {
		rowspans uint64
		index    int
	}
	type want struct {
		value bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "marked column",
			args: args{
				rowspans: 0b1010,
				index:    3,
			},
			want: want{
				value: true,
			},
		},
		{
			name: "unmarked column",
			args: args{
				rowspans: 0b1010,
				index:    2,
			},
			want: want{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := want{
				value: hasRowspan(test.args.rowspans, test.args.index),
			}
			testutil.AssertValue(t, got, test.want, "hasRowspan")
		})
	}
}
