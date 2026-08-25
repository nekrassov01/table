package column

import (
	"math"
	"testing"

	"github.com/nekrassov01/table/internal/testutil"
)

func TestNewSelector(t *testing.T) {
	type args struct {
		indexes []int
		mutate  []int
	}
	type want struct {
		indexes []int
		all     bool
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
				indexes: []int{0, 2},
			},
		},
		{
			name: "sorts unique indexes without deriving storage",
			args: args{
				indexes: []int{math.MaxInt, 2, math.MaxInt, 1},
			},
			want: want{
				indexes: []int{1, 2, math.MaxInt},
			},
		},
		{
			name: "compacts sorted duplicate indexes",
			args: args{
				indexes: []int{1, 1, 2},
			},
			want: want{
				indexes: []int{1, 2},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			selector := NewSelector(test.args.indexes...)
			copy(test.args.indexes, test.args.mutate)
			got := want{
				indexes: selector.indexes,
				all:     selector.all,
			}
			testutil.AssertValue(t, got, test.want, "NewSelector")
		})
	}
}

func TestAll(t *testing.T) {
	type want struct {
		indexes []int
		all     bool
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "selects every column",
			want: want{
				all: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			selector := All()
			got := want{
				indexes: selector.indexes,
				all:     selector.all,
			}
			testutil.AssertValue(t, got, test.want, "All")
		})
	}
}

func TestSelector_denseEnd(t *testing.T) {
	type fields struct {
		indexes []int
	}
	type args struct {
		columnCount int
	}
	type want struct {
		end int
		ok  bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "existing columns",
			fields: fields{
				indexes: []int{0, 1},
			},
			args: args{
				columnCount: 2,
			},
			want: want{
				end: 2,
				ok:  true,
			},
		},
		{
			name: "bounded gaps",
			fields: fields{
				indexes: []int{2, 4},
			},
			args: args{
				columnCount: 2,
			},
			want: want{
				end: 5,
				ok:  true,
			},
		},
		{
			name: "unbounded gaps",
			fields: fields{
				indexes: []int{20},
			},
			want: want{},
		},
		{
			name: "maximum index",
			fields: fields{
				indexes: []int{math.MaxInt},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			selector := Selector{indexes: test.fields.indexes}
			end, ok := selector.denseEnd(test.args.columnCount)
			got := want{
				end: end,
				ok:  ok,
			}
			testutil.AssertValue(t, got, test.want, "denseEnd")
		})
	}
}

func TestSet_Default(t *testing.T) {
	type fields struct {
		state *setState[int]
	}
	type want struct {
		defaults *int
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "unset defaults",
		},
		{
			name: "sparse settings without defaults",
			fields: fields{
				state: setStateOf(nil, value[int]{index: 3, config: 4}),
			},
		},
		{
			name: "configured defaults",
			fields: fields{
				state: setStateOf(new(3)),
			},
			want: want{
				defaults: new(3),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := Set[int]{
				state: test.fields.state,
			}
			got := want{
				defaults: o.Default(),
			}
			testutil.AssertValue(t, got, test.want, "Default")
		})
	}
}

