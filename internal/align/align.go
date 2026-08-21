// Package align defines horizontal alignment values shared by format packages.
package align

// Side controls per-column horizontal text alignment.
type Side uint8

const (
	// Default preserves each package's default behavior.
	Default Side = iota

	// Left left-aligns all values in the column.
	Left

	// Right right-aligns all values in the column.
	Right

	// Center center-aligns all values in the column.
	Center
)

// String returns "left", "right", or "center" for a concrete alignment.
// It returns "default" for the zero value and unknown values because they do
// not request a concrete alignment.
func (o Side) String() string {
	switch o {
	case Left:
		return "left"
	case Right:
		return "right"
	case Center:
		return "center"
	}
	return "default"
}
