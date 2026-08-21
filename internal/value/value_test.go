package value

import (
	"errors"
	"net"
	"reflect"
	"testing"

	"github.com/nekrassov01/table/internal/testutil"
)

func TestNumber(t *testing.T) {
	type args struct {
		x int64
	}
	type want struct {
		val string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "first",
			args: args{
				x: 1,
			},
			want: want{
				val: "1",
			},
		},
		{
			name: "wide",
			args: args{
				x: 1000,
			},
			want: want{
				val: "1000",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var st Store
			testutil.AssertValue(t, Number(&st, test.args.x), test.want.val, "Number")
		})
	}
}

func TestFormat(t *testing.T) {
	type args struct {
		v any
	}
	type want struct {
		text string
	}
	type blob []byte
	type list []int
	type dictionary map[string]int
	type record struct {
		first  int
		second int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "nil is missing",
			args: args{
				v: nil,
			},
			want: want{
				text: "",
			},
		},
		{
			name: "string",
			args: args{
				v: "hello",
			},
			want: want{
				text: "hello",
			},
		},
		{
			name: "empty string is missing",
			args: args{
				v: "",
			},
			want: want{
				text: "",
			},
		},
		{
			name: "int",
			args: args{
				v: -42,
			},
			want: want{
				text: "-42",
			},
		},
		{
			name: "int8",
			args: args{
				v: int8(-8),
			},
			want: want{
				text: "-8",
			},
		},
		{
			name: "int16",
			args: args{
				v: int16(-16),
			},
			want: want{
				text: "-16",
			},
		},
		{
			name: "int32",
			args: args{
				v: int32(-32),
			},
			want: want{
				text: "-32",
			},
		},
		{
			name: "int64",
			args: args{
				v: int64(-64),
			},
			want: want{
				text: "-64",
			},
		},
		{
			name: "uint",
			args: args{
				v: uint(1),
			},
			want: want{
				text: "1",
			},
		},
		{
			name: "uint8",
			args: args{
				v: uint8(255),
			},
			want: want{
				text: "255",
			},
		},
		{
			name: "uint16",
			args: args{
				v: uint16(16),
			},
			want: want{
				text: "16",
			},
		},
		{
			name: "uint32",
			args: args{
				v: uint32(32),
			},
			want: want{
				text: "32",
			},
		},
		{
			name: "uint64",
			args: args{
				v: uint64(64),
			},
			want: want{
				text: "64",
			},
		},
		{
			name: "uintptr",
			args: args{
				v: uintptr(128),
			},
			want: want{
				text: "128",
			},
		},
		{
			name: "float64 shortest",
			args: args{
				v: 1.25,
			},
			want: want{
				text: "1.25",
			},
		},
		{
			name: "float32 rounding",
			args: args{
				v: float32(0.1),
			},
			want: want{
				text: "0.1",
			},
		},
		{
			name: "bool",
			args: args{
				v: true,
			},
			want: want{
				text: "true",
			},
		},
		{
			name: "bytes as text",
			args: args{
				v: []byte("raw"),
			},
			want: want{
				text: "raw",
			},
		},
		{
			name: "empty bytes are missing",
			args: args{
				v: []byte{},
			},
			want: want{
				text: "",
			},
		},
		{
			name: "stringer",
			args: args{
				v: net.IPv4(127, 0, 0, 1),
			},
			want: want{
				text: "127.0.0.1",
			},
		},
		{
			name: "typed-nil stringer is missing",
			args: args{
				v: net.IP(nil),
			},
			want: want{
				text: "",
			},
		},
		{
			name: "error",
			args: args{
				v: errors.New("boom"),
			},
			want: want{
				text: "boom",
			},
		},
		{
			name: "typed-nil error is missing",
			args: args{
				v: (*testutil.PtrError)(nil),
			},
			want: want{
				text: "",
			},
		},
		{
			name: "pointer dereferences",
			args: args{
				v: func() any {
					value := 3
					return &value
				}(),
			},
			want: want{
				text: "3",
			},
		},
		{
			name: "nil pointer is missing",
			args: args{
				v: (*int)(nil),
			},
			want: want{
				text: "",
			},
		},
		{
			name: "empty struct summary",
			args: args{
				v: struct{}{},
			},
			want: want{
				text: "{struct 0 field(s)}",
			},
		},
		{
			name: "anonymous struct summary",
			args: args{
				v: struct {
					A int
					B int
				}{},
			},
			want: want{
				text: "{struct 2 field(s)}",
			},
		},
		{
			name: "named struct summary",
			args: args{
				v: record{},
			},
			want: want{
				text: "{struct 2 field(s)}",
			},
		},
		{
			name: "pointer to struct summary",
			args: args{
				v: &record{},
			},
			want: want{
				text: "{struct 2 field(s)}",
			},
		},
		{
			name: "nil map is missing",
			args: args{
				v: map[string]int(nil),
			},
			want: want{
				text: "",
			},
		},
		{
			name: "empty map summary",
			args: args{
				v: map[string]int{},
			},
			want: want{
				text: "{map 0 key(s)}",
			},
		},
		{
			name: "map with one key summary",
			args: args{
				v: map[string]int{"a": 1},
			},
			want: want{
				text: "{map 1 key(s)}",
			},
		},
		{
			name: "map with multiple keys summary",
			args: args{
				v: map[string]int{
					"a": 1,
					"b": 2,
				},
			},
			want: want{
				text: "{map 2 key(s)}",
			},
		},
		{
			name: "named map summary",
			args: args{
				v: dictionary{
					"a": 1,
				},
			},
			want: want{
				text: "{map 1 key(s)}",
			},
		},
		{
			name: "pointer to map summary",
			args: args{
				v: func() any {
					value := map[string]int{
						"a": 1,
					}
					return &value
				}(),
			},
			want: want{
				text: "{map 1 key(s)}",
			},
		},
		{
			name: "nil slice is missing",
			args: args{
				v: []int(nil),
			},
			want: want{
				text: "",
			},
		},
		{
			name: "slice of ints collapses to a summary",
			args: args{
				v: []int{1, 2, 3},
			},
			want: want{
				text: "[list 3 item(s)]",
			},
		},
		{
			name: "named slice collapses to a summary",
			args: args{
				v: list{1, 2},
			},
			want: want{
				text: "[list 2 item(s)]",
			},
		},
		{
			name: "slice of strings collapses to a summary",
			args: args{
				v: []string{"a", "", "c"},
			},
			want: want{
				text: "[list 3 item(s)]",
			},
		},
		{
			name: "empty slice is missing",
			args: args{
				v: []int{},
			},
			want: want{
				text: "",
			},
		},
		{
			name: "count beyond the pre-rendered summaries is built on demand",
			args: args{
				v: make([]int, summaryCacheSize+1),
			},
			want: want{
				text: "[list 65 item(s)]",
			},
		},
		{
			name: "nested slices collapse to a summary",
			args: args{
				v: [][]int{{1}, {2}},
			},
			want: want{
				text: "[list 2 item(s)]",
			},
		},
		{
			name: "empty array is missing",
			args: args{
				v: [0]int{},
			},
			want: want{
				text: "",
			},
		},
		{
			name: "array collapses to a summary",
			args: args{
				v: [2]int{1, 2},
			},
			want: want{
				text: "[list 2 item(s)]",
			},
		},
		{
			name: "byte array is not special",
			args: args{
				v: [2]byte{7, 200},
			},
			want: want{
				text: "[list 2 item(s)]",
			},
		},
		{
			name: "rune slice is not special",
			args: args{
				v: []rune("ab"),
			},
			want: want{
				text: "[list 2 item(s)]",
			},
		},
		{
			name: "named byte slice as text",
			args: args{
				v: blob("ab"),
			},
			want: want{
				text: "ab",
			},
		},
		{
			name: "slice of any collapses to a summary",
			args: args{
				v: []any{"x", 1, nil},
			},
			want: want{
				text: "[list 3 item(s)]",
			},
		},
		{
			name: "pointer to slice collapses to a summary",
			args: args{
				v: &[]int{4, 5},
			},
			want: want{
				text: "[list 2 item(s)]",
			},
		},
		{
			name: "other kind uses default formatting",
			args: args{
				v: complex(1, 2),
			},
			want: want{
				text: "(1+2i)",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var st Store
			testutil.AssertValue(t, Format(&st, test.args.v), test.want.text, "text")
		})
	}
}

