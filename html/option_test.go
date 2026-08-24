package html

import (
	"testing"

	"github.com/nekrassov01/table/internal/column"
	"github.com/nekrassov01/table/internal/testutil"
)

func TestWithHeader(t *testing.T) {
	type fields struct {
		header [][]string
	}
	type args struct {
		rows [][]string
	}
	type want struct {
		header [][]string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "sets rows",
			fields: fields{
				header: [][]string{{"old"}},
			},
			args: args{
				rows: [][]string{{"top"}, {"bottom"}},
			},
			want: want{
				header: [][]string{{"top"}, {"bottom"}},
			},
		},
		{
			name: "clears rows",
			fields: fields{
				header: [][]string{{"old"}},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &option{
				header: test.fields.header,
			}
			WithHeader(test.args.rows...)(o)
			got := want{
				header: o.header,
			}
			testutil.AssertValue(t, got, test.want, "WithHeader")
		})
	}
}

func TestWithFooter(t *testing.T) {
	type fields struct {
		footer func() [][]string
	}
	type args struct {
		fn func() [][]string
	}
	type want struct {
		isNil bool
		rows  [][]string
	}
	oldFooter := func() [][]string {
		return [][]string{{"old"}}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "sets function",
			fields: fields{
				footer: oldFooter,
			},
			args: args{
				fn: func() [][]string {
					return [][]string{{"total"}}
				},
			},
			want: want{
				rows: [][]string{{"total"}},
			},
		},
		{
			name: "clears function",
			fields: fields{
				footer: oldFooter,
			},
			want: want{
				isNil: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &option{
				footer: test.fields.footer,
			}
			WithFooter(test.args.fn)(o)
			got := want{
				isNil: o.footer == nil,
			}
			if o.footer != nil {
				got.rows = o.footer()
			}
			testutil.AssertValue(t, got, test.want, "WithFooter")
		})
	}
}

func TestWithCaption(t *testing.T) {
	type fields struct {
		caption     string
		captionSide CaptionSide
	}
	type args struct {
		caption string
		side    CaptionSide
	}
	type want struct {
		caption     string
		captionSide CaptionSide
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "sets caption",
			fields: fields{
				caption:     "old",
				captionSide: CaptionBottom,
			},
			args: args{
				caption: "title",
				side:    CaptionTop,
			},
			want: want{
				caption:     "title",
				captionSide: CaptionTop,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &option{
				caption:     test.fields.caption,
				captionSide: test.fields.captionSide,
			}
			WithCaption(test.args.caption, test.args.side)(o)
			got := want{
				caption:     o.caption,
				captionSide: o.captionSide,
			}
			testutil.AssertValue(t, got, test.want, "WithCaption")
		})
	}
}

func TestWithTableAttr(t *testing.T) {
	type args struct {
		attr TableAttr
	}
	type want struct {
		attr TableAttr
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "sets escaped attributes",
			args: args{
				attr: TableAttr{
					Table: Attr{
						Class: `a&b`,
					},
					Body: SectionAttr{
						Cell: Attr{
							Style: `content:"x"`,
						},
					},
				},
			},
			want: want{
				attr: TableAttr{
					Table: Attr{
						Class: `a&amp;b`,
					},
					Body: SectionAttr{
						Cell: Attr{
							Style: `content:&quot;x&quot;`,
						},
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &option{}
			WithTableAttr(test.args.attr)(o)
			got := want{
				attr: o.attrs,
			}
			testutil.AssertValue(t, got, test.want, "WithTableAttr")
		})
	}
}

func TestWithIndex(t *testing.T) {
	type fields struct {
		indexOffset int
	}
	type want struct {
		indexOffset int
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "enables index",
			fields: fields{
				indexOffset: 2,
			},
			want: want{
				indexOffset: 1,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &option{
				indexOffset: test.fields.indexOffset,
			}
			WithIndex()(o)
			got := want{
				indexOffset: o.indexOffset,
			}
			testutil.AssertValue(t, got, test.want, "WithIndex")
		})
	}
}

