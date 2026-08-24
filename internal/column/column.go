// Package column provides shared operations for resolving logical columns.
package column

import "slices"

// allColumnCount distinguishes All without enlarging selectors captured by
// option closures.
const allColumnCount = -1

// Selector identifies explicit column indexes or every column.
type Selector struct {
	indexes     []int
	columnCount int
}

// NewSelector returns a selector for indexes and owns a copy of them.
func NewSelector(indexes ...int) Selector {
	columnCount := 0
	for _, index := range indexes {
		columnCount = max(columnCount, index+1)
	}
	return Selector{
		indexes:     slices.Clone(indexes),
		columnCount: columnCount,
	}
}

// All returns a selector for every column, including columns discovered later.
func All() Selector {
	return Selector{
		columnCount: allColumnCount,
	}
}

// Set holds explicit column settings and the default applied to every column.
type Set[T any] struct {
	Values   []T // Explicit settings in column order.
	Defaults *T  // Settings inherited by every column.
}

// Apply updates every column selected by selector. New explicit columns inherit
// Defaults when present. Otherwise they use the value returned by newDefault,
// or the zero value of T when newDefault is nil.
func (o *Set[T]) Apply(selector Selector, newDefault func() T, fn func(*T)) {
	if selector.columnCount == allColumnCount {
		if o.Defaults == nil {
			var defaults T
			if newDefault != nil {
				defaults = newDefault()
			}
			o.Defaults = &defaults
		}
		fn(o.Defaults)
		for index := range o.Values {
			fn(&o.Values[index])
		}
		return
	}
	existing := len(o.Values)
	if selector.columnCount > existing {
		o.Values = append(o.Values, make([]T, selector.columnCount-existing)...)
		var defaults T
		if newDefault != nil {
			defaults = newDefault()
		}
		if o.Defaults != nil {
			defaults = *o.Defaults
		}
		for index := existing; index < selector.columnCount; index++ {
			o.Values[index] = defaults
		}
	}
	for _, index := range selector.indexes {
		if index >= 0 {
			fn(&o.Values[index])
		}
	}
}

// MaxColumns returns the greatest column count among rows.
func MaxColumns(rows [][]string) int {
	columnCount := 0
	for _, row := range rows {
		columnCount = max(columnCount, len(row))
	}
	return columnCount
}