func Test_formatReflect(t *testing.T) {
	type args struct {
		value any
	}
	type want struct {
		text string
	}
	type namedInt int
	type namedBytes []byte
	type namedList []int
	type namedMap map[string]int
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "invalid value is missing",
			want: want{
				text: "",
			},
		},
		{
			name: "reference resolves",
			args: args{
				value: &testutil.PtrStringer{
					Value: "stringer",
				},
			},
			want: want{
				text: "stringer",
			},
		},
		{
			name: "primitive kind resolves",
			args: args{
				value: namedInt(42),
			},
			want: want{
				text: "42",
			},
		},
		{
			name: "empty struct summarizes fields",
			args: args{
				value: struct{}{},
			},
			want: want{
				text: "{struct 0 field(s)}",
			},
		},
		{
			name: "struct summarizes fields",
			args: args{
				value: struct {
					First  int
					Second int
				}{},
			},
			want: want{
				text: "{struct 2 field(s)}",
			},
		},
		{
			name: "nil map is missing",
			args: args{
				value: map[string]int(nil),
			},
			want: want{
				text: "",
			},
		},
		{
			name: "empty map summarizes keys",
			args: args{
				value: map[string]int{},
			},
			want: want{
				text: "{map 0 key(s)}",
			},
		},
		{
			name: "map summarizes keys",
			args: args{
				value: map[string]int{
					"key": 1,
				},
			},
			want: want{
				text: "{map 1 key(s)}",
			},
		},
		{
			name: "named map summarizes keys",
			args: args{
				value: namedMap{
					"key": 1,
				},
			},
			want: want{
				text: "{map 1 key(s)}",
			},
		},
		{
			name: "nil list is missing",
			args: args{
				value: []int(nil),
			},
			want: want{
				text: "",
			},
		},
		{
			name: "empty list is missing",
			args: args{
				value: []int{},
			},
			want: want{
				text: "",
			},
		},
		{
			name: "slice summarizes items",
			args: args{
				value: []int{1, 2},
			},
			want: want{
				text: "[list 2 item(s)]",
			},
		},
		{
			name: "named slice summarizes items",
			args: args{
				value: namedList{1, 2},
			},
			want: want{
				text: "[list 2 item(s)]",
			},
		},
		{
			name: "named byte slice is text",
			args: args{
				value: namedBytes("bytes"),
			},
			want: want{
				text: "bytes",
			},
		},
		{
			name: "empty array is missing",
			args: args{
				value: [0]int{},
			},
			want: want{
				text: "",
			},
		},
		{
			name: "array summarizes items",
			args: args{
				value: [2]int{1, 2},
			},
			want: want{
				text: "[list 2 item(s)]",
			},
		},
		{
			name: "other kind uses default formatting",
			args: args{
				value: complex(1, 2),
			},
			want: want{
				text: "(1+2i)",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var st Store
			got := formatReflect(&st, test.args.value)
			testutil.AssertValue(t, got, test.want.text, "formatReflect")
		})
	}
}

