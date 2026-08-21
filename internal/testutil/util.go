// Package testutil provides assertions, fixtures, and mocks shared by the
// module's tests.
package testutil

import (
	"bytes"
	"errors"
	"iter"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// goldenDir is where package tests keep the outputs they assert.
const goldenDir = "testdata"

// NewError creates a generic error for tests.
func NewError() error {
	return errors.New("test")
}

// Seq2 yields values in order, followed by err when it is non-nil.
func Seq2[T any](values []T, err error) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		for _, value := range values {
			if !yield(value, nil) {
				return
			}
		}
		if err != nil {
			var zero T
			yield(zero, err)
		}
	}
}

// NewFile returns a temporary file and a function that reads its contents.
// The file also exercises non-terminal behavior for an *os.File and
// is closed during test cleanup.
func NewFile(t *testing.T) (*os.File, func() []byte) {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "render")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = f.Close()
	})
	return f, func() []byte {
		t.Helper()
		b, err := os.ReadFile(f.Name())
		if err != nil {
			t.Fatal(err)
		}
		return b
	}
}

// AssertGolden compares got with the named golden file. It reports got rather
// than creating a missing file so new expected output is reviewed explicitly.
func AssertGolden(t *testing.T, name string, got []byte) {
	t.Helper()
	AssertGoldenDiffers(t, name, got)
	path := filepath.Clean(filepath.Join(goldenDir, name+".txt"))
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("%s: %v\nwrite the file from the output below, having read it:\n%s", name, err, got)
	}
	AssertBytes(t, got, want, name)
}

// AssertGoldenDiffers reports duplicate table and stream golden files. Equal
// output belongs in one common golden so later changes cannot update only one
// copy.
func AssertGoldenDiffers(t *testing.T, name string, got []byte) {
	t.Helper()
	var twin, shared string
	switch {
	case strings.HasPrefix(name, "table_"):
		shared = strings.TrimPrefix(name, "table_")
		twin = "stream_" + shared
	case strings.HasPrefix(name, "stream_"):
		shared = strings.TrimPrefix(name, "stream_")
		twin = "table_" + shared
	default:
		return
	}
	b, err := os.ReadFile(filepath.Clean(filepath.Join(goldenDir, twin+".txt")))
	if err != nil {
		return
	}
	if bytes.Equal(got, b) {
		t.Fatalf("%s holds what %s holds: they are one golden, named common_%s", name, twin, shared)
	}
}

// AssertBytes asserts the given byte slices are equal and reports both as
// textual table output.
func AssertBytes(t *testing.T, got, want []byte, prefix string) {
	t.Helper()
	if !bytes.Equal(got, want) {
		t.Errorf("%s mismatch\ngot:\n%s\nwant:\n%s", prefix, got, want)
	}
}

// AssertValue asserts the given values are equal.
func AssertValue(t *testing.T, got, want any, prefix string) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%s value mismatch\ngot:\n%v\nwant:\n%v\n", prefix, got, want)
	}
}
