package markdown

import (
	"errors"
	"testing"

	"github.com/nekrassov01/table"
	"github.com/nekrassov01/table/internal/column"
	"github.com/nekrassov01/table/internal/testutil"
)

func Test_config_prepare(t *testing.T) {
	type fields struct {
		output configResult
		state  configState
	}
	type want struct {
		columns       []columnConfig
		output        configResult
		isHeaderError bool
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "required header",
			fields: fields{
				output: configResult{
					option: &option{},
				},
			},
			want: want{
				isHeaderError: true,
			},
		},
		{
			name: "defaults and explicit columns",
			fields: fields{
				output: func() configResult {
					defaults := columnConfig{
						align:       AlignLeft,
						placeholder: placeholder,
					}
					return configResult{
						option: &option{
							columns: columnSetOf(
								[]columnConfig{
									{
										align:       AlignRight,
										placeholder: placeholder,
									},
								},
								&defaults,
							),
						},
						header:   []string{"a", "b"},
						bodyRows: 2,
					}
				}(),
				state: configState{
					columns: make([]columnConfig, 1, 4),
				},
			},
			want: want{
				columns: []columnConfig{
					{
						align:       AlignRight,
						placeholder: placeholder,
					},
					{
						align:       AlignLeft,
						placeholder: placeholder,
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			o := &config{
				output: test.fields.output,
				state:  &state,
			}
			o.prepare()
			expected := test.want
			expected.output = test.fields.output
			if !expected.isHeaderError {
				expected.output.columns = test.want.columns
			}
			got := want{
				columns:       state.columns,
				output:        o.output,
				isHeaderError: errors.Is(o.err, table.ErrHeaderRequired),
			}
			testutil.AssertValue(t, got, expected, "prepare")
		})
	}
}

func Test_option_apply(t *testing.T) {
	type fields struct {
		placeholder string
	}
	type args struct {
		opts []Option
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
			name:   "defaults",
			fields: fields{placeholder: placeholder},
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
					WithPlaceholder(AllColumns(), "first"),
					WithPlaceholder(AllColumns(), "last"),
				},
			},
			want: want{
				placeholder: "last",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &option{}
			WithPlaceholder(AllColumns(), test.fields.placeholder)(o)
			o.apply(test.args.opts...)
			got := want{
				placeholder: o.columns.resolve(nil, 1)[0].placeholder,
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
	setRowspan := func(column *columnConfig) {
		column.rowspan = true
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "all columns creates and updates defaults",
			fields: fields{
				values: []columnConfig{
					{
						colspan:     true,
						placeholder: placeholder,
					},
				},
			},
			args: args{
				selector: AllColumns(),
				fn:       setRowspan,
			},
			want: want{
				values: []columnConfig{
					{
						rowspan:     true,
						colspan:     true,
						placeholder: placeholder,
					},
				},
				defaults: &columnConfig{
					rowspan:     true,
					placeholder: placeholder,
				},
			},
		},
		{
			name: "selected columns inherit defaults",
			fields: fields{
				values: []columnConfig{
					{
						align:       AlignRight,
						placeholder: placeholder,
					},
				},
				defaults: &columnConfig{
					colspan:     true,
					placeholder: placeholder,
				},
			},
			args: args{
				selector: Columns(-1, 2),
				fn:       setRowspan,
			},
			want: want{
				values: []columnConfig{
					{
						align:       AlignRight,
						placeholder: placeholder,
					},
					{
						colspan:     true,
						placeholder: placeholder,
					},
					{
						rowspan:     true,
						colspan:     true,
						placeholder: placeholder,
					},
				},
				defaults: &columnConfig{
					colspan:     true,
					placeholder: placeholder,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := columnSetOf(test.fields.values, test.fields.defaults)
			o.apply(test.args.selector, test.args.fn)
			got := want{
				values:   o.resolve(nil, len(test.want.values)),
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