func Test_resolveReference(t *testing.T) {
	type args struct {
		value func() reflect.Value
	}
	type want struct {
		kind     reflect.Kind
		text     string
		resolved bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "nil pointer is missing",
			args: args{
				value: func() reflect.Value {
					return reflect.ValueOf((*int)(nil))
				},
			},
			want: want{
				kind:     reflect.Pointer,
				resolved: true,
			},
		},
		{
			name: "nil map after dereference is missing",
			args: args{
				value: func() reflect.Value {
					value := map[string]int(nil)
					return reflect.ValueOf(&value)
				},
			},
			want: want{
				kind:     reflect.Map,
				resolved: true,
			},
		},
		{
			name: "pointer stringer resolves before dereference",
			args: args{
				value: func() reflect.Value {
					return reflect.ValueOf(&testutil.PtrStringer{
						Value: "stringer",
					})
				},
			},
			want: want{
				kind:     reflect.Pointer,
				text:     "stringer",
				resolved: true,
			},
		},
		{
			name: "interface and pointer chain unwraps",
			args: args{
				value: func() reflect.Value {
					value := 42
					var wrapped any = &value
					return reflect.ValueOf(&wrapped).Elem()
				},
			},
			want: want{
				kind: reflect.Int,
			},
		},
		{
			name: "value stringer resolves after unwrapping",
			args: args{
				value: func() reflect.Value {
					return reflect.ValueOf(testutil.Stringer{
						Value: "stringer",
					})
				},
			},
			want: want{
				kind:     reflect.Struct,
				text:     "stringer",
				resolved: true,
			},
		},
		{
			name: "ordinary value remains unresolved",
			args: args{
				value: func() reflect.Value {
					return reflect.ValueOf(42)
				},
			},
			want: want{
				kind: reflect.Int,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotValue, gotText, gotResolved := resolveReference(test.args.value())
			got := want{
				kind:     gotValue.Kind(),
				text:     gotText,
				resolved: gotResolved,
			}
			testutil.AssertValue(t, got, test.want, "resolveReference")
		})
	}
}

