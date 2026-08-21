package html

import (
	"bytes"
	"io"
	"testing"

	"github.com/nekrassov01/table/internal/testutil"
)

func Test_painter_prepare(t *testing.T) {
	type fields struct {
		input solverResult
		state painterState
	}
	type want struct {
		backingCap int
		lineLen    int
		lineCap    int
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "allocates estimated row block",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option:  &option{},
							columns: []column{{}},
						},
						columnSizes: []int{5},
					},
				},
			},
			want: want{
				backingCap: 177,
				lineCap:    177,
			},
		},
		{
			name: "includes the longest section attribute",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: &option{
								attrs: TableAttr{
									Header: SectionAttr{
										Section: Attr{
											Class: "head",
										},
									},
									Body: SectionAttr{
										Section: Attr{
											Class: "longest",
										},
									},
									Footer: SectionAttr{
										Section: Attr{
											Class: "foot",
										},
									},
								},
							},
							columns: []column{{}},
						},
						columnSizes: []int{5},
					},
				},
			},
			want: want{
				backingCap: 184,
				lineCap:    184,
			},
		},
		{
			name: "adopts a grown line",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option:  &option{},
							columns: []column{{}},
						},
						columnSizes: []int{5},
					},
				},
				state: painterState{
					lineBacking: make([]byte, 0, 4),
					line:        make([]byte, 0, 200),
				},
			},
			want: want{
				backingCap: 200,
				lineCap:    200,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			o := &painter{
				input: test.fields.input,
				state: &state,
			}
			o.prepare()
			got := want{
				backingCap: cap(state.lineBacking),
				lineLen:    len(state.line),
				lineCap:    cap(state.line),
			}
			testutil.AssertValue(t, got, test.want, "prepare")
		})
	}
}

func Test_painter_paintHeader(t *testing.T) {
	type fields struct {
		input solverResult
	}
	type want struct {
		output string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "paints table opening",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: &option{},
						},
					},
				},
			},
			want: want{
				output: "<table>\n",
			},
		},
		{
			name: "opens body section",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: &option{},
						},
						body: []row{{}},
					},
				},
			},
			want: want{
				output: "<table>\n  <tbody>\n",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := painterState{}
			output := &bytes.Buffer{}
			o := &painter{
				input: test.fields.input,
				state: &state,
				w:     output,
			}
			o.paintHeader()
			got := want{
				output: output.String(),
			}
			testutil.AssertValue(t, got, test.want, "paintHeader")
		})
	}
}

func Test_painter_paintBody(t *testing.T) {
	type fields struct {
		input solverResult
		err   error
	}
	type want struct {
		output string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "paints body rows",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option:  &option{},
							columns: []column{{}},
						},
						body: []row{
							{
								cells: []cell{{value: "value", colspan: 1}},
							},
						},
					},
				},
			},
			want: want{
				output: "    <tr>\n      <td>value</td>\n    </tr>\n",
			},
		},
		{
			name: "stops on sticky error",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: &option{},
						},
						body: []row{{}},
					},
				},
				err: testutil.NewError(),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := painterState{}
			output := &bytes.Buffer{}
			o := &painter{
				input: test.fields.input,
				state: &state,
				w:     output,
				err:   test.fields.err,
			}
			o.paintBody()
			got := want{
				output: output.String(),
			}
			testutil.AssertValue(t, got, test.want, "paintBody")
		})
	}
}

func Test_painter_paintFooter(t *testing.T) {
	type fields struct {
		input solverResult
	}
	type want struct {
		output string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "closes table without body",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: &option{},
						},
					},
				},
			},
			want: want{
				output: "</table>\n",
			},
		},
		{
			name: "closes previous body and table",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: &option{},
						},
						hasPreviousBody: true,
					},
				},
			},
			want: want{
				output: "  </tbody>\n</table>\n",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := painterState{}
			output := &bytes.Buffer{}
			o := &painter{
				input: test.fields.input,
				state: &state,
				w:     output,
			}
			o.paintFooter()
			got := want{
				output: output.String(),
			}
			testutil.AssertValue(t, got, test.want, "paintFooter")
		})
	}
}

