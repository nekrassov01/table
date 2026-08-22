package markdown

import (
	"testing"

	"github.com/nekrassov01/table/internal/testutil"
)

func Test_escapeCode(t *testing.T) {
	type fields struct {
		escapes []byte
	}
	type args struct {
		value string
	}
	type want struct {
		value   string
		escapes string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "unchanged",
			args: args{
				value: `a\b`,
			},
			want: want{
				value: `a\b`,
			},
		},
		{
			name: "pads edge backtick",
			args: args{
				value: "`x",
			},
			want: want{
				value:   " `x ",
				escapes: " `x ",
			},
		},
		{
			name: "escapes separator and line break",
			fields: fields{
				escapes: []byte("kept:"),
			},
			args: args{
				value: "a|b\r\nc",
			},
			want: want{
				value:   `a\|b c`,
				escapes: `kept:a\|b c`,
			},
		},
		{
			name: "keeps valid utf-8 while escaping",
			args: args{
				value: "日|本",
			},
			want: want{
				value:   `日\|本`,
				escapes: `日\|本`,
			},
		},
		{
			name: "pads spaces at both ends",
			args: args{
				value: " x ",
			},
			want: want{
				value:   "  x  ",
				escapes: "  x  ",
			},
		},
		{
			name: "keeps only spaces unpadded",
			args: args{
				value: "   ",
			},
			want: want{
				value: "   ",
			},
		},
		{
			name: "normalizes all-space line ending",
			args: args{
				value: " \r\n ",
			},
			want: want{
				value:   "   ",
				escapes: "   ",
			},
		},
		{
			name: "replaces nul",
			args: args{
				value: "a\x00b",
			},
			want: want{
				value:   "a�b",
				escapes: "a�b",
			},
		},
		{
			name: "replaces malformed utf-8",
			args: args{
				value: "a\xffb",
			},
			want: want{
				value:   "a�b",
				escapes: "a�b",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value, escapes := escapeCode(test.fields.escapes, test.args.value)
			got := want{
				value:   value,
				escapes: string(escapes),
			}
			testutil.AssertValue(t, got, test.want, "escapeCode")
		})
	}
}

func Test_escapeValue(t *testing.T) {
	type fields struct {
		escapes []byte
	}
	type args struct {
		value string
	}
	type want struct {
		value   string
		escapes string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "unchanged",
			args: args{
				value: " plain  日本語 ",
			},
			want: want{
				value: " plain  日本語 ",
			},
		},
		{
			name: "escapes markdown cell content",
			fields: fields{
				escapes: []byte("kept:"),
			},
			args: args{
				value: " a  b|c\\d\r\ne\x00\xff",
			},
			want: want{
				value:   " a  b\\|c\\\\d<br>e��",
				escapes: "kept: a  b\\|c\\\\d<br>e��",
			},
		},
		{
			name: "literalizes inline markup",
			args: args{
				value: "**bold** _em_ ~~del~~ [link](url) <tag> &nbsp; `code`",
			},
			want: want{
				value:   "\\*\\*bold\\*\\* \\_em\\_ \\~\\~del\\~\\~ \\[link\\](url) \\<tag\\> \\&nbsp; \\`code\\`",
				escapes: "\\*\\*bold\\*\\* \\_em\\_ \\~\\~del\\~\\~ \\[link\\](url) \\<tag\\> \\&nbsp; \\`code\\`",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value, escapes := escapeValue(test.fields.escapes, test.args.value)
			got := want{
				value:   value,
				escapes: string(escapes),
			}
			testutil.AssertValue(t, got, test.want, "escapeValue")
		})
	}
}

func Test_needsEscapeValue(t *testing.T) {
	type args struct {
		value string
	}
	type want struct {
		val bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "plain with interior space",
			args: args{
				value: "plain value",
			},
		},
		{
			name: "markdown syntax",
			args: args{
				value: "\\|`*_~[]<>&",
			},
			want: want{
				val: true,
			},
		},
		{
			name: "keeps edge space",
			args: args{
				value: " value",
			},
		},
		{
			name: "keeps space run",
			args: args{
				value: "a  b",
			},
		},
		{
			name: "nul",
			args: args{
				value: "a\x00b",
			},
			want: want{
				val: true,
			},
		},
		{
			name: "valid utf-8",
			args: args{
				value: "日本語",
			},
		},
		{
			name: "malformed utf-8",
			args: args{
				value: "a\xffb",
			},
			want: want{
				val: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := needsEscapeValue(test.args.value)
			testutil.AssertValue(t, got, test.want.val, "needsEscapeValue")
		})
	}
}

func Test_escapeAttr(t *testing.T) {
	type args struct {
		value string
	}
	type want struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "unchanged",
			args: args{
				value: "red",
			},
			want: want{
				value: "red",
			},
		},
		{
			name: "escapes attribute and table delimiters",
			args: args{
				value: `a&"|b`,
			},
			want: want{
				value: `a&amp;&quot;&#124;b`,
			},
		},
		{
			name: "normalizes line endings without splitting rows",
			args: args{
				value: "a\r\nb\rc\n",
			},
			want: want{
				value: "a b c ",
			},
		},
		{
			name: "preserves allowed whitespace",
			args: args{
				value: "a\t\fb",
			},
			want: want{
				value: "a\t\fb",
			},
		},
		{
			name: "replaces controls and malformed utf-8",
			args: args{
				value: "\x00\x01\x7f\u0080\xff",
			},
			want: want{
				value: "�����",
			},
		},
		{
			name: "replaces noncharacters",
			args: args{
				value: "\ufdd0\ufffe\U0001ffff",
			},
			want: want{
				value: "���",
			},
		},
		{
			name: "preserves valid unicode",
			args: args{
				value: "日本語�",
			},
			want: want{
				value: "日本語�",
			},
		},
		{
			name: "preserves valid unicode while escaping",
			args: args{
				value: "日&本",
			},
			want: want{
				value: "日&amp;本",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := want{
				value: escapeAttr(test.args.value),
			}
			testutil.AssertValue(t, got, test.want, "escapeAttr")
		})
	}
}

func Test_needsEscapeAttr(t *testing.T) {
	type args struct {
		value string
	}
	type want struct {
		val bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "plain",
			args: args{
				value: "plain 日本語 �",
			},
		},
		{
			name: "allowed whitespace",
			args: args{
				value: "\t\f ",
			},
		},
		{
			name: "attribute and table delimiters",
			args: args{
				value: `&"|`,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "line endings",
			args: args{
				value: "\r\n",
			},
			want: want{
				val: true,
			},
		},
		{
			name: "controls",
			args: args{
				value: "\x00\x01\x7f\u0080",
			},
			want: want{
				val: true,
			},
		},
		{
			name: "malformed utf-8",
			args: args{
				value: "\xff",
			},
			want: want{
				val: true,
			},
		},
		{
			name: "noncharacter",
			args: args{
				value: "\ufdd0",
			},
			want: want{
				val: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := needsEscapeAttr(test.args.value)
			testutil.AssertValue(t, got, test.want.val, "needsEscapeAttr")
		})
	}
}
