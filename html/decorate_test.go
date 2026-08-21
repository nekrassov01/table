package html

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
		val *Decoration
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "empty prefix",
		},
		{
			name: "markers",
			args: args{
				prefix: "<mark>",
				suffix: "</mark>",
			},
			want: want{
				val: &Decoration{
					Prefix: "<mark>",
					Suffix: "</mark>",
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := NewDecoration(test.args.prefix, test.args.suffix)
			testutil.AssertValue(t, got, test.want.val, "NewDecoration")
		})
	}
}