func Test_painter_paintCaption(t *testing.T) {
	type fields struct {
		input solverResult
		line  []byte
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
			name: "paints caption",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: &option{
								caption:     "title",
								captionSide: CaptionBottom,
								attrs: TableAttr{
									Caption: Attr{
										Class: "caption",
									},
								},
							},
						},
						caption: "title",
					},
				},
				line: []byte("<table>\n"),
			},
			want: want{
				line: "<table>\n  <caption class=\"caption\" style=\"caption-side:bottom\">title</caption>\n",
			},
		},
		{
			name: "keeps line without caption",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: &option{},
						},
					},
				},
				line: []byte("before"),
			},
			want: want{
				line: "before",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := painterState{
				line: test.fields.line,
			}
			o := &painter{
				input: test.fields.input,
				state: &state,
			}
			o.paintCaption()
			got := want{
				line: string(state.line),
			}
			testutil.AssertValue(t, got, test.want, "paintCaption")
		})
	}
}

func Test_painter_paintBand(t *testing.T) {
	type fields struct {
		input solverResult
		err   error
	}
	type args struct {
		rows  []row
		scope Scope
	}
	type want struct {
		output string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "paints header band",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option:  &option{},
							columns: []column{{}},
						},
					},
				},
			},
			args: args{
				rows: []row{
					{
						cells: []cell{{value: "label", colspan: 1}},
					},
				},
				scope: ScopeHeader,
			},
			want: want{
				output: "  <thead>\n    <tr>\n      <th>label</th>\n    </tr>\n  </thead>\n",
			},
		},
		{
			name: "omits empty band",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: &option{},
						},
					},
				},
			},
		},
		{
			name: "stops on sticky error",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: &option{},
						},
					},
				},
				err: testutil.NewError(),
			},
			args: args{
				rows:  []row{{}},
				scope: ScopeFooter,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := painterState{}
			output := &bytes.Buffer{}
			o := &painter{
				input: test.fields.input,
				state: &state,
				w:     output,
				err:   test.fields.err,
			}
			o.paintBand(test.args.rows, test.args.scope)
			got := want{
				output: output.String(),
			}
			testutil.AssertValue(t, got, test.want, "paintBand")
		})
	}
}

func Test_painter_paintOpenSection(t *testing.T) {
	type args struct {
		section section
	}
	type want struct {
		output string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "paints opening tag",
			args: args{
				section: section{
					section: element{
						tag: "tbody",
						attr: Attr{
							Class: "body",
						},
					},
				},
			},
			want: want{
				output: "  <tbody class=\"body\">\n",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := painterState{}
			output := &bytes.Buffer{}
			o := &painter{
				state: &state,
				w:     output,
			}
			o.paintOpenSection(test.args.section)
			got := want{
				output: output.String(),
			}
			testutil.AssertValue(t, got, test.want, "paintOpenSection")
		})
	}
}

func Test_painter_paintCloseSection(t *testing.T) {
	type args struct {
		section section
	}
	type want struct {
		output string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "paints closing tag",
			args: args{
				section: section{
					section: element{
						tag: "tbody",
					},
				},
			},
			want: want{
				output: "  </tbody>\n",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := painterState{}
			output := &bytes.Buffer{}
			o := &painter{
				state: &state,
				w:     output,
			}
			o.paintCloseSection(test.args.section)
			got := want{
				output: output.String(),
			}
			testutil.AssertValue(t, got, test.want, "paintCloseSection")
		})
	}
}

func Test_painter_paintRow(t *testing.T) {
	type fields struct {
		input solverResult
	}
	type args struct {
		row     row
		section section
	}
	type want struct {
		output string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "paints visible cells",
			fields: fields{
				input: func() solverResult {
					columns := make([]column, 3)
					columns[2].aligns.Set(ScopeBody, AlignCenter)
					return solverResult{
						compilerResult: compilerResult{
							configResult: configResult{
								option: &option{
									indexOffset: 1,
								},
								columns: columns,
							},
						},
					}
				}(),
			},
			args: args{
				row: row{
					cells: []cell{
						{value: "1", colspan: 1},
						{},
						{value: "value", colspan: 1},
					},
				},
				section: section{
					row: element{
						tag: "tr",
					},
					cell: element{
						tag: "td",
					},
					scope: ScopeBody,
				},
			},
			want: want{
				output: "    <tr>\n      <td style=\"text-align:right\">1</td>\n      <td style=\"text-align:center\">value</td>\n    </tr>\n",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := painterState{}
			output := &bytes.Buffer{}
			o := &painter{
				input: test.fields.input,
				state: &state,
				w:     output,
			}
			o.paintRow(test.args.row, test.args.section)
			got := want{
				output: output.String(),
			}
			testutil.AssertValue(t, got, test.want, "paintRow")
		})
	}
}

