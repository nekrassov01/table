package html

import (
	"testing"

	"github.com/nekrassov01/table/internal/testutil"
)

func TestAttr_escape(t *testing.T) {
	type fields struct {
		Class string
		Style string
	}
	type want struct {
		val Attr
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "escapes class and style",
			fields: fields{
				Class: `a&"b`,
				Style: `x&"y`,
			},
			want: want{
				val: Attr{
					Class: "a&amp;&quot;b",
					Style: "x&amp;&quot;y",
				},
			},
		},
		{
			name: "empty",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := Attr{
				Class: test.fields.Class,
				Style: test.fields.Style,
			}
			got := o.escape()
			testutil.AssertValue(t, got, test.want.val, "escape")
		})
	}
}

func TestAttr_size(t *testing.T) {
	type fields struct {
		Class string
		Style string
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
			name: "class and style",
			fields: fields{
				Class: "class",
				Style: "style",
			},
			want: want{
				val: 10,
			},
		},
		{
			name: "empty",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := Attr{
				Class: test.fields.Class,
				Style: test.fields.Style,
			}
			got := o.size()
			testutil.AssertValue(t, got, test.want.val, "size")
		})
	}
}

func TestTableAttr_escape(t *testing.T) {
	type fields struct {
		Table   Attr
		Caption Attr
		Header  SectionAttr
		Body    SectionAttr
		Footer  SectionAttr
	}
	type want struct {
		val TableAttr
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "escapes the attribute tree",
			fields: fields{
				Table: Attr{
					Class: `a&"b`,
				},
				Caption: Attr{
					Class: `a&"b`,
				},
				Header: SectionAttr{
					Section: Attr{
						Class: `a&"b`,
					},
				},
				Body: SectionAttr{
					Row: Attr{
						Class: `a&"b`,
					},
				},
				Footer: SectionAttr{
					Cell: Attr{
						Class: `a&"b`,
					},
				},
			},
			want: want{
				val: func() TableAttr {
					clean := Attr{
						Class: "a&amp;&quot;b",
					}
					return TableAttr{
						Table:   clean,
						Caption: clean,
						Header: SectionAttr{
							Section: clean,
						},
						Body: SectionAttr{
							Row: clean,
						},
						Footer: SectionAttr{
							Cell: clean,
						},
					}
				}(),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := TableAttr{
				Table:   test.fields.Table,
				Caption: test.fields.Caption,
				Header:  test.fields.Header,
				Body:    test.fields.Body,
				Footer:  test.fields.Footer,
			}
			got := o.escape()
			testutil.AssertValue(t, got, test.want.val, "escape")
		})
	}
}

func TestTableAttr_resolve(t *testing.T) {
	type fields struct {
		Header SectionAttr
		Body   SectionAttr
		Footer SectionAttr
	}
	type args struct {
		scope Scope
	}
	type want struct {
		val SectionAttr
	}
	header := SectionAttr{
		Section: Attr{
			Class: "header",
		},
	}
	body := SectionAttr{
		Section: Attr{
			Class: "body",
		},
	}
	footer := SectionAttr{
		Section: Attr{
			Class: "footer",
		},
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
				Header: header,
				Body:   body,
				Footer: footer,
			},
			args: args{
				scope: ScopeHeader,
			},
			want: want{
				val: header,
			},
		},
		{
			name: "body",
			fields: fields{
				Header: header,
				Body:   body,
				Footer: footer,
			},
			args: args{
				scope: ScopeBody,
			},
			want: want{
				val: body,
			},
		},
		{
			name: "footer",
			fields: fields{
				Header: header,
				Body:   body,
				Footer: footer,
			},
			args: args{
				scope: ScopeFooter,
			},
			want: want{
				val: footer,
			},
		},
		{
			name: "combined",
			fields: fields{
				Header: header,
				Body:   body,
				Footer: footer,
			},
			args: args{
				scope: ScopeHeader | ScopeBody,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &TableAttr{
				Header: test.fields.Header,
				Body:   test.fields.Body,
				Footer: test.fields.Footer,
			}
			got := o.resolve(test.args.scope)
			testutil.AssertValue(t, got, test.want.val, "resolve")
		})
	}
}

func TestSectionAttr_escape(t *testing.T) {
	type fields struct {
		Section Attr
		Row     Attr
		Cell    Attr
	}
	type want struct {
		val SectionAttr
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "escapes every element",
			fields: fields{
				Section: Attr{
					Class: `s&"`,
				},
				Row: Attr{
					Class: `r&"`,
				},
				Cell: Attr{
					Class: `c&"`,
				},
			},
			want: want{
				val: SectionAttr{
					Section: Attr{
						Class: "s&amp;&quot;",
					},
					Row: Attr{
						Class: "r&amp;&quot;",
					},
					Cell: Attr{
						Class: "c&amp;&quot;",
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := SectionAttr{
				Section: test.fields.Section,
				Row:     test.fields.Row,
				Cell:    test.fields.Cell,
			}
			got := o.escape()
			testutil.AssertValue(t, got, test.want.val, "escape")
		})
	}
}
