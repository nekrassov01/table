// Package column provides shared operations for resolving logical columns.
package column

import (
	"math"
	"slices"
)

// sparseGapThreshold keeps nearby settings dense without following distant indexes.
const sparseGapThreshold = 8

// Selector identifies explicit column indexes or every column.
type Selector struct {
	indexes []int
	all     bool
}

// NewSelector returns an owned selector for the nonnegative indexes.
func NewSelector(indexes ...int) Selector {
	indexes = slices.Clone(indexes)
	selected := indexes[:0]
	sorted := true
	unique := true
	for _, index := range indexes {
		if index < 0 {
			continue
		}
		if len(selected) > 0 {
			previous := selected[len(selected)-1]
			sorted = sorted && previous <= index
			unique = unique && previous != index
		}
		selected = append(selected, index)
	}
	if len(selected) == 0 {
		return Selector{}
	}
	if !sorted {
		slices.Sort(selected)
		selected = slices.Compact(selected)
		return Selector{indexes: selected}
	}
	if !unique {
		selected = slices.Compact(selected)
	}
	return Selector{indexes: selected}
}

// All returns a selector for every column, including columns discovered later.
func All() Selector {
	return Selector{all: true}
}

// denseEnd returns a dense prefix when its gaps cost less than sparse storage.
func (o Selector) denseEnd(columnCount int) (int, bool) {
	firstNew := 0
	for firstNew < len(o.indexes) && o.indexes[firstNew] < columnCount {
		firstNew++
	}
	newCount := len(o.indexes) - firstNew
	if newCount == 0 {
		return columnCount, true
	}
	last := o.indexes[len(o.indexes)-1]
	if last == math.MaxInt {
		return 0, false
	}
	holes := last - columnCount - newCount + 1
	if holes >= sparseGapThreshold && holes > newCount {
		return 0, false
	}
	return last + 1, true
}

// Set holds explicit column settings and the default applied to every column.
type Set[T any] struct {
	Values []T // Explicit settings for the contiguous input-column prefix.
	state  *setState[T]
}

// Default returns settings inherited by every input column, or nil when unset.
func (o *Set[T]) Default() *T {
	if o.state == nil || !o.state.hasDefaults {
		return nil
	}
	return &o.state.defaults
}

// Apply updates every column selected by selector. New explicit columns inherit
// the configured default when present. Otherwise they use the value returned by
// newDefault, or the zero value of T when newDefault is nil.
func (o *Set[T]) Apply(selector Selector, newDefault func() T, fn func(*T)) {
	if selector.all {
		o.applyAll(newDefault, fn)
		return
	}
	if len(selector.indexes) == 0 {
		return
	}
	if o.state == nil || len(o.state.values) == 0 {
		if denseEnd, ok := selector.denseEnd(len(o.Values)); ok {
			o.applyDense(selector, denseEnd, newDefault, fn)
			return
		}
	}
	o.applySparse(selector, newDefault, fn)
}

// applyAll updates the defaults and every explicit column setting.
func (o *Set[T]) applyAll(newDefault func() T, fn func(*T)) {
	if o.state == nil {
		o.state = &setState[T]{}
	}
	if !o.state.hasDefaults {
		if newDefault != nil {
			o.state.defaults = newDefault()
		}
		o.state.hasDefaults = true
	}
	fn(&o.state.defaults)
	o.Range(fn)
}

// applyDense updates selected columns in a bounded dense prefix.
func (o *Set[T]) applyDense(selector Selector, denseEnd int, newDefault func() T, fn func(*T)) {
	existing := len(o.Values)
	if denseEnd > existing {
		var config T
		if o.state != nil && o.state.hasDefaults {
			config = o.state.defaults
		} else if newDefault != nil {
			config = newDefault()
		}
		o.Values = append(o.Values, make([]T, denseEnd-existing)...)
		for index := existing; index < denseEnd; index++ {
			o.Values[index] = config
		}
	}
	for _, selected := range selector.indexes {
		fn(&o.Values[selected])
	}
}