func TestSet_Apply(t *testing.T) {
	type fields struct {
		values   []int
		defaults *int
		sparse   *setState[int]
	}
	type args struct {
		selector      Selector
		defaultValue  int
		hasNewDefault bool
		fn            func(*int)
	}
	type want struct {
		values       []int
		defaults     *int
		sparse       []value[int]
		defaultCalls int
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
				selector:      All(),
				defaultValue:  3,
				hasNewDefault: true,
				fn:            increment,
			},
			want: want{
				values:       []int{3},
				defaults:     new(4),
				defaultCalls: 1,
			},
		},
		{
			name: "all columns updates existing defaults",
			fields: fields{
				values:   []int{2},
				defaults: new(5),
			},
			args: args{
				selector: All(),
				fn:       increment,
			},
			want: want{
				values:   []int{3},
				defaults: new(6),
			},
		},
		{
			name: "selected columns inherit defaults",
			fields: fields{
				values:   []int{2},
				defaults: new(5),
			},
			args: args{
				selector:      NewSelector(-1, 2),
				defaultValue:  9,
				hasNewDefault: true,
				fn:            increment,
			},
			want: want{
				values:   []int{2, 5, 6},
				defaults: new(5),
			},
		},
		{
			name: "distant columns inherit defaults",
			fields: fields{
				values:   []int{2},
				defaults: new(5),
			},
			args: args{
				selector: NewSelector(12),
				fn:       increment,
			},
			want: want{
				values:   []int{2},
				defaults: new(5),
				sparse:   []value[int]{{index: 12, config: 6}},
			},
		},
		{
			name: "selected columns use supplied defaults",
			args: args{
				selector:      NewSelector(1),
				defaultValue:  3,
				hasNewDefault: true,
				fn:            increment,
			},
			want: want{
				values:       []int{3, 4},
				defaultCalls: 1,
			},
		},
		{
			name: "selected columns remain sparse and ordered",
			args: args{
				selector:      NewSelector(math.MaxInt, 1, math.MaxInt),
				defaultValue:  3,
				hasNewDefault: true,
				fn:            increment,
			},
			want: want{
				sparse: []value[int]{
					{index: 1, config: 4},
					{index: math.MaxInt, config: 4},
				},
				defaultCalls: 1,
			},
		},
		{
			name: "selected columns fill a dense prefix",
			args: args{
				selector:      NewSelector(2, 0, 1),
				defaultValue:  3,
				hasNewDefault: true,
				fn:            increment,
			},
			want: want{
				values:       []int{4, 4, 4},
				defaultCalls: 1,
			},
		},
		{
			name: "selected columns retain a sparse value after extending prefix",
			args: args{
				selector:      NewSelector(0, 20),
				defaultValue:  3,
				hasNewDefault: true,
				fn:            increment,
			},
			want: want{
				sparse: []value[int]{
					{index: 0, config: 4},
					{index: 20, config: 4},
				},
				defaultCalls: 1,
			},
		},
		{
			name: "selected columns extend through an existing sparse value",
			fields: fields{
				sparse: setStateOf(nil, value[int]{index: 1, config: 5}),
			},
			args: args{
				selector:      NewSelector(0),
				defaultValue:  3,
				hasNewDefault: true,
				fn:            increment,
			},
			want: want{
				sparse: []value[int]{
					{index: 0, config: 4},
					{index: 1, config: 5},
				},
				defaultCalls: 1,
			},
		},
		{
			name: "selected columns update existing dense and sparse values",
			fields: fields{
				values: []int{2},
				sparse: setStateOf(nil, value[int]{index: 2, config: 5}),
			},
			args: args{
				selector: NewSelector(0, 2),
				fn:       increment,
			},
			want: want{
				values: []int{3},
				sparse: []value[int]{{index: 2, config: 6}},
			},
		},
		{
			name: "selected columns add to existing sparse values",
			fields: fields{
				sparse: setStateOf(nil, value[int]{index: 1, config: 5}),
			},
			args: args{
				selector:      NewSelector(3),
				defaultValue:  3,
				hasNewDefault: true,
				fn:            increment,
			},
			want: want{
				sparse: []value[int]{
					{index: 1, config: 5},
					{index: 3, config: 4},
				},
				defaultCalls: 1,
			},
		},
		{
			name: "selected columns update an existing distant value",
			fields: fields{
				sparse: setStateOf(nil, value[int]{index: 3, config: 5}),
			},
			args: args{
				selector: NewSelector(3),
				fn:       increment,
			},
			want: want{
				sparse: []value[int]{{index: 3, config: 6}},
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
			defaultCalls := 0
			var newDefault func() int
			if test.args.hasNewDefault {
				newDefault = func() int {
					defaultCalls++
					return test.args.defaultValue
				}
			}
			o := Set[int]{
				Values: test.fields.values,
				state:  mergeSetState(test.fields.defaults, test.fields.sparse),
			}
			o.Apply(test.args.selector, newDefault, test.args.fn)
			var sparse []value[int]
			if o.state != nil {
				sparse = o.state.values
			}
			got := want{
				values:       o.Values,
				defaults:     setDefaults(o.state),
				sparse:       sparse,
				defaultCalls: defaultCalls,
			}
			testutil.AssertValue(t, got, test.want, "Apply")
		})
	}
}

func TestSet_applyAll(t *testing.T) {
	type fields struct {
		values   []int
		defaults *int
	}
	type args struct {
		defaultValue int
		fn           func(*int)
	}
	type want struct {
		values   []int
		defaults *int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "creates defaults and updates every value",
			fields: fields{
				values: []int{1},
			},
			args: args{
				defaultValue: 2,
				fn: func(value *int) {
					(*value)++
				},
			},
			want: want{
				values:   []int{2},
				defaults: new(3),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := Set[int]{
				Values: test.fields.values,
				state:  mergeSetState(test.fields.defaults, nil),
			}
			o.applyAll(func() int {
				return test.args.defaultValue
			}, test.args.fn)
			got := want{
				values:   o.Values,
				defaults: setDefaults(o.state),
			}
			testutil.AssertValue(t, got, test.want, "applyAll")
		})
	}
}

