package text

import (
	"testing"

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
				rows: [][]string{
					{"top"},
					{"bottom"},
				},
			},
			want: want{
				header: [][]string{
					{"top"},
					{"bottom"},
				},
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
		s    string
		side CaptionSide
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
				s:    "title",
				side: CaptionTop,
			},
			want: want{
				caption:     "title",
				captionSide: CaptionTop,
			},
		},
		{
			name: "clears caption",
			fields: fields{
				caption:     "old",
				captionSide: CaptionTop,
			},
			args: args{
				side: CaptionBottom,
			},
			want: want{
				captionSide: CaptionBottom,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &option{
				caption:     test.fields.caption,
				captionSide: test.fields.captionSide,
			}
			WithCaption(test.args.s, test.args.side)(o)
			got := want{
				caption:     o.caption,
				captionSide: o.captionSide,
			}
			testutil.AssertValue(t, got, test.want, "WithCaption")
		})
	}
}

func TestWithStyle(t *testing.T) {
	type fields struct {
		style Style
	}
	type args struct {
		style Style
	}
	type want struct {
		style Style
	}
	attr := NewAttr(CodeBold)
	style := Style{
		Border: BorderStyle{
			Top: &Horizontal{
				Fill: "-",
			},
		},
		Content: ContentStyle{
			Body: attr,
		},
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "sets style",
			fields: fields{
				style: StyleLight,
			},
			args: args{
				style: style,
			},
			want: want{
				style: style,
			},
		},
		{
			name: "clears style",
			fields: fields{
				style: style,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &option{
				style: test.fields.style,
			}
			WithStyle(test.args.style)(o)
			got := want{
				style: o.style,
			}
			testutil.AssertValue(t, got, test.want, "WithStyle")
		})
	}
}

func TestWithCompact(t *testing.T) {
	type fields struct {
		compact bool
	}
	type want struct {
		compact bool
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "enables compact",
			want: want{
				compact: true,
			},
		},
		{
			name: "keeps compact enabled",
			fields: fields{
				compact: true,
			},
			want: want{
				compact: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &option{
				compact: test.fields.compact,
			}
			WithCompact()(o)
			got := want{
				compact: o.compact,
			}
			testutil.AssertValue(t, got, test.want, "WithCompact")
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
			want: want{
				indexOffset: 1,
			},
		},
		{
			name: "normalizes offset",
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

func TestWithIndexWidth(t *testing.T) {
	type fields struct {
		indexOffset int
		indexWidth  int
	}
	type args struct {
		n int
	}
	type want struct {
		indexOffset int
		indexWidth  int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "sets width",
			fields: fields{
				indexWidth: 2,
			},
			args: args{
				n: 6,
			},
			want: want{
				indexOffset: 1,
				indexWidth:  6,
			},
		},
		{
			name: "zero keeps width",
			fields: fields{
				indexWidth: 3,
			},
			want: want{
				indexOffset: 1,
				indexWidth:  3,
			},
		},
		{
			name: "negative keeps width",
			fields: fields{
				indexWidth: 4,
			},
			args: args{
				n: -1,
			},
			want: want{
				indexOffset: 1,
				indexWidth:  4,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &option{
				indexOffset: test.fields.indexOffset,
				indexWidth:  test.fields.indexWidth,
			}
			WithIndexWidth(test.args.n)(o)
			got := want{
				indexOffset: o.indexOffset,
				indexWidth:  o.indexWidth,
			}
			testutil.AssertValue(t, got, test.want, "WithIndexWidth")
		})
	}
}

func TestWithAutoFit(t *testing.T) {
	type fields struct {
		autoFit bool
	}
	type want struct {
		autoFit bool
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "enables auto fit",
			want: want{
				autoFit: true,
			},
		},
		{
			name: "keeps auto fit enabled",
			fields: fields{
				autoFit: true,
			},
			want: want{
				autoFit: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &option{
				autoFit: test.fields.autoFit,
			}
			WithAutoFit()(o)
			got := want{
				autoFit: o.autoFit,
			}
			testutil.AssertValue(t, got, test.want, "WithAutoFit")
		})
	}
}

func TestWithPlaceholder(t *testing.T) {
	type fields struct {
		placeholder string
	}
	type args struct {
		s string
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
				s: "-",
			},
			want: want{
				placeholder: "-",
			},
		},
		{
			name: "clears placeholder",
			fields: fields{
				placeholder: "old",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &option{
				placeholder: test.fields.placeholder,
			}
			WithPlaceholder(test.args.s)(o)
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
			column := o.columns.values[1]
			got := want{
				columns: len(o.columns.values),
				header:  column.aligns.Resolve(ScopeHeader),
				body:    column.aligns.Resolve(ScopeBody),
				footer:  column.aligns.Resolve(ScopeFooter),
			}
			testutil.AssertValue(t, got, test.want, "WithAlign")
		})
	}
}

func TestWithWidth(t *testing.T) {
	type fields struct {
		limit int
	}
	type args struct {
		columns ColumnSelector
		width   int
	}
	type want struct {
		limit int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "sets width",
			args: args{
				columns: Columns(0),
				width:   8,
			},
			want: want{
				limit: 8,
			},
		},
		{
			name: "zero clears width",
			fields: fields{
				limit: 8,
			},
			args: args{
				columns: Columns(0),
			},
		},
		{
			name: "negative clears width",
			fields: fields{
				limit: 8,
			},
			args: args{
				columns: Columns(0),
				width:   -1,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &option{
				columns: columnSet{
					values: []column{
						{
							limit: test.fields.limit,
						},
					},
				},
			}
			WithWidth(test.args.columns, test.args.width)(o)
			got := want{
				limit: o.columns.values[0].limit,
			}
			testutil.AssertValue(t, got, test.want, "WithWidth")
		})
	}
}

