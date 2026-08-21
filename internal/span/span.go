// Package span detects consecutive equal display values for vertical and
// horizontal cell spanning.
//
// The package returns bit masks describing cells that continue a vertical
// run or are absorbed by a horizontal run. Format packages translate those
// masks into the representation supported by their output format, such as
// blank cells, HTML span counts, or merged text borders. Vertical detection
// uses [PreviousRow] to retain values from the preceding row. Horizontal
// detection requires only the current row.
package span

import "slices"

// PreviousRow retains copies of the values used for the preceding row. Copies
// are required because arena storage may invalidate the original string
// views between rows.
//
// Reset preserves the copies' capacity so pooled arena state can reuse it.
type PreviousRow struct {
	values  [][]byte // The previous row's values, one per rendered position.
	started bool     // Whether a row has been compared since Reset.
}

// Reset clears the comparison state without discarding reusable storage.
func (o *PreviousRow) Reset() {
	o.started = false
}

// Rowspans detects vertical spans in one row. Bit i is set when selected value
// i equals the value above it. Among selected positions, a mismatch prevents
// every position to its right from spanning, and the first row after Reset
// never spans. Values after a mismatch are retained for the next comparison.
//
// The uint64 masks limit span tracking to the first 64 positions.
func Rowspans(rowspan uint64, values []string, previous *PreviousRow) uint64 {
	if rowspan == 0 {
		return 0
	}
	if missing := len(values) - len(previous.values); missing > 0 {
		previous.values = slices.Grow(previous.values, missing)
		previous.values = previous.values[:len(values)]
	}
	var spanned uint64
	spanning := previous.started
	for i, value := range values {
		if rowspan&(1<<uint(i)) == 0 {
			continue
		}
		if spanning && string(previous.values[i]) == value {
			spanned |= 1 << uint(i)
			continue
		}
		spanning = false
		previous.values[i] = append(previous.values[i][:0], value...)
	}
	previous.started = true
	return spanned
}

// Colspans detects horizontal spans in one row. Bit i is set when selected
// value i equals and is absorbed into its left neighbor. Both positions must
// be selected and neither may be present in taken. Equal runs can chain; the
// caller counts consecutive bits to size the leading cell.
//
// The result depends only on values and taken. A caller may reject spans for
// format-specific reasons; HTML does so when a merge would not be rectangular.
func Colspans(colspan uint64, values []string, taken uint64) uint64 {
	if colspan == 0 {
		return 0
	}
	var absorbed uint64
	for j := len(values) - 1; j > 0; j-- {
		if colspan&(1<<uint(j)) == 0 || colspan&(1<<uint(j-1)) == 0 {
			continue
		}
		if taken&(1<<uint(j)) != 0 || taken&(1<<uint(j-1)) != 0 {
			continue
		}
		if values[j] != values[j-1] {
			continue
		}
		absorbed |= 1 << uint(j)
	}
	return absorbed
}
