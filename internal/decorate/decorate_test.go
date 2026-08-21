package decorate

import (
	"testing"

	"github.com/nekrassov01/table/internal/testutil"
)

func TestNew(t *testing.T) {
	type args struct {
		prefix string
		suffix string
	}
	type want struct {
		val *Decoration
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "empty prefix returns nil",
			args: args{
				prefix: "",
				suffix: "</strong>",
			},
			want: want{
				val: nil,
			},
		},
		{
			name: "non-empty prefix returns decoration",
			args: args{
				prefix: "<strong>",
				suffix: "</strong>",
			},
			want: want{
				val: &Decoration{
					Prefix: "<strong>",
					Suffix: "</strong>",
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := New(test.args.prefix, test.args.suffix)
			testutil.AssertValue(t, got, test.want.val, "New")
		})
	}
}

func TestDecoration_IsZero(t *testing.T) {
	type fields struct {
		Prefix string
		Suffix string
	}
	type want struct {
		val bool
	}
	tests := []struct {
		name   string
		fields *fields
		want   want
	}{
		{
			name:   "nil receiver returns true",
			fields: nil,
			want: want{
				val: true,
			},
		},
		{
			name:   "empty prefix returns true",
			fields: &fields{},
			want: want{
				val: true,
			},
		},
		{
			name: "non-empty prefix returns false",
			fields: &fields{
				Prefix: "**",
				Suffix: "**",
			},
			want: want{
				val: false,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var o *Decoration
			if test.fields != nil {
				o = &Decoration{
					Prefix: test.fields.Prefix,
					Suffix: test.fields.Suffix,
				}
			}
			got := o.IsZero()
			testutil.AssertValue(t, got, test.want.val, "IsZero")
		})
	}
}
