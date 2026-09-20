package html

import (
	"bytes"
	"strings"
	"testing"

	"github.com/nekrassov01/table"
	"github.com/nekrassov01/table/internal/param"
	"github.com/nekrassov01/table/internal/testutil"
)

func TestGolden_TableAlignColspan(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), AlignRight),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1)),
	)
	if err := tb.Render([][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_align_colspan", buf.Bytes())
}

func TestGolden_StreamAlignColspan(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), AlignRight),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1)),
	)
	for _, r := range [][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_align_colspan", buf.Bytes())
}

func TestGolden_TableAlignRow(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"Left", "Right", "Center"}),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), AlignRight),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(2), AlignCenter),
	)
	if err := tb.Render([][]table.Value{
		{table.String("a"), table.Int(100), table.String("x")},
		{table.String("b"), table.Int(200), table.String("y")},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_align_row", buf.Bytes())
}

func TestGolden_StreamAlignRow(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"Left", "Right", "Center"}),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), AlignRight),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(2), AlignCenter),
	)
	if err := s.Render([]table.Value{table.String("a"), table.Int(100), table.String("x")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("b"), table.Int(200), table.String("y")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_align_row", buf.Bytes())
}

func TestGolden_TableAlignScope(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
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

func TestGolden_TableAlignTransformer(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), AlignRight),
		WithTransformer(Columns(1), func(v table.Value) (string, *Color, *Decoration) {
			if s, ok := v.AsAny().(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
		}),
	)
	if err := tb.Render([][]table.Value{{table.String("x"), table.String("raw")}, {table.String("p"), table.String("q")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_align_transformer", buf.Bytes())
}

func TestGolden_StreamAlignTransformer(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), AlignRight),
		WithTransformer(Columns(1), func(v table.Value) (string, *Color, *Decoration) {
			if s, ok := v.AsAny().(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
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
	testutil.AssertGolden(t, "common_align_transformer", buf.Bytes())
}

func TestGolden_TableAllPlaceholder(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
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

func TestGolden_TableBold(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithDecoration(ScopeBody, Columns(0), DecorationBold),
		WithHeader([]string{"Name", "Value"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("foo"), table.Int(1)},
		{table.String("bar"), table.Int(2)},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_bold", buf.Bytes())
}

func TestGolden_StreamBold(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithDecoration(ScopeBody, Columns(0), DecorationBold),
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
	testutil.AssertGolden(t, "common_bold", buf.Bytes())
}

func TestGolden_TableCaptionBottom(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithCaption("summary", CaptionBottom),
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]table.Value{{table.String("x"), table.String("y")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_caption_bottom", buf.Bytes())
}

func TestGolden_StreamCaptionBottom(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithCaption("summary", CaptionBottom),
		WithHeader([]string{"A", "B"}),
	)
	if err := s.Render([]table.Value{table.String("x"), table.String("y")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_caption_bottom", buf.Bytes())
}

func TestGolden_TableCaptionCellAttr(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithCaption("cap", CaptionBottom),
		WithCellAttr(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), Attr{Class: "c"}),
	)
	if err := tb.Render([][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_caption_cellattr", buf.Bytes())
}

func TestGolden_StreamCaptionCellAttr(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
		WithCaption("cap", CaptionBottom),
		WithCellAttr(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), Attr{Class: "c"}),
	)
	for _, r := range [][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_caption_cellattr", buf.Bytes())
}

func TestGolden_TableCaptionColor(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithCaption("cap", CaptionBottom),
		WithColor(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), ColorFgRed),
	)
	if err := tb.Render([][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_caption_color", buf.Bytes())
}

func TestGolden_StreamCaptionColor(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
		WithCaption("cap", CaptionBottom),
		WithColor(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), ColorFgRed),
	)
	for _, r := range [][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_caption_color", buf.Bytes())
}

func TestGolden_TableCaptionColspan(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithCaption("cap", CaptionBottom),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1)),
	)
	if err := tb.Render([][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_caption_colspan", buf.Bytes())
}

func TestGolden_StreamCaptionColspan(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
		WithCaption("cap", CaptionBottom),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1)),
	)
	for _, r := range [][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_caption_colspan", buf.Bytes())
}

func TestGolden_TableCaptionDecoration(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithCaption("cap", CaptionBottom),
		WithDecoration(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), DecorationBold),
	)
	if err := tb.Render([][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_caption_decoration", buf.Bytes())
}

func TestGolden_StreamCaptionDecoration(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
		WithCaption("cap", CaptionBottom),
		WithDecoration(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), DecorationBold),
	)
	for _, r := range [][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_caption_decoration", buf.Bytes())
}

func TestGolden_TableCaptionEscape(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A"}),
		WithCaption("A & B <C>", CaptionDefault),
	)
	if err := tb.Render([][]table.Value{{table.String("x")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_caption_escape", buf.Bytes())
}

func TestGolden_StreamCaptionEscape(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A"}),
		WithCaption("A & B <C>", CaptionDefault),
	)
	if err := s.Render([]table.Value{table.String("x")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_caption_escape", buf.Bytes())
}

func TestGolden_TableCaptionPlaceholder(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithCaption("cap", CaptionBottom),
		WithPlaceholder("-"),
	)
	if err := tb.Render([][]table.Value{{table.String("x"), table.Any(nil)}, {table.Any(nil), table.String("q")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_caption_placeholder", buf.Bytes())
}

func TestGolden_StreamCaptionPlaceholder(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
		WithCaption("cap", CaptionBottom),
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
	testutil.AssertGolden(t, "common_caption_placeholder", buf.Bytes())
}

func TestGolden_TableCaptionTop(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
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

func TestGolden_TableCaptionTransformer(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithCaption("cap", CaptionBottom),
		WithTransformer(Columns(1), func(v table.Value) (string, *Color, *Decoration) {
			if s, ok := v.AsAny().(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
		}),
	)
	if err := tb.Render([][]table.Value{{table.String("x"), table.String("raw")}, {table.String("p"), table.String("q")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_caption_transformer", buf.Bytes())
}

func TestGolden_StreamCaptionTransformer(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
		WithCaption("cap", CaptionBottom),
		WithTransformer(Columns(1), func(v table.Value) (string, *Color, *Decoration) {
			if s, ok := v.AsAny().(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
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
	testutil.AssertGolden(t, "common_caption_transformer", buf.Bytes())
}

func TestGolden_TableCellAttrScope(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"Name", "Path"}),
		WithFooter(func() [][]string {
			return [][]string{{"", "2 rows"}}
		}),
		WithCellAttr(ScopeBody, Columns(1), Attr{Class: "mono", Style: "white-space:pre-wrap"}),
		WithCellAttr(ScopeFooter, Columns(1), Attr{Class: "total"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("a"), table.String("  /etc/hosts")},
		{table.String("b"), table.String("  /etc/passwd")},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_cell_attr_scope", buf.Bytes())
}

func TestGolden_StreamCellAttrScope(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"Name", "Path"}),
		WithFooter(func() [][]string {
			return [][]string{{"", "2 rows"}}
		}),
		WithCellAttr(ScopeBody, Columns(1), Attr{Class: "mono", Style: "white-space:pre-wrap"}),
		WithCellAttr(ScopeFooter, Columns(1), Attr{Class: "total"}),
	)
	for _, row := range [][]table.Value{{table.String("a"), table.String("  /etc/hosts")}, {table.String("b"), table.String("  /etc/passwd")}} {
		if err := s.Render(row); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_cell_attr_scope", buf.Bytes())
}

func TestGolden_TableCellClass(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithTableAttr(TableAttr{
			Body: SectionAttr{Cell: Attr{Class: "base"}},
		}),
		WithCellAttr(ScopeHeader|ScopeBody|ScopeFooter, Columns(0), Attr{Class: "col-a"}),
		WithCellAttr(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), Attr{Class: "col-b"}),
	)
	if err := tb.Render([][]table.Value{{table.String("x"), table.String("y")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_cell_class", buf.Bytes())
}

func TestGolden_StreamCellClass(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
		WithTableAttr(TableAttr{
			Body: SectionAttr{Cell: Attr{Class: "base"}},
		}),
		WithCellAttr(ScopeHeader|ScopeBody|ScopeFooter, Columns(0), Attr{Class: "col-a"}),
		WithCellAttr(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), Attr{Class: "col-b"}),
	)
	if err := s.Render([]table.Value{table.String("x"), table.String("y")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_cell_class", buf.Bytes())
}

func TestGolden_TableCellStyle(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithCellAttr(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), Attr{Class: "score", Style: "color:#333"}),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), AlignRight),
		WithHeader([]string{"Name", "Score"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("alice"), table.Int(100)},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_cell_style", buf.Bytes())
}

func TestGolden_StreamCellStyle(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithCellAttr(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), Attr{Class: "score", Style: "color:#333"}),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), AlignRight),
		WithHeader([]string{"Name", "Score"}),
	)
	if err := s.Render([]table.Value{table.String("alice"), table.Int(100)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_cell_style", buf.Bytes())
}

func TestGolden_TableCellAttrTransformer(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithCellAttr(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), Attr{Class: "c"}),
		WithTransformer(Columns(1), func(v table.Value) (string, *Color, *Decoration) {
			if s, ok := v.AsAny().(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
		}),
	)
	if err := tb.Render([][]table.Value{{table.String("x"), table.String("raw")}, {table.String("p"), table.String("q")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_cellattr_transformer", buf.Bytes())
}

func TestGolden_StreamCellAttrTransformer(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
		WithCellAttr(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), Attr{Class: "c"}),
		WithTransformer(Columns(1), func(v table.Value) (string, *Color, *Decoration) {
			if s, ok := v.AsAny().(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
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
	testutil.AssertGolden(t, "common_cellattr_transformer", buf.Bytes())
}

func TestGolden_TableCJK(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"日本語", "ASCII"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("あいう"), table.String("abc")},
		{table.String("日本"), table.String("longer-text")},
		{table.String("テスト"), table.String("x")},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_cjk", buf.Bytes())
}

func TestGolden_StreamCJK(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
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
	testutil.AssertGolden(t, "common_cjk", buf.Bytes())
}

func TestGolden_TableClass(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithTableAttr(TableAttr{
			Table:   Attr{Class: "tbl"},
			Caption: Attr{Class: "cap"},
			Header:  SectionAttr{Section: Attr{Class: "hd"}, Row: Attr{Class: "hr"}, Cell: Attr{Class: "hc"}},
			Body:    SectionAttr{Section: Attr{Class: "bd"}, Row: Attr{Class: "br"}, Cell: Attr{Class: "bc"}},
		}),
		WithCaption("Title", CaptionDefault),
	)
	if err := tb.Render([][]table.Value{{table.String("x"), table.String("y")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_class", buf.Bytes())
}

func TestGolden_StreamClass(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
		WithTableAttr(TableAttr{
			Table:   Attr{Class: "tbl"},
			Caption: Attr{Class: "cap"},
			Header:  SectionAttr{Section: Attr{Class: "hd"}, Row: Attr{Class: "hr"}, Cell: Attr{Class: "hc"}},
			Body:    SectionAttr{Section: Attr{Class: "bd"}, Row: Attr{Class: "br"}, Cell: Attr{Class: "bc"}},
		}),
		WithCaption("Title", CaptionDefault),
	)
	if err := s.Render([]table.Value{table.String("x"), table.String("y")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_class", buf.Bytes())
}

func TestGolden_TableClassEscape(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A"}),
		WithTableAttr(TableAttr{Table: Attr{Class: `x" onmouseover="alert(1)`}}),
		WithCellAttr(ScopeHeader|ScopeBody|ScopeFooter, Columns(0), Attr{Style: `color:red" onclick="evil()`}),
	)
	if err := tb.Render([][]table.Value{{table.String("v")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_class_escape", buf.Bytes())
}

func TestGolden_StreamClassEscape(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A"}),
		WithTableAttr(TableAttr{Table: Attr{Class: `x" onmouseover="alert(1)`}}),
		WithCellAttr(ScopeHeader|ScopeBody|ScopeFooter, Columns(0), Attr{Style: `color:red" onclick="evil()`}),
	)
	if err := s.Render([]table.Value{table.String("v")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_class_escape", buf.Bytes())
}

func TestGolden_TableCode(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithDecoration(ScopeBody, Columns(1), DecorationCode),
		WithHeader([]string{"Name", "ID"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("alice"), table.String("id-001")},
		{table.String("bob"), table.String("id-002")},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_code", buf.Bytes())
}

func TestGolden_StreamCode(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithDecoration(ScopeBody, Columns(1), DecorationCode),
		WithHeader([]string{"Name", "ID"}),
	)
	if err := s.Render([]table.Value{table.String("alice"), table.String("id-001")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("bob"), table.String("id-002")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_code", buf.Bytes())
}

func TestGolden_TableColor(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithColor(ScopeBody, Columns(1), NewColor("red", "")),
		WithHeader([]string{"Name", "Value"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("foo"), table.String("hello")},
		{table.String("bar"), table.String("world")},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_color", buf.Bytes())
}

func TestGolden_StreamColor(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithColor(ScopeBody, Columns(1), NewColor("red", "")),
		WithHeader([]string{"Name", "Value"}),
	)
	if err := s.Render([]table.Value{table.String("foo"), table.String("hello")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("bar"), table.String("world")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_color", buf.Bytes())
}

func TestGolden_TableColorAlign(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithColor(ScopeBody, Columns(1), ColorFgRed),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), AlignRight),
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("x"), table.String("y"), table.Int(1)},
		{table.String("p"), table.String("qq"), table.Int(22)},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_color_align", buf.Bytes())
}

func TestGolden_StreamColorAlign(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithColor(ScopeBody, Columns(1), ColorFgRed),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), AlignRight),
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := s.Render([]table.Value{table.String("x"), table.String("y"), table.Int(1)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("p"), table.String("qq"), table.Int(22)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_color_align", buf.Bytes())
}

func TestGolden_TableColorAttrEscape(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A"}),
		WithColor(ScopeBody, Columns(0), NewColor(`red" onmouseover="alert(1)`, `blue&x`)),
	)
	if err := tb.Render([][]table.Value{
		{table.String("x")},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_color_attr_escape", buf.Bytes())
}

func TestGolden_StreamColorAttrEscape(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A"}),
		WithColor(ScopeBody, Columns(0), NewColor(`red" onmouseover="alert(1)`, `blue&x`)),
	)
	if err := s.Render([]table.Value{table.String("x")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_color_attr_escape", buf.Bytes())
}

func TestGolden_TableColorBg(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithColor(ScopeBody, Columns(1), NewColor("", "red")),
		WithHeader([]string{"Name", "Value"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("foo"), table.String("hello")},
		{table.String("bar"), table.String("world")},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_color_bg", buf.Bytes())
}

func TestGolden_StreamColorBg(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithColor(ScopeBody, Columns(1), NewColor("", "red")),
		WithHeader([]string{"Name", "Value"}),
	)
	if err := s.Render([]table.Value{table.String("foo"), table.String("hello")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("bar"), table.String("world")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_color_bg", buf.Bytes())
}

func TestGolden_TableColorCellAttr(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithColor(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), ColorFgRed),
		WithCellAttr(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), Attr{Class: "c"}),
	)
	if err := tb.Render([][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_color_cellattr", buf.Bytes())
}

func TestGolden_StreamColorCellAttr(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
		WithColor(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), ColorFgRed),
		WithCellAttr(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), Attr{Class: "c"}),
	)
	for _, r := range [][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_color_cellattr", buf.Bytes())
}

func TestGolden_TableColorCode(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithDecoration(ScopeBody, Columns(1), DecorationCode),
		WithColor(ScopeBody, Columns(1), NewColor("red", "")),
		WithHeader([]string{"Type", "Value"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("text"), table.String("hello")},
		{table.String("slice"), table.Any([]int{1, 2, 3})},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_color_code", buf.Bytes())
}

func TestGolden_StreamColorCode(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithDecoration(ScopeBody, Columns(1), DecorationCode),
		WithColor(ScopeBody, Columns(1), NewColor("red", "")),
		WithHeader([]string{"Type", "Value"}),
	)
	if err := s.Render([]table.Value{table.String("text"), table.String("hello")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("slice"), table.Any([]int{1, 2, 3})}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_color_code", buf.Bytes())
}

func TestGolden_TableColorEscape(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithColor(ScopeBody, Columns(0), ColorFgRed),
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("<b>"), table.String("&amp;")},
		{table.String("a|b"), table.String("c\"d")},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_color_escape", buf.Bytes())
}

func TestGolden_StreamColorEscape(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithColor(ScopeBody, Columns(0), ColorFgRed),
		WithHeader([]string{"A", "B"}),
	)
	if err := s.Render([]table.Value{table.String("<b>"), table.String("&amp;")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("a|b"), table.String("c\"d")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_color_escape", buf.Bytes())
}

func TestGolden_TableColorFgBg(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithColor(ScopeBody, Columns(1), NewColor("red", "blue")),
		WithHeader([]string{"Name", "Value"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("foo"), table.String("hello")},
		{table.String("bar"), table.String("world")},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_color_fg_bg", buf.Bytes())
}

func TestGolden_StreamColorFgBg(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithColor(ScopeBody, Columns(1), NewColor("red", "blue")),
		WithHeader([]string{"Name", "Value"}),
	)
	if err := s.Render([]table.Value{table.String("foo"), table.String("hello")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("bar"), table.String("world")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_color_fg_bg", buf.Bytes())
}

func TestGolden_TableColorNil(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithColor(ScopeBody, Columns(1), NewColor("red", "")),
		WithPlaceholder("N/A"),
		WithHeader([]string{"Key", "Value"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("a"), table.Any(nil)},
		{table.String("b"), table.String("")},
		{table.String("c"), table.String("ok")},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_color_nil", buf.Bytes())
}

func TestGolden_StreamColorNil(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithColor(ScopeBody, Columns(1), NewColor("red", "")),
		WithPlaceholder("N/A"),
		WithHeader([]string{"Key", "Value"}),
	)
	if err := s.Render([]table.Value{table.String("a"), table.Any(nil)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("b"), table.String("")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("c"), table.String("ok")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_color_nil", buf.Bytes())
}

func TestGolden_TableColorPlaceholder(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithColor(ScopeBody, Columns(1), ColorFgRed),
		WithPlaceholder("N/A"),
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("x"), table.String(""), table.Int(1)},
		{table.String("y"), table.String("b"), table.Int(2)},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_color_placeholder", buf.Bytes())
}

func TestGolden_StreamColorPlaceholder(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithColor(ScopeBody, Columns(1), ColorFgRed),
		WithPlaceholder("N/A"),
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := s.Render([]table.Value{table.String("x"), table.String(""), table.Int(1)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("y"), table.String("b"), table.Int(2)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_color_placeholder", buf.Bytes())
}

func TestGolden_TableColorScope(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"Name", "Value"}),
		WithFooter(func() [][]string {
			return [][]string{{"total", "3"}}
		}),
		WithColor(ScopeHeader, Columns(0), ColorFgRed),
		WithColor(ScopeBody, Columns(1), ColorFgBlue),
		WithColor(ScopeFooter, Columns(0, 1), ColorFgGreen),
		WithDecoration(ScopeHeader|ScopeFooter, Columns(1), DecorationBold),
	)
	if err := tb.Render([][]table.Value{
		{table.String("foo"), table.Int(1)},
		{table.String("bar"), table.Int(2)},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_color_scope", buf.Bytes())
}

func TestGolden_StreamColorScope(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"Name", "Value"}),
		WithFooter(func() [][]string {
			return [][]string{{"total", "3"}}
		}),
		WithColor(ScopeHeader, Columns(0), ColorFgRed),
		WithColor(ScopeBody, Columns(1), ColorFgBlue),
		WithColor(ScopeFooter, Columns(0, 1), ColorFgGreen),
		WithDecoration(ScopeHeader|ScopeFooter, Columns(1), DecorationBold),
	)
	for _, row := range [][]table.Value{{table.String("foo"), table.Int(1)}, {table.String("bar"), table.Int(2)}} {
		if err := s.Render(row); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_color_scope", buf.Bytes())
}

func TestGolden_TableColspan(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"Category", "Q1", "Q2", "Q3", "Q4"}),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(1, 2, 3, 4)),
	)
	if err := tb.Render([][]table.Value{
		{table.String("Sales"), table.Int(100), table.Int(100), table.Int(100), table.Int(200)},
		{table.String("Cost"), table.Int(50), table.Int(50), table.Int(50), table.Int(50)},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_colspan", buf.Bytes())
}

func TestGolden_StreamColspan(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"Category", "Q1", "Q2", "Q3", "Q4"}),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(1, 2, 3, 4)),
	)
	if err := s.Render([]table.Value{table.String("Sales"), table.Int(100), table.Int(100), table.Int(100), table.Int(200)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("Cost"), table.Int(50), table.Int(50), table.Int(50), table.Int(50)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_colspan", buf.Bytes())
}

func TestGolden_TableColspanCellAttr(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1, 2)),
		WithCellAttr(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), Attr{Class: "c", Style: "width:2em"}),
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := tb.Render([][]table.Value{{table.String("x"), table.String("x"), table.String("y")}, {table.String("p"), table.String("q"), table.String("q")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_colspan_cellattr", buf.Bytes())
}

func TestGolden_StreamColspanCellAttr(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1, 2)),
		WithCellAttr(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), Attr{Class: "c", Style: "width:2em"}),
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
	testutil.AssertGolden(t, "common_colspan_cellattr", buf.Bytes())
}

func TestGolden_TableColspanColor(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1)),
		WithColor(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), ColorFgRed),
	)
	if err := tb.Render([][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_colspan_color", buf.Bytes())
}

func TestGolden_StreamColspanColor(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1)),
		WithColor(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), ColorFgRed),
	)
	for _, r := range [][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_colspan_color", buf.Bytes())
}

func TestGolden_TableColspanDecoration(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1, 2)),
		WithDecoration(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), DecorationBold),
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("x"), table.String("x"), table.String("y")},
		{table.String("p"), table.String("q"), table.String("q")},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_colspan_decoration", buf.Bytes())
}

func TestGolden_StreamColspanDecoration(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1, 2)),
		WithDecoration(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), DecorationBold),
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
	testutil.AssertGolden(t, "common_colspan_decoration", buf.Bytes())
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

func TestGolden_TableColspanScope(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
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

func TestGolden_TableColspanTransformer(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1)),
		WithTransformer(Columns(1), func(v table.Value) (string, *Color, *Decoration) {
			if s, ok := v.AsAny().(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
		}),
	)
	if err := tb.Render([][]table.Value{{table.String("T"), table.String("raw")}, {table.String("p"), table.String("q")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_colspan_transformer", buf.Bytes())
}

func TestGolden_StreamColspanTransformer(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1)),
		WithTransformer(Columns(1), func(v table.Value) (string, *Color, *Decoration) {
			if s, ok := v.AsAny().(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
		}),
	)
	for _, r := range [][]table.Value{{table.String("T"), table.String("raw")}, {table.String("p"), table.String("q")}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_colspan_transformer", buf.Bytes())
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

func TestGolden_TableDecoAlign(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithDecoration(ScopeBody, Columns(1), DecorationBold),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), AlignCenter),
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("x"), table.String("y"), table.Int(1)},
		{table.String("p"), table.String("qq"), table.Int(22)},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_deco_align", buf.Bytes())
}

func TestGolden_StreamDecoAlign(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithDecoration(ScopeBody, Columns(1), DecorationBold),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), AlignCenter),
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := s.Render([]table.Value{table.String("x"), table.String("y"), table.Int(1)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("p"), table.String("qq"), table.Int(22)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_deco_align", buf.Bytes())
}

func TestGolden_TableDecoEscape(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithDecoration(ScopeBody, Columns(1), DecorationCode),
		WithHeader([]string{"Key", "Value"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("amp"), table.String("a&b")},
		{table.String("newline"), table.String("line1\nline2")},
		{table.String("tag"), table.String("<div>")},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_deco_escape", buf.Bytes())
}

func TestGolden_StreamDecoEscape(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithDecoration(ScopeBody, Columns(1), DecorationCode),
		WithHeader([]string{"Key", "Value"}),
	)
	if err := s.Render([]table.Value{table.String("amp"), table.String("a&b")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("newline"), table.String("line1\nline2")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("tag"), table.String("<div>")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_deco_escape", buf.Bytes())
}

func TestGolden_TableDecoMultiCol(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithDecoration(ScopeBody, Columns(0), DecorationCode),
		WithDecoration(ScopeBody, Columns(2), DecorationBold),
		WithHeader([]string{"ID", "Name", "Status"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("id-1"), table.String("alice"), table.String("active")},
		{table.String("id-2"), table.String("bob"), table.String("inactive")},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_deco_multi_col", buf.Bytes())
}

func TestGolden_StreamDecoMultiCol(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithDecoration(ScopeBody, Columns(0), DecorationCode),
		WithDecoration(ScopeBody, Columns(2), DecorationBold),
		WithHeader([]string{"ID", "Name", "Status"}),
	)
	if err := s.Render([]table.Value{table.String("id-1"), table.String("alice"), table.String("active")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("id-2"), table.String("bob"), table.String("inactive")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_deco_multi_col", buf.Bytes())
}

func TestGolden_TableDecoNil(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithDecoration(ScopeBody, Columns(1), DecorationCode),
		WithPlaceholder("N/A"),
		WithHeader([]string{"Key", "Value"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("a"), table.Any(nil)},
		{table.String("b"), table.String("")},
		{table.String("c"), table.String("ok")},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_deco_nil", buf.Bytes())
}

func TestGolden_StreamDecoNil(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithDecoration(ScopeBody, Columns(1), DecorationCode),
		WithPlaceholder("N/A"),
		WithHeader([]string{"Key", "Value"}),
	)
	if err := s.Render([]table.Value{table.String("a"), table.Any(nil)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("b"), table.String("")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("c"), table.String("ok")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_deco_nil", buf.Bytes())
}

func TestGolden_TableDecoShortRow(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithDecoration(ScopeBody, Columns(1), DecorationCode),
		WithPlaceholder("-"),
		WithHeader([]string{"Key", "Value"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("full"), table.String("ok")},
		{table.String("short")},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_deco_short_row", buf.Bytes())
}

func TestGolden_StreamDecoShortRow(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithDecoration(ScopeBody, Columns(1), DecorationCode),
		WithPlaceholder("-"),
		WithHeader([]string{"Key", "Value"}),
	)
	if err := s.Render([]table.Value{table.String("full"), table.String("ok")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("short")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_deco_short_row", buf.Bytes())
}

func TestGolden_TableDecorationCellAttr(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithDecoration(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), DecorationBold),
		WithCellAttr(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), Attr{Class: "c"}),
	)
	if err := tb.Render([][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_decoration_cellattr", buf.Bytes())
}

func TestGolden_StreamDecorationCellAttr(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
		WithDecoration(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), DecorationBold),
		WithCellAttr(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), Attr{Class: "c"}),
	)
	for _, r := range [][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_decoration_cellattr", buf.Bytes())
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

func TestGolden_TableEmptyVsNil(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
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
	testutil.AssertGolden(t, "common_empty_vs_nil", buf.Bytes())
}

func TestGolden_StreamEmptyVsNil(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
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
	testutil.AssertGolden(t, "common_empty_vs_nil", buf.Bytes())
}

func TestGolden_TableEscapeCRLF(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"Key", "Value"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("crlf"), table.String("a\r\nb")},
		{table.String("cr"), table.String("a\rb")},
		{table.String("lf"), table.String("a\nb")},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_escape_crlf", buf.Bytes())
}

func TestGolden_StreamEscapeCRLF(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"Key", "Value"}),
	)
	if err := s.Render([]table.Value{table.String("crlf"), table.String("a\r\nb")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("cr"), table.String("a\rb")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("lf"), table.String("a\nb")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_escape_crlf", buf.Bytes())
}

func TestGolden_TableEscapeSpace(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"Key", "Value"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("interior"), table.String("a b")},
		{table.String("run"), table.String("a  b")},
		{table.String("lead"), table.String("  a")},
		{table.String("trail"), table.String("a ")},
		{table.String("json"), table.String("{\n  \"k\": 1\n}")},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_escape_space", buf.Bytes())
}

func TestGolden_StreamEscapeSpace(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"Key", "Value"}),
	)
	if err := s.Render([]table.Value{table.String("interior"), table.String("a b")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("run"), table.String("a  b")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("lead"), table.String("  a")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("trail"), table.String("a ")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("json"), table.String("{\n  \"k\": 1\n}")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_escape_space", buf.Bytes())
}

func TestGolden_TableFooter(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"Name", "Score"}),
		WithFooter(func() [][]string {
			return [][]string{{"Total", "300"}}
		}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("alice"), table.Int(100)},
		{table.String("bob"), table.Int(200)},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_footer", buf.Bytes())
}

func TestGolden_StreamFooter(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"Name", "Score"}),
		WithFooter(func() [][]string {
			return [][]string{{"Total", "300"}}
		}),
	)
	if err := s.Render([]table.Value{table.String("alice"), table.Int(100)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("bob"), table.Int(200)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_footer", buf.Bytes())
}

func TestGolden_TableFooterCaption(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithFooter(func() [][]string {
			return [][]string{{"t", "u"}}
		}),
		WithCaption("cap", CaptionBottom),
	)
	if err := tb.Render([][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_footer_caption", buf.Bytes())
}

func TestGolden_StreamFooterCaption(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
		WithFooter(func() [][]string {
			return [][]string{{"t", "u"}}
		}),
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
	testutil.AssertGolden(t, "common_footer_caption", buf.Bytes())
}

func TestGolden_TableFooterEmptyBody(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
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
		WithFooter(func() [][]string {
			return [][]string{{"Total", "300"}}
		}),
	)
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_footer_no_header", buf.Bytes())
}

func TestGolden_TableFooterPlaceholder(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
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

func TestGolden_TableFooterTransformer(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithFooter(func() [][]string {
			return [][]string{{"t", "raw"}}
		}),
		WithTransformer(Columns(1), func(v table.Value) (string, *Color, *Decoration) {
			if s, ok := v.AsAny().(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
		}),
	)
	if err := tb.Render([][]table.Value{{table.String("x"), table.String("raw")}, {table.String("p"), table.String("q")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_footer_transformer", buf.Bytes())
}

func TestGolden_StreamFooterTransformer(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
		WithFooter(func() [][]string {
			return [][]string{{"t", "raw"}}
		}),
		WithTransformer(Columns(1), func(v table.Value) (string, *Color, *Decoration) {
			if s, ok := v.AsAny().(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
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
	testutil.AssertGolden(t, "common_footer_transformer", buf.Bytes())
}

func TestGolden_TableHeaderOnly(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
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
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_header_only", buf.Bytes())
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

func TestGolden_TableItalic(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithDecoration(ScopeBody, Columns(1), DecorationItalic),
		WithHeader([]string{"Key", "Note"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("a"), table.String("important")},
		{table.String("b"), table.String("optional")},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_italic", buf.Bytes())
}

func TestGolden_StreamItalic(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithDecoration(ScopeBody, Columns(1), DecorationItalic),
		WithHeader([]string{"Key", "Note"}),
	)
	if err := s.Render([]table.Value{table.String("a"), table.String("important")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("b"), table.String("optional")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_italic", buf.Bytes())
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
	testutil.AssertGolden(t, "common_long_value", buf.Bytes())
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
	testutil.AssertGolden(t, "common_long_value", buf.Bytes())
}

func TestGolden_TableNilInNumeric(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithPlaceholder("N/A"),
		WithHeader([]string{"N", "V"}),
	)
	if err := tb.Render([][]table.Value{
		{table.Int(1), table.Int(100)},
		{table.Int(2), table.Any(nil)},
		{table.Int(3), table.Any(-5)},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_nil_in_numeric", buf.Bytes())
}

func TestGolden_StreamNilInNumeric(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithPlaceholder("N/A"),
		WithHeader([]string{"N", "V"}),
	)
	if err := s.Render([]table.Value{table.Int(1), table.Int(100)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.Int(2), table.Any(nil)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.Int(3), table.Any(-5)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_nil_in_numeric", buf.Bytes())
}

func TestGolden_TableNoHeader(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf)
	if err := tb.Render([][]table.Value{
		{table.String("a"), table.Int(1)},
		{table.String("b"), table.Int(2)},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_no_header", buf.Bytes())
}

func TestGolden_StreamNoHeader(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf)
	if err := s.Render([]table.Value{table.String("a"), table.Int(1)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("b"), table.Int(2)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_no_header", buf.Bytes())
}

func TestGolden_TablePlaceholder(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
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

func TestGolden_TablePlaceholderAlign(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithPlaceholder("-"),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), AlignRight),
	)
	if err := tb.Render([][]table.Value{{table.String("x"), table.Any(nil)}, {table.Any(nil), table.String("q")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_placeholder_align", buf.Bytes())
}

func TestGolden_StreamPlaceholderAlign(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
		WithPlaceholder("-"),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), AlignRight),
	)
	for _, r := range [][]table.Value{{table.String("x"), table.Any(nil)}, {table.Any(nil), table.String("q")}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_placeholder_align", buf.Bytes())
}

func TestGolden_TablePlaceholderCellAttr(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithPlaceholder("-"),
		WithCellAttr(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), Attr{Class: "c"}),
	)
	if err := tb.Render([][]table.Value{{table.String("x"), table.Any(nil)}, {table.Any(nil), table.String("q")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_placeholder_cellattr", buf.Bytes())
}

func TestGolden_StreamPlaceholderCellAttr(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
		WithPlaceholder("-"),
		WithCellAttr(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), Attr{Class: "c"}),
	)
	for _, r := range [][]table.Value{{table.String("x"), table.Any(nil)}, {table.Any(nil), table.String("q")}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_placeholder_cellattr", buf.Bytes())
}

func TestGolden_TablePlaceholderColspan(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithPlaceholder("-"),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1, 2)),
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := tb.Render([][]table.Value{{table.String("x"), table.Any(nil), table.Any(nil)}, {table.Any(nil), table.Any(nil), table.String("q")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_placeholder_colspan", buf.Bytes())
}

func TestGolden_StreamPlaceholderColspan(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithPlaceholder("-"),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1, 2)),
		WithHeader([]string{"A", "B", "C"}),
	)
	for _, r := range [][]table.Value{{table.String("x"), table.Any(nil), table.Any(nil)}, {table.Any(nil), table.Any(nil), table.String("q")}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_placeholder_colspan", buf.Bytes())
}

func TestGolden_TablePlaceholderTransformer(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithPlaceholder("-"),
		WithTransformer(Columns(1), func(v table.Value) (string, *Color, *Decoration) {
			if s, ok := v.AsAny().(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
		}),
	)
	if err := tb.Render([][]table.Value{{table.String("x"), table.Any(nil)}, {table.String("p"), table.String("raw")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_placeholder_transformer", buf.Bytes())
}

func TestGolden_StreamPlaceholderTransformer(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
		WithPlaceholder("-"),
		WithTransformer(Columns(1), func(v table.Value) (string, *Color, *Decoration) {
			if s, ok := v.AsAny().(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
		}),
	)
	for _, r := range [][]table.Value{{table.String("x"), table.Any(nil)}, {table.String("p"), table.String("raw")}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_placeholder_transformer", buf.Bytes())
}

func TestGolden_TablePointer(t *testing.T) {
	s := testutil.Stringer{Value: "y"}
	str := "alive"
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"a", "b"}),
	)
	if err := tb.Render([][]table.Value{
		{table.Any(&str), table.Any(&s)},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_pointer", buf.Bytes())
}

func TestGolden_StreamPointer(t *testing.T) {
	stringer := testutil.Stringer{Value: "y"}
	str := "alive"
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"a", "b"}),
	)
	if err := s.Render([]table.Value{table.Any(&str), table.Any(&stringer)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_pointer", buf.Bytes())
}

func TestGolden_TablePreformatted(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithColor(ScopeBody, Columns(1), ColorFgRed),
		WithDecoration(ScopeBody, Columns(1), DecorationPreformatted),
		WithHeader([]string{"Name", "Content"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("alpha"), table.String("first\nsecond")},
		{table.String("beta"), table.String("  <value>")},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_preformatted", buf.Bytes())
}

func TestGolden_StreamPreformatted(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithColor(ScopeBody, Columns(1), ColorFgRed),
		WithDecoration(ScopeBody, Columns(1), DecorationPreformatted),
		WithHeader([]string{"Name", "Content"}),
	)
	if err := s.Render([]table.Value{table.String("alpha"), table.String("first\nsecond")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("beta"), table.String("  <value>")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_preformatted", buf.Bytes())
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

func TestGolden_TableSingleCell(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
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

func TestGolden_TableSlice(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"Type", "Value"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("[]int"), table.Any([]int{1, 2, 3})},
		{table.String("[]string"), table.Any([]string{"a", "b"})},
		{table.String("[]bool"), table.Any([]bool{true, false})},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_slice", buf.Bytes())
}

func TestGolden_StreamSlice(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"Type", "Value"}),
	)
	if err := s.Render([]table.Value{table.String("[]int"), table.Any([]int{1, 2, 3})}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("[]string"), table.Any([]string{"a", "b"})}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("[]bool"), table.Any([]bool{true, false})}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_slice", buf.Bytes())
}

func TestGolden_TableSpecialChars(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"Char", "Value"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("less-than"), table.String("<")},
		{table.String("greater-than"), table.String(">")},
		{table.String("double-quote"), table.String("\"")},
		{table.String("single-quote"), table.String("'")},
		{table.String("ampersand"), table.String("&")},
		{table.String("space"), table.String(" ")},
		{table.String("asterisk"), table.String("*")},
		{table.String("backslash"), table.String("\\")},
		{table.String("underscore"), table.String("_")},
		{table.String("pipe"), table.String("|")},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_special_chars", buf.Bytes())
}

func TestGolden_StreamSpecialChars(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"Char", "Value"}),
	)
	if err := s.Render([]table.Value{table.String("less-than"), table.String("<")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("greater-than"), table.String(">")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("double-quote"), table.String("\"")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("single-quote"), table.String("'")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("ampersand"), table.String("&")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("space"), table.String(" ")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("asterisk"), table.String("*")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("backslash"), table.String("\\")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("underscore"), table.String("_")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("pipe"), table.String("|")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_special_chars", buf.Bytes())
}

func TestGolden_TableStrikethrough(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithDecoration(ScopeBody, Columns(1), DecorationStrikethrough),
		WithHeader([]string{"Feature", "Status"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("login"), table.String("deprecated")},
		{table.String("signup"), table.String("active")},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_strikethrough", buf.Bytes())
}

func TestGolden_StreamStrikethrough(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithDecoration(ScopeBody, Columns(1), DecorationStrikethrough),
		WithHeader([]string{"Feature", "Status"}),
	)
	if err := s.Render([]table.Value{table.String("login"), table.String("deprecated")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("signup"), table.String("active")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_strikethrough", buf.Bytes())
}

func TestGolden_TableStringerError(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"Stringer", "Error"}),
	)
	if err := tb.Render([][]table.Value{
		{table.Any(testutil.Stringer{Value: "x"}), table.Any(testutil.Error{Value: "boom"})},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_stringer_error", buf.Bytes())
}

func TestGolden_StreamStringerError(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"Stringer", "Error"}),
	)
	if err := s.Render([]table.Value{table.Any(testutil.Stringer{Value: "x"}), table.Any(testutil.Error{Value: "boom"})}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_stringer_error", buf.Bytes())
}

func TestGolden_TableStyle(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithTableAttr(TableAttr{
			Table:   Attr{Style: "border-collapse:collapse"},
			Caption: Attr{Style: "font-style:italic"},
			Header:  SectionAttr{Row: Attr{Style: "background:#eee"}, Cell: Attr{Style: "font-weight:bold"}},
			Body:    SectionAttr{Cell: Attr{Style: "padding:4px"}},
		}),
		WithCaption("Styled", CaptionDefault),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), AlignRight),
		WithHeader([]string{"Name", "Score"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("alice"), table.Int(100)},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_style", buf.Bytes())
}

func TestGolden_StreamStyle(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithTableAttr(TableAttr{
			Table:   Attr{Style: "border-collapse:collapse"},
			Caption: Attr{Style: "font-style:italic"},
			Header:  SectionAttr{Row: Attr{Style: "background:#eee"}, Cell: Attr{Style: "font-weight:bold"}},
			Body:    SectionAttr{Cell: Attr{Style: "padding:4px"}},
		}),
		WithCaption("Styled", CaptionDefault),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), AlignRight),
		WithHeader([]string{"Name", "Score"}),
	)
	if err := s.Render([]table.Value{table.String("alice"), table.Int(100)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_style", buf.Bytes())
}

func TestGolden_TableStyleJoined(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithTableAttr(TableAttr{
			Body: SectionAttr{Cell: Attr{Style: "padding:4px"}},
		}),
		WithCellAttr(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), Attr{Style: "color:#333"}),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), AlignLeft),
		WithHeader([]string{"Name", "Score"}),
		WithFooter(func() [][]string {
			return [][]string{{"total"}}
		}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("alice"), table.Int(100)},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_style_joined", buf.Bytes())
}

func TestGolden_StreamStyleJoined(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithTableAttr(TableAttr{
			Body: SectionAttr{Cell: Attr{Style: "padding:4px"}},
		}),
		WithCellAttr(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), Attr{Style: "color:#333"}),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), AlignLeft),
		WithHeader([]string{"Name", "Score"}),
		WithFooter(func() [][]string {
			return [][]string{{"total"}}
		}),
	)
	if err := s.Render([]table.Value{table.String("alice"), table.Int(100)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_style_joined", buf.Bytes())
}

func TestGolden_TableTableAttrColor(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithTableAttr(TableAttr{
			Table: Attr{Class: "t"},
			Body:  SectionAttr{Cell: Attr{Class: "b"}},
		}),
		WithColor(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), ColorFgRed),
	)
	if err := tb.Render([][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_tableattr_color", buf.Bytes())
}

func TestGolden_StreamTableAttrColor(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
		WithTableAttr(TableAttr{
			Table: Attr{Class: "t"},
			Body:  SectionAttr{Cell: Attr{Class: "b"}},
		}),
		WithColor(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), ColorFgRed),
	)
	for _, r := range [][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_tableattr_color", buf.Bytes())
}

func TestGolden_TableTableAttrColspan(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithTableAttr(TableAttr{
			Table: Attr{Class: "t"},
			Body:  SectionAttr{Cell: Attr{Class: "b"}},
		}),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1)),
	)
	if err := tb.Render([][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_tableattr_colspan", buf.Bytes())
}

func TestGolden_StreamTableAttrColspan(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
		WithTableAttr(TableAttr{
			Table: Attr{Class: "t"},
			Body:  SectionAttr{Cell: Attr{Class: "b"}},
		}),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1)),
	)
	for _, r := range [][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_tableattr_colspan", buf.Bytes())
}

func TestGolden_TableTableAttrDecoration(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithTableAttr(TableAttr{
			Table: Attr{Class: "t"},
			Body:  SectionAttr{Cell: Attr{Class: "b"}},
		}),
		WithDecoration(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), DecorationBold),
	)
	if err := tb.Render([][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_tableattr_decoration", buf.Bytes())
}

func TestGolden_StreamTableAttrDecoration(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
		WithTableAttr(TableAttr{
			Table: Attr{Class: "t"},
			Body:  SectionAttr{Cell: Attr{Class: "b"}},
		}),
		WithDecoration(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), DecorationBold),
	)
	for _, r := range [][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_tableattr_decoration", buf.Bytes())
}

