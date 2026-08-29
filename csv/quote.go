package csv

import (
	"encoding/binary"
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
		//nolint:gosec // The delimiter is restricted to ASCII above.
		delimiterByte := byte(delimiter)
		if len(s) < 8 {
			for index := 0; index < len(s); index++ {
				value := s[index]
				if value == delimiterByte || value == '"' || value == '\n' || value == '\r' {
					return index
				}
			}
			return -1
		}
		return indexQuoteASCII(s, delimiterByte)
	}
	index := strings.IndexAny(s, "\"\r\n")
	delimiterIndex := strings.IndexRune(s, delimiter)
	if delimiterIndex >= 0 && (index < 0 || delimiterIndex < index) {
		return delimiterIndex
	}
	return index
}

// indexQuoteASCII scans eight-byte blocks for an ASCII delimiter or CSV
// quoting byte, then verifies a matching block byte by byte.
func indexQuoteASCII(s string, delimiter byte) int {
	const lowBits uint64 = 0x0101010101010101
	delimiterBits := uint64(delimiter) * lowBits
	index := 0
	for ; index <= len(s)-8; index += 8 {
		value := binary.NativeEndian.Uint64([]byte(s[index:]))
		matches := zeroByteBits(value ^ delimiterBits)
		matches |= zeroByteBits(value ^ ('"' * lowBits))
		matches |= zeroByteBits(value ^ ('\n' * lowBits))
		matches |= zeroByteBits(value ^ ('\r' * lowBits))
		if matches == 0 {
			continue
		}
		for candidate := index; candidate < index+8; candidate++ {
			value := s[candidate]
			if value == delimiter || value == '"' || value == '\n' || value == '\r' {
				return candidate
			}
		}
	}
	for ; index < len(s); index++ {
		value := s[index]
		if value == delimiter || value == '"' || value == '\n' || value == '\r' {
			return index
		}
	}
	return -1
}

// zeroByteBits marks byte positions that may be zero. Subtraction can also
// mark neighboring bytes, so callers must verify matching blocks.
func zeroByteBits(value uint64) uint64 {
	const lowBits uint64 = 0x0101010101010101
	const highBits uint64 = 0x8080808080808080
	return (value - lowBits) & ^value & highBits
}
