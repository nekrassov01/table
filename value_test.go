package table

import (
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/nekrassov01/table/internal/testutil"
)

func TestString(t *testing.T) {
	type args struct {
		v string
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
			name: "empty",
			args: args{
				v: "",
			},
			want: want{
				value: "",
			},
		},
		{
			name: "text",
			args: args{
				v: "日本語\x00text",
			},
			want: want{
				value: "日本語\x00text",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := String(test.args.v)
			testutil.AssertValue(t, got.AsAny(), test.want.value, "String")
		})
	}
}

func TestBytes(t *testing.T) {
	type args struct {
		v []byte
	}
	type want struct {
		value    []byte
		capacity int
		borrowed bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "nil",
			args: args{
				v: nil,
			},
			want: want{
				value:    nil,
				capacity: 0,
				borrowed: true,
			},
		},
		{
			name: "empty",
			args: args{
				v: []byte{},
			},
			want: want{
				value:    []byte{},
				capacity: 0,
				borrowed: true,
			},
		},
		{
			name: "text",
			args: args{
				v: []byte("abc"),
			},
			want: want{
				value:    []byte("abc"),
				capacity: 3,
				borrowed: true,
			},
		},
		{
			name: "spare capacity",
			args: args{
				v: []byte{1, 2, 3, 4}[:2],
			},
			want: want{
				value:    []byte{1, 2},
				capacity: 2,
				borrowed: true,
			},
		},
		{
			name: "empty with capacity",
			args: args{
				v: make([]byte, 0, 8),
			},
			want: want{
				value:    []byte{},
				capacity: 0,
				borrowed: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value := Bytes(test.args.v)
			bytes := value.AsBytes()
			got := want{
				value:    bytes,
				capacity: cap(bytes),
				borrowed: reflect.ValueOf(bytes).Pointer() == reflect.ValueOf(test.args.v).Pointer(),
			}
			testutil.AssertValue(t, got, test.want, "Bytes")
			testutil.AssertValue(t, value.AsAny(), test.want.value, "Bytes AsAny")
		})
	}
}

func TestInt(t *testing.T) {
	type args struct {
		v int
	}
	type want struct {
		value int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "zero",
			args: args{
				v: 0,
			},
			want: want{
				value: 0,
			},
		},
		{
			name: "minimum",
			args: args{
				v: math.MinInt,
			},
			want: want{
				value: math.MinInt,
			},
		},
		{
			name: "maximum",
			args: args{
				v: math.MaxInt,
			},
			want: want{
				value: math.MaxInt,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Int(test.args.v)
			testutil.AssertValue(t, got.AsAny(), test.want.value, "Int")
		})
	}
}

func TestInt8(t *testing.T) {
	type args struct {
		v int8
	}
	type want struct {
		value int8
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "zero",
			args: args{
				v: 0,
			},
			want: want{
				value: 0,
			},
		},
		{
			name: "minimum",
			args: args{
				v: math.MinInt8,
			},
			want: want{
				value: math.MinInt8,
			},
		},
		{
			name: "maximum",
			args: args{
				v: math.MaxInt8,
			},
			want: want{
				value: math.MaxInt8,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Int8(test.args.v)
			testutil.AssertValue(t, got.AsAny(), test.want.value, "Int8")
		})
	}
}

func TestInt16(t *testing.T) {
	type args struct {
		v int16
	}
	type want struct {
		value int16
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "zero",
			args: args{
				v: 0,
			},
			want: want{
				value: 0,
			},
		},
		{
			name: "minimum",
			args: args{
				v: math.MinInt16,
			},
			want: want{
				value: math.MinInt16,
			},
		},
		{
			name: "maximum",
			args: args{
				v: math.MaxInt16,
			},
			want: want{
				value: math.MaxInt16,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Int16(test.args.v)
			testutil.AssertValue(t, got.AsAny(), test.want.value, "Int16")
		})
	}
}

func TestInt32(t *testing.T) {
	type args struct {
		v int32
	}
	type want struct {
		value int32
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "zero",
			args: args{
				v: 0,
			},
			want: want{
				value: 0,
			},
		},
		{
			name: "minimum",
			args: args{
				v: math.MinInt32,
			},
			want: want{
				value: math.MinInt32,
			},
		},
		{
			name: "maximum",
			args: args{
				v: math.MaxInt32,
			},
			want: want{
				value: math.MaxInt32,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Int32(test.args.v)
			testutil.AssertValue(t, got.AsAny(), test.want.value, "Int32")
		})
	}
}

func TestInt64(t *testing.T) {
	type args struct {
		v int64
	}
	type want struct {
		value int64
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "zero",
			args: args{
				v: 0,
			},
			want: want{
				value: 0,
			},
		},
		{
			name: "minimum",
			args: args{
				v: math.MinInt64,
			},
			want: want{
				value: math.MinInt64,
			},
		},
		{
			name: "maximum",
			args: args{
				v: math.MaxInt64,
			},
			want: want{
				value: math.MaxInt64,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Int64(test.args.v)
			testutil.AssertValue(t, got.AsAny(), test.want.value, "Int64")
		})
	}
}

func TestUint(t *testing.T) {
	type args struct {
		v uint
	}
	type want struct {
		value uint
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "zero",
			args: args{
				v: 0,
			},
			want: want{
				value: 0,
			},
		},
		{
			name: "maximum",
			args: args{
				v: math.MaxUint,
			},
			want: want{
				value: math.MaxUint,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Uint(test.args.v)
			testutil.AssertValue(t, got.AsAny(), test.want.value, "Uint")
		})
	}
}

