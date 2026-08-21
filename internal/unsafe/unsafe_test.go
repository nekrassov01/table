package unsafe

import (
	"testing"

	"github.com/nekrassov01/table/internal/testutil"
)

func TestView(t *testing.T) {
	type args struct {
		b []byte
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
			name: "bytes",
			args: args{
				b: []byte("abc"),
			},
			want: want{
				val: "abc",
			},
		},
		{
			name: "empty",
			args: args{
				b: nil,
			},
			want: want{
				val: "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testutil.AssertValue(t, View(test.args.b), test.want.val, "View")
		})
	}
}
