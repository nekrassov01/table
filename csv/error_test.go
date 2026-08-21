package csv

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
				message: "csv: write failed: test",
				write:   true,
				cause:   true,
				pkg:     "csv",
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
		closed  bool
		pkg     string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "wraps closed error",
			want: want{
				message: "csv: render after close",
				closed:  true,
				pkg:     "csv",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := newClosedError()
			var wrapped *table.Error
			got := want{
				message: err.Error(),
				closed:  errors.Is(err, table.ErrClosed),
			}
			if errors.As(err, &wrapped) {
				got.pkg = wrapped.Pkg
			}
			testutil.AssertValue(t, got, test.want, "newClosedError")
		})
	}
}

func Test_newDelimiterError(t *testing.T) {
	type args struct {
		delimiter rune
	}
	type want struct {
		message   string
		delimiter bool
		pkg       string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "names invalid delimiter",
			args: args{
				delimiter: '"',
			},
			want: want{
				message:   "csv: invalid delimiter: '\"'",
				delimiter: true,
				pkg:       "csv",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := newDelimiterError(test.args.delimiter)
			tableError, ok := err.(*table.Error)
			if !ok {
				t.Fatalf("newDelimiterError type = %T, want *table.Error", err)
			}
			got := want{
				message:   err.Error(),
				delimiter: errors.Is(err, table.ErrDelimiter),
				pkg:       tableError.Pkg,
			}
			testutil.AssertValue(t, got, test.want, "newDelimiterError")
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
		columns bool
		pkg     string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "names received and expected counts",
			args: args{
				got:  3,
				want: 2,
			},
			want: want{
				message: "csv: column count exceeded: got 3, want 2",
				columns: true,
				pkg:     "csv",
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
				columns: errors.Is(err, table.ErrColumnCount),
				pkg:     tableError.Pkg,
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
		message string
		cause   bool
		pkg     string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "wraps cause with package",
			args: args{
				err: testutil.NewError(),
			},
			want: want{
				message: "csv: test",
				cause:   true,
				pkg:     "csv",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := newError(test.args.err)
			var wrapped *table.Error
			got := want{
				message: err.Error(),
				cause:   errors.Is(err, test.args.err),
			}
			if errors.As(err, &wrapped) {
				got.pkg = wrapped.Pkg
			}
			testutil.AssertValue(t, got, test.want, "newError")
		})
	}
}
