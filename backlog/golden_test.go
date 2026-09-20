package backlog

import (
	"bytes"
	"strings"
	"testing"

	"github.com/nekrassov01/table"
	"github.com/nekrassov01/table/internal/param"
	"github.com/nekrassov01/table/internal/testutil"
)

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
	testutil.AssertGolden(t, "stream_all_placeholder", buf.Bytes())
}

func TestGolden_StreamBandBlank(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"Name", "A", "B", "Total"}),
		WithFooter(func() [][]string {
			return [][]string{{"Total", ""}}
		}),
		WithPlaceholder("-"),
	)
	for _, row := range [][]table.Value{{table.String("x"), table.Int(1), table.Int(2), table.Int(3)}, {table.String("y"), table.Any(nil), table.Int(5)}} {
		if err := s.Render(row); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_band_blank", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_basic", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_bold", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_cjk", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_code", buf.Bytes())
}

func TestGolden_StreamCodeDecoOwnMarker(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
		WithColor(ScopeBody, Columns(0, 1), ColorFgRed),
		WithDecoration(ScopeBody, Columns(0), NewDecoration("{code}", "{/code}")),
		WithDecoration(ScopeBody, Columns(1), NewDecoration("''", "''")),
	)
	if err := s.Render([]table.Value{table.String("x"), table.String("y")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_code_deco_own_marker", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_color", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_color_bg", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_color_code", buf.Bytes())
}

func TestGolden_StreamColorEscape(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithColor(ScopeBody, Columns(1), NewColor("red", "")),
		WithHeader([]string{"Key", "Value"}),
	)
	if err := s.Render([]table.Value{table.String("space"), table.String("a b")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("pipe"), table.String("a|b")}); err != nil {
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
	if err := s.Render([]table.Value{table.String("foo"), table.String("hello")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("bar"), table.String("world")}); err != nil {
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
	testutil.AssertGolden(t, "stream_color_nil", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_color_placeholder", buf.Bytes())
}

func TestGolden_StreamColorRejected(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"A", "B"}),
		WithColor(ScopeBody, Columns(0), NewColor(`red){injected}&color(blue`, "")),
		WithColor(ScopeBody, Columns(1), ColorFgRed),
	)
	if err := s.Render([]table.Value{table.String("x"), table.String("y")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_color_rejected", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_color_scope", buf.Bytes())
}

func TestGolden_StreamColspan(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1, 2)),
		WithHeader([]string{"A", "B", "C"}),
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
	testutil.AssertGolden(t, "stream_colspan", buf.Bytes())
}

func TestGolden_StreamColspanColor(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1, 2)),
		WithColor(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), ColorFgRed),
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
	testutil.AssertGolden(t, "stream_colspan_color", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_colspan_decoration", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_colspan_edges", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_colspan_scope", buf.Bytes())
}

func TestGolden_StreamColspanTransformer(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1, 2)),
		WithTransformer(Columns(1), func(v table.Value) (string, *Color, *Decoration) {
			if s, ok := v.AsAny().(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
		}),
		WithHeader([]string{"A", "B", "C"}),
	)
	for _, r := range [][]table.Value{{table.String("T"), table.String("raw"), table.String("y")}, {table.String("p"), table.String("q"), table.String("q")}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_colspan_transformer", buf.Bytes())
}

