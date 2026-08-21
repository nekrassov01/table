// Package decorate defines the wrapper used by format packages to attach
// format-specific decoration markup to cell values.
package decorate

// Decoration holds the markers wrapping a decorated cell value.
type Decoration struct {
	Prefix string // Markup written before a value.
	Suffix string // Markup written after a value.
}

// New returns a Decoration that wraps values between prefix and suffix. It
// returns nil when prefix is empty. Both strings bypass cell escaping, so
// callers must provide valid markup for the target format.
func New(prefix, suffix string) *Decoration {
	if prefix == "" {
		return nil
	}
	return &Decoration{
		Prefix: prefix,
		Suffix: suffix,
	}
}

// IsZero reports whether the receiver is nil or has no prefix.
func (o *Decoration) IsZero() bool {
	return o == nil || o.Prefix == ""
}
