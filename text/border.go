package text

const (
	// armUp marks an upward arm.
	armUp = 1 << (3 - iota)

	// armDown marks a downward arm.
	armDown

	// armLeft marks a leftward arm.
	armLeft

	// armRight marks a rightward arm.
	armRight
)

// Horizontal defines the joints and fill of a horizontal border.
type Horizontal struct {
	Inner Joints // The joints crossing a column separator.
	Outer Joints // The joints meeting the outer border, at either end.
	Fill  string // The body of the line between joints.
}

// maxGlyphLen returns the longest byte length of any horizontal border glyph,
// or 0 if the receiver is nil.
func (o *Horizontal) maxGlyphLen() int {
	if o == nil {
		return 0
	}
	return max(
		o.Inner.maxGlyphLen(),
		o.Outer.maxGlyphLen(),
		len(o.Fill),
	)
}

// Vertical defines outer and inner vertical border glyphs.
type Vertical struct {
	Outer string // The border beside the outermost columns.
	Inner string // The separator between two columns.
}

// maxGlyphLen returns the longest byte length of any vertical border glyph,
// or 0 if the receiver is nil.
func (o *Vertical) maxGlyphLen() int {
	if o == nil {
		return 0
	}
	return max(
		len(o.Outer),
		len(o.Inner),
	)
}

// Joints defines glyphs for all combinations of up, down, left, and right arms.
type Joints struct {
	UDLR, UDLX, UDXR, UDXX string // ┼ ┤ ├ │
	XDLR, XDLX, XDXR, XDXX string // ┬ ┐ ┌ ╷
	UXLR, UXLX, UXXR, UXXX string // ┴ ┘ └ ╵
	XXLR, XXLX, XXXR, XXXX string // ─ ╴ ╶ (blank)
}

// maxGlyphLen returns the longest byte length of any joint glyph.
func (o *Joints) maxGlyphLen() int {
	return max(
		len(o.UDLR), len(o.UDLX), len(o.UDXR), len(o.UDXX),
		len(o.XDLR), len(o.XDLX), len(o.XDXR), len(o.XDXX),
		len(o.UXLR), len(o.UXLX), len(o.UXXR), len(o.UXXX),
		len(o.XXLR), len(o.XXLX), len(o.XXXR), len(o.XXXX),
	)
}

// resolve returns the joint glyph with the specified arms.
func (o *Joints) resolve(up, down, left, right bool) string {
	arms := 0
	if up {
		arms |= armUp
	}
	if down {
		arms |= armDown
	}
	if left {
		arms |= armLeft
	}
	if right {
		arms |= armRight
	}
	switch arms {
	case armUp | armDown | armLeft | armRight:
		return o.UDLR
	case armUp | armDown | armLeft:
		return o.UDLX
	case armUp | armDown | armRight:
		return o.UDXR
	case armUp | armDown:
		return o.UDXX
	case armDown | armLeft | armRight:
		return o.XDLR
	case armDown | armLeft:
		return o.XDLX
	case armDown | armRight:
		return o.XDXR
	case armDown:
		return o.XDXX
	case armUp | armLeft | armRight:
		return o.UXLR
	case armUp | armLeft:
		return o.UXLX
	case armUp | armRight:
		return o.UXXR
	case armUp:
		return o.UXXX
	case armLeft | armRight:
		return o.XXLR
	case armLeft:
		return o.XXLX
	case armRight:
		return o.XXXR
	default:
		return o.XXXX
	}
}
