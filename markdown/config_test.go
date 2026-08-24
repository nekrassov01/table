package markdown

import (
	"errors"
	"testing"

	"github.com/nekrassov01/table"
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
			name: "index defaults and explicit columns",
			fields: fields{
				output: func() configResult {
					defaults := columnConfig{
						align: AlignLeft,
					}
					return configResult{
						option: &option{
							columns: columnSet{
								Values: []columnConfig{
									{
										align: AlignRight,
									},
								},
								Defaults: &defaults,
							},
							indexOffset: 1,
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
					{},
					{
						align: AlignRight,
					},
					{
						align: AlignLeft,
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
						colspan: true,
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
						rowspan: true,
						colspan: true,
					},
				},
				defaults: &columnConfig{
					rowspan: true,
				},
			},
		},
		{
			name: "selected columns inherit defaults",
			fields: fields{
				values: []columnConfig{
					{
						align: AlignRight,
					},
				},
				defaults: &columnConfig{
					colspan: true,
				},
			},
			args: args{
				selector: Columns(-1, 2),
				fn:       setRowspan,
			},
			want: want{
				values: []columnConfig{
					{
						align: AlignRight,
					},
					{
						colspan: true,
					},
					{
						rowspan: true,
						colspan: true,
					},
				},
				defaults: &columnConfig{
					colspan: true,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &columnSet{
				Values:   test.fields.values,
				Defaults: test.fields.defaults,
			}
			o.apply(test.args.selector, test.args.fn)
			got := want{
				values:   o.Values,
				defaults: o.Defaults,
			}
			testutil.AssertValue(t, got, test.want, "apply")
		})
	}
}
