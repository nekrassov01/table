package value

import (
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/nekrassov01/table/internal/testutil"
)

func TestValue(t *testing.T) {
	type args struct {
		value    Value
		original any
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
			name: "zero",
			args: args{
				value:    Value{},
				original: nil,
			},
			want: want{
				value: "",
			},
		},
		{
			name: "String",
			args: args{
				value:    String("abc"),
				original: "abc",
			},
			want: want{
				value: "abc",
			},
		},
		{
			name: "Bytes",
			args: args{
				value:    Bytes([]byte("abc")),
				original: []byte("abc"),
			},
			want: want{
				value: "abc",
			},
		},
		{
			name: "Int",
			args: args{
				value:    Int(-1000),
				original: int(-1000),
			},
			want: want{
				value: "-1000",
			},
		},
		{
			name: "Int8",
			args: args{
				value:    Int8(-100),
				original: int8(-100),
			},
			want: want{
				value: "-100",
			},
		},
		{
			name: "Int16",
			args: args{
				value:    Int16(-1000),
				original: int16(-1000),
			},
			want: want{
				value: "-1000",
			},
		},
		{
			name: "Int32",
			args: args{
				value:    Int32(-1000),
				original: int32(-1000),
			},
			want: want{
				value: "-1000",
			},
		},
		{
			name: "Int64",
			args: args{
				value:    Int64(-1000),
				original: int64(-1000),
			},
			want: want{
				value: "-1000",
			},
		},
		{
			name: "Uint",
			args: args{
				value:    Uint(1000),
				original: uint(1000),
			},
			want: want{
				value: "1000",
			},
		},
		{
			name: "Uint8",
			args: args{
				value:    Uint8(200),
				original: uint8(200),
			},
			want: want{
				value: "200",
			},
		},
		{
			name: "Uint16",
			args: args{
				value:    Uint16(1000),
				original: uint16(1000),
			},
			want: want{
				value: "1000",
			},
		},
		{
			name: "Uint32",
			args: args{
				value:    Uint32(1000),
				original: uint32(1000),
			},
			want: want{
				value: "1000",
			},
		},
		{
			name: "Uint64",
			args: args{
				value:    Uint64(1 << 63),
				original: uint64(1 << 63),
			},
			want: want{
				value: "9223372036854775808",
			},
		},
		{
			name: "Uintptr",
			args: args{
				value:    Uintptr(1000),
				original: uintptr(1000),
			},
			want: want{
				value: "1000",
			},
		},
		{
			name: "Float32",
			args: args{
				value:    Float32(1.2),
				original: float32(1.2),
			},
			want: want{
				value: "1.2",
			},
		},
		{
			name: "Float64",
			args: args{
				value:    Float64(1.2),
				original: float64(1.2),
			},
			want: want{
				value: "1.2",
			},
		},
		{
			name: "Bool",
			args: args{
				value:    Bool(true),
				original: true,
			},
			want: want{
				value: "true",
			},
		},
		{
			name: "nil",
			args: args{
				value:    Any(nil),
				original: nil,
			},
			want: want{
				value: "",
			},
		},
		{
			name: "Stringer",
			args: args{
				value: Any(testutil.Stringer{
					Value: "named",
				}),
				original: testutil.Stringer{
					Value: "named",
				},
			},
			want: want{
				value: "named",
			},
		},
		{
			name: "nil bytes",
			args: args{
				value:    Bytes([]byte(nil)),
				original: []byte(nil),
			},
			want: want{
				value: "",
			},
		},
		{
			name: "empty bytes",
			args: args{
				value:    Bytes([]byte{}),
				original: []byte{},
			},
			want: want{
				value: "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var actual Store
			got := test.args.value.AsAny()
			testutil.AssertValue(t, reflect.TypeOf(got), reflect.TypeOf(test.args.original), "type")
			testutil.AssertValue(t, got, test.args.original, "value")
			testutil.AssertValue(t, Format(&actual, test.args.value), test.want.value, "format")
			testutil.AssertValue(t, Format(&actual, Any(test.args.original)), test.want.value, "Any format")
		})
	}
}

