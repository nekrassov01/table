package backlog

import (
	"testing"

	"github.com/nekrassov01/table/internal/column"
	"github.com/nekrassov01/table/internal/testutil"
)

func Test_config_prepare(t *testing.T) {
	type fields struct {
		bodyColumns int
		output      configResult
		state       configState
	}
	type want struct {
		columns       []columnConfig
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
					configured := columnConfig{
						rowspan: ScopeBody,
					}
					defaults := columnConfig{
						colspan: ScopeFooter,
					}
					return configResult{
						option: &option{
							columns:     columnSetOf([]columnConfig{configured}, &defaults),
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
				columns: func() []columnConfig {
					configured := columnConfig{
						rowspan: ScopeBody,
					}
					defaults := columnConfig{
						colspan: ScopeFooter,
					}
					return []columnConfig{
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
					columns: make([]columnConfig, 1, 4),
				},
			},
			want: want{
				columns:       make([]columnConfig, 3),
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
					configured := columnConfig{
						rowspan: ScopeBody,
					}
					return configState{
						columns: []columnConfig{configured},
					}
				}(),
			},
			want: want{
				columns: []columnConfig{},
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
			want: want{
				placeholder: placeholder,
			},
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
		values   []columnConfig
		defaults *columnConfig
	}
	type args struct {
		selector ColumnSelector
		fn       func(*columnConfig)
	}
	type want struct {
		values   []columnConfig
		defaults *columnConfig
	}
	setBody := func(column *columnConfig) {
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
				values: []columnConfig{{
					rowspan: ScopeHeader,
				}},
			},
			args: args{
				selector: AllColumns(),
				fn:       setBody,
			},
			want: want{
				values: []columnConfig{{
					rowspan: ScopeHeader | ScopeBody,
				}},
				defaults: &columnConfig{
					rowspan: ScopeBody,
				},
			},
		},
		{
			name: "all columns updates defaults",
			fields: fields{
				values: []columnConfig{{}},
				defaults: &columnConfig{
					rowspan: ScopeHeader,
				},
			},
			args: args{
				selector: AllColumns(),
				fn:       setBody,
			},
			want: want{
				values: []columnConfig{{
					rowspan: ScopeBody,
				}},
				defaults: &columnConfig{
					rowspan: ScopeHeader | ScopeBody,
				},
			},
		},
		{
			name: "indexes extend default columns",
			fields: fields{
				values: []columnConfig{{
					rowspan: ScopeHeader,
				}},
			},
			args: args{
				selector: Columns(2, -1, 0),
				fn:       setBody,
			},
			want: want{
				values: []columnConfig{
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
				defaults: &columnConfig{
					colspan: ScopeFooter,
				},
			},
			args: args{
				selector: Columns(1),
				fn:       setBody,
			},
			want: want{
				values: []columnConfig{
					{
						colspan: ScopeFooter,
					},
					{
						rowspan: ScopeBody,
						colspan: ScopeFooter,
					},
				},
				defaults: &columnConfig{
					colspan: ScopeFooter,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := columnSetOf(test.fields.values, test.fields.defaults)
			o.apply(test.args.selector, test.args.fn)
			got := want{
				values:   o.resolve(nil, len(test.want.values), 0),
				defaults: (*column.Set[columnConfig])(&o).Default(),
			}
			testutil.AssertValue(t, got, test.want, "apply")
		})
	}
}

func columnSetOf(values []columnConfig, defaults *columnConfig) columnSet {
	set := columnSet{}
	if defaults != nil {
		(*column.Set[columnConfig])(&set).Apply(column.All(), nil, func(config *columnConfig) {
			*config = *defaults
		})
	}
	set.Values = values
	return set
}
