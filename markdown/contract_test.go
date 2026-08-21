package markdown

import (
	"bytes"
	"errors"
	"io"
	"slices"
	"sync"
	"testing"

	"github.com/nekrassov01/table"
	"github.com/nekrassov01/table/internal/testutil"
)

const separatorLine = 1

type contractCase struct {
	name   string
	opts   []Option
	header []string
	rows   [][]any
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

func TestContract_TableHeaderRequired(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf)
	if err := tb.Render([][]any{{"x"}}); !errors.Is(err, table.ErrHeaderRequired) {
		t.Fatalf("expected table.ErrHeaderRequired, got %v", err)
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
		t.Fatal(err)
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
		t.Fatalf("expected nil, got %v", err)
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

func TestContract_StreamHeaderRequired(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf)
	if err := s.Render([]any{"x"}); !errors.Is(err, table.ErrHeaderRequired) {
		t.Fatalf("expected table.ErrHeaderRequired, got %v", err)
	}
}

func TestContract_DecorationWithoutPrefix(t *testing.T) {
	if deco := NewDecoration("", "suffix"); deco != nil {
		t.Fatalf("expected nil decoration, got %v", deco)
	}
	rows := [][]any{{"x", "y"}}
	var dressed, plain bytes.Buffer
	if err := NewTable(&dressed,
		WithHeader([]string{"A", "B"}),
		WithDecoration(ScopeHeader|ScopeBody, Columns(0), NewDecoration("", "suffix")),
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
		t.Run(c.name, func(t *testing.T) {
			table := normalizeContract(renderContractTable(c))
			stream := normalizeContract(renderContractStream(c))
			testutil.AssertBytes(t, stream, table, "stream against table")
		})
	}
}

func TestContract_ReuseMatchesFresh(t *testing.T) {
	for _, c := range contractCases() {
		t.Run(c.name, func(t *testing.T) {
			want := renderContractTable(c)
			var buf bytes.Buffer
			tb := NewTable(&buf, append([]Option{WithHeader(c.header)}, c.opts...)...)
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

func TestContract_StreamCloseWithoutHeader(t *testing.T) {
	var buf bytes.Buffer
	s := NewStream(&buf)
	first := s.Close()
	if first == nil {
		t.Fatal("expected header error, got nil")
	}
	if err := s.Close(); err != first {
		t.Fatalf("second Close: want the latched error, got %v", err)
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

func TestContract_TableLineCapacity(t *testing.T) {
	rows := [][]any{
		{"group-1", "日本語の長い値です | escaped", 100, "line1\nline2\nline3"},
		{"group-1", "second\nline", 99, "x"},
		{"group-2", "x", 98, "y"},
	}
	o := NewTable(io.Discard,
		WithHeader([]string{"Group", "Message", "Score", "Snippet"}),
		WithIndex(),
		WithRowspan(Columns(0)),
		WithColor(ScopeBody, Columns(1), ColorFgRed),
		WithDecoration(ScopeBody, Columns(1), DecorationBold),
		WithDecoration(ScopeBody, Columns(3), DecorationCode),
		WithAlign(Columns(2), AlignRight),
	)
	a := acquireArena()
	config := a.newConfig(&o.option, len(rows))
	config.prepare()
	compiler := a.newCompiler(config.output)
	compiler.prepare()
	compiler.compileHeader()
	compiler.compileBody(rows)
	solver := a.newSolver(compiler.output)
	solver.prepare()
	solver.solve()
	painter := a.newPainter(solver.output, o.w)
	painter.prepare()
	lineCap := cap(a.painter.line)
	painter.paintHeader()
	painter.paintBody()
	got := cap(a.painter.line)
	a.release()
	if got > lineCap {
		t.Errorf("line capacity grew during rendering: got %d, prepared %d", got, lineCap)
	}
}

func TestContract_StreamLineCapacity(t *testing.T) {
	s := NewStream(io.Discard,
		WithHeader([]string{"Message"}),
		WithColor(ScopeBody, Columns(0), ColorFgRed),
		WithDecoration(ScopeBody, Columns(0), DecorationBold),
	)
	if err := s.Render([]any{"short"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Render([]any{"日本語の長い値です | escaped\nnext"}); err != nil {
		t.Fatal(err)
	}
	lineCap := cap(s.arena.painter.line)
	backingCap := cap(s.arena.painter.lineBacking)
	if lineCap > backingCap {
		t.Errorf("line capacity grew during rendering: got %d, prepared %d", lineCap, backingCap)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestContract_ArenaReleasedAfterClose(t *testing.T) {
	var buf bytes.Buffer
	stream := NewStream(&buf, WithHeader([]string{"A", "B"}))
	if err := stream.Render([]any{"a", "b"}); err != nil {
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
		WithHeader([]string{"A", "B"}),
		WithColor(ScopeBody, AllColumns(), ColorFgRed),
		WithDecoration(ScopeBody, Columns(0), DecorationBold),
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
	configured.apply(
		WithHeader([]string{"A", "B", "C", "D", "E"}),
		WithIndex(),
		WithAlign(AllColumns(), AlignRight),
		WithAlign(selector, AlignCenter),
		WithRowspan(AllColumns()),
		WithColspan(Columns(2)),
	)
	a := arena{}
	config := a.newConfig(&configured, 0)
	config.prepare()
	columns := config.output.columns
	cases := []struct {
		name    string
		index   int
		align   AlignSide
		rowspan bool
		colspan bool
	}{
		{name: "index", index: 0},
		{name: "first input column", index: 1, align: AlignRight, rowspan: true},
		{name: "explicit alignment", index: 2, align: AlignCenter, rowspan: true},
		{name: "explicit colspan", index: 3, align: AlignRight, rowspan: true, colspan: true},
		{name: "future input column", index: 5, align: AlignRight, rowspan: true},
	}
	for _, test := range cases {
		column := columns[test.index]
		if column.align != test.align || column.rowspan != test.rowspan || column.colspan != test.colspan {
			t.Errorf("%s: got align=%v rowspan=%v colspan=%v, want align=%v rowspan=%v colspan=%v",
				test.name, column.align, column.rowspan, column.colspan,
				test.align, test.rowspan, test.colspan)
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
			name: "index",
			opts: []Option{
				WithIndex(),
			},
			header: []string{"A", "B", "C"},
			rows:   [][]any{{"x", "y", "z"}, {"p", "q", "r"}},
		},

		{
			name:   "ragged rows",
			opts:   []Option{},
			header: []string{"A", "B", "C"},
			rows:   [][]any{{"x"}, {"p", "q", "r"}, {}},
		},

		{
			name:   "control chars",
			opts:   []Option{},
			header: []string{"A", "B", "C"},
			rows:   [][]any{{"a\tb", "c\vd", "e\x00f"}},
		},

		{
			name:   "invalid utf8",
			opts:   []Option{},
			header: []string{"A", "B", "C"},
			rows:   [][]any{{"a\xffb", "\xfe", "ok"}},
		},

		{
			name:   "emoji",
			opts:   []Option{},
			header: []string{"A", "B", "C"},
			rows:   [][]any{{"\U0001F600", "\U0001F469\u200D\U0001F4BB", "e\u0301"}},
		},

		{
			name:   "plain",
			header: []string{"A", "B", "C"},
			rows:   [][]any{{"xxx", "yyy", "zzz"}, {"aaa", "bbb", "ccc"}},
		},

		{
			name:   "numeric",
			header: []string{"Int", "Float"},
			rows:   [][]any{{100, 1.25}, {200, 2.50}, {300, 3.75}},
		},

		{
			name:   "rowspan string",
			opts:   []Option{WithRowspan(Columns(0))},
			header: []string{"Group", "Item"},
			rows:   [][]any{{"aaa", "x"}, {"aaa", "y"}, {"bbb", "z"}},
		},

		{
			name:   "rowspan numeric",
			opts:   []Option{WithRowspan(Columns(0))},
			header: []string{"ID", "Name"},
			rows:   [][]any{{100, "a"}, {100, "b"}, {200, "c"}, {300, "d"}},
		},

		{
			name:   "rowspan multi",
			opts:   []Option{WithRowspan(Columns(0, 1))},
			header: []string{"Reg", "Zone", "Host"},
			rows: [][]any{
				{"jp", "1a", "h1"},
				{"jp", "1a", "h2"},
				{"jp", "1c", "h3"},
				{"us", "1a", "h4"},
			},
		},

		{
			name: "colspan",
			opts: []Option{
				WithColspan(Columns(0, 1, 2)),
			},
			header: []string{"A", "B", "C"},
			rows:   [][]any{{"x", "x", "y"}, {"p", "q", "q"}},
		},

		{
			name:   "color",
			opts:   []Option{WithColor(ScopeBody, Columns(0), ColorFgRed)},
			header: []string{"A", "B"},
			rows:   [][]any{{"xxx", "yyy"}, {"aaa", "bbb"}},
		},

		{
			name:   "decoration code",
			opts:   []Option{WithDecoration(ScopeBody, Columns(1), DecorationCode)},
			header: []string{"Key", "Value"},
			rows:   [][]any{{"k1", "v1"}, {"k2", "v2"}},
		},

		{
			name:   "decoration bold",
			opts:   []Option{WithDecoration(ScopeBody, Columns(0), DecorationBold)},
			header: []string{"Name", "Score"},
			rows:   [][]any{{"alice", 100}, {"bob", 200}},
		},

		{
			name: "color and decoration",
			opts: []Option{
				WithColor(ScopeBody, Columns(1), ColorFgBlue),
				WithDecoration(ScopeBody, Columns(1), DecorationBold),
				WithRowspan(Columns(0)),
			},
			header: []string{"Group", "Value"},
			rows:   [][]any{{"A", "x"}, {"A", "y"}, {"B", "z"}},
		},

		{
			name:   "placeholder",
			opts:   []Option{WithPlaceholder("N/A")},
			header: []string{"A", "B"},
			rows:   [][]any{{"x", nil}, {nil, "y"}},
		},
	}
}

func renderContractTable(c contractCase) []byte {
	var buf bytes.Buffer
	tb := NewTable(&buf, append([]Option{WithHeader(c.header)}, c.opts...)...)
	if err := tb.Render(c.rows); err != nil {
		panic(err)
	}
	return append([]byte(nil), buf.Bytes()...)
}

func renderContractStream(c contractCase) []byte {
	var buf bytes.Buffer
	s := NewStream(&buf, append([]Option{WithHeader(c.header)}, c.opts...)...)
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

func normalizeContract(b []byte) []byte {
	lines := bytes.Split(b, []byte("\n"))
	for i, line := range lines {
		fields := bytes.Split(line, []byte("|"))
		for j := range fields {
			fields[j] = bytes.TrimSpace(fields[j])
			if i == separatorLine {
				fields[j] = normalizeSeparator(fields[j])
			}
		}
		lines[i] = bytes.Join(fields, []byte("|"))
	}
	return bytes.Join(lines, []byte("\n"))
}

func normalizeSeparator(f []byte) []byte {
	body := bytes.Trim(f, ":")
	if len(body) < 3 || len(bytes.Trim(body, "-")) != 0 {
		return f
	}
	out := make([]byte, 0, len(f))
	if bytes.HasPrefix(f, []byte(":")) {
		out = append(out, ':')
	}
	out = append(out, "---"...)
	if bytes.HasSuffix(f, []byte(":")) {
		out = append(out, ':')
	}
	return out
}
