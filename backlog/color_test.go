package backlog

import (
	"testing"

	"github.com/nekrassov01/table/internal/testutil"
)

func TestNewColor(t *testing.T) {
	type args struct {
		fg string
		bg string
	}
	type want struct {
		val *Color
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "both empty returns nil",
			args: args{
				fg: "",
				bg: "",
			},
			want: want{
				val: nil,
			},
		},
		{
			name: "closing paren is rejected",
			args: args{
				fg: "red){injected}&color(blue",
				bg: "",
			},
			want: want{
				val: nil,
			},
		},
		{
			name: "comma is rejected",
			args: args{
				fg: "red, green",
				bg: "",
			},
			want: want{
				val: nil,
			},
		},
		{
			name: "unsafe background is rejected",
			args: args{
				fg: "red",
				bg: "blue){x",
			},
			want: want{
				val: nil,
			},
		},
		{
			name: "cell separator is rejected",
			args: args{
				fg: "red|x",
				bg: "",
			},
			want: want{
				val: nil,
			},
		},
		{
			name: "line break is rejected",
			args: args{
				fg: "red\nx",
				bg: "",
			},
			want: want{
				val: nil,
			},
		},
		{
			name: "foreground only",
			args: args{
				fg: "red",
				bg: "",
			},
			want: want{
				val: &Color{
					Prefix: "&color(red){",
					Suffix: "}",
				},
			},
		},
		{
			name: "background only",
			args: args{
				fg: "",
				bg: "red",
			},
			want: want{
				val: &Color{
					Prefix: `&color("", red){`,
					Suffix: "}",
				},
			},
		},
		{
			name: "foreground and background",
			args: args{
				fg: "red",
				bg: "blue",
			},
			want: want{
				val: &Color{
					Prefix: "&color(red, blue){",
					Suffix: "}",
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := NewColor(test.args.fg, test.args.bg)
			testutil.AssertValue(t, got, test.want.val, "NewColor")
		})
	}
}
