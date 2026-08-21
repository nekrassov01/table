package markdown

import (
	"strings"

	"github.com/nekrassov01/table/internal/decorate"
)

// Decoration holds Markdown or HTML markup that decorates a cell value.
type Decoration = decorate.Decoration

// Decoration presets provide common Markdown and HTML markup.
var (
	// DecorationBold wraps values in bold markup.
	DecorationBold = NewDecoration("**", "**")

	// DecorationUnderline wraps values in u elements.
	DecorationUnderline = NewDecoration("<u>", "</u>")

	// DecorationItalic wraps values in italic markup.
	DecorationItalic = NewDecoration("*", "*")

	// DecorationStrikethrough wraps values in strikethrough markup.
	DecorationStrikethrough = NewDecoration("~~", "~~")

	// DecorationCode wraps values in inline code markup.
	DecorationCode = NewDecoration("`", "`")

	// DecorationPreformatted wraps values in pre elements to preserve
	// whitespace.
	DecorationPreformatted = NewDecoration("<pre>", "</pre>")
)

// NewDecoration returns a Decoration that wraps values between prefix and
// suffix. It returns nil when prefix is empty. Both strings are emitted
// verbatim, so callers must provide valid, trusted Markdown or HTML markup.
func NewDecoration(prefix, suffix string) *Decoration {
	return decorate.New(prefix, suffix)
}

// resolveTicks returns a code-span fence longer than every backtick run in s,
// or zero when decoration is not a backtick fence.
func resolveTicks(decoration *Decoration, s string) int {
	if decoration.IsZero() || decoration.Prefix != decoration.Suffix ||
		strings.Count(decoration.Prefix, "`") != len(decoration.Prefix) {
		return 0
	}
	ticks := len(decoration.Prefix)
	for index := 0; index < len(s); {
		if s[index] != '`' {
			index++
			continue
		}
		run := 1
		for index+run < len(s) && s[index+run] == '`' {
			run++
		}
		if run >= ticks {
			ticks = run + 1
		}
		index += run
	}
	return ticks
}
