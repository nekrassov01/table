package column

import (
	"testing"

	"github.com/nekrassov01/table/internal/testutil"
)

func TestNewSelector(t *testing.T) {
	type args struct {
		indexes []int
		mutate  []int
	}
	type want struct {
		indexes     []int
		columnCount int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "no indexes",
		},
		{
			name: "owns indexes",
			args: args{
				indexes: []int{0, -1, 2},
				mutate:  []int{3, 4, 5},
			},
			want: want{
				indexes:     []int{0, -1, 2},
				columnCount: 3,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			selector := NewSelector(test.args.indexes...)
			copy(test.args.indexes, test.args.mutate)
			got := want{
				indexes:     selector.indexes,
				columnCount: selector.columnCount,
			}
			testutil.AssertValue(t, got, test.want, "NewSelector")
		})
	}
}

func TestAll(t *testing.T) {
	type want struct {
		indexes     []int
		columnCount int
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "selects every column",
			want: want{
				columnCount: allColumnCount,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			selector := All()
			got := want{
				indexes:     selector.indexes,
				columnCount: selector.columnCount,
			}
			testutil.AssertValue(t, got, test.want, "All")
		})
	}
}

func TestSet_Apply(t *testing.T) {
	type fields struct {
		values   []int
		defaults *int
	}
	type args struct {
		selector   Selector
		newDefault func() int
		fn         func(*int)
	}
	type want struct {
		values   []int
		defaults *int
	}
	increment := func(value *int) {
		(*value)++
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "all columns creates and updates defaults",
			fields: fields{
				values: []int{2},
			},
			args: args{
				selector: All(),
				newDefault: func() int {
					return 3
				},
				fn: increment,
			},
			want: want{
				values:   []int{3},
				defaults: intPointer(4),
			},
		},
		{
			name: "all columns updates existing defaults",
			fields: fields{
				values:   []int{2},
				defaults: intPointer(5),
			},
			args: args{
				selector: All(),
				fn:       increment,
			},
			want: want{
				values:   []int{3},
				defaults: intPointer(6),
			},
		},
		{
			name: "selected columns inherit defaults",
			fields: fields{
				values:   []int{2},
				defaults: intPointer(5),
			},
			args: args{
				selector: NewSelector(-1, 2),
				fn:       increment,
			},
			want: want{
				values:   []int{2, 5, 6},
				defaults: intPointer(5),
			},
		},
		{
			name: "selected columns use supplied defaults",
			args: args{
				selector: NewSelector(1),
				newDefault: func() int {
					return 3
				},
				fn: increment,
			},
			want: want{
				values: []int{3, 4},
			},
		},
		{
			name: "no selected columns changes nothing",
			args: args{
				selector: NewSelector(),
				fn:       increment,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := Set[int]{
				Values:   test.fields.values,
				Defaults: test.fields.defaults,
			}
			o.Apply(test.args.selector, test.args.newDefault, test.args.fn)
			got := want{
				values:   o.Values,
				defaults: o.Defaults,
			}
			testutil.AssertValue(t, got, test.want, "Apply")
		})
	}
}

func TestMaxColumns(t *testing.T) {
	type args struct {
		rows [][]string
	}
	type want struct {
		count int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "no rows",
		},
		{
			name: "empty rows",
			args: args{
				rows: [][]string{
					{},
					{},
				},
			},
		},
		{
			name: "widest row",
			args: args{
				rows: [][]string{
					{"a"},
					{"a", "b", "c"},
					{"a", "b"},
				},
			},
			want: want{
				count: 3,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := want{
				count: MaxColumns(test.args.rows),
			}
			testutil.AssertValue(t, got, test.want, "MaxColumns")
		})
	}
}

func intPointer(value int) *int {
	return &value
}
