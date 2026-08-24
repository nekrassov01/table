package csv

import (
	"errors"
	"testing"
	"unicode/utf8"

	"github.com/nekrassov01/table"
	"github.com/nekrassov01/table/internal/testutil"
)

func Test_config_prepare(t *testing.T) {
	type fields struct {
		bodyColumns int
		output      configResult
		state       configState
	}
	type want struct {
		columnCount   int
		footerColumns int
		transformed   []string
		isDelimiter   bool
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "rejects invalid delimiter",
			fields: fields{
				output: configResult{
					option: &option{
						delimiter: '"',
					},
					header: []string{"A"},
				},
			},
			want: want{
				isDelimiter: true,
			},
		},
		{
			name: "header fixes columns and applies index defaults and overrides",
			fields: fields{
				bodyColumns: 5,
				output: configResult{
					option: &option{
						delimiter:   ',',
						indexOffset: 1,
						columns: columnSet{
							Values: []columnConfig{
								{},
								{
									transformer: func(any) string {
										return "explicit"
									},
								},
							},
							Defaults: &columnConfig{
								transformer: func(any) string {
									return "default"
								},
							},
						},
					},
					header: []string{"A", "B", "C"},
					footer: [][]string{{"x", "y", "z", "ignored"}},
				},
				state: configState{
					columns: make([]columnConfig, 0, 4),
				},
			},
			want: want{
				columnCount:   4,
				footerColumns: 4,
				transformed:   []string{"", "", "explicit", "default"},
			},
		},
		{
			name: "body and footer establish headerless columns",
			fields: fields{
				bodyColumns: 2,
				output: configResult{
					option: &option{
						delimiter: '\t',
					},
					footer: [][]string{{"a"}, {"a", "b", "c"}},
				},
			},
			want: want{
				columnCount:   3,
				footerColumns: 3,
				transformed:   []string{"", "", ""},
			},
		},
		{
			name: "keeps zero columns",
			fields: fields{
				output: configResult{
					option: &option{
						delimiter: '\t',
					},
				},
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
			var transformed []string
			if len(o.output.columns) > 0 {
				transformed = make([]string, len(o.output.columns))
				for index := range o.output.columns {
					if fn := o.output.columns[index].transformer; fn != nil {
						transformed[index] = fn(nil)
					}
				}
			}
			got := want{
				columnCount:   len(o.output.columns),
				footerColumns: o.output.footerColumns,
				transformed:   transformed,
				isDelimiter:   errors.Is(o.err, table.ErrDelimiter),
			}
			testutil.AssertValue(t, got, test.want, "prepare")
		})
	}
}

func Test_option_apply(t *testing.T) {
	type args struct {
		opts []Option
	}
	type want struct {
		delimiter   rune
		placeholder string
		header      []string
		crlf        bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "sets defaults",
			want: want{
				delimiter:   '\t',
				placeholder: DefaultPlaceholder,
			},
		},
		{
			name: "applies options in order",
			args: args{
				opts: []Option{
					WithDelimiter(','),
					WithPlaceholder("old"),
					WithPlaceholder("-"),
					WithHeader([]string{"A"}),
					WithCRLF(),
				},
			},
			want: want{
				delimiter:   ',',
				placeholder: "-",
				header:      []string{"A"},
				crlf:        true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &option{}
			o.apply(test.args.opts...)
			got := want{
				delimiter:   o.delimiter,
				placeholder: o.placeholder,
				header:      o.header,
				crlf:        o.crlf,
			}
			testutil.AssertValue(t, got, test.want, "apply")
		})
	}
}

func Test_columnSet_apply(t *testing.T) {
	type fields struct {
		values       []columnConfig
		defaultValue string
	}
	type args struct {
		selector ColumnSelector
		value    string
	}
	type want struct {
		values       []string
		defaultValue string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "applies to all existing and future columns",
			fields: fields{
				values: []columnConfig{{}, {}},
			},
			args: args{
				selector: AllColumns(),
				value:    "all",
			},
			want: want{
				values:       []string{"all", "all"},
				defaultValue: "all",
			},
		},
		{
			name: "replaces existing default",
			fields: fields{
				values:       []columnConfig{{}},
				defaultValue: "old",
			},
			args: args{
				selector: AllColumns(),
				value:    "new",
			},
			want: want{
				values:       []string{"new"},
				defaultValue: "new",
			},
		},
		{
			name: "grows explicit columns from default and ignores negative index",
			fields: fields{
				values:       []columnConfig{{}},
				defaultValue: "default",
			},
			args: args{
				selector: Columns(-1, 2),
				value:    "explicit",
			},
			want: want{
				values:       []string{"", "default", "explicit"},
				defaultValue: "default",
			},
		},
		{
			name: "updates existing explicit column without growing",
			fields: fields{
				values: []columnConfig{{}, {}},
			},
			args: args{
				selector: Columns(0),
				value:    "first",
			},
			want: want{
				values: []string{"first", ""},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			set := columnSet{
				Values: test.fields.values,
			}
			if test.fields.defaultValue != "" {
				set.Defaults = &columnConfig{
					transformer: func(any) string {
						return test.fields.defaultValue
					},
				}
			}
			set.apply(test.args.selector, func(c *columnConfig) {
				value := test.args.value
				c.transformer = func(any) string {
					return value
				}
			})
			values := make([]string, len(set.Values))
			for index := range set.Values {
				if fn := set.Values[index].transformer; fn != nil {
					values[index] = fn(nil)
				}
			}
			defaultValue := ""
			if set.Defaults != nil && set.Defaults.transformer != nil {
				defaultValue = set.Defaults.transformer(nil)
			}
			got := want{
				values:       values,
				defaultValue: defaultValue,
			}
			testutil.AssertValue(t, got, test.want, "apply")
		})
	}
}

func Test_validDelimiter(t *testing.T) {
	type args struct {
		delimiter rune
	}
	type want struct {
		valid bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "tab",
			args: args{
				delimiter: '\t',
			},
			want: want{
				valid: true,
			},
		},
		{
			name: "space",
			args: args{
				delimiter: ' ',
			},
			want: want{
				valid: true,
			},
		},
		{
			name: "tilde",
			args: args{
				delimiter: '~',
			},
			want: want{
				valid: true,
			},
		},
		{
			name: "Unicode",
			args: args{
				delimiter: '・',
			},
			want: want{
				valid: true,
			},
		},
		{
			name: "vertical tab",
			args: args{
				delimiter: '\v',
			},
			want: want{
				valid: true,
			},
		},
		{
			name: "quote",
			args: args{
				delimiter: '"',
			},
		},
		{
			name: "line feed",
			args: args{
				delimiter: '\n',
			},
		},
		{
			name: "carriage return",
			args: args{
				delimiter: '\r',
			},
		},
		{
			name: "NUL",
		},
		{
			name: "replacement rune",
			args: args{
				delimiter: utf8.RuneError,
			},
		},
		{
			name: "negative rune",
			args: args{
				delimiter: -1,
			},
		},
		{
			name: "surrogate rune",
			args: args{
				delimiter: 0xd800,
			},
		},
		{
			name: "rune above maximum",
			args: args{
				delimiter: utf8.MaxRune + 1,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := want{
				valid: validDelimiter(test.args.delimiter),
			}
			testutil.AssertValue(t, got, test.want, "validDelimiter")
		})
	}
}