func TestWithPlaceholder(t *testing.T) {
	type fields struct {
		placeholder string
	}
	type args struct {
		placeholder string
	}
	type want struct {
		placeholder string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "sets placeholder",
			fields: fields{
				placeholder: "old",
			},
			args: args{
				placeholder: "-",
			},
			want: want{
				placeholder: "-",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &option{
				placeholder: test.fields.placeholder,
			}
			WithPlaceholder(test.args.placeholder)(o)
			got := want{
				placeholder: o.placeholder,
			}
			testutil.AssertValue(t, got, test.want, "WithPlaceholder")
		})
	}
}

func TestWithAlign(t *testing.T) {
	type args struct {
		scopes  Scope
		columns ColumnSelector
		align   AlignSide
	}
	type want struct {
		columns int
		header  AlignSide
		body    AlignSide
		footer  AlignSide
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "sets selected scopes",
			args: args{
				scopes:  ScopeHeader | ScopeFooter,
				columns: Columns(1),
				align:   AlignRight,
			},
			want: want{
				columns: 2,
				header:  AlignRight,
				footer:  AlignRight,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &option{}
			WithAlign(test.args.scopes, test.args.columns, test.args.align)(o)
			column := o.columns.Values[1]
			got := want{
				columns: len(o.columns.Values),
				header:  column.aligns.Resolve(ScopeHeader),
				body:    column.aligns.Resolve(ScopeBody),
				footer:  column.aligns.Resolve(ScopeFooter),
			}
			testutil.AssertValue(t, got, test.want, "WithAlign")
		})
	}
}

func TestWithRowspan(t *testing.T) {
	type fields struct {
		rowspan Scope
	}
	type args struct {
		scopes  Scope
		columns ColumnSelector
	}
	type want struct {
		rowspan Scope
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "adds scopes",
			fields: fields{
				rowspan: ScopeHeader,
			},
			args: args{
				scopes:  ScopeBody | ScopeFooter,
				columns: Columns(0),
			},
			want: want{
				rowspan: ScopeHeader | ScopeBody | ScopeFooter,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &option{
				columns: columnSet{
					Values: []columnConfig{
						{
							rowspan: test.fields.rowspan,
						},
					},
				},
			}
			WithRowspan(test.args.scopes, test.args.columns)(o)
			got := want{
				rowspan: o.columns.Values[0].rowspan,
			}
			testutil.AssertValue(t, got, test.want, "WithRowspan")
		})
	}
}

func TestWithColspan(t *testing.T) {
	type fields struct {
		colspan Scope
	}
	type args struct {
		scopes  Scope
		columns ColumnSelector
	}
	type want struct {
		colspan Scope
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "adds scopes",
			fields: fields{
				colspan: ScopeHeader,
			},
			args: args{
				scopes:  ScopeBody | ScopeFooter,
				columns: Columns(0),
			},
			want: want{
				colspan: ScopeHeader | ScopeBody | ScopeFooter,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &option{
				columns: columnSet{
					Values: []columnConfig{
						{
							colspan: test.fields.colspan,
						},
					},
				},
			}
			WithColspan(test.args.scopes, test.args.columns)(o)
			got := want{
				colspan: o.columns.Values[0].colspan,
			}
			testutil.AssertValue(t, got, test.want, "WithColspan")
		})
	}
}

func TestWithColor(t *testing.T) {
	type args struct {
		scopes  Scope
		columns ColumnSelector
		color   *Color
	}
	type want struct {
		header *Color
		body   *Color
		footer *Color
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "sets selected scopes",
			args: args{
				scopes:  ScopeHeader | ScopeFooter,
				columns: Columns(0),
				color:   ColorFgRed,
			},
			want: want{
				header: ColorFgRed,
				footer: ColorFgRed,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &option{}
			WithColor(test.args.scopes, test.args.columns, test.args.color)(o)
			colors := &o.columns.Values[0].transformer.colors
			got := want{
				header: colors.Resolve(ScopeHeader),
				body:   colors.Resolve(ScopeBody),
				footer: colors.Resolve(ScopeFooter),
			}
			testutil.AssertValue(t, got, test.want, "WithColor")
		})
	}
}