func TestSet_applyDense(t *testing.T) {
	type fields struct {
		values   []int
		defaults *int
	}
	type args struct {
		selector     Selector
		denseEnd     int
		defaultValue int
		fn           func(*int)
	}
	type want struct {
		values []int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "extends defaults and updates selected values",
			fields: fields{
				values: []int{1},
			},
			args: args{
				selector:     NewSelector(2),
				denseEnd:     3,
				defaultValue: 2,
				fn: func(value *int) {
					(*value)++
				},
			},
			want: want{
				values: []int{1, 2, 3},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := Set[int]{
				Values: test.fields.values,
				state:  mergeSetState(test.fields.defaults, nil),
			}
			o.applyDense(test.args.selector, test.args.denseEnd, func() int {
				return test.args.defaultValue
			}, test.args.fn)
			got := want{
				values: o.Values,
			}
			testutil.AssertValue(t, got, test.want, "applyDense")
		})
	}
}

func TestSet_applySparse(t *testing.T) {
	type fields struct {
		values   []int
		defaults *int
		sparse   *setState[int]
	}
	type args struct {
		selector     Selector
		defaultValue int
		fn           func(*int)
	}
	type want struct {
		values []int
		sparse []value[int]
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "updates dense and sparse values",
			fields: fields{
				values: []int{1},
				sparse: setStateOf(nil, value[int]{index: 3, config: 3}),
			},
			args: args{
				selector:     NewSelector(0, 3, 5),
				defaultValue: 2,
				fn: func(value *int) {
					(*value)++
				},
			},
			want: want{
				values: []int{2},
				sparse: []value[int]{
					{index: 3, config: 4},
					{index: 5, config: 3},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := Set[int]{
				Values: test.fields.values,
				state:  mergeSetState(test.fields.defaults, test.fields.sparse),
			}
			o.applySparse(test.args.selector, func() int {
				return test.args.defaultValue
			}, test.args.fn)
			var sparse []value[int]
			if o.state != nil {
				sparse = o.state.values
			}
			got := want{
				values: o.Values,
				sparse: sparse,
			}
			testutil.AssertValue(t, got, test.want, "applySparse")
		})
	}
}

func TestSet_Range(t *testing.T) {
	type fields struct {
		values []int
		sparse *setState[int]
	}
	type want struct {
		values []int
		sparse []value[int]
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "no values",
		},
		{
			name: "dense and sparse values",
			fields: fields{
				values: []int{1, 2},
				sparse: setStateOf(nil, value[int]{index: 4, config: 3}),
			},
			want: want{
				values: []int{2, 3},
				sparse: []value[int]{{index: 4, config: 4}},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := Set[int]{
				Values: test.fields.values,
				state:  test.fields.sparse,
			}
			o.Range(func(value *int) {
				(*value)++
			})
			var sparse []value[int]
			if o.state != nil {
				sparse = o.state.values
			}
			got := want{
				values: o.Values,
				sparse: sparse,
			}
			testutil.AssertValue(t, got, test.want, "Range")
		})
	}
}

func TestSet_Any(t *testing.T) {
	type fields struct {
		values []int
		sparse *setState[int]
	}
	type args struct {
		value int
	}
	type want struct {
		matched bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "no match",
			fields: fields{
				values: []int{1},
			},
			args: args{
				value: 2,
			},
		},
		{
			name: "dense match",
			fields: fields{
				values: []int{1, 2},
				sparse: setStateOf(nil, value[int]{index: 4, config: 3}),
			},
			args: args{
				value: 2,
			},
			want: want{
				matched: true,
			},
		},
		{
			name: "sparse match",
			fields: fields{
				values: []int{1},
				sparse: setStateOf(nil, value[int]{index: 4, config: 3}),
			},
			args: args{
				value: 3,
			},
			want: want{
				matched: true,
			},
		},
		{
			name: "sparse miss",
			fields: fields{
				values: []int{1},
				sparse: setStateOf(nil, value[int]{index: 4, config: 3}),
			},
			args: args{
				value: 4,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := Set[int]{
				Values: test.fields.values,
				state:  test.fields.sparse,
			}
			got := want{
				matched: o.Any(func(value *int) bool {
					return *value == test.args.value
				}),
			}
			testutil.AssertValue(t, got, test.want, "Any")
		})
	}
}

