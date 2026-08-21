package html

import "github.com/nekrassov01/table/internal/decorate"

// Decoration holds HTML markup that decorates a cell value.
type Decoration = decorate.Decoration

// Decoration presets provide common semantic markup.
var (
	// DecorationBold wraps values in strong elements.
	DecorationBold = NewDecoration("<strong>", "</strong>")

	// DecorationUnderline wraps values in u elements.
	DecorationUnderline = NewDecoration("<u>", "</u>")

	// DecorationItalic wraps values in em elements.
	DecorationItalic = NewDecoration("<em>", "</em>")

	// DecorationStrikethrough wraps values in del elements.
	DecorationStrikethrough = NewDecoration("<del>", "</del>")

	// DecorationCode wraps values in code elements.
	DecorationCode = NewDecoration("<code>", "</code>")

	// DecorationPreformatted wraps values in pre elements.
	DecorationPreformatted = NewDecoration("<pre>", "</pre>")
)

// NewDecoration returns a Decoration that wraps values between prefix and
// suffix. It returns nil when prefix is empty. Both strings are emitted
// verbatim, so callers must provide valid, trusted HTML.
func NewDecoration(prefix, suffix string) *Decoration {
	return decorate.New(prefix, suffix)
}
