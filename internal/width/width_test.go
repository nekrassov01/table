package width

import (
	"testing"

	"github.com/nekrassov01/table/internal/testutil"
)

func TestScanner_Next(t *testing.T) {
	type fields struct {
		text string
	}
	type want struct {
		units [][3]int
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "empty",
			want: want{},
		},
		{
			name: "ASCII",
			fields: fields{
				text: "a\tb",
			},
			want: want{
				units: [][3]int{
					{0, 1, 1},
					{1, 2, 0},
					{2, 3, 1},
				},
			},
		},
		{
			name: "grapheme clusters",
			fields: fields{
				text: "a👩‍💻b",
			},
			want: want{
				units: [][3]int{
					{0, 1, 1},
					{1, 12, 2},
					{12, 13, 1},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := NewScanner(test.fields.text)
			var units [][3]int
			for start, end, displayWidth, ok := o.Next(); ok; start, end, displayWidth, ok = o.Next() {
				units = append(units, [3]int{start, end, displayWidth})
			}
			got := want{
				units: units,
			}
			testutil.AssertValue(t, got, test.want, "Next")
		})
	}
}

func TestStringWidth(t *testing.T) {
	type args struct {
		s string
	}
	type want struct {
		val int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "empty",
			args: args{
				s: "",
			},
			want: want{
				val: 0,
			},
		},
		{
			name: "ascii",
			args: args{
				s: "hello",
			},
			want: want{
				val: 5,
			},
		},
		{
			name: "single-fullwidth",
			args: args{
				s: "あ",
			},
			want: want{
				val: 2,
			},
		},
		{
			name: "mixed-ascii-and-fullwidth",
			args: args{
				s: "aあb",
			},
			want: want{
				val: 4,
			},
		},
		{
			name: "multiple-fullwidth",
			args: args{
				s: "あいう",
			},
			want: want{
				val: 6,
			},
		},
		{
			name: "control-in-ascii",
			args: args{
				s: "a\tb",
			},
			want: want{
				val: 2,
			},
		},
		{
			name: "control-with-fullwidth",
			args: args{
				s: "あ\tb",
			},
			want: want{
				val: 3,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := StringWidth(test.args.s)
			testutil.AssertValue(t, got, test.want.val, "StringWidth")
		})
	}
}
