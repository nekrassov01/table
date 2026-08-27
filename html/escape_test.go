package html

import (
	"testing"

	"github.com/nekrassov01/table/internal/testutil"
)

func Test_escapeValue(t *testing.T) {
	type fields struct {
		escapes []byte
	}
	type args struct {
		value string
	}
	type want struct {
		val     string
		escapes string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "plain value",
			args: args{
				value: "plain",
			},
			want: want{
				val: "plain",
			},
		},
		{
			name: "escapes markup and line break",
			fields: fields{
				escapes: []byte("old"),
			},
			args: args{
				value: "<&\n",
			},
			want: want{
				val:     "&lt;&amp;<br>",
				escapes: "old&lt;&amp;<br>",
			},
		},
		{
			name: "escapes every markup delimiter",
			fields: fields{
				escapes: []byte("before:"),
			},
			args: args{
				value: `&<>"`,
			},
			want: want{
				val:     "&amp;&lt;&gt;&quot;",
				escapes: "before:&amp;&lt;&gt;&quot;",
			},
		},
		{
			name: "normalizes line breaks",
			args: args{
				value: "a\r\nb\rc\n",
			},
			want: want{
				val:     "a<br>b<br>c<br>",
				escapes: "a<br>b<br>c<br>",
			},
		},
		{
			name: "replaces controls and invalid UTF-8",
			args: args{
				value: "\t\x01\x7f\xff",
			},
			want: want{
				val:     "\t���",
				escapes: "\t���",
			},
		},
		{
			name: "retains valid UTF-8",
			args: args{
				value: "世",
			},
			want: want{
				val: "世",
			},
		},
		{
			name: "retains valid UTF-8 after escaped markup",
			args: args{
				value: "a&世界",
			},
			want: want{
				val:     "a&amp;世界",
				escapes: "a&amp;世界",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value, escapes := escapeValue(test.fields.escapes, test.args.value)
			got := want{
				val:     value,
				escapes: string(escapes),
			}
			testutil.AssertValue(t, got, test.want, "escapeValue")
		})
	}
}

func Test_indexEscapeValue(t *testing.T) {
	type args struct {
		value string
	}
	type want struct {
		index int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "plain ASCII",
			args: args{
				value: "plain\ttext",
			},
			want: want{
				index: -1,
			},
		},
		{
			name: "markup",
			args: args{
				value: "<&>\"",
			},
			want: want{
				index: 0,
			},
		},
		{
			name: "line break",
			args: args{
				value: "a\nb",
			},
			want: want{
				index: 1,
			},
		},
		{
			name: "control",
			args: args{
				value: "\x00",
			},
			want: want{
				index: 0,
			},
		},
		{
			name: "valid UTF-8",
			args: args{
				value: "日本語",
			},
			want: want{
				index: -1,
			},
		},
		{
			name: "invalid UTF-8",
			args: args{
				value: string([]byte{0xff}),
			},
			want: want{
				index: 0,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := indexEscapeValue(test.args.value)
			testutil.AssertValue(t, got, test.want.index, "indexEscapeValue")
		})
	}
}

func Test_escapeAttr(t *testing.T) {
	type args struct {
		value string
	}
	type want struct {
		val string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "plain",
			args: args{
				value: "plain",
			},
			want: want{
				val: "plain",
			},
		},
		{
			name: "ampersand and quote",
			args: args{
				value: `a&"b`,
			},
			want: want{
				val: "a&amp;&quot;b",
			},
		},
		{
			name: "normalizes line endings",
			args: args{
				value: "a\r\nb\rc\n",
			},
			want: want{
				val: "a\nb\nc\n",
			},
		},
		{
			name: "preserves allowed whitespace",
			args: args{
				value: "a\t\n\fb",
			},
			want: want{
				val: "a\t\n\fb",
			},
		},
		{
			name: "replaces controls and malformed utf-8",
			args: args{
				value: "\x00\x01\x7f\u0080\xff",
			},
			want: want{
				val: "�����",
			},
		},
		{
			name: "replaces noncharacters",
			args: args{
				value: "\ufdd0\ufffe\U0001ffff",
			},
			want: want{
				val: "���",
			},
		},
		{
			name: "preserves valid unicode",
			args: args{
				value: "日本語�",
			},
			want: want{
				val: "日本語�",
			},
		},
		{
			name: "preserves valid unicode while escaping",
			args: args{
				value: "日&本",
			},
			want: want{
				val: "日&amp;本",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := escapeAttr(test.args.value)
			testutil.AssertValue(t, got, test.want.val, "escapeAttr")
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
				value: "\t\n\f ",
			},
		},
		{
			name: "attribute delimiters",
			args: args{
				value: `&"`,
			},
			want: want{
				val: true,
			},
		},
		{
			name: "carriage return",
			args: args{
				value: "\r",
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
