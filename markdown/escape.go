package markdown

import (
	"strings"
	"unicode/utf8"

	"github.com/nekrassov01/table/internal/unsafe"
)

const (
	// br is the HTML tag for a line break.
	br = "<br>"

	// replacement substitutes for invalid characters and malformed UTF-8.
	replacement = "\uFFFD"
)

// escapeCode prepares s for a GFM code span, appending changed content to
// escapes, and returns the value together with the updated storage.
func escapeCode(escapes []byte, s string) (string, []byte) {
	pad := strings.HasPrefix(s, "`") || strings.HasSuffix(s, "`")
	if !pad && strings.HasPrefix(s, " ") && strings.HasSuffix(s, " ") {
		for index := 0; index < len(s); index++ {
			if s[index] != ' ' && s[index] != '\r' && s[index] != '\n' {
				pad = true
				break
			}
		}
	}
	if !pad && !strings.ContainsAny(s, "|\r\n\x00") && utf8.ValidString(s) {
		return s, escapes
	}
	start := len(escapes)
	if pad {
		escapes = append(escapes, ' ')
	}
	for index := 0; index < len(s); index++ {
		switch {
		case s[index] == '|':
			escapes = append(escapes, '\\', '|')
		case s[index] == '\r' || s[index] == '\n':
			if s[index] == '\r' && index+1 < len(s) && s[index+1] == '\n' {
				index++
			}
			escapes = append(escapes, ' ')
		case s[index] == 0:
			escapes = append(escapes, replacement...)
		case s[index] < utf8.RuneSelf:
			escapes = append(escapes, s[index])
		default:
			_, size := utf8.DecodeRuneInString(s[index:])
			if size == 1 {
				escapes = append(escapes, replacement...)
				continue
			}
			escapes = append(escapes, s[index:index+size]...)
			index += size - 1
		}
	}
	if pad {
		escapes = append(escapes, ' ')
	}
	return unsafe.View(escapes[start:]), escapes
}

// escapeValue literalizes GFM and HTML syntax, normalizes invalid input, and
// appends changed content to escapes. It returns the value together with the
// updated storage.
func escapeValue(escapes []byte, s string) (string, []byte) {
	if !needsEscapeValue(s) {
		return s, escapes
	}
	start := len(escapes)
	for index := 0; index < len(s); index++ {
		switch s[index] {
		case '\\', '|', '`', '*', '_', '~', '[', ']', '<', '>', '&':
			escapes = append(escapes, '\\', s[index])
		case '\r', '\n':
			if s[index] == '\r' && index+1 < len(s) && s[index+1] == '\n' {
				index++
			}
			escapes = append(escapes, br...)
		case 0:
			escapes = append(escapes, replacement...)
		default:
			if s[index] < utf8.RuneSelf {
				escapes = append(escapes, s[index])
				continue
			}
			_, size := utf8.DecodeRuneInString(s[index:])
			if size == 1 {
				escapes = append(escapes, replacement...)
				continue
			}
			escapes = append(escapes, s[index:index+size]...)
			index += size - 1
		}
	}
	return unsafe.View(escapes[start:]), escapes
}

// needsEscapeValue reports whether s requires GFM literalization or UTF-8
// normalization.
func needsEscapeValue(s string) bool {
	nonASCII := false
	for index := 0; index < len(s); index++ {
		if s[index] >= utf8.RuneSelf {
			nonASCII = true
			continue
		}
		switch s[index] {
		case '\\', '|', '`', '*', '_', '~', '[', ']', '<', '>', '&', '\r', '\n', 0:
			return true
		}
	}
	return nonASCII && !utf8.ValidString(s)
}

// escapeAttr produces a normalized double-quoted HTML attribute value without
// introducing a GFM line ending. It returns already-safe strings unchanged
// without allocating.
func escapeAttr(s string) string {
	if !needsEscapeAttr(s) {
		return s
	}
	var b strings.Builder
	b.Grow(len(s) + 8)
	for index := 0; index < len(s); {
		c := s[index]
		switch {
		case c == '&':
			b.WriteString("&amp;")
			index++
		case c == '"':
			b.WriteString("&quot;")
			index++
		case c == '\r' || c == '\n':
			if c == '\r' && index+1 < len(s) && s[index+1] == '\n' {
				index++
			}
			b.WriteByte(' ')
			index++
		case c < 0x20 && c != '\t' && c != '\f', c == 0x7f:
			b.WriteString(replacement)
			index++
		case c < utf8.RuneSelf:
			b.WriteByte(c)
			index++
		default:
			r, size := utf8.DecodeRuneInString(s[index:])
			if r == utf8.RuneError && size == 1 ||
				r >= 0x80 && r <= 0x9f ||
				r >= 0xfdd0 && r <= 0xfdef ||
				r&0xfffe == 0xfffe {
				b.WriteString(replacement)
			} else {
				b.WriteString(s[index : index+size])
			}
			index += size
		}
	}
	return b.String()
}

// needsEscapeAttr reports whether s requires escaping or normalization for a
// double-quoted HTML attribute in a GFM table row.
func needsEscapeAttr(s string) bool {
	for index := 0; index < len(s); {
		c := s[index]
		if c < utf8.RuneSelf {
			switch {
			case c == '&' || c == '"' || c == '\r' || c == '\n':
				return true
			case c < 0x20 && c != '\t' && c != '\f', c == 0x7f:
				return true
			}
			index++
			continue
		}
		r, size := utf8.DecodeRuneInString(s[index:])
		if r == utf8.RuneError && size == 1 ||
			r >= 0x80 && r <= 0x9f ||
			r >= 0xfdd0 && r <= 0xfdef ||
			r&0xfffe == 0xfffe {
			return true
		}
		index += size
	}
	return false
}
