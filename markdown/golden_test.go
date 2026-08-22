package markdown

import (
	"bytes"
	"strings"
	"testing"

	"github.com/nekrassov01/table/internal/param"
	"github.com/nekrassov01/table/internal/testutil"
)

func TestGolden_StreamAlign(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithAlign(Columns(0), AlignRight),
		WithAlign(Columns(1), AlignCenter),
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
	testutil.AssertGolden(t, "stream_align", buf.Bytes())
}

func TestGolden_StreamAlignCenter(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithAlign(Columns(1), AlignCenter),
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
	testutil.AssertGolden(t, "stream_align_center", buf.Bytes())
}

func TestGolden_StreamAlignColspan(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithColspan(Columns(0, 1, 2)),
		WithAlign(Columns(1), AlignRight),
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
	testutil.AssertGolden(t, "stream_align_colspan", buf.Bytes())
}

func TestGolden_StreamAlignLeft(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithAlign(Columns(1), AlignLeft),
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
	testutil.AssertGolden(t, "stream_align_left", buf.Bytes())
}

func TestGolden_StreamAlignMixed(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithAlign(Columns(0), AlignRight),
		WithAlign(Columns(1), AlignCenter),
		WithAlign(Columns(2), AlignLeft),
		WithHeader([]string{"Right", "Center", "Left"}),
	)
	if err := s.Render([]any{"a", "b", "c"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"longer", "y", "x"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_align_mixed", buf.Bytes())
}

func TestGolden_StreamAlignRight(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithAlign(Columns(1), AlignRight),
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
	testutil.AssertGolden(t, "stream_align_right", buf.Bytes())
}

func TestGolden_StreamAlignTransformer(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithAlign(Columns(1), AlignCenter),
		WithTransformer(Columns(1), func(v any) (string, *Color, *Decoration) {
			if s, ok := v.(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
		}),

		WithHeader([]string{"A", "B"}),
	)
	for _, r := range [][]any{{"x", "raw"}, {"p", "q"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_align_transformer", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_all_placeholder", buf.Bytes())
}

func TestGolden_StreamAutolink(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"Kind", "Value"}),
	)
	for _, r := range [][]any{
		{"www", "www.commonmark.org"},
		{"https", "https://example.com/path"},
		{"email", "foo@bar.baz"},
		{"mailto", "mailto:foo@bar.baz"},
	} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_autolink", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_basic", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_bold", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_cjk", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_code", buf.Bytes())
}

func TestGolden_StreamCodeSpanEdges(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"Case", "Value"}),
		WithDecoration(ScopeBody, Columns(1), DecorationCode),
	)
	if err := s.Render([]any{"backslash", `a\b`}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"one backtick", "a`b"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"two backticks", "a``b"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"leading backtick", "`x"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"trailing backtick", "x`"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"only backtick", "`"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"pipe", `a|b`}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"backslash then pipe", `a\|b`}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"padded spaces", " x "}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"only spaces", "   "}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_code_span_edges", buf.Bytes())
}

func TestGolden_StreamColor(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithColor(ScopeBody, Columns(1), ColorFgRed),
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
	testutil.AssertGolden(t, "stream_color", buf.Bytes())
}

func TestGolden_StreamColorAlign(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithColor(ScopeBody, Columns(1), ColorFgRed),
		WithAlign(Columns(1), AlignRight),
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
	testutil.AssertGolden(t, "stream_color_align", buf.Bytes())
}

func TestGolden_StreamColorAttrEscape(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A"}),
		WithColor(ScopeBody, Columns(0), NewColor(`red" onmouseover="alert(1)`, `blue&x|y`)),
	)
	if err := s.Render([]any{"x"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_color_attr_escape", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_color_bg", buf.Bytes())
}

func TestGolden_StreamColorCode(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithDecoration(ScopeBody, Columns(1), DecorationCode),
		WithColor(ScopeBody, Columns(1), ColorFgRed),
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
	testutil.AssertGolden(t, "stream_color_code", buf.Bytes())
}

func TestGolden_StreamColorEscape(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithColor(ScopeBody, Columns(1), ColorFgRed),
		WithHeader([]string{"Key", "Value"}),
	)
	if err := s.Render([]any{"space", "a b"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"pipe", "a|b"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_color_escape", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_color_fg_bg", buf.Bytes())
}

