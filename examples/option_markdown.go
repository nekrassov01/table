package examples

import (
	"fmt"
	"strings"

	"github.com/nekrassov01/table/markdown"
)

// MarkdownOptionSimple configures the simple Markdown table example.
var MarkdownOptionSimple = []markdown.Option{
	markdown.WithHeader(SimpleData.Header[0]),
}

// MarkdownOptionRowspan configures the Markdown table example with row spans.
var MarkdownOptionRowspan = []markdown.Option{
	markdown.WithHeader(RowspanData.Header[0]),
	markdown.WithRowspan(markdown.Columns(0, 1, 2)),
	markdown.WithAlign(markdown.Columns(4, 5), markdown.AlignRight),
}

// MarkdownOptionColspan configures the Markdown table example with column spans.
var MarkdownOptionColspan = []markdown.Option{
	markdown.WithHeader(ColspanData.Header[0]),
	markdown.WithColspan(markdown.Columns(0, 1, 2, 3)),
}

// MarkdownOptionTransformer configures the transformed Markdown table example.
var MarkdownOptionTransformer = []markdown.Option{
	markdown.WithHeader(FooterData.Header[0]),
	markdown.WithRowspan(markdown.Columns(0)),
	markdown.WithColspan(markdown.Columns(0, 1, 2, 3)),
	markdown.WithAlign(markdown.Columns(3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13), markdown.AlignRight),
	markdown.WithTransformer(markdown.Columns(5), func(v any) (string, *markdown.Color, *markdown.Decoration) {
		n, ok := v.(int)
		if !ok {
			return "", nil, nil
		}
		if n >= 3000 {
			return fmt.Sprintf("*%d", n), markdown.ColorFgRed, markdown.DecorationBold
		}
		return "", nil, nil
	}),
	markdown.WithTransformer(markdown.Columns(9), func(v any) (string, *markdown.Color, *markdown.Decoration) {
		if n, ok := v.(int); ok && n >= 90 {
			return fmt.Sprintf("*%d", n), markdown.ColorFgYellow, markdown.DecorationBold
		}
		return "", nil, nil
	}),
	markdown.WithTransformer(markdown.Columns(13), func(v any) (string, *markdown.Color, *markdown.Decoration) {
		if n, ok := v.(int); ok && n >= 60 {
			return fmt.Sprintf("*%d", n), markdown.ColorFgGreen, markdown.DecorationBold
		}
		return "", nil, nil
	}),
}

// MarkdownOptionComplex configures the complex Markdown table example.
var MarkdownOptionComplex = []markdown.Option{
	markdown.WithHeader(ComplexData.Header[0]),
	markdown.WithColor(markdown.ScopeBody, markdown.Columns(8, 9, 10), markdown.ColorFgBlack),
	markdown.WithDecoration(markdown.ScopeBody, markdown.Columns(11), markdown.DecorationUnderline),
	markdown.WithTransformer(markdown.Columns(5), func(v any) (string, *markdown.Color, *markdown.Decoration) {
		values, ok := v.([]string)
		if !ok {
			return "", nil, nil
		}
		tokens := make([]string, len(values))
		for i, s := range values {
			tokens[i] = fmt.Sprintf("%d/%d: %s", i+1, len(values), s)
		}
		return strings.Join(tokens, "\n"), markdown.ColorBgGreen, markdown.DecorationBold
	}),
	markdown.WithTransformer(markdown.Columns(6), func(v any) (string, *markdown.Color, *markdown.Decoration) {
		values, ok := v.([3]string)
		if !ok {
			return "", nil, nil
		}
		tokens := make([]string, len(values))
		for i, s := range values {
			tokens[i] = fmt.Sprintf("%d/%d: %s", i+1, len(values), s)
		}
		return strings.Join(tokens, "\n"), markdown.ColorBgMagenta, markdown.DecorationItalic
	}),
	markdown.WithTransformer(markdown.Columns(7), func(v any) (string, *markdown.Color, *markdown.Decoration) {
		values, ok := v.([]int)
		if !ok {
			return "", nil, nil
		}
		sum := 0
		for _, value := range values {
			sum += value
		}
		return fmt.Sprintf("sum=%d", sum), markdown.ColorFgCyan, markdown.DecorationBold
	}),
}
