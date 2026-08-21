package html

import (
	"strings"
	"unicode/utf8"

	"github.com/nekrassov01/table/internal/unsafe"
)

// br is the HTML tag for a line break.
const br = "<br>"

// replacement substitutes for invalid characters and malformed UTF-8.
const replacement = "\uFFFD"

// escapeValue appends an escaped display value to escapes and returns the
// value together with the updated storage.
func escapeValue(escapes []byte, s string) (string, []byte) {
	if !needsEscapeValue(s) {
		return s, escapes
	}
	start := len(escapes)
	for index := 0; index < len(s); index++ {
		switch {
		case s[index] == '&':
			escapes = append(escapes, "&amp;"...)
		case s[index] == '<':
			escapes = append(escapes, "&lt;"...)
		case s[index] == '>':
			escapes = append(escapes, "&gt;"...)
		case s[index] == '"':
			escapes = append(escapes, "&quot;"...)
		case s[index] == '\r' || s[index] == '\n':
			if s[index] == '\r' && index+1 < len(s) && s[index+1] == '\n' {
				index++
			}
			escapes = append(escapes, br...)
		case s[index] < 0x20 && s[index] != '\t', s[index] == 0x7f:
			escapes = append(escapes, replacement...)
		case s[index] < utf8.RuneSelf:
			escapes = append(escapes, s[index])
		default:
			r, size := utf8.DecodeRuneInString(s[index:])
			if r == utf8.RuneError && size == 1 {
				escapes = append(escapes, replacement...)
			} else {
				escapes = append(escapes, s[index:index+size]...)
				index += size - 1
			}
		}
	}
	return unsafe.View(escapes[start:]), escapes
}

// needsEscapeValue reports whether s requires HTML escaping or UTF-8 normalization.
// It validates UTF-8 only after observing a non-ASCII byte.
func needsEscapeValue(s string) bool {
	nonASCII := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= utf8.RuneSelf {
			nonASCII = true
			continue
		}
		switch {
		case c == '&' || c == '<' || c == '>' || c == '"' || c == '\r' || c == '\n':
			return true
		case c < 0x20 && c != '\t', c == 0x7f:
			return true
		}
	}
	return nonASCII && !utf8.ValidString(s)
}

// escapeAttr produces a normalized double-quoted HTML attribute value. It
// returns already-safe strings unchanged without allocating.
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
		case c == '\r':
			if index+1 < len(s) && s[index+1] == '\n' {
				index++
			}
			b.WriteByte('\n')
			index++
		case c < 0x20 && c != '\t' && c != '\n' && c != '\f', c == 0x7f:
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
// double-quoted HTML attribute.
func needsEscapeAttr(s string) bool {
	for index := 0; index < len(s); {
		c := s[index]
		if c < utf8.RuneSelf {
			switch {
			case c == '&' || c == '"' || c == '\r':
				return true
			case c < 0x20 && c != '\t' && c != '\n' && c != '\f', c == 0x7f:
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
