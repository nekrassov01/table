package text

import (
	"testing"

	"github.com/nekrassov01/table/internal/testutil"
)

func Test_measureLine(t *testing.T) {
	type args struct {
		s string
	}
	type want struct {
		width    int
		hasBreak bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "empty",
			args: args{
				s: "",
			},
			want: want{
				width:    0,
				hasBreak: false,
			},
		},
		{
			name: "ascii",
			args: args{
				s: "hello",
			},
			want: want{
				width:    5,
				hasBreak: false,
			},
		},
		{
			name: "newline",
			args: args{
				s: "hello\nworld",
			},
			want: want{
				width:    0,
				hasBreak: true,
			},
		},
		{
			name: "carriage-return",
			args: args{
				s: "hello\rworld",
			},
			want: want{
				width:    0,
				hasBreak: true,
			},
		},
		{
			name: "crlf",
			args: args{
				s: "hello\r\nworld",
			},
			want: want{
				width:    0,
				hasBreak: true,
			},
		},
		{
			name: "newline-at-start",
			args: args{
				s: "\nhello",
			},
			want: want{
				width:    0,
				hasBreak: true,
			},
		},
		{
			name: "newline-at-end",
			args: args{
				s: "hello\n",
			},
			want: want{
				width:    0,
				hasBreak: true,
			},
		},
		{
			name: "control-in-ascii",
			args: args{
				s: "a\tb",
			},
			want: want{
				width:    2,
				hasBreak: false,
			},
		},
		{
			name: "control-then-break",
			args: args{
				s: "a\tb\nc",
			},
			want: want{
				width:    0,
				hasBreak: true,
			},
		},
		{
			name: "fullwidth",
			args: args{
				s: "あいう",
			},
			want: want{
				width:    6,
				hasBreak: false,
			},
		},
		{
			name: "mixed-ascii-fullwidth",
			args: args{
				s: "aあb",
			},
			want: want{
				width:    4,
				hasBreak: false,
			},
		},
		{
			name: "fullwidth-with-newline",
			args: args{
				s: "あ\nい",
			},
			want: want{
				width:    0,
				hasBreak: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotWidth, gotBreak := measureLine(test.args.s)
			testutil.AssertValue(t, gotWidth, test.want.width, "measureLine width")
			testutil.AssertValue(t, gotBreak, test.want.hasBreak, "measureLine hasBreak")
		})
	}
}

func Test_scanLine(t *testing.T) {
	type args struct {
		value string
		start int
	}
	type want struct {
		line     string
		width    int
		next     int
		hasBreak bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "ascii",
			args: args{
				value: "hello\nworld",
			},
			want: want{
				line:     "hello",
				width:    5,
				next:     6,
				hasBreak: true,
			},
		},
		{
			name: "fullwidth",
			args: args{
				value: "あいう\rnext",
			},
			want: want{
				line:     "あいう",
				width:    6,
				next:     10,
				hasBreak: true,
			},
		},
		{
			name: "crlf",
			args: args{
				value: "hello\r\nworld",
			},
			want: want{
				line:     "hello",
				width:    5,
				next:     7,
				hasBreak: true,
			},
		},
		{
			name: "offset",
			args: args{
				value: "hello\nworld",
				start: 6,
			},
			want: want{
				line:  "world",
				width: 5,
				next:  11,
			},
		},
		{
			name: "empty-after-break",
			args: args{
				value: "hello\n",
				start: 6,
			},
			want: want{
				next: 6,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotLine, gotWidth, gotNext, gotBreak := scanLine(test.args.value, test.args.start)
			testutil.AssertValue(t, gotLine, test.want.line, "scanLine line")
			testutil.AssertValue(t, gotWidth, test.want.width, "scanLine width")
			testutil.AssertValue(t, gotNext, test.want.next, "scanLine next")
			testutil.AssertValue(t, gotBreak, test.want.hasBreak, "scanLine hasBreak")
		})
	}
}
