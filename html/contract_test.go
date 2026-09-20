package html

import (
	"bytes"
	"errors"
	"io"
	"math"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/nekrassov01/table"
	"github.com/nekrassov01/table/internal/testutil"
)

type contractCase struct {
	name          string
	opts          []Option
	header        []string
	rows          [][]table.Value
	omitHeader    bool
	streamDiffers bool
}

func TestContract_TableDeterministic(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	header := []string{"A", "B", "C"}
	rows := [][]table.Value{{table.String("x"), table.Int(1), table.Bool(true)}, {table.String("y"), table.Int(2), table.Bool(false)}}
	t1 := NewTable(&buf1, WithHeader(header))
	if err := t1.Render(rows); err != nil {
		t.Fatal(err)
	}
	t2 := NewTable(&buf2, WithHeader(header))
	if err := t2.Render(rows); err != nil {
		t.Fatal(err)
	}
	testutil.AssertBytes(t, buf2.Bytes(), buf1.Bytes(), "second render")
}

func TestContract_TableReusable(t *testing.T) {
	var buf bytes.Buffer
	header := []string{"A", "B"}
	rows := [][]table.Value{{table.String("x"), table.Int(1)}, {table.String("y"), table.Int(2)}}
	tb := NewTable(&buf, WithHeader(header))
	if err := tb.Render(rows); err != nil {
		t.Fatal(err)
	}
	first := make([]byte, buf.Len())
	copy(first, buf.Bytes())
	buf.Reset()
	if err := tb.Render(rows); err != nil {
		t.Fatal(err)
	}
	testutil.AssertBytes(t, buf.Bytes(), first, "second render")
}

func TestContract_TableFooterTiming(t *testing.T) {
	var buf bytes.Buffer
	calls := 0
	total := "before"
	tb := NewTable(&buf,
		WithHeader([]string{"Item", "Total"}),
		WithFooter(func() [][]string {
			calls++
			return [][]string{{"sum", total}}
		}),
	)
	if calls != 0 {
		t.Fatalf("constructor called footer %d times", calls)
	}
	total = "after"
	if err := tb.Render([][]table.Value{{table.String("a"), table.Int(1)}}); err != nil {
		t.Fatal(err)
	}
	if calls != 1 {
		t.Fatalf("Render called footer %d times, want 1", calls)
	}
	if !strings.Contains(buf.String(), total) {
		t.Fatalf("output does not contain resolved footer %q:\n%s", total, buf.String())
	}
}

func TestContract_TableWriterError(t *testing.T) {
	cause := testutil.NewError()
	tb := NewTable(&testutil.ErrorWriter{
		Err: cause,
	},
		WithHeader([]string{"A"}),
	)
	err := tb.Render([][]table.Value{{table.String("x")}})
	if _, ok := err.(*table.Error); !ok {
		t.Fatalf("expected outer *table.Error, got %T", err)
	}
	if !errors.Is(err, table.ErrWriteFailed) || !errors.Is(err, cause) {
		t.Fatalf("expected table.ErrWriteFailed and writer error, got %v", err)
	}
}