func TestGolden_StreamColorNil(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithColor(ScopeBody, Columns(1), ColorFgRed),
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
	testutil.AssertGolden(t, "stream_color_nil", buf.Bytes())
}

func TestGolden_StreamColorPlaceholder(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithColor(ScopeBody, Columns(1), ColorFgCyan),
		WithPlaceholder("N/A"),
		WithHeader([]string{"Key", "Value"}),
	)
	if err := s.Render([]any{"a", nil}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"b", "ok"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_color_placeholder", buf.Bytes())
}

func TestGolden_StreamColorRowspan(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithColor(ScopeBody, Columns(1), ColorFgBlue),
		WithRowspan(Columns(0)),
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
	testutil.AssertGolden(t, "stream_color_rowspan", buf.Bytes())
}

func TestGolden_StreamColorScope(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"Name", "Value"}),
		WithColor(ScopeHeader, Columns(0), ColorFgRed),
		WithColor(ScopeBody, Columns(1), ColorFgBlue),
		WithDecoration(ScopeHeader|ScopeBody, Columns(1), DecorationBold),
	)
	for _, row := range [][]any{{"foo", 1}, {"bar", 2}} {
		if err := s.Render(row); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_color_scope", buf.Bytes())
}

func TestGolden_StreamColspan(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithColspan(Columns(0, 1, 2)),
		WithHeader([]string{"A", "B", "C"}),
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
	testutil.AssertGolden(t, "stream_colspan", buf.Bytes())
}

func TestGolden_StreamColspanColor(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithColspan(Columns(0, 1, 2)),
		WithColor(ScopeHeader|ScopeBody, Columns(1), ColorFgRed),
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
	testutil.AssertGolden(t, "stream_colspan_color", buf.Bytes())
}

func TestGolden_StreamColspanDecoration(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithColspan(Columns(0, 1, 2)),
		WithDecoration(ScopeHeader|ScopeBody, Columns(1), DecorationBold),
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
	testutil.AssertGolden(t, "stream_colspan_decoration", buf.Bytes())
}

func TestGolden_StreamColspanEdges(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B", "C", "D"}),
		WithColspan(Columns(0, 1)),
		WithColspan(Columns(2, 3)),
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
	testutil.AssertGolden(t, "stream_colspan_edges", buf.Bytes())
}

func TestGolden_StreamColspanTransformer(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithColspan(Columns(0, 1, 2)),
		WithTransformer(Columns(1), func(v any) (string, *Color, *Decoration) {
			if s, ok := v.(string); ok && s == "raw" {
				return "x", ColorFgRed, nil
			}
			return "", nil, nil
		}),

		WithHeader([]string{"A", "B", "C"}),
	)
	for _, r := range [][]any{{"x", "raw", "y"}, {"p", "raw", "z"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_colspan_transformer", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_control_chars", buf.Bytes())
}

func TestGolden_StreamDecoAlign(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithDecoration(ScopeBody, Columns(1), DecorationBold),
		WithAlign(Columns(1), AlignRight),
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
	testutil.AssertGolden(t, "stream_deco_align", buf.Bytes())
}

func TestGolden_StreamDecoColorAlign(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithDecoration(ScopeBody, Columns(1), DecorationItalic),
		WithColor(ScopeBody, Columns(1), ColorFgGreen),
		WithAlign(Columns(1), AlignCenter),
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
	testutil.AssertGolden(t, "stream_deco_color_align", buf.Bytes())
}

func TestGolden_StreamDecoColorRowspan(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithDecoration(ScopeBody, Columns(1), DecorationBold),
		WithColor(ScopeBody, Columns(1), ColorFgRed),
		WithRowspan(Columns(0)),
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
	testutil.AssertGolden(t, "stream_deco_color_rowspan", buf.Bytes())
}

func TestGolden_StreamDecoEscape(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithDecoration(ScopeBody, Columns(1), DecorationCode),
		WithHeader([]string{"Key", "Value"}),
	)
	if err := s.Render([]any{"pipe", "a|b"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"newline", "line1\nline2"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"space", "a b"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_deco_escape", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_deco_multi_col", buf.Bytes())
}