func TestUint8(t *testing.T) {
	type args struct {
		v uint8
	}
	type want struct {
		value uint8
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "zero",
			args: args{
				v: 0,
			},
			want: want{
				value: 0,
			},
		},
		{
			name: "maximum",
			args: args{
				v: math.MaxUint8,
			},
			want: want{
				value: math.MaxUint8,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Uint8(test.args.v)
			testutil.AssertValue(t, got.AsAny(), test.want.value, "Uint8")
		})
	}
}

func TestUint16(t *testing.T) {
	type args struct {
		v uint16
	}
	type want struct {
		value uint16
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "zero",
			args: args{
				v: 0,
			},
			want: want{
				value: 0,
			},
		},
		{
			name: "maximum",
			args: args{
				v: math.MaxUint16,
			},
			want: want{
				value: math.MaxUint16,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Uint16(test.args.v)
			testutil.AssertValue(t, got.AsAny(), test.want.value, "Uint16")
		})
	}
}

func TestUint32(t *testing.T) {
	type args struct {
		v uint32
	}
	type want struct {
		value uint32
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "zero",
			args: args{
				v: 0,
			},
			want: want{
				value: 0,
			},
		},
		{
			name: "maximum",
			args: args{
				v: math.MaxUint32,
			},
			want: want{
				value: math.MaxUint32,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Uint32(test.args.v)
			testutil.AssertValue(t, got.AsAny(), test.want.value, "Uint32")
		})
	}
}

func TestUint64(t *testing.T) {
	type args struct {
		v uint64
	}
	type want struct {
		value uint64
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "zero",
			args: args{
				v: 0,
			},
			want: want{
				value: 0,
			},
		},
		{
			name: "maximum",
			args: args{
				v: math.MaxUint64,
			},
			want: want{
				value: math.MaxUint64,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Uint64(test.args.v)
			testutil.AssertValue(t, got.AsAny(), test.want.value, "Uint64")
		})
	}
}

func TestUintptr(t *testing.T) {
	type args struct {
		v uintptr
	}
	type want struct {
		value uintptr
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "zero",
			args: args{
				v: 0,
			},
			want: want{
				value: 0,
			},
		},
		{
			name: "maximum",
			args: args{
				v: ^uintptr(0),
			},
			want: want{
				value: ^uintptr(0),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Uintptr(test.args.v)
			testutil.AssertValue(t, got.AsAny(), test.want.value, "Uintptr")
		})
	}
}

func TestFloat32(t *testing.T) {
	type args struct {
		v float32
	}
	type want struct {
		value float32
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "zero",
			args: args{
				v: 0,
			},
			want: want{
				value: 0,
			},
		},
		{
			name: "fraction",
			args: args{
				v: -1.25,
			},
			want: want{
				value: -1.25,
			},
		},
		{
			name: "maximum",
			args: args{
				v: math.MaxFloat32,
			},
			want: want{
				value: math.MaxFloat32,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Float32(test.args.v)
			testutil.AssertValue(t, got.AsAny(), test.want.value, "Float32")
		})
	}
}

func TestFloat64(t *testing.T) {
	type args struct {
		v float64
	}
	type want struct {
		value float64
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "zero",
			args: args{
				v: 0,
			},
			want: want{
				value: 0,
			},
		},
		{
			name: "fraction",
			args: args{
				v: -1.25,
			},
			want: want{
				value: -1.25,
			},
		},
		{
			name: "maximum",
			args: args{
				v: math.MaxFloat64,
			},
			want: want{
				value: math.MaxFloat64,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Float64(test.args.v)
			testutil.AssertValue(t, got.AsAny(), test.want.value, "Float64")
		})
	}
}

func TestBool(t *testing.T) {
	type args struct {
		v bool
	}
	type want struct {
		value bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "false",
			args: args{
				v: false,
			},
			want: want{
				value: false,
			},
		},
		{
			name: "true",
			args: args{
				v: true,
			},
			want: want{
				value: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Bool(test.args.v)
			testutil.AssertValue(t, got.AsAny(), test.want.value, "Bool")
		})
	}
}

func TestAny(t *testing.T) {
	type args struct {
		v any
	}
	type want struct {
		value any
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "nil",
			args: args{
				v: nil,
			},
			want: want{
				value: nil,
			},
		},
		{
			name: "typed nil",
			args: args{
				v: (*int)(nil),
			},
			want: want{
				value: (*int)(nil),
			},
		},
		{
			name: "named type",
			args: args{
				v: time.Duration(123),
			},
			want: want{
				value: time.Duration(123),
			},
		},
		{
			name: "stringer",
			args: args{
				v: testutil.Stringer{Value: "text"},
			},
			want: want{
				value: testutil.Stringer{Value: "text"},
			},
		},
		{
			name: "map",
			args: args{
				v: map[string]int{"key": 42},
			},
			want: want{
				value: map[string]int{"key": 42},
			},
		},
		{
			name: "bytes",
			args: args{
				v: []byte{1, 2, 3},
			},
			want: want{
				value: []byte{1, 2, 3},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Any(test.args.v)
			testutil.AssertValue(t, got.AsAny(), test.want.value, "Any")
		})
	}
}
