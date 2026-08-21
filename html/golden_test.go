package html

import (
	"bytes"
	"strings"
	"testing"

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
	if err := tb.Render([][]any{{"s", "s"}, {"s", "s"}}); err != nil {
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
	for _, r := range [][]any{{"s", "s"}, {"s", "s"}} {
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
	if err := tb.Render([][]any{
		{"a", 100, "x"},
		{"b", 200, "y"},
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
	if err := s.Render([]any{"a", 100, "x"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"b", 200, "y"}); err != nil {
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

func TestGolden_TableAlignTransformer(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), AlignRight),
		WithTransformer(Columns(1), func(v any) (string, *Color, *Decoration) {
			if s, ok := v.(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
		}),
	)
	if err := tb.Render([][]any{{"x", "raw"}, {"p", "q"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_align_transformer", buf.Bytes())
}

func TestGolden_StreamAlignTransformer(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), AlignRight),
		WithTransformer(Columns(1), func(v any) (string, *Color, *Decoration) {
			if s, ok := v.(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
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
	testutil.AssertGolden(t, "common_align_transformer", buf.Bytes())
}

func TestGolden_TableAllPlaceholder(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
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

func TestGolden_TableBasic(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
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

func TestGolden_TableBold(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithDecoration(ScopeBody, Columns(0), DecorationBold),
		WithHeader([]string{"Name", "Value"}),
	)
	if err := tb.Render([][]any{
		{"foo", 1},
		{"bar", 2},
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
	if err := s.Render([]any{"foo", 1}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"bar", 2}); err != nil {
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
	if err := tb.Render([][]any{{"x", "y"}}); err != nil {
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
	if err := s.Render([]any{"x", "y"}); err != nil {
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
	if err := tb.Render([][]any{{"s", "s"}, {"s", "s"}}); err != nil {
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
	for _, r := range [][]any{{"s", "s"}, {"s", "s"}} {
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
	if err := tb.Render([][]any{{"s", "s"}, {"s", "s"}}); err != nil {
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
	for _, r := range [][]any{{"s", "s"}, {"s", "s"}} {
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
	if err := tb.Render([][]any{{"s", "s"}, {"s", "s"}}); err != nil {
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
	for _, r := range [][]any{{"s", "s"}, {"s", "s"}} {
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
	if err := tb.Render([][]any{{"s", "s"}, {"s", "s"}}); err != nil {
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
	for _, r := range [][]any{{"s", "s"}, {"s", "s"}} {
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
	if err := tb.Render([][]any{{"x"}}); err != nil {
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
	if err := s.Render([]any{"x"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_caption_escape", buf.Bytes())
}

func TestGolden_TableCaptionIndex(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithCaption("cap", CaptionBottom),
		WithIndex(),
	)
	if err := tb.Render([][]any{{"s", "s"}, {"s", "s"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_caption_index", buf.Bytes())
}

func TestGolden_StreamCaptionIndex(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
		WithCaption("cap", CaptionBottom),
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
	testutil.AssertGolden(t, "common_caption_index", buf.Bytes())
}

func TestGolden_TableCaptionPlaceholder(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithCaption("cap", CaptionBottom),
		WithPlaceholder("-"),
	)
	if err := tb.Render([][]any{{"x", nil}, {nil, "q"}}); err != nil {
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
	for _, r := range [][]any{{"x", nil}, {nil, "q"}} {
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
	if err := tb.Render([][]any{{"x", "y"}}); err != nil {
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
	if err := s.Render([]any{"x", "y"}); err != nil {
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
		WithTransformer(Columns(1), func(v any) (string, *Color, *Decoration) {
			if s, ok := v.(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
		}),
	)
	if err := tb.Render([][]any{{"x", "raw"}, {"p", "q"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_caption_transformer", buf.Bytes())
}

func TestGolden_StreamCaptionTransformer(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
		WithCaption("cap", CaptionBottom),
		WithTransformer(Columns(1), func(v any) (string, *Color, *Decoration) {
			if s, ok := v.(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
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
	if err := tb.Render([][]any{
		{"a", "  /etc/hosts"},
		{"b", "  /etc/passwd"},
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
	for _, row := range [][]any{{"a", "  /etc/hosts"}, {"b", "  /etc/passwd"}} {
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
	if err := tb.Render([][]any{{"x", "y"}}); err != nil {
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
	if err := s.Render([]any{"x", "y"}); err != nil {
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
	if err := tb.Render([][]any{
		{"alice", 100},
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
	if err := s.Render([]any{"alice", 100}); err != nil {
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
		WithTransformer(Columns(1), func(v any) (string, *Color, *Decoration) {
			if s, ok := v.(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
		}),
	)
	if err := tb.Render([][]any{{"x", "raw"}, {"p", "q"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_cellattr_transformer", buf.Bytes())
}

func TestGolden_StreamCellAttrTransformer(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
		WithCellAttr(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), Attr{Class: "c"}),
		WithTransformer(Columns(1), func(v any) (string, *Color, *Decoration) {
			if s, ok := v.(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
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
	testutil.AssertGolden(t, "common_cellattr_transformer", buf.Bytes())
}

func TestGolden_TableCJK(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"日本語", "ASCII"}),
	)
	if err := tb.Render([][]any{
		{"あいう", "abc"},
		{"日本", "longer-text"},
		{"テスト", "x"},
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
	if err := tb.Render([][]any{{"x", "y"}}); err != nil {
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
	if err := s.Render([]any{"x", "y"}); err != nil {
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
	if err := tb.Render([][]any{{"v"}}); err != nil {
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
	if err := s.Render([]any{"v"}); err != nil {
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
	if err := tb.Render([][]any{
		{"alice", "id-001"},
		{"bob", "id-002"},
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
	if err := s.Render([]any{"alice", "id-001"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"bob", "id-002"}); err != nil {
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
	if err := tb.Render([][]any{
		{"foo", "hello"},
		{"bar", "world"},
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
	if err := s.Render([]any{"foo", "hello"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"bar", "world"}); err != nil {
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
	if err := tb.Render([][]any{
		{"x", "y", 1},
		{"p", "qq", 22},
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
	if err := s.Render([]any{"x", "y", 1}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"p", "qq", 22}); err != nil {
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
	if err := tb.Render([][]any{
		{"x"},
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
	if err := s.Render([]any{"x"}); err != nil {
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
	if err := tb.Render([][]any{
		{"foo", "hello"},
		{"bar", "world"},
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
	if err := s.Render([]any{"foo", "hello"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"bar", "world"}); err != nil {
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
	if err := tb.Render([][]any{{"s", "s"}, {"s", "s"}}); err != nil {
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
	for _, r := range [][]any{{"s", "s"}, {"s", "s"}} {
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
	if err := tb.Render([][]any{
		{"text", "hello"},
		{"slice", []int{1, 2, 3}},
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
	if err := s.Render([]any{"text", "hello"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"slice", []int{1, 2, 3}}); err != nil {
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
	if err := tb.Render([][]any{
		{"<b>", "&amp;"},
		{"a|b", "c\"d"},
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
	if err := s.Render([]any{"<b>", "&amp;"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"a|b", "c\"d"}); err != nil {
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
	if err := tb.Render([][]any{
		{"foo", "hello"},
		{"bar", "world"},
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
	if err := s.Render([]any{"foo", "hello"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"bar", "world"}); err != nil {
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
	if err := tb.Render([][]any{
		{"a", nil},
		{"b", ""},
		{"c", "ok"},
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
	if err := s.Render([]any{"a", nil}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"b", ""}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"c", "ok"}); err != nil {
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
	if err := tb.Render([][]any{
		{"x", "", 1},
		{"y", "b", 2},
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
	if err := s.Render([]any{"x", "", 1}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"y", "b", 2}); err != nil {
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
	if err := tb.Render([][]any{
		{"foo", 1},
		{"bar", 2},
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
	for _, row := range [][]any{{"foo", 1}, {"bar", 2}} {
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
	if err := tb.Render([][]any{
		{"Sales", 100, 100, 100, 200},
		{"Cost", 50, 50, 50, 50},
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
	if err := s.Render([]any{"Sales", 100, 100, 100, 200}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"Cost", 50, 50, 50, 50}); err != nil {
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
	if err := tb.Render([][]any{{"x", "x", "y"}, {"p", "q", "q"}}); err != nil {
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
	for _, r := range [][]any{{"x", "x", "y"}, {"p", "q", "q"}} {
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
	if err := tb.Render([][]any{{"s", "s"}, {"s", "s"}}); err != nil {
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
	for _, r := range [][]any{{"s", "s"}, {"s", "s"}} {
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
	if err := tb.Render([][]any{
		{"x", "x", "y"},
		{"p", "q", "q"},
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
	for _, r := range [][]any{{"x", "x", "y"}, {"p", "q", "q"}} {
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

func TestGolden_TableColspanScope(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
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

func TestGolden_TableColspanTransformer(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1)),
		WithTransformer(Columns(1), func(v any) (string, *Color, *Decoration) {
			if s, ok := v.(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
		}),
	)
	if err := tb.Render([][]any{{"T", "raw"}, {"p", "q"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_colspan_transformer", buf.Bytes())
}

func TestGolden_StreamColspanTransformer(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1)),
		WithTransformer(Columns(1), func(v any) (string, *Color, *Decoration) {
			if s, ok := v.(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
		}),
	)
	for _, r := range [][]any{{"T", "raw"}, {"p", "q"}} {
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

func TestGolden_TableDecoAlign(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithDecoration(ScopeBody, Columns(1), DecorationBold),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), AlignCenter),
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := tb.Render([][]any{
		{"x", "y", 1},
		{"p", "qq", 22},
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
	if err := s.Render([]any{"x", "y", 1}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"p", "qq", 22}); err != nil {
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
	if err := tb.Render([][]any{
		{"amp", "a&b"},
		{"newline", "line1\nline2"},
		{"tag", "<div>"},
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
	if err := s.Render([]any{"amp", "a&b"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"newline", "line1\nline2"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"tag", "<div>"}); err != nil {
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
	if err := tb.Render([][]any{
		{"id-1", "alice", "active"},
		{"id-2", "bob", "inactive"},
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
	if err := s.Render([]any{"id-1", "alice", "active"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"id-2", "bob", "inactive"}); err != nil {
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
	if err := tb.Render([][]any{
		{"a", nil},
		{"b", ""},
		{"c", "ok"},
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
	if err := s.Render([]any{"a", nil}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"b", ""}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"c", "ok"}); err != nil {
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
	if err := tb.Render([][]any{
		{"full", "ok"},
		{"short"},
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
	if err := s.Render([]any{"full", "ok"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"short"}); err != nil {
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
	if err := tb.Render([][]any{{"s", "s"}, {"s", "s"}}); err != nil {
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
	for _, r := range [][]any{{"s", "s"}, {"s", "s"}} {
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

func TestGolden_TableEmptyVsNil(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
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
	testutil.AssertGolden(t, "common_empty_vs_nil", buf.Bytes())
}

func TestGolden_StreamEmptyVsNil(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
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
	testutil.AssertGolden(t, "common_empty_vs_nil", buf.Bytes())
}

func TestGolden_TableEscapeCRLF(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"Key", "Value"}),
	)
	if err := tb.Render([][]any{
		{"crlf", "a\r\nb"},
		{"cr", "a\rb"},
		{"lf", "a\nb"},
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
	if err := s.Render([]any{"crlf", "a\r\nb"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"cr", "a\rb"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"lf", "a\nb"}); err != nil {
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
	if err := tb.Render([][]any{
		{"interior", "a b"},
		{"run", "a  b"},
		{"lead", "  a"},
		{"trail", "a "},
		{"json", "{\n  \"k\": 1\n}"},
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
	if err := s.Render([]any{"interior", "a b"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"run", "a  b"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"lead", "  a"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"trail", "a "}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"json", "{\n  \"k\": 1\n}"}); err != nil {
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
	if err := tb.Render([][]any{
		{"alice", 100},
		{"bob", 200},
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
	if err := s.Render([]any{"alice", 100}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"bob", 200}); err != nil {
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
	if err := tb.Render([][]any{{"s", "s"}, {"s", "s"}}); err != nil {
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
	for _, r := range [][]any{{"s", "s"}, {"s", "s"}} {
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

func TestGolden_TableFooterIndex(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithFooter(func() [][]string {
			return [][]string{{"t", "u"}}
		}),
		WithIndex(),
	)
	if err := tb.Render([][]any{{"s", "s"}, {"s", "s"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_footer_index", buf.Bytes())
}

func TestGolden_StreamFooterIndex(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
		WithFooter(func() [][]string {
			return [][]string{{"t", "u"}}
		}),
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
	testutil.AssertGolden(t, "common_footer_index", buf.Bytes())
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
	if err := tb.Render([][]any{{"x", nil}, {nil, "q"}}); err != nil {
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

func TestGolden_TableFooterTransformer(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithFooter(func() [][]string {
			return [][]string{{"t", "raw"}}
		}),
		WithTransformer(Columns(1), func(v any) (string, *Color, *Decoration) {
			if s, ok := v.(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
		}),
	)
	if err := tb.Render([][]any{{"x", "raw"}, {"p", "q"}}); err != nil {
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
		WithTransformer(Columns(1), func(v any) (string, *Color, *Decoration) {
			if s, ok := v.(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
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

func TestGolden_TableIndexAlign(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithIndex(),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), AlignRight),
	)
	if err := tb.Render([][]any{{"s", "s"}, {"s", "s"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_index_align", buf.Bytes())
}

func TestGolden_StreamIndexAlign(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
		WithIndex(),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), AlignRight),
	)
	for _, r := range [][]any{{"s", "s"}, {"s", "s"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_index_align", buf.Bytes())
}

func TestGolden_TableIndexCellAttr(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithIndex(),
		WithCellAttr(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), Attr{Class: "c"}),
	)
	if err := tb.Render([][]any{{"s", "s"}, {"s", "s"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_index_cellattr", buf.Bytes())
}

func TestGolden_StreamIndexCellAttr(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
		WithIndex(),
		WithCellAttr(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), Attr{Class: "c"}),
	)
	for _, r := range [][]any{{"s", "s"}, {"s", "s"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_index_cellattr", buf.Bytes())
}

func TestGolden_TableIndexColor(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithIndex(),
		WithColor(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), ColorFgRed),
	)
	if err := tb.Render([][]any{{"s", "s"}, {"s", "s"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_index_color", buf.Bytes())
}

func TestGolden_StreamIndexColor(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
		WithIndex(),
		WithColor(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), ColorFgRed),
	)
	for _, r := range [][]any{{"s", "s"}, {"s", "s"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_index_color", buf.Bytes())
}

func TestGolden_TableIndexColspan(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithIndex(),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1, 2)),
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := tb.Render([][]any{{"x", "x", "y"}, {"p", "q", "q"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_index_colspan", buf.Bytes())
}

func TestGolden_StreamIndexColspan(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithIndex(),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1, 2)),
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
	testutil.AssertGolden(t, "common_index_colspan", buf.Bytes())
}

func TestGolden_TableIndexDecoration(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithIndex(),
		WithDecoration(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), DecorationBold),
	)
	if err := tb.Render([][]any{{"s", "s"}, {"s", "s"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_index_decoration", buf.Bytes())
}

func TestGolden_StreamIndexDecoration(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
		WithIndex(),
		WithDecoration(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), DecorationBold),
	)
	for _, r := range [][]any{{"s", "s"}, {"s", "s"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_index_decoration", buf.Bytes())
}

func TestGolden_TableIndexPlaceholder(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithIndex(),
		WithPlaceholder("-"),
	)
	if err := tb.Render([][]any{{"x", nil}, {nil, "q"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_index_placeholder", buf.Bytes())
}

func TestGolden_StreamIndexPlaceholder(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
		WithIndex(),
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
	testutil.AssertGolden(t, "common_index_placeholder", buf.Bytes())
}

func TestGolden_TableIndexTransformer(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithIndex(),
		WithTransformer(Columns(1), func(v any) (string, *Color, *Decoration) {
			if s, ok := v.(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
		}),
	)
	if err := tb.Render([][]any{{"x", "raw"}, {"p", "q"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_index_transformer", buf.Bytes())
}

func TestGolden_StreamIndexTransformer(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
		WithIndex(),
		WithTransformer(Columns(1), func(v any) (string, *Color, *Decoration) {
			if s, ok := v.(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
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
	testutil.AssertGolden(t, "common_index_transformer", buf.Bytes())
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

func TestGolden_TableItalic(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithDecoration(ScopeBody, Columns(1), DecorationItalic),
		WithHeader([]string{"Key", "Note"}),
	)
	if err := tb.Render([][]any{
		{"a", "important"},
		{"b", "optional"},
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
	if err := s.Render([]any{"a", "important"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"b", "optional"}); err != nil {
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
	if err := tb.Render([][]any{
		{"s", 1},
		{strings.Repeat("longvalue", 40), 2},
		{"t", 3},
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
	testutil.AssertGolden(t, "common_long_value", buf.Bytes())
}

func TestGolden_TableNilInNumeric(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithPlaceholder("N/A"),
		WithHeader([]string{"N", "V"}),
	)
	if err := tb.Render([][]any{
		{1, 100},
		{2, nil},
		{3, -5},
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
	if err := s.Render([]any{1, 100}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{2, nil}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{3, -5}); err != nil {
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
	if err := tb.Render([][]any{
		{"a", 1},
		{"b", 2},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_no_header", buf.Bytes())
}

func TestGolden_StreamNoHeader(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf)
	if err := s.Render([]any{"a", 1}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"b", 2}); err != nil {
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

func TestGolden_TablePlaceholderAlign(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithPlaceholder("-"),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), AlignRight),
	)
	if err := tb.Render([][]any{{"x", nil}, {nil, "q"}}); err != nil {
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
	for _, r := range [][]any{{"x", nil}, {nil, "q"}} {
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
	if err := tb.Render([][]any{{"x", nil}, {nil, "q"}}); err != nil {
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
	for _, r := range [][]any{{"x", nil}, {nil, "q"}} {
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
	if err := tb.Render([][]any{{"x", nil, nil}, {nil, nil, "q"}}); err != nil {
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
	for _, r := range [][]any{{"x", nil, nil}, {nil, nil, "q"}} {
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
		WithTransformer(Columns(1), func(v any) (string, *Color, *Decoration) {
			if s, ok := v.(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
		}),
	)
	if err := tb.Render([][]any{{"x", nil}, {"p", "raw"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_placeholder_transformer", buf.Bytes())
}

func TestGolden_StreamPlaceholderTransformer(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
		WithPlaceholder("-"),
		WithTransformer(Columns(1), func(v any) (string, *Color, *Decoration) {
			if s, ok := v.(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
		}),
	)
	for _, r := range [][]any{{"x", nil}, {"p", "raw"}} {
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
	if err := tb.Render([][]any{
		{&str, &s},
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
	if err := s.Render([]any{&str, &stringer}); err != nil {
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
	if err := tb.Render([][]any{
		{"alpha", "first\nsecond"},
		{"beta", "  <value>"},
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
	if err := s.Render([]any{"alpha", "first\nsecond"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"beta", "  <value>"}); err != nil {
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

func TestGolden_TableSingleCell(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
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

func TestGolden_TableSlice(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"Type", "Value"}),
	)
	if err := tb.Render([][]any{
		{"[]int", []int{1, 2, 3}},
		{"[]string", []string{"a", "b"}},
		{"[]bool", []bool{true, false}},
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
	if err := s.Render([]any{"[]int", []int{1, 2, 3}}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"[]string", []string{"a", "b"}}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"[]bool", []bool{true, false}}); err != nil {
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
	if err := tb.Render([][]any{
		{"less-than", "<"},
		{"greater-than", ">"},
		{"double-quote", "\""},
		{"single-quote", "'"},
		{"ampersand", "&"},
		{"space", " "},
		{"asterisk", "*"},
		{"backslash", "\\"},
		{"underscore", "_"},
		{"pipe", "|"},
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
	if err := s.Render([]any{"less-than", "<"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"greater-than", ">"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"double-quote", "\""}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"single-quote", "'"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"ampersand", "&"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"space", " "}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"asterisk", "*"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"backslash", "\\"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"underscore", "_"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"pipe", "|"}); err != nil {
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
	if err := tb.Render([][]any{
		{"login", "deprecated"},
		{"signup", "active"},
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
	if err := s.Render([]any{"login", "deprecated"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"signup", "active"}); err != nil {
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
	if err := tb.Render([][]any{
		{testutil.Stringer{Value: "x"}, testutil.Error{Value: "boom"}},
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
	if err := s.Render([]any{testutil.Stringer{Value: "x"}, testutil.Error{Value: "boom"}}); err != nil {
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
	if err := tb.Render([][]any{
		{"alice", 100},
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
	if err := s.Render([]any{"alice", 100}); err != nil {
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
	if err := tb.Render([][]any{
		{"alice", 100},
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
	if err := s.Render([]any{"alice", 100}); err != nil {
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
	if err := tb.Render([][]any{{"s", "s"}, {"s", "s"}}); err != nil {
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
	for _, r := range [][]any{{"s", "s"}, {"s", "s"}} {
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
	if err := tb.Render([][]any{{"s", "s"}, {"s", "s"}}); err != nil {
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
	for _, r := range [][]any{{"s", "s"}, {"s", "s"}} {
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
	if err := tb.Render([][]any{{"s", "s"}, {"s", "s"}}); err != nil {
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
	for _, r := range [][]any{{"s", "s"}, {"s", "s"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_tableattr_decoration", buf.Bytes())
}

func TestGolden_TableTableAttrIndex(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithTableAttr(TableAttr{
			Table: Attr{Class: "t"},
			Body:  SectionAttr{Cell: Attr{Class: "b"}},
		}),
		WithIndex(),
	)
	if err := tb.Render([][]any{{"s", "s"}, {"s", "s"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_tableattr_index", buf.Bytes())
}

func TestGolden_StreamTableAttrIndex(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
		WithTableAttr(TableAttr{
			Table: Attr{Class: "t"},
			Body:  SectionAttr{Cell: Attr{Class: "b"}},
		}),
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
	testutil.AssertGolden(t, "common_tableattr_index", buf.Bytes())
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
	if err := tb.Render([][]any{{"x", nil}, {nil, "q"}}); err != nil {
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
	for _, r := range [][]any{{"x", nil}, {nil, "q"}} {
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
		WithTransformer(Columns(1), func(v any) (string, *Color, *Decoration) {
			if s, ok := v.(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
		}),
	)
	if err := tb.Render([][]any{{"x", "raw"}, {"p", "q"}}); err != nil {
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
		WithTransformer(Columns(1), func(v any) (string, *Color, *Decoration) {
			if s, ok := v.(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
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
	testutil.AssertGolden(t, "common_tableattr_transformer", buf.Bytes())
}

func TestGolden_TableTransformer(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithTransformer(Columns(1), func(v any) (string, *Color, *Decoration) {
			n, ok := v.(int)
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
		WithTransformer(Columns(1), func(v any) (string, *Color, *Decoration) {
			n, ok := v.(int)
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
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithColor(ScopeBody, Columns(1), ColorFgBlue),
		WithDecoration(ScopeBody, Columns(1), DecorationBold),
		WithTransformer(Columns(1), func(v any) (string, *Color, *Decoration) {
			if s, ok := v.(string); ok && s == "warn" {
				return "", ColorFgYellow, DecorationItalic
			}
			return "", nil, nil
		}),
		WithHeader([]string{"Level", "Message"}),
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
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithColor(ScopeBody, Columns(1), ColorFgBlue),
		WithDecoration(ScopeBody, Columns(1), DecorationBold),
		WithTransformer(Columns(1), func(v any) (string, *Color, *Decoration) {
			if s, ok := v.(string); ok && s == "warn" {
				return "", ColorFgYellow, DecorationItalic
			}
			return "", nil, nil
		}),
		WithHeader([]string{"Level", "Message"}),
	)
	if err := s.Render([]any{"1", "ok"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"2", "warn"}); err != nil {
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
	if err := tb.Render([][]any{{float32(3.14), float64(2.71828)}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_type_float", buf.Bytes())
}

func TestGolden_StreamTypeFloat(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
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
		WithHeader([]string{
			"i", "i8", "i16", "i32", "i64",
			"u", "u8", "u16", "u32", "u64",
		}),
	)
	if err := s.Render([]any{int(-1), int8(-2), int16(-3), int32(-4), int64(-5),
		uint(1), uint8(2), uint16(3), uint32(4), uint64(5)}); err != nil {
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
	if err := tb.Render([][]any{
		{nilStringer, nilError},
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
	if err := s.Render([]any{nilStringer, nilError}); err != nil {
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

func TestGolden_TableValueEqualsPlaceholderColor(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A"}),
		WithPlaceholder("N/A"),
		WithColor(ScopeBody, Columns(0), ColorFgRed),
	)
	if err := tb.Render([][]any{
		{"N/A"},
		{""},
		{"x"},
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
	if err := s.Render([]any{"N/A"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{""}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"x"}); err != nil {
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
	if err := tb.Render([][]any{
		{int64(9223372036854775807), -9223372036854775808, 1.7976931348623157e+308},
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
	if err := s.Render([]any{int64(9223372036854775807), -9223372036854775808, 1.7976931348623157e+308}); err != nil {
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

func TestGolden_StreamCaption(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A"}),
		WithCaption("Stream Caption", CaptionDefault),
	)
	if err := s.Render([]any{"x"}); err != nil {
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
	for _, r := range [][]any{{"s", "s"}, {"s", "s"}} {
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
	if err := s.Render([]any{"g", "x", 1}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"g", "y", 2}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"h", "z", 3}); err != nil {
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
	testutil.AssertGolden(t, "stream_deco_rowspan", buf.Bytes())
}

func TestGolden_StreamDecoSlice(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithDecoration(ScopeBody, Columns(1), DecorationCode),
		WithHeader([]string{"Type", "Values"}),
	)
	if err := s.Render([]any{"ints", []int{1, 2, 3}}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"mixed", []string{"a", "", "c"}}); err != nil {
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
	if err := s.Render([]any{"amp", "a&b"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"newline", "a\nb"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"multiple&", `<one> & "two"`}); err != nil {
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
	for _, r := range [][]any{{"g", "x"}, {"g", "y"}, {"h", "z"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_footer_rowspan", buf.Bytes())
}

func TestGolden_StreamIndex(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithIndex(),
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
	testutil.AssertGolden(t, "stream_index", buf.Bytes())
}

func TestGolden_StreamRowspan(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
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
	testutil.AssertGolden(t, "stream_rowspan", buf.Bytes())
}

func TestGolden_StreamRowspanAlign(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(0), AlignRight),
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := s.Render([]any{"g", "x", 1}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"g", "y", 2}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"h", "z", 3}); err != nil {
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
	for _, r := range [][]any{{"s", "s"}, {"s", "s"}} {
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
	if err := s.Render([]any{"a", "x", "x"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"a", "x", "x"}); err != nil {
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
	testutil.AssertGolden(t, "stream_rowspan_colspan_edge", buf.Bytes())
}

func TestGolden_StreamRowspanEscape(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
		WithHeader([]string{"A", "B"}),
	)
	if err := s.Render([]any{"<x>", 1}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"<x>", 2}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"&y&", 3}); err != nil {
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
	testutil.AssertGolden(t, "stream_rowspan_missing_kinds", buf.Bytes())
}

func TestGolden_StreamRowspanPlaceholder(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
		WithPlaceholder("N/A"),
		WithHeader([]string{"Group", "Value"}),
	)
	if err := s.Render([]any{"A", nil}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"A", "y"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"B", nil}); err != nil {
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
		WithTransformer(Columns(1), func(v any) (string, *Color, *Decoration) {
			if s, ok := v.(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
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
	testutil.AssertGolden(t, "stream_rowspan_transformer", buf.Bytes())
}

func TestGolden_StreamSingleRow(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
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
	for _, r := range [][]any{{"s", "s"}, {"s", "s"}} {
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
	if err := s.Render([]any{"text", "hello"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"number", 42}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"nil", nil}); err != nil {
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

func TestGolden_TableBandColspan(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"Pair", "Pair", "Solo"}),
		WithFooter(func() [][]string {
			return [][]string{{"Sum", "Sum", "Rest"}}
		}),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1, 2)),
	)
	if err := tb.Render([][]any{
		{"x", "x", "z"},
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
	if err := tb.Render([][]any{
		{"foo", 1},
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
	if err := tb.Render([][]any{{"s", "s"}, {"s", "s"}}); err != nil {
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
	if err := tb.Render([][]any{
		{"g", "x", 1},
		{"g", "y", 2},
		{"h", "z", 3},
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
	if err := tb.Render([][]any{
		{"A", "x"},
		{"A", "y"},
		{"B", "z"},
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
	if err := tb.Render([][]any{
		{"ints", []int{1, 2, 3}},
		{"strings", []string{"a", "b"}},
		{"empty-elem", []string{"x", "", "z"}},
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
	if err := tb.Render([][]any{
		{"ampersand", "a&b"},
		{"less-than", "a<b"},
		{"greater-than", "a>b"},
		{"double-quote", `a"b`},
		{"newline", "a\nb"},
		{"combined", "a<b&c\nd>e"},
		{"multiple&", `<one> & "two"`},
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
	if err := tb.Render([][]any{{"g", "x"}, {"g", "y"}, {"h", "z"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_footer_rowspan", buf.Bytes())
}

func TestGolden_TableIndex(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithIndex(),
		WithHeader([]string{"Name", "Score"}),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
	)
	if err := tb.Render([][]any{
		{"alice", 100},
		{"alice", 99},
		{"bob", 98},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_index", buf.Bytes())
}

func TestGolden_TableNoHeaderRagged(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf)
	if err := tb.Render([][]any{
		{"a", 1},
		{"b"},
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
	if err := tb.Render([][]any{
		{"A", "x"},
		{"A", "y"},
		{"B", "z"},
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
	if err := tb.Render([][]any{
		{"g", "x", 1},
		{"g", "y", 2},
		{"h", "z", 3},
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
	if err := tb.Render([][]any{{"s", "s"}, {"s", "s"}}); err != nil {
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
	if err := tb.Render([][]any{
		{"x", "x", "x"},
		{"x", "x", "y"},
		{"z", "z", "y"},
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
	if err := tb.Render([][]any{
		{"g", "x", "x"},
		{"g", "y", "y"},
		{"h", "z", "w"},
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
	if err := tb.Render([][]any{
		{"<x>", 1},
		{"<x>", 2},
		{"&y&", 3},
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
	if err := tb.Render([][]any{
		{1, "X"},
		{2},
		{3, nil},
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
	if err := tb.Render([][]any{
		{"A", nil},
		{"A", "y"},
		{"B", nil},
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
		WithTransformer(Columns(1), func(v any) (string, *Color, *Decoration) {
			if s, ok := v.(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
		}),
	)
	if err := tb.Render([][]any{{"x", "raw"}, {"p", "raw"}, {"q", "z"}}); err != nil {
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
	if err := tb.Render([][]any{{"s", "s"}, {"s", "s"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_tableattr_rowspan", buf.Bytes())
}

func TestGolden_TableTypeMixed(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
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
