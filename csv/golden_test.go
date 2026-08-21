package csv

import (
	"bytes"
	"strings"
	"testing"

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

func TestGolden_TableBandBlank(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"Name", "A", "B", "Total"}),
		WithFooter(func() [][]string {
			return [][]string{{"Total", ""}}
		}),
		WithPlaceholder("-"),
	)
	if err := tb.Render([][]any{
		{"x", 1, 2, 3},
		{"y", nil, 5},
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
	for _, row := range [][]any{{"x", 1, 2, 3}, {"y", nil, 5}} {
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

func TestGolden_TableDelimiterIndex(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithDelimiter(';'),
		WithIndex(),
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]any{{"x", "y"}, {"p", "q"}}); err != nil {
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
	for _, r := range [][]any{{"x", "y"}, {"p", "q"}} {
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
	if err := tb.Render([][]any{{"alpha", "x・y"}, {"beta", "a,b"}}); err != nil {
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
	for _, r := range [][]any{{"alpha", "x・y"}, {"beta", "a,b"}} {
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
	if err := tb.Render([][]any{{"x", nil}, {nil, "q"}}); err != nil {
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
	for _, r := range [][]any{{"x", nil}, {nil, "q"}} {
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
		WithTransformer(Columns(1), func(any) string {
			return `a;b"c`
		}),
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]any{{"x", "y"}, {"p", "q"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_delimiter_transformer", buf.Bytes())
}

func TestGolden_StreamDelimiterTransformer(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithDelimiter(';'),
		WithTransformer(Columns(1), func(any) string {
			return `a;b"c`
		}),
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
	testutil.AssertGolden(t, "common_delimiter_transformer", buf.Bytes())
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

func TestGolden_TableFooterDelimiter(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithFooter(func() [][]string {
			return [][]string{{"t", "u"}}
		}),
		WithDelimiter(';'),
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]any{{"x", "y"}}); err != nil {
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
	for _, r := range [][]any{{"x", "y"}} {
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
		WithTransformer(Columns(1), func(v any) string {
			if s, ok := v.(string); ok && s == "raw" {
				return "T"
			}
			return ""
		}),
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]any{{"x", "raw"}, {"p", "q"}}); err != nil {
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
		WithTransformer(Columns(1), func(v any) string {
			if s, ok := v.(string); ok && s == "raw" {
				return "T"
			}
			return ""
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
	testutil.AssertGolden(t, "common_index_placeholder", buf.Bytes())
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
	testutil.AssertGolden(t, "common_index_placeholder", buf.Bytes())
}

func TestGolden_TableIndexTransformer(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithIndex(),
		WithTransformer(Columns(1), func(v any) string {
			if s, ok := v.(string); ok && s == "raw" {
				return "T"
			}
			return ""
		}),
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]any{{"x", "raw"}, {"p", "q"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_index_transformer", buf.Bytes())
}

func TestGolden_StreamIndexTransformer(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithIndex(),
		WithTransformer(Columns(1), func(v any) string {
			if s, ok := v.(string); ok && s == "raw" {
				return "T"
			}
			return ""
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

func TestGolden_TablePlaceholderTransformer(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithPlaceholder("-"),
		WithTransformer(Columns(1), func(v any) string {
			if s, ok := v.(string); ok && s == "raw" {
				return "T"
			}
			return ""
		}),
		WithHeader([]string{"A", "B"}),
	)
	if err := tb.Render([][]any{{"x", nil}, {"p", "raw"}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertGolden(t, "common_placeholder_transformer", buf.Bytes())
}

func TestGolden_StreamPlaceholderTransformer(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithPlaceholder("-"),
		WithTransformer(Columns(1), func(v any) string {
			if s, ok := v.(string); ok && s == "raw" {
				return "T"
			}
			return ""
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
	testutil.AssertGolden(t, "common_placeholder_transformer", buf.Bytes())
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
	testutil.AssertGolden(t, "common_pointer", buf.Bytes())
}

func TestGolden_StreamPointer(t *testing.T) {
	var buf bytes.Buffer
	st := testutil.Stringer{Value: "y"}
	str := "alive"
	s := NewStream(&buf,
		WithHeader([]string{"a", "b"}),
	)
	if err := s.Render([]any{&str, &st}); err != nil {
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
	if err := tb.Render([][]any{
		{"alice", "hello, world"},
		{"bob", "no comma"},
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
	if err := s.Render([]any{"alice", "hello, world"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"bob", "no comma"}); err != nil {
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
	if err := tb.Render([][]any{
		{"alice", `say "hello"`},
		{"bob", "no quote"},
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
	if err := s.Render([]any{"alice", `say "hello"`}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"bob", "no quote"}); err != nil {
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
	if err := tb.Render([][]any{
		{"alice", "has, comma\nand newline"},
		{"bob", `"quoted"`},
		{"carol", "plain"},
		{"dave", " leading"},
		{"erin", `\.`},
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
	if err := s.Render([]any{"alice", "has, comma\nand newline"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"bob", `"quoted"`}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"carol", "plain"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"dave", " leading"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"erin", `\.`}); err != nil {
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
	if err := tb.Render([][]any{
		{"alice", "line1\nline2"},
		{"bob", "single"},
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
	if err := s.Render([]any{"alice", "line1\nline2"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"bob", "single"}); err != nil {
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

func TestGolden_TableStringerError(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"Stringer", "Error"}),
	)
	if err := tb.Render([][]any{{testutil.Stringer{Value: "x"}, testutil.Error{Value: "boom"}}}); err != nil {
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

func TestGolden_TableTransformer(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithTransformer(Columns(1), func(v any) string {
			n, ok := v.(int)
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
		WithTransformer(Columns(1), func(v any) string {
			n, ok := v.(int)
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

func TestGolden_TableTypeIndirect(t *testing.T) {
	var buf bytes.Buffer
	stringer := &testutil.PtrStringer{Value: "from Stringer"}
	err := &testutil.PtrError{Value: "from error"}
	tb := NewTable(&buf,
		WithHeader([]string{"Stringer", "Error"}),
	)
	if err := tb.Render([][]any{{&stringer, &err}}); err != nil {
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
	if err := s.Render([]any{&stringer, &err}); err != nil {
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

func TestGolden_TableTypeNamed(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"bool", "uint", "float32", "float64", "uintptr", "complex"}),
	)
	if err := tb.Render([][]any{
		{flag(true), size(7), ratio(1.5), amount(2.25), uintptr(42), complex(1, 2)},
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
	if err := s.Render([]any{flag(true), size(7), ratio(1.5), amount(2.25), uintptr(42), complex(1, 2)}); err != nil {
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
	if err := tb.Render([][]any{{nilStringer, nilError}}); err != nil {
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
	testutil.AssertGolden(t, "common_wide_number", buf.Bytes())
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

func TestGolden_StreamCSV(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithDelimiter(','),
		WithHeader([]string{"Name", "Score"}),
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
	testutil.AssertGolden(t, "stream_csv", buf.Bytes())
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

func TestGolden_StreamQuote(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf,
		WithDelimiter(','),
		WithHeader([]string{"Key", "Value"}),
	)
	if err := s.Render([]any{"comma", "a, b"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"newline", "a\nb"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"quote", `say "hi"`}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"plain", "ok"}); err != nil {
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
	if err := s.Render([]any{"crlf", "a\r\nb"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"cr", "a\rb"}); err != nil {
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
	if err := s.Render([]any{"tab-value", "a\tb"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"plain", "ok"}); err != nil {
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
	if err := s.Render([]any{"foo", 1}); err != nil {
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

func TestGolden_TableCSV(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithDelimiter(','),
		WithHeader([]string{"Name", "Score", "Note"}),
	)
	if err := tb.Render([][]any{
		{"alice", 100, "good"},
		{"bob", 200, "needs work"},
		{"carol", 300, "excellent"},
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
	if err := tb.Render([][]any{
		{"alice", 100},
		{"bob", 99},
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

func TestGolden_TableQuoteCRLF(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"Key", "Value"}),
	)
	if err := tb.Render([][]any{
		{"crlf", "line1\r\nline2"},
		{"cr-only", "a\rb"},
		{"lf-only", "a\nb"},
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
	if err := tb.Render([][]any{
		{"embedded-tab", "a\tb"},
		{"no-tab", "plain"},
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
