package backlog

import (
	"strings"

	"github.com/nekrassov01/table/internal/unsafe"
)

// br is the Backlog notation for a line break.
const br = "&br;"

// escapeValue replaces characters that break Backlog table cells, appending
// changed content to escapes, and returns the value together with the updated
// storage.
func escapeValue(escapes []byte, s string) (string, []byte) {
	if !needsEscapeValue(s) {
		return s, escapes
	}
	start := len(escapes)
	for index := 0; index < len(s); index++ {
		switch {
		case s[index] == '\\':
			escapes = append(escapes, '\\', '\\', '\\')
		case s[index] == '|':
			escapes = append(escapes, '\\', '\\', '|')
		case s[index] == '\r' || s[index] == '\n':
			if s[index] == '\r' && index+1 < len(s) && s[index+1] == '\n' {
				index++
			}
			escapes = append(escapes, br...)
		default:
			escapes = append(escapes, s[index])
		}
	}
	return unsafe.View(escapes[start:]), escapes
}

// needsEscapeValue reports whether s requires Backlog table escaping.
func needsEscapeValue(s string) bool {
	return strings.ContainsAny(s, "\\|\r\n")
}
