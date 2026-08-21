package markdown

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/nekrassov01/table/internal/testutil"
)

func Test_painter_prepare(t *testing.T) {
	type fields struct {
		input solverResult
		state painterState
	}
	type want struct {
		lineLen    int
		lineCap    int
		backingCap int
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "adopts growth and uses current row size",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						header: row{
							cells: []cell{
								{
									width: 4,
									size:  7,
								},
							},
						},
					},
					metrics: []columnMetric{
						{
							box: box{
								width: 3,
							},
							separator: separator{
								width: 5,
							},
						},
					},
				},
				state: painterState{
					lineBacking: make([]byte, 0, 1),
					line:        make([]byte, 0, 8),
				},
			},
			want: want{
				lineCap:    12,
				backingCap: 12,
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
				lineLen:    len(state.line),
				lineCap:    cap(state.line),
				backingCap: cap(state.lineBacking),
			}
			testutil.AssertValue(t, got, test.want, "prepare")
		})
	}
}

func Test_painter_lineSize(t *testing.T) {
	type fields struct {
		input solverResult
	}
	type args struct {
		row row
	}
	type want struct {
		size int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "adds solved padding",
			fields: fields{
				input: solverResult{
					metrics: []columnMetric{
						{
							box: box{
								width: 5,
							},
						},
					},
				},
			},
			args: args{
				row: row{
					cells: []cell{
						{
							width: 2,
							size:  4,
						},
					},
				},
			},
			want: want{
				size: 12,
			},
		},
		{
			name: "uses unpadded stream cell",
			fields: fields{
				input: solverResult{
					metrics: make([]columnMetric, 1),
				},
			},
			args: args{
				row: row{
					cells: []cell{
						{
							width: 4,
							size:  7,
						},
					},
				},
			},
			want: want{
				size: 12,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &painter{
				input: test.fields.input,
			}
			r := test.args.row
			got := want{
				size: o.lineSize(&r),
			}
			testutil.AssertValue(t, got, test.want, "lineSize")
		})
	}
}

func Test_painter_paintHeader(t *testing.T) {
	type fields struct {
		input solverResult
		state painterState
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
			name: "writes header and delimiter",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						header: row{
							cells: []cell{
								{
									value: "A",
									width: 1,
								},
							},
						},
					},
					metrics: []columnMetric{
						{
							separator: separator{
								width: 3,
							},
						},
					},
				},
				state: painterState{
					line: make([]byte, 0, 16),
				},
			},
			want: want{
				output: "| A |\n| --- |\n",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			var output bytes.Buffer
			o := &painter{
				input: test.fields.input,
				state: &state,
				w:     &output,
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
		state painterState
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
			name: "writes body rows",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						body: []row{
							{
								cells: []cell{
									{
										value: "x",
										width: 1,
									},
								},
							},
							{
								cells: []cell{
									{
										value: "y",
										width: 1,
									},
								},
							},
						},
					},
					metrics: []columnMetric{{}},
				},
				state: painterState{
					line: make([]byte, 0, 16),
				},
			},
			want: want{
				output: "| x |\n| y |\n",
			},
		},
		{
			name: "stops on sticky error",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						body: []row{
							{
								cells: []cell{{}},
							},
						},
					},
					metrics: []columnMetric{{}},
				},
				err: testutil.NewError(),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			var output bytes.Buffer
			o := &painter{
				input: test.fields.input,
				state: &state,
				w:     &output,
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

func Test_painter_paintSeparator(t *testing.T) {
	type fields struct {
		input solverResult
		state painterState
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
			name: "writes aligned delimiter row",
			fields: fields{
				input: solverResult{
					metrics: []columnMetric{
						{
							box: box{
								width: 5,
								align: AlignCenter,
							},
							separator: separator{
								width: 5,
								lead:  true,
								trail: true,
							},
						},
					},
				},
				state: painterState{
					line: make([]byte, 0, 16),
				},
			},
			want: want{
				output: "| :---: |\n",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			var output bytes.Buffer
			o := &painter{
				input: test.fields.input,
				state: &state,
				w:     &output,
			}
			o.paintSeparator()
			got := want{
				output: output.String(),
			}
			testutil.AssertValue(t, got, test.want, "paintSeparator")
		})
	}
}

func Test_painter_paintSeparatorCell(t *testing.T) {
	type fields struct {
		state painterState
	}
	type args struct {
		metric columnMetric
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
			name: "writes trailing marker",
			fields: fields{
				state: painterState{
					line: []byte("prefix:"),
				},
			},
			args: args{
				metric: columnMetric{
					box: box{
						align: AlignRight,
					},
					separator: separator{
						width: 4,
						trail: true,
					},
				},
			},
			want: want{
				line: "prefix:---:",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			o := &painter{
				state: &state,
			}
			metric := test.args.metric
			o.paintSeparatorCell(&metric)
			got := want{
				line: string(state.line),
			}
			testutil.AssertValue(t, got, test.want, "paintSeparatorCell")
		})
	}
}