func TestContract_TableNilRows(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf, WithHeader([]string{"A"}))
	if err := tb.Render(nil); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestContract_StreamLifecycle(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf, WithHeader([]string{"A", "B"}))
	for _, row := range [][]table.Value{{table.String("x"), table.Int(1)}, {table.String("y"), table.Int(2)}} {
		if err := s.Render(row); err != nil {
			t.Fatalf("Render: %v", err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	if buf.Len() == 0 {
		t.Fatal("expected output, got empty")
	}
}

func TestContract_StreamRenderAfterClose(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf, WithHeader([]string{"A"}))
	if err := s.Render([]table.Value{table.String("x")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("y")}); !errors.Is(err, table.ErrClosed) {
		t.Fatalf("expected table.ErrClosed, got %v", err)
	}
}

func TestContract_StreamCloseAfterClose(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf, WithHeader([]string{"A"}))
	if err := s.Render([]table.Value{table.String("x")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestContract_StreamFooterTiming(t *testing.T) {
	var buf bytes.Buffer
	calls := 0
	total := "before"
	s := NewStream(&buf,
		WithHeader([]string{"Item", "Total"}),
		WithFooter(func() [][]string {
			calls++
			return [][]string{{"sum", total}}
		}),
	)
	if err := s.Render([]table.Value{table.String("a"), table.Int(1)}); err != nil {
		t.Fatal(err)
	}
	if calls != 0 {
		t.Fatalf("Render called footer %d times", calls)
	}
	total = "after"
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	if calls != 1 {
		t.Fatalf("Close called footer %d times, want 1", calls)
	}
	if !strings.Contains(buf.String(), total) {
		t.Fatalf("output does not contain resolved footer %q:\n%s", total, buf.String())
	}
}

func TestContract_FooterCanWidenTableButNotStream(t *testing.T) {
	footer := func() [][]string {
		return [][]string{{"sum", "footer-only"}}
	}
	rows := [][]table.Value{{table.Int(1)}}
	var tableOutput bytes.Buffer
	if err := NewTable(&tableOutput, WithFooter(footer)).Render(rows); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(tableOutput.String(), "footer-only") {
		t.Fatalf("Table dropped the wider footer:\n%s", tableOutput.String())
	}
	var streamOutput bytes.Buffer
	stream := NewStream(&streamOutput, WithFooter(footer))
	if err := stream.Render(rows[0]); err != nil {
		t.Fatal(err)
	}
	first := stream.Close()
	if !errors.Is(first, table.ErrColumnCount) {
		t.Fatalf("Close: expected table.ErrColumnCount, got %v", first)
	}
	if err := stream.Close(); err != first {
		t.Fatalf("second Close: expected %v, got %v", first, err)
	}
	if !strings.HasSuffix(streamOutput.String(), "  </tbody>\n</table>\n") {
		t.Fatalf("Stream did not close the table after the footer error:\n%s", streamOutput.String())
	}
	if strings.Contains(streamOutput.String(), "footer-only") {
		t.Fatalf("Stream partially wrote the wider footer:\n%s", streamOutput.String())
	}
}

func TestContract_StreamCloseErrorOrder(t *testing.T) {
	cause := testutil.NewError()
	footer := func() [][]string {
		return [][]string{{"sum", "extra"}}
	}
	var footerOutput bytes.Buffer
	footerStream := NewStream(&footerOutput,
		WithHeader([]string{"A"}),
		WithFooter(footer),
	)
	if err := footerStream.Render([]table.Value{table.String("a")}); err != nil {
		t.Fatal(err)
	}
	footerStream.w = &testutil.ErrorWriter{Err: cause}
	footerErr := footerStream.Close()
	if !errors.Is(footerErr, table.ErrColumnCount) || errors.Is(footerErr, table.ErrWriteFailed) {
		t.Fatalf("footer and close errors: expected only table.ErrColumnCount, got %v", footerErr)
	}
	if err := footerStream.Close(); err != footerErr {
		t.Fatalf("second footer Close: expected %v, got %v", footerErr, err)
	}

	var writeOutput bytes.Buffer
	writeStream := NewStream(&writeOutput,
		WithHeader([]string{"A"}),
		WithFooter(func() [][]string {
			return [][]string{{"sum"}}
		}),
	)
	if err := writeStream.Render([]table.Value{table.String("a")}); err != nil {
		t.Fatal(err)
	}
	writeStream.w = &testutil.ErrorWriter{Err: cause}
	writeErr := writeStream.Close()
	if !errors.Is(writeErr, table.ErrWriteFailed) || !errors.Is(writeErr, cause) {
		t.Fatalf("close error: expected table.ErrWriteFailed and writer error, got %v", writeErr)
	}
	if err := writeStream.Close(); err != writeErr {
		t.Fatalf("second writer Close: expected %v, got %v", writeErr, err)
	}
}

func TestContract_StreamWriterError(t *testing.T) {
	cause := testutil.NewError()
	s := NewStream(&testutil.ErrorWriter{
		Err: cause,
	},
		WithHeader([]string{"A"}),
	)
	err := s.Render([]table.Value{table.String("x")})
	if _, ok := err.(*table.Error); !ok {
		t.Fatalf("expected outer *table.Error, got %T", err)
	}
	if !errors.Is(err, table.ErrWriteFailed) || !errors.Is(err, cause) {
		t.Fatalf("expected table.ErrWriteFailed and writer error, got %v", err)
	}
}

func TestContract_StreamCloseWithoutRender(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf, WithHeader([]string{"A"}))
	if err := s.Close(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestContract_StreamWriterErrorSticky(t *testing.T) {
	s := NewStream(&testutil.ErrorWriter{
		Err: testutil.NewError(),
	},
		WithHeader([]string{"A"}),
	)
	first := s.Render([]table.Value{table.String("x")})
	if first == nil {
		t.Fatal("expected error, got nil")
	}
	if err := s.Render([]table.Value{table.String("y")}); err == nil || err != first {
		t.Fatalf("second Render: want the latched error, got %v", err)
	}
	if err := s.Close(); err == nil || err != first {
		t.Fatalf("Close: want the latched error, got %v", err)
	}
}

func TestContract_TableZeroColumns(t *testing.T) {
	var buf bytes.Buffer
	if err := NewTable(&buf).Render([][]table.Value{}); err != nil {
		t.Fatal(err)
	}
	if buf.Len() != 0 {
		t.Fatalf("expected no output, got:\n%s", buf.String())
	}
}

func TestContract_StreamZeroColumns(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf)
	if err := s.Render(nil); err != nil {
		t.Fatal(err)
	}
	if s.arena != nil {
		t.Fatal("zero-column render retained an arena")
	}
	if buf.Len() != 0 {
		t.Fatalf("expected no output, got:\n%s", buf.String())
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestContract_DecorationWithoutPrefix(t *testing.T) {
	if deco := NewDecoration("", "suffix"); deco != nil {
		t.Fatalf("expected nil decoration, got %v", deco)
	}
	rows := [][]table.Value{{table.String("x"), table.String("y")}}
	var dressed, plain bytes.Buffer
	if err := NewTable(&dressed,
		WithHeader([]string{"A", "B"}),
		WithDecoration(ScopeHeader|ScopeBody|ScopeFooter, Columns(0), NewDecoration("", "suffix")),
	).Render(rows); err != nil {
		t.Fatal(err)
	}
	if err := NewTable(&plain, WithHeader([]string{"A", "B"})).Render(rows); err != nil {
		t.Fatal(err)
	}
	testutil.AssertBytes(t, dressed.Bytes(), plain.Bytes(), "table with an empty decoration")
}

func TestContract_TableStreamAgree(t *testing.T) {
	for _, c := range contractCases() {
		if c.streamDiffers {
			continue
		}
		t.Run(c.name, func(t *testing.T) {
			table := renderContractTable(c)
			stream := renderContractStream(c)
			testutil.AssertBytes(t, stream, table, "stream against table")
		})
	}
}

func TestContract_ReuseMatchesFresh(t *testing.T) {
	for _, c := range contractCases() {
		t.Run(c.name, func(t *testing.T) {
			want := renderContractTable(c)
			var buf bytes.Buffer
			opts := c.opts
			if !c.omitHeader {
				opts = append([]Option{WithHeader(c.header)}, opts...)
			}
			tb := NewTable(&buf, opts...)
			for pass := range 3 {
				buf.Reset()
				if err := tb.Render(c.rows); err != nil {
					t.Fatalf("pass %d: %v", pass, err)
				}
				if !bytes.Equal(want, buf.Bytes()) {
					t.Fatalf("pass %d drifted from fresh\nwant:\n%s\ngot:\n%s", pass, want, buf.String())
				}
			}
		})
	}
}

func TestContract_PoolIsolation(t *testing.T) {
	cases := contractCases()
	tableWant := make([][]byte, len(cases))
	streamWant := make([][]byte, len(cases))
	for i, c := range cases {
		tableWant[i] = renderContractTable(c)
		streamWant[i] = renderContractStream(c)
	}
	for round := range 3 {
		for i, c := range slices.Backward(cases) {
			if got := renderContractStream(c); !bytes.Equal(streamWant[i], got) {
				t.Fatalf("round %d %q stream leaked\nwant:\n%s\ngot:\n%s", round, c.name, streamWant[i], got)
			}
			if got := renderContractTable(c); !bytes.Equal(tableWant[i], got) {
				t.Fatalf("round %d %q table leaked\nwant:\n%s\ngot:\n%s", round, c.name, tableWant[i], got)
			}
		}
	}
}

func TestContract_TableRowWiderThanCount(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf)
	if err := tb.Render([][]table.Value{nil, {table.String("a")}, {table.String("b"), table.String("c"), table.String("d")}}); !errors.Is(err, table.ErrColumnCount) {
		t.Fatalf("expected table.ErrColumnCount, got %v", err)
	}
	if buf.Len() != 0 {
		t.Fatalf("expected no output, got:\n%s", buf.String())
	}
}

func TestContract_StreamRowWiderThanCount(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf)
	if err := s.Render(nil); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("a")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("b"), table.String("c"), table.String("d")}); !errors.Is(err, table.ErrColumnCount) {
		t.Fatalf("expected table.ErrColumnCount, got %v", err)
	}
}

func TestContract_TableHeaderRejectsOverflow(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf, WithHeader([]string{"A"}))
	err := tb.Render([][]table.Value{{table.String("a")}, {table.String("b"), table.String("c"), table.String("d")}})
	if _, ok := err.(*table.Error); !ok {
		t.Fatalf("expected outer *table.Error, got %T", err)
	}
	if !errors.Is(err, table.ErrColumnCount) {
		t.Fatalf("expected table.ErrColumnCount, got %v", err)
	}
	if buf.Len() != 0 {
		t.Fatalf("expected no output, got:\n%s", buf.String())
	}
}

