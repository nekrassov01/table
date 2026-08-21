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
		columns       []column
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
					defaults := column{
						align: AlignLeft,
					}
					return configResult{
						option: &option{
							columns: columnSet{
								values: []column{
									{
										align: AlignRight,
									},
								},
								defaults: &defaults,
							},
							indexOffset: 1,
						},
						header:   []string{"a", "b"},
						bodyRows: 2,
					}
				}(),
				state: configState{
					columns: make([]column, 1, 4),
				},
			},
			want: want{
				columns: []column{
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
				placeholder: DefaultPlaceholder,
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
	setRowspan := func(column *column) {
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
				values: []column{
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
				values: []column{
					{
						rowspan: true,
						colspan: true,
					},
				},
				defaults: &column{
					rowspan: true,
				},
			},
		},
		{
			name: "selected columns inherit defaults",
			fields: fields{
				values: []column{
					{
						align: AlignRight,
					},
				},
				defaults: &column{
					colspan: true,
				},
			},
			args: args{
				selector: Columns(-1, 2),
				fn:       setRowspan,
			},
			want: want{
				values: []column{
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
				defaults: &column{
					colspan: true,
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
