package value

import "strconv"

// summaryCacheSize is the largest element count with a pre-rendered summary.
const summaryCacheSize = 64

// summary pre-renders composite-type placeholders for counts from zero through
// summaryCacheSize. Counts beyond the limit are built on demand.
type summary struct {
	entries [summaryCacheSize + 1]string
	prefix  string
	suffix  string
}

var (
	// summaryStruct is the placeholder for struct values.
	summaryStruct = newSummary("{struct ", " field(s)}")

	// summaryMap is the placeholder for map values.
	summaryMap = newSummary("{map ", " key(s)}")

	// summaryList is the placeholder for slice and array values.
	summaryList = newSummary("[list ", " item(s)]")
)

// newSummary returns a summary with all entries pre-rendered.
func newSummary(prefix, suffix string) summary {
	var s summary
	s.prefix = prefix
	s.suffix = suffix
	for i := range s.entries {
		s.entries[i] = s.build(i)
	}
	return s
}

// format returns the placeholder string for count n.
func (o *summary) format(n int) string {
	if n >= 0 && n <= summaryCacheSize {
		return o.entries[n]
	}
	return o.build(n)
}

// build renders a single placeholder string for count n.
func (o *summary) build(n int) string {
	b := make([]byte, 0, len(o.prefix)+20+len(o.suffix))
	b = append(b, o.prefix...)
	b = strconv.AppendInt(b, int64(n), 10)
	b = append(b, o.suffix...)
	return string(b)
}