func Test_painter_resolveSection(t *testing.T) {
	type fields struct {
		input solverResult
	}
	type args struct {
		scope Scope
	}
	type want struct {
		section section
	}
	option := &option{
		attrs: TableAttr{
			Header: SectionAttr{
				Section: Attr{
					Class: "head",
				},
			},
			Body: SectionAttr{
				Row: Attr{
					Class: "body",
				},
			},
			Footer: SectionAttr{
				Cell: Attr{
					Class: "foot",
				},
			},
		},
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "header",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: option,
						},
					},
				},
			},
			args: args{
				scope: ScopeHeader,
			},
			want: want{
				section: section{
					section: element{
						tag: "thead",
						attr: Attr{
							Class: "head",
						},
					},
					row: element{
						tag: "tr",
					},
					cell: element{
						tag: "th",
					},
					scope: ScopeHeader,
				},
			},
		},
		{
			name: "body",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: option,
						},
					},
				},
			},
			args: args{
				scope: ScopeBody,
			},
			want: want{
				section: section{
					section: element{
						tag: "tbody",
					},
					row: element{
						tag: "tr",
						attr: Attr{
							Class: "body",
						},
					},
					cell: element{
						tag: "td",
					},
					scope: ScopeBody,
				},
			},
		},
		{
			name: "footer",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						configResult: configResult{
							option: option,
						},
					},
				},
			},
			args: args{
				scope: ScopeFooter,
			},
			want: want{
				section: section{
					section: element{
						tag: "tfoot",
					},
					row: element{
						tag: "tr",
					},
					cell: element{
						tag: "td",
						attr: Attr{
							Class: "foot",
						},
					},
					scope: ScopeFooter,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &painter{
				input: test.fields.input,
			}
			got := want{
				section: o.resolveSection(test.args.scope),
			}
			testutil.AssertValue(t, got, test.want, "resolveSection")
		})
	}
}

func Test_painter_resetLine(t *testing.T) {
	type fields struct {
		line []byte
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
			name: "resets length and retains capacity",
			fields: fields{
				line: []byte("line"),
			},
			want: want{
				lineCap: 4,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := painterState{
				line: test.fields.line,
			}
			o := &painter{
				state: &state,
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
		line   []byte
		w      io.Writer
		output *bytes.Buffer
		err    error
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
			name: "writes line",
			fields: fields{
				line:   []byte("line"),
				output: &bytes.Buffer{},
			},
			want: want{
				output: "line",
			},
		},
		{
			name: "wraps writer error",
			fields: fields{
				line: []byte("line"),
				w: &testutil.ErrorWriter{
					Err: testutil.NewError(),
				},
				output: &bytes.Buffer{},
			},
			want: want{
				err: "html: write failed: test",
			},
		},
		{
			name: "keeps sticky error",
			fields: fields{
				line:   []byte("line"),
				output: &bytes.Buffer{},
				err:    testutil.NewError(),
			},
			want: want{
				err: "test",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := painterState{
				line: test.fields.line,
			}
			writer := test.fields.w
			if writer == nil {
				writer = test.fields.output
			}
			o := &painter{
				state: &state,
				w:     writer,
				err:   test.fields.err,
			}
			o.flushLine()
			got := want{
				output: test.fields.output.String(),
			}
			if o.err != nil {
				got.err = o.err.Error()
			}
			testutil.AssertValue(t, got, test.want, "flushLine")
		})
	}
}

func Test_painter_writeOpenTag(t *testing.T) {
	type fields struct {
		line []byte
	}
	type args struct {
		element    element
		extraStyle string
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
			name: "appends tag and attributes",
			fields: fields{
				line: []byte("before"),
			},
			args: args{
				element: element{
					tag: "table",
					attr: Attr{
						Class: "data",
						Style: "width:100%",
					},
				},
				extraStyle: "caption-side:bottom",
			},
			want: want{
				line: `before<table class="data" style="width:100%;caption-side:bottom">`,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := painterState{
				line: test.fields.line,
			}
			o := &painter{
				state: &state,
			}
			o.writeOpenTag(test.args.element, test.args.extraStyle)
			got := want{
				line: string(state.line),
			}
			testutil.AssertValue(t, got, test.want, "writeOpenTag")
		})
	}
}

