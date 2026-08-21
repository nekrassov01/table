package markdown

import (
	"fmt"

	"github.com/nekrassov01/table"
)

// newWriteError wraps table.ErrWriteFailed and the error returned by the
// output writer.
func newWriteError(err error) error {
	return newError(fmt.Errorf("%w: %w", table.ErrWriteFailed, err))
}

// newClosedError wraps table.ErrClosed for a Render call after Close.
func newClosedError() error {
	return newError(table.ErrClosed)
}

// newColumnCountError wraps table.ErrColumnCount with the received and
// expected counts.
func newColumnCountError(got, want int) error {
	return newError(fmt.Errorf("%w: got %d, want %d", table.ErrColumnCount, got, want))
}

// newHeaderError wraps table.ErrHeaderRequired: GFM tables require a header.
func newHeaderError() error {
	return newError(table.ErrHeaderRequired)
}

// newError wraps err with the package name.
func newError(err error) error {
	return &table.Error{
		Pkg: "markdown",
		Err: err,
	}
}