// applySparse updates selected columns while retaining distant settings.
func (o *Set[T]) applySparse(selector Selector, newDefault func() T, fn func(*T)) {
	missing := 0
	for _, selected := range selector.indexes {
		if selected < len(o.Values) {
			continue
		}
		if o.state == nil {
			missing++
			continue
		}
		if _, configured := findValue(o.state.values, selected); !configured {
			missing++
		}
	}
	var config T
	if missing > 0 {
		if o.state != nil && o.state.hasDefaults {
			config = o.state.defaults
		} else if newDefault != nil {
			config = newDefault()
		}
		if o.state == nil {
			o.state = &setState[T]{}
		}
		o.state.values = slices.Grow(o.state.values, missing)
	}
	for _, selected := range selector.indexes {
		if selected < len(o.Values) {
			fn(&o.Values[selected])
			continue
		}
		index, configured := findValue(o.state.values, selected)
		if configured {
			fn(&o.state.values[index].config)
			continue
		}
		o.state.values = append(o.state.values, value[T]{})
		copy(o.state.values[index+1:], o.state.values[index:])
		o.state.values[index] = value[T]{
			index:  selected,
			config: config,
		}
		fn(&o.state.values[index].config)
	}
}

// Range calls fn for every explicit column setting.
func (o *Set[T]) Range(fn func(*T)) {
	for index := range o.Values {
		fn(&o.Values[index])
	}
	if o.state == nil {
		return
	}
	for index := range o.state.values {
		fn(&o.state.values[index].config)
	}
}

// Any reports whether fn matches an explicit column setting.
func (o *Set[T]) Any(fn func(*T) bool) bool {
	for index := range o.Values {
		if fn(&o.Values[index]) {
			return true
		}
	}
	if o.state == nil {
		return false
	}
	for index := range o.state.values {
		if fn(&o.state.values[index].config) {
			return true
		}
	}
	return false
}

// Resolve applies defaults and explicit settings to logical columns.
// indexOffset reserves leading synthetic columns that selectors cannot target.
func (o *Set[T]) Resolve(columns []T, columnCount, indexOffset int, defaults T) []T {
	if o.state == nil {
		return o.resolveValues(columns, columnCount, indexOffset, defaults)
	}
	return o.resolveState(columns, columnCount, indexOffset, defaults)
}

// resolveValues applies defaults and contiguous explicit settings.
func (o *Set[T]) resolveValues(columns []T, columnCount, indexOffset int, defaults T) []T {
	if cap(columns) < columnCount {
		columns = make([]T, columnCount)
	} else {
		columns = columns[:columnCount]
	}
	offset := min(indexOffset, columnCount)
	for index := range offset {
		columns[index] = defaults
	}
	for index := offset; index < columnCount; index++ {
		columns[index] = defaults
	}
	copy(columns[offset:], o.Values)
	return columns
}

// resolveState resolves Values using state defaults when present, then applies sparse settings.
func (o *Set[T]) resolveState(columns []T, columnCount, indexOffset int, defaults T) []T {
	state := o.state
	inputDefaults := defaults
	if state.hasDefaults {
		inputDefaults = state.defaults
	}
	columns = o.resolveValues(columns, columnCount, indexOffset, inputDefaults)
	offset := min(indexOffset, columnCount)
	if state.hasDefaults {
		for index := range offset {
			columns[index] = defaults
		}
	}
	inputColumns := columns[offset:]
	for _, value := range state.values {
		if value.index >= len(inputColumns) {
			break
		}
		inputColumns[value.index] = value.config
	}
	return columns
}

// setState holds settings that are absent from the common dense representation.
type setState[T any] struct {
	values      []value[T]
	defaults    T
	hasDefaults bool
}

// value pairs an input column index with its explicit settings.
type value[T any] struct {
	index  int
	config T
}

// findValue returns the sorted position of index in values.
func findValue[T any](values []value[T], index int) (int, bool) {
	position := 0
	for position < len(values) && values[position].index < index {
		position++
	}
	return position, position < len(values) && values[position].index == index
}

// MaxColumns returns the greatest column count among rows.
func MaxColumns(rows [][]string) int {
	columnCount := 0
	for _, row := range rows {
		columnCount = max(columnCount, len(row))
	}
	return columnCount
}