func TestWithDecoration(t *testing.T) {
	type args struct {
		scopes     Scope
		columns    ColumnSelector
		decoration *Decoration
	}
	type want struct {
		header *Decoration
		body   *Decoration
		footer *Decoration
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "sets selected scopes",
			args: args{
				scopes:     ScopeHeader | ScopeFooter,
				columns:    Columns(0),
				decoration: DecorationBold,
			},
			want: want{
				header: DecorationBold,
				footer: DecorationBold,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &option{}
			WithDecoration(test.args.scopes, test.args.columns, test.args.decoration)(o)
			decorations := &o.columns.Values[0].transformer.decorations
			got := want{
				header: decorations.Resolve(ScopeHeader),
				body:   decorations.Resolve(ScopeBody),
				footer: decorations.Resolve(ScopeFooter),
			}
			testutil.AssertValue(t, got, test.want, "WithDecoration")
		})
	}
}

func TestWithCellAttr(t *testing.T) {
	type args struct {
		scopes  Scope
		columns ColumnSelector
		attr    Attr
	}
	type want struct {
		header Attr
		body   Attr
		footer Attr
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "sets escaped attributes in selected scopes",
			args: args{
				scopes:  ScopeHeader | ScopeFooter,
				columns: Columns(0),
				attr: Attr{
					Class: `a&b`,
					Style: `content:"x"`,
				},
			},
			want: want{
				header: Attr{
					Class: `a&amp;b`,
					Style: `content:&quot;x&quot;`,
				},
				footer: Attr{
					Class: `a&amp;b`,
					Style: `content:&quot;x&quot;`,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &option{}
			WithCellAttr(test.args.scopes, test.args.columns, test.args.attr)(o)
			attrs := &o.columns.Values[0].attrs
			got := want{
				header: attrs.Resolve(ScopeHeader),
				body:   attrs.Resolve(ScopeBody),
				footer: attrs.Resolve(ScopeFooter),
			}
			testutil.AssertValue(t, got, test.want, "WithCellAttr")
		})
	}
}

func TestWithTransformer(t *testing.T) {
	type fields struct {
		fn func(any) (string, *Color, *Decoration)
	}
	type args struct {
		columns ColumnSelector
		fn      func(any) (string, *Color, *Decoration)
	}
	type want struct {
		isNil      bool
		value      string
		color      *Color
		decoration *Decoration
	}
	oldTransformer := func(any) (string, *Color, *Decoration) {
		return "old", nil, nil
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "sets function",
			fields: fields{
				fn: oldTransformer,
			},
			args: args{
				columns: Columns(0),
				fn: func(any) (string, *Color, *Decoration) {
					return "new", ColorFgRed, DecorationBold
				},
			},
			want: want{
				value:      "new",
				color:      ColorFgRed,
				decoration: DecorationBold,
			},
		},
		{
			name: "clears function",
			fields: fields{
				fn: oldTransformer,
			},
			args: args{
				columns: Columns(0),
			},
			want: want{
				isNil: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &option{
				columns: columnSet{
					Values: []columnConfig{
						{
							transformer: transformer{
								fn: test.fields.fn,
							},
						},
					},
				},
			}
			WithTransformer(test.args.columns, test.args.fn)(o)
			gotFn := o.columns.Values[0].transformer.fn
			got := want{
				isNil: gotFn == nil,
			}
			if gotFn != nil {
				got.value, got.color, got.decoration = gotFn("input")
			}
			testutil.AssertValue(t, got, test.want, "WithTransformer")
		})
	}
}

func TestColumns(t *testing.T) {
	type args struct {
		indexes []int
		mutate  []int
	}
	type want struct {
		val ColumnSelector
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "no indexes",
			want: want{
				val: ColumnSelector{},
			},
		},
		{
			name: "clones indexes",
			args: args{
				indexes: []int{0, -1, 2},
				mutate:  []int{3, 4, 5},
			},
			want: want{
				val: ColumnSelector{
					selector: column.NewSelector(0, -1, 2),
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Columns(test.args.indexes...)
			copy(test.args.indexes, test.args.mutate)
			testutil.AssertValue(t, got, test.want.val, "Columns")
		})
	}
}

func TestAllColumns(t *testing.T) {
	type want struct {
		val ColumnSelector
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "all columns",
			want: want{
				val: ColumnSelector{
					selector: column.All(),
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := AllColumns()
			testutil.AssertValue(t, got, test.want.val, "AllColumns")
		})
	}
}
