package markdown

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
			name: "foreground only",
			args: args{
				fg: "red",
				bg: "",
			},
			want: want{
				val: &Color{
					Prefix: `<span style="color:red">`,
					Suffix: `</span>`,
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
					Prefix: `<span style="background-color:red">`,
					Suffix: `</span>`,
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
					Prefix: `<span style="color:red;background-color:blue">`,
					Suffix: `</span>`,
				},
			},
		},
		{
			name: "normalizes color values",
			args: args{
				fg: "r\r\n\x00ed",
				bg: "\u0080blue",
			},
			want: want{
				val: &Color{
					Prefix: `<span style="color:r �ed;background-color:�blue">`,
					Suffix: `</span>`,
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
