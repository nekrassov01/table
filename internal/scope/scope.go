// Package scope identifies the header, body, and footer of a table and stores
// values or column masks selected independently for those parts.
package scope

// Scope selects the table parts to which an option applies. Values may be
// combined; the zero value selects no parts.
type Scope uint8

const (
	// Header selects the header rows.
	Header Scope = 1 << iota

	// Body selects the data rows.
	Body

	// Footer selects the footer rows.
	Footer
)

// Scopes holds one independently configurable value per table part.
type Scopes[T any] struct {
	header T
	body   T
	footer T
}

// Set assigns the value to every part selected by sc.
func (o *Scopes[T]) Set(sc Scope, v T) {
	if sc&Header != 0 {
		o.header = v
	}
	if sc&Body != 0 {
		o.body = v
	}
	if sc&Footer != 0 {
		o.footer = v
	}
}

// Resolve returns the value for sc. It returns the zero value of T unless sc
// identifies exactly one table part.
func (o *Scopes[T]) Resolve(sc Scope) T {
	switch sc {
	case Header:
		return o.header
	case Body:
		return o.body
	case Footer:
		return o.footer
	}
	var zero T
	return zero
}

// Masks holds one bit mask of selected positions per table part.
//
// Masks is separate from Scopes[uint64] because Mark accumulates bits, whereas
// Scopes.Set replaces the stored value.
type Masks struct {
	header uint64
	body   uint64
	footer uint64
}

// Mark marks the position in every part selected by sc.
func (o *Masks) Mark(sc Scope, position int) {
	bit := uint64(1) << uint(position)
	if sc&Header != 0 {
		o.header |= bit
	}
	if sc&Body != 0 {
		o.body |= bit
	}
	if sc&Footer != 0 {
		o.footer |= bit
	}
}

// Resolve returns the mask for sc. It returns zero unless sc identifies
// exactly one table part.
func (o *Masks) Resolve(sc Scope) uint64 {
	switch sc {
	case Header:
		return o.header
	case Body:
		return o.body
	case Footer:
		return o.footer
	}
	return 0
}
