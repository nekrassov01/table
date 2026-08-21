package examples

import (
	"fmt"
	"strings"

	"github.com/nekrassov01/table/text"
)

var (
	textFgRedBold        = text.NewAttr(text.CodeFgRed, text.CodeBold)
	textFgYellowBold     = text.NewAttr(text.CodeFgYellow, text.CodeBold)
	textFgGreenBold      = text.NewAttr(text.CodeFgGreen, text.CodeBold)
	textBgGreenUnderline = text.NewAttr(text.CodeBgGreen, text.CodeUnderline)
	textBgMagentaItalic  = text.NewAttr(text.CodeBgMagenta, text.CodeItalic)
	textFgCyanBold       = text.NewAttr(text.CodeFgCyan, text.CodeBold)
)

// TextOptionASCII configures the ASCII table example.
var TextOptionASCII = []text.Option{
	text.WithStyle(text.StyleASCII),
	text.WithHeader(SimpleData.Header...),
}

// TextOptionSimple configures the simple text table example.
var TextOptionSimple = []text.Option{
	text.WithStyle(text.StyleLight),
	text.WithHeader(SimpleData.Header...),
}

// TextOptionCompact configures the compact text table example.
var TextOptionCompact = []text.Option{
	text.WithStyle(text.StyleColoredRounded),
	text.WithHeader(CompactData.Header...),
	text.WithRowspan(text.ScopeBody, text.Columns(0)),
	text.WithAlign(text.ScopeBody, text.Columns(3), text.AlignRight),
	text.WithCompact(),
	text.WithWidth(text.Columns(0), 11),
	text.WithWidth(text.Columns(1), 29),
	text.WithWidth(text.Columns(2), 36),
	text.WithWidth(text.Columns(3), 6),
}

// TextOptionRowspan configures the text table example with row spans.
var TextOptionRowspan = []text.Option{
	text.WithStyle(text.StyleColoredHeavy),
	text.WithHeader(RowspanData.Header...),
	text.WithRowspan(text.ScopeBody, text.Columns(0, 1, 2)),
	text.WithAlign(text.ScopeBody, text.Columns(4, 5), text.AlignRight),
	text.WithAutoFit(),
}

// TextOptionColspan configures the text table example with column spans.
var TextOptionColspan = []text.Option{
	text.WithStyle(text.StyleColoredLight),
	text.WithHeader(ColspanData.Header...),
	text.WithColspan(text.ScopeBody, text.Columns(0, 1, 2, 3)),
}

// TextOptionFooter configures the text table example with a footer.
var TextOptionFooter = []text.Option{
	text.WithStyle(text.StyleColoredLight),
	text.WithHeader(FooterData.Header...),
	text.WithFooter(FooterData.Footer),
	text.WithRowspan(text.ScopeBody, text.Columns(0)),
	text.WithColspan(text.ScopeFooter, text.Columns(0, 1, 2, 3)),
	text.WithAlign(text.ScopeBody, text.Columns(3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13), text.AlignRight),
	text.WithAlign(text.ScopeFooter, text.Columns(4, 5, 6, 7, 8, 9, 10, 11, 12, 13), text.AlignRight),
	text.WithAlign(text.ScopeFooter, text.Columns(0), text.AlignCenter),
	text.WithAutoFit(),
}

// TextOptionTransformer configures the transformed text table example.
var TextOptionTransformer = []text.Option{
	text.WithStyle(text.StyleColoredLight),
	text.WithHeader(FooterData.Header...),
	text.WithFooter(FooterData.Footer),
	text.WithRowspan(text.ScopeBody, text.Columns(0)),
	text.WithColspan(text.ScopeFooter, text.Columns(0, 1, 2, 3)),
	text.WithAlign(text.ScopeBody, text.Columns(3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13), text.AlignRight),
	text.WithAlign(text.ScopeFooter, text.Columns(4, 5, 6, 7, 8, 9, 10, 11, 12, 13), text.AlignRight),
	text.WithAlign(text.ScopeFooter, text.Columns(0), text.AlignCenter),
	text.WithAutoFit(),
	text.WithTransformer(text.Columns(5), func(v any) (string, *text.Attr) {
		n, ok := v.(int)
		if !ok {
			return "", nil
		}
		if n >= 3000 {
			return fmt.Sprintf("*%d", n), textFgRedBold
		}
		return "", nil
	}),
	text.WithTransformer(text.Columns(9), func(v any) (string, *text.Attr) {
		if n, ok := v.(int); ok && n >= 90 {
			return fmt.Sprintf("*%d", n), textFgYellowBold
		}
		return "", nil
	}),
	text.WithTransformer(text.Columns(13), func(v any) (string, *text.Attr) {
		if n, ok := v.(int); ok && n >= 60 {
			return fmt.Sprintf("*%d", n), textFgGreenBold
		}
		return "", nil
	}),
}

// TextOptionComplex configures the complex text table example.
var TextOptionComplex = []text.Option{
	text.WithStyle(text.StyleColoredDouble),
	text.WithHeader(ComplexData.Header...),
	text.WithAutoFit(),
	text.WithAttr(text.ScopeBody, text.Columns(8, 9, 10), text.ColorFgBlack),
	text.WithTransformer(text.Columns(5), func(v any) (string, *text.Attr) {
		values, ok := v.([]string)
		if !ok {
			return "", nil
		}
		tokens := make([]string, len(values))
		for i, s := range values {
			tokens[i] = fmt.Sprintf("%d/%d: %s", i+1, len(values), s)
		}
		return strings.Join(tokens, "\n"), textBgGreenUnderline
	}),
	text.WithTransformer(text.Columns(6), func(v any) (string, *text.Attr) {
		values, ok := v.([3]string)
		if !ok {
			return "", nil
		}
		tokens := make([]string, len(values))
		for i, s := range values {
			tokens[i] = fmt.Sprintf("%d/%d: %s", i+1, len(values), s)
		}
		return strings.Join(tokens, "\n"), textBgMagentaItalic
	}),
	text.WithTransformer(text.Columns(7), func(v any) (string, *text.Attr) {
		values, ok := v.([]int)
		if !ok {
			return "", nil
		}
		sum := 0
		for _, value := range values {
			sum += value
		}
		return fmt.Sprintf("sum=%d", sum), textFgCyanBold
	}),
	text.WithCaption("⚡️ Rendered by <github.com/nekrassov01/table/text>", text.CaptionDefault),
}

// TextOptionStackedHeader configures the text table example with a stacked header.
var TextOptionStackedHeader = []text.Option{
	text.WithStyle(text.StyleColoredLight),
	text.WithHeader(StackedHeaderData.Header...),
	text.WithRowspan(text.ScopeHeader, text.Columns(4)),
	text.WithColspan(text.ScopeHeader, text.Columns(0, 1, 2, 3)),
}
