package text

import (
	"testing"

	"github.com/nekrassov01/table/internal/testutil"
)

func TestBorderStyle_maxGlyphLen(t *testing.T) {
	type fields struct {
		Top      *Horizontal
		Header   *Horizontal
		Body     *Horizontal
		Footer   *Horizontal
		Bottom   *Horizontal
		Vertical *Vertical
		Attr     *Attr
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
		},
		{
			name: "vertical",
			fields: fields{
				Vertical: &Vertical{
					Outer: "123456",
				},
			},
			want: want{
				val: 6,
			},
		},
		{
			name: "top",
			fields: fields{
				Top: &Horizontal{
					Fill: "1",
				},
			},
			want: want{
				val: 1,
			},
		},
		{
			name: "header",
			fields: fields{
				Header: &Horizontal{
					Fill: "12",
				},
			},
			want: want{
				val: 2,
			},
		},
		{
			name: "body",
			fields: fields{
				Body: &Horizontal{
					Fill: "123",
				},
			},
			want: want{
				val: 3,
			},
		},
		{
			name: "footer",
			fields: fields{
				Footer: &Horizontal{
					Fill: "1234",
				},
			},
			want: want{
				val: 4,
			},
		},
		{
			name: "bottom",
			fields: fields{
				Bottom: &Horizontal{
					Fill: "12345",
				},
			},
			want: want{
				val: 5,
			},
		},
		{
			name: "attribute excluded",
			fields: fields{
				Attr: &Attr{
					Prefix: []byte("attribute"),
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &BorderStyle{
				Top:      test.fields.Top,
				Header:   test.fields.Header,
				Body:     test.fields.Body,
				Footer:   test.fields.Footer,
				Bottom:   test.fields.Bottom,
				Vertical: test.fields.Vertical,
				Attr:     test.fields.Attr,
			}
			got := o.maxGlyphLen()
			testutil.AssertValue(t, got, test.want.val, "maxGlyphLen")
		})
	}
}

func TestContentStyle_maxAttrLen(t *testing.T) {
	type fields struct {
		Header  *Attr
		Body    *Attr
		Footer  *Attr
		Caption *Attr
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
		},
		{
			name: "header",
			fields: fields{
				Header: &Attr{
					Prefix: []byte("12"),
					Suffix: []byte("34"),
				},
			},
			want: want{
				val: 4,
			},
		},
		{
			name: "body",
			fields: fields{
				Body: &Attr{
					Prefix: []byte("123"),
					Suffix: []byte("45"),
				},
			},
			want: want{
				val: 5,
			},
		},
		{
			name: "footer",
			fields: fields{
				Footer: &Attr{
					Prefix: []byte("1234"),
					Suffix: []byte("56"),
				},
			},
			want: want{
				val: 6,
			},
		},
		{
			name: "caption excluded",
			fields: fields{
				Caption: &Attr{
					Prefix: []byte("caption"),
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &ContentStyle{
				Header:  test.fields.Header,
				Body:    test.fields.Body,
				Footer:  test.fields.Footer,
				Caption: test.fields.Caption,
			}
			got := o.maxAttrLen()
			testutil.AssertValue(t, got, test.want.val, "maxAttrLen")
		})
	}
}

func TestContentStyle_resolve(t *testing.T) {
	type fields struct {
		Header  *Attr
		Body    *Attr
		Footer  *Attr
		Caption *Attr
	}
	type args struct {
		sc Scope
	}
	type want struct {
		val *Attr
	}
	content := fields{
		Header:  NewAttr(CodeBold),
		Body:    NewAttr(CodeItalic),
		Footer:  NewAttr(CodeUnderline),
		Caption: NewAttr(CodeFaint),
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name:   "header",
			fields: content,
			args: args{
				sc: ScopeHeader,
			},
			want: want{
				val: NewAttr(CodeBold),
			},
		},
		{
			name:   "body",
			fields: content,
			args: args{
				sc: ScopeBody,
			},
			want: want{
				val: NewAttr(CodeItalic),
			},
		},
		{
			name:   "footer",
			fields: content,
			args: args{
				sc: ScopeFooter,
			},
			want: want{
				val: NewAttr(CodeUnderline),
			},
		},
		{
			name:   "no scope",
			fields: content,
		},
		{
			name:   "combined scope",
			fields: content,
			args: args{
				sc: ScopeHeader | ScopeBody,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &ContentStyle{
				Header:  test.fields.Header,
				Body:    test.fields.Body,
				Footer:  test.fields.Footer,
				Caption: test.fields.Caption,
			}
			got := o.resolve(test.args.sc)
			testutil.AssertValue(t, got, test.want.val, "resolve")
		})
	}
}

func Test_colored(t *testing.T) {
	type args struct {
		s Style
	}
	type want struct {
		style  Style
		source Style
	}
	borderAttr := NewAttr(CodeFaint)
	headerAttr := NewAttr(CodeBold)
	captionAttr := NewAttr(CodeFaint)
	coloredContent := ContentStyle{
		Header:  headerAttr,
		Footer:  headerAttr,
		Caption: captionAttr,
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "empty style",
			want: want{
				style: Style{
					Border: BorderStyle{
						Attr: borderAttr,
					},
					Content: coloredContent,
				},
			},
		},
		{
			name: "preserves glyphs",
			args: args{
				s: Style{
					Border: BorderStyle{
						Top: &Horizontal{
							Fill: "-",
						},
						Vertical: &Vertical{
							Outer: "|",
							Inner: "|",
						},
						Attr: NewAttr(CodeBlinkSlow),
					},
					Content: ContentStyle{
						Body: NewAttr(CodeItalic),
					},
				},
			},
			want: want{
				style: func() Style {
					style := Style{
						Border: BorderStyle{
							Top: &Horizontal{
								Fill: "-",
							},
							Vertical: &Vertical{
								Outer: "|",
								Inner: "|",
							},
							Attr: NewAttr(CodeBlinkSlow),
						},
						Content: ContentStyle{
							Body: NewAttr(CodeItalic),
						},
					}
					return Style{
						Border: BorderStyle{
							Top:      style.Border.Top,
							Vertical: style.Border.Vertical,
							Attr:     borderAttr,
						},
						Content: coloredContent,
					}
				}(),
				source: Style{
					Border: BorderStyle{
						Top: &Horizontal{
							Fill: "-",
						},
						Vertical: &Vertical{
							Outer: "|",
							Inner: "|",
						},
						Attr: NewAttr(CodeBlinkSlow),
					},
					Content: ContentStyle{
						Body: NewAttr(CodeItalic),
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			source := test.args.s
			got := want{
				style:  colored(source),
				source: source,
			}
			testutil.AssertValue(t, got, test.want, "colored")
		})
	}
}
