package csv

import (
	"bytes"
	"io"
	"testing"

	"github.com/nekrassov01/table/internal/testutil"
)

func Test_painter_prepare(t *testing.T) {
	type fields struct {
		input compilerResult
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
			name: "sizes for largest body row",
			fields: fields{
				input: compilerResult{
					configResult: configResult{
						option: &option{},
					},
					header: row{
						cells: []cell{{value: "A"}},
					},
					body: []row{
						{
							cells: []cell{{value: "a"}, {value: "long"}},
						},
					},
					footer: []row{
						{
							cells: []cell{{value: "f"}},
						},
					},
				},
			},
			want: want{
				lineCap: 7,
			},
		},
		{
			name: "keeps larger current line as backing",
			fields: fields{
				input: compilerResult{
					configResult: configResult{
						option: &option{},
					},
					header: row{
						cells: []cell{{value: "A"}},
					},
				},
				state: painterState{
					lineBacking: make([]byte, 0, 2),
					line:        make([]byte, 0, 8),
				},
			},
			want: want{
				lineCap: 8,
			},
		},
		{
			name: "accepts empty input",
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
				lineLen: len(state.line),
				lineCap: cap(state.line),
			}
			testutil.AssertValue(t, got, test.want, "prepare")
		})
	}
}

func Test_painter_lineSize(t *testing.T) {
	type fields struct {
		delimiter rune
		crlf      bool
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
			name: "empty row",
		},
		{
			name: "counts values delimiters and newline",
			fields: fields{
				delimiter: ',',
			},
			args: args{
				row: row{
					cells: []cell{{value: "a"}, {value: `"b,c"`}},
				},
			},
			want: want{
				size: 8,
			},
		},
		{
			name: "counts CRLF terminator",
			fields: fields{
				delimiter: ',',
				crlf:      true,
			},
			args: args{
				row: row{
					cells: []cell{{value: "a"}, {value: `"b,c"`}},
				},
			},
			want: want{
				size: 9,
			},
		},
		{
			name: "counts Unicode delimiter bytes",
			fields: fields{
				delimiter: '・',
			},
			args: args{
				row: row{
					cells: []cell{{value: "a"}, {value: `"b,c"`}},
				},
			},
			want: want{
				size: 10,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &painter{
				input: compilerResult{
					configResult: configResult{
						option: &option{
							delimiter: test.fields.delimiter,
							crlf:      test.fields.crlf,
						},
					},
				},
			}
			got := want{
				size: o.lineSize(&test.args.row),
			}
			testutil.AssertValue(t, got, test.want, "lineSize")
		})
	}
}

func Test_painter_paintHeader(t *testing.T) {
	type fields struct {
		input compilerResult
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
			name: "writes header",
			fields: fields{
				input: compilerResult{
					configResult: configResult{
						option: &option{
							delimiter: ',',
						},
					},
					header: row{
						cells: []cell{{value: "A"}, {value: "B"}},
					},
				},
			},
			want: want{
				output: "A,B\n",
			},
		},
		{
			name: "skips absent header",
			fields: fields{
				input: compilerResult{
					configResult: configResult{
						option: &option{},
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var buf bytes.Buffer
			state := painterState{}
			o := &painter{
				input: test.fields.input,
				state: &state,
				w:     &buf,
			}
			o.paintHeader()
			got := want{
				output: buf.String(),
			}
			testutil.AssertValue(t, got, test.want, "paintHeader")
		})
	}
}

func Test_painter_paintBody(t *testing.T) {
	type fields struct {
		input compilerResult
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
			name: "writes rows in order",
			fields: fields{
				input: compilerResult{
					configResult: configResult{
						option: &option{
							delimiter: '\t',
						},
					},
					body: []row{
						{cells: []cell{{value: "a"}, {value: "1"}}},
						{cells: []cell{{value: "b"}, {value: "2"}}},
					},
				},
			},
			want: want{
				output: "a\t1\nb\t2\n",
			},
		},
		{
			name: "stops after error",
			fields: fields{
				input: compilerResult{
					body: []row{{cells: []cell{{value: "a"}}}},
				},
				err: testutil.NewError(),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var buf bytes.Buffer
			state := painterState{}
			o := &painter{
				input: test.fields.input,
				state: &state,
				w:     &buf,
				err:   test.fields.err,
			}
			o.paintBody()
			got := want{
				output: buf.String(),
			}
			testutil.AssertValue(t, got, test.want, "paintBody")
		})
	}
}

