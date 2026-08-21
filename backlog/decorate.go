package backlog

import "github.com/nekrassov01/table/internal/decorate"

// Decoration holds Backlog markup that decorates a cell value.
type Decoration = decorate.Decoration

// Decoration presets provide common Backlog markup.
var (
	// DecorationBold wraps values in bold markup.
	DecorationBold = NewDecoration("''", "''")

	// DecorationItalic wraps values in italic markup.
	DecorationItalic = NewDecoration("'''", "'''")

	// DecorationStrikethrough wraps values in strikethrough markup.
	DecorationStrikethrough = NewDecoration("%%", "%%")

	// DecorationCode wraps values in code markup.
	DecorationCode = NewDecoration("{code}", "{/code}")
)

// NewDecoration returns a Decoration that wraps values between prefix and
// suffix. It returns nil when prefix is empty. Both strings are emitted
// verbatim, so callers must provide valid, trusted Backlog markup.
func NewDecoration(prefix, suffix string) *Decoration {
	return decorate.New(prefix, suffix)
}
