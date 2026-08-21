package backlog

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
				message: "backlog: write failed: test",
				write:   true,
				cause:   true,
				pkg:     "backlog",
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
		message string
		match   bool
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "closed",
			want: want{
				message: "backlog: render after close",
				match:   true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotError := newClosedError()
			got := want{
				message: gotError.Error(),
				match:   errors.Is(gotError, table.ErrClosed),
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
			name: "wider row",
			args: args{
				got:  3,
				want: 2,
			},
			want: want{
				message: "backlog: column count exceeded: got 3, want 2",
				match:   true,
				pkg:     "backlog",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotError := newColumnCountError(test.args.got, test.args.want)
			gotTableError, ok := gotError.(*table.Error)
			if !ok {
				t.Fatalf("newColumnCountError type = %T, want *table.Error", gotError)
			}
			got := want{
				message: gotError.Error(),
				match:   errors.Is(gotError, table.ErrColumnCount),
				pkg:     gotTableError.Pkg,
			}
			testutil.AssertValue(t, got, test.want, "newColumnCountError")
		})
	}
}

func Test_newError(t *testing.T) {
	type args struct {
		err error
	}
	type want struct {
		val error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "wraps error",
			args: args{
				err: table.ErrClosed,
			},
			want: want{
				val: &table.Error{
					Pkg: "backlog",
					Err: table.ErrClosed,
				},
			},
		},
		{
			name: "nil error",
			want: want{
				val: &table.Error{
					Pkg: "backlog",
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := newError(test.args.err)
			testutil.AssertValue(t, got, test.want.val, "newError")
		})
	}
}
