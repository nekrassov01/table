package csv

import (
	"bytes"
	"strings"
	"testing"

	"github.com/nekrassov01/table"
	"github.com/nekrassov01/table/internal/testutil"
)

type (
	flag   bool
	size   uint
	ratio  float32
	amount float64
)

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
	testutil.AssertGolden(t, "common_band_blank", buf.Bytes())
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
	testutil.AssertGolden(t, "common_band_blank", buf.Bytes())
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

func TestGolden_TableCRLFQuote(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithCRLF(),
		WithHeader([]string{"Key", "Value"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("crlf"), table.String("a\r\nb")},
		{table.String("cr"), table.String("a\rb")},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_crlf_quote", buf.Bytes())
}

func TestGolden_StreamCRLFQuote(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithCRLF(),
		WithHeader([]string{"Key", "Value"}),
	)
	if err := s.Render([]table.Value{table.String("crlf"), table.String("a\r\nb")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("cr"), table.String("a\rb")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_crlf_quote", buf.Bytes())
}

func TestGolden_TableDelimiterIndex(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithDelimiter(';'),
		WithIndex(),
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]table.Value{{table.String("x"), table.String("y")}, {table.String("p"), table.String("q")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_delimiter_index", buf.Bytes())
}

func TestGolden_StreamDelimiterIndex(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithDelimiter(';'),
		WithIndex(),
		WithHeader([]string{"A", "B"}),
	)
	for _, r := range [][]table.Value{{table.String("x"), table.String("y")}, {table.String("p"), table.String("q")}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_delimiter_index", buf.Bytes())
}

func TestGolden_TableDelimiterUnicode(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithDelimiter('・'),
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]table.Value{{table.String("alpha"), table.String("x・y")}, {table.String("beta"), table.String("a,b")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_delimiter_unicode", buf.Bytes())
}

func TestGolden_StreamDelimiterUnicode(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithDelimiter('・'),
		WithHeader([]string{"A", "B"}),
	)
	for _, r := range [][]table.Value{{table.String("alpha"), table.String("x・y")}, {table.String("beta"), table.String("a,b")}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_delimiter_unicode", buf.Bytes())
}

func TestGolden_TableDelimiterPlaceholder(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithDelimiter(';'),
		WithPlaceholder("-"),
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]table.Value{{table.String("x"), table.Any(nil)}, {table.Any(nil), table.String("q")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_delimiter_placeholder", buf.Bytes())
}

func TestGolden_StreamDelimiterPlaceholder(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithDelimiter(';'),
		WithPlaceholder("-"),
		WithHeader([]string{"A", "B"}),
	)
	for _, r := range [][]table.Value{{table.String("x"), table.Any(nil)}, {table.Any(nil), table.String("q")}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_delimiter_placeholder", buf.Bytes())
}

func TestGolden_TableDelimiterTransformer(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithDelimiter(';'),
		WithTransformer(Columns(1), func(table.Value) string {
			return `a;b"c`
		}),
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]table.Value{{table.String("x"), table.String("y")}, {table.String("p"), table.String("q")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_delimiter_transformer", buf.Bytes())
}

func TestGolden_StreamDelimiterTransformer(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithDelimiter(';'),
		WithTransformer(Columns(1), func(table.Value) string {
			return `a;b"c`
		}),
		WithHeader([]string{"A", "B"}),
	)
	for _, r := range [][]table.Value{{table.String("x"), table.String("y")}, {table.String("p"), table.String("q")}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_delimiter_transformer", buf.Bytes())
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

func TestGolden_TableFooterDelimiter(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithFooter(func() [][]string {
			return [][]string{{"t", "u"}}
		}),
		WithDelimiter(';'),
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]table.Value{{table.String("x"), table.String("y")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_footer_delimiter", buf.Bytes())
}

func TestGolden_StreamFooterDelimiter(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithFooter(func() [][]string {
			return [][]string{{"t", "u"}}
		}),
		WithDelimiter(';'),
		WithHeader([]string{"A", "B"}),
	)
	for _, r := range [][]table.Value{{table.String("x"), table.String("y")}} {
		if err := s.Render(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_footer_delimiter", buf.Bytes())
}

func TestGolden_TableFooterEmptyBody(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A", "B"}),
		WithFooter(func() [][]string {
			return [][]string{{"x", "y"}}
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
			return [][]string{{"x", "y"}}
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

func TestGolden_TableFooterTransformer(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithFooter(func() [][]string {
			return [][]string{{"t", "raw"}}
		}),
		WithTransformer(Columns(1), func(v table.Value) string {
			if s, ok := v.AsAny().(string); ok && s == "raw" {
				return "T"
			}
			return ""
		}),
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]table.Value{{table.String("x"), table.String("raw")}, {table.String("p"), table.String("q")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_footer_transformer", buf.Bytes())
}

func TestGolden_StreamFooterTransformer(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithFooter(func() [][]string {
			return [][]string{{"t", "raw"}}
		}),
		WithTransformer(Columns(1), func(v table.Value) string {
			if s, ok := v.AsAny().(string); ok && s == "raw" {
				return "T"
			}
			return ""
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

func TestGolden_TableIndexPlaceholder(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithIndex(),
		WithPlaceholder("-"),
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]table.Value{{table.String("x"), table.Any(nil)}, {table.Any(nil), table.String("q")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_index_placeholder", buf.Bytes())
}

func TestGolden_StreamIndexPlaceholder(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithIndex(),
		WithPlaceholder("-"),
		WithHeader([]string{"A", "B"}),
	)
	for _, r := range [][]table.Value{{table.String("x"), table.Any(nil)}, {table.Any(nil), table.String("q")}} {
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
		WithIndex(),
		WithTransformer(Columns(1), func(v table.Value) string {
			if s, ok := v.AsAny().(string); ok && s == "raw" {
				return "T"
			}
			return ""
		}),
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]table.Value{{table.String("x"), table.String("raw")}, {table.String("p"), table.String("q")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_index_transformer", buf.Bytes())
}

func TestGolden_StreamIndexTransformer(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithIndex(),
		WithTransformer(Columns(1), func(v table.Value) string {
			if s, ok := v.AsAny().(string); ok && s == "raw" {
				return "T"
			}
			return ""
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
	testutil.AssertGolden(t, "common_index_transformer", buf.Bytes())
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

func TestGolden_TablePlaceholderTransformer(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithPlaceholder("-"),
		WithTransformer(Columns(1), func(v table.Value) string {
			if s, ok := v.AsAny().(string); ok && s == "raw" {
				return "T"
			}
			return ""
		}),
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]table.Value{{table.String("x"), table.Any(nil)}, {table.String("p"), table.String("raw")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_placeholder_transformer", buf.Bytes())
}

func TestGolden_StreamPlaceholderTransformer(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithPlaceholder("-"),
		WithTransformer(Columns(1), func(v table.Value) string {
			if s, ok := v.AsAny().(string); ok && s == "raw" {
				return "T"
			}
			return ""
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
	testutil.AssertGolden(t, "common_placeholder_transformer", buf.Bytes())
}

func TestGolden_TablePointer(t *testing.T) {
	var buf bytes.Buffer
	s := testutil.Stringer{Value: "y"}
	str := "alive"
	tb := NewTable(&buf,
		WithHeader([]string{"a", "b"}),
	)
	if err := tb.Render([][]table.Value{{table.Any(&str), table.Any(&s)}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_pointer", buf.Bytes())
}

func TestGolden_StreamPointer(t *testing.T) {
	var buf bytes.Buffer
	st := testutil.Stringer{Value: "y"}
	str := "alive"
	s := NewStream(&buf,
		WithHeader([]string{"a", "b"}),
	)
	if err := s.Render([]table.Value{table.Any(&str), table.Any(&st)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_pointer", buf.Bytes())
}

func TestGolden_TableQuoteComma(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithDelimiter(','),
		WithHeader([]string{"Name", "Note"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("alice"), table.String("hello, world")},
		{table.String("bob"), table.String("no comma")},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_quote_comma", buf.Bytes())
}

func TestGolden_StreamQuoteComma(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithDelimiter(','),
		WithHeader([]string{"Name", "Note"}),
	)
	if err := s.Render([]table.Value{table.String("alice"), table.String("hello, world")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("bob"), table.String("no comma")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_quote_comma", buf.Bytes())
}

func TestGolden_TableQuoteDoubleQuote(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"Name", "Note"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("alice"), table.String(`say "hello"`)},
		{table.String("bob"), table.String("no quote")},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_quote_double_quote", buf.Bytes())
}

func TestGolden_StreamQuoteDoubleQuote(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"Name", "Note"}),
	)
	if err := s.Render([]table.Value{table.String("alice"), table.String(`say "hello"`)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("bob"), table.String("no quote")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_quote_double_quote", buf.Bytes())
}

func TestGolden_TableQuoteMixed(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithDelimiter(','),
		WithHeader([]string{"Name", "Note"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("alice"), table.String("has, comma\nand newline")},
		{table.String("bob"), table.String(`"quoted"`)},
		{table.String("carol"), table.String("plain")},
		{table.String("dave"), table.String(" leading")},
		{table.String("erin"), table.String(`\.`)},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_quote_mixed", buf.Bytes())
}

func TestGolden_StreamQuoteMixed(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithDelimiter(','),
		WithHeader([]string{"Name", "Note"}),
	)
	if err := s.Render([]table.Value{table.String("alice"), table.String("has, comma\nand newline")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("bob"), table.String(`"quoted"`)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("carol"), table.String("plain")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("dave"), table.String(" leading")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("erin"), table.String(`\.`)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_quote_mixed", buf.Bytes())
}

func TestGolden_TableQuoteNewline(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"Name", "Note"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("alice"), table.String("line1\nline2")},
		{table.String("bob"), table.String("single")},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_quote_newline", buf.Bytes())
}

func TestGolden_StreamQuoteNewline(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"Name", "Note"}),
	)
	if err := s.Render([]table.Value{table.String("alice"), table.String("line1\nline2")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("bob"), table.String("single")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_quote_newline", buf.Bytes())
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

func TestGolden_TableStringerError(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"Stringer", "Error"}),
	)
	if err := tb.Render([][]table.Value{{table.Any(testutil.Stringer{Value: "x"}), table.Any(testutil.Error{Value: "boom"})}}); err != nil {
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

func TestGolden_TableTransformer(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithTransformer(Columns(1), func(v table.Value) string {
			n, ok := v.AsAny().(int)
			if !ok {
				return ""
			}
			if n >= 100 {
				return "high"
			}
			return ""
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
		WithTransformer(Columns(1), func(v table.Value) string {
			n, ok := v.AsAny().(int)
			if !ok {
				return ""
			}
			if n >= 100 {
				return "high"
			}
			return ""
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

func TestGolden_TableTypeIndirect(t *testing.T) {
	var buf bytes.Buffer
	stringer := &testutil.PtrStringer{Value: "from Stringer"}
	err := &testutil.PtrError{Value: "from error"}
	tb := NewTable(&buf,
		WithHeader([]string{"Stringer", "Error"}),
	)
	if err := tb.Render([][]table.Value{{table.Any(&stringer), table.Any(&err)}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_type_indirect", buf.Bytes())
}

func TestGolden_StreamTypeIndirect(t *testing.T) {
	var buf bytes.Buffer
	stringer := &testutil.PtrStringer{Value: "from Stringer"}
	err := &testutil.PtrError{Value: "from error"}
	s := NewStream(&buf,
		WithHeader([]string{"Stringer", "Error"}),
	)
	if err := s.Render([]table.Value{table.Any(&stringer), table.Any(&err)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_type_indirect", buf.Bytes())
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

func TestGolden_TableTypeNamed(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"bool", "uint", "float32", "float64", "uintptr", "complex"}),
	)
	if err := tb.Render([][]table.Value{
		{table.Any(flag(true)), table.Any(size(7)), table.Any(ratio(1.5)), table.Any(amount(2.25)), table.Uintptr(uintptr(42)), table.Any(complex(1, 2))},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_type_named", buf.Bytes())
}

func TestGolden_StreamTypeNamed(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"bool", "uint", "float32", "float64", "uintptr", "complex"}),
	)
	if err := s.Render([]table.Value{table.Any(flag(true)), table.Any(size(7)), table.Any(ratio(1.5)), table.Any(amount(2.25)), table.Uintptr(uintptr(42)), table.Any(complex(1, 2))}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_type_named", buf.Bytes())
}

func TestGolden_TableTypedNil(t *testing.T) {
	var buf bytes.Buffer
	var nilStringer *testutil.PtrStringer
	var nilError *testutil.PtrError
	tb := NewTable(&buf,
		WithPlaceholder("<nil>"),
		WithHeader([]string{"Stringer", "Error"}),
	)
	if err := tb.Render([][]table.Value{{table.Any(nilStringer), table.Any(nilError)}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_typed_nil", buf.Bytes())
}

func TestGolden_StreamTypedNil(t *testing.T) {
	var buf bytes.Buffer
	var nilStringer *testutil.PtrStringer
	var nilError *testutil.PtrError
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
	testutil.AssertGolden(t, "common_wide_number", buf.Bytes())
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

func TestGolden_StreamCSV(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithDelimiter(','),
		WithHeader([]string{"Name", "Score"}),
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
	testutil.AssertGolden(t, "stream_csv", buf.Bytes())
}

func TestGolden_StreamIndex(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithIndex(),
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
	testutil.AssertGolden(t, "stream_index", buf.Bytes())
}

func TestGolden_StreamQuote(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithDelimiter(','),
		WithHeader([]string{"Key", "Value"}),
	)
	if err := s.Render([]table.Value{table.String("comma"), table.String("a, b")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("newline"), table.String("a\nb")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("quote"), table.String(`say "hi"`)}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("plain"), table.String("ok")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_quote", buf.Bytes())
}

func TestGolden_StreamQuoteCRLF(t *testing.T) {
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
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_quote_crlf", buf.Bytes())
}

func TestGolden_StreamQuoteTab(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithHeader([]string{"Key", "Value"}),
	)
	if err := s.Render([]table.Value{table.String("tab-value"), table.String("a\tb")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("plain"), table.String("ok")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "stream_quote_tab", buf.Bytes())
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

func TestGolden_TableCSV(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithDelimiter(','),
		WithHeader([]string{"Name", "Score", "Note"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("alice"), table.Int(100), table.String("good")},
		{table.String("bob"), table.Int(200), table.String("needs work")},
		{table.String("carol"), table.Int(300), table.String("excellent")},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_csv", buf.Bytes())
}

func TestGolden_TableIndex(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithIndex(),
		WithHeader([]string{"Name", "Score"}),
		WithFooter(func() [][]string {
			return [][]string{{"Total", "199"}}
		}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("alice"), table.Int(100)},
		{table.String("bob"), table.Int(99)},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_index", buf.Bytes())
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

func TestGolden_TableQuoteCRLF(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"Key", "Value"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("crlf"), table.String("line1\r\nline2")},
		{table.String("cr-only"), table.String("a\rb")},
		{table.String("lf-only"), table.String("a\nb")},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_quote_crlf", buf.Bytes())
}

func TestGolden_TableQuoteTab(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"Key", "Value"}),
	)
	if err := tb.Render([][]table.Value{
		{table.String("embedded-tab"), table.String("a\tb")},
		{table.String("no-tab"), table.String("plain")},
	}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "table_quote_tab", buf.Bytes())
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
