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
	index := o.indexQuoteValue(s)
	if index < 0 {
		return s
	}
	state := o.state
	start := len(state.quotes)
	state.quotes = append(state.quotes, '"')
	state.quotes = append(state.quotes, s[:index]...)
	crlf := o.input.option.crlf
	for ; index < len(s); index++ {
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

// indexQuoteValue returns the first byte from which encoding/csv-compatible
// quoting must proceed, or -1 when s can be used unchanged.
func (o *compiler) indexQuoteValue(s string) int {
	if s == "" {
		return -1
	}
	if s == `\.` {
		return 0
	}
	firstByte := s[0]
	if firstByte == ' ' || (firstByte >= '\t' && firstByte <= '\r') {
		return 0
	}
	if firstByte >= utf8.RuneSelf {
		first, _ := utf8.DecodeRuneInString(s)
		if unicode.IsSpace(first) {
			return 0
		}
	}
	delimiter := o.input.option.delimiter
	if delimiter < utf8.RuneSelf {
		for index := 0; index < len(s); index++ {
			value := s[index]
			if rune(value) == delimiter || value == '"' || value == '\n' || value == '\r' {
				return index
			}
		}
		return -1
	}
	index := strings.IndexAny(s, "\"\r\n")
	delimiterIndex := strings.IndexRune(s, delimiter)
	if delimiterIndex >= 0 && (index < 0 || delimiterIndex < index) {
		return delimiterIndex
	}
	return index
}
