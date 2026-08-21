package table

import "errors"

var (
	// ErrWriteFailed reports an error returned by the output writer.
	// Callers can match it with errors.Is.
	ErrWriteFailed = errors.New("write failed")

	// ErrClosed reports a Render call after Close on a stream.
	// Callers can match it with errors.Is.
	ErrClosed = errors.New("render after close")

	// ErrColumnCount reports a body or footer row containing more cells than
	// the established column count. Callers can match it with errors.Is.
	ErrColumnCount = errors.New("column count exceeded")

	// ErrHeaderRequired reports a missing header for a format that requires
	// one. Callers can match it with errors.Is.
	ErrHeaderRequired = errors.New("header is required")

	// ErrDelimiter reports an unsupported field delimiter.
	// Callers can match it with errors.Is.
	ErrDelimiter = errors.New("invalid delimiter")
)

// Error is the outermost wrapper returned by a format package. It identifies
// the package associated with an underlying error.
type Error struct {
	Pkg string // Package name, such as "text".
	Err error  // The underlying error.
}

// Error returns the package name followed by the underlying error. If Err is
// nil, it returns only the package name.
func (o *Error) Error() string {
	if o.Err == nil {
		return o.Pkg
	}
	return o.Pkg + ": " + o.Err.Error()
}

// Unwrap returns the underlying error.
func (o *Error) Unwrap() error {
	return o.Err
}