func Test_resolveStringer(t *testing.T) {
	type args struct {
		value func() reflect.Value
	}
	type want struct {
		text     string
		resolved bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "inaccessible value remains unresolved",
			args: args{
				value: func() reflect.Value {
					return reflect.ValueOf(struct {
						value string
					}{
						value: "hidden",
					}).Field(0)
				},
			},
		},
		{
			name: "stringer resolves",
			args: args{
				value: func() reflect.Value {
					return reflect.ValueOf(testutil.Stringer{
						Value: "stringer",
					})
				},
			},
			want: want{
				text:     "stringer",
				resolved: true,
			},
		},
		{
			name: "error resolves",
			args: args{
				value: func() reflect.Value {
					return reflect.ValueOf(testutil.Error{
						Value: "error",
					})
				},
			},
			want: want{
				text:     "error",
				resolved: true,
			},
		},
		{
			name: "ordinary value remains unresolved",
			args: args{
				value: func() reflect.Value {
					return reflect.ValueOf(42)
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotText, gotResolved := resolveStringer(test.args.value())
			got := want{
				text:     gotText,
				resolved: gotResolved,
			}
			testutil.AssertValue(t, got, test.want, "resolveStringer")
		})
	}
}

func Test_resolveKind(t *testing.T) {
	type args struct {
		value any
	}
	type want struct {
		text     string
		resolved bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "string",
			args: args{
				value: "value",
			},
			want: want{
				text:     "value",
				resolved: true,
			},
		},
		{
			name: "bool",
			args: args{
				value: true,
			},
			want: want{
				text:     "true",
				resolved: true,
			},
		},
		{
			name: "int",
			args: args{
				value: int16(-16),
			},
			want: want{
				text:     "-16",
				resolved: true,
			},
		},
		{
			name: "uint",
			args: args{
				value: uint16(16),
			},
			want: want{
				text:     "16",
				resolved: true,
			},
		},
		{
			name: "float32",
			args: args{
				value: float32(0.1),
			},
			want: want{
				text:     "0.1",
				resolved: true,
			},
		},
		{
			name: "float64",
			args: args{
				value: 1.25,
			},
			want: want{
				text:     "1.25",
				resolved: true,
			},
		},
		{
			name: "other kind remains unresolved",
			args: args{
				value: complex(1, 2),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var st Store
			gotText, gotResolved := resolveKind(&st, reflect.ValueOf(test.args.value))
			got := want{
				text:     gotText,
				resolved: gotResolved,
			}
			testutil.AssertValue(t, got, test.want, "resolveKind")
		})
	}
}

func Test_appendBytes(t *testing.T) {
	type args struct {
		value []byte
	}
	type want struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "appends bytes",
			args: args{
				value: []byte("value"),
			},
			want: want{
				text: "value",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var st Store
			got := appendBytes(&st, test.args.value)
			testutil.AssertValue(t, got, test.want.text, "appendBytes")
		})
	}
}

func Test_appendInt(t *testing.T) {
	type args struct {
		value int64
	}
	type want struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "appends integer",
			args: args{
				value: -42,
			},
			want: want{
				text: "-42",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var st Store
			got := appendInt(&st, test.args.value)
			testutil.AssertValue(t, got, test.want.text, "appendInt")
		})
	}
}

func Test_appendUint(t *testing.T) {
	type args struct {
		value uint64
	}
	type want struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "appends unsigned integer",
			args: args{
				value: 42,
			},
			want: want{
				text: "42",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var st Store
			got := appendUint(&st, test.args.value)
			testutil.AssertValue(t, got, test.want.text, "appendUint")
		})
	}
}

func Test_appendFloat(t *testing.T) {
	type args struct {
		value   float64
		bitSize int
	}
	type want struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "appends floating point",
			args: args{
				value:   1.25,
				bitSize: 64,
			},
			want: want{
				text: "1.25",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var st Store
			got := appendFloat(&st, test.args.value, test.args.bitSize)
			testutil.AssertValue(t, got, test.want.text, "appendFloat")
		})
	}
}

func Test_isTypedNil(t *testing.T) {
	type args struct {
		value any
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
			name: "nil interface is not typed nil",
			want: want{
				val: false,
			},
		},
		{
			name: "nil pointer",
			args: args{
				value: (*int)(nil),
			},
			want: want{
				val: true,
			},
		},
		{
			name: "nil map",
			args: args{
				value: map[string]int(nil),
			},
			want: want{
				val: true,
			},
		},
		{
			name: "nil slice",
			args: args{
				value: []int(nil),
			},
			want: want{
				val: true,
			},
		},
		{
			name: "nil channel",
			args: args{
				value: chan int(nil),
			},
			want: want{
				val: true,
			},
		},
		{
			name: "nil function",
			args: args{
				value: (func())(nil),
			},
			want: want{
				val: true,
			},
		},
		{
			name: "non-nil pointer",
			args: args{
				value: new(int),
			},
			want: want{
				val: false,
			},
		},
		{
			name: "ordinary value",
			args: args{
				value: 42,
			},
			want: want{
				val: false,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := isTypedNil(test.args.value)
			testutil.AssertValue(t, got, test.want.val, "isTypedNil")
		})
	}
}