func TestGolden_StreamDecoNil(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithDecoration(ScopeBody, Columns(1), DecorationBold),
		WithPlaceholder("-"),
		WithHeader([]string{"Key", "Value"}),
	)
	if err := s.Render([]any{"a", nil}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"b", "ok"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_deco_nil", buf.Bytes())
}

func TestGolden_StreamDecoRowspan(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithDecoration(ScopeBody, Columns(0), DecorationBold),
		WithRowspan(Columns(0)),
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
	testutil.AssertGolden(t, "stream_deco_short_row", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_emoji", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_empty_header_label", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_empty_vs_nil", buf.Bytes())
}

func TestGolden_StreamEscape(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"Key", "Value"}),
	)
	if err := s.Render([]any{"pipe", "a|b"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"newline", "a\nb"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"a|b", "*c* & <d>"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_escape", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_escape_crlf", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_escape_space", buf.Bytes())
}

func TestGolden_StreamHeaderOnly(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_header_only", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_header_wider_than_rows", buf.Bytes())
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

func TestGolden_StreamIndexAlign(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithIndex(),
		WithAlign(Columns(1), AlignRight),
		WithHeader([]string{"A", "B"}),
	)
	for _, r := range [][]any{{"x", "y"}, {"p", "q"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_index_align", buf.Bytes())
}

func TestGolden_StreamIndexColor(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithIndex(),
		WithColor(ScopeHeader|ScopeBody, Columns(1), ColorFgRed),
		WithHeader([]string{"A", "B"}),
	)
	for _, r := range [][]any{{"x", "y"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_index_color", buf.Bytes())
}

func TestGolden_StreamIndexColspan(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithIndex(),
		WithColspan(Columns(0, 1, 2)),
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
	testutil.AssertGolden(t, "stream_index_colspan", buf.Bytes())
}

func TestGolden_StreamIndexDecoration(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithIndex(),
		WithDecoration(ScopeHeader|ScopeBody, Columns(1), DecorationBold),
		WithHeader([]string{"A", "B"}),
	)
	for _, r := range [][]any{{"x", "y"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_index_decoration", buf.Bytes())
}

func TestGolden_StreamIndexPlaceholder(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithIndex(),
		WithPlaceholder("-"),
		WithHeader([]string{"A", "B"}),
	)
	for _, r := range [][]any{{"x", nil}, {nil, "q"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_index_placeholder", buf.Bytes())
}

func TestGolden_StreamIndexRowspan(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithIndex(),
		WithRowspan(Columns(0)),
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
	testutil.AssertGolden(t, "stream_index_rowspan", buf.Bytes())
}

func TestGolden_StreamIndexTransformer(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithIndex(),
		WithTransformer(Columns(1), func(v any) (string, *Color, *Decoration) {
			if s, ok := v.(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
		}),

		WithHeader([]string{"A", "B"}),
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
	testutil.AssertGolden(t, "stream_invalid_utf8", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_italic", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_nil_in_numeric", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_placeholder", buf.Bytes())
}

func TestGolden_StreamPlaceholderColspan(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithPlaceholder("-"),
		WithColspan(Columns(0, 1, 2)),
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
	testutil.AssertGolden(t, "stream_placeholder_colspan", buf.Bytes())
}

func TestGolden_StreamPlaceholderTransformer(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithPlaceholder("-"),
		WithTransformer(Columns(1), func(v any) (string, *Color, *Decoration) {
			if s, ok := v.(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
		}),

		WithHeader([]string{"A", "B"}),
	)
	for _, r := range [][]any{{"x", nil}, {"p", "raw"}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_placeholder_transformer", buf.Bytes())
}

func TestGolden_StreamPointer(t *testing.T) {
	var buf bytes.Buffer
	s := testutil.Stringer{Value: "y"}
	str := "alive"
	st := NewStream(&buf,
		WithHeader([]string{"a", "b"}),
	)
	if err := st.Render([]any{&str, &s}); err != nil {
		t.Fatal(err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_pointer", buf.Bytes())
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
	if err := s.Render([]any{"gamma", "{\n  \"k\": 1\n}"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_preformatted", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_ragged_rows", buf.Bytes())
}

func TestGolden_StreamRowspan(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithRowspan(Columns(0)),
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
		WithRowspan(Columns(0)),
		WithAlign(Columns(1), AlignRight),
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
	testutil.AssertGolden(t, "stream_rowspan_align", buf.Bytes())
}

func TestGolden_StreamRowspanAlignPlaceholder(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithRowspan(Columns(0)),
		WithAlign(Columns(1), AlignCenter),
		WithPlaceholder("-"),
		WithHeader([]string{"Group", "Score"}),
	)
	if err := s.Render([]any{"A", 100}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"A", nil}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"B", 200}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_rowspan_align_placeholder", buf.Bytes())
}

func TestGolden_StreamRowspanColspanEdge(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B", "C"}),
		WithRowspan(Columns(0)),
		WithColspan(Columns(1, 2)),
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
		WithRowspan(Columns(0)),
		WithHeader([]string{"Group", "Value"}),
	)
	if err := s.Render([]any{"a|b", "x"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"a|b", "y"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"c\\d", "z"}); err != nil {
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
		WithRowspan(Columns(1)),
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
		WithRowspan(Columns(0)),
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
		WithRowspan(Columns(1)),
		WithTransformer(Columns(1), func(v any) (string, *Color, *Decoration) {
			if s, ok := v.(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
		}),

		WithHeader([]string{"A", "B"}),
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
	testutil.AssertGolden(t, "stream_single_cell", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_slice", buf.Bytes())
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
		WithColspan(Columns(cols...)),
		WithRowspan(Columns(cols...)),
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
	if err := s.Render([]any{"backtick", "`"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"backslash", "\\"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"left-bracket", "["}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"right-bracket", "]"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"tilde", "~"}); err != nil {
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
	testutil.AssertGolden(t, "stream_special_chars", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_strikethrough", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_stringer_error", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_transformer", buf.Bytes())
}

func TestGolden_StreamTransformerColumnOverride(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithColor(ScopeBody, Columns(1), ColorFgBlue),
		WithDecoration(ScopeBody, Columns(1), DecorationCode),
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
	testutil.AssertGolden(t, "stream_transformer_column_override", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_type_float", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_type_integer", buf.Bytes())
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

func TestGolden_StreamTypedNil(t *testing.T) {
	var buf bytes.Buffer
	var nilStringer *testutil.PtrStringer
	var nilError *testutil.PtrError
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
	testutil.AssertGolden(t, "stream_typed_nil", buf.Bytes())
}

func TestGolden_StreamUnderline(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithDecoration(ScopeBody, Columns(1), DecorationUnderline),
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
	testutil.AssertGolden(t, "stream_underline", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_value_equals_placeholder", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_value_equals_placeholder_color", buf.Bytes())
}

func TestGolden_StreamWideNumber(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
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
	testutil.AssertGolden(t, "stream_zero_width", buf.Bytes())
}

func TestGolden_TableAlignCenter(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithAlign(Columns(1), AlignCenter),
		WithHeader([]string{"Name", "Score"}),
	)
	if err := tb.Render([][]any{
		{"alice", 100},
		{"bob", 99},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_align_center", buf.Bytes())
}

func TestGolden_TableAlignColspan(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithColspan(Columns(0, 1, 2)),
		WithAlign(Columns(1), AlignRight),
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := tb.Render([][]any{{"x", "x", "y"}, {"p", "q", "q"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_align_colspan", buf.Bytes())
}

func TestGolden_TableAlignLeft(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithAlign(Columns(1), AlignLeft),
		WithHeader([]string{"Name", "Score"}),
	)
	if err := tb.Render([][]any{
		{"alice", 100},
		{"bob", 99},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_align_left", buf.Bytes())
}

func TestGolden_TableAlignMixed(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithAlign(Columns(0), AlignRight),
		WithAlign(Columns(1), AlignCenter),
		WithAlign(Columns(2), AlignLeft),
		WithHeader([]string{"Right", "Center", "Left"}),
	)
	if err := tb.Render([][]any{
		{"a", "b", "c"},
		{"longer", "y", "x"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_align_mixed", buf.Bytes())
}

func TestGolden_TableAlignRight(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithAlign(Columns(1), AlignRight),
		WithHeader([]string{"Name", "Score"}),
	)
	if err := tb.Render([][]any{
		{"alice", 100},
		{"bob", 99},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_align_right", buf.Bytes())
}

func TestGolden_TableAlignTransformer(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithAlign(Columns(1), AlignCenter),
		WithTransformer(Columns(1), func(v any) (string, *Color, *Decoration) {
			if s, ok := v.(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
		}),

		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]any{{"x", "raw"}, {"p", "q"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_align_transformer", buf.Bytes())
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
	testutil.AssertGolden(t, "table_all_placeholder", buf.Bytes())
}

func TestGolden_TableAutolink(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"Kind", "Value"}),
	)
	if err := tb.Render([][]any{
		{"www", "www.commonmark.org"},
		{"https", "https://example.com/path"},
		{"email", "foo@bar.baz"},
		{"mailto", "mailto:foo@bar.baz"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_autolink", buf.Bytes())
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
	testutil.AssertGolden(t, "table_basic", buf.Bytes())
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
	testutil.AssertGolden(t, "table_bold", buf.Bytes())
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
	testutil.AssertGolden(t, "table_cjk", buf.Bytes())
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
	testutil.AssertGolden(t, "table_code", buf.Bytes())
}

func TestGolden_TableCodeSpanEdges(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"Case", "Value"}),
		WithDecoration(ScopeBody, Columns(1), DecorationCode),
	)
	if err := tb.Render([][]any{
		{"backslash", `a\b`},
		{"one backtick", "a`b"},
		{"two backticks", "a``b"},
		{"leading backtick", "`x"},
		{"trailing backtick", "x`"},
		{"only backtick", "`"},
		{"pipe", `a|b`},
		{"backslash then pipe", `a\|b`},
		{"padded spaces", " x "},
		{"only spaces", "   "},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_code_span_edges", buf.Bytes())
}

func TestGolden_TableColor(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithColor(ScopeBody, Columns(1), ColorFgRed),
		WithHeader([]string{"Name", "Value"}),
	)
	if err := tb.Render([][]any{
		{"foo", "hello"},
		{"bar", "world"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_color", buf.Bytes())
}

func TestGolden_TableColorAlign(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithColor(ScopeBody, Columns(1), ColorFgRed),
		WithAlign(Columns(1), AlignRight),
		WithHeader([]string{"Name", "Score"}),
	)
	if err := tb.Render([][]any{
		{"alice", 100},
		{"bob", 99},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_color_align", buf.Bytes())
}

func TestGolden_TableColorAttrEscape(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A"}),
		WithColor(ScopeBody, Columns(0), NewColor(`red" onmouseover="alert(1)`, `blue&x|y`)),
	)
	if err := tb.Render([][]any{
		{"x"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_color_attr_escape", buf.Bytes())
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
	testutil.AssertGolden(t, "table_color_bg", buf.Bytes())
}

func TestGolden_TableColorCode(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithDecoration(ScopeBody, Columns(1), DecorationCode),
		WithColor(ScopeBody, Columns(1), ColorFgRed),
		WithHeader([]string{"Type", "Value"}),
	)
	if err := tb.Render([][]any{
		{"text", "hello"},
		{"slice", []int{1, 2, 3}},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_color_code", buf.Bytes())
}

func TestGolden_TableColorEscape(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithColor(ScopeBody, Columns(1), ColorFgRed),
		WithHeader([]string{"Key", "Value"}),
	)
	if err := tb.Render([][]any{
		{"space", "a b"},
		{"pipe", "a|b"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_color_escape", buf.Bytes())
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
	testutil.AssertGolden(t, "table_color_fg_bg", buf.Bytes())
}

func TestGolden_TableColorNil(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithColor(ScopeBody, Columns(1), ColorFgRed),
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
	testutil.AssertGolden(t, "table_color_nil", buf.Bytes())
}

func TestGolden_TableColorPlaceholder(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithColor(ScopeBody, Columns(1), ColorFgCyan),
		WithPlaceholder("N/A"),
		WithHeader([]string{"Key", "Value"}),
	)
	if err := tb.Render([][]any{
		{"a", nil},
		{"b", "ok"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_color_placeholder", buf.Bytes())
}

func TestGolden_TableColorRowspan(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithColor(ScopeBody, Columns(1), ColorFgBlue),
		WithRowspan(Columns(0)),
		WithHeader([]string{"Group", "Value"}),
	)
	if err := tb.Render([][]any{
		{"A", "x"},
		{"A", "y"},
		{"B", "z"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_color_rowspan", buf.Bytes())
}

func TestGolden_TableColorScope(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"Name", "Value"}),
		WithColor(ScopeHeader, Columns(0), ColorFgRed),
		WithColor(ScopeBody, Columns(1), ColorFgBlue),
		WithDecoration(ScopeHeader|ScopeBody, Columns(1), DecorationBold),
	)
	if err := tb.Render([][]any{
		{"foo", 1},
		{"bar", 2},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_color_scope", buf.Bytes())
}

func TestGolden_TableColspan(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithColspan(Columns(0, 1, 2)),
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := tb.Render([][]any{
		{"x", "x", "y"},
		{"p", "q", "q"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_colspan", buf.Bytes())
}

func TestGolden_TableColspanColor(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithColspan(Columns(0, 1, 2)),
		WithColor(ScopeHeader|ScopeBody, Columns(1), ColorFgRed),
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := tb.Render([][]any{{"x", "x", "y"}, {"p", "q", "q"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_colspan_color", buf.Bytes())
}

func TestGolden_TableColspanDecoration(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithColspan(Columns(0, 1, 2)),
		WithDecoration(ScopeHeader|ScopeBody, Columns(1), DecorationBold),
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := tb.Render([][]any{
		{"x", "x", "y"},
		{"p", "q", "q"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_colspan_decoration", buf.Bytes())
}

func TestGolden_TableColspanEdges(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B", "C", "D"}),
		WithColspan(Columns(0, 1)),
		WithColspan(Columns(2, 3)),
	)
	if err := tb.Render([][]any{
		{"x", "x", "y", "y"},
		{"p", "q", "r", "r"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_colspan_edges", buf.Bytes())
}

func TestGolden_TableColspanTransformer(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithColspan(Columns(0, 1, 2)),
		WithTransformer(Columns(1), func(v any) (string, *Color, *Decoration) {
			if s, ok := v.(string); ok && s == "raw" {
				return "x", ColorFgRed, nil
			}
			return "", nil, nil
		}),

		WithHeader([]string{"A", "B", "C"}),
	)
	if err := tb.Render([][]any{{"x", "raw", "y"}, {"p", "raw", "z"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_colspan_transformer", buf.Bytes())
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
	testutil.AssertGolden(t, "table_control_chars", buf.Bytes())
}

func TestGolden_TableDecoAlign(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithDecoration(ScopeBody, Columns(1), DecorationBold),
		WithAlign(Columns(1), AlignRight),
		WithHeader([]string{"Name", "Score"}),
	)
	if err := tb.Render([][]any{
		{"alice", 100},
		{"bob", 99},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_deco_align", buf.Bytes())
}

func TestGolden_TableDecoColorAlign(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithDecoration(ScopeBody, Columns(1), DecorationItalic),
		WithColor(ScopeBody, Columns(1), ColorFgMagenta),
		WithAlign(Columns(1), AlignCenter),
		WithHeader([]string{"Key", "Note"}),
	)
	if err := tb.Render([][]any{
		{"a", "important"},
		{"b", "optional"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_deco_color_align", buf.Bytes())
}

func TestGolden_TableDecoColorRowspan(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithDecoration(ScopeBody, Columns(1), DecorationBold),
		WithColor(ScopeBody, Columns(1), ColorFgRed),
		WithRowspan(Columns(0)),
		WithHeader([]string{"Group", "Value"}),
	)
	if err := tb.Render([][]any{
		{"A", "x"},
		{"A", "y"},
		{"B", "z"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_deco_color_rowspan", buf.Bytes())
}

func TestGolden_TableDecoEscape(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithDecoration(ScopeBody, Columns(1), DecorationCode),
		WithHeader([]string{"Key", "Value"}),
	)
	if err := tb.Render([][]any{
		{"pipe", "a|b"},
		{"newline", "line1\nline2"},
		{"space", "a b"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_deco_escape", buf.Bytes())
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
	testutil.AssertGolden(t, "table_deco_multi_col", buf.Bytes())
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
	testutil.AssertGolden(t, "table_deco_nil", buf.Bytes())
}

func TestGolden_TableDecoRowspan(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithDecoration(ScopeBody, Columns(0), DecorationBold),
		WithRowspan(Columns(0)),
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
	testutil.AssertGolden(t, "table_deco_short_row", buf.Bytes())
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
	testutil.AssertGolden(t, "table_emoji", buf.Bytes())
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
	testutil.AssertGolden(t, "table_empty_header_label", buf.Bytes())
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
	testutil.AssertGolden(t, "table_empty_vs_nil", buf.Bytes())
}

func TestGolden_TableEscape(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"Key", "Value"}),
	)
	if err := tb.Render([][]any{
		{"pipe", "a|b"},
		{"backslash", "a\\b"},
		{"newline", "a\nb"},
		{"combined", "a|b\nc\\d"},
		{"a|b", "*c* & <d>"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_escape", buf.Bytes())
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
	testutil.AssertGolden(t, "table_escape_crlf", buf.Bytes())
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
	testutil.AssertGolden(t, "table_escape_space", buf.Bytes())
}

func TestGolden_TableHeaderOnly(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := tb.Render(nil); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_header_only", buf.Bytes())
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
	testutil.AssertGolden(t, "table_header_wider_than_rows", buf.Bytes())
}

func TestGolden_TableIndex(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithIndex(),
		WithHeader([]string{"Name", "Score"}),
	)
	if err := tb.Render([][]any{
		{"alice", 100},
		{"bob", 99},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_index", buf.Bytes())
}

func TestGolden_TableIndexAlign(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithIndex(),
		WithAlign(Columns(1), AlignRight),
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]any{{"x", "y"}, {"p", "q"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_index_align", buf.Bytes())
}

func TestGolden_TableIndexColor(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithIndex(),
		WithColor(ScopeHeader|ScopeBody, Columns(1), ColorFgRed),
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]any{{"x", "y"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_index_color", buf.Bytes())
}

func TestGolden_TableIndexColspan(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithIndex(),
		WithColspan(Columns(0, 1, 2)),
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := tb.Render([][]any{{"x", "x", "y"}, {"p", "q", "q"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_index_colspan", buf.Bytes())
}

func TestGolden_TableIndexDecoration(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithIndex(),
		WithDecoration(ScopeHeader|ScopeBody, Columns(1), DecorationBold),
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]any{{"x", "y"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_index_decoration", buf.Bytes())
}

func TestGolden_TableIndexPlaceholder(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithIndex(),
		WithPlaceholder("-"),
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]any{{"x", nil}, {nil, "q"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_index_placeholder", buf.Bytes())
}

func TestGolden_TableIndexRowspan(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithIndex(),
		WithRowspan(Columns(0)),
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]any{{"g", "x"}, {"g", "y"}, {"h", "z"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_index_rowspan", buf.Bytes())
}

func TestGolden_TableIndexTransformer(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithIndex(),
		WithTransformer(Columns(1), func(v any) (string, *Color, *Decoration) {
			if s, ok := v.(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
		}),

		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]any{{"x", "raw"}, {"p", "q"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_index_transformer", buf.Bytes())
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
	testutil.AssertGolden(t, "table_invalid_utf8", buf.Bytes())
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
	testutil.AssertGolden(t, "table_italic", buf.Bytes())
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
	testutil.AssertGolden(t, "table_nil_in_numeric", buf.Bytes())
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
	testutil.AssertGolden(t, "table_placeholder", buf.Bytes())
}

func TestGolden_TablePlaceholderColspan(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithPlaceholder("-"),
		WithColspan(Columns(0, 1, 2)),
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := tb.Render([][]any{{"x", nil, nil}, {nil, nil, "q"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_placeholder_colspan", buf.Bytes())
}

func TestGolden_TablePlaceholderTransformer(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithPlaceholder("-"),
		WithTransformer(Columns(1), func(v any) (string, *Color, *Decoration) {
			if s, ok := v.(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
		}),

		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]any{{"x", nil}, {"p", "raw"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_placeholder_transformer", buf.Bytes())
}

func TestGolden_TablePointer(t *testing.T) {
	var buf bytes.Buffer
	s := testutil.Stringer{Value: "y"}
	str := "alive"
	tb := NewTable(&buf,
		WithHeader([]string{"a", "b"}),
	)
	if err := tb.Render([][]any{{&str, &s}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_pointer", buf.Bytes())
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
		{"gamma", "{\n  \"k\": 1\n}"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_preformatted", buf.Bytes())
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
	testutil.AssertGolden(t, "table_ragged_rows", buf.Bytes())
}

func TestGolden_TableRowspan(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithRowspan(Columns(0)),
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
		WithRowspan(Columns(0)),
		WithAlign(Columns(1), AlignRight),
		WithHeader([]string{"Group", "Score"}),
	)
	if err := tb.Render([][]any{
		{"A", 100},
		{"A", 200},
		{"B", 300},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_rowspan_align", buf.Bytes())
}

func TestGolden_TableRowspanAlignPlaceholder(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithRowspan(Columns(0)),
		WithAlign(Columns(1), AlignCenter),
		WithPlaceholder("-"),
		WithHeader([]string{"Group", "Score"}),
	)
	if err := tb.Render([][]any{
		{"A", 100},
		{"A", nil},
		{"B", 200},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_rowspan_align_placeholder", buf.Bytes())
}

func TestGolden_TableRowspanColspanEdge(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B", "C"}),
		WithRowspan(Columns(0)),
		WithColspan(Columns(1, 2)),
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
		WithRowspan(Columns(0)),
		WithHeader([]string{"Group", "Value"}),
	)
	if err := tb.Render([][]any{
		{"a|b", "x"},
		{"a|b", "y"},
		{"c\\d", "z"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_rowspan_escape", buf.Bytes())
}

func TestGolden_TableRowspanMissingKinds(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithRowspan(Columns(1)),
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
		WithRowspan(Columns(0)),
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
		WithRowspan(Columns(1)),
		WithTransformer(Columns(1), func(v any) (string, *Color, *Decoration) {
			if s, ok := v.(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
		}),

		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]any{{"x", "raw"}, {"p", "raw"}, {"q", "z"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_rowspan_transformer", buf.Bytes())
}

func TestGolden_TableSingleCell(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"only"}),
	)
	if err := tb.Render([][]any{{"v"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_single_cell", buf.Bytes())
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
	testutil.AssertGolden(t, "table_slice", buf.Bytes())
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
		WithColspan(Columns(cols...)),
		WithRowspan(Columns(cols...)),
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
		{"backtick", "`"},
		{"backslash", "\\"},
		{"left-bracket", "["},
		{"right-bracket", "]"},
		{"tilde", "~"},
		{"underscore", "_"},
		{"pipe", "|"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_special_chars", buf.Bytes())
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
	testutil.AssertGolden(t, "table_strikethrough", buf.Bytes())
}

func TestGolden_TableStringerError(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"Stringer", "Error"}),
	)
	if err := tb.Render([][]any{{testutil.Stringer{Value: "x"}, testutil.Error{Value: "boom"}}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_stringer_error", buf.Bytes())
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
	testutil.AssertGolden(t, "table_transformer", buf.Bytes())
}

func TestGolden_TableTransformerColumnOverride(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithColor(ScopeBody, Columns(1), ColorFgBlue),
		WithDecoration(ScopeBody, Columns(1), DecorationCode),
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
	testutil.AssertGolden(t, "table_transformer_column_override", buf.Bytes())
}

func TestGolden_TableTypeFloat(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"f32", "f64"}),
	)
	if err := tb.Render([][]any{{float32(3.14), float64(2.71828)}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_type_float", buf.Bytes())
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
	testutil.AssertGolden(t, "table_type_integer", buf.Bytes())
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

func TestGolden_TableTypedNil(t *testing.T) {
	var buf bytes.Buffer
	var nilStringer *testutil.PtrStringer
	var nilError *testutil.PtrError
	tb := NewTable(&buf,
		WithPlaceholder("<nil>"),
		WithHeader([]string{"Stringer", "Error"}),
	)
	if err := tb.Render([][]any{{nilStringer, nilError}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_typed_nil", buf.Bytes())
}

func TestGolden_TableUnderline(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithDecoration(ScopeBody, Columns(1), DecorationUnderline),
		WithHeader([]string{"Key", "Note"}),
	)
	if err := tb.Render([][]any{
		{"a", "important"},
		{"b", "optional"},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_underline", buf.Bytes())
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
	testutil.AssertGolden(t, "table_value_equals_placeholder", buf.Bytes())
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
	testutil.AssertGolden(t, "table_value_equals_placeholder_color", buf.Bytes())
}

func TestGolden_TableWideNumber(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
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
	testutil.AssertGolden(t, "table_zero_width", buf.Bytes())
}