func TestValueBorrowedStorage(t *testing.T) {
	type args struct {
		makeValue func() Value
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
			name: "string remains reachable",
			args: args{
				makeValue: func() Value {
					return String(strings.Repeat("ab", 1000))
				},
			},
			want: want{
				value: strings.Repeat("ab", 1000),
			},
		},
		{
			name: "bytes remain reachable",
			args: args{
				makeValue: func() Value {
					return Bytes([]byte(strings.Repeat("cd", 1000)))
				},
			},
			want: want{
				value: []byte(strings.Repeat("cd", 1000)),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value := test.args.makeValue()
			runtime.GC()
			testutil.AssertValue(t, value.AsAny(), test.want.value, "retained value")
		})
	}
}

func TestValueAccessors(t *testing.T) {
	type namedInt int
	type namedString string
	type fields struct {
		value Value
	}
	type args struct {
		read func(Value) any
	}
	type want struct {
		value  any
		panics bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "String/primitive",
			fields: fields{
				value: String("abc"),
			},
			args: args{
				read: func(v Value) any {
					return v.AsString()
				},
			},
			want: want{
				value:  "abc",
				panics: false,
			},
		},
		{
			name: "String/fallback",
			fields: fields{
				value: Any("abc"),
			},
			args: args{
				read: func(v Value) any {
					return v.AsString()
				},
			},
			want: want{
				value:  "abc",
				panics: false,
			},
		},
		{
			name: "String/wrong type",
			fields: fields{
				value: Int(0),
			},
			args: args{
				read: func(v Value) any {
					return v.AsString()
				},
			},
			want: want{
				value:  nil,
				panics: true,
			},
		},
		{
			name: "String/zero",
			fields: fields{
				value: Value{},
			},
			args: args{
				read: func(v Value) any {
					return v.AsString()
				},
			},
			want: want{
				value:  nil,
				panics: true,
			},
		},
		{
			name: "Bytes/primitive",
			fields: fields{
				value: Bytes([]byte("abc")),
			},
			args: args{
				read: func(v Value) any {
					return v.AsBytes()
				},
			},
			want: want{
				value:  []byte("abc"),
				panics: false,
			},
		},
		{
			name: "Bytes/fallback",
			fields: fields{
				value: Any([]byte("abc")),
			},
			args: args{
				read: func(v Value) any {
					return v.AsBytes()
				},
			},
			want: want{
				value:  []byte("abc"),
				panics: false,
			},
		},
		{
			name: "Bytes/wrong type",
			fields: fields{
				value: Int(0),
			},
			args: args{
				read: func(v Value) any {
					return v.AsBytes()
				},
			},
			want: want{
				value:  nil,
				panics: true,
			},
		},
		{
			name: "Bytes/zero",
			fields: fields{
				value: Value{},
			},
			args: args{
				read: func(v Value) any {
					return v.AsBytes()
				},
			},
			want: want{
				value:  nil,
				panics: true,
			},
		},
		{
			name: "Int/positive primitive",
			fields: fields{
				value: Int(1000),
			},
			args: args{
				read: func(v Value) any {
					return v.AsInt()
				},
			},
			want: want{
				value: 1000,
			},
		},
		{
			name: "Int/positive fallback",
			fields: fields{
				value: Any(1000),
			},
			args: args{
				read: func(v Value) any {
					return v.AsInt()
				},
			},
			want: want{
				value: 1000,
			},
		},
		{
			name: "Int/primitive",
			fields: fields{
				value: Int(-1000),
			},
			args: args{
				read: func(v Value) any {
					return v.AsInt()
				},
			},
			want: want{
				value:  int(-1000),
				panics: false,
			},
		},
		{
			name: "Int/fallback",
			fields: fields{
				value: Any(int(-1000)),
			},
			args: args{
				read: func(v Value) any {
					return v.AsInt()
				},
			},
			want: want{
				value:  int(-1000),
				panics: false,
			},
		},
		{
			name: "Int/wrong type",
			fields: fields{
				value: Int8(0),
			},
			args: args{
				read: func(v Value) any {
					return v.AsInt()
				},
			},
			want: want{
				value:  nil,
				panics: true,
			},
		},
		{
			name: "Int/zero",
			fields: fields{
				value: Value{},
			},
			args: args{
				read: func(v Value) any {
					return v.AsInt()
				},
			},
			want: want{
				value:  nil,
				panics: true,
			},
		},
		{
			name: "Int8/primitive",
			fields: fields{
				value: Int8(-100),
			},
			args: args{
				read: func(v Value) any {
					return v.AsInt8()
				},
			},
			want: want{
				value:  int8(-100),
				panics: false,
			},
		},
		{
			name: "Int8/fallback",
			fields: fields{
				value: Any(int8(-100)),
			},
			args: args{
				read: func(v Value) any {
					return v.AsInt8()
				},
			},
			want: want{
				value:  int8(-100),
				panics: false,
			},
		},
		{
			name: "Int8/wrong type",
			fields: fields{
				value: Int(0),
			},
			args: args{
				read: func(v Value) any {
					return v.AsInt8()
				},
			},
			want: want{
				value:  nil,
				panics: true,
			},
		},
		{
			name: "Int8/zero",
			fields: fields{
				value: Value{},
			},
			args: args{
				read: func(v Value) any {
					return v.AsInt8()
				},
			},
			want: want{
				value:  nil,
				panics: true,
			},
		},
		{
			name: "Int16/primitive",
			fields: fields{
				value: Int16(-1000),
			},
			args: args{
				read: func(v Value) any {
					return v.AsInt16()
				},
			},
			want: want{
				value:  int16(-1000),
				panics: false,
			},
		},
		{
			name: "Int16/fallback",
			fields: fields{
				value: Any(int16(-1000)),
			},
			args: args{
				read: func(v Value) any {
					return v.AsInt16()
				},
			},
			want: want{
				value:  int16(-1000),
				panics: false,
			},
		},
		{
			name: "Int16/wrong type",
			fields: fields{
				value: Int(0),
			},
			args: args{
				read: func(v Value) any {
					return v.AsInt16()
				},
			},
			want: want{
				value:  nil,
				panics: true,
			},
		},
		{
			name: "Int16/zero",
			fields: fields{
				value: Value{},
			},
			args: args{
				read: func(v Value) any {
					return v.AsInt16()
				},
			},
			want: want{
				value:  nil,
				panics: true,
			},
		},
		{
			name: "Int32/primitive",
			fields: fields{
				value: Int32(-1000),
			},
			args: args{
				read: func(v Value) any {
					return v.AsInt32()
				},
			},
			want: want{
				value:  int32(-1000),
				panics: false,
			},
		},
		{
			name: "Int32/fallback",
			fields: fields{
				value: Any(int32(-1000)),
			},
			args: args{
				read: func(v Value) any {
					return v.AsInt32()
				},
			},
			want: want{
				value:  int32(-1000),
				panics: false,
			},
		},
		{
			name: "Int32/wrong type",
			fields: fields{
				value: Int(0),
			},
			args: args{
				read: func(v Value) any {
					return v.AsInt32()
				},
			},
			want: want{
				value:  nil,
				panics: true,
			},
		},
		{
			name: "Int32/zero",
			fields: fields{
				value: Value{},
			},
			args: args{
				read: func(v Value) any {
					return v.AsInt32()
				},
			},
			want: want{
				value:  nil,
				panics: true,
			},
		},
		{
			name: "Int64/primitive",
			fields: fields{
				value: Int64(-1000),
			},
			args: args{
				read: func(v Value) any {
					return v.AsInt64()
				},
			},
			want: want{
				value:  int64(-1000),
				panics: false,
			},
		},
		{
			name: "Int64/fallback",
			fields: fields{
				value: Any(int64(-1000)),
			},
			args: args{
				read: func(v Value) any {
					return v.AsInt64()
				},
			},
			want: want{
				value:  int64(-1000),
				panics: false,
			},
		},
		{
			name: "Int64/wrong type",
			fields: fields{
				value: Int(0),
			},
			args: args{
				read: func(v Value) any {
					return v.AsInt64()
				},
			},
			want: want{
				value:  nil,
				panics: true,
			},
		},
		{
			name: "Int64/zero",
			fields: fields{
				value: Value{},
			},
			args: args{
				read: func(v Value) any {
					return v.AsInt64()
				},
			},
			want: want{
				value:  nil,
				panics: true,
			},
		},
		{
			name: "Uint/primitive",
			fields: fields{
				value: Uint(1000),
			},
			args: args{
				read: func(v Value) any {
					return v.AsUint()
				},
			},
			want: want{
				value:  uint(1000),
				panics: false,
			},
		},
		{
			name: "Uint/fallback",
			fields: fields{
				value: Any(uint(1000)),
			},
			args: args{
				read: func(v Value) any {
					return v.AsUint()
				},
			},
			want: want{
				value:  uint(1000),
				panics: false,
			},
		},
		{
			name: "Uint/wrong type",
			fields: fields{
				value: Int(0),
			},
			args: args{
				read: func(v Value) any {
					return v.AsUint()
				},
			},
			want: want{
				value:  nil,
				panics: true,
			},
		},
		{
			name: "Uint/zero",
			fields: fields{
				value: Value{},
			},
			args: args{
				read: func(v Value) any {
					return v.AsUint()
				},
			},
			want: want{
				value:  nil,
				panics: true,
			},
		},
		{
			name: "Uint8/primitive",
			fields: fields{
				value: Uint8(200),
			},
			args: args{
				read: func(v Value) any {
					return v.AsUint8()
				},
			},
			want: want{
				value:  uint8(200),
				panics: false,
			},
		},
		{
			name: "Uint8/fallback",
			fields: fields{
				value: Any(uint8(200)),
			},
			args: args{
				read: func(v Value) any {
					return v.AsUint8()
				},
			},
			want: want{
				value:  uint8(200),
				panics: false,
			},
		},
		{
			name: "Uint8/wrong type",
			fields: fields{
				value: Int(0),
			},
			args: args{
				read: func(v Value) any {
					return v.AsUint8()
				},
			},
			want: want{
				value:  nil,
				panics: true,
			},
		},
		{
			name: "Uint8/zero",
			fields: fields{
				value: Value{},
			},
			args: args{
				read: func(v Value) any {
					return v.AsUint8()
				},
			},
			want: want{
				value:  nil,
				panics: true,
			},
		},
		{
			name: "Uint16/primitive",
			fields: fields{
				value: Uint16(1000),
			},
			args: args{
				read: func(v Value) any {
					return v.AsUint16()
				},
			},
			want: want{
				value:  uint16(1000),
				panics: false,
			},
		},
		{
			name: "Uint16/fallback",
			fields: fields{
				value: Any(uint16(1000)),
			},
			args: args{
				read: func(v Value) any {
					return v.AsUint16()
				},
			},
			want: want{
				value:  uint16(1000),
				panics: false,
			},
		},
		{
			name: "Uint16/wrong type",
			fields: fields{
				value: Int(0),
			},
			args: args{
				read: func(v Value) any {
					return v.AsUint16()
				},
			},
			want: want{
				value:  nil,
				panics: true,
			},
		},
		{
			name: "Uint16/zero",
			fields: fields{
				value: Value{},
			},
			args: args{
				read: func(v Value) any {
					return v.AsUint16()
				},
			},
			want: want{
				value:  nil,
				panics: true,
			},
		},
		{
			name: "Uint32/primitive",
			fields: fields{
				value: Uint32(1000),
			},
			args: args{
				read: func(v Value) any {
					return v.AsUint32()
				},
			},
			want: want{
				value:  uint32(1000),
				panics: false,
			},
		},
		{
			name: "Uint32/fallback",
			fields: fields{
				value: Any(uint32(1000)),
			},
			args: args{
				read: func(v Value) any {
					return v.AsUint32()
				},
			},
			want: want{
				value:  uint32(1000),
				panics: false,
			},
		},
		{
			name: "Uint32/wrong type",
			fields: fields{
				value: Int(0),
			},
			args: args{
				read: func(v Value) any {
					return v.AsUint32()
				},
			},
			want: want{
				value:  nil,
				panics: true,
			},
		},
		{
			name: "Uint32/zero",
			fields: fields{
				value: Value{},
			},
			args: args{
				read: func(v Value) any {
					return v.AsUint32()
				},
			},
			want: want{
				value:  nil,
				panics: true,
			},
		},
		{
			name: "Uint64/primitive",
			fields: fields{
				value: Uint64(1 << 63),
			},
			args: args{
				read: func(v Value) any {
					return v.AsUint64()
				},
			},
			want: want{
				value:  uint64(1 << 63),
				panics: false,
			},
		},
		{
			name: "Uint64/fallback",
			fields: fields{
				value: Any(uint64(1 << 63)),
			},
			args: args{
				read: func(v Value) any {
					return v.AsUint64()
				},
			},
			want: want{
				value:  uint64(1 << 63),
				panics: false,
			},
		},
		{
			name: "Uint64/wrong type",
			fields: fields{
				value: Int(0),
			},
			args: args{
				read: func(v Value) any {
					return v.AsUint64()
				},
			},
			want: want{
				value:  nil,
				panics: true,
			},
		},
		{
			name: "Uint64/zero",
			fields: fields{
				value: Value{},
			},
			args: args{
				read: func(v Value) any {
					return v.AsUint64()
				},
			},
			want: want{
				value:  nil,
				panics: true,
			},
		},
		{
			name: "Uintptr/primitive",
			fields: fields{
				value: Uintptr(1000),
			},
			args: args{
				read: func(v Value) any {
					return v.AsUintptr()
				},
			},
			want: want{
				value:  uintptr(1000),
				panics: false,
			},
		},
		{
			name: "Uintptr/fallback",
			fields: fields{
				value: Any(uintptr(1000)),
			},
			args: args{
				read: func(v Value) any {
					return v.AsUintptr()
				},
			},
			want: want{
				value:  uintptr(1000),
				panics: false,
			},
		},
		{
			name: "Uintptr/wrong type",
			fields: fields{
				value: Int(0),
			},
			args: args{
				read: func(v Value) any {
					return v.AsUintptr()
				},
			},
			want: want{
				value:  nil,
				panics: true,
			},
		},
		{
			name: "Uintptr/zero",
			fields: fields{
				value: Value{},
			},
			args: args{
				read: func(v Value) any {
					return v.AsUintptr()
				},
			},
			want: want{
				value:  nil,
				panics: true,
			},
		},
		{
			name: "Float32/primitive",
			fields: fields{
				value: Float32(1.25),
			},
			args: args{
				read: func(v Value) any {
					return v.AsFloat32()
				},
			},
			want: want{
				value:  float32(1.25),
				panics: false,
			},
		},
		{
			name: "Float32/fallback",
			fields: fields{
				value: Any(float32(1.25)),
			},
			args: args{
				read: func(v Value) any {
					return v.AsFloat32()
				},
			},
			want: want{
				value:  float32(1.25),
				panics: false,
			},
		},
		{
			name: "Float32/wrong type",
			fields: fields{
				value: Int(0),
			},
			args: args{
				read: func(v Value) any {
					return v.AsFloat32()
				},
			},
			want: want{
				value:  nil,
				panics: true,
			},
		},
		{
			name: "Float32/zero",
			fields: fields{
				value: Value{},
			},
			args: args{
				read: func(v Value) any {
					return v.AsFloat32()
				},
			},
			want: want{
				value:  nil,
				panics: true,
			},
		},
		{
			name: "Float64/primitive",
			fields: fields{
				value: Float64(1.25),
			},
			args: args{
				read: func(v Value) any {
					return v.AsFloat64()
				},
			},
			want: want{
				value:  float64(1.25),
				panics: false,
			},
		},
		{
			name: "Float64/fallback",
			fields: fields{
				value: Any(float64(1.25)),
			},
			args: args{
				read: func(v Value) any {
					return v.AsFloat64()
				},
			},
			want: want{
				value:  float64(1.25),
				panics: false,
			},
		},
		{
			name: "Float64/wrong type",
			fields: fields{
				value: Int(0),
			},
			args: args{
				read: func(v Value) any {
					return v.AsFloat64()
				},
			},
			want: want{
				value:  nil,
				panics: true,
			},
		},
		{
			name: "Float64/zero",
			fields: fields{
				value: Value{},
			},
			args: args{
				read: func(v Value) any {
					return v.AsFloat64()
				},
			},
			want: want{
				value:  nil,
				panics: true,
			},
		},
		{
			name: "Bool/primitive",
			fields: fields{
				value: Bool(true),
			},
			args: args{
				read: func(v Value) any {
					return v.AsBool()
				},
			},
			want: want{
				value:  true,
				panics: false,
			},
		},
		{
			name: "Bool/fallback",
			fields: fields{
				value: Any(true),
			},
			args: args{
				read: func(v Value) any {
					return v.AsBool()
				},
			},
			want: want{
				value:  true,
				panics: false,
			},
		},
		{
			name: "Bool/wrong type",
			fields: fields{
				value: Int(0),
			},
			args: args{
				read: func(v Value) any {
					return v.AsBool()
				},
			},
			want: want{
				value:  nil,
				panics: true,
			},
		},
		{
			name: "Bool/zero",
			fields: fields{
				value: Value{},
			},
			args: args{
				read: func(v Value) any {
					return v.AsBool()
				},
			},
			want: want{
				value:  nil,
				panics: true,
			},
		},
		{
			name: "Int/narrow integer",
			fields: fields{
				value: Int16(1),
			},
			args: args{
				read: func(v Value) any {
					return v.AsInt()
				},
			},
			want: want{
				value:  nil,
				panics: true,
			},
		},
		{
			name: "Int/fallback mismatch",
			fields: fields{
				value: Any(int64(1)),
			},
			args: args{
				read: func(v Value) any {
					return v.AsInt()
				},
			},
			want: want{
				value:  nil,
				panics: true,
			},
		},
		{
			name: "String/fallback mismatch",
			fields: fields{
				value: Any([]byte("x")),
			},
			args: args{
				read: func(v Value) any {
					return v.AsString()
				},
			},
			want: want{
				value:  nil,
				panics: true,
			},
		},
		{
			name: "Bytes/nil",
			fields: fields{
				value: Bytes(nil),
			},
			args: args{
				read: func(v Value) any {
					return v.AsBytes()
				},
			},
			want: want{
				value:  []byte(nil),
				panics: false,
			},
		},
		{
			name: "Bytes/fallback nil",
			fields: fields{
				value: Any([]byte(nil)),
			},
			args: args{
				read: func(v Value) any {
					return v.AsBytes()
				},
			},
			want: want{
				value:  []byte(nil),
				panics: false,
			},
		},
		{
			name: "Int/named type",
			fields: fields{
				value: Any(namedInt(1)),
			},
			args: args{
				read: func(v Value) any {
					return v.AsInt()
				},
			},
			want: want{
				panics: true,
			},
		},
		{
			name: "String/named type",
			fields: fields{
				value: Any(namedString("x")),
			},
			args: args{
				read: func(v Value) any {
					return v.AsString()
				},
			},
			want: want{
				panics: true,
			},
		},
		{
			name: "Bytes/primitive capacity",
			fields: fields{
				value: Bytes(make([]byte, 2, 8)),
			},
			args: args{
				read: func(v Value) any {
					return cap(v.AsBytes())
				},
			},
			want: want{
				value: 2,
			},
		},
		{
			name: "Bytes/fallback capacity",
			fields: fields{
				value: Any(make([]byte, 2, 8)),
			},
			args: args{
				read: func(v Value) any {
					return cap(v.AsBytes())
				},
			},
			want: want{
				value: 8,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, panicked := func() (got any, panicked bool) {
				defer func() {
					panicked = recover() != nil
				}()
				return test.args.read(test.fields.value), false
			}()
			testutil.AssertValue(t, got, test.want.value, "value")
			testutil.AssertValue(t, panicked, test.want.panics, "panic")
		})
	}
}
