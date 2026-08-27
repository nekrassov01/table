package backlog

import (
	"strings"

	"github.com/nekrassov01/table/internal/unsafe"
)

// br is the Backlog notation for a line break.
const br = "&br;"

// escapeValue literalizes Backlog notation and normalizes line breaks,
// appending changed content to escapes.
func escapeValue(escapes []byte, s string) (string, []byte) {
	index := indexEscapeValue(s)
	if index < 0 {
		return s, escapes
	}
	start := len(escapes)
	escapes = append(escapes, s[:index]...)
	for ; index < len(s); index++ {
		escape := false
		switch s[index] {
		case '\\':
			escapes = append(escapes, '\\', '\\', '\\')
			continue
		case '|':
			escape = true
		case '\r', '\n':
			if s[index] == '\r' && index+1 < len(s) && s[index+1] == '\n' {
				index++
			}
			escapes = append(escapes, br...)
			continue
		case '\'', '%', '[', ']':
			escape = index > 0 && s[index-1] == s[index] ||
				index+1 < len(s) && s[index+1] == s[index]
		case '&':
			escape = strings.HasPrefix(s[index:], "&br;") ||
				strings.HasPrefix(s[index:], "&color(")
		case '{':
			escape = strings.HasPrefix(s[index:], "{quote}") ||
				strings.HasPrefix(s[index:], "{/quote}") ||
				strings.HasPrefix(s[index:], "{code}") ||
				strings.HasPrefix(s[index:], "{code:") ||
				strings.HasPrefix(s[index:], "{/code}")
		case '#':
			escape = strings.HasPrefix(s[index:], "#attach(") ||
				strings.HasPrefix(s[index:], "#image(") ||
				strings.HasPrefix(s[index:], "#thumbnail(") ||
				strings.HasPrefix(s[index:], "#rev(") ||
				strings.HasPrefix(s[index:], "#contents")
		}
		if escape {
			escapes = append(escapes, '\\', '\\')
		}
		escapes = append(escapes, s[index])
	}
	return unsafe.View(escapes[start:]), escapes
}

// indexEscapeValue returns the first byte that requires Backlog literalization
// or line-break normalization, or -1 when s can be used unchanged.
func indexEscapeValue(s string) int {
	for index := 0; index < len(s); index++ {
		switch s[index] {
		case '\\', '|', '\r', '\n':
			return index
		case '\'', '%', '[', ']':
			if index > 0 && s[index-1] == s[index] ||
				index+1 < len(s) && s[index+1] == s[index] {
				return index
			}
		case '&':
			if strings.HasPrefix(s[index:], "&br;") ||
				strings.HasPrefix(s[index:], "&color(") {
				return index
			}
		case '{':
			if strings.HasPrefix(s[index:], "{quote}") ||
				strings.HasPrefix(s[index:], "{/quote}") ||
				strings.HasPrefix(s[index:], "{code}") ||
				strings.HasPrefix(s[index:], "{code:") ||
				strings.HasPrefix(s[index:], "{/code}") {
				return index
			}
		case '#':
			if strings.HasPrefix(s[index:], "#attach(") ||
				strings.HasPrefix(s[index:], "#image(") ||
				strings.HasPrefix(s[index:], "#thumbnail(") ||
				strings.HasPrefix(s[index:], "#rev(") ||
				strings.HasPrefix(s[index:], "#contents") {
				return index
			}
		}
	}
	return -1
}
