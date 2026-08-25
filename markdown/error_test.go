package markdown

import (
	"errors"
	"testing"

	"github.com/nekrassov01/table"
	"github.com/nekrassov01/table/internal/testutil"
)

func Test_newWriteError(t *testing.T) {
	type args struct {
		err error
	}
	type want struct {
		message string
		write   bool
		cause   bool
		pkg     string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "writer error",
			args: args{
				err: testutil.NewError(),
			},
			want: want{
				message: "markdown: write failed: test",
				write:   true,
				cause:   true,
				pkg:     "markdown",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := newWriteError(test.args.err)
			tableError, ok := err.(*table.Error)
			if !ok {
				t.Fatalf("newWriteError type = %T, want *table.Error", err)
			}
			got := want{
				message: err.Error(),
				write:   errors.Is(err, table.ErrWriteFailed),
				cause:   errors.Is(err, test.args.err),
				pkg:     tableError.Pkg,
			}
			testutil.AssertValue(t, got, test.want, "newWriteError")
		})
	}
}

func Test_newClosedError(t *testing.T) {
	type want struct {
		pkg      string
		isClosed bool
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "wraps closed error",
			want: want{
				pkg:      "markdown",
				isClosed: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := newClosedError()
			got := want{
				pkg:      err.(*table.Error).Pkg,
				isClosed: errors.Is(err, table.ErrClosed),
			}
			testutil.AssertValue(t, got, test.want, "newClosedError")
		})
	}
}

func Test_newColumnCountError(t *testing.T) {
	type args struct {
		got  int
		want int
	}
	type want struct {
		message string
		match   bool
		pkg     string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "wraps column count error",
			args: args{
				got:  3,
				want: 2,
			},
			want: want{
				message: "markdown: column count exceeded: got 3, want 2",
				match:   true,
				pkg:     "markdown",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := newColumnCountError(test.args.got, test.args.want)
			tableError, ok := err.(*table.Error)
			if !ok {
				t.Fatalf("newColumnCountError type = %T, want *table.Error", err)
			}
			got := want{
				message: err.Error(),
				match:   errors.Is(err, table.ErrColumnCount),
				pkg:     tableError.Pkg,
			}
			testutil.AssertValue(t, got, test.want, "newColumnCountError")
		})
	}
}

func Test_newHeaderError(t *testing.T) {
	type want struct {
		pkg      string
		isHeader bool
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "wraps header error",
			want: want{
				pkg:      "markdown",
				isHeader: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := newHeaderError()
			got := want{
				pkg:      err.(*table.Error).Pkg,
				isHeader: errors.Is(err, table.ErrHeaderRequired),
			}
			testutil.AssertValue(t, got, test.want, "newHeaderError")
		})
	}
}

func Test_newError(t *testing.T) {
	type args struct {
		err error
	}
	type want struct {
		pkg    string
		isSame bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "wraps source error",
			args: args{
				err: testutil.NewError(),
			},
			want: want{
				pkg:    "markdown",
				isSame: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := newError(test.args.err)
			got := want{
				pkg:    err.(*table.Error).Pkg,
				isSame: errors.Is(err, test.args.err),
			}
			testutil.AssertValue(t, got, test.want, "newError")
		})
	}
}
