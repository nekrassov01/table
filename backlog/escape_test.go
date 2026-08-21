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
		escape bool
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
		},
		{
			name: "backslash",
			args: args{
				value: `a\b`,
			},
			want: want{
				escape: true,
			},
		},
		{
			name: "separator",
			args: args{
				value: "a|b",
			},
			want: want{
				escape: true,
			},
		},
		{
			name: "carriage return",
			args: args{
				value: "a\rb",
			},
			want: want{
				escape: true,
			},
		},
		{
			name: "line feed",
			args: args{
				value: "a\nb",
			},
			want: want{
				escape: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := want{
				escape: needsEscapeValue(test.args.value),
			}
			testutil.AssertValue(t, got, test.want, "needsEscape")
		})
	}
}
