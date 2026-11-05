package markdown

import (
	"testing"

	"github.com/nekrassov01/table"
	"github.com/nekrassov01/table/internal/column"
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

func TestWithAlign(t *testing.T) {
	type fields struct {
		columns columnSet
	}
	type args struct {
		columns ColumnSelector
		align   AlignSide
	}
	type want struct {
		columns []columnConfig
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
				columns: []columnConfig{
					{placeholder: placeholder},
					{
						align:       AlignRight,
						placeholder: placeholder,
					},
				},
			},
		},
		{
			name: "clears alignment",
			fields: fields{
				columns: columnSet{
					Values: []columnConfig{
						{
							align:       AlignRight,
							placeholder: placeholder,
						},
					},
				},
			},
			args: args{
				columns: Columns(0),
			},
			want: want{
				columns: []columnConfig{{placeholder: placeholder}},
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
				columns: o.columns.resolve(nil, len(test.want.columns)),
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
		columns []columnConfig
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
				columns: []columnConfig{
					{
						rowspan:     true,
						placeholder: placeholder,
					},
					{
						rowspan:     true,
						placeholder: placeholder,
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
				columns: o.columns.resolve(nil, len(test.want.columns)),
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
		columns []columnConfig
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
				columns: []columnConfig{
					{placeholder: placeholder},
					{
						colspan:     true,
						placeholder: placeholder,
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
				columns: o.columns.resolve(nil, len(test.want.columns)),
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
					Values: func() []columnConfig {
						configured := columnConfig{placeholder: placeholder}
						configured.transformer.colors.Set(ScopeHeader|ScopeBody, ColorFgBlue)
						return []columnConfig{configured}
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
			colors := o.columns.Values[0].transformer.colors
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
			decorations := o.columns.Values[0].transformer.decorations
			got := want{
				header: decorations.Resolve(ScopeHeader),
				body:   decorations.Resolve(ScopeBody),
			}
			testutil.AssertValue(t, got, test.want, "WithDecoration")
		})
	}
}

func TestWithPlaceholder(t *testing.T) {
	type fields struct {
		columns columnSet
	}
	type args struct {
		columns ColumnSelector
		s       string
	}
	type want struct {
		placeholders []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "selected column",
			args: args{
				columns: Columns(1),
				s:       "-",
			},
			want: want{
				placeholders: []string{placeholder, "-", placeholder},
			},
		},
		{
			name: "all columns",
			args: args{
				columns: AllColumns(),
				s:       "-",
			},
			want: want{
				placeholders: []string{"-", "-", "-"},
			},
		},
		{
			name: "explicit empty",
			fields: fields{
				columns: columnSetOf(nil, &columnConfig{placeholder: "-"}),
			},
			args: args{
				columns: Columns(1),
				s:       "",
			},
			want: want{
				placeholders: []string{"-", "", "-"},
			},
		},
		{
			name: "all overrides selected",
			fields: fields{
				columns: columnSet{
					Values: []columnConfig{{}, {placeholder: "old"}},
				},
			},
			args: args{
				columns: AllColumns(),
				s:       "new",
			},
			want: want{
				placeholders: []string{"new", "new", "new"},
			},
		},
		{
			name: "empty selector",
			args: args{
				columns: Columns(),
				s:       "-",
			},
			want: want{
				placeholders: []string{placeholder, placeholder, placeholder},
			},
		},
		{
			name: "outside resolved columns",
			args: args{
				columns: Columns(1000000),
				s:       "-",
			},
			want: want{
				placeholders: []string{placeholder, placeholder, placeholder},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &option{
				columns: test.fields.columns,
			}
			WithPlaceholder(test.args.columns, test.args.s)(o)
			columns := o.columns.resolve(nil, 3)
			got := want{
				placeholders: make([]string, len(columns)),
			}
			for index := range columns {
				got.placeholders[index] = columns[index].placeholder
			}
			testutil.AssertValue(t, got, test.want, "WithPlaceholder")
		})
	}
}

func TestWithTransformer(t *testing.T) {
	type args struct {
		columns     ColumnSelector
		transformer func(table.Value) (string, *Color, *Decoration)
		value       table.Value
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
				transformer: func(table.Value) (string, *Color, *Decoration) {
					return "transformed", ColorFgRed, DecorationBold
				},
				value: table.String("raw"),
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
			text, color, decoration := o.columns.Values[0].transformer.fn(test.args.value)
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
		selector ColumnSelector
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
				selector: ColumnSelector{
					selector: column.NewSelector(0, 2),
				},
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
				selector: selector,
			}
			testutil.AssertValue(t, got, test.want, "Columns")
		})
	}
}

func TestAllColumns(t *testing.T) {
	type want struct {
		selector ColumnSelector
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "selects every column",
			want: want{
				selector: ColumnSelector{
					selector: column.All(),
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			selector := AllColumns()
			got := want{
				selector: selector,
			}
			testutil.AssertValue(t, got, test.want, "AllColumns")
		})
	}
}
