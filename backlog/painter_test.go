package backlog

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
		lineCap int
		lineLen int
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "sizes maximum row",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						header: []row{{cells: []cell{{value: "h", width: 1, size: 1}}}},
						body:   []row{{cells: []cell{{value: "body", width: 4, size: 4}}}},
						footer: []row{{cells: []cell{{value: "f", width: 1, size: 1}}}},
					},
					metrics: []columnMetric{{box: box{width: 4}}},
				},
			},
			want: want{
				lineCap: 9,
			},
		},
		{
			name: "adopts grown line",
			fields: fields{
				state: painterState{
					lineBacking: make([]byte, 0, 1),
					line:        make([]byte, 0, 8),
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
				input: test.fields.input,
				state: &state,
			}
			o.prepare()
			got := want{
				lineCap: cap(state.line),
				lineLen: len(state.line),
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
		row  row
		band bool
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
			name: "body padding and byte size",
			fields: fields{
				input: solverResult{
					metrics: []columnMetric{{box: box{width: 5}}},
				},
			},
			args: args{
				row: row{
					cells: []cell{{value: "界", width: 2, size: 3}},
				},
			},
			want: want{
				size: 11,
			},
		},
		{
			name: "band marker",
			fields: fields{
				input: solverResult{
					metrics: []columnMetric{{box: box{width: 3}}},
				},
			},
			args: args{
				row: row{
					cells: []cell{{}},
				},
				band: true,
			},
			want: want{
				size: 8,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &painter{
				input: test.fields.input,
			}
			got := want{
				size: o.lineSize(&test.args.row, test.args.band),
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
			name: "header rows",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						header: []row{{cells: []cell{{value: "H", width: 1, size: 1}}}},
					},
					metrics: []columnMetric{{}},
				},
			},
			want: want{
				output: "|~H  |\n",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var output bytes.Buffer
			state := test.fields.state
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
			name: "body rows",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						body: []row{
							{cells: []cell{{value: "a", width: 1, size: 1}}},
							{cells: []cell{{value: "b", width: 1, size: 1}}},
						},
					},
					metrics: []columnMetric{{}},
				},
			},
			want: want{
				output: "| a |\n| b |\n",
			},
		},
		{
			name: "sticky error",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						body: []row{{cells: []cell{{value: "a", width: 1, size: 1}}}},
					},
					metrics: []columnMetric{{}},
				},
				err: testutil.NewError(),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var output bytes.Buffer
			state := test.fields.state
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

func Test_painter_paintFooter(t *testing.T) {
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
			name: "footer rows",
			fields: fields{
				input: solverResult{
					compilerResult: compilerResult{
						footer: []row{{cells: []cell{{value: "F", width: 1, size: 1}}}},
					},
					metrics: []columnMetric{{}},
				},
			},
			want: want{
				output: "|~F  |\n",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var output bytes.Buffer
			state := test.fields.state
			o := &painter{
				input: test.fields.input,
				state: &state,
				w:     &output,
			}
			o.paintFooter()
			got := want{
				output: output.String(),
			}
			testutil.AssertValue(t, got, test.want, "paintFooter")
		})
	}
}

func Test_painter_paintBand(t *testing.T) {
	type fields struct {
		input solverResult
		state painterState
		err   error
	}
	type args struct {
		rows []row
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
			name: "rows",
			fields: fields{
				input: solverResult{
					metrics: []columnMetric{{}},
				},
			},
			args: args{
				rows: []row{
					{cells: []cell{{value: "a", width: 1, size: 1}}},
					{cells: []cell{{value: "b", width: 1, size: 1}}},
				},
			},
			want: want{
				output: "|~a  |\n|~b  |\n",
			},
		},
		{
			name: "sticky error",
			fields: fields{
				input: solverResult{
					metrics: []columnMetric{{}},
				},
				err: testutil.NewError(),
			},
			args: args{
				rows: []row{{cells: []cell{{value: "a", width: 1, size: 1}}}},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var output bytes.Buffer
			state := test.fields.state
			o := &painter{
				input: test.fields.input,
				state: &state,
				w:     &output,
				err:   test.fields.err,
			}
			o.paintBand(test.args.rows)
			got := want{
				output: output.String(),
			}
			testutil.AssertValue(t, got, test.want, "paintBand")
		})
	}
}

