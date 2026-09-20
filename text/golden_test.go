package text

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/nekrassov01/table"
	"github.com/nekrassov01/table/internal/param"
	"github.com/nekrassov01/table/internal/testutil"
)

func TestGolden_TableAlignAttrPaddingTransformerTruncate(t *testing.T) {
	restore := isTerminal
	isTerminal = func(io.Writer) bool { return true }
	t.Cleanup(func() { isTerminal = restore })
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithWidth(Columns(1), 5),
		WithPadding(Columns(1), 2, 1),
		WithAlign(ScopeBody, Columns(1), AlignRight),
		WithAttr(ScopeBody, Columns(1), ColorFgRed),
		WithTransformer(Columns(1), func(v any) (string, *Attr) {
			if v == "raw" {
				return "transformed", nil
			}
			return "", nil
		}),
		WithTruncate(Columns(1)),
	)
	if err := tb.Render([][]table.Value{{table.String("x"), table.String("raw")}, {table.String("y"), table.String("ok")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_align_attr_padding_transformer_truncate", buf.Bytes())
}

func TestGolden_StreamAlignAttrPaddingTransformerTruncate(t *testing.T) {
	restore := isTerminal
	isTerminal = func(io.Writer) bool { return true }
	t.Cleanup(func() { isTerminal = restore })
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithWidth(Columns(1), 5),
		WithPadding(Columns(1), 2, 1),
		WithAlign(ScopeBody, Columns(1), AlignRight),
		WithAttr(ScopeBody, Columns(1), ColorFgRed),
		WithTransformer(Columns(1), func(v any) (string, *Attr) {
			if v == "raw" {
				return "transformed", nil
			}
			return "", nil
		}),
		WithTruncate(Columns(1)),
	)
	for _, row := range [][]table.Value{{table.String("x"), table.String("raw")}, {table.String("y"), table.String("ok")}} {
		if err := s.Render(row); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_align_attr_padding_transformer_truncate", buf.Bytes())
}

func TestGolden_TableAlignAttrPlaceholder(t *testing.T) {
	restore := isTerminal
	isTerminal = func(io.Writer) bool { return true }
	t.Cleanup(func() { isTerminal = restore })
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithPlaceholder("N/A"),
		WithAlign(ScopeBody, Columns(1), AlignRight),
		WithAttr(ScopeBody, Columns(1), ColorFgRed),
	)
	if err := tb.Render([][]table.Value{{table.String("x"), table.Any(nil)}, {table.String("y"), table.String("z")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_align_attr_placeholder", buf.Bytes())
}

func TestGolden_StreamAlignAttrPlaceholder(t *testing.T) {
	restore := isTerminal
	isTerminal = func(io.Writer) bool { return true }
	t.Cleanup(func() { isTerminal = restore })
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithPlaceholder("N/A"),
		WithAlign(ScopeBody, Columns(1), AlignRight),
		WithAttr(ScopeBody, Columns(1), ColorFgRed),
	)
	for _, row := range [][]table.Value{{table.String("x"), table.Any(nil)}, {table.String("y"), table.String("z")}} {
		if err := s.Render(row); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_align_attr_placeholder", buf.Bytes())
}

func TestGolden_TableAlignMultiline(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), AlignCenter),
		WithHeader([]string{"ID", "Lines"}),
	)
	if err := tb.Render([][]table.Value{
		{table.Int(1), table.String("short\nlonger-text\nx")},
		{table.Int(2), table.String("only")},
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
	if err := s.Render([]table.Value{table.Int(1), table.String("short\nlonger-text\nx")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.Int(2), table.String("only")}); err != nil {
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
	if err := tb.Render([][]table.Value{
		{table.String("widget"), table.Int(100)},
		{table.String("gadget"), table.Int(200)},
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
	if err := s.Render([]table.Value{table.String("widget"), table.Int(100)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("gadget"), table.Int(200)}); err != nil {
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
	if err := tb.Render([][]table.Value{
		{table.Any(nil), table.Any(nil), table.Any(nil)},
		{table.Any(nil), table.Any(nil), table.Any(nil)},
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
	if err := s.Render([]table.Value{table.Any(nil), table.Any(nil), table.Any(nil)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.Any(nil), table.Any(nil), table.Any(nil)}); err != nil {
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
	if err := tb.Render([][]table.Value{
		{table.String("foo"), table.Int(1)},
		{table.String("bar"), table.Int(2)},
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
	if err := s.Render([]table.Value{table.String("foo"), table.Int(1)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("bar"), table.Int(2)}); err != nil {
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
	if err := tb.Render([][]table.Value{{table.String("foo"), table.Int(1)}, {table.String("bar"), table.Int(2)}}); err != nil {
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
	for _, row := range [][]table.Value{{table.String("foo"), table.Int(1)}, {table.String("bar"), table.Int(2)}} {
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
	if err := tb.Render([][]table.Value{
		{table.String("foo"), table.Int(1)},
		{table.String("bar"), table.Any(nil)},
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
	for _, row := range [][]table.Value{{table.String("foo"), table.Int(1)}, {table.String("bar"), table.Any(nil)}} {
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
	if err := tb.Render([][]table.Value{
		{table.String("foo"), table.Int(1)},
		{table.String("bar"), table.Int(2)},
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
	for _, row := range [][]table.Value{{table.String("foo"), table.Int(1)}, {table.String("bar"), table.Int(2)}} {
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
	if err := tb.Render([][]table.Value{{table.String("a fairly long value"), table.String("yy")}, {table.String("another long one"), table.String("z")}}); err != nil {
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
	for _, r := range [][]table.Value{{table.String("a fairly long value"), table.String("yy")}, {table.String("another long one"), table.String("z")}} {
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
	if err := tb.Render([][]table.Value{{table.String("a fairly long value"), table.String("x")}, {table.String("another long one"), table.String("z")}}); err != nil {
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
	for _, r := range [][]table.Value{{table.String("a fairly long value"), table.String("x")}, {table.String("another long one"), table.String("z")}} {
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
	if err := tb.Render([][]table.Value{{table.String("a fairly long value"), table.String("x")}, {table.String("another long one"), table.String("z")}}); err != nil {
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
	for _, r := range [][]table.Value{{table.String("a fairly long value"), table.String("x")}, {table.String("another long one"), table.String("z")}} {
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
	if err := tb.Render([][]table.Value{{table.String("xxxxx"), table.String("yyyyy")}}); err != nil {
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
	if err := s.Render([]table.Value{table.String("xxxxx"), table.String("yyyyy")}); err != nil {
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
	if err := tb.Render([][]table.Value{{table.String("a fairly long value"), table.String("x")}, {table.String("another long one"), table.String("z")}}); err != nil {
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
	for _, r := range [][]table.Value{{table.String("a fairly long value"), table.String("x")}, {table.String("another long one"), table.String("z")}} {
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
	if err := tb.Render([][]table.Value{{table.String("a fairly long value"), table.Any(nil)}, {table.String("another long one"), table.String("z")}}); err != nil {
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
	for _, r := range [][]table.Value{{table.String("a fairly long value"), table.Any(nil)}, {table.String("another long one"), table.String("z")}} {
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
	if err := tb.Render([][]table.Value{
		{table.String("same-and-rather-long-label"), table.Int(1)},
		{table.String("same-and-rather-long-label"), table.Int(2)},
		{table.String("other"), table.Int(3)},
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
	for _, row := range [][]table.Value{
		{table.String("same-and-rather-long-label"), table.Int(1)},
		{table.String("same-and-rather-long-label"), table.Int(2)},
		{table.String("other"), table.Int(3)},
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
	if err := tb.Render([][]table.Value{{table.String("a fairly long value"), table.String("raw")}, {table.String("another long one"), table.String("z")}}); err != nil {
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
	for _, r := range [][]table.Value{{table.String("a fairly long value"), table.String("raw")}, {table.String("another long one"), table.String("z")}} {
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
	if err := tb.Render([][]table.Value{{table.String("xxxxx"), table.Any(strings.Repeat("y", 20))}}); err != nil {
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
	if err := s.Render([]table.Value{table.String("xxxxx"), table.Any(strings.Repeat("y", 20))}); err != nil {
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
	if err := tb.Render([][]table.Value{{table.String("xxxxx"), table.Any(strings.Repeat("y", 20))}}); err != nil {
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
	if err := s.Render([]table.Value{table.String("xxxxx"), table.Any(strings.Repeat("y", 20))}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_autofit_with_width", buf.Bytes())
}

func TestGolden_TableBandRowspanColspanBoundary(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader(
			[]string{"X", "A"},
			[]string{"A", "A"},
		),
		WithFooter(func() [][]string {
			return [][]string{{"A", "A"}, {"X", "A"}}
		}),
		WithRowspan(ScopeHeader|ScopeFooter, Columns(1)),
		WithColspan(ScopeHeader|ScopeFooter, Columns(0, 1)),
	)
	if err := tb.Render(nil); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_band_rowspan_colspan_boundary", buf.Bytes())
}

func TestGolden_StreamBandRowspanColspanBoundary(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader(
			[]string{"X", "A"},
			[]string{"A", "A"},
		),
		WithFooter(func() [][]string {
			return [][]string{{"A", "A"}, {"X", "A"}}
		}),
		WithRowspan(ScopeHeader|ScopeFooter, Columns(1)),
		WithColspan(ScopeHeader|ScopeFooter, Columns(0, 1)),
	)
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_band_rowspan_colspan_boundary", buf.Bytes())
}

func TestGolden_TableBasic(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"Name", "Value"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("foo"), table.Int(1)},
		{table.String("bar"), table.Int(2)},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_basic", buf.Bytes())
}
func TestGolden_StreamBasic(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"Name", "Value"}),
	)
	if err := s.Render([]table.Value{table.String("foo"), table.Int(1)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("bar"), table.Int(2)}); err != nil {
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
	if err := tb.Render([][]table.Value{
		{table.String("x"), table.Int(1)},
		{table.String("y"), table.Int(2)},
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
	if err := s.Render([]table.Value{table.String("x"), table.Int(1)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("y"), table.Int(2)}); err != nil {
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
	if err := tb.Render([][]table.Value{{table.String("x"), table.String("y")}}); err != nil {
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
	if err := s.Render([]table.Value{table.String("x"), table.String("y")}); err != nil {
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
	if err := tb.Render([][]table.Value{
		{table.String("x"), table.String("x"), table.String("y")},
		{table.String("p"), table.String("q"), table.String("q")},
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
	if err := s.Render([]table.Value{table.String("x"), table.String("x"), table.String("y")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("p"), table.String("q"), table.String("q")}); err != nil {
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
	if err := tb.Render([][]table.Value{{table.String("x"), table.String("x"), table.String("y")}, {table.String("p"), table.String("q"), table.String("q")}}); err != nil {
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
	for _, r := range [][]table.Value{{table.String("x"), table.String("x"), table.String("y")}, {table.String("p"), table.String("q"), table.String("q")}} {
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
	if err := tb.Render([][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}}); err != nil {
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
	for _, r := range [][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}} {
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
	if err := tb.Render([][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}}); err != nil {
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
	for _, r := range [][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}} {
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
	if err := tb.Render([][]table.Value{
		{table.String("x"), table.String("x"), table.String("y"), table.String("y")},
		{table.String("p"), table.String("q"), table.String("r"), table.String("r")},
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
	if err := s.Render([]table.Value{table.String("x"), table.String("x"), table.String("y"), table.String("y")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("p"), table.String("q"), table.String("r"), table.String("r")}); err != nil {
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
	if err := tb.Render([][]table.Value{{table.String("x"), table.String("x"), table.String("y")}, {table.String("p"), table.String("q"), table.String("q")}}); err != nil {
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
	for _, r := range [][]table.Value{{table.String("x"), table.String("x"), table.String("y")}, {table.String("p"), table.String("q"), table.String("q")}} {
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
	if err := tb.Render([][]table.Value{{table.String("x"), table.Any(nil)}, {table.Any(nil), table.String("q")}}); err != nil {
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
	for _, r := range [][]table.Value{{table.String("x"), table.Any(nil)}, {table.Any(nil), table.String("q")}} {
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
	if err := tb.Render([][]table.Value{
		{table.String("A"), table.String("A"), table.Int(1)},
		{table.String("B"), table.String("B"), table.Int(2)},
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
	if err := s.Render([]table.Value{table.String("A"), table.String("A"), table.Int(1)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("B"), table.String("B"), table.Int(2)}); err != nil {
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
	if err := tb.Render([][]table.Value{
		{table.String("g"), table.String("x"), table.String("x")},
		{table.String("h"), table.String("y"), table.String("z")},
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
	if err := s.Render([]table.Value{table.String("g"), table.String("x"), table.String("x")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("h"), table.String("y"), table.String("z")}); err != nil {
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
		WithWidth(AllColumns(), 5),
		WithTruncate(AllColumns()),
	)
	if err := tb.Render([][]table.Value{{table.String("a longer value"), table.String("a longer value")}, {table.String("a long value"), table.String("q")}}); err != nil {
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
		WithWidth(AllColumns(), 5),
		WithTruncate(AllColumns()),
	)
	for _, r := range [][]table.Value{{table.String("a longer value"), table.String("a longer value")}, {table.String("a long value"), table.String("q")}} {
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
	if err := tb.Render([][]table.Value{
		{table.String("x"), table.Int(1)},
		{table.String("y"), table.Int(2)},
		{table.String("z"), table.Int(3)},
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
	if err := s.Render([]table.Value{table.String("x"), table.Int(1)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("y"), table.Int(2)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("z"), table.Int(3)}); err != nil {
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
	if err := tb.Render([][]table.Value{
		{table.String("alice"), table.Int(100)},
		{table.String("bob"), table.Int(99)},
		{table.String("carol"), table.Int(1)},
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
	if err := s.Render([]table.Value{table.String("alice"), table.Int(100)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("bob"), table.Int(99)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("carol"), table.Int(1)}); err != nil {
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
	if err := tb.Render([][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}}); err != nil {
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
	for _, r := range [][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}} {
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
	if err := tb.Render([][]table.Value{{table.String("a fairly long value"), table.String("a fairly long value")}, {table.String("another long one"), table.String("z")}}); err != nil {
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
	for _, r := range [][]table.Value{{table.String("a fairly long value"), table.String("a fairly long value")}, {table.String("another long one"), table.String("z")}} {
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
	if err := tb.Render([][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}}); err != nil {
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
	for _, r := range [][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}} {
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
	if err := tb.Render([][]table.Value{
		{table.String("x"), table.String("x"), table.String("y")},
		{table.String("p"), table.String("p"), table.String("q")},
		{table.String("m"), table.String("n"), table.String("n")},
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
	for _, r := range [][]table.Value{
		{table.String("x"), table.String("x"), table.String("y")},
		{table.String("p"), table.String("p"), table.String("q")},
		{table.String("m"), table.String("n"), table.String("n")},
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
	if err := tb.Render([][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}}); err != nil {
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
	for _, r := range [][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}} {
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
	if err := tb.Render([][]table.Value{
		{table.Int(1), table.String("a\nb")},
		{table.Int(2), table.String("c")},
		{table.Int(3), table.String("d\ne\nf")},
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
	if err := s.Render([]table.Value{table.Int(1), table.String("a\nb")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.Int(2), table.String("c")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.Int(3), table.String("d\ne\nf")}); err != nil {
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
	if err := tb.Render([][]table.Value{
		{table.String("x"), table.Int(1)},
		{table.String("y"), table.Int(2)},
		{table.String("z"), table.Int(3)},
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
	if err := s.Render([]table.Value{table.String("x"), table.Int(1)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("y"), table.Int(2)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("z"), table.Int(3)}); err != nil {
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
	if err := tb.Render([][]table.Value{{table.String("x"), table.Any(nil)}, {table.Any(nil), table.String("q")}}); err != nil {
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
	for _, r := range [][]table.Value{{table.String("x"), table.Any(nil)}, {table.Any(nil), table.String("q")}} {
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
	if err := tb.Render([][]table.Value{
		{table.String("A"), table.String("x")},
		{table.String("A"), table.String("y")},
		{table.String("B"), table.String("z")},
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
	if err := s.Render([]table.Value{table.String("A"), table.String("x")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("A"), table.String("y")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("B"), table.String("z")}); err != nil {
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
	if err := tb.Render([][]table.Value{{table.String("x"), table.String("raw")}, {table.String("p"), table.String("q")}}); err != nil {
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
	for _, r := range [][]table.Value{{table.String("x"), table.String("raw")}, {table.String("p"), table.String("q")}} {
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
	if err := tb.Render([][]table.Value{{table.String("x"), table.String("a long value")}, {table.String("y"), table.String("short")}}); err != nil {
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
	for _, r := range [][]table.Value{{table.String("x"), table.String("a long value")}, {table.String("y"), table.String("short")}} {
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
	if err := tb.Render([][]table.Value{{table.String("x"), table.String("a long value")}, {table.String("y"), table.String("short")}}); err != nil {
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
	for _, r := range [][]table.Value{{table.String("x"), table.String("a long value")}, {table.String("y"), table.String("short")}} {
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
	if err := tb.Render([][]table.Value{
		{table.String("a\tb"), table.String("c\vd"), table.String("e\x00f")},
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
	if err := s.Render([]table.Value{table.String("a\tb"), table.String("c\vd"), table.String("e\x00f")}); err != nil {
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
	if err := tb.Render([][]table.Value{
		{table.String("\U0001F600"), table.String("\U0001F469\u200D\U0001F4BB"), table.String("e\u0301")},
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
	if err := s.Render([]table.Value{table.String("\U0001F600"), table.String("\U0001F469\u200D\U0001F4BB"), table.String("e\u0301")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_emoji", buf.Bytes())
}

func TestGolden_TableFlagWidth(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"W", "T"}),
		WithWidth(Columns(0), 2),
		WithWidth(Columns(1), 4),
		WithTruncate(Columns(1)),
	)
	if err := tb.Render([][]table.Value{{table.String("🇯🇵🇺🇸"), table.String("🇯🇵🇺🇸x")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_flag_width", buf.Bytes())
}

func TestGolden_StreamFlagWidth(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"W", "T"}),
		WithWidth(Columns(0), 2),
		WithWidth(Columns(1), 4),
		WithTruncate(Columns(1)),
	)
	if err := s.Render([]table.Value{table.String("🇯🇵🇺🇸"), table.String("🇯🇵🇺🇸x")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_flag_width", buf.Bytes())
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
	if err := tb.Render([][]table.Value{
		{table.String("x"), table.String("y"), table.String("z")},
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
	if err := s.Render([]table.Value{table.String("x"), table.String("y"), table.String("z")}); err != nil {
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
	if err := tb.Render([][]table.Value{
		{table.String("x"), table.String("y"), table.Int(1)},
		{table.String("p"), table.String("q"), table.Int(8)},
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
	if err := s.Render([]table.Value{table.String("x"), table.String("y"), table.Int(1)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("p"), table.String("q"), table.Int(8)}); err != nil {
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
	if err := tb.Render([][]table.Value{{table.String("x"), table.String("y")}, {table.String("p"), table.String("q")}}); err != nil {
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
	for _, r := range [][]table.Value{{table.String("x"), table.String("y")}, {table.String("p"), table.String("q")}} {
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
	if err := tb.Render([][]table.Value{{table.String("x"), table.Any(nil)}, {table.Any(nil), table.String("q")}}); err != nil {
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
	for _, r := range [][]table.Value{{table.String("x"), table.Any(nil)}, {table.Any(nil), table.String("q")}} {
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
	if err := tb.Render([][]table.Value{{table.String("a long value"), table.String("x")}, {table.String("y"), table.String("z")}}); err != nil {
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
	for _, r := range [][]table.Value{{table.String("a long value"), table.String("x")}, {table.String("y"), table.String("z")}} {
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

func TestGolden_TableHeaderlessColspan(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithColspan(ScopeBody, Columns(0, 1)),
	)
	if err := tb.Render([][]table.Value{{table.String("A"), table.String("A")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_headerless_colspan", buf.Bytes())
}

func TestGolden_StreamHeaderlessColspan(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithColspan(ScopeBody, Columns(0, 1)),
	)
	if err := s.Render([]table.Value{table.String("A"), table.String("A")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_headerless_colspan", buf.Bytes())
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
	if err := tb.Render([][]table.Value{
		{table.String("jp"), table.String("prod"), table.String("web")},
		{table.String("jp"), table.String("prod"), table.String("db")},
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
	for _, row := range [][]table.Value{{table.String("jp"), table.String("prod"), table.String("web")}, {table.String("jp"), table.String("prod"), table.String("db")}} {
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
	if err := tb.Render([][]table.Value{
		{table.String("a"), table.String("b"), table.Int(1)},
		{table.String("a"), table.String("b"), table.Int(2)},
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
	for _, row := range [][]table.Value{{table.String("a"), table.String("b"), table.Int(1)}, {table.String("a"), table.String("b"), table.Int(2)}} {
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
	if err := tb.Render([][]table.Value{
		{table.String("x"), table.String("y")},
		{table.String("p"), table.String("q")},
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
	if err := s.Render([]table.Value{table.String("x"), table.String("y")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("p"), table.String("q")}); err != nil {
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
	if err := tb.Render([][]table.Value{{table.String("x"), table.String("y")}, {table.String("p"), table.String("q")}}); err != nil {
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
	for _, r := range [][]table.Value{{table.String("x"), table.String("y")}, {table.String("p"), table.String("q")}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_index_width", buf.Bytes())
}

func TestGolden_TableIndexWidthAutoFit(t *testing.T) {
	restore := terminalWidth
	terminalWidth = func(io.Writer) int { return 20 }
	t.Cleanup(func() { terminalWidth = restore })
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithAutoFit(),
		WithIndexWidth(4),
	)
	if err := tb.Render([][]table.Value{{table.String("a fairly long value"), table.String("x")}, {table.String("short"), table.String("y")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_index_width_autofit", buf.Bytes())
}

func TestGolden_StreamIndexWidthAutoFit(t *testing.T) {
	restore := terminalWidth
	terminalWidth = func(io.Writer) int { return 20 }
	t.Cleanup(func() { terminalWidth = restore })
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithAutoFit(),
		WithIndexWidth(4),
	)
	for _, row := range [][]table.Value{{table.String("a fairly long value"), table.String("x")}, {table.String("short"), table.String("y")}} {
		if err := s.Render(row); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_index_width_autofit", buf.Bytes())
}

func TestGolden_TableInvalidUtf8(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("a\xffb"), table.String("\xfe"), table.String("ok")},
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
	if err := s.Render([]table.Value{table.String("a\xffb"), table.String("\xfe"), table.String("ok")}); err != nil {
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
	if err := tb.Render([][]table.Value{
		{table.String("alice"), table.String("short")},
		{table.String("bob"), table.String("this is a longer description")},
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
	if err := s.Render([]table.Value{table.String("alice"), table.String("short")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("bob"), table.String("this is a longer description")}); err != nil {
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
	if err := tb.Render([][]table.Value{
		{table.String("x"), table.String("あいうえおかきくけこ")},
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
	if err := s.Render([]table.Value{table.String("x"), table.String("あいうえおかきくけこ")}); err != nil {
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
	if err := tb.Render([][]table.Value{
		{table.String("ab"), table.String("xy")},
		{table.String("cdef"), table.String("z")},
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
	if err := s.Render([]table.Value{table.String("ab"), table.String("xy")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("cdef"), table.String("z")}); err != nil {
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
	if err := tb.Render([][]table.Value{
		{table.Int(1), table.String("line1\nthis-is-a-long-line2")},
		{table.Int(2), table.String("short")},
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
	if err := s.Render([]table.Value{table.Int(1), table.String("line1\nthis-is-a-long-line2")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.Int(2), table.String("short")}); err != nil {
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
	if err := tb.Render([][]table.Value{
		{table.String("one"), table.String("a\nb"), table.String("x\ny\nz")},
		{table.String("two\nlines"), table.String("p"), table.String("q\nr")},
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
	if err := s.Render([]table.Value{table.String("one"), table.String("a\nb"), table.String("x\ny\nz")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("two\nlines"), table.String("p"), table.String("q\nr")}); err != nil {
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
	if err := tb.Render([][]table.Value{{table.String("a\r\nb\nc\rd")}}); err != nil {
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
	if err := s.Render([]table.Value{table.String("a\r\nb\nc\rd")}); err != nil {
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
	if err := tb.Render([][]table.Value{
		{table.String("a"), table.String("b"), table.String("c")},
		{table.String("d"), table.String("e"), table.String("f")},
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
	if err := s.Render([]table.Value{table.String("a"), table.String("b"), table.String("c")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("d"), table.String("e"), table.String("f")}); err != nil {
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
	if err := tb.Render([][]table.Value{
		{table.String("alice"), table.Int(100)},
		{table.String("bob"), table.Any(nil)},
		{table.String("carol"), table.Int(99999)},
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
	if err := s.Render([]table.Value{table.String("alice"), table.Int(100)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("bob"), table.Any(nil)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("carol"), table.Int(99999)}); err != nil {
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
	if err := tb.Render([][]table.Value{
		{table.String("x"), table.Int(1)},
		{table.String("y"), table.Int(2)},
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
	if err := s.Render([]table.Value{table.String("x"), table.Int(1)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("y"), table.Int(2)}); err != nil {
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
	if err := tb.Render([][]table.Value{
		{table.Int(1), table.String("a\nb\nc")},
		{table.Int(2), table.String("x")},
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
	if err := s.Render([]table.Value{table.Int(1), table.String("a\nb\nc")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.Int(2), table.String("x")}); err != nil {
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
	if err := tb.Render([][]table.Value{
		{table.String("x"), table.Int(1)},
		{table.String("y"), table.Int(2)},
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
	if err := s.Render([]table.Value{table.String("x"), table.Int(1)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("y"), table.Int(2)}); err != nil {
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
	if err := tb.Render([][]table.Value{
		{table.String("x"), table.String("")},
		{table.Any(nil), table.String("y")},
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
	if err := s.Render([]table.Value{table.String("x"), table.String("")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.Any(nil), table.String("y")}); err != nil {
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
	if err := tb.Render([][]table.Value{{table.String("x"), table.String("")}}); err != nil {
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
	if err := s.Render([]table.Value{table.String("x"), table.String("")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_placeholder_fixed_width", buf.Bytes())
}

func TestGolden_TablePlaceholderTransformerTruncate(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithPlaceholder("missing"),
		WithWidth(Columns(1), 4),
		WithTransformer(Columns(1), func(v any) (string, *Attr) {
			if v == "raw" {
				return "replacement", nil
			}
			return "", nil
		}),
		WithTruncate(Columns(1)),
	)
	if err := tb.Render([][]table.Value{{table.String("empty"), table.String("")}, {table.String("nil"), table.Any(nil)}, {table.String("raw"), table.String("raw")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_placeholder_transformer_truncate", buf.Bytes())
}

func TestGolden_StreamPlaceholderTransformerTruncate(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"A", "B"}),
		WithPlaceholder("missing"),
		WithWidth(Columns(1), 4),
		WithTransformer(Columns(1), func(v any) (string, *Attr) {
			if v == "raw" {
				return "replacement", nil
			}
			return "", nil
		}),
		WithTruncate(Columns(1)),
	)
	for _, row := range [][]table.Value{{table.String("empty"), table.String("")}, {table.String("nil"), table.Any(nil)}, {table.String("raw"), table.String("raw")}} {
		if err := s.Render(row); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_placeholder_transformer_truncate", buf.Bytes())
}

func TestGolden_TablePlaceholderWideBytes(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithStyle(StyleLight),
		WithPlaceholder("\u2014"),
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("x"), table.Any(nil)},
		{table.String("y"), table.String("z")},
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
	for _, row := range [][]table.Value{{table.String("x"), table.Any(nil)}, {table.String("y"), table.String("z")}} {
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
	if err := tb.Render([][]table.Value{
		{table.String("x")},
		{table.String("p"), table.String("q"), table.String("r")},
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
	if err := s.Render([]table.Value{table.String("x")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("p"), table.String("q"), table.String("r")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{}); err != nil {
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
	if err := tb.Render([][]table.Value{
		{table.String("a"), table.Int(1)},
		{table.String("b"), table.Int(2)},
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
	if err := s.Render([]table.Value{table.String("a"), table.Int(1)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("b"), table.Int(2)}); err != nil {
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
	if err := tb.Render([][]table.Value{
		{table.String("A"), table.String("x")},
		{table.String("A"), table.String("y")},
		{table.String("B"), table.String("z")},
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
	if err := s.Render([]table.Value{table.String("A"), table.String("x")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("A"), table.String("y")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("B"), table.String("z")}); err != nil {
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
	if err := tb.Render([][]table.Value{
		{table.String("A"), table.Int(100)},
		{table.String("A"), table.Int(200)},
		{table.String("B"), table.Int(300)},
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
	if err := s.Render([]table.Value{table.String("A"), table.Int(100)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("A"), table.Int(200)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("B"), table.Int(300)}); err != nil {
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
	if err := tb.Render([][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}}); err != nil {
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
	for _, r := range [][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}} {
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
	if err := tb.Render([][]table.Value{{table.String("g"), table.String("a fairly long value")}, {table.String("g"), table.String("another long one")}}); err != nil {
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
	for _, r := range [][]table.Value{{table.String("g"), table.String("a fairly long value")}, {table.String("g"), table.String("another long one")}} {
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
	if err := tb.Render([][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}}); err != nil {
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
	for _, r := range [][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_rowspan_caption", buf.Bytes())
}

func TestGolden_TableRowspanColspanBoundary(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"c0", "c1", "c2"}),
		WithRowspan(ScopeBody, Columns(1, 2)),
		WithColspan(ScopeBody, Columns(0, 1, 2)),
	)
	if err := tb.Render([][]table.Value{
		{table.String("A"), table.String("A"), table.String("A")},
		{table.String("X"), table.String("A"), table.String("A")},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_rowspan_colspan_boundary", buf.Bytes())
}

func TestGolden_StreamRowspanColspanBoundary(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"c0", "c1", "c2"}),
		WithRowspan(ScopeBody, Columns(1, 2)),
		WithColspan(ScopeBody, Columns(0, 1, 2)),
	)
	if err := s.Render([]table.Value{table.String("A"), table.String("A"), table.String("A")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("X"), table.String("A"), table.String("A")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_rowspan_colspan_boundary", buf.Bytes())
}

func TestGolden_TableRowspanColspanEdge(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B", "C"}),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(1, 2)),
	)
	if err := tb.Render([][]table.Value{
		{table.String("g"), table.String("x"), table.String("x")},
		{table.String("g"), table.String("y"), table.String("y")},
		{table.String("h"), table.String("z"), table.String("w")},
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
	if err := s.Render([]table.Value{table.String("g"), table.String("x"), table.String("x")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("g"), table.String("y"), table.String("y")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("h"), table.String("z"), table.String("w")}); err != nil {
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
	if err := tb.Render([][]table.Value{
		{table.Int(1), table.String("X")},
		{table.Int(2)},
		{table.Int(3), table.Any(nil)},
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
	if err := s.Render([]table.Value{table.Int(1), table.String("X")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.Int(2)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.Int(3), table.Any(nil)}); err != nil {
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
	if err := tb.Render([][]table.Value{
		{table.String("A"), table.String("x")},
		{table.String("A"), table.String("y")},
		{table.String("B"), table.String("z")},
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
	if err := s.Render([]table.Value{table.String("A"), table.String("x")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("A"), table.String("y")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("B"), table.String("z")}); err != nil {
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
	if err := tb.Render([][]table.Value{
		{table.String("g"), table.String(""), table.String("x")},
		{table.String("g"), table.String(""), table.String("y")},
		{table.String(""), table.String(""), table.String("")},
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
	if err := s.Render([]table.Value{table.String("g"), table.String(""), table.String("x")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("g"), table.String(""), table.String("y")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String(""), table.String(""), table.String("")}); err != nil {
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
	if err := tb.Render([][]table.Value{
		{table.String("prod"), table.String("web")},
		{table.String("prod"), table.String("db")},
		{table.String("dev"), table.String("web")},
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
	for _, row := range [][]table.Value{
		{table.String("prod"), table.String("web")},
		{table.String("prod"), table.String("db")},
		{table.String("dev"), table.String("web")},
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
	if err := tb.Render([][]table.Value{{table.String("x"), table.String("raw")}, {table.String("p"), table.String("raw")}, {table.String("q"), table.String("z")}}); err != nil {
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
	for _, r := range [][]table.Value{{table.String("x"), table.String("raw")}, {table.String("p"), table.String("raw")}, {table.String("q"), table.String("z")}} {
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
	if err := tb.Render([][]table.Value{{table.String("g"), table.String("a long value")}, {table.String("g"), table.String("short")}}); err != nil {
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
	for _, r := range [][]table.Value{{table.String("g"), table.String("a long value")}, {table.String("g"), table.String("short")}} {
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
	if err := tb.Render([][]table.Value{{table.String("g"), table.String("a long value")}, {table.String("g"), table.String("short")}, {table.String("h"), table.String("x")}}); err != nil {
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
	for _, r := range [][]table.Value{{table.String("g"), table.String("a long value")}, {table.String("g"), table.String("short")}, {table.String("h"), table.String("x")}} {
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
	if err := tb.Render([][]table.Value{{table.String("v")}}); err != nil {
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
	if err := s.Render([]table.Value{table.String("v")}); err != nil {
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
	rows := make([][]table.Value, 2)
	for r := range rows {
		row := make([]table.Value, n)
		for i := range n {
			row[i] = table.String("v")
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
		row := make([]table.Value, n)
		for i := range n {
			row[i] = table.String("v")
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
	if err := tb.Render([][]table.Value{
		{table.String("vpc-1"), table.String("sub-1"), table.String("sg-1"), table.String("nacl-1"), table.String("i-001")},
		{table.String("vpc-2"), table.String("sub-2"), table.String("sg-2"), table.String("nacl-2"), table.String("i-002")},
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
	for _, row := range [][]table.Value{
		{table.String("vpc-1"), table.String("sub-1"), table.String("sg-1"), table.String("nacl-1"), table.String("i-001")},
		{table.String("vpc-2"), table.String("sub-2"), table.String("sg-2"), table.String("nacl-2"), table.String("i-002")},
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
	if err := tb.Render([][]table.Value{{table.Int(1), table.Int(2), table.Int(3), table.Int(4), table.String("x")}}); err != nil {
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
	if err := s.Render([]table.Value{table.Int(1), table.Int(2), table.Int(3), table.Int(4), table.String("x")}); err != nil {
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
	if err := tb.Render([][]table.Value{{table.String("us-east-1a"), table.String("a"), table.String("1")}}); err != nil {
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
	if err := s.Render([]table.Value{table.String("us-east-1a"), table.String("a"), table.String("1")}); err != nil {
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
	if err := tb.Render([][]table.Value{{table.String("us-east-1a"), table.String("a"), table.String("1")}}); err != nil {
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
	if err := s.Render([]table.Value{table.String("us-east-1a"), table.String("a"), table.String("1")}); err != nil {
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
	if err := tb.Render([][]table.Value{
		{table.String("foo"), table.Int(1)},
		{table.String("foo"), table.Int(2)},
		{table.String("bar"), table.Int(3)},
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
	for _, row := range [][]table.Value{{table.String("foo"), table.Int(1)}, {table.String("foo"), table.Int(2)}, {table.String("bar"), table.Int(3)}} {
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
	if err := tb.Render([][]table.Value{
		{table.String("foo"), table.Int(1)},
		{table.String("bar"), table.Int(2)},
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
	for _, row := range [][]table.Value{{table.String("foo"), table.Int(1)}, {table.String("bar"), table.Int(2)}} {
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
	if err := tb.Render([][]table.Value{
		{table.String("foo"), table.Int(1)},
		{table.String("bar"), table.Int(2)},
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
	for _, row := range [][]table.Value{{table.String("foo"), table.Int(1)}, {table.String("bar"), table.Int(2)}} {
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
	if err := tb.Render([][]table.Value{
		{table.String("foo"), table.Int(1)},
		{table.String("bar"), table.Int(2)},
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
	if err := s.Render([]table.Value{table.String("foo"), table.Int(1)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("bar"), table.Int(2)}); err != nil {
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
	if err := tb.Render([][]table.Value{
		{table.String("foo"), table.Int(1)},
		{table.String("bar"), table.Int(2)},
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
	if err := s.Render([]table.Value{table.String("foo"), table.Int(1)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("bar"), table.Int(2)}); err != nil {
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
	if err := tb.Render([][]table.Value{
		{table.String("foo"), table.Int(1)},
		{table.String("bar"), table.Int(2)},
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
	if err := s.Render([]table.Value{table.String("foo"), table.Int(1)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("bar"), table.Int(2)}); err != nil {
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
	if err := tb.Render([][]table.Value{
		{table.String("foo"), table.Int(1), table.String("x")},
		{table.String("bar"), table.Int(2), table.String("y")},
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
	for _, row := range [][]table.Value{{table.String("foo"), table.Int(1), table.String("x")}, {table.String("bar"), table.Int(2), table.String("y")}} {
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
	if err := tb.Render([][]table.Value{
		{table.String("x"), table.String("x"), table.String("y")},
		{table.String("p"), table.String("p"), table.String("q")},
		{table.String("m"), table.String("n"), table.String("n")},
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
	for _, r := range [][]table.Value{
		{table.String("x"), table.String("x"), table.String("y")},
		{table.String("p"), table.String("p"), table.String("q")},
		{table.String("m"), table.String("n"), table.String("n")},
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
	if err := tb.Render([][]table.Value{
		{table.String("foo"), table.Int(1)},
		{table.String("bar"), table.Int(2)},
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
	if err := s.Render([]table.Value{table.String("foo"), table.Int(1)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("bar"), table.Int(2)}); err != nil {
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
	if err := tb.Render([][]table.Value{
		{table.String("alice"), table.Int(100)},
		{table.String("bob"), table.Int(99)},
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
	if err := s.Render([]table.Value{table.String("alice"), table.Int(100)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("bob"), table.Int(99)}); err != nil {
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
	if err := tb.Render([][]table.Value{
		{table.String("1"), table.String("ok")},
		{table.String("2"), table.String("warn")},
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
	for _, r := range [][]table.Value{{table.String("1"), table.String("ok")}, {table.String("2"), table.String("warn")}} {
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
	if err := tb.Render([][]table.Value{
		{table.String("alice"), table.String("short")},
		{table.String("bob"), table.String("this is a longer description")},
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
	if err := s.Render([]table.Value{table.String("alice"), table.String("short")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("bob"), table.String("this is a longer description")}); err != nil {
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
	if err := tb.Render([][]table.Value{
		{table.String("x"), table.String("あいうえおかきくけこ")},
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
	if err := s.Render([]table.Value{table.String("x"), table.String("あいうえおかきくけこ")}); err != nil {
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
	if err := tb.Render([][]table.Value{{table.String("x"), table.String("ab\ncd")}}); err != nil {
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
	if err := s.Render([]table.Value{table.String("x"), table.String("ab\ncd")}); err != nil {
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
	if err := tb.Render([][]table.Value{
		{table.Int(1), table.String("line1\nthis-is-a-long-line2")},
		{table.Int(2), table.String("short")},
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
	if err := s.Render([]table.Value{table.Int(1), table.String("line1\nthis-is-a-long-line2")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.Int(2), table.String("short")}); err != nil {
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
	if err := tb.Render([][]table.Value{
		{table.String("x"), table.String("long value")},
		{table.String("y"), table.String("z")},
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
	for _, row := range [][]table.Value{{table.String("x"), table.String("long value")}, {table.String("y"), table.String("z")}} {
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
	if err := tb.Render([][]table.Value{
		{table.String(""), table.Any(nil), table.String("x")},
		{table.Any(nil), table.String(""), table.String("")},
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
	if err := s.Render([]table.Value{table.String(""), table.Any(nil), table.String("x")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.Any(nil), table.String(""), table.String("")}); err != nil {
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
	if err := tb.Render([][]table.Value{{table.Float32(float32(3.14)), table.Float64(float64(2.71828))}}); err != nil {
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
	if err := s.Render([]table.Value{table.Float32(float32(3.14)), table.Float64(float64(2.71828))}); err != nil {
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
	if err := tb.Render([][]table.Value{
		{table.Int(int(-1)), table.Int8(int8(-2)), table.Int16(int16(-3)), table.Int32(int32(-4)), table.Int64(int64(-5)),
			table.Uint(uint(1)), table.Uint8(uint8(2)), table.Uint16(uint16(3)), table.Uint32(uint32(4)), table.Uint64(uint64(5))},
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
	if err := s.Render([]table.Value{
		table.Int(int(-1)), table.Int8(int8(-2)), table.Int16(int16(-3)), table.Int32(int32(-4)), table.Int64(int64(-5)),
		table.Uint(uint(1)), table.Uint8(uint8(2)), table.Uint16(uint16(3)), table.Uint32(uint32(4)), table.Uint64(uint64(5)),
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
	if err := tb.Render([][]table.Value{{table.Any(&str), table.Any(&s)}}); err != nil {
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
	if err := s.Render([]table.Value{table.Any(&str), table.Any(&st)}); err != nil {
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
	if err := tb.Render([][]table.Value{{table.Any(testutil.Stringer{Value: "x"}), table.Any(testutil.Error{Value: "boom"})}}); err != nil {
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
	if err := s.Render([]table.Value{table.Any(testutil.Stringer{Value: "x"}), table.Any(testutil.Error{Value: "boom"})}); err != nil {
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
	if err := tb.Render([][]table.Value{{table.Any(nilStringer), table.Any(nilError)}}); err != nil {
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
	if err := s.Render([]table.Value{table.Any(nilStringer), table.Any(nilError)}); err != nil {
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
	if err := tb.Render([][]table.Value{
		{table.String("N/A"), table.String("")},
		{table.String("x"), table.String("N/A")},
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
	if err := s.Render([]table.Value{table.String("N/A"), table.String("")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("x"), table.String("N/A")}); err != nil {
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
	if err := tb.Render([][]table.Value{{table.String("a long value"), table.String("yy")}, {table.String("y"), table.String("z")}}); err != nil {
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
	for _, r := range [][]table.Value{{table.String("a long value"), table.String("yy")}, {table.String("y"), table.String("z")}} {
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
	if err := tb.Render([][]table.Value{{table.String("a long value"), table.String("x")}, {table.String("y"), table.String("z")}}); err != nil {
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
	for _, r := range [][]table.Value{{table.String("a long value"), table.String("x")}, {table.String("y"), table.String("z")}} {
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
	if err := tb.Render([][]table.Value{{table.String("a long value"), table.String("x")}, {table.String("y"), table.String("z")}}); err != nil {
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
	for _, r := range [][]table.Value{{table.String("a long value"), table.String("x")}, {table.String("y"), table.String("z")}} {
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
	if err := tb.Render([][]table.Value{{table.String("a long value"), table.String("x")}, {table.String("y"), table.String("z")}}); err != nil {
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
	for _, r := range [][]table.Value{{table.String("a long value"), table.String("x")}, {table.String("y"), table.String("z")}} {
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
	if err := tb.Render([][]table.Value{{table.String("a long value"), table.String("x")}, {table.String("y"), table.String("z")}}); err != nil {
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
	for _, r := range [][]table.Value{{table.String("a long value"), table.String("x")}, {table.String("y"), table.String("z")}} {
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
	if err := tb.Render([][]table.Value{{table.String("a long value"), table.String("raw")}, {table.String("y"), table.String("q")}}); err != nil {
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
	for _, r := range [][]table.Value{{table.String("a long value"), table.String("raw")}, {table.String("y"), table.String("q")}} {
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
	if err := tb.Render([][]table.Value{
		{table.String("a\u200Bb"), table.String("\uFEFF"), table.String("x\u200Dy")},
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
	if err := s.Render([]table.Value{table.String("a\u200Bb"), table.String("\uFEFF"), table.String("x\u200Dy")}); err != nil {
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
	if err := s.Render([]table.Value{table.String("a"), table.String("b"), table.Int(1)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("longer"), table.String("y"), table.Int(22)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("x"), table.String("longest"), table.Int(333)}); err != nil {
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
	if err := s.Render([]table.Value{table.String("xxxxx"), table.Any(strings.Repeat("y", 40))}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("z"), table.Any(strings.Repeat("w", 20))}); err != nil {
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
	if err := s.Render([]table.Value{table.String("あいう"), table.String("abc")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("日本"), table.String("longer-text")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("テスト"), table.String("x")}); err != nil {
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
	for _, r := range [][]table.Value{{table.String("a fairly long value"), table.String("a fairly long value")}, {table.String("another long one"), table.String("z")}} {
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
	for _, r := range [][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}} {
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
	for _, r := range [][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}} {
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
	if err := s.Render([]table.Value{table.String("nil"), table.Any(nil)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("empty"), table.String("")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("space"), table.String(" ")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("text"), table.String("hello")}); err != nil {
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
		WithWidth(Columns(0), 1),
	)
	for _, v := range []string{"x", "\u3042", "y"} {
		if err := s.Render([]table.Value{table.Any(v)}); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_frozen_overflow", buf.Bytes())
}

func TestGolden_StreamFrozenZeroWidth(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithHeader([]string{"", "B"}),
		WithPlaceholder(""),
	)
	if err := s.Render([]table.Value{table.String(""), table.String("x")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("wide"), table.String("y")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_frozen_zero_width", buf.Bytes())
}

func TestGolden_StreamIndex(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithStyle(StyleLight),
		WithIndex(),
		WithHeader([]string{"Name", "Score"}),
	)
	if err := s.Render([]table.Value{table.String("alice"), table.Int(100)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("bob"), table.Int(99)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("carol"), table.Int(98)}); err != nil {
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
	for _, r := range [][]table.Value{{table.String("x"), table.String("yy")}, {table.String("p"), table.String("z")}} {
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
	for _, r := range [][]table.Value{{table.String("x"), table.String("y")}, {table.String("p"), table.String("q")}} {
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
	for _, r := range [][]table.Value{{table.String("x"), table.String("y")}, {table.String("p"), table.String("q")}} {
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
	for _, row := range [][]table.Value{{table.String("alice"), table.Int(100)}, {table.String("bob"), table.Int(99)}, {table.String("carol"), table.Int(98)}} {
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
	for _, r := range [][]table.Value{{table.String("x"), table.String("y")}, {table.String("p"), table.String("q")}} {
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
	for _, r := range [][]table.Value{{table.String("x"), table.String("raw")}, {table.String("p"), table.String("q")}} {
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
	for _, r := range [][]table.Value{{table.String("a long value"), table.String("x")}, {table.String("y"), table.String("z")}} {
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
	if err := s.Render([]table.Value{table.String("s"), table.Int(1)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.Any(strings.Repeat("longvalue", 40)), table.Int(2)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("t"), table.Int(3)}); err != nil {
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
	for _, r := range [][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}} {
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
	if err := s.Render([]table.Value{table.Int(100), table.String("a")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.Int(100), table.String("b")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.Int(200), table.String("c")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.Int(300), table.String("d")}); err != nil {
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
	if err := s.Render([]table.Value{table.String("a")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("bb")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("ccc")}); err != nil {
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
	if err := s.Render([]table.Value{table.String("foo"), table.Int(1)}); err != nil {
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
	if err := s.Render([]table.Value{table.String("less-than"), table.String("<")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("pipe"), table.String("|")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("backslash"), table.String("\\")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("ampersand"), table.String("&")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("asterisk"), table.String("*")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("underscore"), table.String("_")}); err != nil {
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
	if err := s.Render([]table.Value{table.String("text"), table.String("hello")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("number"), table.Int(42)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("big-num"), table.Int(100000)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("empty"), table.Any(nil)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("neg"), table.Any(-7)}); err != nil {
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
	if err := s.Render([]table.Value{table.String("[]int"), table.Any([]int{1, 2, 3})}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("[]float64"), table.Any([]float64{1.1, 2.2})}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("[]bool"), table.Any([]bool{true, false, true})}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("[]byte"), table.Any([]byte("hello"))}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("[]byte (raw)"), table.Any([]byte{0x00, 0xff})}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("[]string"), table.Any([]string{"a", "b"})}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("[]any"), table.Any([]any{"x", 1, nil})}); err != nil {
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
	if err := s.Render([]table.Value{table.Int(1), table.Int(1)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.Int(2), table.Int(1000000)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.Int(3), table.Any(-99)}); err != nil {
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
	for _, r := range [][]table.Value{{table.String("a long value"), table.String("x")}, {table.String("y"), table.String("z")}} {
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
	if err := tb.Render([][]table.Value{
		{table.String("a"), table.String("b"), table.Int(1)},
		{table.String("longer"), table.String("y"), table.Int(22)},
		{table.String("x"), table.String("longest"), table.Int(333)},
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
	if err := tb.Render([][]table.Value{
		{table.String("xxxxx"), table.Any(strings.Repeat("y", 40))},
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
	if err := tb.Render([][]table.Value{
		{table.Any(strings.Repeat("x", 20)), table.Any(strings.Repeat("y", 21))},
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
	if err := tb.Render([][]table.Value{
		{table.String("あいう"), table.String("abc")},
		{table.String("日本"), table.String("longer-text")},
		{table.String("テスト"), table.String("x")},
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
	if err := tb.Render([][]table.Value{{table.String("a fairly long value"), table.String("a fairly long value")}, {table.String("another long one"), table.String("z")}}); err != nil {
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
	if err := tb.Render([][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}}); err != nil {
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
	if err := tb.Render([][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}}); err != nil {
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
	if err := tb.Render([][]table.Value{
		{table.String("nil"), table.Any(nil)},
		{table.String("empty"), table.String("")},
		{table.String("space"), table.String(" ")},
		{table.String("text"), table.String("hello")},
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
	if err := tb.Render([][]table.Value{
		{table.String("alice"), table.Int(100)},
		{table.String("bob"), table.Int(99)},
		{table.String("carol"), table.Int(98)},
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
	if err := tb.Render([][]table.Value{{table.String("x"), table.String("yy")}, {table.String("p"), table.String("z")}}); err != nil {
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
	if err := tb.Render([][]table.Value{{table.String("x"), table.String("y")}, {table.String("p"), table.String("q")}}); err != nil {
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
	if err := tb.Render([][]table.Value{{table.String("x"), table.String("y")}, {table.String("p"), table.String("q")}}); err != nil {
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
	if err := tb.Render([][]table.Value{
		{table.String("alice"), table.Int(100)},
		{table.String("bob"), table.Int(99)},
		{table.String("carol"), table.Int(98)},
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
	if err := tb.Render([][]table.Value{{table.String("x"), table.String("y")}, {table.String("p"), table.String("q")}}); err != nil {
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
	if err := tb.Render([][]table.Value{{table.String("x"), table.String("raw")}, {table.String("p"), table.String("q")}}); err != nil {
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
	if err := tb.Render([][]table.Value{{table.String("a long value"), table.String("x")}, {table.String("y"), table.String("z")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_index_truncate", buf.Bytes())
}

func TestGolden_TableLongValue(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("s"), table.Int(1)},
		{table.Any(strings.Repeat("longvalue", 40)), table.Int(2)},
		{table.String("t"), table.Int(3)},
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
	if err := tb.Render([][]table.Value{
		{table.String("a"), table.Int(1)},
		{table.String("b")},
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
	if err := tb.Render([][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}}); err != nil {
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
	if err := tb.Render([][]table.Value{
		{table.String("a")},
		{table.String("bb")},
		{table.String("ccc")},
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
	if err := tb.Render([][]table.Value{
		{table.String("less-than"), table.String("<")},
		{table.String("pipe"), table.String("|")},
		{table.String("backslash"), table.String("\\")},
		{table.String("ampersand"), table.String("&")},
		{table.String("asterisk"), table.String("*")},
		{table.String("underscore"), table.String("_")},
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
	if err := tb.Render([][]table.Value{
		{table.String("text"), table.String("hello")},
		{table.String("number"), table.Int(42)},
		{table.String("big-num"), table.Int(100000)},
		{table.String("empty"), table.Any(nil)},
		{table.String("neg"), table.Any(-7)},
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
	if err := tb.Render([][]table.Value{
		{table.String("[]int"), table.Any([]int{1, 2, 3})},
		{table.String("[]float64"), table.Any([]float64{1.1, 2.2})},
		{table.String("[]bool"), table.Any([]bool{true, false, true})},
		{table.String("[]byte"), table.Any([]byte("hello"))},
		{table.String("[]byte (raw)"), table.Any([]byte{0x00, 0xff})},
		{table.String("[]string"), table.Any([]string{"a", "b"})},
		{table.String("[]any"), table.Any([]any{"x", 1, nil})},
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
	if err := tb.Render([][]table.Value{
		{table.Int(1), table.Int(1)},
		{table.Int(2), table.Int(1000000)},
		{table.Int(3), table.Any(-99)},
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
	if err := tb.Render([][]table.Value{{table.String("a long value"), table.String("x")}, {table.String("y"), table.String("z")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_width_index", buf.Bytes())
}
