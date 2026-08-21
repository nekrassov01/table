package markdown

import (
	"testing"

	"github.com/nekrassov01/table/internal/testutil"
)

func TestWithHeader(t *testing.T) {
	type fields struct {
		header []string
	}
	type args struct {
		header []string
	}
	type want struct {
		header []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "sets header",
			fields: fields{
				header: []string{"old"},
			},
			args: args{
				header: []string{"a", "b"},
			},
			want: want{
				header: []string{"a", "b"},
			},
		},
		{
			name: "clears header",
			fields: fields{
				header: []string{"old"},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &option{
				header: test.fields.header,
			}
			WithHeader(test.args.header)(o)
			got := want{
				header: o.header,
			}
			testutil.AssertValue(t, got, test.want, "WithHeader")
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
			name: "keeps one index",
			fields: fields{
				indexOffset: 1,
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
			WithPlaceholder(test.args.placeholder)(o)
			got := want{
				placeholder: o.placeholder,
			}
			testutil.AssertValue(t, got, test.want, "WithPlaceholder")
		})
	}
}

func TestWithAlign(t *testing.T) {
	type fields struct {
		columns columnSet
	}
	type args struct {
		columns ColumnSelector
		align   AlignSide
	}
	type want struct {
		columns []column
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "sets selected columns",
			args: args{
				columns: Columns(1),
				align:   AlignRight,
			},
			want: want{
				columns: []column{
					{},
					{
						align: AlignRight,
					},
				},
			},
		},
		{
			name: "clears alignment",
			fields: fields{
				columns: columnSet{
					values: []column{
						{
							align: AlignRight,
						},
					},
				},
			},
			args: args{
				columns: Columns(0),
			},
			want: want{
				columns: []column{{}},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &option{
				columns: test.fields.columns,
			}
			WithAlign(test.args.columns, test.args.align)(o)
			got := want{
				columns: o.columns.values,
			}
			testutil.AssertValue(t, got, test.want, "WithAlign")
		})
	}
}

func TestWithRowspan(t *testing.T) {
	type fields struct {
		columns columnSet
	}
	type args struct {
		columns ColumnSelector
	}
	type want struct {
		columns []column
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "enables selected columns",
			args: args{
				columns: Columns(0, 1),
			},
			want: want{
				columns: []column{
					{
						rowspan: true,
					},
					{
						rowspan: true,
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &option{
				columns: test.fields.columns,
			}
			WithRowspan(test.args.columns)(o)
			got := want{
				columns: o.columns.values,
			}
			testutil.AssertValue(t, got, test.want, "WithRowspan")
		})
	}
}

func TestWithColspan(t *testing.T) {
	type fields struct {
		columns columnSet
	}
	type args struct {
		columns ColumnSelector
	}
	type want struct {
		columns []column
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "enables selected columns",
			args: args{
				columns: Columns(1),
			},
			want: want{
				columns: []column{
					{},
					{
						colspan: true,
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &option{
				columns: test.fields.columns,
			}
			WithColspan(test.args.columns)(o)
			got := want{
				columns: o.columns.values,
			}
			testutil.AssertValue(t, got, test.want, "WithColspan")
		})
	}
}

func TestWithColor(t *testing.T) {
	type fields struct {
		columns columnSet
	}
	type args struct {
		scopes  Scope
		columns ColumnSelector
		color   *Color
	}
	type want struct {
		header *Color
		body   *Color
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "sets scoped color",
			args: args{
				scopes:  ScopeBody,
				columns: Columns(0),
				color:   ColorFgRed,
			},
			want: want{
				body: ColorFgRed,
			},
		},
		{
			name: "clears scoped color",
			fields: fields{
				columns: columnSet{
					values: func() []column {
						configured := column{}
						configured.transformer.colors.Set(ScopeHeader|ScopeBody, ColorFgBlue)
						return []column{configured}
					}(),
				},
			},
			args: args{
				scopes:  ScopeHeader,
				columns: Columns(0),
			},
			want: want{
				body: ColorFgBlue,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &option{
				columns: test.fields.columns,
			}
			WithColor(test.args.scopes, test.args.columns, test.args.color)(o)
			colors := o.columns.values[0].transformer.colors
			got := want{
				header: colors.Resolve(ScopeHeader),
				body:   colors.Resolve(ScopeBody),
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
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "sets scoped decoration",
			args: args{
				scopes:     ScopeHeader | ScopeBody,
				columns:    Columns(0),
				decoration: DecorationBold,
			},
			want: want{
				header: DecorationBold,
				body:   DecorationBold,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &option{}
			WithDecoration(test.args.scopes, test.args.columns, test.args.decoration)(o)
			decorations := o.columns.values[0].transformer.decorations
			got := want{
				header: decorations.Resolve(ScopeHeader),
				body:   decorations.Resolve(ScopeBody),
			}
			testutil.AssertValue(t, got, test.want, "WithDecoration")
		})
	}
}

func TestWithTransformer(t *testing.T) {
	type args struct {
		columns     ColumnSelector
		transformer func(any) (string, *Color, *Decoration)
		value       any
	}
	type want struct {
		text       string
		color      *Color
		decoration *Decoration
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "sets transformer",
			args: args{
				columns: Columns(0),
				transformer: func(any) (string, *Color, *Decoration) {
					return "transformed", ColorFgRed, DecorationBold
				},
				value: "raw",
			},
			want: want{
				text:       "transformed",
				color:      ColorFgRed,
				decoration: DecorationBold,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &option{}
			WithTransformer(test.args.columns, test.args.transformer)(o)
			text, color, decoration := o.columns.values[0].transformer.fn(test.args.value)
			got := want{
				text:       text,
				color:      color,
				decoration: decoration,
			}
			testutil.AssertValue(t, got, test.want, "WithTransformer")
		})
	}
}

func TestColumns(t *testing.T) {
	type args struct {
		indexes []int
	}
	type want struct {
		indexes []int
		all     bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "clones indexes",
			args: args{
				indexes: []int{0, 2},
			},
			want: want{
				indexes: []int{0, 2},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			selector := Columns(test.args.indexes...)
			if len(test.args.indexes) > 0 {
				test.args.indexes[0] = 9
			}
			got := want{
				indexes: selector.indexes,
				all:     selector.all,
			}
			testutil.AssertValue(t, got, test.want, "Columns")
		})
	}
}

func TestAllColumns(t *testing.T) {
	type want struct {
		indexes []int
		all     bool
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "selects every column",
			want: want{
				all: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			selector := AllColumns()
			got := want{
				indexes: selector.indexes,
				all:     selector.all,
			}
			testutil.AssertValue(t, got, test.want, "AllColumns")
		})
	}
}
