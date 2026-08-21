package value

import (
	"testing"

	"github.com/nekrassov01/table/internal/testutil"
)

func TestStore_Mark(t *testing.T) {
	type fields struct {
		buf []byte
	}
	type want struct {
		val int
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "empty",
			want: want{
				val: 0,
			},
		},
		{
			name: "appended bytes",
			fields: fields{
				buf: []byte("abc"),
			},
			want: want{
				val: 3,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &Store{
				buf: test.fields.buf,
			}
			got := o.Mark()
			testutil.AssertValue(t, got, test.want.val, "Mark")
		})
	}
}

func TestStore_Reset(t *testing.T) {
	type fields struct {
		buf []byte
	}
	type want struct {
		length   int
		capacity int
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "empties and retains backing",
			fields: fields{
				buf: make([]byte, 3, 8),
			},
			want: want{
				length:   0,
				capacity: 8,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &Store{
				buf: test.fields.buf,
			}
			o.Reset()
			got := want{
				length:   len(o.buf),
				capacity: cap(o.buf),
			}
			testutil.AssertValue(t, got, test.want, "Reset")
		})
	}
}

func TestStore_Since(t *testing.T) {
	type fields struct {
		buf []byte
	}
	type args struct {
		mark int
	}
	type want struct {
		val string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "suffix from mark",
			fields: fields{
				buf: []byte("abcdef"),
			},
			args: args{
				mark: 3,
			},
			want: want{
				val: "def",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &Store{
				buf: test.fields.buf,
			}
			got := o.Since(test.args.mark)
			testutil.AssertValue(t, got, test.want.val, "Since")
		})
	}
}

func TestStore_AppendString(t *testing.T) {
	type fields struct {
		buf []byte
	}
	type args struct {
		s string
	}
	type want struct {
		val string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "empty",
			args: args{
				s: "",
			},
			want: want{
				val: "",
			},
		},
		{
			name: "single",
			args: args{
				s: "abc",
			},
			want: want{
				val: "abc",
			},
		},
		{
			name: "appends",
			fields: fields{
				buf: []byte("abc"),
			},
			args: args{
				s: "def",
			},
			want: want{
				val: "abcdef",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &Store{
				buf: test.fields.buf,
			}
			o.AppendString(test.args.s)
			testutil.AssertValue(t, string(o.buf), test.want.val, "AppendString")
		})
	}
}

func TestStore_AppendBytes(t *testing.T) {
	type fields struct {
		buf []byte
	}
	type args struct {
		b []byte
	}
	type want struct {
		val string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "bytes",
			fields: fields{
				buf: []byte("a"),
			},
			args: args{
				b: []byte("bc"),
			},
			want: want{
				val: "abc",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &Store{
				buf: test.fields.buf,
			}
			o.AppendBytes(test.args.b)
			testutil.AssertValue(t, string(o.buf), test.want.val, "AppendBytes")
		})
	}
}

func TestStore_AppendInt(t *testing.T) {
	type fields struct {
		buf []byte
	}
	type args struct {
		x int64
	}
	type want struct {
		val string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "negative",
			args: args{
				x: -42,
			},
			want: want{
				val: "-42",
			},
		},
		{
			name: "zero",
			args: args{
				x: 0,
			},
			want: want{
				val: "0",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &Store{
				buf: test.fields.buf,
			}
			o.AppendInt(test.args.x)
			testutil.AssertValue(t, string(o.buf), test.want.val, "AppendInt")
		})
	}
}

func TestStore_AppendUint(t *testing.T) {
	type fields struct {
		buf []byte
	}
	type args struct {
		x uint64
	}
	type want struct {
		val string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "maximum",
			args: args{
				x: 18446744073709551615,
			},
			want: want{
				val: "18446744073709551615",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &Store{
				buf: test.fields.buf,
			}
			o.AppendUint(test.args.x)
			testutil.AssertValue(t, string(o.buf), test.want.val, "AppendUint")
		})
	}
}

func TestStore_AppendFloat(t *testing.T) {
	type fields struct {
		buf []byte
	}
	type args struct {
		x       float64
		bitSize int
	}
	type want struct {
		val string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "shortest",
			args: args{
				x:       1.25,
				bitSize: 64,
			},
			want: want{
				val: "1.25",
			},
		},
		{
			name: "32-bit rounding",
			args: args{
				x:       float64(float32(0.1)),
				bitSize: 32,
			},
			want: want{
				val: "0.1",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &Store{
				buf: test.fields.buf,
			}
			o.AppendFloat(test.args.x, test.args.bitSize)
			testutil.AssertValue(t, string(o.buf), test.want.val, "AppendFloat")
		})
	}
}

func TestStore_grow(t *testing.T) {
	type fields struct {
		buf []byte
	}
	type want struct {
		length   int
		capacity int
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "allocates initial backing",
			want: want{
				length:   0,
				capacity: 128,
			},
		},
		{
			name: "retains existing backing",
			fields: fields{
				buf: make([]byte, 3, 8),
			},
			want: want{
				length:   3,
				capacity: 8,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &Store{
				buf: test.fields.buf,
			}
			o.grow()
			got := want{
				length:   len(o.buf),
				capacity: cap(o.buf),
			}
			testutil.AssertValue(t, got, test.want, "grow")
		})
	}
}