func TestSet_Resolve(t *testing.T) {
	type fields struct {
		values   []int
		defaults *int
		sparse   *setState[int]
	}
	type args struct {
		columns     []int
		columnCount int
		indexOffset int
		defaults    int
	}
	type want struct {
		columns []int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "resolves index defaults and explicit input settings",
			fields: fields{
				values:   []int{7},
				defaults: new(5),
				sparse:   setStateOf(nil, value[int]{index: 2, config: 8}),
			},
			args: args{
				columns:     make([]int, 1, 4),
				columnCount: 4,
				indexOffset: 1,
				defaults:    3,
			},
			want: want{
				columns: []int{3, 7, 5, 8},
			},
		},
		{
			name: "ignores settings beyond resolved input columns",
			fields: fields{
				sparse: setStateOf(nil, value[int]{index: math.MaxInt, config: 9}),
			},
			args: args{
				columnCount: 2,
				defaults:    3,
			},
			want: want{
				columns: []int{3, 3},
			},
		},
		{
			name: "configured defaults override supplied input defaults",
			fields: fields{
				defaults: new(5),
			},
			args: args{
				columnCount: 2,
				defaults:    3,
			},
			want: want{
				columns: []int{5, 5},
			},
		},
		{
			name: "supplied defaults without settings",
			args: args{
				columnCount: 2,
				defaults:    3,
			},
			want: want{
				columns: []int{3, 3},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := Set[int]{
				Values: test.fields.values,
				state:  mergeSetState(test.fields.defaults, test.fields.sparse),
			}
			got := want{
				columns: o.Resolve(test.args.columns, test.args.columnCount, test.args.indexOffset, test.args.defaults),
			}
			testutil.AssertValue(t, got, test.want, "Resolve")
		})
	}
}

func TestSet_resolveValues(t *testing.T) {
	type fields struct {
		values []int
	}
	type args struct {
		columns     []int
		columnCount int
		indexOffset int
		defaults    int
	}
	type want struct {
		columns []int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "resolves contiguous settings",
			fields: fields{
				values: []int{7},
			},
			args: args{
				columns:     make([]int, 1, 3),
				columnCount: 3,
				indexOffset: 1,
				defaults:    3,
			},
			want: want{
				columns: []int{3, 7, 3},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := Set[int]{
				Values: test.fields.values,
			}
			got := want{
				columns: o.resolveValues(test.args.columns, test.args.columnCount, test.args.indexOffset, test.args.defaults),
			}
			testutil.AssertValue(t, got, test.want, "resolveValues")
		})
	}
}

func TestSet_resolveState(t *testing.T) {
	type fields struct {
		values []int
		state  *setState[int]
	}
	type args struct {
		columns     []int
		columnCount int
		indexOffset int
		defaults    int
	}
	type want struct {
		columns []int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "resolves defaults and sparse settings",
			fields: fields{
				values: []int{7},
				state:  setStateOf(new(5), value[int]{index: 2, config: 8}),
			},
			args: args{
				columnCount: 4,
				indexOffset: 1,
				defaults:    3,
			},
			want: want{
				columns: []int{3, 7, 5, 8},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := Set[int]{
				Values: test.fields.values,
				state:  test.fields.state,
			}
			got := want{
				columns: o.resolveState(test.args.columns, test.args.columnCount, test.args.indexOffset, test.args.defaults),
			}
			testutil.AssertValue(t, got, test.want, "resolveState")
		})
	}
}

func Test_findValue(t *testing.T) {
	type args struct {
		values []value[int]
		index  int
	}
	type want struct {
		position int
		found    bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "empty values",
		},
		{
			name: "existing value",
			args: args{
				values: []value[int]{{index: 1}, {index: 3}},
				index:  3,
			},
			want: want{
				position: 1,
				found:    true,
			},
		},
		{
			name: "insertion position",
			args: args{
				values: []value[int]{{index: 1}, {index: 3}},
				index:  2,
			},
			want: want{
				position: 1,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			position, found := findValue(test.args.values, test.args.index)
			got := want{
				position: position,
				found:    found,
			}
			testutil.AssertValue(t, got, test.want, "findValue")
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

func setStateOf(defaults *int, values ...value[int]) *setState[int] {
	state := &setState[int]{
		values: values,
	}
	if defaults != nil {
		state.defaults = *defaults
		state.hasDefaults = true
	}
	return state
}

func mergeSetState(defaults *int, state *setState[int]) *setState[int] {
	if state == nil && defaults == nil {
		return nil
	}
	if state == nil {
		state = &setState[int]{}
	}
	if defaults != nil {
		state.defaults = *defaults
		state.hasDefaults = true
	}
	return state
}

func setDefaults(state *setState[int]) *int {
	if state == nil || !state.hasDefaults {
		return nil
	}
	return &state.defaults
}