func Test_painter_writeCloseTag(t *testing.T) {
	type fields struct {
		line []byte
	}
	type args struct {
		tag string
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
			name: "appends closing tag",
			fields: fields{
				line: []byte("before"),
			},
			args: args{
				tag: "table",
			},
			want: want{
				line: "before</table>",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := painterState{
				line: test.fields.line,
			}
			o := &painter{
				state: &state,
			}
			o.writeCloseTag(test.args.tag)
			got := want{
				line: string(state.line),
			}
			testutil.AssertValue(t, got, test.want, "writeCloseTag")
		})
	}
}

func Test_painter_writeOpenCell(t *testing.T) {
	type fields struct {
		line []byte
	}
	type args struct {
		element    element
		columnAttr Attr
		align      AlignSide
		cell       cell
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
			name: "joins attributes and spans",
			fields: fields{
				line: []byte("before"),
			},
			args: args{
				element: element{
					tag: "td",
					attr: Attr{
						Class: "base",
						Style: "padding:0",
					},
				},
				columnAttr: Attr{
					Class: "column",
					Style: "font-weight:bold",
				},
				align: AlignRight,
				cell: cell{
					rowspan: 2,
					colspan: 3,
				},
			},
			want: want{
				line: `before<td class="base column" style="padding:0;font-weight:bold;text-align:right" rowspan="2" colspan="3">`,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := painterState{
				line: test.fields.line,
			}
			o := &painter{
				state: &state,
			}
			o.writeOpenCell(test.args.element, test.args.columnAttr, test.args.align, &test.args.cell)
			got := want{
				line: string(state.line),
			}
			testutil.AssertValue(t, got, test.want, "writeOpenCell")
		})
	}
}

func Test_painter_writeCloseCell(t *testing.T) {
	type fields struct {
		line []byte
	}
	type args struct {
		element element
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
			name: "cell tag",
			fields: fields{
				line: []byte("value"),
			},
			args: args{
				element: element{
					tag: "td",
				},
			},
			want: want{
				line: "value</td>",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := painterState{
				line: test.fields.line,
			}
			o := &painter{
				state: &state,
			}
			o.writeCloseCell(test.args.element)
			got := want{
				line: string(state.line),
			}
			testutil.AssertValue(t, got, test.want, "writeCloseCell")
		})
	}
}

func Test_painter_writeClass(t *testing.T) {
	type fields struct {
		line []byte
	}
	type args struct {
		class       string
		columnClass string
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
			name: "joins classes",
			fields: fields{
				line: []byte("tag"),
			},
			args: args{
				class:       "base",
				columnClass: "column",
			},
			want: want{
				line: `tag class="base column"`,
			},
		},
		{
			name: "appends column class",
			fields: fields{
				line: []byte("tag"),
			},
			args: args{
				columnClass: "column",
			},
			want: want{
				line: `tag class="column"`,
			},
		},
		{
			name: "omits empty attribute",
			fields: fields{
				line: []byte("tag"),
			},
			want: want{
				line: "tag",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := painterState{
				line: test.fields.line,
			}
			o := &painter{
				state: &state,
			}
			o.writeClass(test.args.class, test.args.columnClass)
			got := want{
				line: string(state.line),
			}
			testutil.AssertValue(t, got, test.want, "writeClass")
		})
	}
}

