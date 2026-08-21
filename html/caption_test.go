package html

import (
	"testing"

	"github.com/nekrassov01/table/internal/testutil"
)

func Test_resolveCaptionSide(t *testing.T) {
	type args struct {
		side CaptionSide
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
			name: "default",
		},
		{
			name: "top",
			args: args{
				side: CaptionTop,
			},
			want: want{
				val: "caption-side:top",
			},
		},
		{
			name: "bottom",
			args: args{
				side: CaptionBottom,
			},
			want: want{
				val: "caption-side:bottom",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := resolveCaptionSide(test.args.side)
			testutil.AssertValue(t, got, test.want.val, "resolveCaptionSide")
		})
	}
}
