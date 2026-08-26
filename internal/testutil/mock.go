package testutil

// Tabular records rows passed to Render and returns Err.
type Tabular struct {
	Rows [][]any
	Err  error
}

// Render records rows and returns Err.
func (o *Tabular) Render(rows [][]any) error {
	o.Rows = rows
	return o.Err
}

// Streamer records rows passed to Render and returns configured errors.
type Streamer struct {
	Rows      [][]any
	RenderErr error
	CloseErr  error
}

// Render records row and returns RenderErr.
func (o *Streamer) Render(row []any) error {
	o.Rows = append(o.Rows, row)
	return o.RenderErr
}

// Close returns CloseErr.
func (o *Streamer) Close() error {
	return o.CloseErr
}

// Error implements error for tests.
type Error struct {
	Value string
}

// Error returns o.Value.
func (o Error) Error() string {
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

// Stringer implements fmt.Stringer for tests.
type Stringer struct {
	Value string
}

// String returns o.Value.
func (o Stringer) String() string {
	return o.Value
}

// PanicStringer panics if String is called.
type PanicStringer struct{}

// String panics to detect unexpected default formatting in tests.
func (PanicStringer) String() string {
	panic("unexpected String call")
}

// PtrStringer implements fmt.Stringer with a pointer receiver for tests.
type PtrStringer struct {
	Value string
}

// String returns o.Value.
func (o *PtrStringer) String() string {
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

// MatchErrorWriter returns Err when a write exactly matches Value.
type MatchErrorWriter struct {
	Value string
	Err   error
}

// Write consumes input other than Value.
func (o *MatchErrorWriter) Write(value []byte) (int, error) {
	if string(value) == o.Value {
		return 0, o.Err
	}
	return len(value), nil
}
