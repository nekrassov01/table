package text

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/nekrassov01/table/internal/testutil"
)

func Test_resolveWriter(t *testing.T) {
	file, _ := testutil.NewFile(t)
	type args struct {
		w           io.Writer
		passthrough bool
	}
	type want struct {
		same    bool
		called  bool
		written string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "non-file writer",
			args: args{
				w: &bytes.Buffer{},
			},
			want: want{
				same: true,
			},
		},
		{
			name: "terminal file",
			args: args{
				w: file,
			},
			want: want{
				called:  true,
				written: "value",
			},
		},
		{
			name: "terminal file with native ANSI support",
			args: args{
				w:           file,
				passthrough: true,
			},
			want: want{
				same:   true,
				called: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := test.args.w
			var adapted bytes.Buffer
			called := false
			original := terminalWriter
			terminalWriter = func(file *os.File) io.Writer {
				called = true
				if test.args.passthrough {
					return file
				}
				return &adapted
			}
			t.Cleanup(func() {
				terminalWriter = original
			})
			gotWriter := resolveWriter(input)
			if !test.want.same {
				_, _ = gotWriter.Write([]byte("value"))
			}
			got := want{
				same:    gotWriter == input,
				called:  called,
				written: adapted.String(),
			}
			testutil.AssertValue(t, got, test.want, "resolveWriter")
		})
	}
}

func Test_resolveTerminalWidth(t *testing.T) {
	file, _ := testutil.NewFile(t)
	type args struct {
		w        io.Writer
		termSize func(int) (int, int, error)
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
			name: "non-file writer",
			args: args{
				w: &bytes.Buffer{},
			},
		},
		{
			name: "regular file",
			args: args{
				w: file,
			},
		},
		{
			name: "terminal width",
			args: args{
				w: file,
				termSize: func(int) (int, int, error) {
					return 80, 24, nil
				},
			},
			want: want{
				val: 80,
			},
		},
		{
			name: "adapted terminal width",
			args: args{
				w: &writer{
					Writer: &bytes.Buffer{},
					file:   file,
				},
				termSize: func(int) (int, int, error) {
					return 80, 24, nil
				},
			},
			want: want{
				val: 80,
			},
		},
		{
			name: "terminal size error",
			args: args{
				w: file,
				termSize: func(int) (int, int, error) {
					return 0, 0, testutil.NewError()
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			original := terminalSize
			if test.args.termSize != nil {
				terminalSize = test.args.termSize
			}
			t.Cleanup(func() {
				terminalSize = original
			})
			got := resolveTerminalWidth(test.args.w)
			testutil.AssertValue(t, got, test.want.val, "resolveTerminalWidth")
		})
	}
}

func Test_resolveIsTerminal(t *testing.T) {
	file, _ := testutil.NewFile(t)
	type args struct {
		w io.Writer
	}
	type want struct {
		val bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "non-file writer",
			args: args{
				w: &bytes.Buffer{},
			},
		},
		{
			name: "regular file",
			args: args{
				w: file,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := resolveIsTerminal(test.args.w)
			testutil.AssertValue(t, got, test.want.val, "resolveIsTerminal")
		})
	}
}
