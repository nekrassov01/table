package text

import (
	"bytes"
	"io"
	"testing"

	"github.com/nekrassov01/table/internal/testutil"
)

func Test_resolveTerminalWidth(t *testing.T) {
	type args struct {
		w       io.Writer
		getSize func(int) (int, int, error)
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
				w: func() io.Writer {
					file, _ := testutil.NewFile(t)
					return file
				}(),
			},
		},
		{
			name: "terminal width",
			args: args{
				w: func() io.Writer {
					file, _ := testutil.NewFile(t)
					return file
				}(),
				getSize: func(int) (int, int, error) {
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
				w: func() io.Writer {
					file, _ := testutil.NewFile(t)
					return file
				}(),
				getSize: func(int) (int, int, error) {
					return 0, 0, testutil.NewError()
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			original := getSize
			if test.args.getSize != nil {
				getSize = test.args.getSize
			}
			t.Cleanup(func() {
				getSize = original
			})
			got := resolveTerminalWidth(test.args.w)
			testutil.AssertValue(t, got, test.want.val, "resolveTerminalWidth")
		})
	}
}

func Test_resolveIsTerminal(t *testing.T) {
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
				w: func() io.Writer {
					file, _ := testutil.NewFile(t)
					return file
				}(),
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
