package text

import "github.com/nekrassov01/table/internal/width"

// measureLine returns value's display width when it contains no line break. If
// hasBreak is true, displayWidth is zero and scanLine must measure each line.
func measureLine(value string) (displayWidth int, hasBreak bool) {
	for i := 0; i < len(value); i++ {
		if value[i]-0x20 < 0x5f {
			continue
		}
		for ; i < len(value); i++ {
			if value[i] == '\n' || value[i] == '\r' {
				return 0, true
			}
		}
		return width.StringWidth(value), false
	}
	return len(value), false
}

// scanLine returns the physical line at start, its display width, and the next
// start position. hasBreak reports whether next skipped a line break.
func scanLine(s string, start int) (line string, displayWidth, next int, hasBreak bool) {
	end := start
	simple := true
	for end < len(s) && s[end] != '\n' && s[end] != '\r' {
		simple = simple && s[end]-0x20 < 0x5f
		end++
	}
	line = s[start:end]
	displayWidth = len(line)
	if !simple {
		displayWidth = width.StringWidth(line)
	}
	next = end
	hasBreak = end < len(s)
	if hasBreak {
		next++
		if s[end] == '\r' && next < len(s) && s[next] == '\n' {
			next++
		}
	}
	return line, displayWidth, next, hasBreak
}