func TestGolden_TableTableAttrPlaceholder(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithTableAttr(TableAttr{
			Table: Attr{Class: "t"},
			Body:  SectionAttr{Cell: Attr{Class: "b"}},
		}),
		WithPlaceholder("-"),
	)
	if err := tb.Render([][]table.Value{{table.String("x"), table.Any(nil)}, {table.Any(nil), table.String("q")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_tableattr_placeholder", buf.Bytes())
}

func TestGolden_StreamTableAttrPlaceholder(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
		WithTableAttr(TableAttr{
			Table: Attr{Class: "t"},
			Body:  SectionAttr{Cell: Attr{Class: "b"}},
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
	testutil.AssertGolden(t, "common_tableattr_placeholder", buf.Bytes())
}

func TestGolden_TableTableAttrTransformer(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithTableAttr(TableAttr{
			Table: Attr{Class: "t"},
			Body:  SectionAttr{Cell: Attr{Class: "b"}},
		}),
		WithTransformer(Columns(1), func(v table.Value) (string, *Color, *Decoration) {
			if s, ok := v.AsAny().(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
		}),
	)
	if err := tb.Render([][]table.Value{{table.String("x"), table.String("raw")}, {table.String("p"), table.String("q")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_tableattr_transformer", buf.Bytes())
}

func TestGolden_StreamTableAttrTransformer(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
		WithTableAttr(TableAttr{
			Table: Attr{Class: "t"},
			Body:  SectionAttr{Cell: Attr{Class: "b"}},
		}),
		WithTransformer(Columns(1), func(v table.Value) (string, *Color, *Decoration) {
			if s, ok := v.AsAny().(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
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
	testutil.AssertGolden(t, "common_tableattr_transformer", buf.Bytes())
}

func TestGolden_TableTransformer(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithTransformer(Columns(1), func(v table.Value) (string, *Color, *Decoration) {
			n, ok := v.AsAny().(int)
			if !ok {
				return "", nil, nil
			}
			if n >= 100 {
				return "high", ColorFgRed, DecorationBold
			}
			return "", nil, nil
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
		WithTransformer(Columns(1), func(v table.Value) (string, *Color, *Decoration) {
			n, ok := v.AsAny().(int)
			if !ok {
				return "", nil, nil
			}
			if n >= 100 {
				return "high", ColorFgRed, DecorationBold
			}
			return "", nil, nil
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
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithColor(ScopeBody, Columns(1), ColorFgBlue),
		WithDecoration(ScopeBody, Columns(1), DecorationBold),
		WithTransformer(Columns(1), func(v table.Value) (string, *Color, *Decoration) {
			if s, ok := v.AsAny().(string); ok && s == "warn" {
				return "", ColorFgYellow, DecorationItalic
			}
			return "", nil, nil
		}),
		WithHeader([]string{"Level", "Message"}),
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
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithColor(ScopeBody, Columns(1), ColorFgBlue),
		WithDecoration(ScopeBody, Columns(1), DecorationBold),
		WithTransformer(Columns(1), func(v table.Value) (string, *Color, *Decoration) {
			if s, ok := v.AsAny().(string); ok && s == "warn" {
				return "", ColorFgYellow, DecorationItalic
			}
			return "", nil, nil
		}),
		WithHeader([]string{"Level", "Message"}),
	)
	if err := s.Render([]table.Value{table.String("1"), table.String("ok")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("2"), table.String("warn")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_transformer_column_override", buf.Bytes())
}

func TestGolden_TableTypeFloat(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
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
		WithHeader([]string{
			"i", "i8", "i16", "i32", "i64",
			"u", "u8", "u16", "u32", "u64",
		}),
	)
	if err := s.Render([]table.Value{table.Int(int(-1)), table.Int8(int8(-2)), table.Int16(int16(-3)), table.Int32(int32(-4)), table.Int64(int64(-5)),
		table.Uint(uint(1)), table.Uint8(uint8(2)), table.Uint16(uint16(3)), table.Uint32(uint32(4)), table.Uint64(uint64(5))}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_type_integer", buf.Bytes())
}

func TestGolden_TableTypedNil(t *testing.T) {
	var nilStringer *testutil.PtrStringer
	var nilError *testutil.PtrError
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithPlaceholder("<nil>"),
		WithHeader([]string{"Stringer", "Error"}),
	)
	if err := tb.Render([][]table.Value{
		{table.Any(nilStringer), table.Any(nilError)},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_typed_nil", buf.Bytes())
}

func TestGolden_StreamTypedNil(t *testing.T) {
	var nilStringer *testutil.PtrStringer
	var nilError *testutil.PtrError
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithPlaceholder("<nil>"),
		WithHeader([]string{"Stringer", "Error"}),
	)
	if err := s.Render([]table.Value{table.Any(nilStringer), table.Any(nilError)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_typed_nil", buf.Bytes())
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

func TestGolden_TableValueEqualsPlaceholderColor(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A"}),
		WithPlaceholder("N/A"),
		WithColor(ScopeBody, Columns(0), ColorFgRed),
	)
	if err := tb.Render([][]table.Value{
		{table.String("N/A")},
		{table.String("")},
		{table.String("x")},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_value_equals_placeholder_color", buf.Bytes())
}

func TestGolden_StreamValueEqualsPlaceholderColor(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A"}),
		WithPlaceholder("N/A"),
		WithColor(ScopeBody, Columns(0), ColorFgRed),
	)
	if err := s.Render([]table.Value{table.String("N/A")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("x")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_value_equals_placeholder_color", buf.Bytes())
}

func TestGolden_TableWideNumber(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"Big", "Neg", "Sci"}),
	)
	if err := tb.Render([][]table.Value{
		{table.Int64(int64(9223372036854775807)), table.Any(-9223372036854775808), table.Float64(1.7976931348623157e+308)},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_wide_number", buf.Bytes())
}

func TestGolden_StreamWideNumber(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"Big", "Neg", "Sci"}),
	)
	if err := s.Render([]table.Value{table.Int64(int64(9223372036854775807)), table.Any(-9223372036854775808), table.Float64(1.7976931348623157e+308)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_wide_number", buf.Bytes())
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

func TestGolden_StreamCaption(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A"}),
		WithCaption("Stream Caption", CaptionDefault),
	)
	if err := s.Render([]table.Value{table.String("x")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_caption", buf.Bytes())
}

func TestGolden_StreamCaptionRowspan(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
		WithCaption("cap", CaptionBottom),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
	)
	for _, r := range [][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_caption_rowspan", buf.Bytes())
}

func TestGolden_StreamColorRowspan(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithColor(ScopeBody, Columns(0), ColorFgBlue),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := s.Render([]table.Value{table.String("g"), table.String("x"), table.Int(1)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("g"), table.String("y"), table.Int(2)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("h"), table.String("z"), table.Int(3)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_color_rowspan", buf.Bytes())
}

func TestGolden_StreamDecoRowspan(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithDecoration(ScopeBody, Columns(0), DecorationBold),
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
	testutil.AssertGolden(t, "stream_deco_rowspan", buf.Bytes())
}

func TestGolden_StreamDecoSlice(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithDecoration(ScopeBody, Columns(1), DecorationCode),
		WithHeader([]string{"Type", "Values"}),
	)
	if err := s.Render([]table.Value{table.String("ints"), table.Any([]int{1, 2, 3})}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("mixed"), table.Any([]string{"a", "", "c"})}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_deco_slice", buf.Bytes())
}

func TestGolden_StreamEscape(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"Key", "Value"}),
	)
	if err := s.Render([]table.Value{table.String("amp"), table.String("a&b")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("newline"), table.String("a\nb")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("multiple&"), table.String(`<one> & "two"`)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_escape", buf.Bytes())
}

func TestGolden_StreamFooterRowspan(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithFooter(func() [][]string {
			return [][]string{
				{"t", "x"},
				{"t", "y"},
				{"u", "z"},
			}
		}),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
		WithHeader([]string{"A", "B"}),
	)
	for _, r := range [][]table.Value{{table.String("g"), table.String("x")}, {table.String("g"), table.String("y")}, {table.String("h"), table.String("z")}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_footer_rowspan", buf.Bytes())
}

func TestGolden_StreamRowspan(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
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
	testutil.AssertGolden(t, "stream_rowspan", buf.Bytes())
}

func TestGolden_StreamRowspanAlign(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(0), AlignRight),
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := s.Render([]table.Value{table.String("g"), table.String("x"), table.Int(1)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("g"), table.String("y"), table.Int(2)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("h"), table.String("z"), table.Int(3)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_rowspan_align", buf.Bytes())
}

func TestGolden_StreamRowspanCellAttr(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
		WithCellAttr(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), Attr{Class: "c"}),
	)
	for _, r := range [][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_rowspan_cellattr", buf.Bytes())
}

func TestGolden_StreamRowspanColspanContinuation(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"c1", "c2", "c3"}),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(1, 2)),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(1, 2)),
	)
	if err := s.Render([]table.Value{table.String("a"), table.String("x"), table.String("x")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("a"), table.String("x"), table.String("x")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_rowspan_colspan_continuation", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_rowspan_colspan_edge", buf.Bytes())
}

func TestGolden_StreamRowspanEscape(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
		WithHeader([]string{"A", "B"}),
	)
	if err := s.Render([]table.Value{table.String("<x>"), table.Int(1)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("<x>"), table.Int(2)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("&y&"), table.Int(3)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_rowspan_escape", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_rowspan_missing_kinds", buf.Bytes())
}

func TestGolden_StreamRowspanPlaceholder(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
		WithPlaceholder("N/A"),
		WithHeader([]string{"Group", "Value"}),
	)
	if err := s.Render([]table.Value{table.String("A"), table.Any(nil)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("A"), table.String("y")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("B"), table.Any(nil)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_rowspan_placeholder", buf.Bytes())
}

func TestGolden_StreamRowspanTransformer(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(1)),
		WithTransformer(Columns(1), func(v table.Value) (string, *Color, *Decoration) {
			if s, ok := v.AsAny().(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
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
	testutil.AssertGolden(t, "stream_rowspan_transformer", buf.Bytes())
}

func TestGolden_StreamSingleRow(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
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

func TestGolden_StreamSpanLimit(t *testing.T) {
	const n = param.SpanLimit + 2
	header := make([]string, n)
	indexes := make([]int, n)
	for i := range n {
		header[i] = "c"
		indexes[i] = i
	}
	cols := Columns(indexes...)
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader(header),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, cols),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, cols),
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
	testutil.AssertGolden(t, "stream_span_limit", buf.Bytes())
}

func TestGolden_StreamTableAttrRowspan(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
		WithTableAttr(TableAttr{
			Table: Attr{Class: "t"},
			Body:  SectionAttr{Cell: Attr{Class: "b"}},
		}),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
	)
	for _, r := range [][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_tableattr_rowspan", buf.Bytes())
}

func TestGolden_StreamTypeMixed(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"Label", "Value"}),
	)
	if err := s.Render([]table.Value{table.String("text"), table.String("hello")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("number"), table.Int(42)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("nil"), table.Any(nil)}); err != nil {
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

func TestGolden_TableBandColspan(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"Pair", "Pair", "Solo"}),
		WithFooter(func() [][]string {
			return [][]string{{"Sum", "Sum", "Rest"}}
		}),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1, 2)),
	)
	if err := tb.Render([][]table.Value{
		{table.String("x"), table.String("x"), table.String("z")},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_band_colspan", buf.Bytes())
}

func TestGolden_TableCaption(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"Name", "Value"}),
		WithCaption("Test Caption", CaptionDefault),
	)
	if err := tb.Render([][]table.Value{
		{table.String("foo"), table.Int(1)},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_caption", buf.Bytes())
}

func TestGolden_TableCaptionRowspan(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithCaption("cap", CaptionBottom),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
	)
	if err := tb.Render([][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_caption_rowspan", buf.Bytes())
}

func TestGolden_TableColorRowspan(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithColor(ScopeBody, Columns(0), ColorFgBlue),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("g"), table.String("x"), table.Int(1)},
		{table.String("g"), table.String("y"), table.Int(2)},
		{table.String("h"), table.String("z"), table.Int(3)},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_color_rowspan", buf.Bytes())
}

func TestGolden_TableDecoRowspan(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithDecoration(ScopeBody, Columns(0), DecorationBold),
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
	testutil.AssertGolden(t, "table_deco_rowspan", buf.Bytes())
}

func TestGolden_TableDecoSlice(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithDecoration(ScopeBody, Columns(1), DecorationCode),
		WithHeader([]string{"Type", "Values"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("ints"), table.Any([]int{1, 2, 3})},
		{table.String("strings"), table.Any([]string{"a", "b"})},
		{table.String("empty-elem"), table.Any([]string{"x", "", "z"})},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_deco_slice", buf.Bytes())
}

func TestGolden_TableEscape(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"Key", "Value"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("ampersand"), table.String("a&b")},
		{table.String("less-than"), table.String("a<b")},
		{table.String("greater-than"), table.String("a>b")},
		{table.String("double-quote"), table.String(`a"b`)},
		{table.String("newline"), table.String("a\nb")},
		{table.String("combined"), table.String("a<b&c\nd>e")},
		{table.String("multiple&"), table.String(`<one> & "two"`)},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_escape", buf.Bytes())
}

func TestGolden_TableFooterRowspan(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithFooter(func() [][]string {
			return [][]string{
				{"t", "x"},
				{"t", "y"},
				{"u", "z"},
			}
		}),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]table.Value{{table.String("g"), table.String("x")}, {table.String("g"), table.String("y")}, {table.String("h"), table.String("z")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_footer_rowspan", buf.Bytes())
}

func TestGolden_TableNoHeaderRagged(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf)
	if err := tb.Render([][]table.Value{
		{table.String("a"), table.Int(1)},
		{table.String("b")},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_no_header_ragged", buf.Bytes())
}

func TestGolden_TableRowspan(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
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
	testutil.AssertGolden(t, "table_rowspan", buf.Bytes())
}

func TestGolden_TableRowspanAlign(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(0), AlignRight),
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("g"), table.String("x"), table.Int(1)},
		{table.String("g"), table.String("y"), table.Int(2)},
		{table.String("h"), table.String("z"), table.Int(3)},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_rowspan_align", buf.Bytes())
}

func TestGolden_TableRowspanCellAttr(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
		WithCellAttr(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), Attr{Class: "c"}),
	)
	if err := tb.Render([][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_rowspan_cellattr", buf.Bytes())
}

func TestGolden_TableRowspanColspan(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B", "C"}),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1)),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1, 2)),
	)
	if err := tb.Render([][]table.Value{
		{table.String("x"), table.String("x"), table.String("x")},
		{table.String("x"), table.String("x"), table.String("y")},
		{table.String("z"), table.String("z"), table.String("y")},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_rowspan_colspan", buf.Bytes())
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
	testutil.AssertGolden(t, "table_rowspan_colspan_edge", buf.Bytes())
}

func TestGolden_TableRowspanEscape(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("<x>"), table.Int(1)},
		{table.String("<x>"), table.Int(2)},
		{table.String("&y&"), table.Int(3)},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_rowspan_escape", buf.Bytes())
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
	testutil.AssertGolden(t, "table_rowspan_missing_kinds", buf.Bytes())
}

func TestGolden_TableRowspanPlaceholder(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
		WithPlaceholder("N/A"),
		WithHeader([]string{"Group", "Value"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("A"), table.Any(nil)},
		{table.String("A"), table.String("y")},
		{table.String("B"), table.Any(nil)},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_rowspan_placeholder", buf.Bytes())
}

func TestGolden_TableRowspanTransformer(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(1)),
		WithTransformer(Columns(1), func(v table.Value) (string, *Color, *Decoration) {
			if s, ok := v.AsAny().(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
		}),
	)
	if err := tb.Render([][]table.Value{{table.String("x"), table.String("raw")}, {table.String("p"), table.String("raw")}, {table.String("q"), table.String("z")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_rowspan_transformer", buf.Bytes())
}

func TestGolden_TableSpanLimit(t *testing.T) {
	const n = param.SpanLimit + 2
	header := make([]string, n)
	indexes := make([]int, n)
	for i := range n {
		header[i] = "c"
		indexes[i] = i
	}
	cols := Columns(indexes...)
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader(header),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, cols),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, cols),
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
	testutil.AssertGolden(t, "table_span_limit", buf.Bytes())
}

func TestGolden_TableTableAttrRowspan(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithTableAttr(TableAttr{
			Table: Attr{Class: "t"},
			Body:  SectionAttr{Cell: Attr{Class: "b"}},
		}),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
	)
	if err := tb.Render([][]table.Value{{table.String("s"), table.String("s")}, {table.String("s"), table.String("s")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_tableattr_rowspan", buf.Bytes())
}

func TestGolden_TableTypeMixed(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
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
