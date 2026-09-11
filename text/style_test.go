package text

import (
	"testing"

	"github.com/nekrassov01/table/internal/testutil"
)

func TestStyle_Clone(t *testing.T) {
	type fields struct {
		style Style
	}
	type want struct {
		style Style
	}
	sharedHorizontal := &Horizontal{Fill: "base"}
	sharedAttr := &Attr{
		Prefix: []byte("prefix"),
		Suffix: []byte("suffix"),
	}
	fullStyle := Style{
		Border: BorderStyle{
			Top:      sharedHorizontal,
			Header:   sharedHorizontal,
			Body:     sharedHorizontal,
			Footer:   sharedHorizontal,
			Bottom:   sharedHorizontal,
			Vertical: &Vertical{Outer: "outer", Inner: "inner"},
			Attr:     sharedAttr,
		},
		Content: ContentStyle{
			Header:  sharedAttr,
			Body:    sharedAttr,
			Footer:  sharedAttr,
			Caption: sharedAttr,
		},
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
			name: "nested values are independent",
			fields: fields{
				style: fullStyle,
			},
			want: want{
				style: fullStyle,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.fields.style.Clone()
			testutil.AssertValue(t, got, test.want.style, "Clone")
			horizontalPairs := [...]struct {
				source *Horizontal
				clone  *Horizontal
			}{
				{test.fields.style.Border.Top, got.Border.Top},
				{test.fields.style.Border.Header, got.Border.Header},
				{test.fields.style.Border.Body, got.Border.Body},
				{test.fields.style.Border.Footer, got.Border.Footer},
				{test.fields.style.Border.Bottom, got.Border.Bottom},
			}
			seenHorizontals := make(map[*Horizontal]struct{}, len(horizontalPairs))
			for i, pair := range horizontalPairs {
				if pair.clone != nil && pair.clone == pair.source {
					t.Errorf("horizontal %d shares the source", i)
				}
				if _, exists := seenHorizontals[pair.clone]; pair.clone != nil && exists {
					t.Errorf("horizontal %d shares another clone", i)
				}
				if pair.clone != nil {
					seenHorizontals[pair.clone] = struct{}{}
				}
			}
			if got.Border.Vertical != nil && got.Border.Vertical == test.fields.style.Border.Vertical {
				t.Error("vertical shares the source")
			}
			attrPairs := [...]struct {
				source *Attr
				clone  *Attr
			}{
				{test.fields.style.Border.Attr, got.Border.Attr},
				{test.fields.style.Content.Header, got.Content.Header},
				{test.fields.style.Content.Body, got.Content.Body},
				{test.fields.style.Content.Footer, got.Content.Footer},
				{test.fields.style.Content.Caption, got.Content.Caption},
			}
			seenAttrs := make(map[*Attr]struct{}, len(attrPairs))
			for i, pair := range attrPairs {
				if pair.clone == nil {
					continue
				}
				if pair.clone == pair.source {
					t.Errorf("attribute %d shares the source", i)
				}
				if len(pair.clone.Prefix) > 0 && len(pair.source.Prefix) > 0 &&
					&pair.clone.Prefix[0] == &pair.source.Prefix[0] {
					t.Errorf("attribute %d shares its prefix", i)
				}
				if len(pair.clone.Suffix) > 0 && len(pair.source.Suffix) > 0 &&
					&pair.clone.Suffix[0] == &pair.source.Suffix[0] {
					t.Errorf("attribute %d shares its suffix", i)
				}
				if _, exists := seenAttrs[pair.clone]; exists {
					t.Errorf("attribute %d shares another clone", i)
				}
				seenAttrs[pair.clone] = struct{}{}
			}
		})
	}
}

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

func Test_clonePointer(t *testing.T) {
	type args struct {
		value *Horizontal
	}
	type want struct {
		value  *Horizontal
		shared bool
	}
	value := &Horizontal{Fill: "-"}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "nil",
		},
		{
			name: "value",
			args: args{
				value: value,
			},
			want: want{
				value: value,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value := clonePointer(test.args.value)
			got := want{
				value:  value,
				shared: value != nil && value == test.args.value,
			}
			testutil.AssertValue(t, got, test.want, "clonePointer")
		})
	}
}

func Test_cloneAttr(t *testing.T) {
	type args struct {
		attr *Attr
	}
	type want struct {
		attr         *Attr
		shared       bool
		prefixShared bool
		suffixShared bool
	}
	empty := &Attr{
		Prefix: []byte{},
		Suffix: []byte{},
	}
	attr := &Attr{
		Prefix: []byte("prefix"),
		Suffix: []byte("suffix"),
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "nil",
		},
		{
			name: "empty",
			args: args{
				attr: empty,
			},
			want: want{
				attr: empty,
			},
		},
		{
			name: "byte slices",
			args: args{
				attr: attr,
			},
			want: want{
				attr: attr,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			attr := cloneAttr(test.args.attr)
			got := want{
				attr:   attr,
				shared: attr != nil && attr == test.args.attr,
			}
			if attr != nil && test.args.attr != nil {
				got.prefixShared = len(attr.Prefix) > 0 && len(test.args.attr.Prefix) > 0 &&
					&attr.Prefix[0] == &test.args.attr.Prefix[0]
				got.suffixShared = len(attr.Suffix) > 0 && len(test.args.attr.Suffix) > 0 &&
					&attr.Suffix[0] == &test.args.attr.Suffix[0]
			}
			testutil.AssertValue(t, got, test.want, "cloneAttr")
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
