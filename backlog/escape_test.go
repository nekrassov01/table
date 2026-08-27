package backlog

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
				value: "plain value",
			},
			want: want{
				value: "plain value",
			},
		},
		{
			name: "escapes cell content",
			fields: fields{
				escapes: []byte("kept:"),
			},
			args: args{
				value: "a\\b|c\r\nd\re\nf",
			},
			want: want{
				value:   `a\\\b\\|c&br;d&br;e&br;f`,
				escapes: `kept:a\\\b\\|c&br;d&br;e&br;f`,
			},
		},
		{
			name: "literalizes inline notation",
			args: args{
				value: `[[link]] ''bold'' %%strike%% {code}code{/code} #contents &br; &color(red){value}`,
			},
			want: want{
				value:   `\\[\\[link\\]\\] \\'\\'bold\\'\\' \\%\\%strike\\%\\% \\{code}code\\{/code} \\#contents \\&br; \\&color(red){value}`,
				escapes: `\\[\\[link\\]\\] \\'\\'bold\\'\\' \\%\\%strike\\%\\% \\{code}code\\{/code} \\#contents \\&br; \\&color(red){value}`,
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
			name: "plain",
			args: args{
				value: "plain value",
			},
			want: want{
				index: -1,
			},
		},
		{
			name: "backslash",
			args: args{
				value: `a\b`,
			},
			want: want{
				index: 1,
			},
		},
		{
			name: "separator",
			args: args{
				value: "a|b",
			},
			want: want{
				index: 1,
			},
		},
		{
			name: "plain punctuation",
			args: args{
				value: "[]'%&{}#",
			},
			want: want{
				index: -1,
			},
		},
		{
			name: "adjacent marker",
			args: args{
				value: "''",
			},
			want: want{
				index: 0,
			},
		},
		{
			name: "line break notation",
			args: args{
				value: "&br;",
			},
			want: want{
				index: 0,
			},
		},
		{
			name: "color notation",
			args: args{
				value: "&color(red){value}",
			},
			want: want{
				index: 0,
			},
		},
		{
			name: "quote notation",
			args: args{
				value: "{quote}",
			},
			want: want{
				index: 0,
			},
		},
		{
			name: "quote close notation",
			args: args{
				value: "{/quote}",
			},
			want: want{
				index: 0,
			},
		},
		{
			name: "code notation",
			args: args{
				value: "{code}",
			},
			want: want{
				index: 0,
			},
		},
		{
			name: "typed code notation",
			args: args{
				value: "{code:go}",
			},
			want: want{
				index: 0,
			},
		},
		{
			name: "code close notation",
			args: args{
				value: "{/code}",
			},
			want: want{
				index: 0,
			},
		},
		{
			name: "attachment notation",
			args: args{
				value: "#attach(file:1)",
			},
			want: want{
				index: 0,
			},
		},
		{
			name: "image notation",
			args: args{
				value: "#image(1)",
			},
			want: want{
				index: 0,
			},
		},
		{
			name: "thumbnail notation",
			args: args{
				value: "#thumbnail(1)",
			},
			want: want{
				index: 0,
			},
		},
		{
			name: "revision notation",
			args: args{
				value: "#rev(1)",
			},
			want: want{
				index: 0,
			},
		},
		{
			name: "contents notation",
			args: args{
				value: "#contents",
			},
			want: want{
				index: 0,
			},
		},
		{
			name: "carriage return",
			args: args{
				value: "a\rb",
			},
			want: want{
				index: 1,
			},
		},
		{
			name: "line feed",
			args: args{
				value: "a\nb",
			},
			want: want{
				index: 1,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := want{
				index: indexEscapeValue(test.args.value),
			}
			testutil.AssertValue(t, got, test.want, "indexEscapeValue")
		})
	}
}
