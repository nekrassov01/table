package text

import (
	"bytes"
	"errors"
	"io"
	"math"
	"os"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/nekrassov01/table"
	"github.com/nekrassov01/table/internal/display"
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
		t.Fatalf("Close: %v", err)
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
	rows := [][]table.Value{{table.String("x"), table.String("y")}}
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
	rows := [][]table.Value{{table.String("xxxxx"), table.Any(strings.Repeat("y", 40))}}
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

func TestContract_WindowsTerminalColor(t *testing.T) {
	file, _ := testutil.NewFile(t)
	var output bytes.Buffer
	originalIsTerminal := isTerminal
	originalTerminalWriter := terminalWriter
	isTerminal = func(w io.Writer) bool {
		return w == file
	}
	terminalWriter = func(*os.File) io.Writer {
		return &output
	}
	t.Cleanup(func() {
		isTerminal = originalIsTerminal
		terminalWriter = originalTerminalWriter
	})
	opts := []Option{
		WithStyle(StyleColoredLight),
	}
	if err := NewTable(file, opts...).Render([][]table.Value{{table.String("x")}}); err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(output.Bytes(), []byte("\x1b[")) {
		t.Fatalf("table output does not contain ANSI attributes:\n%s", output.Bytes())
	}
	output.Reset()
	stream := NewStream(file, opts...)
	if err := stream.Render([]table.Value{table.String("x")}); err != nil {
		t.Fatal(err)
	}
	if err := stream.Close(); err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(output.Bytes(), []byte("\x1b[")) {
		t.Fatalf("stream output does not contain ANSI attributes:\n%s", output.Bytes())
	}
}

func TestContract_MultiColumnFillKeepsFrameWidth(t *testing.T) {
	style := StyleASCII
	top := *style.Border.Top
	top.Fill = "界"
	style.Border.Top = &top
	var buf bytes.Buffer
	if err := NewTable(&buf, WithStyle(style)).Render([][]table.Value{{table.String("a"), table.String("bb")}}); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	want := display.StringWidth(lines[1])
	for i, line := range lines {
		if got := display.StringWidth(line); got != want {
			t.Fatalf("line %d is %d wide, the frame is %d: %s", i+1, got, want, line)
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
	if err := s.Render([]table.Value{table.String("b")}); err != nil {
		t.Fatalf("render after column count error: %v", err)
	}
	if err := s.Close(); err != nil {
		t.Fatalf("close after column count error: %v", err)
	}
}

func TestContract_EmptyHeaderDoesNotPinCount(t *testing.T) {
	rows := [][]table.Value{{table.String("a")}, {table.String("b"), table.String("c"), table.String("d")}}
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
		if err := s.Render([]table.Value{table.Any(i % 10)}); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	want := display.StringWidth(lines[0])
	for i, l := range lines {
		if w := display.StringWidth(l); w != want {
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
			rows:       [][]table.Value{nil, {}, {table.String("a"), table.String("b")}, {}, {table.String("c")}},
		},
		{
			name:          "index",
			streamDiffers: true,
			opts: []Option{
				WithStyle(StyleASCII),
				WithIndex(),
			},
			header: []string{"A", "B", "C"},
			rows:   [][]table.Value{{table.String("x"), table.String("y"), table.String("z")}, {table.String("p"), table.String("q"), table.String("r")}},
		},
		{
			name: "placeholder",
			opts: []Option{
				WithStyle(StyleASCII),
				WithPlaceholder("-"),
			},
			header: []string{"A", "B", "C"},
			rows:   [][]table.Value{{table.String("x"), table.String(""), table.String("z")}, {table.String(""), table.String("q"), table.String("")}},
		},
		{
			name: "ragged rows",
			opts: []Option{
				WithStyle(StyleASCII),
			},
			header: []string{"A", "B", "C"},
			rows:   [][]table.Value{{table.String("x")}, {table.String("p"), table.String("q"), table.String("r")}, {}},
		},
		{
			name: "control chars",
			opts: []Option{
				WithStyle(StyleASCII),
			},
			header: []string{"A", "B", "C"},
			rows:   [][]table.Value{{table.String("a\tb"), table.String("c\vd"), table.String("e\x00f")}},
		},
		{
			name: "invalid utf8",
			opts: []Option{
				WithStyle(StyleASCII),
			},
			header: []string{"A", "B", "C"},
			rows:   [][]table.Value{{table.String("a\xffb"), table.String("\xfe"), table.String("ok")}},
		},
		{
			name: "emoji",
			opts: []Option{
				WithStyle(StyleASCII),
			},
			header: []string{"A", "B", "C"},
			rows:   [][]table.Value{{table.String("\U0001F600"), table.String("\U0001F469\u200D\U0001F4BB"), table.String("e\u0301")}},
		},
		{
			name:   "plain",
			opts:   []Option{WithStyle(StyleASCII)},
			header: []string{"A", "B", "C"},
			rows:   [][]table.Value{{table.String("xxx"), table.String("yyy"), table.String("zzz")}, {table.String("aaa"), table.String("bbb"), table.String("ccc")}},
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
			rows:   [][]table.Value{{table.String("x"), table.String("yy")}},
		},
		{
			name:   "numeric columns",
			opts:   []Option{WithStyle(StyleASCII)},
			header: []string{"Int", "Float"},
			rows:   [][]table.Value{{table.Int(100), table.Float64(1.25)}, {table.Int(200), table.Float64(2.50)}, {table.Int(300), table.Float64(3.75)}},
		},
		{
			name:   "rowspan string",
			opts:   []Option{WithStyle(StyleASCII), WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0))},
			header: []string{"Group", "Item"},
			rows:   [][]table.Value{{table.String("aaa"), table.String("x")}, {table.String("aaa"), table.String("y")}, {table.String("bbb"), table.String("z")}},
		},
		{
			name:   "rowspan numeric",
			opts:   []Option{WithStyle(StyleASCII), WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0))},
			header: []string{"ID", "Name"},
			rows:   [][]table.Value{{table.Int(100), table.String("a")}, {table.Int(100), table.String("b")}, {table.Int(200), table.String("c")}, {table.Int(300), table.String("d")}},
		},
		{
			name:   "rowspan multi",
			opts:   []Option{WithStyle(StyleASCII), WithRowspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1))},
			header: []string{"Reg", "Zone", "Host"},
			rows: [][]table.Value{
				{table.String("jp"), table.String("1a"), table.String("h1")},
				{table.String("jp"), table.String("1a"), table.String("h2")},
				{table.String("jp"), table.String("1c"), table.String("h3")},
				{table.String("us"), table.String("1a"), table.String("h4")},
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
			rows:   [][]table.Value{{table.String("aaa"), table.Int(1)}, {table.String("aaa"), table.Int(2)}, {table.String("bbb"), table.Int(3)}},
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
			rows:   [][]table.Value{{table.String("a"), table.Int(111)}, {table.String("b"), table.Int(222)}, {table.String("c"), table.Int(333)}},
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
			rows:   [][]table.Value{{table.String("a"), table.String("b"), table.Int(1)}, {table.String("a"), table.String("b"), table.Int(2)}},
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
			rows:   [][]table.Value{{table.String("long body merge"), table.String("long body merge"), table.Int(3), table.Int(4), table.String("x")}},
		},
		{
			name: "clamped width and padding",
			opts: []Option{
				WithStyle(StyleASCII),
				WithWidth(Columns(1), -1),
				WithPadding(Columns(0), -1, -2),
			},
			header: []string{"Key", "Value"},
			rows:   [][]table.Value{{table.String("a"), table.Int(111)}, {table.String("b"), table.Int(222)}},
		},
		{
			name: "colspan",
			opts: []Option{
				WithStyle(StyleASCII),
				WithColspan(ScopeHeader|ScopeBody|ScopeFooter, Columns(0, 1, 2)),
			},
			header: []string{"A", "B", "C"},
			rows:   [][]table.Value{{table.String("x"), table.String("x"), table.String("y")}, {table.String("p"), table.String("q"), table.String("q")}},
		},
		{
			name: "rowspan beside colspan",
			opts: []Option{
				WithStyle(StyleASCII),
				WithRowspan(ScopeBody, Columns(1, 2)),
				WithColspan(ScopeBody, Columns(0, 1, 2)),
			},
			header: []string{"c0", "c1", "c2"},
			rows:   [][]table.Value{{table.String("A"), table.String("A"), table.String("A")}, {table.String("X"), table.String("A"), table.String("A")}},
		},
		{
			name:       "headerless initial colspan",
			omitHeader: true,
			opts: []Option{
				WithStyle(StyleASCII),
				WithColspan(ScopeBody, Columns(0, 1)),
			},
			rows: [][]table.Value{{table.String("A"), table.String("A")}},
		},
		{
			name:   "colors stripped off terminal",
			opts:   []Option{WithStyle(StyleColoredLight), WithAttr(ScopeBody, Columns(0), ColorFgRed)},
			header: []string{"A", "B"},
			rows:   [][]table.Value{{table.String("xxx"), table.String("yyy")}, {table.String("aaa"), table.String("bbb")}},
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