func Test_painter_paintFooter(t *testing.T) {
	type fields struct {
		input compilerResult
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
			name: "writes rows in order",
			fields: fields{
				input: compilerResult{
					configResult: configResult{
						option: &option{
							delimiter: ',',
						},
					},
					footer: []row{
						{cells: []cell{{value: "sum"}, {value: "3"}}},
						{cells: []cell{{value: "end"}, {value: ""}}},
					},
				},
			},
			want: want{
				output: "sum,3\nend,\n",
			},
		},
		{
			name: "stops after error",
			fields: fields{
				input: compilerResult{
					footer: []row{{cells: []cell{{value: "sum"}}}},
				},
				err: testutil.NewError(),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var buf bytes.Buffer
			state := painterState{}
			o := &painter{
				input: test.fields.input,
				state: &state,
				w:     &buf,
				err:   test.fields.err,
			}
			o.paintFooter()
			got := want{
				output: buf.String(),
			}
			testutil.AssertValue(t, got, test.want, "paintFooter")
		})
	}
}

func Test_painter_paintRow(t *testing.T) {
	type fields struct {
		delimiter rune
		crlf      bool
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
			name: "writes delimited row",
			fields: fields{
				delimiter: ',',
			},
			args: args{
				row: row{
					cells: []cell{{value: "a"}, {value: ""}, {value: "c"}},
				},
			},
			want: want{
				output: "a,,c\n",
			},
		},
		{
			name: "skips empty row",
			fields: fields{
				delimiter: ',',
			},
		},
		{
			name: "writes CRLF terminator",
			fields: fields{
				delimiter: ',',
				crlf:      true,
			},
			args: args{
				row: row{
					cells: []cell{{value: "a"}, {value: "b"}},
				},
			},
			want: want{
				output: "a,b\r\n",
			},
		},
		{
			name: "writes Unicode delimiter",
			fields: fields{
				delimiter: '・',
			},
			args: args{
				row: row{
					cells: []cell{{value: "a"}, {value: "b"}},
				},
			},
			want: want{
				output: "a・b\n",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var buf bytes.Buffer
			state := painterState{}
			o := &painter{
				input: compilerResult{
					configResult: configResult{
						option: &option{
							delimiter: test.fields.delimiter,
							crlf:      test.fields.crlf,
						},
					},
				},
				state: &state,
				w:     &buf,
			}
			o.paintRow(test.args.row)
			got := want{
				output: buf.String(),
			}
			testutil.AssertValue(t, got, test.want, "paintRow")
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
			name: "resets length and preserves capacity",
			fields: fields{
				line: []byte("value"),
			},
			want: want{
				lineCap: 5,
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
		line []byte
		err  error
		fail bool
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
				line: []byte("value\n"),
			},
			want: want{
				output: "value\n",
			},
		},
		{
			name: "skips empty line",
		},
		{
			name: "skips after error",
			fields: fields{
				line: []byte("value\n"),
				err:  testutil.NewError(),
			},
			want: want{
				hasError: true,
			},
		},
		{
			name: "wraps writer error",
			fields: fields{
				line: []byte("value\n"),
				fail: true,
			},
			want: want{
				hasError: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var buf bytes.Buffer
			state := painterState{
				line: test.fields.line,
			}
			var writer io.Writer = &buf
			if test.fields.fail {
				writer = &testutil.ErrorWriter{
					Err: testutil.NewError(),
				}
			}
			o := &painter{
				state: &state,
				w:     writer,
				err:   test.fields.err,
			}
			o.flushLine()
			got := want{
				output:   buf.String(),
				hasError: o.err != nil,
			}
			testutil.AssertValue(t, got, test.want, "flushLine")
		})
	}
}

func Test_painter_writeNewline(t *testing.T) {
	type fields struct {
		line []byte
		crlf bool
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
				line: []byte("value"),
			},
			want: want{
				line: "value\n",
			},
		},
		{
			name: "appends CRLF",
			fields: fields{
				line: []byte("value"),
				crlf: true,
			},
			want: want{
				line: "value\r\n",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := painterState{
				line: test.fields.line,
			}
			o := &painter{
				input: compilerResult{
					configResult: configResult{
						option: &option{
							crlf: test.fields.crlf,
						},
					},
				},
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
