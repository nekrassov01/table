package examples

import (
	"fmt"
	"strings"

	"github.com/nekrassov01/table/html"
)

// HTMLOptionSimple configures the simple HTML table example.
var HTMLOptionSimple = []html.Option{
	html.WithHeader(SimpleData.Header...),
}

// HTMLOptionRowspan configures the HTML table example with row spans.
var HTMLOptionRowspan = []html.Option{
	html.WithHeader(RowspanData.Header...),
	html.WithRowspan(html.ScopeBody, html.Columns(0, 1, 2)),
	html.WithAlign(html.ScopeBody, html.Columns(4, 5), html.AlignRight),
}

// HTMLOptionColspan configures the HTML table example with column spans.
var HTMLOptionColspan = []html.Option{
	html.WithHeader(ColspanData.Header...),
	html.WithColspan(html.ScopeBody, html.Columns(0, 1, 2, 3)),
}

// HTMLOptionFooter configures the HTML table example with a footer.
var HTMLOptionFooter = []html.Option{
	html.WithHeader(FooterData.Header...),
	html.WithFooter(FooterData.Footer),
	html.WithRowspan(html.ScopeBody, html.Columns(0)),
	html.WithColspan(html.ScopeFooter, html.Columns(0, 1, 2, 3)),
	html.WithAlign(html.ScopeBody, html.Columns(3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13), html.AlignRight),
	html.WithAlign(html.ScopeFooter, html.Columns(4, 5, 6, 7, 8, 9, 10, 11, 12, 13), html.AlignRight),
	html.WithAlign(html.ScopeFooter, html.Columns(0), html.AlignCenter),
}

// HTMLOptionTransformer configures the transformed HTML table example.
var HTMLOptionTransformer = []html.Option{
	html.WithHeader(FooterData.Header...),
	html.WithFooter(FooterData.Footer),
	html.WithRowspan(html.ScopeBody, html.Columns(0)),
	html.WithColspan(html.ScopeFooter, html.Columns(0, 1, 2, 3)),
	html.WithAlign(html.ScopeBody, html.Columns(3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13), html.AlignRight),
	html.WithAlign(html.ScopeFooter, html.Columns(4, 5, 6, 7, 8, 9, 10, 11, 12, 13), html.AlignRight),
	html.WithAlign(html.ScopeFooter, html.Columns(0), html.AlignCenter),
	html.WithTransformer(html.Columns(5), func(v any) (string, *html.Color, *html.Decoration) {
		n, ok := v.(int)
		if !ok {
			return "", nil, nil
		}
		if n >= 3000 {
			return fmt.Sprintf("*%d", n), html.ColorFgRed, html.DecorationBold
		}
		return "", nil, nil
	}),
	html.WithTransformer(html.Columns(9), func(v any) (string, *html.Color, *html.Decoration) {
		if n, ok := v.(int); ok && n >= 90 {
			return fmt.Sprintf("*%d", n), html.ColorFgYellow, html.DecorationBold
		}
		return "", nil, nil
	}),
	html.WithTransformer(html.Columns(13), func(v any) (string, *html.Color, *html.Decoration) {
		if n, ok := v.(int); ok && n >= 60 {
			return fmt.Sprintf("*%d", n), html.ColorFgGreen, html.DecorationBold
		}
		return "", nil, nil
	}),
}

// HTMLOptionComplex configures the complex HTML table example.
var HTMLOptionComplex = []html.Option{
	html.WithHeader(ComplexData.Header...),
	html.WithColor(html.ScopeBody, html.Columns(8, 9, 10), html.ColorFgBlack),
	html.WithDecoration(html.ScopeBody, html.Columns(11), html.DecorationPreformatted),
	html.WithTransformer(html.Columns(5), func(v any) (string, *html.Color, *html.Decoration) {
		values, ok := v.([]string)
		if !ok {
			return "", nil, nil
		}
		tokens := make([]string, len(values))
		for i, s := range values {
			tokens[i] = fmt.Sprintf("%d/%d: %s", i+1, len(values), s)
		}
		return strings.Join(tokens, "\n"), html.ColorBgGreen, html.DecorationUnderline
	}),
	html.WithTransformer(html.Columns(6), func(v any) (string, *html.Color, *html.Decoration) {
		values, ok := v.([3]string)
		if !ok {
			return "", nil, nil
		}
		tokens := make([]string, len(values))
		for i, s := range values {
			tokens[i] = fmt.Sprintf("%d/%d: %s", i+1, len(values), s)
		}
		return strings.Join(tokens, "\n"), html.ColorBgMagenta, html.DecorationItalic
	}),
	html.WithTransformer(html.Columns(7), func(v any) (string, *html.Color, *html.Decoration) {
		values, ok := v.([]int)
		if !ok {
			return "", nil, nil
		}
		sum := 0
		for _, value := range values {
			sum += value
		}
		return fmt.Sprintf("sum=%d", sum), html.ColorFgCyan, html.DecorationBold
	}),
	html.WithCaption("⚡️ Rendered by <github.com/nekrassov01/table/html>", html.CaptionDefault),
}

// HTMLOptionStackedHeader configures the HTML table example with a stacked header.
var HTMLOptionStackedHeader = []html.Option{
	html.WithHeader(StackedHeaderData.Header...),
	html.WithRowspan(html.ScopeHeader, html.Columns(4)),
	html.WithColspan(html.ScopeHeader, html.Columns(0, 1, 2, 3)),
}
