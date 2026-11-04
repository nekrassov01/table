package csv

import (
	"testing"

	"github.com/nekrassov01/table"
	"github.com/nekrassov01/table/internal/column"
	"github.com/nekrassov01/table/internal/testutil"
)

func TestWithHeader(t *testing.T) {
	type args struct {
		header []string
	}
	type want struct {
		header []string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "sets header",
			args: args{
				header: []string{"A", "B"},
			},
			want: want{
				header: []string{"A", "B"},
			},
		},
		{
			name: "clears header",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &option{
				header: []string{"old"},
			}
			WithHeader(test.args.header)(o)
			got := want{
				header: o.header,
			}
			testutil.AssertValue(t, got, test.want, "WithHeader")
		})
	}
}

func TestWithFooter(t *testing.T) {
	type args struct {
		fn func() [][]string
	}
	type want struct {
		rows   [][]string
		called int
		isNil  bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "retains function",
			args: args{
				fn: func() [][]string {
					return [][]string{{"sum", "3"}}
				},
			},
			want: want{
				rows:   [][]string{{"sum", "3"}},
				called: 1,
			},
		},
		{
			name: "clears function",
			want: want{
				isNil: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			called := 0
			fn := test.args.fn
			if fn != nil {
				fn = func() [][]string {
					called++
					return test.args.fn()
				}
			}
			o := &option{
				footer: func() [][]string {
					return [][]string{{"old"}}
				},
			}
			WithFooter(fn)(o)
			var rows [][]string
			if o.footer != nil {
				rows = o.footer()
			}
			got := want{
				rows:   rows,
				called: called,
				isNil:  o.footer == nil,
			}
			testutil.AssertValue(t, got, test.want, "WithFooter")
		})
	}
}

func TestWithDelimiter(t *testing.T) {
	type args struct {
		delimiter rune
	}
	type want struct {
		delimiter rune
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "sets Unicode delimiter",
			args: args{
				delimiter: '・',
			},
			want: want{
				delimiter: '・',
			},
		},
		{
			name: "sets zero",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &option{
				delimiter: '\t',
			}
			WithDelimiter(test.args.delimiter)(o)
			got := want{
				delimiter: o.delimiter,
			}
			testutil.AssertValue(t, got, test.want, "WithDelimiter")
		})
	}
}

func TestWithCRLF(t *testing.T) {
	type fields struct {
		crlf bool
	}
	type want struct {
		crlf bool
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "enables CRLF",
			want: want{
				crlf: true,
			},
		},
		{
			name: "keeps CRLF enabled",
			fields: fields{
				crlf: true,
			},
			want: want{
				crlf: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &option{
				crlf: test.fields.crlf,
			}
			WithCRLF()(o)
			got := want{
				crlf: o.crlf,
			}
			testutil.AssertValue(t, got, test.want, "WithCRLF")
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
		columns ColumnSelector
		fn      func(table.Value) string
		value   table.Value
	}
	type want struct {
		value string
		isNil bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "sets function",
			args: args{
				columns: Columns(0),
				fn: func(table.Value) string {
					return "new"
				},
				value: table.String("raw"),
			},
			want: want{
				value: "new",
			},
		},
		{
			name: "clears function",
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
							transformer: func(table.Value) string {
								return "old"
							},
						},
					},
				},
			}
			WithTransformer(test.args.columns, test.args.fn)(o)
			fn := o.columns.Values[0].transformer
			got := want{
				isNil: fn == nil,
			}
			if fn != nil {
				got.value = fn(test.args.value)
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
		selector ColumnSelector
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "no indexes",
			want: want{
				selector: ColumnSelector{},
			},
		},
		{
			name: "clones indexes",
			args: args{
				indexes: []int{0, -1, 2},
				mutate:  []int{3, 4, 5},
			},
			want: want{
				selector: ColumnSelector{
					selector: column.NewSelector(0, -1, 2),
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			selector := Columns(test.args.indexes...)
			copy(test.args.indexes, test.args.mutate)
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
			name: "selects all columns",
			want: want{
				selector: ColumnSelector{
					selector: column.All(),
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := want{
				selector: AllColumns(),
			}
			testutil.AssertValue(t, got, test.want, "AllColumns")
		})
	}
}
