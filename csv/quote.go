package csv

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/nekrassov01/table/internal/unsafe"
)

// quoteValue wraps s in double quotes when it contains the field delimiter, a
// double quote, a line break, or leading Unicode whitespace. Embedded double
// quotes are doubled. CRLF mode normalizes line breaks like
// encoding/csv.Writer.
func (o *compiler) quoteValue(s string) string {
	if !o.needsQuoteValue(s) {
		return s
	}
	state := o.state
	start := len(state.quotes)
	state.quotes = append(state.quotes, '"')
	crlf := o.input.option.crlf
	for index := 0; index < len(s); index++ {
		switch s[index] {
		case '"':
			state.quotes = append(state.quotes, '"', '"')
		case '\r':
			if !crlf {
				state.quotes = append(state.quotes, '\r')
			}
		case '\n':
			if crlf {
				state.quotes = append(state.quotes, '\r')
			}
			state.quotes = append(state.quotes, '\n')
		default:
			state.quotes = append(state.quotes, s[index])
		}
	}
	state.quotes = append(state.quotes, '"')
	return unsafe.View(state.quotes[start:])
}

// needsQuoteValue reports whether s requires encoding/csv-compatible quoting.
func (o *compiler) needsQuoteValue(s string) bool {
	if s == "" {
		return false
	}
	if s == `\.` {
		return true
	}
	firstByte := s[0]
	if firstByte == ' ' || (firstByte >= '\t' && firstByte <= '\r') {
		return true
	}
	if firstByte >= utf8.RuneSelf {
		first, _ := utf8.DecodeRuneInString(s)
		if unicode.IsSpace(first) {
			return true
		}
	}
	delimiter := o.input.option.delimiter
	if delimiter < utf8.RuneSelf {
		for index := 0; index < len(s); index++ {
			value := s[index]
			if rune(value) == delimiter || value == '"' || value == '\n' || value == '\r' {
				return true
			}
		}
		return false
	}
	return strings.ContainsRune(s, delimiter) || strings.ContainsAny(s, "\"\r\n")
}
