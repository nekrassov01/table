package table

import (
	"iter"
	"testing"

	"github.com/nekrassov01/table/internal/testutil"
)

func TestTableOf(t *testing.T) {
	type args struct {
		values []entry
	}
	type want struct {
		rows [][]Value
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "values",
			args: args{
				values: []entry{
					{
						name: "a",
						size: 1,
					},
					{
						name: "b",
						size: 2,
					},
				},
			},
			want: want{
				rows: [][]Value{{String("a"), Int(1)}, {String("b"), Int(2)}},
			},
		},
		{
			name: "no values",
			args: args{
				values: nil,
			},
			want: want{
				rows: [][]Value{},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := TableOf(test.args.values, entryRow)
			testutil.AssertValue(t, got, test.want.rows, "TableOf")
		})
	}
}

func TestStreamOf(t *testing.T) {
	type args struct {
		values iter.Seq2[entry, error]
		stop   int
	}
	type want struct {
		rows  [][]Value
		isErr bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "values",
			args: args{
				values: testutil.Seq2([]entry{
					{
						name: "a",
						size: 0,
					},
					{
						name: "b",
						size: 1,
					},
				}, nil),
			},
			want: want{
				rows: [][]Value{{String("a"), Int(0)}, {String("b"), Int(1)}},
			},
		},
		{
			name: "error ends the sequence",
			args: args{
				values: testutil.Seq2([]entry{
					{
						name: "a",
						size: 0,
					},
				}, testutil.NewError()),
			},
			want: want{
				rows:  [][]Value{{String("a"), Int(0)}},
				isErr: true,
			},
		},
		{
			name: "consumer stops early",
			args: args{
				values: testutil.Seq2([]entry{
					{
						name: "a",
						size: 0,
					},
					{
						name: "b",
						size: 1,
					},
					{
						name: "c",
						size: 2,
					},
				}, nil),
				stop: 1,
			},
			want: want{
				rows: [][]Value{{String("a"), Int(0)}},
			},
		},
		{
			name: "no values",
			args: args{
				values: testutil.Seq2[entry](nil, nil),
			},
			want: want{
				rows: nil,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var rows [][]Value
			var err error
			for row, e := range StreamOf(test.args.values, entryRow) {
				if e != nil {
					err = e
					break
				}
				rows = append(rows, row)
				if test.args.stop != 0 && len(rows) >= test.args.stop {
					break
				}
			}
			testutil.AssertValue(t, rows, test.want.rows, "StreamOf")
			testutil.AssertValue(t, err != nil, test.want.isErr, "StreamOf error")
		})
	}
}

type entry struct {
	name string
	size int
}

func entryRow(e entry) []Value {
	return []Value{String(e.name), Int(e.size)}
}
