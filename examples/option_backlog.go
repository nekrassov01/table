package examples

import (
	"fmt"
	"strings"

	"github.com/nekrassov01/table/backlog"
)

// BacklogOptionSimple configures the simple Backlog table example.
var BacklogOptionSimple = []backlog.Option{
	backlog.WithHeader(SimpleData.Header...),
}

// BacklogOptionRowspan configures the Backlog table example with row spans.
var BacklogOptionRowspan = []backlog.Option{
	backlog.WithHeader(RowspanData.Header...),
	backlog.WithRowspan(backlog.ScopeBody, backlog.Columns(0, 1, 2)),
}

// BacklogOptionColspan configures the Backlog table example with column spans.
var BacklogOptionColspan = []backlog.Option{
	backlog.WithHeader(ColspanData.Header...),
	backlog.WithColspan(backlog.ScopeBody, backlog.Columns(0, 1, 2, 3)),
}

// BacklogOptionFooter configures the Backlog table example with a footer.
var BacklogOptionFooter = []backlog.Option{
	backlog.WithHeader(FooterData.Header...),
	backlog.WithFooter(FooterData.Footer),
	backlog.WithRowspan(backlog.ScopeBody, backlog.Columns(0)),
	backlog.WithColspan(backlog.ScopeFooter, backlog.Columns(0, 1, 2, 3)),
}

// BacklogOptionTransformer configures the transformed Backlog table example.
var BacklogOptionTransformer = []backlog.Option{
	backlog.WithHeader(FooterData.Header...),
	backlog.WithFooter(FooterData.Footer),
	backlog.WithRowspan(backlog.ScopeBody, backlog.Columns(0)),
	backlog.WithColspan(backlog.ScopeFooter, backlog.Columns(0, 1, 2, 3)),
	backlog.WithTransformer(backlog.Columns(5), func(v any) (string, *backlog.Color, *backlog.Decoration) {
		n, ok := v.(int)
		if !ok {
			return "", nil, nil
		}
		if n >= 3000 {
			return fmt.Sprintf("*%d", n), backlog.ColorFgRed, backlog.DecorationBold
		}
		return "", nil, nil
	}),
	backlog.WithTransformer(backlog.Columns(9), func(v any) (string, *backlog.Color, *backlog.Decoration) {
		if n, ok := v.(int); ok && n >= 90 {
			return fmt.Sprintf("*%d", n), backlog.ColorFgYellow, backlog.DecorationBold
		}
		return "", nil, nil
	}),
	backlog.WithTransformer(backlog.Columns(13), func(v any) (string, *backlog.Color, *backlog.Decoration) {
		if n, ok := v.(int); ok && n >= 60 {
			return fmt.Sprintf("*%d", n), backlog.ColorFgGreen, backlog.DecorationBold
		}
		return "", nil, nil
	}),
}

// BacklogOptionComplex configures the complex Backlog table example.
var BacklogOptionComplex = []backlog.Option{
	backlog.WithHeader(ComplexData.Header...),
	backlog.WithColor(backlog.ScopeBody, backlog.Columns(8, 9, 10), backlog.ColorFgBlack),
	backlog.WithDecoration(backlog.ScopeBody, backlog.Columns(11), backlog.DecorationBold),
	backlog.WithTransformer(backlog.Columns(5), func(v any) (string, *backlog.Color, *backlog.Decoration) {
		values, ok := v.([]string)
		if !ok {
			return "", nil, nil
		}
		tokens := make([]string, len(values))
		for i, s := range values {
			tokens[i] = fmt.Sprintf("%d/%d: %s", i+1, len(values), s)
		}
		return strings.Join(tokens, "\n"), backlog.ColorBgGreen, backlog.DecorationBold
	}),
	backlog.WithTransformer(backlog.Columns(6), func(v any) (string, *backlog.Color, *backlog.Decoration) {
		values, ok := v.([3]string)
		if !ok {
			return "", nil, nil
		}
		tokens := make([]string, len(values))
		for i, s := range values {
			tokens[i] = fmt.Sprintf("%d/%d: %s", i+1, len(values), s)
		}
		return strings.Join(tokens, "\n"), backlog.ColorBgYellow, backlog.DecorationItalic
	}),
	backlog.WithTransformer(backlog.Columns(7), func(v any) (string, *backlog.Color, *backlog.Decoration) {
		values, ok := v.([]int)
		if !ok {
			return "", nil, nil
		}
		sum := 0
		for _, value := range values {
			sum += value
		}
		return fmt.Sprintf("sum=%d", sum), backlog.ColorFgBlue, backlog.DecorationStrikethrough
	}),
}

// BacklogOptionStackedHeader configures the Backlog table example with a stacked header.
var BacklogOptionStackedHeader = []backlog.Option{
	backlog.WithHeader(StackedHeaderData.Header...),
	backlog.WithRowspan(backlog.ScopeHeader, backlog.Columns(4)),
	backlog.WithColspan(backlog.ScopeHeader, backlog.Columns(0, 1, 2, 3)),
}