func Test_painter_paintRow(t *testing.T) {
	type fields struct {
		input solverResult
		state painterState
	}
	type args struct {
		row row
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
			name: "writes padded row",
			fields: fields{
				input: solverResult{
					metrics: []columnMetric{
						{
							box: box{
								width: 3,
								align: AlignRight,
							},
						},
					},
				},
				state: painterState{
					line: make([]byte, 0, 16),
				},
			},
			args: args{
				row: row{
					cells: []cell{
						{
							value: "x",
							width: 1,
						},
					},
				},
			},
			want: want{
				output: "|   x |\n",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			var output bytes.Buffer
			o := &painter{
				input: test.fields.input,
				state: &state,
				w:     &output,
			}
			o.paintRow(test.args.row)
			got := want{
				output: output.String(),
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
		cell cell
		box  box
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
			name: "left",
			args: args{
				cell: cell{
					value: "x",
					width: 1,
				},
				box: box{
					width: 4,
				},
			},
			want: want{
				line: "x   ",
			},
		},
		{
			name: "right",
			args: args{
				cell: cell{
					value: "x",
					width: 1,
				},
				box: box{
					width: 4,
					align: AlignRight,
				},
			},
			want: want{
				line: "   x",
			},
		},
		{
			name: "center",
			args: args{
				cell: cell{
					value: "x",
					width: 1,
				},
				box: box{
					width: 4,
					align: AlignCenter,
				},
			},
			want: want{
				line: " x  ",
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
			box := test.args.box
			o.paintCell(&cell, &box)
			got := want{
				line: string(state.line),
			}
			testutil.AssertValue(t, got, test.want, "paintCell")
		})
	}
}

func Test_painter_resetLine(t *testing.T) {
	type fields struct {
		state painterState
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
			name: "retains capacity",
			fields: fields{
				state: painterState{
					line: func() []byte {
						line := make([]byte, 8)
						return line[:5]
					}(),
				},
			},
			want: want{
				lineCap: 8,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
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
		state  painterState
		writer func(*bytes.Buffer) io.Writer
	}
	type want struct {
		output   string
		hasError bool
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "writes line",
			fields: fields{
				state: painterState{
					line: []byte("line"),
				},
				writer: func(output *bytes.Buffer) io.Writer {
					return output
				},
			},
			want: want{
				output: "line",
			},
		},
		{
			name: "wraps writer error",
			fields: fields{
				state: painterState{
					line: []byte("line"),
				},
				writer: func(*bytes.Buffer) io.Writer {
					return &testutil.ErrorWriter{
						Err: testutil.NewError(),
					}
				},
			},
			want: want{
				hasError: true,
			},
		},
		{
			name: "skips empty line",
			fields: fields{
				writer: func(output *bytes.Buffer) io.Writer {
					return output
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			var output bytes.Buffer
			writer := test.fields.writer(&output)
			o := &painter{
				state: &state,
				w:     writer,
			}
			o.flushLine()
			got := want{
				output:   output.String(),
				hasError: o.err != nil,
			}
			testutil.AssertValue(t, got, test.want, "flushLine")
		})
	}
}

func Test_painter_writeValue(t *testing.T) {
	type fields struct {
		state painterState
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
			name: "empty",
			fields: fields{
				state: painterState{
					line: []byte("prefix:"),
				},
			},
			want: want{
				line: "prefix:",
			},
		},
		{
			name: "decoration contains color",
			args: args{
				cell: cell{
					value:      "x",
					color:      ColorFgRed,
					decoration: DecorationBold,
				},
			},
			want: want{
				line: `**<span style="color:red">x</span>**`,
			},
		},
		{
			name: "block decoration contains color",
			args: args{
				cell: cell{
					value:      "x",
					color:      ColorFgRed,
					decoration: DecorationPreformatted,
				},
			},
			want: want{
				line: `<pre><span style="color:red">x</span></pre>`,
			},
		},
		{
			name: "color contains code fence",
			args: args{
				cell: cell{
					value: "x",
					color: ColorFgRed,
					ticks: 2,
				},
			},
			want: want{
				line: "<span style=\"color:red\">``x``</span>",
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
			o.writeValue(&cell)
			got := want{
				line: string(state.line),
			}
			testutil.AssertValue(t, got, test.want, "writeValue")
		})
	}
}

func Test_painter_writeBackticks(t *testing.T) {
	type fields struct {
		state painterState
	}
	type args struct {
		n int
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
			name: "appends requested count",
			fields: fields{
				state: painterState{
					line: []byte("prefix:"),
				},
			},
			args: args{
				n: 3,
			},
			want: want{
				line: "prefix:```",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			o := &painter{
				state: &state,
			}
			o.writeBackticks(test.args.n)
			got := want{
				line: string(state.line),
			}
			testutil.AssertValue(t, got, test.want, "writeBackticks")
		})
	}
}

func Test_painter_writeDashes(t *testing.T) {
	type fields struct {
		state painterState
	}
	type args struct {
		n int
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
			name: "non-positive",
			fields: fields{
				state: painterState{
					line: []byte("prefix:"),
				},
			},
			want: want{
				line: "prefix:",
			},
		},
		{
			name: "larger than reusable block",
			args: args{
				n: 257,
			},
			want: want{
				line: strings.Repeat("-", 257),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			o := &painter{
				state: &state,
			}
			o.writeDashes(test.args.n)
			got := want{
				line: string(state.line),
			}
			testutil.AssertValue(t, got, test.want, "writeDashes")
		})
	}
}

func Test_painter_writeSpaces(t *testing.T) {
	type fields struct {
		state painterState
	}
	type args struct {
		n int
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
			name: "non-positive",
			fields: fields{
				state: painterState{
					line: []byte("prefix:"),
				},
			},
			want: want{
				line: "prefix:",
			},
		},
		{
			name: "larger than reusable block",
			args: args{
				n: 257,
			},
			want: want{
				line: strings.Repeat(" ", 257),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			o := &painter{
				state: &state,
			}
			o.writeSpaces(test.args.n)
			got := want{
				line: string(state.line),
			}
			testutil.AssertValue(t, got, test.want, "writeSpaces")
		})
	}
}

func Test_painter_writeNewline(t *testing.T) {
	type fields struct {
		state painterState
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
					line: []byte("line"),
				},
			},
			want: want{
				line: "line\n",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
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
