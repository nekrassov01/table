package table

import (
	"errors"
	"testing"

	"github.com/nekrassov01/table/internal/testutil"
)

func TestError_Error(t *testing.T) {
	type fields struct {
		Pkg string
		Err error
	}
	type want struct {
		val string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "with message",
			fields: fields{
				Pkg: "markdown",
				Err: errors.New("header is required"),
			},
			want: want{
				val: "markdown: header is required",
			},
		},
		{
			name: "with wrapped sentinel",
			fields: fields{
				Pkg: "text",
				Err: ErrClosed,
			},
			want: want{
				val: "text: render after close",
			},
		},
		{
			name: "with write sentinel",
			fields: fields{
				Pkg: "html",
				Err: ErrWriteFailed,
			},
			want: want{
				val: "html: write failed",
			},
		},
		{
			name: "nil inner error",
			fields: fields{
				Pkg: "backlog",
				Err: nil,
			},
			want: want{
				val: "backlog",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := &Error{
				Pkg: test.fields.Pkg,
				Err: test.fields.Err,
			}
			got := e.Error()
			testutil.AssertValue(t, got, test.want.val, "Error")
		})
	}
}

func TestError_Unwrap(t *testing.T) {
	type fields struct {
		Pkg string
		Err error
	}
	type want struct {
		val error
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "with inner error",
			fields: fields{
				Pkg: "text",
				Err: testutil.NewError(),
			},
			want: want{
				val: testutil.NewError(),
			},
		},
		{
			name: "nil inner error",
			fields: fields{
				Pkg: "text",
				Err: nil,
			},
			want: want{
				val: nil,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := &Error{
				Pkg: test.fields.Pkg,
				Err: test.fields.Err,
			}
			got := e.Unwrap()
			testutil.AssertValue(t, got, test.want.val, "Unwrap")
		})
	}
}
