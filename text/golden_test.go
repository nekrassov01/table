package text

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/nekrassov01/table/internal/param"
	"github.com/nekrassov01/table/internal/testutil"
)

func TestGolden_TableAlignMultiline(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), AlignCenter),
		WithHeader([]string{"ID", "Lines"}),
	)
	if err := tb.Render([][]any{
		{1, "short\nlonger-text\nx"},
		{2, "only"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_align_multiline", buf.Bytes())
}

func TestGolden_StreamAlignMultiline(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), AlignCenter),
		WithHeader([]string{"ID", "Lines"}),
	)
	if err := s.Render([]any{1, "short\nlonger-text\nx"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{2, "only"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_align_multiline", buf.Bytes())
}

func TestGolden_TableAlignScope(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"Item", "Amount"}),
		WithFooter(func() [][]string {
			return [][]string{{"Total", "300"}}
		}),
		WithAlign(ScopeHeader, Columns(0, 1), AlignCenter),
		WithAlign(ScopeBody|ScopeFooter, Columns(1), AlignRight),
	)
	if err := tb.Render([][]any{
		{"widget", 100},
		{"gadget", 200},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_align_scope", buf.Bytes())
}

func TestGolden_StreamAlignScope(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"Item", "Amount"}),
		WithFooter(func() [][]string {
			return [][]string{{"Total", "300"}}
		}),
		WithAlign(ScopeHeader, Columns(0, 1), AlignCenter),
		WithAlign(ScopeBody|ScopeFooter, Columns(1), AlignRight),
	)
	if err := s.Render([]any{"widget", 100}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"gadget", 200}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_align_scope", buf.Bytes())
}

func TestGolden_TableAllPlaceholder(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := tb.Render([][]any{
		{nil, nil, nil},
		{nil, nil, nil},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_all_placeholder", buf.Bytes())
}

func TestGolden_StreamAllPlaceholder(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := s.Render([]any{nil, nil, nil}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{nil, nil, nil}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_all_placeholder", buf.Bytes())
}

func TestGolden_TableAttr(t *testing.T) {
	restore := isTerminal
	isTerminal = func(io.Writer) bool { return true }
	t.Cleanup(func() { isTerminal = restore })
	var buf bytes.Buffer
	style := StyleLight
	style.Border.Attr = NewAttr(CodeFaint)
	style.Content = ContentStyle{
		Header:  NewAttr(CodeFgCyan, CodeBold),
		Body:    NewAttr(CodeItalic),
		Footer:  NewAttr(CodeFgYellow),
		Caption: NewAttr(CodeUnderline),
	}
	tb := NewTable(&buf,
		WithStyle(style),
		WithHeader([]string{"Name", "Value"}),
		WithFooter(func() [][]string {
			return [][]string{{"total", "3"}}
		}),
		WithCaption("caption", CaptionDefault),
		WithAttr(ScopeBody, Columns(1), ColorFgRed),
		WithTransformer(Columns(0), func(v any) (string, *Attr) {
			if s, ok := v.(string); ok && s == "bar" {
				return "", NewAttr(CodeFgGreen)
			}
			return "", nil
		}),
	)
	if err := tb.Render([][]any{
		{"foo", 1},
		{"bar", 2},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_attr", buf.Bytes())
}

func TestGolden_StreamAttr(t *testing.T) {
	restore := isTerminal
	isTerminal = func(io.Writer) bool { return true }
	t.Cleanup(func() { isTerminal = restore })
	var buf bytes.Buffer
	style := StyleLight
	style.Border.Attr = NewAttr(CodeFaint)
	style.Content = ContentStyle{
		Header:  NewAttr(CodeFgCyan, CodeBold),
		Body:    NewAttr(CodeItalic),
		Footer:  NewAttr(CodeFgYellow),
		Caption: NewAttr(CodeUnderline),
	}
	s := NewStream(&buf,
		WithStyle(style),
		WithHeader([]string{"Name", "Value"}),
		WithFooter(func() [][]string {
			return [][]string{{"total", "3"}}
		}),
		WithCaption("caption", CaptionDefault),
		WithAttr(ScopeBody, Columns(1), ColorFgRed),
		WithTransformer(Columns(0), func(v any) (string, *Attr) {
			if s, ok := v.(string); ok && s == "bar" {
				return "", NewAttr(CodeFgGreen)
			}
			return "", nil
		}),
	)
	if err := s.Render([]any{"foo", 1}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"bar", 2}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_attr", buf.Bytes())
}

func TestGolden_TableAttrBandPlaceholder(t *testing.T) {
	restore := isTerminal
	isTerminal = func(io.Writer) bool { return true }
	t.Cleanup(func() { isTerminal = restore })
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"Name", "Value"}, []string{"Sub"}),
		WithFooter(func() [][]string {
			return [][]string{{"total"}}
		}),
		WithAttr(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), ColorFgRed),
	)
	if err := tb.Render([][]any{{"foo", 1}, {"bar", 2}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_attr_band_placeholder", buf.Bytes())
}

func TestGolden_StreamAttrBandPlaceholder(t *testing.T) {
	restore := isTerminal
	isTerminal = func(io.Writer) bool { return true }
	t.Cleanup(func() { isTerminal = restore })
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"Name", "Value"}, []string{"Sub"}),
		WithFooter(func() [][]string {
			return [][]string{{"total"}}
		}),
		WithAttr(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), ColorFgRed),
	)
	for _, row := range [][]any{{"foo", 1}, {"bar", 2}} {
		if err := s.Render(row); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_attr_band_placeholder", buf.Bytes())
}

func TestGolden_TableAttrPlaceholder(t *testing.T) {
	restore := isTerminal
	isTerminal = func(io.Writer) bool { return true }
	t.Cleanup(func() { isTerminal = restore })
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"Name", "Value"}),
		WithAttr(ScopeBody, Columns(1), ColorFgRed),
	)
	if err := tb.Render([][]any{
		{"foo", 1},
		{"bar", nil},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_attr_placeholder", buf.Bytes())
}

func TestGolden_StreamAttrPlaceholder(t *testing.T) {
	restore := isTerminal
	isTerminal = func(io.Writer) bool { return true }
	t.Cleanup(func() { isTerminal = restore })
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"Name", "Value"}),
		WithAttr(ScopeBody, Columns(1), ColorFgRed),
	)
	for _, row := range [][]any{{"foo", 1}, {"bar", nil}} {
		if err := s.Render(row); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_attr_placeholder", buf.Bytes())
}

func TestGolden_TableAttrScope(t *testing.T) {
	restore := isTerminal
	isTerminal = func(io.Writer) bool { return true }
	t.Cleanup(func() { isTerminal = restore })
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"Name", "Value"}),
		WithFooter(func() [][]string {
			return [][]string{{"total", "3"}}
		}),
		WithAttr(ScopeBody, Columns(1), ColorFgRed),
		WithAttr(ScopeFooter, Columns(0, 1), ColorFgYellow),
	)
	if err := tb.Render([][]any{
		{"foo", 1},
		{"bar", 2},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_attr_scope", buf.Bytes())
}

func TestGolden_StreamAttrScope(t *testing.T) {
	restore := isTerminal
	isTerminal = func(io.Writer) bool { return true }
	t.Cleanup(func() { isTerminal = restore })
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"Name", "Value"}),
		WithFooter(func() [][]string {
			return [][]string{{"total", "3"}}
		}),
		WithAttr(ScopeBody, Columns(1), ColorFgRed),
		WithAttr(ScopeFooter, Columns(0, 1), ColorFgYellow),
	)
	for _, row := range [][]any{{"foo", 1}, {"bar", 2}} {
		if err := s.Render(row); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_attr_scope", buf.Bytes())
}

