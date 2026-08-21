// Package caption defines caption positions shared by format packages.
package caption

// Side controls which side of the table the caption renders on.
type Side uint8

const (
	// Default preserves each package's default side.
	Default Side = iota

	// Top renders the caption above the table.
	Top

	// Bottom renders the caption below the table.
	Bottom
)
