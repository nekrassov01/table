package text

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
	"github.com/nekrassov01/table/internal/width"
)

type contractCase struct {
	name          string
	opts          []Option
	header        []string
	rows          [][]any
	omitHeader    bool
	streamDiffers bool
}

func TestContract_TableDeterministic(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	header := []string{"A", "B", "C"}
	rows := [][]any{{"x", 1, true}, {"y", 2, false}}
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
	rows := [][]any{{"x", 1}, {"y", 2}}
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
	if err := tb.Render([][]any{{"a", 1}}); err != nil {
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
	err := tb.Render([][]any{{"x"}})
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
	for _, row := range [][]any{{"x", 1}, {"y", 2}} {
		if err := s.Render(row); err != nil {
			t.Fatalf("Render: %v", err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	if buf.Len() == 0 {
		t.Fatal("expected output, got empty")
	}
}

func TestContract_StreamRenderAfterClose(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf, WithHeader([]string{"A"}))
	if err := s.Render([]any{"x"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"y"}); !errors.Is(err, table.ErrClosed) {
		t.Fatalf("expected table.ErrClosed, got %v", err)
	}
}

func TestContract_StreamCloseAfterClose(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf, WithHeader([]string{"A"}))
	if err := s.Render([]any{"x"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatalf("expected nil on double close, got %v", err)
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
	if err := s.Render([]any{"a", 1}); err != nil {
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
	rows := [][]any{{1}}
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
	if strings.Contains(streamOutput.String(), "footer-only") {
		t.Fatalf("Stream partially wrote the wider footer:\n%s", streamOutput.String())
	}
}

func TestContract_StreamWriterError(t *testing.T) {
	cause := testutil.NewError()
	s := NewStream(&testutil.ErrorWriter{
		Err: cause,
	},
		WithHeader([]string{"A"}),
	)
	err := s.Render([]any{"x"})
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
			tb := NewTable(&buf, contractOptions(c)...)
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

func TestContract_StreamWriterErrorSticky(t *testing.T) {
	s := NewStream(&testutil.ErrorWriter{
		Err: testutil.NewError(),
	},
		WithHeader([]string{"A"}),
	)
	first := s.Render([]any{"x"})
	if first == nil {
		t.Fatal("expected error, got nil")
	}
	if err := s.Render([]any{"y"}); err == nil || err != first {
		t.Fatalf("second Render: want the latched error, got %v", err)
	}
	if err := s.Close(); err == nil || err != first {
		t.Fatalf("Close: want the latched error, got %v", err)
	}
}

func TestContract_TableZeroColumns(t *testing.T) {
	var buf bytes.Buffer
	if err := NewTable(&buf).Render([][]any{}); err != nil {
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

func TestContract_AttrWithoutCodes(t *testing.T) {
	if attr := NewAttr(); attr != nil {
		t.Fatalf("expected nil attribute, got %v", attr)
	}
	restore := isTerminal
	isTerminal = func(io.Writer) bool {
		return true
	}
	t.Cleanup(func() {
		isTerminal = restore
	})
	rows := [][]any{{"x", "y"}}
	var dressed, plain bytes.Buffer
	if err := NewTable(&dressed,
		WithHeader([]string{"A", "B"}),
		WithAttr(ScopeHeader|ScopeBody|ScopeFooter, Columns(0), NewAttr()),
	).Render(rows); err != nil {
		t.Fatal(err)
	}
	if err := NewTable(&plain, WithHeader([]string{"A", "B"})).Render(rows); err != nil {
		t.Fatal(err)
	}
	testutil.AssertBytes(t, dressed.Bytes(), plain.Bytes(), "table with an empty attribute")
}

func TestContract_FileIsNotATerminal(t *testing.T) {
	opts := []Option{
		WithStyle(StyleColoredLight),
		WithAutoFit(),
		WithHeader([]string{"A", "B"}),
	}
	rows := [][]any{{"xxxxx", strings.Repeat("y", 40)}}
	f, read := testutil.NewFile(t)
	if err := NewTable(f, opts...).Render(rows); err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := NewTable(&buf, opts...).Render(rows); err != nil {
		t.Fatal(err)
	}
	testutil.AssertBytes(t, read(), buf.Bytes(), "render into a file")
	if strings.Contains(buf.String(), "\x1b[") {
		t.Fatal("wrote ANSI escapes to output that is not a terminal")
	}
}

func TestContract_MultiColumnFillKeepsFrameWidth(t *testing.T) {
	style := StyleASCII
	top := *style.Border.Top
	top.Fill = "界"
	style.Border.Top = &top
	var buf bytes.Buffer
	if err := NewTable(&buf, WithStyle(style)).Render([][]any{{"a", "bb"}}); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	want := width.StringWidth(lines[1])
	for i, line := range lines {
		if got := width.StringWidth(line); got != want {
			t.Fatalf("line %d is %d wide, the frame is %d: %s", i+1, got, want, line)
		}
	}
}

func TestContract_TableRowWiderThanCount(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf)
	if err := tb.Render([][]any{nil, {"a"}, {"b", "c", "d"}}); !errors.Is(err, table.ErrColumnCount) {
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
	if err := s.Render([]any{"a"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"b", "c", "d"}); !errors.Is(err, table.ErrColumnCount) {
		t.Fatalf("expected table.ErrColumnCount, got %v", err)
	}
	if err := s.Render([]any{"b"}); err != nil {
		t.Fatalf("render after column count error: %v", err)
	}
	if err := s.Close(); err != nil {
		t.Fatalf("close after column count error: %v", err)
	}
}

func TestContract_EmptyHeaderDoesNotPinCount(t *testing.T) {
	rows := [][]any{{"a"}, {"b", "c", "d"}}
	var tableOutput bytes.Buffer
	if err := NewTable(&tableOutput, WithHeader([]string{})).Render(rows); !errors.Is(err, table.ErrColumnCount) {
		t.Fatalf("Table: expected table.ErrColumnCount, got %v", err)
	}
	if tableOutput.Len() != 0 {
		t.Fatalf("Table: expected no output, got:\n%s", tableOutput.String())
	}
	var streamOutput bytes.Buffer
	stream := NewStream(&streamOutput, WithHeader([]string{}))
	if err := stream.Render(rows[0]); err != nil {
		t.Fatal(err)
	}
	if err := stream.Render(rows[1]); !errors.Is(err, table.ErrColumnCount) {
		t.Fatalf("Stream: expected table.ErrColumnCount, got %v", err)
	}
	if err := stream.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestContract_TableHeaderRejectsOverflow(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf, WithHeader([]string{"A"}))
	err := tb.Render([][]any{{"a"}, {"b", "c", "d"}})
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
	if err := s.Render([]any{"a", "b"}); !errors.Is(err, table.ErrColumnCount) {
		t.Fatalf("initial Render: expected table.ErrColumnCount, got %v", err)
	}
	if buf.Len() != 0 {
		t.Fatalf("initial Render wrote output:\n%s", buf.String())
	}
	if err := s.Render([]any{"a"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"b", "c", "d"}); !errors.Is(err, table.ErrColumnCount) {
		t.Fatalf("expected table.ErrColumnCount, got %v", err)
	}
	if err := s.Render([]any{"b"}); err != nil {
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
	).Render([][]any{{"a"}}); !errors.Is(err, table.ErrColumnCount) {
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
	if err := stream.Render([]any{"a"}); err != nil {
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

func TestContract_ArenaReleaseDropsViews(t *testing.T) {
	attr := &Attr{
		Prefix: []byte("attr"),
	}
	a := &arena{
		compiler: compilerState{
			cells: []cell{
				{
					value: "value",
					attr:  attr,
				},
			},
			spanValues: []string{"value"},
			rows: []row{
				{
					cells: []cell{
						{
							value: "row",
						},
					},
				},
			},
		},
		painter: painterState{
			line:    make([]byte, 0, 8),
			horizon: []byte("horizon"),
			layouts: []layout{
				{
					value: "cell",
					attr:  attr,
				},
			},
			segments: []segment{
				{
					value: "segment",
				},
			},
		},
	}
	a.release()
	if a.painter.line != nil || a.painter.horizon != nil {
		t.Fatal("temporary view remained after arena release")
	}
	if a.compiler.cells[0].value != "" || a.compiler.cells[0].attr != nil ||
		a.compiler.spanValues[0] != "" || a.compiler.rows[0].cells != nil ||
		a.painter.layouts[0].value != "" || a.painter.segments[0].value != "" {
		t.Fatal("scratch element retained a reference after release")
	}
}

func TestContract_OptionIsShareable(t *testing.T) {
	opts := []Option{
		WithStyle(StyleASCII),
		WithWidth(Columns(1), -1),
		WithPadding(Columns(0), -1, -2),
	}
	var wg sync.WaitGroup
	for range 8 {
		wg.Go(func() {
			for range 4 {
				var buf bytes.Buffer
				if err := NewTable(&buf, opts...).Render([][]any{{"a", 1}}); err != nil {
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
	configured.apply(io.Discard, 0,
		WithIndex(),
		WithWidth(AllColumns(), 7),
		WithWidth(selector, 3),
		WithPadding(Columns(2), 5, 6),
		WithPadding(AllColumns(), 1, 2),
		WithPadding(Columns(2), 3, 4),
		WithWidth(Columns(math.MaxInt/2), 98),
		WithWidth(Columns(math.MaxInt), 99),
	)
	cases := []struct {
		name  string
		index int
		width int
		left  int
		right int
	}{
		{name: "index", index: 0, width: 0, left: 1, right: 1},
		{name: "first data column", index: 1, width: 7, left: 1, right: 2},
		{name: "explicit width", index: 2, width: 3, left: 1, right: 2},
		{name: "explicit padding", index: 3, width: 7, left: 3, right: 4},
		{name: "future data column", index: 5, width: 7, left: 1, right: 2},
	}
	a := arena{}
	config := a.newConfig(&configured, nil, 0, 5)
	config.prepare()
	columns := config.output.columns
	for _, tc := range cases {
		column := columns[tc.index]
		if column.limit != tc.width || column.lPad != tc.left || column.rPad != tc.right {
			t.Errorf("%s: got width=%d padding=%d/%d, want width=%d padding=%d/%d",
				tc.name, column.limit, column.lPad, column.rPad, tc.width, tc.left, tc.right)
		}
	}
	autoFit := option{}
	autoFit.apply(io.Discard, 0, WithAutoFit(), WithWidth(AllColumns(), 7))
	if autoFit.autoFit {
		t.Error("all-column width did not disable automatic fitting")
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

func TestContract_IndexWidthHoldsPastTheFloor(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf, WithStyle(StyleLight), WithHeader([]string{"A"}), WithIndexWidth(4))
	for i := range 1002 {
		if err := s.Render([]any{i % 10}); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	want := width.StringWidth(lines[0])
	for i, l := range lines {
		if w := width.StringWidth(l); w != want {
			t.Fatalf("line %d is %d wide, the frame is %d: %s", i+1, w, want, l)
		}
	}
	last := lines[len(lines)-2]
	if !strings.Contains(last, "1002") {
		t.Fatalf("the last row does not carry its number whole: %s", last)
	}
}

func contractCases() []contractCase {
	return []contractCase{
		{
			name:   "empty header",
			header: []string{},
		},
		{
			name:   "blank header cell",
			header: []string{""},
			opts:   []Option{WithPlaceholder("")},
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
			name:       "nil footer",
			omitHeader: true,
			opts:       []Option{WithFooter(nil)},
		},
		{
			name:       "leading zero-column rows",
			omitHeader: true,
			rows:       [][]any{nil, {}, {"a", "b"}, {}, {"c"}},
		},
		{
			name:          "index",
			streamDiffers: true,
			opts: []Option{
				WithStyle(StyleASCII),
				WithIndex(),
			},
			header: []string{"A", "B", "C"},
			rows:   [][]any{{"x", "y", "z"}, {"p", "q", "r"}},
		},
		{
			name: "placeholder",
			opts: []Option{
				WithStyle(StyleASCII),
				WithPlaceholder("-"),
			},
			header: []string{"A", "B", "C"},
			rows:   [][]any{{"x", "", "z"}, {"", "q", ""}},
		},
		{
			name: "ragged rows",
			opts: []Option{
				WithStyle(StyleASCII),
			},
			header: []string{"A", "B", "C"},
			rows:   [][]any{{"x"}, {"p", "q", "r"}, {}},
		},
		{
			name: "control chars",
			opts: []Option{
				WithStyle(StyleASCII),
			},
			header: []string{"A", "B", "C"},
			rows:   [][]any{{"a\tb", "c\vd", "e\x00f"}},
		},
		{
			name: "invalid utf8",
			opts: []Option{
				WithStyle(StyleASCII),
			},
			header: []string{"A", "B", "C"},
			rows:   [][]any{{"a\xffb", "\xfe", "ok"}},
		},
		{
			name: "emoji",
			opts: []Option{
				WithStyle(StyleASCII),
			},
			header: []string{"A", "B", "C"},
			rows:   [][]any{{"\U0001F600", "\U0001F469\u200D\U0001F4BB", "e\u0301"}},
		},
		{
			name:   "plain",
			opts:   []Option{WithStyle(StyleASCII)},
			header: []string{"A", "B", "C"},
			rows:   [][]any{{"xxx", "yyy", "zzz"}, {"aaa", "bbb", "ccc"}},
		},
		{
			name: "all input columns",
			opts: []Option{
				WithStyle(StyleASCII),
				WithIndexWidth(3),
				WithWidth(AllColumns(), 1),
				WithPadding(AllColumns(), 0, 0),
			},
			header: []string{"A", "B"},
			rows:   [][]any{{"x", "yy"}},
		},
		{
			name:   "numeric columns",
			opts:   []Option{WithStyle(StyleASCII)},
			header: []string{"Int", "Float"},
			rows:   [][]any{{100, 1.25}, {200, 2.50}, {300, 3.75}},
		},
		{
			name:   "rowspan string",
			opts:   []Option{WithStyle(StyleASCII), WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0))},
			header: []string{"Group", "Item"},
			rows:   [][]any{{"aaa", "x"}, {"aaa", "y"}, {"bbb", "z"}},
		},
		{
			name:   "rowspan numeric",
			opts:   []Option{WithStyle(StyleASCII), WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0))},
			header: []string{"ID", "Name"},
			rows:   [][]any{{100, "a"}, {100, "b"}, {200, "c"}, {300, "d"}},
		},
		{
			name:   "rowspan multi",
			opts:   []Option{WithStyle(StyleASCII), WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1))},
			header: []string{"Reg", "Zone", "Host"},
			rows: [][]any{
				{"jp", "1a", "h1"},
				{"jp", "1a", "h2"},
				{"jp", "1c", "h3"},
				{"us", "1a", "h4"},
			},
		},
		{
			name: "footer",
			opts: []Option{
				WithStyle(StyleASCII),
				WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0)),
				WithFooter(func() [][]string {
					return [][]string{{"tot", "9"}}
				}),
			},
			header: []string{"Grp", "Num"},
			rows:   [][]any{{"aaa", 1}, {"aaa", 2}, {"bbb", 3}},
		},
		{
			name: "footer colspan without body",
			opts: []Option{
				WithStyle(StyleASCII),
				WithFooter(func() [][]string {
					return [][]string{{"F", "F"}}
				}),
				WithColspan(ScopeFooter, Columns(0, 1)),
			},
			header: []string{"A", "B"},
		},
		{
			name:          "truncated footer label",
			streamDiffers: true,
			opts: []Option{
				WithStyle(StyleASCII),
				WithFooter(func() [][]string {
					return [][]string{{"total", "LongFooterLabel"}}
				}),
				WithWidth(Columns(1), 8),
				WithTruncate(Columns(1)),
			},
			header: []string{"Key", "Value"},
			rows:   [][]any{{"a", 111}, {"b", 222}, {"c", 333}},
		},
		{
			name: "spanned band with footer",
			opts: []Option{
				WithStyle(StyleASCII),
				WithFooter(func() [][]string {
					return [][]string{{"Sum", "Sum", "3"}, {"Sum", "Sum", "3"}}
				}),
				WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1, 2)),
				WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1, 2)),
			},
			header: []string{"Group", "Group", "Val"},
			rows:   [][]any{{"a", "b", 1}, {"a", "b", 2}},
		},
		{
			name: "widened merges in two parts",
			opts: []Option{
				WithStyle(StyleASCII),
				WithFooter(func() [][]string {
					return [][]string{{"long footer label", "long footer label", "long footer label", "long footer label", "F"}}
				}),
				WithColspan(ScopeBody, Columns(0, 1)),
				WithColspan(ScopeFooter, Columns(0, 1, 2, 3)),
			},
			header: []string{"a", "b", "c", "d", "ID"},
			rows:   [][]any{{"long body merge", "long body merge", 3, 4, "x"}},
		},
		{
			name: "clamped width and padding",
			opts: []Option{
				WithStyle(StyleASCII),
				WithWidth(Columns(1), -1),
				WithPadding(Columns(0), -1, -2),
			},
			header: []string{"Key", "Value"},
			rows:   [][]any{{"a", 111}, {"b", 222}},
		},
		{
			name: "colspan",
			opts: []Option{
				WithStyle(StyleASCII),
				WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1, 2)),
			},
			header: []string{"A", "B", "C"},
			rows:   [][]any{{"x", "x", "y"}, {"p", "q", "q"}},
		},
		{
			name: "rowspan beside colspan",
			opts: []Option{
				WithStyle(StyleASCII),
				WithRowspan(ScopeBody, Columns(1, 2)),
				WithColspan(ScopeBody, Columns(0, 1, 2)),
			},
			header: []string{"c0", "c1", "c2"},
			rows:   [][]any{{"A", "A", "A"}, {"X", "A", "A"}},
		},
		{
			name:       "headerless initial colspan",
			omitHeader: true,
			opts: []Option{
				WithStyle(StyleASCII),
				WithColspan(ScopeBody, Columns(0, 1)),
			},
			rows: [][]any{{"A", "A"}},
		},
		{
			name:   "colors stripped off terminal",
			opts:   []Option{WithStyle(StyleColoredLight), WithAttr(ScopeBody, Columns(0), ColorFgRed)},
			header: []string{"A", "B"},
			rows:   [][]any{{"xxx", "yyy"}, {"aaa", "bbb"}},
		},
	}
}

func contractOptions(c contractCase) []Option {
	opts := make([]Option, 0, len(c.opts)+1)
	if !c.omitHeader {
		opts = append(opts, WithHeader(c.header))
	}
	return append(opts, c.opts...)
}

func renderContractTable(c contractCase) []byte {
	var buf bytes.Buffer
	tb := NewTable(&buf, contractOptions(c)...)
	if err := tb.Render(c.rows); err != nil {
		panic(err)
	}
	return append([]byte(nil), buf.Bytes()...)
}

func renderContractStream(c contractCase) []byte {
	var buf bytes.Buffer
	s := NewStream(&buf, contractOptions(c)...)
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
