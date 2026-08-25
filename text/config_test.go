package text

import (
	"bytes"
	"io"
	"testing"

	"github.com/nekrassov01/table/internal/column"
	"github.com/nekrassov01/table/internal/testutil"
)

func Test_config_prepare(t *testing.T) {
	type fields struct {
		bodyColumns int
		state       configState
		output      configResult
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
			name: "header with index and column settings",
			fields: fields{
				output: func() configResult {
					defaults := columnConfig{
						limit: 9,
						lPad:  2,
						rPad:  3,
					}
					return configResult{
						option: &option{
							columns: columnSetOf(
								[]columnConfig{
									{
										limit: 7,
										lPad:  4,
										rPad:  5,
									},
								},
								&defaults,
							),
							indexOffset: 1,
						},
						header:   [][]string{{"a", "b"}},
						footer:   [][]string{{"a", "b", "c", "d"}},
						bodyRows: 2,
					}
				}(),
				bodyColumns: 3,
				state: configState{
					columns: make([]columnConfig, 3, 4),
				},
			},
			want: want{
				columns: []columnConfig{
					defaultColumn(),
					{
						limit: 7,
						lPad:  4,
						rPad:  5,
					},
					{
						limit: 9,
						lPad:  2,
						rPad:  3,
					},
				},
				footerColumns: 4,
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
			},
			want: want{
				columns: []columnConfig{
					defaultColumn(),
					defaultColumn(),
					defaultColumn(),
				},
				footerColumns: 3,
			},
		},
		{
			name: "header is authoritative",
			fields: fields{
				bodyColumns: 3,
				output: configResult{
					option: &option{},
					header: [][]string{{"header"}},
					footer: [][]string{{"a", "b", "c", "d"}},
				},
			},
			want: want{
				columns: []columnConfig{
					defaultColumn(),
				},
				footerColumns: 4,
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
				state: configState{
					columns: []columnConfig{{
						limit: 9,
					}},
				},
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
				state:       &state,
				output:      test.fields.output,
			}
			o.prepare()
			expectedOutput := test.fields.output
			expectedOutput.columns = test.want.columns
			expectedOutput.footerColumns = test.want.footerColumns
			got := want{
				columns:       state.columns,
				footerColumns: o.output.footerColumns,
				output:        o.output,
			}
			expected := test.want
			expected.output = expectedOutput
			testutil.AssertValue(t, got, expected, "prepare")
		})
	}
}

