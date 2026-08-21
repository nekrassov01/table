// Package width measures terminal display width and scans strings by display
// units.
//
// ASCII text uses a direct byte-width path. Other text is segmented into
// grapheme clusters and measured with go-runewidth so that all format packages
// use the same width model. The package does not interpret line breaks or
// combine lines; the text package applies those layout rules.
package width

import (
	"github.com/clipperhouse/uax29/v2/graphemes"
	"github.com/mattn/go-runewidth"
)

// Scanner walks a string by terminal display units.
type Scanner struct {
	text     string
	clusters graphemes.Iterator[string]
	next     int
	ascii    bool
}

// Next returns the byte range and display width of the next unit.
func (o *Scanner) Next() (start, end, displayWidth int, ok bool) {
	if o.ascii {
		if o.next == len(o.text) {
			return 0, 0, 0, false
		}
		start = o.next
		o.next++
		value := o.text[start]
		if value >= 0x20 && value != 0x7f {
			displayWidth = 1
		}
		return start, o.next, displayWidth, true
	}
	if !o.clusters.Next() {
		return 0, 0, 0, false
	}
	for _, r := range o.clusters.Value() {
		displayWidth = runewidth.RuneWidth(r)
		if displayWidth > 0 {
			break
		}
	}
	return o.clusters.Start(), o.clusters.End(), displayWidth, true
}

// NewScanner returns a Scanner for s.
func NewScanner(s string) Scanner {
	o := Scanner{
		text:  s,
		ascii: true,
	}
	for i := 0; i < len(s); i++ {
		if s[i] < 0x80 {
			continue
		}
		o.clusters = *graphemes.FromString(s)
		o.ascii = false
		break
	}
	return o
}

// StringWidth returns the display width of s in terminal columns.
// Strings of printable ASCII use len(s) directly to avoid runewidth
// overhead. Multibyte runes and control bytes take the measured path;
// runewidth counts the latter as width 0. Both paths therefore produce the
// same result for the same content.
func StringWidth(s string) int {
	for i := 0; i < len(s); i++ {
		if s[i]-0x20 >= 0x5f {
			return runewidth.StringWidth(s)
		}
	}
	return len(s)
}
