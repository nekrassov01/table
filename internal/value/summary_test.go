package value

import (
	"testing"

	"github.com/nekrassov01/table/internal/testutil"
)

func Test_newSummary(t *testing.T) {
	type args struct {
		prefix string
		suffix string
	}
	type want struct {
		prefix string
		suffix string
		first  string
		middle string
		last   string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "pre-renders the cache range",
			args: args{
				prefix: "<",
				suffix: ">",
			},
			want: want{
				prefix: "<",
				suffix: ">",
				first:  "<0>",
				middle: "<32>",
				last:   "<64>",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := newSummary(test.args.prefix, test.args.suffix)
			got := want{
				prefix: o.prefix,
				suffix: o.suffix,
				first:  o.entries[0],
				middle: o.entries[summaryCacheSize/2],
				last:   o.entries[summaryCacheSize],
			}
			testutil.AssertValue(t, got, test.want, "newSummary")
		})
	}
}

func Test_summary_format(t *testing.T) {
	type fields struct {
		value summary
	}
	type args struct {
		count int
	}
	type want struct {
		val string
	}
	cachedSummaryFields := fields{
		value: newSummary("<", ">"),
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "struct without fields",
			fields: fields{
				value: summaryStruct,
			},
			args: args{
				count: 0,
			},
			want: want{
				val: "{struct 0 field(s)}",
			},
		},
		{
			name: "struct with fields",
			fields: fields{
				value: summaryStruct,
			},
			args: args{
				count: 2,
			},
			want: want{
				val: "{struct 2 field(s)}",
			},
		},
		{
			name: "empty map",
			fields: fields{
				value: summaryMap,
			},
			args: args{
				count: 0,
			},
			want: want{
				val: "{map 0 key(s)}",
			},
		},
		{
			name: "map with keys",
			fields: fields{
				value: summaryMap,
			},
			args: args{
				count: 2,
			},
			want: want{
				val: "{map 2 key(s)}",
			},
		},
		{
			name: "empty list",
			fields: fields{
				value: summaryList,
			},
			args: args{
				count: 0,
			},
			want: want{
				val: "[list 0 item(s)]",
			},
		},
		{
			name: "list with items",
			fields: fields{
				value: summaryList,
			},
			args: args{
				count: 2,
			},
			want: want{
				val: "[list 2 item(s)]",
			},
		},
		{
			name:   "first cached entry",
			fields: cachedSummaryFields,
			args: args{
				count: 0,
			},
			want: want{
				val: "<0>",
			},
		},
		{
			name:   "last cached entry",
			fields: cachedSummaryFields,
			args: args{
				count: summaryCacheSize,
			},
			want: want{
				val: "<64>",
			},
		},
		{
			name:   "above cache range",
			fields: cachedSummaryFields,
			args: args{
				count: summaryCacheSize + 1,
			},
			want: want{
				val: "<65>",
			},
		},
		{
			name:   "below cache range",
			fields: cachedSummaryFields,
			args: args{
				count: -1,
			},
			want: want{
				val: "<-1>",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := test.fields.value
			got := o.format(test.args.count)
			testutil.AssertValue(t, got, test.want.val, "format")
		})
	}
}

func Test_summary_build(t *testing.T) {
	type fields struct {
		prefix string
		suffix string
	}
	type args struct {
		count int
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
			name: "joins markers and count",
			fields: fields{
				prefix: "prefix",
				suffix: "suffix",
			},
			args: args{
				count: 42,
			},
			want: want{
				val: "prefix42suffix",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &summary{
				prefix: test.fields.prefix,
				suffix: test.fields.suffix,
			}
			got := o.build(test.args.count)
			testutil.AssertValue(t, got, test.want.val, "build")
		})
	}
}