func Test_painter_writeStyle(t *testing.T) {
	type fields struct {
		line []byte
	}
	type args struct {
		style string
		extra string
		align AlignSide
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
			name: "joins declarations",
			fields: fields{
				line: []byte("tag"),
			},
			args: args{
				style: "padding:0",
				extra: "font-weight:bold",
				align: AlignCenter,
			},
			want: want{
				line: `tag style="padding:0;font-weight:bold;text-align:center"`,
			},
		},
		{
			name: "appends alignment",
			fields: fields{
				line: []byte("tag"),
			},
			args: args{
				align: AlignRight,
			},
			want: want{
				line: `tag style="text-align:right"`,
			},
		},
		{
			name: "omits empty attribute",
			fields: fields{
				line: []byte("tag"),
			},
			want: want{
				line: "tag",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := painterState{
				line: test.fields.line,
			}
			o := &painter{
				state: &state,
			}
			o.writeStyle(test.args.style, test.args.extra, test.args.align)
			got := want{
				line: string(state.line),
			}
			testutil.AssertValue(t, got, test.want, "writeStyle")
		})
	}
}

func Test_painter_writeSpan(t *testing.T) {
	type fields struct {
		line []byte
	}
	type args struct {
		name  string
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
			name: "appends span count",
			fields: fields{
				line: []byte("td"),
			},
			args: args{
				name:  ` rowspan="`,
				count: 12,
			},
			want: want{
				line: `td rowspan="12"`,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := painterState{
				line: test.fields.line,
			}
			o := &painter{
				state: &state,
			}
			o.writeSpan(test.args.name, test.args.count)
			got := want{
				line: string(state.line),
			}
			testutil.AssertValue(t, got, test.want, "writeSpan")
		})
	}
}

func Test_painter_writeCellValue(t *testing.T) {
	type fields struct {
		line []byte
	}
	type args struct {
		cell cell
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
			name: "inner markup",
			fields: fields{
				line: []byte("before"),
			},
			args: args{
				cell: cell{
					value:      "value",
					color:      ColorFgRed,
					decoration: DecorationBold,
				},
			},
			want: want{
				line: `before<strong><span style="color:red">value</span></strong>`,
			},
		},
		{
			name: "block decoration contains color",
			fields: fields{
				line: []byte("before"),
			},
			args: args{
				cell: cell{
					value:      "value",
					color:      ColorFgRed,
					decoration: DecorationPreformatted,
				},
			},
			want: want{
				line: `before<pre><span style="color:red">value</span></pre>`,
			},
		},
		{
			name: "empty value",
			fields: fields{
				line: []byte("before"),
			},
			want: want{
				line: "before",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := painterState{
				line: test.fields.line,
			}
			o := &painter{
				state: &state,
			}
			o.writeCellValue(&test.args.cell)
			got := want{
				line: string(state.line),
			}
			testutil.AssertValue(t, got, test.want, "writeCellValue")
		})
	}
}

func Test_painter_writeIndent(t *testing.T) {
	type fields struct {
		line []byte
	}
	type args struct {
		level int
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
			name: "appends indentation",
			fields: fields{
				line: []byte("before"),
			},
			args: args{
				level: 2,
			},
			want: want{
				line: "before    ",
			},
		},
		{
			name: "zero appends nothing",
			fields: fields{
				line: []byte("before"),
			},
			want: want{
				line: "before",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := painterState{
				line: test.fields.line,
			}
			o := &painter{
				state: &state,
			}
			o.writeIndent(test.args.level)
			got := want{
				line: string(state.line),
			}
			testutil.AssertValue(t, got, test.want, "writeIndent")
		})
	}
}

func Test_painter_writeNewline(t *testing.T) {
	type fields struct {
		line []byte
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
				line: []byte("before"),
			},
			want: want{
				line: "before\n",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := painterState{
				line: test.fields.line,
			}
			o := &painter{
				state: &state,
			}
			o.writeNewline()
			got := want{
				line: string(state.line),
			}
			testutil.AssertValue(t, got, test.want, "writeNewline")
		})
	}
}

func Test_maxAttrLen(t *testing.T) {
	type args struct {
		a Attr
		b Attr
		c Attr
	}
	type want struct {
		val int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "greatest combined length",
			args: args{
				a: Attr{
					Class: "class",
					Style: "a",
				},
				b: Attr{
					Class: "b",
					Style: "long-style",
				},
				c: Attr{
					Class: "cc",
					Style: "cc",
				},
			},
			want: want{
				val: 11,
			},
		},
		{
			name: "empty",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := maxAttrLen(test.args.a, test.args.b, test.args.c)
			testutil.AssertValue(t, got, test.want.val, "maxAttrLen")
		})
	}
}
