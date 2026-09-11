package text

import (
	"io"
	"os"

	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"golang.org/x/term"
)

var (
	// terminalSize reads terminal dimensions from a file descriptor.
	terminalSize = term.GetSize

	// terminalWriter adapts ANSI attributes to a Windows terminal.
	terminalWriter = colorable.NewColorable

	// terminalWidth resolves the writer's terminal width.
	terminalWidth = resolveTerminalWidth

	// isTerminal reports whether the writer is a terminal.
	isTerminal = resolveIsTerminal
)

// writer retains the terminal file used by an adapted output writer.
type writer struct {
	io.Writer

	file *os.File
}

// Fd returns the original terminal file descriptor.
func (o *writer) Fd() uintptr {
	return o.file.Fd()
}

// resolveWriter adapts a terminal file for Windows color output.
func resolveWriter(w io.Writer) io.Writer {
	file, ok := w.(*os.File)
	if !ok {
		return w
	}
	resolved := terminalWriter(file)
	if resolvedFile, ok := resolved.(*os.File); ok && resolvedFile == file {
		return file
	}
	return &writer{
		Writer: resolved,
		file:   file,
	}
}

// resolveTerminalWidth returns the terminal width of w, or 0 when w is not a
// terminal.
func resolveTerminalWidth(w io.Writer) int {
	descriptor, ok := w.(interface{ Fd() uintptr })
	if !ok {
		return 0
	}
	width, _, err := terminalSize(int(descriptor.Fd()))
	if err != nil {
		return 0
	}
	return width
}

// resolveIsTerminal reports whether w is a terminal device.
func resolveIsTerminal(w io.Writer) bool {
	file, ok := w.(*os.File)
	if !ok {
		return false
	}
	fd := file.Fd()
	return isatty.IsTerminal(fd) || isatty.IsCygwinTerminal(fd)
}
