package html

import (
	"testing"

	"github.com/nekrassov01/table/internal/testutil"
)

func Test_config_prepare(t *testing.T) {
	type fields struct {
		bodyColumns int
		output      configResult
		state       configState
	}
	type want struct {
		columns       []column
		footerColumns int
		output        configResult
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "header with index and configured column",
			fields: fields{
				output: func() configResult {
					configured := column{
						rowspan: ScopeBody,
					}
					defaults := column{
						colspan: ScopeFooter,
					}
					return configResult{
						option: &option{
							columns: columnSet{
								values:   []column{configured},
								defaults: &defaults,
							},
							indexOffset: 1,
						},
						header:   [][]string{{"a", "b"}},
						footer:   [][]string{{"a", "b", "c"}},
						bodyRows: 2,
					}
				}(),
				bodyColumns: 4,
			},
			want: want{
				columns: func() []column {
					configured := column{
						rowspan: ScopeBody,
					}
					defaults := column{
						colspan: ScopeFooter,
					}
					return []column{
						{},
						configured,
						defaults,
					}
				}(),
				footerColumns: 3,
			},
		},
		{
			name: "body and footer fallback",
			fields: fields{
				bodyColumns: 2,
				output: configResult{
					option: &option{},
					footer: [][]string{{"a", "b", "c"}},
				},
				state: configState{
					columns: make([]column, 1, 4),
				},
			},
			want: want{
				columns:       make([]column, 3),
				footerColumns: 3,
			},
		},
		{
			name: "empty input",
			fields: fields{
				output: configResult{
					option: &option{
						indexOffset: 1,
					},
				},
				state: func() configState {
					configured := column{
						rowspan: ScopeBody,
					}
					return configState{
						columns: []column{configured},
					}
				}(),
			},
			want: want{
				columns: []column{},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			o := &config{
				bodyColumns: test.fields.bodyColumns,
				output:      test.fields.output,
				state:       &state,
			}
			o.prepare()
			expected := test.want
			expected.output = test.fields.output
			expected.output.columns = test.want.columns
			expected.output.footerColumns = test.want.footerColumns
			got := want{
				columns:       state.columns,
				footerColumns: o.output.footerColumns,
				output:        o.output,
			}
			testutil.AssertValue(t, got, expected, "prepare")
		})
	}
}

func Test_option_apply(t *testing.T) {
	type fields struct {
		placeholder string
		indexOffset int
	}
	type args struct {
		opts []Option
	}
	type want struct {
		placeholder string
		indexOffset int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "defaults",
		},
		{
			name: "applies in order",
			fields: fields{
				placeholder: "old",
			},
			args: args{
				opts: []Option{
					WithPlaceholder("first"),
					WithPlaceholder("last"),
					WithIndex(),
				},
			},
			want: want{
				placeholder: "last",
				indexOffset: 1,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &option{
				placeholder: test.fields.placeholder,
				indexOffset: test.fields.indexOffset,
			}
			o.apply(test.args.opts...)
			got := want{
				placeholder: o.placeholder,
				indexOffset: o.indexOffset,
			}
			testutil.AssertValue(t, got, test.want, "apply")
		})
	}
}

func Test_columnSet_apply(t *testing.T) {
	type fields struct {
		values   []column
		defaults *column
	}
	type args struct {
		selector ColumnSelector
		fn       func(*column)
	}
	type want struct {
		values   []column
		defaults *column
	}
	setBody := func(column *column) {
		column.rowspan |= ScopeBody
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "all columns creates defaults",
			fields: fields{
				values: []column{{
					rowspan: ScopeHeader,
				}},
			},
			args: args{
				selector: AllColumns(),
				fn:       setBody,
			},
			want: want{
				values: []column{{
					rowspan: ScopeHeader | ScopeBody,
				}},
				defaults: &column{
					rowspan: ScopeBody,
				},
			},
		},
		{
			name: "all columns updates defaults",
			fields: fields{
				values: []column{{}},
				defaults: &column{
					rowspan: ScopeHeader,
				},
			},
			args: args{
				selector: AllColumns(),
				fn:       setBody,
			},
			want: want{
				values: []column{{
					rowspan: ScopeBody,
				}},
				defaults: &column{
					rowspan: ScopeHeader | ScopeBody,
				},
			},
		},
		{
			name: "indexes extend default columns",
			fields: fields{
				values: []column{{
					rowspan: ScopeHeader,
				}},
			},
			args: args{
				selector: Columns(2, -1, 0),
				fn:       setBody,
			},
			want: want{
				values: []column{
					{
						rowspan: ScopeHeader | ScopeBody,
					},
					{},
					{
						rowspan: ScopeBody,
					},
				},
			},
		},
		{
			name: "indexes inherit configured defaults",
			fields: fields{
				defaults: &column{
					colspan: ScopeFooter,
				},
			},
			args: args{
				selector: Columns(1),
				fn:       setBody,
			},
			want: want{
				values: []column{
					{
						colspan: ScopeFooter,
					},
					{
						rowspan: ScopeBody,
						colspan: ScopeFooter,
					},
				},
				defaults: &column{
					colspan: ScopeFooter,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &columnSet{
				values:   test.fields.values,
				defaults: test.fields.defaults,
			}
			o.apply(test.args.selector, test.args.fn)
			got := want{
				values:   o.values,
				defaults: o.defaults,
			}
			testutil.AssertValue(t, got, test.want, "apply")
		})
	}
}

func Test_maxColumns(t *testing.T) {
	type args struct {
		rows [][]string
	}
	type want struct {
		val int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "widest row",
			args: args{
				rows: [][]string{{"a"}, {"a", "b", "c"}, {"a", "b"}},
			},
			want: want{
				val: 3,
			},
		},
		{
			name: "empty",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := maxColumns(test.args.rows)
			testutil.AssertValue(t, got, test.want.val, "maxColumns")
		})
	}
}
