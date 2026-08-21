package span

import (
	"testing"

	"github.com/nekrassov01/table/internal/testutil"
)

func TestPreviousRow_Reset(t *testing.T) {
	type fields struct {
		values  [][]byte
		started bool
	}
	type want struct {
		retainsCapacity bool
		started         bool
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "retains backing without continuing the previous part",
			fields: fields{
				values:  [][]byte{make([]byte, 6, 8)},
				started: true,
			},
			want: want{
				retainsCapacity: true,
				started:         false,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &PreviousRow{
				values:  test.fields.values,
				started: test.fields.started,
			}
			o.Reset()
			gotCapacity := cap(o.values[0]) == cap(test.fields.values[0])
			testutil.AssertValue(t, gotCapacity, test.want.retainsCapacity, "entry capacity")
			testutil.AssertValue(t, o.started, test.want.started, "started")
		})
	}
}

func TestRowspans(t *testing.T) {
	type args struct {
		rowspan uint64
		rows    [][]string
	}
	type want struct {
		bits []uint64
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "no spanning positions",
			args: args{
				rowspan: 0,
				rows:    [][]string{{"a"}, {"a"}},
			},
			want: want{
				bits: []uint64{0, 0},
			},
		},
		{
			name: "first row never spans",
			args: args{
				rowspan: 1,
				rows:    [][]string{{"a"}},
			},
			want: want{
				bits: []uint64{0},
			},
		},
		{
			name: "single column groups and breaks",
			args: args{
				rowspan: 1,
				rows:    [][]string{{"a"}, {"a"}, {"b"}, {"b"}},
			},
			want: want{
				bits: []uint64{0, 1, 0, 1},
			},
		},
		{
			name: "a break on the left breaks the right",
			args: args{
				rowspan: 0b11,
				rows:    [][]string{{"a", "x"}, {"a", "x"}, {"b", "x"}},
			},
			want: want{
				bits: []uint64{0, 0b11, 0},
			},
		},
		{
			name: "a break on the right keeps the left",
			args: args{
				rowspan: 0b11,
				rows:    [][]string{{"a", "x"}, {"a", "y"}},
			},
			want: want{
				bits: []uint64{0, 0b01},
			},
		},
		{
			name: "the previous row refreshes even after a break",
			args: args{
				rowspan: 0b11,
				rows:    [][]string{{"a", "x"}, {"b", "x"}, {"b", "x"}},
			},
			want: want{
				bits: []uint64{0, 0, 0b11},
			},
		},
		{
			name: "non-grouping positions are ignored",
			args: args{
				rowspan: 0b10,
				rows:    [][]string{{"a", "x"}, {"b", "x"}},
			},
			want: want{
				bits: []uint64{0, 0b10},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var previous PreviousRow
			previous.Reset()
			got := make([]uint64, 0, len(test.args.rows))
			for _, row := range test.args.rows {
				got = append(got, Rowspans(test.args.rowspan, row, &previous))
			}
			testutil.AssertValue(t, got, test.want.bits, "Rows")
		})
	}
}

func TestColspans(t *testing.T) {
	type args struct {
		colspan uint64
		values  []string
		taken   uint64
	}
	type want struct {
		bits uint64
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "no spanning positions",
			args: args{
				colspan: 0,
				values:  []string{"a", "a"},
			},
			want: want{
				bits: 0,
			},
		},
		{
			name: "adjacent equals absorb rightward",
			args: args{
				colspan: 0b11,
				values:  []string{"a", "a"},
			},
			want: want{
				bits: 0b10,
			},
		},
		{
			name: "runs chain",
			args: args{
				colspan: 0b111,
				values:  []string{"a", "a", "a"},
			},
			want: want{
				bits: 0b110,
			},
		},
		{
			name: "unequal text stays",
			args: args{
				colspan: 0b11,
				values:  []string{"a", "b"},
			},
			want: want{
				bits: 0,
			},
		},
		{
			name: "a gap in the positions blocks the pair",
			args: args{
				colspan: 0b101,
				values:  []string{"a", "a", "a"},
			},
			want: want{
				bits: 0,
			},
		},
		{
			name: "an excluded position never merges",
			args: args{
				colspan: 0b11,
				values:  []string{"a", "a"},
				taken:   0b10,
			},
			want: want{
				bits: 0,
			},
		},
		{
			name: "empty text merges like any text",
			args: args{
				colspan: 0b11,
				values:  []string{"", ""},
			},
			want: want{
				bits: 0b10,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Colspans(test.args.colspan, test.args.values, test.args.taken)
			testutil.AssertValue(t, got, test.want.bits, "Cols")
		})
	}
}