func TestContract_StreamHeaderRejectsOverflow(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf, WithHeader([]string{"A"}))
	if err := s.Render([]table.Value{table.String("a"), table.String("b")}); !errors.Is(err, table.ErrColumnCount) {
		t.Fatalf("initial Render: expected table.ErrColumnCount, got %v", err)
	}
	if buf.Len() != 0 {
		t.Fatalf("initial Render wrote output:\n%s", buf.String())
	}
	if err := s.Render([]table.Value{table.String("a")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]table.Value{table.String("b"), table.String("c"), table.String("d")}); !errors.Is(err, table.ErrColumnCount) {
		t.Fatalf("expected table.ErrColumnCount, got %v", err)
	}
	if err := s.Render([]table.Value{table.String("b")}); err != nil {
		t.Fatalf("render after column count error: %v", err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestContract_FooterCannotExceedHeader(t *testing.T) {
	footer := func() [][]string {
		return [][]string{{"sum", "extra"}}
	}
	var tableOutput bytes.Buffer
	if err := NewTable(&tableOutput,
		WithHeader([]string{"A"}),
		WithFooter(footer),
	).Render([][]table.Value{{table.String("a")}}); !errors.Is(err, table.ErrColumnCount) {
		t.Fatalf("Table: expected table.ErrColumnCount, got %v", err)
	}
	if tableOutput.Len() != 0 {
		t.Fatalf("Table: expected no output, got:\n%s", tableOutput.String())
	}
	var streamOutput bytes.Buffer
	stream := NewStream(&streamOutput,
		WithHeader([]string{"A"}),
		WithFooter(footer),
	)
	if err := stream.Render([]table.Value{table.String("a")}); err != nil {
		t.Fatal(err)
	}
	first := stream.Close()
	if !errors.Is(first, table.ErrColumnCount) {
		t.Fatalf("Stream Close: expected table.ErrColumnCount, got %v", first)
	}
	if err := stream.Close(); err != first {
		t.Fatalf("second Close: expected %v, got %v", first, err)
	}
	var emptyOutput bytes.Buffer
	emptyStream := NewStream(&emptyOutput,
		WithHeader([]string{"A"}),
		WithFooter(footer),
	)
	first = emptyStream.Close()
	if !errors.Is(first, table.ErrColumnCount) {
		t.Fatalf("empty Stream Close: expected table.ErrColumnCount, got %v", first)
	}
	if err := emptyStream.Close(); err != first {
		t.Fatalf("empty Stream second Close: expected %v, got %v", first, err)
	}
	if emptyOutput.Len() != 0 {
		t.Fatalf("empty Stream wrote output:\n%s", emptyOutput.String())
	}
}

func TestContract_TableLineCapacity(t *testing.T) {
	rows := [][]table.Value{
		{table.String("group-1"), table.String("日本語の長い値です <escaped>"), table.Int(100), table.String("line1\nline2\nline3")},
		{table.String("group-1"), table.String("second value"), table.Int(99), table.String("x")},
		{table.String("group-2"), table.String("x"), table.Int(98), table.String("y")},
	}
	o := NewTable(io.Discard,
		WithHeader([]string{"Group", "Message", "Score", "Snippet"}),
		WithIndex(),
		WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
		WithColor(ScopeBody, Columns(1), ColorFgRed),
		WithDecoration(ScopeBody, Columns(1), DecorationBold),
		WithDecoration(ScopeBody, Columns(3), DecorationCode),
		WithAlign(ScopeHeader|ScopeBody|ScopeFooter, Columns(2), AlignRight),
		WithFooter(func() [][]string {
			return [][]string{{"", "", "297", ""}}
		}),
		WithCaption("caption", CaptionBottom),
		WithTableAttr(TableAttr{
			Table: Attr{
				Class: "tbl",
			},
			Caption: Attr{
				Class: "cap",
			},
			Header: SectionAttr{
				Section: Attr{
					Class: strings.Repeat("header-section-", 128),
				},
			},
			Body: SectionAttr{
				Section: Attr{
					Class: strings.Repeat("body-section-", 128),
				},
				Row: Attr{
					Class: "r",
				},
				Cell: Attr{
					Class: "c",
				},
			},
			Footer: SectionAttr{
				Section: Attr{
					Class: strings.Repeat("footer-section-", 128),
				},
			},
		}),
		WithCellAttr(ScopeHeader|ScopeBody|ScopeFooter, Columns(1), Attr{
			Style: "white-space:nowrap",
		}),
	)
	footer := o.option.footer()
	a := acquireArena()
	config := a.newConfig(&o.option, footer, len(rows), len(rows[0]))
	config.prepare()
	compiler := a.newCompiler(config.output)
	compiler.prepare()
	compiler.compileHeader()
	compiler.compileBody(rows)
	compiler.compileFooter()
	solver := a.newSolver(compiler.output)
	solver.solve()
	painter := a.newPainter(solver.output, o.w)
	painter.prepare()
	lineCap := cap(a.painter.line)
	painter.paintHeader()
	painter.paintBody()
	painter.paintFooter()
	got := cap(a.painter.line)
	a.release()
	if got > lineCap {
		t.Errorf("line capacity grew during rendering: got %d, prepared %d", got, lineCap)
	}
}

func TestContract_ArenaReleasedAfterClose(t *testing.T) {
	var buf bytes.Buffer
	stream := NewStream(&buf, WithHeader([]string{"A", "B"}))
	if err := stream.Render([]table.Value{table.String("a"), table.String("b")}); err != nil {
		t.Fatal(err)
	}
	a := stream.arena
	if err := stream.Close(); err != nil {
		t.Fatal(err)
	}
	if a.painter.line != nil {
		t.Fatal("line outlived the arena it was carved from")
	}
	if stream.arena != nil {
		t.Fatal("arena outlived the render")
	}
}

func TestContract_OptionIsShareable(t *testing.T) {
	opts := []Option{
		WithColor(ScopeBody, AllColumns(), ColorFgRed),
		WithTableAttr(TableAttr{Table: Attr{Class: `a&"b`}}),
		WithCellAttr(ScopeBody, Columns(0), Attr{Class: `a&"b`}),
	}
	var wg sync.WaitGroup
	for range 8 {
		wg.Go(func() {
			for range 4 {
				var buf bytes.Buffer
				if err := NewTable(&buf, opts...).Render([][]table.Value{{table.String("a"), table.Int(1)}}); err != nil {
					t.Error(err)
					return
				}
			}
		})
	}
	wg.Wait()
}

func TestContract_ColumnSelectors(t *testing.T) {
	indexes := []int{1}
	selector := Columns(indexes...)
	indexes[0] = 0
	configured := option{}
	configured.apply(
		WithIndex(),
		WithAlign(ScopeBody, AllColumns(), AlignRight),
		WithAlign(ScopeBody, selector, AlignCenter),
		WithRowspan(ScopeHeader, AllColumns()),
		WithRowspan(ScopeFooter, Columns(2)),
		WithAlign(ScopeBody, Columns(math.MaxInt/2), AlignCenter),
		WithAlign(ScopeBody, Columns(math.MaxInt), AlignLeft),
	)
	a := arena{}
	config := a.newConfig(&configured, nil, 0, 5)
	config.prepare()
	columns := config.output.columns
	cases := []struct {
		name    string
		index   int
		align   AlignSide
		rowspan Scope
	}{
		{name: "index", index: 0, align: AlignDefault},
		{name: "first input column", index: 1, align: AlignRight, rowspan: ScopeHeader},
		{name: "explicit alignment", index: 2, align: AlignCenter, rowspan: ScopeHeader},
		{name: "explicit rowspan", index: 3, align: AlignRight, rowspan: ScopeHeader | ScopeFooter},
		{name: "future input column", index: 5, align: AlignRight, rowspan: ScopeHeader},
	}
	for _, test := range cases {
		column := columns[test.index]
		align := column.aligns.Resolve(ScopeBody)
		if align != test.align || column.rowspan != test.rowspan {
			t.Errorf("%s: got align=%v rowspan=%v, want align=%v rowspan=%v",
				test.name, align, column.rowspan, test.align, test.rowspan)
		}
	}
}

func TestContract_ConcurrentInstances(t *testing.T) {
	cases := contractCases()
	wantTable := make([][]byte, len(cases))
	wantStream := make([][]byte, len(cases))
	for i, c := range cases {
		wantTable[i] = renderContractTable(c)
		wantStream[i] = renderContractStream(c)
	}
	var wg sync.WaitGroup
	for g := range 8 {
		wg.Go(func() {
			for pass := range 4 {
				for i, c := range cases {
					if got := renderContractTable(c); !bytes.Equal(wantTable[i], got) {
						t.Errorf("goroutine %d pass %d %q: table leaked\nwant:\n%s\ngot:\n%s",
							g, pass, c.name, wantTable[i], got)
						return
					}
					if got := renderContractStream(c); !bytes.Equal(wantStream[i], got) {
						t.Errorf("goroutine %d pass %d %q: stream leaked\nwant:\n%s\ngot:\n%s",
							g, pass, c.name, wantStream[i], got)
						return
					}
				}
			}
		})
	}
	wg.Wait()
}

func contractCases() []contractCase {
	return []contractCase{
		{
			name:   "empty header",
			header: []string{},
		},
		{
			name:       "empty footer",
			omitHeader: true,
			opts: []Option{
				WithFooter(func() [][]string {
					return [][]string{{}}
				}),
			},
		},
		{
			name:       "leading zero-column rows",
			omitHeader: true,
			rows:       [][]table.Value{nil, {}, {table.String("a"), table.String("b")}, {}, {table.String("c")}},
		},
		{
			name: "placeholder",
			opts: []Option{
				WithPlaceholder("-"),
			},
			header: []string{"A", "B", "C"},
			rows:   [][]table.Value{{table.String("x"), table.String(""), table.String("z")}, {table.String(""), table.String("q"), table.String("")}},
		},
		{
			name: "footer",
			opts: []Option{
				WithFooter(func() [][]string {
					return [][]string{{"sum", "", "9"}}
				}),
			},
			header: []string{"A", "B", "C"},
			rows:   [][]table.Value{{table.String("x"), table.String("y"), table.Int(1)}, {table.String("p"), table.String("q"), table.Int(8)}},
		},
		{
			name:   "ragged rows",
			opts:   []Option{},
			header: []string{"A", "B", "C"},
			rows:   [][]table.Value{{table.String("x")}, {table.String("p"), table.String("q"), table.String("r")}, {}},
		},
		{
			name:   "control chars",
			opts:   []Option{},
			header: []string{"A", "B", "C"},
			rows:   [][]table.Value{{table.String("a\tb"), table.String("c\vd"), table.String("e\x00f")}},
		},
		{
			name:   "invalid utf8",
			opts:   []Option{},
			header: []string{"A", "B", "C"},
			rows:   [][]table.Value{{table.String("a\xffb"), table.String("\xfe"), table.String("ok")}},
		},
		{
			name:   "emoji",
			opts:   []Option{},
			header: []string{"A", "B", "C"},
			rows:   [][]table.Value{{table.String("\U0001F600"), table.String("\U0001F469\u200D\U0001F4BB"), table.String("e\u0301")}},
		},
		{
			name:   "plain",
			header: []string{"A", "B", "C"},
			rows:   [][]table.Value{{table.String("xxx"), table.String("yyy"), table.String("zzz")}, {table.String("aaa"), table.String("bbb"), table.String("ccc")}},
		},
		{
			name:   "numeric",
			header: []string{"Int", "Float"},
			rows:   [][]table.Value{{table.Int(100), table.Float64(1.25)}, {table.Int(200), table.Float64(2.50)}, {table.Int(300), table.Float64(3.75)}},
		},
		{
			name:          "rowspan string",
			streamDiffers: true,
			opts:          []Option{WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0))},
			header:        []string{"Group", "Item"},
			rows:          [][]table.Value{{table.String("aaa"), table.String("x")}, {table.String("aaa"), table.String("y")}, {table.String("bbb"), table.String("z")}},
		},
		{
			name:          "rowspan multi",
			streamDiffers: true,
			opts:          []Option{WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1))},
			header:        []string{"Reg", "Zone", "Host"},
			rows: [][]table.Value{
				{table.String("jp"), table.String("1a"), table.String("h1")},
				{table.String("jp"), table.String("1a"), table.String("h2")},
				{table.String("jp"), table.String("1c"), table.String("h3")},
				{table.String("us"), table.String("1a"), table.String("h4")},
			},
		},
		{
			name: "colspan",
			opts: []Option{
				WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1, 2)),
			},
			header: []string{"A", "B", "C"},
			rows:   [][]table.Value{{table.String("x"), table.String("x"), table.String("y")}, {table.String("p"), table.String("q"), table.String("q")}},
		},
		{
			name:   "index",
			opts:   []Option{WithIndex()},
			header: []string{"Name", "Score"},
			rows:   [][]table.Value{{table.String("alice"), table.Int(100)}, {table.String("bob"), table.Int(200)}},
		},
		{
			name:   "color",
			opts:   []Option{WithColor(ScopeBody, Columns(0), ColorFgRed)},
			header: []string{"A", "B"},
			rows:   [][]table.Value{{table.String("xxx"), table.String("yyy")}, {table.String("aaa"), table.String("bbb")}},
		},
		{
			name:   "decoration code",
			opts:   []Option{WithDecoration(ScopeBody, Columns(1), DecorationCode)},
			header: []string{"Key", "Value"},
			rows:   [][]table.Value{{table.String("k1"), table.String("v1")}, {table.String("k2"), table.String("v2")}},
		},
		{
			name:          "color and decoration with rowspan",
			streamDiffers: true,
			opts: []Option{
				WithColor(ScopeBody, Columns(1), ColorFgBlue),
				WithDecoration(ScopeBody, Columns(1), DecorationBold),
				WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
			},
			header: []string{"Group", "Value"},
			rows:   [][]table.Value{{table.String("aaa"), table.String("x")}, {table.String("aaa"), table.String("y")}, {table.String("bbb"), table.String("z")}},
		},
		{
			name:          "rowspan and colspan",
			streamDiffers: true,
			opts: []Option{
				WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1)),
				WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1, 2)),
			},
			header: []string{"A", "B", "C"},
			rows:   [][]table.Value{{table.String("x"), table.String("x"), table.String("x")}, {table.String("x"), table.String("x"), table.String("y")}, {table.String("z"), table.String("z"), table.String("y")}},
		},
	}
}

func renderContractTable(c contractCase) []byte {
	var buf bytes.Buffer
	opts := c.opts
	if !c.omitHeader {
		opts = append([]Option{WithHeader(c.header)}, opts...)
	}
	tb := NewTable(&buf, opts...)
	if err := tb.Render(c.rows); err != nil {
		panic(err)
	}
	return append([]byte(nil), buf.Bytes()...)
}

func renderContractStream(c contractCase) []byte {
	var buf bytes.Buffer
	opts := c.opts
	if !c.omitHeader {
		opts = append([]Option{WithHeader(c.header)}, opts...)
	}
	s := NewStream(&buf, opts...)
	for _, row := range c.rows {
		if err := s.Render(row); err != nil {
			panic(err)
		}
	}
	if err := s.Close(); err != nil {
		panic(err)
	}
	return append([]byte(nil), buf.Bytes()...)
}