func Test_painter_paintRow(t *testing.T) {
	type fields struct {
		input solverResult
		state painterState
	}
	type args struct {
		row  row
		band bool
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
			name: "body row with padding",
			fields: fields{
				input: solverResult{
					metrics: []columnMetric{{box: box{width: 3}}},
				},
			},
			args: args{
				row: row{
					cells: []cell{{value: "x", width: 1, size: 1}},
				},
			},
			want: want{
				output: "| x   |\n",
			},
		},
		{
			name: "band row",
			fields: fields{
				input: solverResult{
					metrics: []columnMetric{{box: box{width: 2}}},
				},
			},
			args: args{
				row: row{
					cells: []cell{{value: "x", width: 1, size: 1}},
				},
				band: true,
			},
			want: want{
				output: "|~x  |\n",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var output bytes.Buffer
			state := test.fields.state
			o := &painter{
				input: test.fields.input,
				state: &state,
				w:     &output,
			}
			o.paintRow(test.args.row, test.args.band)
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
		band bool
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
			name: "body value and padding",
			args: args{
				cell: cell{value: "x", width: 1},
				box:  box{width: 3},
			},
			want: want{
				line: "x  ",
			},
		},
		{
			name: "empty band cell",
			args: args{
				box:  box{width: 2},
				band: true,
			},
			want: want{
				line: "~ ",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			o := &painter{
				state: &state,
			}
			o.paintCell(&test.args.cell, &test.args.box, test.args.band)
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
			name: "retains backing",
			fields: fields{
				state: painterState{
					line: []byte("value"),
				},
			},
			want: want{
				lineCap: 5,
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
		state       painterState
		err         error
		writerError error
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
			},
			want: want{
				output: "line",
			},
		},
		{
			name: "empty line",
		},
		{
			name: "sticky error",
			fields: fields{
				state: painterState{
					line: []byte("line"),
				},
				err: testutil.NewError(),
			},
			want: want{
				hasError: true,
			},
		},
		{
			name: "writer error",
			fields: fields{
				state: painterState{
					line: []byte("line"),
				},
				writerError: testutil.NewError(),
			},
			want: want{
				hasError: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var output bytes.Buffer
			var writer io.Writer = &output
			if test.fields.writerError != nil {
				writer = &testutil.ErrorWriter{Err: test.fields.writerError}
			}
			state := test.fields.state
			o := &painter{
				state: &state,
				w:     writer,
				err:   test.fields.err,
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
		},
		{
			name: "value",
			args: args{
				cell: cell{value: "value"},
			},
			want: want{
				line: "value",
			},
		},
		{
			name: "color outside decoration",
			args: args{
				cell: cell{
					value:      "value",
					color:      NewColor("red", ""),
					decoration: NewDecoration("<", ">"),
				},
			},
			want: want{
				line: "&color(red){<value>}",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			o := &painter{
				state: &state,
			}
			o.writeValue(&test.args.cell)
			got := want{
				line: string(state.line),
			}
			testutil.AssertValue(t, got, test.want, "writeValue")
		})
	}
}

func Test_painter_writeSpaces(t *testing.T) {
	type fields struct {
		state painterState
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
			name: "non-positive",
			args: args{
				count: -1,
			},
		},
		{
			name: "single block",
			fields: fields{
				state: painterState{
					line: []byte("x"),
				},
			},
			args: args{
				count: 2,
			},
			want: want{
				line: "x  ",
			},
		},
		{
			name: "multiple blocks",
			args: args{
				count: 257,
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