func TestGolden_StreamDecoEscape(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithDecoration(ScopeBody, Columns(1), DecorationCode),
		WithHeader([]string{"Key", "Value"}),
	)
	if err := s.Render([]table.Value{table.String("pipe"), table.String("a|b")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("newline"), table.String("line1\nline2")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("space"), table.String("a b")}); err != nil {
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
	if err := s.Render([]table.Value{table.String("id-1"), table.String("alice"), table.String("active")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("id-2"), table.String("bob"), table.String("inactive")}); err != nil {
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
	if err := s.Render([]table.Value{table.String("a"), table.Any(nil)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("b"), table.String("ok")}); err != nil {
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
	testutil.AssertGolden(t, "stream_deco_short_row", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_emoji", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_empty_header_label", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_empty_vs_nil", buf.Bytes())
}

func TestGolden_StreamEscape(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"Key", "Value"}),
	)
	if err := s.Render([]table.Value{table.String("pipe"), table.String("a|b")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("newline"), table.String("a\nb")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("a|b"), table.String("c\\d\ne")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_escape", buf.Bytes())
}

func TestGolden_StreamEscapeNotation(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"Notation", "Value"}),
	)
	for _, row := range [][]table.Value{
		{table.String("link"), table.String("[[BLG-87]]")},
		{table.String("bold"), table.String("''bold''")},
		{table.String("italic"), table.String("'''italic'''")},
		{table.String("strikethrough"), table.String("%%strike%%")},
		{table.String("color"), table.String("&color(red){value}")},
		{table.String("quote"), table.String("{quote}value{/quote}")},
		{table.String("code"), table.String("{code}value{/code}")},
		{table.String("typed code"), table.String("{code:go}value{/code}")},
		{table.String("attachment"), table.String("#attach(sample.zip:11)")},
		{table.String("image"), table.String("#image(11)")},
		{table.String("thumbnail"), table.String("#thumbnail(11)")},
		{table.String("revision"), table.String("#rev(11)")},
		{table.String("contents"), table.String("#contents")},
		{table.String("line break"), table.String("&br;")},
	} {
		if err := s.Render(row); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_escape_notation", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_escape_crlf", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_escape_space", buf.Bytes())
}

func TestGolden_StreamFooterColspan(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithFooter(func() [][]string {
			return [][]string{{"t", "t", "u"}}
		}),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1, 2)),
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
	testutil.AssertGolden(t, "stream_footer_colspan", buf.Bytes())
}