func Test_option_apply(t *testing.T) {
	type fields struct {
		style       Style
		placeholder string
		header      [][]string
		footer      func() [][]string
		caption     string
		columns     columnSet
		indexOffset int
		indexWidth  int
		captionSide CaptionSide
		compact     bool
		autoFit     bool
		plain       bool
	}
	type args struct {
		w             io.Writer
		minIndexWidth int
		opts          []Option
		terminal      bool
	}
	type want struct {
		placeholder string
		columns     columnSet
		indexOffset int
		indexWidth  int
		plain       bool
		autoFit     bool
		borderAttr  *Attr
		content     ContentStyle
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "defaults on non-terminal",
			args: args{
				w: &bytes.Buffer{},
			},
			want: want{
				placeholder: placeholder,
				plain:       true,
			},
		},
		{
			name: "terminal keeps attributes and auto fit",
			args: args{
				w:             &bytes.Buffer{},
				minIndexWidth: 3,
				opts: []Option{
					WithStyle(StyleColoredLight),
					WithPlaceholder("-"),
					WithIndexWidth(2),
					WithAutoFit(),
					WithAttr(ScopeBody, AllColumns(), NewAttr(CodeUnderline)),
				},
				terminal: true,
			},
			want: want{
				placeholder: "-",
				columns: func() columnSet {
					attr := NewAttr(CodeUnderline)
					configured := defaultColumn()
					configured.transformer.attrs.Set(ScopeBody, attr)
					return columnSetOf(nil, &configured)
				}(),
				indexOffset: 1,
				indexWidth:  3,
				autoFit:     true,
				borderAttr:  StyleColoredLight.Border.Attr,
				content:     StyleColoredLight.Content,
			},
		},
		{
			name: "non-terminal strips attributes",
			args: args{
				w: &bytes.Buffer{},
				opts: func() []Option {
					attr := NewAttr(CodeUnderline)
					return []Option{
						WithStyle(StyleColoredLight),
						WithAttr(ScopeBody, AllColumns(), attr),
						WithAttr(ScopeHeader, Columns(0), attr),
					}
				}(),
			},
			want: want{
				placeholder: placeholder,
				columns: func() columnSet {
					plainDefault := defaultColumn()
					return columnSetOf([]columnConfig{defaultColumn()}, &plainDefault)
				}(),
				plain: true,
			},
		},
		{
			name: "default width disables auto fit",
			args: args{
				w: &bytes.Buffer{},
				opts: []Option{
					WithAutoFit(),
					WithWidth(AllColumns(), 4),
				},
				terminal: true,
			},
			want: want{
				placeholder: placeholder,
				columns: func() columnSet {
					configured := defaultColumn()
					configured.limit = 4
					return columnSetOf(nil, &configured)
				}(),
			},
		},
		{
			name: "default truncation disables auto fit",
			args: args{
				w: &bytes.Buffer{},
				opts: []Option{
					WithAutoFit(),
					WithTruncate(AllColumns()),
				},
				terminal: true,
			},
			want: want{
				placeholder: placeholder,
				columns: func() columnSet {
					configured := defaultColumn()
					configured.truncate = true
					return columnSetOf(nil, &configured)
				}(),
			},
		},
		{
			name: "column width disables auto fit",
			args: args{
				w: &bytes.Buffer{},
				opts: []Option{
					WithAutoFit(),
					WithWidth(Columns(0), 4),
				},
				terminal: true,
			},
			want: want{
				placeholder: placeholder,
				columns: func() columnSet {
					configured := defaultColumn()
					configured.limit = 4
					return columnSet{
						Values: []columnConfig{configured},
					}
				}(),
			},
		},
		{
			name: "column truncation disables auto fit",
			args: args{
				w: &bytes.Buffer{},
				opts: []Option{
					WithAutoFit(),
					WithTruncate(Columns(0)),
				},
				terminal: true,
			},
			want: want{
				placeholder: placeholder,
				columns: func() columnSet {
					configured := defaultColumn()
					configured.truncate = true
					return columnSet{
						Values: []columnConfig{configured},
					}
				}(),
			},
		},
	}
	restore := isTerminal
	t.Cleanup(func() {
		isTerminal = restore
	})
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isTerminal = func(io.Writer) bool {
				return test.args.terminal
			}
			o := &option{
				style:       test.fields.style,
				placeholder: test.fields.placeholder,
				header:      test.fields.header,
				footer:      test.fields.footer,
				caption:     test.fields.caption,
				columns:     test.fields.columns,
				indexOffset: test.fields.indexOffset,
				indexWidth:  test.fields.indexWidth,
				captionSide: test.fields.captionSide,
				compact:     test.fields.compact,
				autoFit:     test.fields.autoFit,
				plain:       test.fields.plain,
			}
			o.apply(test.args.w, test.args.minIndexWidth, test.args.opts...)
			got := want{
				placeholder: o.placeholder,
				columns:     o.columns,
				indexOffset: o.indexOffset,
				indexWidth:  o.indexWidth,
				plain:       o.plain,
				autoFit:     o.autoFit,
				borderAttr:  o.style.Border.Attr,
				content:     o.style.Content,
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
	increment := func(column *columnConfig) {
		column.limit++
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
					limit: 1,
				}},
			},
			args: args{
				selector: AllColumns(),
				fn:       increment,
			},
			want: want{
				values: []columnConfig{{
					limit: 2,
				}},
				defaults: &columnConfig{
					limit: 1,
					lPad:  1,
					rPad:  1,
				},
			},
		},
		{
			name: "all columns updates defaults",
			fields: fields{
				values: []columnConfig{{}},
				defaults: &columnConfig{
					limit: 2,
					lPad:  1,
					rPad:  1,
				},
			},
			args: args{
				selector: AllColumns(),
				fn:       increment,
			},
			want: want{
				values: []columnConfig{{
					limit: 1,
				}},
				defaults: &columnConfig{
					limit: 3,
					lPad:  1,
					rPad:  1,
				},
			},
		},
		{
			name: "indexes extend default columns",
			fields: fields{
				values: []columnConfig{{
					limit: 4,
				}},
			},
			args: args{
				selector: Columns(2, -1, 2),
				fn:       increment,
			},
			want: want{
				values: []columnConfig{
					{
						limit: 4,
					},
					defaultColumn(),
					{
						limit: 1,
						lPad:  1,
						rPad:  1,
					},
				},
			},
		},
		{
			name: "indexes inherit configured defaults",
			fields: fields{
				defaults: &columnConfig{
					limit: 5,
					lPad:  2,
					rPad:  3,
				},
			},
			args: args{
				selector: Columns(1),
				fn:       increment,
			},
			want: want{
				values: []columnConfig{
					{
						limit: 5,
						lPad:  2,
						rPad:  3,
					},
					{
						limit: 6,
						lPad:  2,
						rPad:  3,
					},
				},
				defaults: &columnConfig{
					limit: 5,
					lPad:  2,
					rPad:  3,
				},
			},
		},
		{
			name: "negative index is ignored",
			fields: fields{
				values: []columnConfig{{}},
			},
			args: args{
				selector: Columns(-1, 0),
				fn:       increment,
			},
			want: want{
				values: []columnConfig{{
					limit: 1,
				}},
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

func Test_defaultColumn(t *testing.T) {
	type want struct {
		val columnConfig
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "default padding",
			want: want{
				val: columnConfig{
					lPad: 1,
					rPad: 1,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := defaultColumn()
			testutil.AssertValue(t, got, test.want.val, "defaultColumn")
		})
	}
}

func columnSetOf(values []columnConfig, defaults *columnConfig) columnSet {
	set := columnSet{}
	if defaults != nil {
		(*column.Set[columnConfig])(&set).Apply(column.All(), defaultColumn, func(config *columnConfig) {
			*config = *defaults
		})
	}
	set.Values = values
	return set
}