func TestGolden_TableAutoFitAlign(t *testing.T) {
	restoreW := terminalWidth
	terminalWidth = func(io.Writer) int { return 20 }
	t.Cleanup(func() { terminalWidth = restoreW })
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithAutoFit(),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), AlignRight),
	)
	if err := tb.Render([][]any{{"a fairly long value", "yy"}, {"another long one", "z"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_autofit_align", buf.Bytes())
}

func TestGolden_StreamAutoFitAlign(t *testing.T) {
	restoreW := terminalWidth
	terminalWidth = func(io.Writer) int { return 20 }
	t.Cleanup(func() { terminalWidth = restoreW })
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithAutoFit(),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), AlignRight),
	)
	for _, r := range [][]any{{"a fairly long value", "yy"}, {"another long one", "z"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_autofit_align", buf.Bytes())
}

func TestGolden_TableAutoFitAttr(t *testing.T) {
	restore := isTerminal
	isTerminal = func(io.Writer) bool { return true }
	t.Cleanup(func() { isTerminal = restore })
	restoreW := terminalWidth
	terminalWidth = func(io.Writer) int { return 20 }
	t.Cleanup(func() { terminalWidth = restoreW })
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithAutoFit(),
		WithAttr(ScopeBody, Columns(1), ColorFgRed),
	)
	if err := tb.Render([][]any{{"a fairly long value", "x"}, {"another long one", "z"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_autofit_attr", buf.Bytes())
}

func TestGolden_StreamAutoFitAttr(t *testing.T) {
	restore := isTerminal
	isTerminal = func(io.Writer) bool { return true }
	t.Cleanup(func() { isTerminal = restore })
	restoreW := terminalWidth
	terminalWidth = func(io.Writer) int { return 20 }
	t.Cleanup(func() { terminalWidth = restoreW })
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithAutoFit(),
		WithAttr(ScopeBody, Columns(1), ColorFgRed),
	)
	for _, r := range [][]any{{"a fairly long value", "x"}, {"another long one", "z"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_autofit_attr", buf.Bytes())
}

func TestGolden_TableAutoFitCaption(t *testing.T) {
	restoreW := terminalWidth
	terminalWidth = func(io.Writer) int { return 20 }
	t.Cleanup(func() { terminalWidth = restoreW })
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithAutoFit(),
		WithCaption("cap", CaptionBottom),
	)
	if err := tb.Render([][]any{{"a fairly long value", "x"}, {"another long one", "z"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_autofit_caption", buf.Bytes())
}

func TestGolden_StreamAutoFitCaption(t *testing.T) {
	restoreW := terminalWidth
	terminalWidth = func(io.Writer) int { return 20 }
	t.Cleanup(func() { terminalWidth = restoreW })
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithAutoFit(),
		WithCaption("cap", CaptionBottom),
	)
	for _, r := range [][]any{{"a fairly long value", "x"}, {"another long one", "z"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_autofit_caption", buf.Bytes())
}

func TestGolden_TableAutoFitFits(t *testing.T) {
	restore := terminalWidth
	terminalWidth = func(io.Writer) int { return 80 }
	t.Cleanup(func() { terminalWidth = restore })
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithAutoFit(),
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]any{{"xxxxx", "yyyyy"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_autofit_fits", buf.Bytes())
}

func TestGolden_StreamAutoFitFits(t *testing.T) {
	restore := terminalWidth
	terminalWidth = func(io.Writer) int { return 80 }
	t.Cleanup(func() { terminalWidth = restore })
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithAutoFit(),
		WithHeader([]string{"A", "B"}),
	)
	if err := s.Render([]any{"xxxxx", "yyyyy"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_autofit_fits", buf.Bytes())
}

func TestGolden_TableAutoFitPadding(t *testing.T) {
	restoreW := terminalWidth
	terminalWidth = func(io.Writer) int { return 20 }
	t.Cleanup(func() { terminalWidth = restoreW })
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithAutoFit(),
		WithPadding(Columns(1), 3, 0),
	)
	if err := tb.Render([][]any{{"a fairly long value", "x"}, {"another long one", "z"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_autofit_padding", buf.Bytes())
}

func TestGolden_StreamAutoFitPadding(t *testing.T) {
	restoreW := terminalWidth
	terminalWidth = func(io.Writer) int { return 20 }
	t.Cleanup(func() { terminalWidth = restoreW })
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithAutoFit(),
		WithPadding(Columns(1), 3, 0),
	)
	for _, r := range [][]any{{"a fairly long value", "x"}, {"another long one", "z"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_autofit_padding", buf.Bytes())
}

func TestGolden_TableAutoFitPlaceholder(t *testing.T) {
	restoreW := terminalWidth
	terminalWidth = func(io.Writer) int { return 20 }
	t.Cleanup(func() { terminalWidth = restoreW })
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithAutoFit(),
		WithPlaceholder("-"),
	)
	if err := tb.Render([][]any{{"a fairly long value", nil}, {"another long one", "z"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_autofit_placeholder", buf.Bytes())
}

func TestGolden_StreamAutoFitPlaceholder(t *testing.T) {
	restoreW := terminalWidth
	terminalWidth = func(io.Writer) int { return 20 }
	t.Cleanup(func() { terminalWidth = restoreW })
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithAutoFit(),
		WithPlaceholder("-"),
	)
	for _, r := range [][]any{{"a fairly long value", nil}, {"another long one", "z"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_autofit_placeholder", buf.Bytes())
}

func TestGolden_TableAutoFitRowspan(t *testing.T) {
	restore := terminalWidth
	terminalWidth = func(io.Writer) int { return 30 }
	t.Cleanup(func() { terminalWidth = restore })
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"Name", "Value"}),
		WithFooter(func() [][]string {
			return [][]string{{"total", "3"}}
		}),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
		WithAutoFit(),
	)
	if err := tb.Render([][]any{
		{"same-and-rather-long-label", 1},
		{"same-and-rather-long-label", 2},
		{"other", 3},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_autofit_rowspan", buf.Bytes())
}

func TestGolden_StreamAutoFitRowspan(t *testing.T) {
	restore := terminalWidth
	terminalWidth = func(io.Writer) int { return 30 }
	t.Cleanup(func() { terminalWidth = restore })
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"Name", "Value"}),
		WithFooter(func() [][]string {
			return [][]string{{"total", "3"}}
		}),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
		WithAutoFit(),
	)
	for _, row := range [][]any{
		{"same-and-rather-long-label", 1},
		{"same-and-rather-long-label", 2},
		{"other", 3},
	} {
		if err := s.Render(row); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_autofit_rowspan", buf.Bytes())
}

func TestGolden_TableAutoFitTransformer(t *testing.T) {
	restoreW := terminalWidth
	terminalWidth = func(io.Writer) int { return 20 }
	t.Cleanup(func() { terminalWidth = restoreW })
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithAutoFit(),
		WithTransformer(Columns(1), func(v any) (string, *Attr) {
			if s, ok := v.(string); ok && s == "raw" {
				return "T", nil
			}
			return "", nil
		}),
	)
	if err := tb.Render([][]any{{"a fairly long value", "raw"}, {"another long one", "z"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_autofit_transformer", buf.Bytes())
}

func TestGolden_StreamAutoFitTransformer(t *testing.T) {
	restoreW := terminalWidth
	terminalWidth = func(io.Writer) int { return 20 }
	t.Cleanup(func() { terminalWidth = restoreW })
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithAutoFit(),
		WithTransformer(Columns(1), func(v any) (string, *Attr) {
			if s, ok := v.(string); ok && s == "raw" {
				return "T", nil
			}
			return "", nil
		}),
	)
	for _, r := range [][]any{{"a fairly long value", "raw"}, {"another long one", "z"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_autofit_transformer", buf.Bytes())
}

func TestGolden_TableAutoFitTruncate(t *testing.T) {
	restore := terminalWidth
	terminalWidth = func(io.Writer) int { return 24 }
	t.Cleanup(func() { terminalWidth = restore })
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithAutoFit(),
		WithTruncate(Columns(1)),
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]any{{"xxxxx", strings.Repeat("y", 20)}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_autofit_truncate", buf.Bytes())
}

func TestGolden_StreamAutoFitTruncate(t *testing.T) {
	restore := terminalWidth
	terminalWidth = func(io.Writer) int { return 24 }
	t.Cleanup(func() { terminalWidth = restore })
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithAutoFit(),
		WithTruncate(Columns(1)),
		WithHeader([]string{"A", "B"}),
	)
	if err := s.Render([]any{"xxxxx", strings.Repeat("y", 20)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_autofit_truncate", buf.Bytes())
}

func TestGolden_TableAutoFitWithWidth(t *testing.T) {
	restore := terminalWidth
	terminalWidth = func(io.Writer) int { return 24 }
	t.Cleanup(func() { terminalWidth = restore })
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithAutoFit(),
		WithWidth(Columns(0), 3),
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]any{{"xxxxx", strings.Repeat("y", 20)}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_autofit_with_width", buf.Bytes())
}

func TestGolden_StreamAutoFitWithWidth(t *testing.T) {
	restore := terminalWidth
	terminalWidth = func(io.Writer) int { return 24 }
	t.Cleanup(func() { terminalWidth = restore })
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithAutoFit(),
		WithWidth(Columns(0), 3),
		WithHeader([]string{"A", "B"}),
	)
	if err := s.Render([]any{"xxxxx", strings.Repeat("y", 20)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_autofit_with_width", buf.Bytes())
}

func TestGolden_TableBasic(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"Name", "Value"}),
	)
	if err := tb.Render([][]any{
		{"foo", 1},
		{"bar", 2},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_basic", buf.Bytes())
}
func TestGolden_StreamBasic(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"Name", "Value"}),
	)
	if err := s.Render([]any{"foo", 1}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"bar", 2}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_basic", buf.Bytes())
}
func TestGolden_TableCaption(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithCaption("Table 1", CaptionDefault),
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]any{
		{"x", 1},
		{"y", 2},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_caption", buf.Bytes())
}

func TestGolden_StreamCaption(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithCaption("Table 1", CaptionDefault),
		WithHeader([]string{"A", "B"}),
	)
	if err := s.Render([]any{"x", 1}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"y", 2}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_caption", buf.Bytes())
}

func TestGolden_TableCaptionTop(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithCaption("summary", CaptionTop),
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]any{{"x", "y"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_caption_top", buf.Bytes())
}

func TestGolden_StreamCaptionTop(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithCaption("summary", CaptionTop),
		WithHeader([]string{"A", "B"}),
	)
	if err := s.Render([]any{"x", "y"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_caption_top", buf.Bytes())
}

func TestGolden_TableColspan(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B", "C"}),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1, 2)),
	)
	if err := tb.Render([][]any{
		{"x", "x", "y"},
		{"p", "q", "q"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_colspan", buf.Bytes())
}

func TestGolden_StreamColspan(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B", "C"}),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1, 2)),
	)
	if err := s.Render([]any{"x", "x", "y"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"p", "q", "q"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_colspan", buf.Bytes())
}

func TestGolden_TableColspanAlign(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1, 2)),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(0), AlignCenter),
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := tb.Render([][]any{{"x", "x", "y"}, {"p", "q", "q"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_colspan_align", buf.Bytes())
}

func TestGolden_StreamColspanAlign(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1, 2)),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(0), AlignCenter),
		WithHeader([]string{"A", "B", "C"}),
	)
	for _, r := range [][]any{{"x", "x", "y"}, {"p", "q", "q"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_colspan_align", buf.Bytes())
}

func TestGolden_TableColspanAttr(t *testing.T) {
	restore := isTerminal
	isTerminal = func(io.Writer) bool { return true }
	t.Cleanup(func() { isTerminal = restore })
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1)),
		WithAttr(ScopeBody, Columns(1), ColorFgRed),
	)
	if err := tb.Render([][]any{{"s", "s"}, {"s", "s"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_colspan_attr", buf.Bytes())
}

func TestGolden_StreamColspanAttr(t *testing.T) {
	restore := isTerminal
	isTerminal = func(io.Writer) bool { return true }
	t.Cleanup(func() { isTerminal = restore })
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1)),
		WithAttr(ScopeBody, Columns(1), ColorFgRed),
	)
	for _, r := range [][]any{{"s", "s"}, {"s", "s"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_colspan_attr", buf.Bytes())
}

func TestGolden_TableColspanCaption(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1)),
		WithCaption("cap", CaptionBottom),
	)
	if err := tb.Render([][]any{{"s", "s"}, {"s", "s"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_colspan_caption", buf.Bytes())
}

func TestGolden_StreamColspanCaption(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1)),
		WithCaption("cap", CaptionBottom),
	)
	for _, r := range [][]any{{"s", "s"}, {"s", "s"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_colspan_caption", buf.Bytes())
}

func TestGolden_TableColspanEdges(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B", "C", "D"}),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1)),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(2, 3)),
	)
	if err := tb.Render([][]any{
		{"x", "x", "y", "y"},
		{"p", "q", "r", "r"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_colspan_edges", buf.Bytes())
}

func TestGolden_StreamColspanEdges(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B", "C", "D"}),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1)),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(2, 3)),
	)
	if err := s.Render([]any{"x", "x", "y", "y"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"p", "q", "r", "r"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_colspan_edges", buf.Bytes())
}

func TestGolden_TableColspanPadding(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1, 2)),
		WithPadding(Columns(1), 3, 0),
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := tb.Render([][]any{{"x", "x", "y"}, {"p", "q", "q"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_colspan_padding", buf.Bytes())
}

func TestGolden_StreamColspanPadding(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1, 2)),
		WithPadding(Columns(1), 3, 0),
		WithHeader([]string{"A", "B", "C"}),
	)
	for _, r := range [][]any{{"x", "x", "y"}, {"p", "q", "q"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_colspan_padding", buf.Bytes())
}

func TestGolden_TableColspanPlaceholder(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1)),
		WithPlaceholder("-"),
	)
	if err := tb.Render([][]any{{"x", nil}, {nil, "q"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_colspan_placeholder", buf.Bytes())
}

func TestGolden_StreamColspanPlaceholder(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1)),
		WithPlaceholder("-"),
	)
	for _, r := range [][]any{{"x", nil}, {nil, "q"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_colspan_placeholder", buf.Bytes())
}

func TestGolden_TableColspanScope(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"Group", "Group", "Value"}),
		WithColspan(ScopeBody, Columns(0, 1)),
	)
	if err := tb.Render([][]any{
		{"A", "A", 1},
		{"B", "B", 2},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_colspan_scope", buf.Bytes())
}

func TestGolden_StreamColspanScope(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"Group", "Group", "Value"}),
		WithColspan(ScopeBody, Columns(0, 1)),
	)
	if err := s.Render([]any{"A", "A", 1}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"B", "B", 2}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_colspan_scope", buf.Bytes())
}

func TestGolden_TableColspanTransformerAttr(t *testing.T) {
	restore := isTerminal
	isTerminal = func(io.Writer) bool { return true }
	t.Cleanup(func() { isTerminal = restore })
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B", "C"}),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(1, 2)),
		WithTransformer(Columns(1, 2), func(_ any) (string, *Attr) {
			return "", ColorFgRed
		}),
	)
	if err := tb.Render([][]any{
		{"g", "x", "x"},
		{"h", "y", "z"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_colspan_transformer_attr", buf.Bytes())
}

func TestGolden_StreamColspanTransformerAttr(t *testing.T) {
	restore := isTerminal
	isTerminal = func(io.Writer) bool { return true }
	t.Cleanup(func() { isTerminal = restore })
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B", "C"}),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(1, 2)),
		WithTransformer(Columns(1, 2), func(_ any) (string, *Attr) {
			return "", ColorFgRed
		}),
	)
	if err := s.Render([]any{"g", "x", "x"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"h", "y", "z"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_colspan_transformer_attr", buf.Bytes())
}

func TestGolden_TableColspanTruncate(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1)),
		WithWidth(Columns(0), 5),
		WithTruncate(Columns(0)),
	)
	if err := tb.Render([][]any{{"a long value", "a long value"}, {"a long value", "q"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_colspan_truncate", buf.Bytes())
}

func TestGolden_StreamColspanTruncate(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1)),
		WithWidth(Columns(0), 5),
		WithTruncate(Columns(0)),
	)
	for _, r := range [][]any{{"a long value", "a long value"}, {"a long value", "q"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_colspan_truncate", buf.Bytes())
}

func TestGolden_TableCompact(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithCompact(),
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]any{
		{"x", 1},
		{"y", 2},
		{"z", 3},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_compact", buf.Bytes())
}

func TestGolden_StreamCompact(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithCompact(),
		WithHeader([]string{"A", "B"}),
	)
	if err := s.Render([]any{"x", 1}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"y", 2}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"z", 3}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_compact", buf.Bytes())
}

func TestGolden_TableCompactAlign(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithCompact(),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), AlignRight),
		WithHeader([]string{"Name", "Score"}),
	)
	if err := tb.Render([][]any{
		{"alice", 100},
		{"bob", 99},
		{"carol", 1},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_compact_align", buf.Bytes())
}

func TestGolden_StreamCompactAlign(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithCompact(),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), AlignRight),
		WithHeader([]string{"Name", "Score"}),
	)
	if err := s.Render([]any{"alice", 100}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"bob", 99}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"carol", 1}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_compact_align", buf.Bytes())
}

func TestGolden_TableCompactAttr(t *testing.T) {
	restore := isTerminal
	isTerminal = func(io.Writer) bool { return true }
	t.Cleanup(func() { isTerminal = restore })
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithCompact(),
		WithAttr(ScopeBody, Columns(1), ColorFgRed),
	)
	if err := tb.Render([][]any{{"s", "s"}, {"s", "s"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_compact_attr", buf.Bytes())
}

func TestGolden_StreamCompactAttr(t *testing.T) {
	restore := isTerminal
	isTerminal = func(io.Writer) bool { return true }
	t.Cleanup(func() { isTerminal = restore })
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithCompact(),
		WithAttr(ScopeBody, Columns(1), ColorFgRed),
	)
	for _, r := range [][]any{{"s", "s"}, {"s", "s"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_compact_attr", buf.Bytes())
}

func TestGolden_TableCompactAutoFit(t *testing.T) {
	restore := terminalWidth
	terminalWidth = func(io.Writer) int { return 20 }
	t.Cleanup(func() { terminalWidth = restore })
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithCompact(),
		WithAutoFit(),
	)
	if err := tb.Render([][]any{{"a fairly long value", "a fairly long value"}, {"another long one", "z"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_compact_autofit", buf.Bytes())
}

func TestGolden_StreamCompactAutoFit(t *testing.T) {
	restore := terminalWidth
	terminalWidth = func(io.Writer) int { return 20 }
	t.Cleanup(func() { terminalWidth = restore })
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithCompact(),
		WithAutoFit(),
	)
	for _, r := range [][]any{{"a fairly long value", "a fairly long value"}, {"another long one", "z"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_compact_autofit", buf.Bytes())
}

func TestGolden_TableCompactCaption(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithCompact(),
		WithCaption("cap", CaptionBottom),
	)
	if err := tb.Render([][]any{{"s", "s"}, {"s", "s"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_compact_caption", buf.Bytes())
}

func TestGolden_StreamCompactCaption(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithCompact(),
		WithCaption("cap", CaptionBottom),
	)
	for _, r := range [][]any{{"s", "s"}, {"s", "s"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_compact_caption", buf.Bytes())
}

func TestGolden_TableCompactColspan(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithCompact(),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1, 2)),
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := tb.Render([][]any{
		{"x", "x", "y"},
		{"p", "p", "q"},
		{"m", "n", "n"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_compact_colspan", buf.Bytes())
}

func TestGolden_StreamCompactColspan(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithCompact(),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1, 2)),
		WithHeader([]string{"A", "B", "C"}),
	)
	for _, r := range [][]any{
		{"x", "x", "y"},
		{"p", "p", "q"},
		{"m", "n", "n"},
	} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_compact_colspan", buf.Bytes())
}

func TestGolden_TableCompactFooter(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithCompact(),
		WithFooter(func() [][]string {
			return [][]string{{"t", "u"}}
		}),
	)
	if err := tb.Render([][]any{{"s", "s"}, {"s", "s"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_compact_footer", buf.Bytes())
}

func TestGolden_StreamCompactFooter(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithCompact(),
		WithFooter(func() [][]string {
			return [][]string{{"t", "u"}}
		}),
	)
	for _, r := range [][]any{{"s", "s"}, {"s", "s"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_compact_footer", buf.Bytes())
}

func TestGolden_TableCompactMultiline(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithCompact(),
		WithHeader([]string{"ID", "Data"}),
	)
	if err := tb.Render([][]any{
		{1, "a\nb"},
		{2, "c"},
		{3, "d\ne\nf"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_compact_multiline", buf.Bytes())
}

func TestGolden_StreamCompactMultiline(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithCompact(),
		WithHeader([]string{"ID", "Data"}),
	)
	if err := s.Render([]any{1, "a\nb"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{2, "c"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{3, "d\ne\nf"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_compact_multiline", buf.Bytes())
}

func TestGolden_TableCompactPadding(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithCompact(),
		WithPadding(Columns(1), 2, 2),
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]any{
		{"x", 1},
		{"y", 2},
		{"z", 3},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_compact_padding", buf.Bytes())
}

func TestGolden_StreamCompactPadding(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithCompact(),
		WithPadding(Columns(1), 2, 2),
		WithHeader([]string{"A", "B"}),
	)
	if err := s.Render([]any{"x", 1}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"y", 2}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"z", 3}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_compact_padding", buf.Bytes())
}

func TestGolden_TableCompactPlaceholder(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithCompact(),
		WithPlaceholder("-"),
	)
	if err := tb.Render([][]any{{"x", nil}, {nil, "q"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_compact_placeholder", buf.Bytes())
}

func TestGolden_StreamCompactPlaceholder(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithCompact(),
		WithPlaceholder("-"),
	)
	for _, r := range [][]any{{"x", nil}, {nil, "q"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_compact_placeholder", buf.Bytes())
}

func TestGolden_TableCompactRowspan(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithCompact(),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
		WithHeader([]string{"Group", "Item"}),
	)
	if err := tb.Render([][]any{
		{"A", "x"},
		{"A", "y"},
		{"B", "z"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_compact_rowspan", buf.Bytes())
}

func TestGolden_StreamCompactRowspan(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithCompact(),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
		WithHeader([]string{"Group", "Item"}),
	)
	if err := s.Render([]any{"A", "x"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"A", "y"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"B", "z"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_compact_rowspan", buf.Bytes())
}

func TestGolden_TableCompactTransformer(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithCompact(),
		WithTransformer(Columns(1), func(v any) (string, *Attr) {
			if s, ok := v.(string); ok && s == "raw" {
				return "T", nil
			}
			return "", nil
		}),
	)
	if err := tb.Render([][]any{{"x", "raw"}, {"p", "q"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_compact_transformer", buf.Bytes())
}

func TestGolden_StreamCompactTransformer(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithCompact(),
		WithTransformer(Columns(1), func(v any) (string, *Attr) {
			if s, ok := v.(string); ok && s == "raw" {
				return "T", nil
			}
			return "", nil
		}),
	)
	for _, r := range [][]any{{"x", "raw"}, {"p", "q"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_compact_transformer", buf.Bytes())
}

func TestGolden_TableCompactTruncate(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithCompact(),
		WithWidth(Columns(1), 5),
		WithTruncate(Columns(1)),
	)
	if err := tb.Render([][]any{{"x", "a long value"}, {"y", "short"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_compact_truncate", buf.Bytes())
}

func TestGolden_StreamCompactTruncate(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithCompact(),
		WithWidth(Columns(1), 5),
		WithTruncate(Columns(1)),
	)
	for _, r := range [][]any{{"x", "a long value"}, {"y", "short"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_compact_truncate", buf.Bytes())
}

func TestGolden_TableCompactWidth(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithCompact(),
		WithWidth(Columns(1), 5),
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]any{{"x", "a long value"}, {"y", "short"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_compact_width", buf.Bytes())
}

func TestGolden_StreamCompactWidth(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithCompact(),
		WithWidth(Columns(1), 5),
		WithHeader([]string{"A", "B"}),
	)
	for _, r := range [][]any{{"x", "a long value"}, {"y", "short"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_compact_width", buf.Bytes())
}

func TestGolden_TableControlChars(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := tb.Render([][]any{
		{"a\tb", "c\vd", "e\x00f"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_control_chars", buf.Bytes())
}

func TestGolden_StreamControlChars(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := s.Render([]any{"a\tb", "c\vd", "e\x00f"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_control_chars", buf.Bytes())
}

func TestGolden_TableEmoji(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := tb.Render([][]any{
		{"\U0001F600", "\U0001F469\u200D\U0001F4BB", "e\u0301"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_emoji", buf.Bytes())
}

func TestGolden_StreamEmoji(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := s.Render([]any{"\U0001F600", "\U0001F469\u200D\U0001F4BB", "e\u0301"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_emoji", buf.Bytes())
}

func TestGolden_TableEmpty(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
	)
	if err := tb.Render(nil); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_empty", buf.Bytes())
}

func TestGolden_StreamEmpty(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
	)
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_empty", buf.Bytes())
}

func TestGolden_TableEmptyHeaderLabel(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"", "B", ""}),
	)
	if err := tb.Render([][]any{
		{"x", "y", "z"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_empty_header_label", buf.Bytes())
}

func TestGolden_StreamEmptyHeaderLabel(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"", "B", ""}),
	)
	if err := s.Render([]any{"x", "y", "z"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_empty_header_label", buf.Bytes())
}

func TestGolden_TableFooter(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B", "C"}),
		WithFooter(func() [][]string {
			return [][]string{{"sum", "", "9"}}
		}),
	)
	if err := tb.Render([][]any{
		{"x", "y", 1},
		{"p", "q", 8},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_footer", buf.Bytes())
}

func TestGolden_StreamFooter(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B", "C"}),
		WithFooter(func() [][]string {
			return [][]string{{"sum", "", "9"}}
		}),
	)
	if err := s.Render([]any{"x", "y", 1}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"p", "q", 8}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_footer", buf.Bytes())
}

func TestGolden_TableFooterEmptyBody(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithFooter(func() [][]string {
			return [][]string{{"f1", "f2"}}
		}),
	)
	if err := tb.Render(nil); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_footer_empty_body", buf.Bytes())
}

func TestGolden_StreamFooterEmptyBody(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithFooter(func() [][]string {
			return [][]string{{"f1", "f2"}}
		}),
	)
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_footer_empty_body", buf.Bytes())
}

func TestGolden_TableFooterNoHeader(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithFooter(func() [][]string {
			return [][]string{{"Total", "300"}}
		}),
	)
	if err := tb.Render(nil); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_footer_no_header", buf.Bytes())
}

func TestGolden_StreamFooterNoHeader(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithFooter(func() [][]string {
			return [][]string{{"Total", "300"}}
		}),
	)
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_footer_no_header", buf.Bytes())
}

func TestGolden_TableFooterPadding(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithFooter(func() [][]string {
			return [][]string{{"t", "u"}}
		}),
		WithPadding(Columns(1), 3, 0),
	)
	if err := tb.Render([][]any{{"x", "y"}, {"p", "q"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_footer_padding", buf.Bytes())
}

func TestGolden_StreamFooterPadding(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithFooter(func() [][]string {
			return [][]string{{"t", "u"}}
		}),
		WithPadding(Columns(1), 3, 0),
	)
	for _, r := range [][]any{{"x", "y"}, {"p", "q"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_footer_padding", buf.Bytes())
}

func TestGolden_TableFooterPlaceholder(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithFooter(func() [][]string {
			return [][]string{{"t", "u"}}
		}),
		WithPlaceholder("-"),
	)
	if err := tb.Render([][]any{{"x", nil}, {nil, "q"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_footer_placeholder", buf.Bytes())
}

func TestGolden_StreamFooterPlaceholder(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithFooter(func() [][]string {
			return [][]string{{"t", "u"}}
		}),
		WithPlaceholder("-"),
	)
	for _, r := range [][]any{{"x", nil}, {nil, "q"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_footer_placeholder", buf.Bytes())
}

func TestGolden_TableFooterTruncate(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithFooter(func() [][]string {
			return [][]string{{"t", "u"}}
		}),
		WithWidth(Columns(0), 5),
		WithTruncate(Columns(0)),
	)
	if err := tb.Render([][]any{{"a long value", "x"}, {"y", "z"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_footer_truncate", buf.Bytes())
}

func TestGolden_StreamFooterTruncate(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithFooter(func() [][]string {
			return [][]string{{"t", "u"}}
		}),
		WithWidth(Columns(0), 5),
		WithTruncate(Columns(0)),
	)
	for _, r := range [][]any{{"a long value", "x"}, {"y", "z"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_footer_truncate", buf.Bytes())
}

func TestGolden_TableHeaderOnly(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := tb.Render(nil); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_header_only", buf.Bytes())
}

func TestGolden_StreamHeaderOnly(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_header_only", buf.Bytes())
}

func TestGolden_TableHeaderRowspan(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader(
			[]string{"Region", "Env", "Tier"},
			[]string{"Region", "Env", "Kind"},
		),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1, 2)),
	)
	if err := tb.Render([][]any{
		{"jp", "prod", "web"},
		{"jp", "prod", "db"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_header_rowspan", buf.Bytes())
}

func TestGolden_StreamHeaderRowspan(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader(
			[]string{"Region", "Env", "Tier"},
			[]string{"Region", "Env", "Kind"},
		),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1, 2)),
	)
	for _, row := range [][]any{{"jp", "prod", "web"}, {"jp", "prod", "db"}} {
		if err := s.Render(row); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_header_rowspan", buf.Bytes())
}

func TestGolden_TableHeaderRowspanColspan(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader(
			[]string{"Group", "Group", "Val"},
			[]string{"Group", "Group", "Val"},
		),
		WithFooter(func() [][]string {
			return [][]string{{"Sum", "Sum", "3"}, {"Sum", "Sum", "3"}}
		}),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1, 2)),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1, 2)),
	)
	if err := tb.Render([][]any{
		{"a", "b", 1},
		{"a", "b", 2},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_header_rowspan_colspan", buf.Bytes())
}

func TestGolden_StreamHeaderRowspanColspan(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader(
			[]string{"Group", "Group", "Val"},
			[]string{"Group", "Group", "Val"},
		),
		WithFooter(func() [][]string {
			return [][]string{{"Sum", "Sum", "3"}, {"Sum", "Sum", "3"}}
		}),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1, 2)),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1, 2)),
	)
	for _, row := range [][]any{{"a", "b", 1}, {"a", "b", 2}} {
		if err := s.Render(row); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_header_rowspan_colspan", buf.Bytes())
}

func TestGolden_TableHeaderWiderThanRows(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B", "C", "D", "E"}),
	)
	if err := tb.Render([][]any{
		{"x", "y"},
		{"p", "q"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_header_wider_than_rows", buf.Bytes())
}

func TestGolden_StreamHeaderWiderThanRows(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B", "C", "D", "E"}),
	)
	if err := s.Render([]any{"x", "y"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"p", "q"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_header_wider_than_rows", buf.Bytes())
}

func TestGolden_TableIndexWidth(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithIndexWidth(5),
	)
	if err := tb.Render([][]any{{"x", "y"}, {"p", "q"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_index_width", buf.Bytes())
}

func TestGolden_StreamIndexWidth(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithIndexWidth(5),
	)
	for _, r := range [][]any{{"x", "y"}, {"p", "q"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_index_width", buf.Bytes())
}

func TestGolden_TableInvalidUtf8(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := tb.Render([][]any{
		{"a\xffb", "\xfe", "ok"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_invalid_utf8", buf.Bytes())
}

func TestGolden_StreamInvalidUtf8(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := s.Render([]any{"a\xffb", "\xfe", "ok"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_invalid_utf8", buf.Bytes())
}

func TestGolden_TableMaxWidth(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithWidth(Columns(1), 15),
		WithHeader([]string{"Name", "Description"}),
	)
	if err := tb.Render([][]any{
		{"alice", "short"},
		{"bob", "this is a longer description"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_max_width", buf.Bytes())
}

func TestGolden_StreamMaxWidth(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithWidth(Columns(1), 15),
		WithHeader([]string{"Name", "Description"}),
	)
	if err := s.Render([]any{"alice", "short"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"bob", "this is a longer description"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_max_width", buf.Bytes())
}

func TestGolden_TableMaxWidthCJK(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithWidth(Columns(1), 8),
		WithHeader([]string{"Label", "日本語"}),
	)
	if err := tb.Render([][]any{
		{"x", "あいうえおかきくけこ"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_max_width_cjk", buf.Bytes())
}

func TestGolden_StreamMaxWidthCJK(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithWidth(Columns(1), 8),
		WithHeader([]string{"Label", "日本語"}),
	)
	if err := s.Render([]any{"x", "あいうえおかきくけこ"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_max_width_cjk", buf.Bytes())
}

func TestGolden_TableMaxWidthHeader(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithWidth(Columns(0), 5),
		WithWidth(Columns(1), 5),
		WithHeader([]string{"LongHeader", "Value"}),
	)
	if err := tb.Render([][]any{
		{"ab", "xy"},
		{"cdef", "z"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_max_width_header", buf.Bytes())
}

func TestGolden_StreamMaxWidthHeader(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithWidth(Columns(0), 5),
		WithWidth(Columns(1), 5),
		WithHeader([]string{"LongHeader", "Value"}),
	)
	if err := s.Render([]any{"ab", "xy"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"cdef", "z"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_max_width_header", buf.Bytes())
}

func TestGolden_TableMaxWidthMultiline(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithWidth(Columns(1), 10),
		WithHeader([]string{"ID", "Data"}),
	)
	if err := tb.Render([][]any{
		{1, "line1\nthis-is-a-long-line2"},
		{2, "short"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_max_width_multiline", buf.Bytes())
}

func TestGolden_StreamMaxWidthMultiline(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithWidth(Columns(1), 10),
		WithHeader([]string{"ID", "Data"}),
	)
	if err := s.Render([]any{1, "line1\nthis-is-a-long-line2"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{2, "short"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_max_width_multiline", buf.Bytes())
}

func TestGolden_TableMultilineCell(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"single", "double", "triple"}),
	)
	if err := tb.Render([][]any{
		{"one", "a\nb", "x\ny\nz"},
		{"two\nlines", "p", "q\nr"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_multiline_cell", buf.Bytes())
}

func TestGolden_StreamMultilineCell(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"single", "double", "triple"}),
	)
	if err := s.Render([]any{"one", "a\nb", "x\ny\nz"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"two\nlines", "p", "q\nr"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_multiline_cell", buf.Bytes())
}

func TestGolden_TableMultilineCRLF(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"line"}),
	)
	if err := tb.Render([][]any{{"a\r\nb\nc\rd"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_multiline_crlf", buf.Bytes())
}

func TestGolden_StreamMultilineCRLF(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"line"}),
	)
	if err := s.Render([]any{"a\r\nb\nc\rd"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_multiline_crlf", buf.Bytes())
}

func TestGolden_TableMultilineHeader(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"single", "two\nlines", "three\nlines\nhere"}),
	)
	if err := tb.Render([][]any{
		{"a", "b", "c"},
		{"d", "e", "f"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_multiline_header", buf.Bytes())
}

func TestGolden_StreamMultilineHeader(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"single", "two\nlines", "three\nlines\nhere"}),
	)
	if err := s.Render([]any{"a", "b", "c"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"d", "e", "f"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_multiline_header", buf.Bytes())
}

func TestGolden_TableNilInNumeric(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"Name", "Score"}),
	)
	if err := tb.Render([][]any{
		{"alice", 100},
		{"bob", nil},
		{"carol", 99999},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_nil_in_numeric", buf.Bytes())
}

func TestGolden_StreamNilInNumeric(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"Name", "Score"}),
	)
	if err := s.Render([]any{"alice", 100}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"bob", nil}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"carol", 99999}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_nil_in_numeric", buf.Bytes())
}

func TestGolden_TablePadding(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithPadding(Columns(1), 3, 3),
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]any{
		{"x", 1},
		{"y", 2},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_padding", buf.Bytes())
}

func TestGolden_StreamPadding(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithPadding(Columns(1), 3, 3),
		WithHeader([]string{"A", "B"}),
	)
	if err := s.Render([]any{"x", 1}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"y", 2}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_padding", buf.Bytes())
}

func TestGolden_TablePaddingMultiline(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithPadding(Columns(1), 3, 3),
		WithHeader([]string{"ID", "Lines"}),
	)
	if err := tb.Render([][]any{
		{1, "a\nb\nc"},
		{2, "x"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_padding_multiline", buf.Bytes())
}

func TestGolden_StreamPaddingMultiline(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithPadding(Columns(1), 3, 3),
		WithHeader([]string{"ID", "Lines"}),
	)
	if err := s.Render([]any{1, "a\nb\nc"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{2, "x"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_padding_multiline", buf.Bytes())
}

func TestGolden_TablePaddingZero(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithPadding(Columns(0), 0, 0),
		WithPadding(Columns(1), 0, 0),
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]any{
		{"x", 1},
		{"y", 2},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_padding_zero", buf.Bytes())
}

func TestGolden_StreamPaddingZero(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithPadding(Columns(0), 0, 0),
		WithPadding(Columns(1), 0, 0),
		WithHeader([]string{"A", "B"}),
	)
	if err := s.Render([]any{"x", 1}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"y", 2}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_padding_zero", buf.Bytes())
}

func TestGolden_TablePlaceholder(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithPlaceholder("N/A"),
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]any{
		{"x", ""},
		{nil, "y"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_placeholder", buf.Bytes())
}

func TestGolden_StreamPlaceholder(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithPlaceholder("N/A"),
		WithHeader([]string{"A", "B"}),
	)
	if err := s.Render([]any{"x", ""}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{nil, "y"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_placeholder", buf.Bytes())
}

func TestGolden_TablePlaceholderFixedWidth(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithPlaceholder("N/A"),
		WithHeader([]string{"A", "B"}),
		WithWidth(Columns(1), 2),
	)
	if err := tb.Render([][]any{{"x", ""}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_placeholder_fixed_width", buf.Bytes())
}

func TestGolden_StreamPlaceholderFixedWidth(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithPlaceholder("N/A"),
		WithHeader([]string{"A", "B"}),
		WithWidth(Columns(1), 2),
	)
	if err := s.Render([]any{"x", ""}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_placeholder_fixed_width", buf.Bytes())
}

func TestGolden_TablePlaceholderWideBytes(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithPlaceholder("\u2014"),
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]any{
		{"x", nil},
		{"y", "z"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_placeholder_wide_bytes", buf.Bytes())
}

func TestGolden_StreamPlaceholderWideBytes(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithPlaceholder("\u2014"),
		WithHeader([]string{"A", "B"}),
	)
	for _, row := range [][]any{{"x", nil}, {"y", "z"}} {
		if err := s.Render(row); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_placeholder_wide_bytes", buf.Bytes())
}

func TestGolden_TableRaggedRows(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := tb.Render([][]any{
		{"x"},
		{"p", "q", "r"},
		{},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_ragged_rows", buf.Bytes())
}

func TestGolden_StreamRaggedRows(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := s.Render([]any{"x"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"p", "q", "r"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_ragged_rows", buf.Bytes())
}

func TestGolden_TableRowsOnly(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
	)
	if err := tb.Render([][]any{
		{"a", 1},
		{"b", 2},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_rows_only", buf.Bytes())
}

func TestGolden_StreamRowsOnly(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
	)
	if err := s.Render([]any{"a", 1}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"b", 2}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_rows_only", buf.Bytes())
}

func TestGolden_TableRowspan(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
		WithHeader([]string{"Group", "Item"}),
	)
	if err := tb.Render([][]any{
		{"A", "x"},
		{"A", "y"},
		{"B", "z"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_rowspan", buf.Bytes())
}

func TestGolden_StreamRowspan(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
		WithHeader([]string{"Group", "Item"}),
	)
	if err := s.Render([]any{"A", "x"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"A", "y"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"B", "z"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_rowspan", buf.Bytes())
}

func TestGolden_TableRowspanAlign(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), AlignRight),
		WithHeader([]string{"Group", "Score"}),
	)
	if err := tb.Render([][]any{
		{"A", 100},
		{"A", 200},
		{"B", 300},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_rowspan_align", buf.Bytes())
}

func TestGolden_StreamRowspanAlign(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), AlignRight),
		WithHeader([]string{"Group", "Score"}),
	)
	if err := s.Render([]any{"A", 100}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"A", 200}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"B", 300}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_rowspan_align", buf.Bytes())
}

func TestGolden_TableRowspanAttr(t *testing.T) {
	restore := isTerminal
	isTerminal = func(io.Writer) bool { return true }
	t.Cleanup(func() { isTerminal = restore })
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
		WithAttr(ScopeBody, Columns(1), ColorFgRed),
	)
	if err := tb.Render([][]any{{"s", "s"}, {"s", "s"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_rowspan_attr", buf.Bytes())
}

func TestGolden_StreamRowspanAttr(t *testing.T) {
	restore := isTerminal
	isTerminal = func(io.Writer) bool { return true }
	t.Cleanup(func() { isTerminal = restore })
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
		WithAttr(ScopeBody, Columns(1), ColorFgRed),
	)
	for _, r := range [][]any{{"s", "s"}, {"s", "s"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_rowspan_attr", buf.Bytes())
}

func TestGolden_TableRowspanAutoFit(t *testing.T) {
	restore := terminalWidth
	terminalWidth = func(io.Writer) int { return 20 }
	t.Cleanup(func() { terminalWidth = restore })
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
		WithAutoFit(),
	)
	if err := tb.Render([][]any{{"g", "a fairly long value"}, {"g", "another long one"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_rowspan_autofit", buf.Bytes())
}

func TestGolden_StreamRowspanAutoFit(t *testing.T) {
	restore := terminalWidth
	terminalWidth = func(io.Writer) int { return 20 }
	t.Cleanup(func() { terminalWidth = restore })
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
		WithAutoFit(),
	)
	for _, r := range [][]any{{"g", "a fairly long value"}, {"g", "another long one"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_rowspan_autofit", buf.Bytes())
}

func TestGolden_TableRowspanCaption(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
		WithCaption("cap", CaptionBottom),
	)
	if err := tb.Render([][]any{{"s", "s"}, {"s", "s"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_rowspan_caption", buf.Bytes())
}

func TestGolden_StreamRowspanCaption(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
		WithCaption("cap", CaptionBottom),
	)
	for _, r := range [][]any{{"s", "s"}, {"s", "s"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_rowspan_caption", buf.Bytes())
}

func TestGolden_TableRowspanColspanEdge(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B", "C"}),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(1, 2)),
	)
	if err := tb.Render([][]any{
		{"g", "x", "x"},
		{"g", "y", "y"},
		{"h", "z", "w"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_rowspan_colspan_edge", buf.Bytes())
}

func TestGolden_StreamRowspanColspanEdge(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B", "C"}),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(1, 2)),
	)
	if err := s.Render([]any{"g", "x", "x"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"g", "y", "y"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"h", "z", "w"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_rowspan_colspan_edge", buf.Bytes())
}

func TestGolden_TableRowspanMissingKinds(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(1)),
		WithPlaceholder("X"),
	)
	if err := tb.Render([][]any{
		{1, "X"},
		{2},
		{3, nil},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_rowspan_missing_kinds", buf.Bytes())
}

func TestGolden_StreamRowspanMissingKinds(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(1)),
		WithPlaceholder("X"),
	)
	if err := s.Render([]any{1, "X"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{2}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{3, nil}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_rowspan_missing_kinds", buf.Bytes())
}

func TestGolden_TableRowspanPadding(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
		WithPadding(Columns(1), 2, 2),
		WithHeader([]string{"Group", "Value"}),
	)
	if err := tb.Render([][]any{
		{"A", "x"},
		{"A", "y"},
		{"B", "z"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_rowspan_padding", buf.Bytes())
}

func TestGolden_StreamRowspanPadding(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
		WithPadding(Columns(1), 2, 2),
		WithHeader([]string{"Group", "Value"}),
	)
	if err := s.Render([]any{"A", "x"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"A", "y"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"B", "z"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_rowspan_padding", buf.Bytes())
}

func TestGolden_TableRowspanPlaceholder(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B", "C"}),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
		WithPlaceholder("-"),
	)
	if err := tb.Render([][]any{
		{"g", "", "x"},
		{"g", "", "y"},
		{"", "", ""},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_rowspan_placeholder", buf.Bytes())
}

func TestGolden_StreamRowspanPlaceholder(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B", "C"}),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
		WithPlaceholder("-"),
	)
	if err := s.Render([]any{"g", "", "x"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"g", "", "y"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"", "", ""}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_rowspan_placeholder", buf.Bytes())
}

func TestGolden_TableRowspanScope(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"Env", "Tier"}, []string{"Env", "Tier"}),
		WithRowspan(ScopeBody, Columns(0)),
	)
	if err := tb.Render([][]any{
		{"prod", "web"},
		{"prod", "db"},
		{"dev", "web"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_rowspan_scope", buf.Bytes())
}

func TestGolden_StreamRowspanScope(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"Env", "Tier"}, []string{"Env", "Tier"}),
		WithRowspan(ScopeBody, Columns(0)),
	)
	for _, row := range [][]any{
		{"prod", "web"},
		{"prod", "db"},
		{"dev", "web"},
	} {
		if err := s.Render(row); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_rowspan_scope", buf.Bytes())
}

func TestGolden_TableRowspanTransformer(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(1)),
		WithTransformer(Columns(1), func(v any) (string, *Attr) {
			if s, ok := v.(string); ok && s == "raw" {
				return "T", nil
			}
			return "", nil
		}),
	)
	if err := tb.Render([][]any{{"x", "raw"}, {"p", "raw"}, {"q", "z"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_rowspan_transformer", buf.Bytes())
}

func TestGolden_StreamRowspanTransformer(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(1)),
		WithTransformer(Columns(1), func(v any) (string, *Attr) {
			if s, ok := v.(string); ok && s == "raw" {
				return "T", nil
			}
			return "", nil
		}),
	)
	for _, r := range [][]any{{"x", "raw"}, {"p", "raw"}, {"q", "z"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_rowspan_transformer", buf.Bytes())
}

func TestGolden_TableRowspanTruncate(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
		WithWidth(Columns(1), 5),
		WithTruncate(Columns(1)),
	)
	if err := tb.Render([][]any{{"g", "a long value"}, {"g", "short"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_rowspan_truncate", buf.Bytes())
}

func TestGolden_StreamRowspanTruncate(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
		WithWidth(Columns(1), 5),
		WithTruncate(Columns(1)),
	)
	for _, r := range [][]any{{"g", "a long value"}, {"g", "short"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_rowspan_truncate", buf.Bytes())
}

func TestGolden_TableRowspanWidth(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
		WithWidth(Columns(1), 6),
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]any{{"g", "a long value"}, {"g", "short"}, {"h", "x"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_rowspan_width", buf.Bytes())
}

func TestGolden_StreamRowspanWidth(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
		WithWidth(Columns(1), 6),
		WithHeader([]string{"A", "B"}),
	)
	for _, r := range [][]any{{"g", "a long value"}, {"g", "short"}, {"h", "x"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_rowspan_width", buf.Bytes())
}

func TestGolden_TableSingleCell(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"only"}),
	)
	if err := tb.Render([][]any{{"v"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_single_cell", buf.Bytes())
}

func TestGolden_StreamSingleCell(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"only"}),
	)
	if err := s.Render([]any{"v"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_single_cell", buf.Bytes())
}

func TestGolden_TableSpanLimit(t *testing.T) {
	const n = param.SpanLimit + 2
	header := make([]string, n)
	cols := make([]int, n)
	for i := range n {
		header[i] = "c"
		cols[i] = i
	}
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader(header),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(cols...)),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(cols...)),
	)
	rows := make([][]any, 2)
	for r := range rows {
		row := make([]any, n)
		for i := range n {
			row[i] = "v"
		}
		rows[r] = row
	}
	if err := tb.Render(rows); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_span_limit", buf.Bytes())
}

func TestGolden_StreamSpanLimit(t *testing.T) {
	const n = param.SpanLimit + 2
	header := make([]string, n)
	cols := make([]int, n)
	for i := range n {
		header[i] = "c"
		cols[i] = i
	}
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader(header),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(cols...)),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(cols...)),
	)
	for range 2 {
		row := make([]any, n)
		for i := range n {
			row[i] = "v"
		}
		if err := s.Render(row); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_span_limit", buf.Bytes())
}

func TestGolden_TableStackedHeader(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader(
			[]string{"Resources", "Resources", "Resources", "Resources", "ID"},
			[]string{"Network", "Network", "Security", "Security", "ID"},
			[]string{"VPC", "Subnet", "SG", "NACL", "ID"},
		),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1, 2, 3)),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(4)),
	)
	if err := tb.Render([][]any{
		{"vpc-1", "sub-1", "sg-1", "nacl-1", "i-001"},
		{"vpc-2", "sub-2", "sg-2", "nacl-2", "i-002"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_stacked_header", buf.Bytes())
}

func TestGolden_StreamStackedHeaderStack(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader(
			[]string{"Resources", "Resources", "Resources", "Resources", "ID"},
			[]string{"Network", "Network", "Security", "Security", "ID"},
			[]string{"VPC", "Subnet", "SG", "NACL", "ID"},
		),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1, 2, 3)),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(4)),
	)
	for _, row := range [][]any{
		{"vpc-1", "sub-1", "sg-1", "nacl-1", "i-001"},
		{"vpc-2", "sub-2", "sg-2", "nacl-2", "i-002"},
	} {
		if err := s.Render(row); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_stacked_header", buf.Bytes())
}

func TestGolden_TableStackedHeaderNested(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader(
			[]string{"Very long top label", "Very long top label", "Very long top label", "Very long top label", "ID"},
			[]string{"Mid label", "Mid label", "Other", "Other", "ID"},
			[]string{"a", "b", "c", "d", "ID"},
		),
		WithColspan(ScopeHeader, Columns(0, 1, 2, 3)),
	)
	if err := tb.Render([][]any{{1, 2, 3, 4, "x"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_stacked_header_nested", buf.Bytes())
}

func TestGolden_StreamStackedHeaderNested(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader(
			[]string{"Very long top label", "Very long top label", "Very long top label", "Very long top label", "ID"},
			[]string{"Mid label", "Mid label", "Other", "Other", "ID"},
			[]string{"a", "b", "c", "d", "ID"},
		),
		WithColspan(ScopeHeader, Columns(0, 1, 2, 3)),
	)
	if err := s.Render([]any{1, 2, 3, 4, "x"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_stacked_header_nested", buf.Bytes())
}

func TestGolden_TableStackedHeaderWideLabel(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader(
			[]string{"Availability zone", "Availability zone", "ID"},
			[]string{"Zone", "Z", "ID"},
		),
		WithColspan(ScopeHeader, Columns(0, 1)),
	)
	if err := tb.Render([][]any{{"us-east-1a", "a", "1"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_stacked_header_wide_label", buf.Bytes())
}

func TestGolden_StreamStackedHeaderWideLabel(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader(
			[]string{"Availability zone", "Availability zone", "ID"},
			[]string{"Zone", "Z", "ID"},
		),
		WithColspan(ScopeHeader, Columns(0, 1)),
	)
	if err := s.Render([]any{"us-east-1a", "a", "1"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_stacked_header_wide_label", buf.Bytes())
}

func TestGolden_TableStackedHeaderWideLabelFixed(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader(
			[]string{"Availability zone", "Availability zone", "ID"},
			[]string{"Zone", "Z", "ID"},
		),
		WithColspan(ScopeHeader, Columns(0, 1)),
		WithWidth(Columns(0, 1), 4),
	)
	if err := tb.Render([][]any{{"us-east-1a", "a", "1"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_stacked_header_wide_label_fixed", buf.Bytes())
}

func TestGolden_StreamStackedHeaderWideLabelFixed(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader(
			[]string{"Availability zone", "Availability zone", "ID"},
			[]string{"Zone", "Z", "ID"},
		),
		WithColspan(ScopeHeader, Columns(0, 1)),
		WithWidth(Columns(0, 1), 4),
	)
	if err := s.Render([]any{"us-east-1a", "a", "1"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_stacked_header_wide_label_fixed", buf.Bytes())
}

func TestGolden_TableStyleASCII(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleASCII),
		WithHeader([]string{"Name", "Value"}),
		WithFooter(func() [][]string {
			return [][]string{{"total", "3"}}
		}),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
	)
	if err := tb.Render([][]any{
		{"foo", 1},
		{"foo", 2},
		{"bar", 3},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_style_ascii", buf.Bytes())
}

func TestGolden_StreamStyleASCII(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleASCII),
		WithHeader([]string{"Name", "Value"}),
		WithFooter(func() [][]string {
			return [][]string{{"total", "3"}}
		}),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
	)
	for _, row := range [][]any{{"foo", 1}, {"foo", 2}, {"bar", 3}} {
		if err := s.Render(row); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_style_ascii", buf.Bytes())
}

func TestGolden_TableStyleBorderless(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(Style{}),
		WithHeader([]string{"Name", "Value"}),
		WithFooter(func() [][]string {
			return [][]string{{"total", "3"}}
		}),
	)
	if err := tb.Render([][]any{
		{"foo", 1},
		{"bar", 2},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_style_borderless", buf.Bytes())
}

func TestGolden_StreamStyleBorderless(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(Style{}),
		WithHeader([]string{"Name", "Value"}),
		WithFooter(func() [][]string {
			return [][]string{{"total", "3"}}
		}),
	)
	for _, row := range [][]any{{"foo", 1}, {"bar", 2}} {
		if err := s.Render(row); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_style_borderless", buf.Bytes())
}

func TestGolden_TableStyleColored(t *testing.T) {
	restore := isTerminal
	isTerminal = func(io.Writer) bool { return true }
	t.Cleanup(func() { isTerminal = restore })
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleColoredLight),
		WithHeader([]string{"Name", "Value"}),
		WithFooter(func() [][]string {
			return [][]string{{"total", "3"}}
		}),
		WithCaption("caption", CaptionDefault),
	)
	if err := tb.Render([][]any{
		{"foo", 1},
		{"bar", 2},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_style_colored", buf.Bytes())
}

func TestGolden_StreamStyleColored(t *testing.T) {
	restore := isTerminal
	isTerminal = func(io.Writer) bool { return true }
	t.Cleanup(func() { isTerminal = restore })
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleColoredLight),
		WithHeader([]string{"Name", "Value"}),
		WithFooter(func() [][]string {
			return [][]string{{"total", "3"}}
		}),
		WithCaption("caption", CaptionDefault),
	)
	for _, row := range [][]any{{"foo", 1}, {"bar", 2}} {
		if err := s.Render(row); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_style_colored", buf.Bytes())
}

func TestGolden_TableStyleDouble(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleDouble),
		WithHeader([]string{"Name", "Value"}),
	)
	if err := tb.Render([][]any{
		{"foo", 1},
		{"bar", 2},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_style_double", buf.Bytes())
}

func TestGolden_StreamStyleDouble(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleDouble),
		WithHeader([]string{"Name", "Value"}),
	)
	if err := s.Render([]any{"foo", 1}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"bar", 2}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_style_double", buf.Bytes())
}

func TestGolden_TableStyleHeavy(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleHeavy),
		WithHeader([]string{"Name", "Value"}),
	)
	if err := tb.Render([][]any{
		{"foo", 1},
		{"bar", 2},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_style_heavy", buf.Bytes())
}

func TestGolden_StreamStyleHeavy(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleHeavy),
		WithHeader([]string{"Name", "Value"}),
	)
	if err := s.Render([]any{"foo", 1}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"bar", 2}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_style_heavy", buf.Bytes())
}

func TestGolden_TableStyleLight(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"Name", "Value"}),
	)
	if err := tb.Render([][]any{
		{"foo", 1},
		{"bar", 2},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_style_light", buf.Bytes())
}

func TestGolden_StreamStyleLight(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"Name", "Value"}),
	)
	if err := s.Render([]any{"foo", 1}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"bar", 2}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_style_light", buf.Bytes())
}

func TestGolden_TableStyleNoBodyHorizon(t *testing.T) {
	style := StyleLight
	style.Border.Body = nil
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(style),
		WithHeader(
			[]string{"Group", "Group", "ID"},
			[]string{"A", "B", "ID"},
		),
		WithColspan(ScopeHeader, Columns(0, 1)),
	)
	if err := tb.Render([][]any{
		{"foo", 1, "x"},
		{"bar", 2, "y"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_style_no_body_horizon", buf.Bytes())
}

func TestGolden_StreamStyleNoBodyHorizon(t *testing.T) {
	style := StyleLight
	style.Border.Body = nil
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(style),
		WithHeader(
			[]string{"Group", "Group", "ID"},
			[]string{"A", "B", "ID"},
		),
		WithColspan(ScopeHeader, Columns(0, 1)),
	)
	for _, row := range [][]any{{"foo", 1, "x"}, {"bar", 2, "y"}} {
		if err := s.Render(row); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_style_no_body_horizon", buf.Bytes())
}

func TestGolden_TableStyleNoBodyHorizonColspan(t *testing.T) {
	style := StyleLight
	style.Border.Body = nil
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(style),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1, 2)),
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := tb.Render([][]any{
		{"x", "x", "y"},
		{"p", "p", "q"},
		{"m", "n", "n"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_style_no_body_horizon_colspan", buf.Bytes())
}

func TestGolden_StreamStyleNoBodyHorizonColspan(t *testing.T) {
	style := StyleLight
	style.Border.Body = nil
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(style),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1, 2)),
		WithHeader([]string{"A", "B", "C"}),
	)
	for _, r := range [][]any{
		{"x", "x", "y"},
		{"p", "p", "q"},
		{"m", "n", "n"},
	} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_style_no_body_horizon_colspan", buf.Bytes())
}

func TestGolden_TableStyleRounded(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleRounded),
		WithHeader([]string{"Name", "Value"}),
	)
	if err := tb.Render([][]any{
		{"foo", 1},
		{"bar", 2},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_style_rounded", buf.Bytes())
}

func TestGolden_StreamStyleRounded(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleRounded),
		WithHeader([]string{"Name", "Value"}),
	)
	if err := s.Render([]any{"foo", 1}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"bar", 2}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_style_rounded", buf.Bytes())
}

func TestGolden_TableTransformer(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithTransformer(Columns(1), func(v any) (string, *Attr) {
			n, ok := v.(int)
			if !ok {
				return "", nil
			}
			if n >= 100 {
				return "high", ColorFgRed
			}
			return "", nil
		}),
		WithHeader([]string{"Name", "Score"}),
	)
	if err := tb.Render([][]any{
		{"alice", 100},
		{"bob", 99},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_transformer", buf.Bytes())
}

func TestGolden_StreamTransformer(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithTransformer(Columns(1), func(v any) (string, *Attr) {
			n, ok := v.(int)
			if !ok {
				return "", nil
			}
			if n >= 100 {
				return "high", ColorFgRed
			}
			return "", nil
		}),
		WithHeader([]string{"Name", "Score"}),
	)
	if err := s.Render([]any{"alice", 100}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"bob", 99}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_transformer", buf.Bytes())
}

func TestGolden_TableTransformerColumnOverride(t *testing.T) {
	restore := isTerminal
	isTerminal = func(io.Writer) bool { return true }
	t.Cleanup(func() { isTerminal = restore })
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"Level", "Message"}),
		WithAttr(ScopeBody, Columns(1), ColorFgBlue),
		WithTransformer(Columns(1), func(v any) (string, *Attr) {
			if s, ok := v.(string); ok && s == "warn" {
				return "", ColorFgYellow
			}
			return "", nil
		}),
	)
	if err := tb.Render([][]any{
		{"1", "ok"},
		{"2", "warn"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_transformer_column_override", buf.Bytes())
}

func TestGolden_StreamTransformerColumnOverride(t *testing.T) {
	restore := isTerminal
	isTerminal = func(io.Writer) bool { return true }
	t.Cleanup(func() { isTerminal = restore })
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"Level", "Message"}),
		WithAttr(ScopeBody, Columns(1), ColorFgBlue),
		WithTransformer(Columns(1), func(v any) (string, *Attr) {
			if s, ok := v.(string); ok && s == "warn" {
				return "", ColorFgYellow
			}
			return "", nil
		}),
	)
	for _, r := range [][]any{{"1", "ok"}, {"2", "warn"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_transformer_column_override", buf.Bytes())
}

func TestGolden_TableTruncate(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithWidth(Columns(1), 15),
		WithTruncate(Columns(1)),
		WithHeader([]string{"Name", "Description"}),
	)
	if err := tb.Render([][]any{
		{"alice", "short"},
		{"bob", "this is a longer description"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_truncate", buf.Bytes())
}

func TestGolden_StreamTruncate(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithWidth(Columns(1), 15),
		WithTruncate(Columns(1)),
		WithHeader([]string{"Name", "Description"}),
	)
	if err := s.Render([]any{"alice", "short"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"bob", "this is a longer description"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_truncate", buf.Bytes())
}

func TestGolden_TableTruncateCJK(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithWidth(Columns(1), 8),
		WithTruncate(Columns(1)),
		WithHeader([]string{"Label", "日本語"}),
	)
	if err := tb.Render([][]any{
		{"x", "あいうえおかきくけこ"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_truncate_cjk", buf.Bytes())
}

func TestGolden_StreamTruncateCJK(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithWidth(Columns(1), 8),
		WithTruncate(Columns(1)),
		WithHeader([]string{"Label", "日本語"}),
	)
	if err := s.Render([]any{"x", "あいうえおかきくけこ"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_truncate_cjk", buf.Bytes())
}

func TestGolden_TableTruncateFittingLines(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithWidth(Columns(1), 5),
		WithTruncate(Columns(1)),
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]any{{"x", "ab\ncd"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_truncate_fitting_lines", buf.Bytes())
}

func TestGolden_StreamTruncateFittingLines(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithWidth(Columns(1), 5),
		WithTruncate(Columns(1)),
		WithHeader([]string{"A", "B"}),
	)
	if err := s.Render([]any{"x", "ab\ncd"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_truncate_fitting_lines", buf.Bytes())
}

func TestGolden_TableTruncateMultiline(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithWidth(Columns(1), 10),
		WithTruncate(Columns(1)),
		WithHeader([]string{"ID", "Data"}),
	)
	if err := tb.Render([][]any{
		{1, "line1\nthis-is-a-long-line2"},
		{2, "short"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_truncate_multiline", buf.Bytes())
}

func TestGolden_StreamTruncateMultiline(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithWidth(Columns(1), 10),
		WithTruncate(Columns(1)),
		WithHeader([]string{"ID", "Data"}),
	)
	if err := s.Render([]any{1, "line1\nthis-is-a-long-line2"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{2, "short"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_truncate_multiline", buf.Bytes())
}

func TestGolden_TableTruncateToEllipsis(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithWidth(Columns(1), 1),
		WithTruncate(Columns(1)),
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]any{
		{"x", "long value"},
		{"y", "z"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_truncate_to_ellipsis", buf.Bytes())
}

func TestGolden_StreamTruncateToEllipsis(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithWidth(Columns(1), 1),
		WithTruncate(Columns(1)),
		WithHeader([]string{"A", "B"}),
	)
	for _, row := range [][]any{{"x", "long value"}, {"y", "z"}} {
		if err := s.Render(row); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_truncate_to_ellipsis", buf.Bytes())
}

func TestGolden_TableTypeEmptyNil(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"a", "b", "c"}),
	)
	if err := tb.Render([][]any{
		{"", nil, "x"},
		{nil, "", ""},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_type_empty_nil", buf.Bytes())
}

func TestGolden_StreamTypeEmptyNil(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"a", "b", "c"}),
	)
	if err := s.Render([]any{"", nil, "x"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{nil, "", ""}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_type_empty_nil", buf.Bytes())
}

func TestGolden_TableTypeFloat(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"f32", "f64"}),
	)
	if err := tb.Render([][]any{{float32(3.14), float64(2.71828)}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_type_float", buf.Bytes())
}

func TestGolden_StreamTypeFloat(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"f32", "f64"}),
	)
	if err := s.Render([]any{float32(3.14), float64(2.71828)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_type_float", buf.Bytes())
}

func TestGolden_TableTypeInteger(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{
			"i", "i8", "i16", "i32", "i64",
			"u", "u8", "u16", "u32", "u64",
		}),
	)
	if err := tb.Render([][]any{
		{int(-1), int8(-2), int16(-3), int32(-4), int64(-5),
			uint(1), uint8(2), uint16(3), uint32(4), uint64(5)},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_type_integer", buf.Bytes())
}

func TestGolden_StreamTypeInteger(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{
			"i", "i8", "i16", "i32", "i64",
			"u", "u8", "u16", "u32", "u64",
		}),
	)
	if err := s.Render([]any{
		int(-1), int8(-2), int16(-3), int32(-4), int64(-5),
		uint(1), uint8(2), uint16(3), uint32(4), uint64(5),
	}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_type_integer", buf.Bytes())
}

func TestGolden_TableTypePointer(t *testing.T) {
	var buf bytes.Buffer
	s := testutil.Stringer{Value: "y"}
	str := "alive"
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"a", "b"}),
	)
	if err := tb.Render([][]any{{&str, &s}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_type_pointer", buf.Bytes())
}

func TestGolden_StreamTypePointer(t *testing.T) {
	var buf bytes.Buffer
	st := testutil.Stringer{Value: "y"}
	str := "alive"
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"a", "b"}),
	)
	if err := s.Render([]any{&str, &st}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_type_pointer", buf.Bytes())
}

func TestGolden_TableTypeStringerError(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"Stringer", "Error"}),
	)
	if err := tb.Render([][]any{{testutil.Stringer{Value: "x"}, testutil.Error{Value: "boom"}}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_type_stringer_error", buf.Bytes())
}

func TestGolden_StreamTypeStringerError(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"Stringer", "Error"}),
	)
	if err := s.Render([]any{testutil.Stringer{Value: "x"}, testutil.Error{Value: "boom"}}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_type_stringer_error", buf.Bytes())
}

func TestGolden_TableTypeTypedNil(t *testing.T) {
	var buf bytes.Buffer
	var nilStringer *testutil.PtrStringer
	var nilError *testutil.PtrError
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithPlaceholder("<nil>"),
		WithHeader([]string{"Stringer", "Error"}),
	)
	if err := tb.Render([][]any{{nilStringer, nilError}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_type_typed_nil", buf.Bytes())
}

func TestGolden_StreamTypeTypedNil(t *testing.T) {
	var buf bytes.Buffer
	var nilStringer *testutil.PtrStringer
	var nilError *testutil.PtrError
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithPlaceholder("<nil>"),
		WithHeader([]string{"Stringer", "Error"}),
	)
	if err := s.Render([]any{nilStringer, nilError}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_type_typed_nil", buf.Bytes())
}

func TestGolden_TableValueEqualsPlaceholder(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithPlaceholder("N/A"),
	)
	if err := tb.Render([][]any{
		{"N/A", ""},
		{"x", "N/A"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_value_equals_placeholder", buf.Bytes())
}

func TestGolden_StreamValueEqualsPlaceholder(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
		WithPlaceholder("N/A"),
	)
	if err := s.Render([]any{"N/A", ""}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"x", "N/A"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_value_equals_placeholder", buf.Bytes())
}

func TestGolden_TableWidthAlign(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithWidth(Columns(0), 5),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), AlignRight),
	)
	if err := tb.Render([][]any{{"a long value", "yy"}, {"y", "z"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_width_align", buf.Bytes())
}

func TestGolden_StreamWidthAlign(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithWidth(Columns(0), 5),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), AlignRight),
	)
	for _, r := range [][]any{{"a long value", "yy"}, {"y", "z"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_width_align", buf.Bytes())
}

func TestGolden_TableWidthAttr(t *testing.T) {
	restore := isTerminal
	isTerminal = func(io.Writer) bool { return true }
	t.Cleanup(func() { isTerminal = restore })
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithWidth(Columns(0), 5),
		WithAttr(ScopeBody, Columns(1), ColorFgRed),
	)
	if err := tb.Render([][]any{{"a long value", "x"}, {"y", "z"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_width_attr", buf.Bytes())
}

func TestGolden_StreamWidthAttr(t *testing.T) {
	restore := isTerminal
	isTerminal = func(io.Writer) bool { return true }
	t.Cleanup(func() { isTerminal = restore })
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithWidth(Columns(0), 5),
		WithAttr(ScopeBody, Columns(1), ColorFgRed),
	)
	for _, r := range [][]any{{"a long value", "x"}, {"y", "z"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_width_attr", buf.Bytes())
}

func TestGolden_TableWidthCaption(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithWidth(Columns(0), 5),
		WithCaption("cap", CaptionBottom),
	)
	if err := tb.Render([][]any{{"a long value", "x"}, {"y", "z"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_width_caption", buf.Bytes())
}

func TestGolden_StreamWidthCaption(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithWidth(Columns(0), 5),
		WithCaption("cap", CaptionBottom),
	)
	for _, r := range [][]any{{"a long value", "x"}, {"y", "z"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_width_caption", buf.Bytes())
}

func TestGolden_TableWidthFooter(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithWidth(Columns(0), 5),
		WithFooter(func() [][]string {
			return [][]string{{"t", "u"}}
		}),
	)
	if err := tb.Render([][]any{{"a long value", "x"}, {"y", "z"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_width_footer", buf.Bytes())
}

func TestGolden_StreamWidthFooter(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithWidth(Columns(0), 5),
		WithFooter(func() [][]string {
			return [][]string{{"t", "u"}}
		}),
	)
	for _, r := range [][]any{{"a long value", "x"}, {"y", "z"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_width_footer", buf.Bytes())
}

func TestGolden_TableWidthPadding(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithWidth(Columns(0), 5),
		WithPadding(Columns(1), 3, 0),
	)
	if err := tb.Render([][]any{{"a long value", "x"}, {"y", "z"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_width_padding", buf.Bytes())
}

func TestGolden_StreamWidthPadding(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithWidth(Columns(0), 5),
		WithPadding(Columns(1), 3, 0),
	)
	for _, r := range [][]any{{"a long value", "x"}, {"y", "z"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_width_padding", buf.Bytes())
}

func TestGolden_TableWidthTransformer(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithWidth(Columns(0), 5),
		WithTransformer(Columns(1), func(v any) (string, *Attr) {
			if s, ok := v.(string); ok && s == "raw" {
				return "T", nil
			}
			return "", nil
		}),
	)
	if err := tb.Render([][]any{{"a long value", "raw"}, {"y", "q"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_width_transformer", buf.Bytes())
}

func TestGolden_StreamWidthTransformer(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithWidth(Columns(0), 5),
		WithTransformer(Columns(1), func(v any) (string, *Attr) {
			if s, ok := v.(string); ok && s == "raw" {
				return "T", nil
			}
			return "", nil
		}),
	)
	for _, r := range [][]any{{"a long value", "raw"}, {"y", "q"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_width_transformer", buf.Bytes())
}

func TestGolden_TableZeroWidth(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := tb.Render([][]any{
		{"a\u200Bb", "\uFEFF", "x\u200Dy"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_zero_width", buf.Bytes())
}

func TestGolden_StreamZeroWidth(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := s.Render([]any{"a\u200Bb", "\uFEFF", "x\u200Dy"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_zero_width", buf.Bytes())
}

func TestGolden_StreamAlignRow(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(0), AlignRight),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), AlignCenter),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(2), AlignLeft),
		WithHeader([]string{"R", "C", "L"}),
	)
	if err := s.Render([]any{"a", "b", 1}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"longer", "y", 22}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"x", "longest", 333}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_align_row", buf.Bytes())
}

func TestGolden_StreamAutoFit(t *testing.T) {
	restore := terminalWidth
	terminalWidth = func(io.Writer) int { return 24 }
	t.Cleanup(func() { terminalWidth = restore })
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithAutoFit(),
		WithHeader([]string{"A", "B"}),
	)
	if err := s.Render([]any{"xxxxx", strings.Repeat("y", 40)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"z", strings.Repeat("w", 20)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_autofit", buf.Bytes())
}

func TestGolden_StreamAutoFitNoRows(t *testing.T) {
	restore := terminalWidth
	terminalWidth = func(io.Writer) int { return 20 }
	t.Cleanup(func() { terminalWidth = restore })
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithAutoFit(),
		WithHeader([]string{"A very long label", "B very long label"}),
		WithFooter(func() [][]string {
			return [][]string{{"total", "0"}}
		}),
	)
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_autofit_no_rows", buf.Bytes())
}

func TestGolden_StreamCJK(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"日本語", "ASCII"}),
	)
	if err := s.Render([]any{"あいう", "abc"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"日本", "longer-text"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"テスト", "x"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_cjk", buf.Bytes())
}

func TestGolden_StreamColspanAutoFit(t *testing.T) {
	restore := terminalWidth
	terminalWidth = func(io.Writer) int { return 20 }
	t.Cleanup(func() { terminalWidth = restore })
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1)),
		WithAutoFit(),
	)
	for _, r := range [][]any{{"a fairly long value", "a fairly long value"}, {"another long one", "z"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_colspan_autofit", buf.Bytes())
}

func TestGolden_StreamColspanIndex(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1)),
		WithIndex(),
	)
	for _, r := range [][]any{{"s", "s"}, {"s", "s"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_colspan_index", buf.Bytes())
}

func TestGolden_StreamCompactIndex(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithCompact(),
		WithIndex(),
	)
	for _, r := range [][]any{{"s", "s"}, {"s", "s"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_compact_index", buf.Bytes())
}

func TestGolden_StreamEmptyVsNil(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"Kind", "Value"}),
	)
	if err := s.Render([]any{"nil", nil}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"empty", ""}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"space", " "}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"text", "hello"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_empty_vs_nil", buf.Bytes())
}

func TestGolden_StreamFrozenOverflow(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"a"}),
	)
	for _, v := range []string{"x", "\u3042", "y"} {
		if err := s.Render([]any{v}); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_frozen_overflow", buf.Bytes())
}

func TestGolden_StreamIndex(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithIndex(),
		WithHeader([]string{"Name", "Score"}),
	)
	if err := s.Render([]any{"alice", 100}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"bob", 99}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"carol", 98}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_index", buf.Bytes())
}

func TestGolden_StreamIndexAlign(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithIndex(),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), AlignRight),
	)
	for _, r := range [][]any{{"x", "yy"}, {"p", "z"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_index_align", buf.Bytes())
}

func TestGolden_StreamIndexAttr(t *testing.T) {
	restore := isTerminal
	isTerminal = func(io.Writer) bool { return true }
	t.Cleanup(func() { isTerminal = restore })
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithIndex(),
		WithAttr(ScopeBody, Columns(1), ColorFgRed),
	)
	for _, r := range [][]any{{"x", "y"}, {"p", "q"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_index_attr", buf.Bytes())
}

func TestGolden_StreamIndexCaption(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithIndex(),
		WithCaption("cap", CaptionBottom),
	)
	for _, r := range [][]any{{"x", "y"}, {"p", "q"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_index_caption", buf.Bytes())
}

func TestGolden_StreamIndexFooter(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithIndex(),
		WithHeader([]string{"Name", "Score"}),
		WithFooter(func() [][]string {
			return [][]string{{"total", "297"}}
		}),
	)
	for _, row := range [][]any{{"alice", 100}, {"bob", 99}, {"carol", 98}} {
		if err := s.Render(row); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_index_footer", buf.Bytes())
}

func TestGolden_StreamIndexPadding(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithIndex(),
		WithPadding(Columns(1), 3, 0),
	)
	for _, r := range [][]any{{"x", "y"}, {"p", "q"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_index_padding", buf.Bytes())
}

func TestGolden_StreamIndexTransformer(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithIndex(),
		WithTransformer(Columns(1), func(v any) (string, *Attr) {
			if s, ok := v.(string); ok && s == "raw" {
				return "T", nil
			}
			return "", nil
		}),
	)
	for _, r := range [][]any{{"x", "raw"}, {"p", "q"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_index_transformer", buf.Bytes())
}

func TestGolden_StreamIndexTruncate(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithIndex(),
		WithWidth(Columns(0), 5),
		WithTruncate(Columns(0)),
	)
	for _, r := range [][]any{{"a long value", "x"}, {"y", "z"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_index_truncate", buf.Bytes())
}

func TestGolden_StreamLongValue(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
	)
	if err := s.Render([]any{"s", 1}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{strings.Repeat("longvalue", 40), 2}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"t", 3}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_long_value", buf.Bytes())
}

func TestGolden_StreamRowspanIndex(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
		WithIndex(),
	)
	for _, r := range [][]any{{"s", "s"}, {"s", "s"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_rowspan_index", buf.Bytes())
}

func TestGolden_StreamRowspanNumeric(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
		WithHeader([]string{"ID", "Name"}),
	)
	if err := s.Render([]any{100, "a"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{100, "b"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{200, "c"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{300, "d"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_rowspan_numeric", buf.Bytes())
}

func TestGolden_StreamSingleColumn(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"x"}),
	)
	if err := s.Render([]any{"a"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"bb"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"ccc"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_single_column", buf.Bytes())
}

func TestGolden_StreamSingleRow(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"Name", "Value"}),
	)
	if err := s.Render([]any{"foo", 1}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_single_row", buf.Bytes())
}

func TestGolden_StreamSpecialChars(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"Char", "Value"}),
	)
	if err := s.Render([]any{"less-than", "<"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"pipe", "|"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"backslash", "\\"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"ampersand", "&"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"asterisk", "*"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"underscore", "_"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_special_chars", buf.Bytes())
}

func TestGolden_StreamTypeMixed(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"Label", "Value"}),
	)
	if err := s.Render([]any{"text", "hello"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"number", 42}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"big-num", 100000}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"empty", nil}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"neg", -7}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_type_mixed", buf.Bytes())
}

func TestGolden_StreamTypeSlice(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"Type", "Value"}),
	)
	if err := s.Render([]any{"[]int", []int{1, 2, 3}}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"[]float64", []float64{1.1, 2.2}}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"[]bool", []bool{true, false, true}}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"[]byte", []byte("hello")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"[]byte (raw)", []byte{0x00, 0xff}}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"[]string", []string{"a", "b"}}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"[]any", []any{"x", 1, nil}}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_type_slice", buf.Bytes())
}

func TestGolden_StreamWideNumber(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"N", "V"}),
	)
	if err := s.Render([]any{1, 1}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{2, 1000000}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{3, -99}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_wide_number", buf.Bytes())
}

func TestGolden_StreamWidthIndex(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithWidth(Columns(0), 5),
		WithIndex(),
	)
	for _, r := range [][]any{{"a long value", "x"}, {"y", "z"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_width_index", buf.Bytes())
}

func TestGolden_TableAlignRow(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(0), AlignRight),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), AlignCenter),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(2), AlignLeft),
		WithHeader([]string{"R", "C", "L"}),
	)
	if err := tb.Render([][]any{
		{"a", "b", 1},
		{"longer", "y", 22},
		{"x", "longest", 333},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_align_row", buf.Bytes())
}

func TestGolden_TableAutoFit(t *testing.T) {
	restore := terminalWidth
	terminalWidth = func(io.Writer) int { return 24 }
	t.Cleanup(func() { terminalWidth = restore })
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithAutoFit(),
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]any{
		{"xxxxx", strings.Repeat("y", 40)},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_autofit", buf.Bytes())
}

func TestGolden_TableAutoFitIndex(t *testing.T) {
	restore := terminalWidth
	terminalWidth = func(io.Writer) int { return 30 }
	t.Cleanup(func() { terminalWidth = restore })
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithAutoFit(),
		WithIndex(),
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]any{
		{strings.Repeat("x", 20), strings.Repeat("y", 21)},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_autofit_index", buf.Bytes())
}

func TestGolden_TableCJK(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"日本語", "ASCII"}),
	)
	if err := tb.Render([][]any{
		{"あいう", "abc"},
		{"日本", "longer-text"},
		{"テスト", "x"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_cjk", buf.Bytes())
}

func TestGolden_TableColspanAutoFit(t *testing.T) {
	restore := terminalWidth
	terminalWidth = func(io.Writer) int { return 20 }
	t.Cleanup(func() { terminalWidth = restore })
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1)),
		WithAutoFit(),
	)
	if err := tb.Render([][]any{{"a fairly long value", "a fairly long value"}, {"another long one", "z"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_colspan_autofit", buf.Bytes())
}

func TestGolden_TableColspanIndex(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1)),
		WithIndex(),
	)
	if err := tb.Render([][]any{{"s", "s"}, {"s", "s"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_colspan_index", buf.Bytes())
}

func TestGolden_TableCompactIndex(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithCompact(),
		WithIndex(),
	)
	if err := tb.Render([][]any{{"s", "s"}, {"s", "s"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_compact_index", buf.Bytes())
}

func TestGolden_TableEmptyVsNil(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"Kind", "Value"}),
	)
	if err := tb.Render([][]any{
		{"nil", nil},
		{"empty", ""},
		{"space", " "},
		{"text", "hello"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_empty_vs_nil", buf.Bytes())
}

func TestGolden_TableIndex(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithIndex(),
		WithHeader([]string{"Name", "Score"}),
	)
	if err := tb.Render([][]any{
		{"alice", 100},
		{"bob", 99},
		{"carol", 98},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_index", buf.Bytes())
}

func TestGolden_TableIndexAlign(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithIndex(),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), AlignRight),
	)
	if err := tb.Render([][]any{{"x", "yy"}, {"p", "z"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_index_align", buf.Bytes())
}

func TestGolden_TableIndexAttr(t *testing.T) {
	restore := isTerminal
	isTerminal = func(io.Writer) bool { return true }
	t.Cleanup(func() { isTerminal = restore })
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithIndex(),
		WithAttr(ScopeBody, Columns(1), ColorFgRed),
	)
	if err := tb.Render([][]any{{"x", "y"}, {"p", "q"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_index_attr", buf.Bytes())
}

func TestGolden_TableIndexCaption(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithIndex(),
		WithCaption("cap", CaptionBottom),
	)
	if err := tb.Render([][]any{{"x", "y"}, {"p", "q"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_index_caption", buf.Bytes())
}

func TestGolden_TableIndexFooter(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithIndex(),
		WithHeader([]string{"Name", "Score"}),
		WithFooter(func() [][]string {
			return [][]string{{"total", "297"}}
		}),
	)
	if err := tb.Render([][]any{
		{"alice", 100},
		{"bob", 99},
		{"carol", 98},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_index_footer", buf.Bytes())
}

func TestGolden_TableIndexPadding(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithIndex(),
		WithPadding(Columns(1), 3, 0),
	)
	if err := tb.Render([][]any{{"x", "y"}, {"p", "q"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_index_padding", buf.Bytes())
}

func TestGolden_TableIndexTransformer(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithIndex(),
		WithTransformer(Columns(1), func(v any) (string, *Attr) {
			if s, ok := v.(string); ok && s == "raw" {
				return "T", nil
			}
			return "", nil
		}),
	)
	if err := tb.Render([][]any{{"x", "raw"}, {"p", "q"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_index_transformer", buf.Bytes())
}

func TestGolden_TableIndexTruncate(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithIndex(),
		WithWidth(Columns(0), 5),
		WithTruncate(Columns(0)),
	)
	if err := tb.Render([][]any{{"a long value", "x"}, {"y", "z"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_index_truncate", buf.Bytes())
}

func TestGolden_TableLongValue(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]any{
		{"s", 1},
		{strings.Repeat("longvalue", 40), 2},
		{"t", 3},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_long_value", buf.Bytes())
}

func TestGolden_TableNoHeaderRagged(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
	)
	if err := tb.Render([][]any{
		{"a", 1},
		{"b"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_no_header_ragged", buf.Bytes())
}

func TestGolden_TableRowspanIndex(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
		WithIndex(),
	)
	if err := tb.Render([][]any{{"s", "s"}, {"s", "s"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_rowspan_index", buf.Bytes())
}

func TestGolden_TableSingleColumn(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"x"}),
	)
	if err := tb.Render([][]any{
		{"a"},
		{"bb"},
		{"ccc"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_single_column", buf.Bytes())
}

func TestGolden_TableSpecialChars(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"Char", "Value"}),
	)
	if err := tb.Render([][]any{
		{"less-than", "<"},
		{"pipe", "|"},
		{"backslash", "\\"},
		{"ampersand", "&"},
		{"asterisk", "*"},
		{"underscore", "_"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_special_chars", buf.Bytes())
}

func TestGolden_TableTypeMixed(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"Label", "Value"}),
	)
	if err := tb.Render([][]any{
		{"text", "hello"},
		{"number", 42},
		{"big-num", 100000},
		{"empty", nil},
		{"neg", -7},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_type_mixed", buf.Bytes())
}

func TestGolden_TableTypeSlice(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"Type", "Value"}),
	)
	if err := tb.Render([][]any{
		{"[]int", []int{1, 2, 3}},
		{"[]float64", []float64{1.1, 2.2}},
		{"[]bool", []bool{true, false, true}},
		{"[]byte", []byte("hello")},
		{"[]byte (raw)", []byte{0x00, 0xff}},
		{"[]string", []string{"a", "b"}},
		{"[]any", []any{"x", 1, nil}},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_type_slice", buf.Bytes())
}

func TestGolden_TableWideNumber(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"N", "V"}),
	)
	if err := tb.Render([][]any{
		{1, 1},
		{2, 1000000},
		{3, -99},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_wide_number", buf.Bytes())
}

func TestGolden_TableWidthIndex(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithWidth(Columns(0), 5),
		WithIndex(),
	)
	if err := tb.Render([][]any{{"a long value", "x"}, {"y", "z"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_width_index", buf.Bytes())
}
