package scope

import (
	"testing"

	"github.com/nekrassov01/table/internal/testutil"
)

func TestScopes_Set(t *testing.T) {
	type fields struct {
		header string
		body   string
		footer string
	}
	type args struct {
		scope Scope
		value string
	}
	type want struct {
		header string
		body   string
		footer string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "no part changes nothing",
			fields: fields{
				header: "header",
				body:   "body",
				footer: "footer",
			},
			args: args{
				scope: 0,
				value: "value",
			},
			want: want{
				header: "header",
				body:   "body",
				footer: "footer",
			},
		},
		{
			name: "one part changes alone",
			fields: fields{
				header: "header",
				body:   "body",
				footer: "footer",
			},
			args: args{
				scope: Body,
				value: "value",
			},
			want: want{
				header: "header",
				body:   "value",
				footer: "footer",
			},
		},
		{
			name: "combination changes each named part",
			fields: fields{
				header: "header",
				body:   "body",
				footer: "footer",
			},
			args: args{
				scope: Header | Footer,
				value: "value",
			},
			want: want{
				header: "value",
				body:   "body",
				footer: "value",
			},
		},
		{
			name: "every part changes",
			args: args{
				scope: Header | Body | Footer,
				value: "value",
			},
			want: want{
				header: "value",
				body:   "value",
				footer: "value",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := Scopes[string]{
				header: test.fields.header,
				body:   test.fields.body,
				footer: test.fields.footer,
			}
			o.Set(test.args.scope, test.args.value)
			got := want{
				header: o.header,
				body:   o.body,
				footer: o.footer,
			}
			testutil.AssertValue(t, got, test.want, "Set")
		})
	}
}

func TestScopes_Resolve(t *testing.T) {
	type fields struct {
		header string
		body   string
		footer string
	}
	type args struct {
		scope Scope
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
			name: "header",
			fields: fields{
				header: "header",
				body:   "body",
				footer: "footer",
			},
			args: args{
				scope: Header,
			},
			want: want{
				val: "header",
			},
		},
		{
			name: "body",
			fields: fields{
				header: "header",
				body:   "body",
				footer: "footer",
			},
			args: args{
				scope: Body,
			},
			want: want{
				val: "body",
			},
		},
		{
			name: "footer",
			fields: fields{
				header: "header",
				body:   "body",
				footer: "footer",
			},
			args: args{
				scope: Footer,
			},
			want: want{
				val: "footer",
			},
		},
		{
			name: "combination returns zero",
			fields: fields{
				header: "header",
				body:   "body",
				footer: "footer",
			},
			args: args{
				scope: Header | Body,
			},
			want: want{
				val: "",
			},
		},
		{
			name: "no part returns zero",
			fields: fields{
				header: "header",
				body:   "body",
				footer: "footer",
			},
			want: want{
				val: "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := Scopes[string]{
				header: test.fields.header,
				body:   test.fields.body,
				footer: test.fields.footer,
			}
			got := o.Resolve(test.args.scope)
			testutil.AssertValue(t, got, test.want.val, "Resolve")
		})
	}
}

func TestMasks_Mark(t *testing.T) {
	type fields struct {
		header uint64
		body   uint64
		footer uint64
	}
	type args struct {
		scope    Scope
		position int
	}
	type want struct {
		header uint64
		body   uint64
		footer uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "no part reaches nothing",
			args: args{
				scope:    0,
				position: 0,
			},
			want: want{
				header: 0,
				body:   0,
				footer: 0,
			},
		},
		{
			name: "one part takes the bit alone",
			args: args{
				scope:    Body,
				position: 2,
			},
			want: want{
				header: 0,
				body:   0b100,
				footer: 0,
			},
		},
		{
			name: "a combination takes the bit in each part it names",
			args: args{
				scope:    Header | Footer,
				position: 1,
			},
			want: want{
				header: 0b10,
				body:   0,
				footer: 0b10,
			},
		},
		{
			name: "every part takes the bit everywhere",
			args: args{
				scope:    Header | Body | Footer,
				position: 0,
			},
			want: want{
				header: 0b1,
				body:   0b1,
				footer: 0b1,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m := Masks{
				header: test.fields.header,
				body:   test.fields.body,
				footer: test.fields.footer,
			}
			m.Mark(test.args.scope, test.args.position)
			testutil.AssertValue(t, m.Resolve(Header), test.want.header, "header")
			testutil.AssertValue(t, m.Resolve(Body), test.want.body, "body")
			testutil.AssertValue(t, m.Resolve(Footer), test.want.footer, "footer")
		})
	}
}

func TestMasks_Resolve(t *testing.T) {
	type fields struct {
		header uint64
		body   uint64
		footer uint64
	}
	type args struct {
		scope Scope
	}
	type want struct {
		mask uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "each part answers with its own mask",
			fields: fields{
				header: 0b1,
				body:   0b10,
				footer: 0b100,
			},
			args: args{
				scope: Body,
			},
			want: want{
				mask: 0b10,
			},
		},
		{
			name: "a combination names no single part and reaches nothing",
			fields: fields{
				header: 0b1,
				body:   0b10,
				footer: 0b100,
			},
			args: args{
				scope: Header | Body,
			},
			want: want{
				mask: 0,
			},
		},
		{
			name: "no part reaches nothing",
			fields: fields{
				header: 0b1,
				body:   0b10,
				footer: 0b100,
			},
			args: args{
				scope: 0,
			},
			want: want{
				mask: 0,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m := Masks{
				header: test.fields.header,
				body:   test.fields.body,
				footer: test.fields.footer,
			}
			testutil.AssertValue(t, m.Resolve(test.args.scope), test.want.mask, "Resolve")
		})
	}
}