func TestWithTruncate(t *testing.T) {
	type fields struct {
		truncate bool
	}
	type args struct {
		columns ColumnSelector
	}
	type want struct {
		truncate bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "enables truncation",
			args: args{
				columns: Columns(0),
			},
			want: want{
				truncate: true,
			},
		},
		{
			name: "keeps truncation enabled",
			fields: fields{
				truncate: true,
			},
			args: args{
				columns: Columns(0),
			},
			want: want{
				truncate: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &option{
				columns: columnSet{
					values: []column{
						{
							truncate: test.fields.truncate,
						},
					},
				},
			}
			WithTruncate(test.args.columns)(o)
			got := want{
				truncate: o.columns.values[0].truncate,
			}
			testutil.AssertValue(t, got, test.want, "WithTruncate")
		})
	}
}

func TestWithPadding(t *testing.T) {
	type fields struct {
		lPad int
		rPad int
	}
	type args struct {
		columns ColumnSelector
		left    int
		right   int
	}
	type want struct {
		lPad int
		rPad int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "sets padding",
			args: args{
				columns: Columns(0),
				left:    2,
				right:   3,
			},
			want: want{
				lPad: 2,
				rPad: 3,
			},
		},
		{
			name: "clamps negative padding",
			fields: fields{
				lPad: 4,
				rPad: 5,
			},
			args: args{
				columns: Columns(0),
				left:    -1,
				right:   -2,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &option{
				columns: columnSet{
					values: []column{
						{
							lPad: test.fields.lPad,
							rPad: test.fields.rPad,
						},
					},
				},
			}
			WithPadding(test.args.columns, test.args.left, test.args.right)(o)
			got := want{
				lPad: o.columns.values[0].lPad,
				rPad: o.columns.values[0].rPad,
			}
			testutil.AssertValue(t, got, test.want, "WithPadding")
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
					values: []column{
						{
							rowspan: test.fields.rowspan,
						},
					},
				},
			}
			WithRowspan(test.args.scopes, test.args.columns)(o)
			got := want{
				rowspan: o.columns.values[0].rowspan,
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
					values: []column{
						{
							colspan: test.fields.colspan,
						},
					},
				},
			}
			WithColspan(test.args.scopes, test.args.columns)(o)
			got := want{
				colspan: o.columns.values[0].colspan,
			}
			testutil.AssertValue(t, got, test.want, "WithColspan")
		})
	}
}

func TestWithAttr(t *testing.T) {
	type fields struct {
		header *Attr
		body   *Attr
		footer *Attr
	}
	type args struct {
		scopes  Scope
		columns ColumnSelector
		attr    *Attr
	}
	type want struct {
		header *Attr
		body   *Attr
		footer *Attr
	}
	oldAttr := NewAttr(CodeFaint)
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "sets selected scopes",
			fields: fields{
				body: oldAttr,
			},
			args: args{
				scopes:  ScopeHeader | ScopeFooter,
				columns: Columns(0),
				attr:    NewAttr(CodeBold),
			},
			want: want{
				header: NewAttr(CodeBold),
				body:   oldAttr,
				footer: NewAttr(CodeBold),
			},
		},
		{
			name: "clears selected scope",
			fields: fields{
				body: oldAttr,
			},
			args: args{
				scopes:  ScopeBody,
				columns: Columns(0),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &option{
				columns: columnSet{
					values: []column{{}},
				},
			}
			attrs := &o.columns.values[0].transformer.attrs
			attrs.Set(ScopeHeader, test.fields.header)
			attrs.Set(ScopeBody, test.fields.body)
			attrs.Set(ScopeFooter, test.fields.footer)
			WithAttr(test.args.scopes, test.args.columns, test.args.attr)(o)
			got := want{
				header: attrs.Resolve(ScopeHeader),
				body:   attrs.Resolve(ScopeBody),
				footer: attrs.Resolve(ScopeFooter),
			}
			testutil.AssertValue(t, got, test.want, "WithAttr")
		})
	}
}

func TestWithTransformer(t *testing.T) {
	type fields struct {
		fn func(any) (string, *Attr)
	}
	type args struct {
		columns ColumnSelector
		fn      func(any) (string, *Attr)
	}
	type want struct {
		isNil bool
		value string
		attr  *Attr
	}
	oldTransformer := func(any) (string, *Attr) {
		return "old", nil
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
				fn: func() func(any) (string, *Attr) {
					attr := NewAttr(CodeBold)
					return func(any) (string, *Attr) {
						return "new", attr
					}
				}(),
			},
			want: want{
				value: "new",
				attr:  NewAttr(CodeBold),
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
					values: []column{
						{
							transformer: transformer{
								fn: test.fields.fn,
							},
						},
					},
				},
			}
			WithTransformer(test.args.columns, test.args.fn)(o)
			gotFn := o.columns.values[0].transformer.fn
			got := want{
				isNil: gotFn == nil,
			}
			if gotFn != nil {
				got.value, got.attr = gotFn("input")
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
					indexes: []int{0, -1, 2},
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
					all: true,
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
