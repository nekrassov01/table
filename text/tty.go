package text

import (
	"io"
	"os"

	"github.com/mattn/go-isatty"
	"golang.org/x/term"
)

var (
	// termSize reads terminal dimensions from a file descriptor.
	termSize = term.GetSize

	// terminalWidth resolves the writer's terminal width.
	terminalWidth = resolveTerminalWidth

	// isTerminal reports whether the writer is a terminal.
	isTerminal = resolveIsTerminal
)

// resolveTerminalWidth returns the terminal width of w, or 0 when w is not a
// terminal.
func resolveTerminalWidth(w io.Writer) int {
	f, ok := w.(*os.File)
	if !ok {
		return 0
	}
	width, _, err := termSize(int(f.Fd()))
	if err != nil {
		return 0
	}
	return width
}

// resolveIsTerminal reports whether w is a terminal device.
func resolveIsTerminal(w io.Writer) bool {
	if f, ok := w.(*os.File); ok {
		fd := f.Fd()
		return isatty.IsTerminal(fd) || isatty.IsCygwinTerminal(fd)
	}
	return false
}
