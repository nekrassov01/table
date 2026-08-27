package csv

import (
	"testing"

	"github.com/nekrassov01/table/internal/testutil"
)

func Test_quoteValue(t *testing.T) {
	type fields struct {
		delimiter rune
		quotes    []byte
		crlf      bool
	}
	type args struct {
		value string
	}
	type want struct {
		value  string
		quotes string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "keeps plain value",
			fields: fields{
				delimiter: ',',
			},
			args: args{
				value: "plain",
			},
			want: want{
				value: "plain",
			},
		},
		{
			name: "quotes delimiter",
			fields: fields{
				delimiter: ',',
				quotes:    []byte("old"),
			},
			args: args{
				value: "a,b",
			},
			want: want{
				value:  `"a,b"`,
				quotes: `old"a,b"`,
			},
		},
		{
			name: "doubles quotes and preserves line breaks",
			fields: fields{
				delimiter: '\t',
			},
			args: args{
				value: "a\"b\r\nc",
			},
			want: want{
				value:  "\"a\"\"b\r\nc\"",
				quotes: "\"a\"\"b\r\nc\"",
			},
		},
		{
			name: "normalizes line breaks in CRLF mode",
			fields: fields{
				delimiter: '\t',
				crlf:      true,
			},
			args: args{
				value: "a\r\nb\rc\nd",
			},
			want: want{
				value:  "\"a\r\nbc\r\nd\"",
				quotes: "\"a\r\nbc\r\nd\"",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := compilerState{
				quotes: test.fields.quotes,
			}
			o := &compiler{
				input: configResult{
					option: &option{
						delimiter: test.fields.delimiter,
						crlf:      test.fields.crlf,
					},
				},
				state: &state,
			}
			got := want{
				value:  o.quoteValue(test.args.value),
				quotes: string(state.quotes),
			}
			testutil.AssertValue(t, got, test.want, "quoteValue")
		})
	}
}

func Test_compiler_indexQuoteValue(t *testing.T) {
	type fields struct {
		delimiter rune
	}
	type args struct {
		value string
	}
	type want struct {
		index int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "plain",
			fields: fields{
				delimiter: ',',
			},
			args: args{
				value: "plain",
			},
			want: want{
				index: -1,
			},
		},
		{
			name: "delimiter",
			fields: fields{
				delimiter: ',',
			},
			args: args{
				value: "a,b",
			},
			want: want{
				index: 1,
			},
		},
		{
			name: "Unicode delimiter",
			fields: fields{
				delimiter: '・',
			},
			args: args{
				value: "a・b",
			},
			want: want{
				index: 1,
			},
		},
		{
			name: "quote",
			fields: fields{
				delimiter: '\t',
			},
			args: args{
				value: `a"b`,
			},
			want: want{
				index: 1,
			},
		},
		{
			name: "line feed",
			fields: fields{
				delimiter: '\t',
			},
			args: args{
				value: "a\nb",
			},
			want: want{
				index: 1,
			},
		},
		{
			name: "carriage return",
			fields: fields{
				delimiter: '\t',
			},
			args: args{
				value: "a\rb",
			},
			want: want{
				index: 1,
			},
		},
		{
			name: "leading ASCII space",
			fields: fields{
				delimiter: ',',
			},
			args: args{
				value: " value",
			},
			want: want{
				index: 0,
			},
		},
		{
			name: "leading Unicode space",
			fields: fields{
				delimiter: ',',
			},
			args: args{
				value: "\u00a0value",
			},
			want: want{
				index: 0,
			},
		},
		{
			name: "leading ASCII vertical tab",
			fields: fields{
				delimiter: ',',
			},
			args: args{
				value: "\vvalue",
			},
			want: want{
				index: 0,
			},
		},
		{
			name: "trailing space",
			fields: fields{
				delimiter: ',',
			},
			args: args{
				value: "value ",
			},
			want: want{
				index: -1,
			},
		},
		{
			name: "PostgreSQL terminator",
			fields: fields{
				delimiter: ',',
			},
			args: args{
				value: `\.`,
			},
			want: want{
				index: 0,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &compiler{
				input: configResult{
					option: &option{
						delimiter: test.fields.delimiter,
					},
				},
			}
			got := want{
				index: o.indexQuoteValue(test.args.value),
			}
			testutil.AssertValue(t, got, test.want, "indexQuoteValue")
		})
	}
}
