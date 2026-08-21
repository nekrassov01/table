package text

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
				message: "text: write failed: test",
				write:   true,
				cause:   true,
				pkg:     "text",
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
		message    string
		matches    bool
		found      bool
		tableError *table.Error
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "closed",
			want: want{
				message: "text: render after close",
				matches: true,
				found:   true,
				tableError: &table.Error{
					Pkg: "text",
					Err: table.ErrClosed,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotError := newClosedError()
			var gotTableError *table.Error
			gotFound := errors.As(gotError, &gotTableError)
			got := want{
				message:    gotError.Error(),
				matches:    errors.Is(gotError, table.ErrClosed),
				found:      gotFound,
				tableError: gotTableError,
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
		matches bool
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
				message: "text: column count exceeded: got 3, want 2",
				matches: true,
				pkg:     "text",
			},
		},
		{
			name: "zero counts",
			want: want{
				message: "text: column count exceeded: got 0, want 0",
				matches: true,
				pkg:     "text",
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
				matches: errors.Is(gotError, table.ErrColumnCount),
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
			name: "with error",
			args: args{
				err: table.ErrHeaderRequired,
			},
			want: want{
				val: &table.Error{
					Pkg: "text",
					Err: table.ErrHeaderRequired,
				},
			},
		},
		{
			name: "nil error",
			want: want{
				val: &table.Error{
					Pkg: "text",
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
