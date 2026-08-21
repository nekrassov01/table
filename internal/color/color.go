// Package color defines the wrapper used by format packages to attach
// format-specific color markup to cell values.
package color

// Color holds the precomputed markup that wraps a colored value.
type Color struct {
	Prefix string // Markup written before a value.
	Suffix string // Markup written after a value.
}

// New returns a Color that wraps values between prefix and suffix. It returns
// nil when prefix is empty. Both strings bypass cell escaping, so callers must
// provide valid markup for the target format.
func New(prefix, suffix string) *Color {
	if prefix == "" {
		return nil
	}
	return &Color{
		Prefix: prefix,
		Suffix: suffix,
	}
}

// IsZero reports whether the receiver is nil or has no prefix.
func (o *Color) IsZero() bool {
	return o == nil || o.Prefix == ""
}
