package testutil

// Stringer implements fmt.Stringer for tests.
type Stringer struct {
	Value string
}

// String returns o.Value.
func (o Stringer) String() string {
	return o.Value
}

// Error implements error for tests.
type Error struct {
	Value string
}

// Error returns o.Value.
func (o Error) Error() string {
	return o.Value
}

// PtrStringer implements fmt.Stringer with a pointer receiver for tests.
type PtrStringer struct {
	Value string
}

// String returns o.Value.
func (o *PtrStringer) String() string {
	return o.Value
}

// PtrError implements error with a pointer receiver for tests.
type PtrError struct {
	Value string
}

// Error returns o.Value.
func (o *PtrError) Error() string {
	return o.Value
}

// ErrorWriter implements io.Writer by returning a configured error.
type ErrorWriter struct {
	Err error
}

// Write returns o.Err without consuming the input.
func (o *ErrorWriter) Write([]byte) (int, error) {
	return 0, o.Err
}
