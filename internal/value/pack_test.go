package value

import (
	"testing"

	"github.com/nekrassov01/table/internal/testutil"
)

func Test_unpackString(t *testing.T) {
	type args struct {
		value string
	}
	type want struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "text",
			args: args{
				value: "abc",
			},
			want: want{
				value: "abc",
			},
		},
		{
			name: "empty",
			args: args{
				value: "",
			},
			want: want{
				value: "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := unpackString(packString(test.args.value))
			testutil.AssertValue(t, got, test.want.value, "string")
		})
	}
}

func Test_unpackBytes(t *testing.T) {
	type args struct {
		value []byte
	}
	type want struct {
		value    []byte
		capacity int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "text",
			args: args{
				value: []byte("abc"),
			},
			want: want{
				value:    []byte("abc"),
				capacity: 3,
			},
		},
		{
			name: "nil",
			args: args{
				value: nil,
			},
			want: want{
				value:    nil,
				capacity: 0,
			},
		},
		{
			name: "empty",
			args: args{
				value: []byte{},
			},
			want: want{
				value:    []byte{},
				capacity: 0,
			},
		},
		{
			name: "empty with capacity",
			args: args{
				value: make([]byte, 0, 8),
			},
			want: want{
				value:    []byte{},
				capacity: 0,
			},
		},
		{
			name: "spare capacity",
			args: args{
				value: make([]byte, 2, 8),
			},
			want: want{
				value:    []byte{0, 0},
				capacity: 2,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := unpackBytes(packBytes(test.args.value))
			testutil.AssertValue(t, got, test.want.value, "bytes")
			testutil.AssertValue(t, cap(got), test.want.capacity, "capacity")
		})
	}
}