func TestGolden_StreamFooterRowspan(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithFooter(func() [][]string {
			return [][]string{{"t", "t"}}
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

func TestGolden_StreamFooterTransformer(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithFooter(func() [][]string {
			return [][]string{{"t", "raw"}}
		}),
		WithTransformer(Columns(1), func(v table.Value) (string, *Color, *Decoration) {
			if s, ok := v.AsAny().(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
		}),
		WithHeader([]string{"A", "B"}),
	)
	for _, r := range [][]table.Value{{table.String("x"), table.String("raw")}, {table.String("p"), table.String("q")}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_footer_transformer", buf.Bytes())
}

func TestGolden_StreamHeaderRowspan(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
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
	testutil.AssertGolden(t, "stream_header_rowspan", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_header_wider_than_rows", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_invalid_utf8", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_italic", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_nil_in_numeric", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_placeholder", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_placeholder_colspan", buf.Bytes())
}

func TestGolden_StreamPlaceholderTransformer(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithPlaceholder("-"),
		WithTransformer(Columns(1), func(v table.Value) (string, *Color, *Decoration) {
			if s, ok := v.AsAny().(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
		}),
		WithHeader([]string{"A", "B"}),
	)
	for _, r := range [][]table.Value{{table.String("x"), table.Any(nil)}, {table.String("p"), table.String("raw")}} {
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
	s := testutil.Stringer{Value: "y"}
	str := "alive"
	var buf bytes.Buffer
	st := NewStream(&buf,
		WithHeader([]string{"a", "b"}),
	)
	if err := st.Render([]table.Value{table.Any(&str), table.Any(&s)}); err != nil {
		t.Fatal(err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_pointer", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_ragged_rows", buf.Bytes())
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
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(1)),
		WithTransformer(Columns(1), func(v table.Value) (string, *Color, *Decoration) {
			if s, ok := v.AsAny().(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
		}),
		WithHeader([]string{"A", "B"}),
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
	testutil.AssertGolden(t, "stream_single_cell", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_span_limit", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_special_chars", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_strikethrough", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_stringer_error", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_transformer", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_transformer_column_override", buf.Bytes())
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
	if err := s.Render([]table.Value{table.Int(int(-1)), table.Int8(int8(-2)), table.Int16(int16(-3)), table.Int32(int32(-4)), table.Int64(int64(-5)),
		table.Uint(uint(1)), table.Uint8(uint8(2)), table.Uint16(uint16(3)), table.Uint32(uint32(4)), table.Uint64(uint64(5))}); err != nil {
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
	testutil.AssertGolden(t, "stream_typed_nil", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_value_equals_placeholder", buf.Bytes())
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
	testutil.AssertGolden(t, "stream_value_equals_placeholder_color", buf.Bytes())
}

func TestGolden_StreamWideNumber(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
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
	testutil.AssertGolden(t, "stream_zero_width", buf.Bytes())
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
	testutil.AssertGolden(t, "table_all_placeholder", buf.Bytes())
}

func TestGolden_TableBandBlank(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"Name", "A", "B", "Total"}),
		WithFooter(func() [][]string {
			return [][]string{{"Total", ""}}
		}),
		WithPlaceholder("-"),
	)
	if err := tb.Render([][]table.Value{
		{table.String("x"), table.Int(1), table.Int(2), table.Int(3)},
		{table.String("y"), table.Any(nil), table.Int(5)},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_band_blank", buf.Bytes())
}

func TestGolden_TableBandBlankColored(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"Name", "A"}),
		WithFooter(func() [][]string {
			return [][]string{{"Total", ""}}
		}),
		WithColor(ScopeFooter, Columns(1), ColorFgRed),
	)
	if err := tb.Render([][]table.Value{{table.String("x"), table.Int(1)}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_band_blank_colored", buf.Bytes())
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
	testutil.AssertGolden(t, "table_basic", buf.Bytes())
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
	testutil.AssertGolden(t, "table_bold", buf.Bytes())
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
	testutil.AssertGolden(t, "table_cjk", buf.Bytes())
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
	testutil.AssertGolden(t, "table_code", buf.Bytes())
}

func TestGolden_TableCodeDecoOwnMarker(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithColor(ScopeBody, Columns(0, 1), ColorFgRed),
		WithDecoration(ScopeBody, Columns(0), NewDecoration("{code}", "{/code}")),
		WithDecoration(ScopeBody, Columns(1), NewDecoration("''", "''")),
	)
	if err := tb.Render([][]table.Value{
		{table.String("x"), table.String("y")},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_code_deco_own_marker", buf.Bytes())
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
	testutil.AssertGolden(t, "table_color", buf.Bytes())
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
	testutil.AssertGolden(t, "table_color_bg", buf.Bytes())
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
	testutil.AssertGolden(t, "table_color_code", buf.Bytes())
}

func TestGolden_TableColorEscape(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithColor(ScopeBody, Columns(1), NewColor("red", "")),
		WithHeader([]string{"Key", "Value"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("space"), table.String("a b")},
		{table.String("pipe"), table.String("a|b")},
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
	if err := tb.Render([][]table.Value{
		{table.String("foo"), table.String("hello")},
		{table.String("bar"), table.String("world")},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_color_fg_bg", buf.Bytes())
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
	testutil.AssertGolden(t, "table_color_nil", buf.Bytes())
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
	testutil.AssertGolden(t, "table_color_placeholder", buf.Bytes())
}

func TestGolden_TableColorRejected(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithColor(ScopeBody, Columns(0), NewColor(`red){injected}&color(blue`, "")),
		WithColor(ScopeBody, Columns(1), ColorFgRed),
	)
	if err := tb.Render([][]table.Value{
		{table.String("x"), table.String("y")},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_color_rejected", buf.Bytes())
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
	testutil.AssertGolden(t, "table_color_scope", buf.Bytes())
}

func TestGolden_TableColspan(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1, 2)),
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("x"), table.String("x"), table.String("y")},
		{table.String("p"), table.String("q"), table.String("q")},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_colspan", buf.Bytes())
}

func TestGolden_TableColspanColor(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1, 2)),
		WithColor(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), ColorFgRed),
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := tb.Render([][]table.Value{{table.String("x"), table.String("x"), table.String("y")}, {table.String("p"), table.String("q"), table.String("q")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_colspan_color", buf.Bytes())
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
	testutil.AssertGolden(t, "table_colspan_decoration", buf.Bytes())
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
	testutil.AssertGolden(t, "table_colspan_edges", buf.Bytes())
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
	testutil.AssertGolden(t, "table_colspan_scope", buf.Bytes())
}

func TestGolden_TableColspanTransformer(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1, 2)),
		WithTransformer(Columns(1), func(v table.Value) (string, *Color, *Decoration) {
			if s, ok := v.AsAny().(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
		}),
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := tb.Render([][]table.Value{{table.String("T"), table.String("raw"), table.String("y")}, {table.String("p"), table.String("q"), table.String("q")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_colspan_transformer", buf.Bytes())
}

func TestGolden_TableDecoEscape(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithDecoration(ScopeBody, Columns(1), DecorationCode),
		WithHeader([]string{"Key", "Value"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("pipe"), table.String("a|b")},
		{table.String("newline"), table.String("line1\nline2")},
		{table.String("space"), table.String("a b")},
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
	if err := tb.Render([][]table.Value{
		{table.String("id-1"), table.String("alice"), table.String("active")},
		{table.String("id-2"), table.String("bob"), table.String("inactive")},
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
	if err := tb.Render([][]table.Value{
		{table.String("a"), table.Any(nil)},
		{table.String("b"), table.String("")},
		{table.String("c"), table.String("ok")},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_deco_nil", buf.Bytes())
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
	testutil.AssertGolden(t, "table_deco_short_row", buf.Bytes())
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
	testutil.AssertGolden(t, "table_emoji", buf.Bytes())
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
	testutil.AssertGolden(t, "table_empty_header_label", buf.Bytes())
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
	testutil.AssertGolden(t, "table_empty_vs_nil", buf.Bytes())
}

func TestGolden_TableEscape(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"Key", "Value"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("pipe"), table.String("a|b")},
		{table.String("backslash"), table.String("a\\b")},
		{table.String("newline"), table.String("a\nb")},
		{table.String("combined"), table.String("a|b\nc\\d")},
		{table.String("a|b"), table.String("c\\d\ne")},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_escape", buf.Bytes())
}

func TestGolden_TableEscapeNotation(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"Notation", "Value"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("link"), table.String("[[BLG-87]]")},
		{table.String("bold"), table.String("''bold''")},
		{table.String("italic"), table.String("'''italic'''")},
		{table.String("strikethrough"), table.String("%%strike%%")},
		{table.String("color"), table.String("&color(red){value}")},
		{table.String("quote"), table.String("{quote}value{/quote}")},
		{table.String("code"), table.String("{code}value{/code}")},
		{table.String("typed code"), table.String("{code:go}value{/code}")},
		{table.String("attachment"), table.String("#attach(sample.zip:11)")},
		{table.String("image"), table.String("#image(11)")},
		{table.String("thumbnail"), table.String("#thumbnail(11)")},
		{table.String("revision"), table.String("#rev(11)")},
		{table.String("contents"), table.String("#contents")},
		{table.String("line break"), table.String("&br;")},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_escape_notation", buf.Bytes())
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
	testutil.AssertGolden(t, "table_escape_crlf", buf.Bytes())
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
	testutil.AssertGolden(t, "table_escape_space", buf.Bytes())
}

func TestGolden_TableFooterColspan(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithFooter(func() [][]string {
			return [][]string{{"t", "t", "u"}}
		}),
		WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1, 2)),
		WithHeader([]string{"A", "B", "C"}),
	)
	if err := tb.Render([][]table.Value{{table.String("x"), table.String("x"), table.String("y")}, {table.String("p"), table.String("q"), table.String("q")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_footer_colspan", buf.Bytes())
}

func TestGolden_TableFooterRowspan(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithFooter(func() [][]string {
			return [][]string{{"t", "t"}}
		}),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]table.Value{{table.String("g"), table.String("x")}, {table.String("g"), table.String("y")}, {table.String("h"), table.String("z")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_footer_rowspan", buf.Bytes())
}

func TestGolden_TableFooterTransformer(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithFooter(func() [][]string {
			return [][]string{{"t", "raw"}}
		}),
		WithTransformer(Columns(1), func(v table.Value) (string, *Color, *Decoration) {
			if s, ok := v.AsAny().(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
		}),
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]table.Value{{table.String("x"), table.String("raw")}, {table.String("p"), table.String("q")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_footer_transformer", buf.Bytes())
}

func TestGolden_TableHeaderRowspan(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
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
	testutil.AssertGolden(t, "table_header_rowspan", buf.Bytes())
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
	testutil.AssertGolden(t, "table_header_wider_than_rows", buf.Bytes())
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
	testutil.AssertGolden(t, "table_invalid_utf8", buf.Bytes())
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
	testutil.AssertGolden(t, "table_italic", buf.Bytes())
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
	testutil.AssertGolden(t, "table_nil_in_numeric", buf.Bytes())
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
	testutil.AssertGolden(t, "table_placeholder", buf.Bytes())
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
	testutil.AssertGolden(t, "table_placeholder_colspan", buf.Bytes())
}

func TestGolden_TablePlaceholderTransformer(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithPlaceholder("-"),
		WithTransformer(Columns(1), func(v table.Value) (string, *Color, *Decoration) {
			if s, ok := v.AsAny().(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
		}),
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]table.Value{{table.String("x"), table.Any(nil)}, {table.String("p"), table.String("raw")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_placeholder_transformer", buf.Bytes())
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
	testutil.AssertGolden(t, "table_pointer", buf.Bytes())
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
	testutil.AssertGolden(t, "table_ragged_rows", buf.Bytes())
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
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(1)),
		WithTransformer(Columns(1), func(v table.Value) (string, *Color, *Decoration) {
			if s, ok := v.AsAny().(string); ok && s == "raw" {
				return "T", nil, nil
			}
			return "", nil, nil
		}),
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]table.Value{{table.String("x"), table.String("raw")}, {table.String("p"), table.String("raw")}, {table.String("q"), table.String("z")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_rowspan_transformer", buf.Bytes())
}

func TestGolden_TableSingleCell(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"only"}),
	)
	if err := tb.Render([][]table.Value{{table.String("v")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_single_cell", buf.Bytes())
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
	testutil.AssertGolden(t, "table_span_limit", buf.Bytes())
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
	testutil.AssertGolden(t, "table_special_chars", buf.Bytes())
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
	testutil.AssertGolden(t, "table_strikethrough", buf.Bytes())
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
	testutil.AssertGolden(t, "table_stringer_error", buf.Bytes())
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
	testutil.AssertGolden(t, "table_transformer", buf.Bytes())
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
	testutil.AssertGolden(t, "table_transformer_column_override", buf.Bytes())
}

func TestGolden_TableTypeFloat(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"f32", "f64"}),
	)
	if err := tb.Render([][]table.Value{{table.Float32(float32(3.14)), table.Float64(float64(2.71828))}}); err != nil {
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
	if err := tb.Render([][]table.Value{
		{table.Int(int(-1)), table.Int8(int8(-2)), table.Int16(int16(-3)), table.Int32(int32(-4)), table.Int64(int64(-5)),
			table.Uint(uint(1)), table.Uint8(uint8(2)), table.Uint16(uint16(3)), table.Uint32(uint32(4)), table.Uint64(uint64(5))},
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
	testutil.AssertGolden(t, "table_typed_nil", buf.Bytes())
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
	testutil.AssertGolden(t, "table_value_equals_placeholder", buf.Bytes())
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
	testutil.AssertGolden(t, "table_value_equals_placeholder_color", buf.Bytes())
}

func TestGolden_TableWideNumber(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
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
	testutil.AssertGolden(t, "table_zero_width", buf.Bytes())
}
