package markdown

import (
	"testing"

	"github.com/nekrassov01/table/internal/testutil"
)

func TestNewDecoration(t *testing.T) {
	type args struct {
		prefix string
		suffix string
	}
	type want struct {
		decoration *Decoration
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "empty prefix",
			args: args{
				suffix: "suffix",
			},
		},
		{
			name: "markers",
			args: args{
				prefix: "<",
				suffix: ">",
			},
			want: want{
				decoration: &Decoration{
					Prefix: "<",
					Suffix: ">",
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := want{
				decoration: NewDecoration(test.args.prefix, test.args.suffix),
			}
			testutil.AssertValue(t, got, test.want, "NewDecoration")
		})
	}
}

func TestResolveTicks(t *testing.T) {
	type args struct {
		decoration *Decoration
		value      string
	}
	type want struct {
		ticks int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "nil decoration",
		},
		{
			name: "non-code decoration",
			args: args{
				decoration: DecorationBold,
				value:      "value",
			},
		},
		{
			name: "configured fence",
			args: args{
				decoration: NewDecoration("``", "``"),
				value:      "value",
			},
			want: want{
				ticks: 2,
			},
		},
		{
			name: "grows beyond content run",
			args: args{
				decoration: DecorationCode,
				value:      "a```b`c",
			},
			want: want{
				ticks: 4,
			},
		},
		{
			name: "different suffix",
			args: args{
				decoration: NewDecoration("`", "``"),
				value:      "value",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := want{
				ticks: resolveTicks(test.args.decoration, test.args.value),
			}
			testutil.AssertValue(t, got, test.want, "resolveTicks")
		})
	}
}
