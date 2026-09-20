package csv

import (
	"bytes"
	"errors"
	"io"
	"math"
	"slices"
	"strings"
	"sync"
	"testing"
	"unicode/utf8"

	"github.com/nekrassov01/table"
	"github.com/nekrassov01/table/internal/testutil"
)

type contractCase struct {
	name   string
	opts   []Option
	header []string
	rows   [][]table.Value
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
	empty := NewStream(io.Discard)
	if err := empty.Close(); err != nil {
		t.Fatalf("empty stream: expected nil, got %v", err)
	}
	zeroColumns := NewStream(io.Discard, WithFooter(func() [][]string {
		return [][]string{{}}
	}))
	if err := zeroColumns.Close(); err != nil {
		t.Fatalf("zero-column footer: expected nil, got %v", err)
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

func TestContract_TableStreamAgree(t *testing.T) {
	for _, c := range contractCases() {
		t.Run(c.name, func(t *testing.T) {
			table := renderContractTable(c)
			stream := renderContractStream(c)
			testutil.AssertBytes(t, stream, table, "stream against table")
		})
	}
}

func TestContract_CRLF(t *testing.T) {
	footer := func() [][]string {
		return [][]string{{"total", " done\nnow"}}
	}
	rows := [][]table.Value{
		{table.String("crlf"), table.String("a\r\nb")},
		{table.String("cr"), table.String("a\rb")},
		{table.String("lf"), table.String("a\nb")},
		{table.String("leading"), table.String(" value")},
		{table.String("postgres"), table.String(`\.`)},
	}
	want := []byte(
		"Key,Value\r\n" +
			"crlf,\"a\r\nb\"\r\n" +
			"cr,\"ab\"\r\n" +
			"lf,\"a\r\nb\"\r\n" +
			"leading,\" value\"\r\n" +
			"postgres,\"\\.\"\r\n" +
			"total,\" done\r\nnow\"\r\n",
	)
	var tableOutput bytes.Buffer
	if err := NewTable(&tableOutput,
		WithDelimiter(','),
		WithCRLF(),
		WithHeader([]string{"Key", "Value"}),
		WithFooter(footer),
	).Render(rows); err != nil {
		t.Fatal(err)
	}
	testutil.AssertBytes(t, tableOutput.Bytes(), want, "Table CRLF")

	var streamOutput bytes.Buffer
	stream := NewStream(&streamOutput,
		WithDelimiter(','),
		WithCRLF(),
		WithHeader([]string{"Key", "Value"}),
		WithFooter(footer),
	)
	for _, row := range rows {
		if err := stream.Render(row); err != nil {
			t.Fatal(err)
		}
	}
	if err := stream.Close(); err != nil {
		t.Fatal(err)
	}
	testutil.AssertBytes(t, streamOutput.Bytes(), want, "Stream CRLF")
}

func TestContract_TableInvalidDelimiter(t *testing.T) {
	for _, d := range []rune{'"', '\n', '\r', 0, utf8.RuneError, -1, 0xd800, utf8.MaxRune + 1} {
		var buf bytes.Buffer
		tb := NewTable(&buf, WithDelimiter(d), WithHeader([]string{"A", "B"}))
		err := tb.Render([][]table.Value{{table.String("x"), table.String("y")}})
		if _, ok := err.(*table.Error); !ok {
			t.Fatalf("delimiter %q: expected outer *table.Error, got %T", d, err)
		}
		if !errors.Is(err, table.ErrDelimiter) {
			t.Fatalf("delimiter %q: expected table.ErrDelimiter, got %v", d, err)
		}
		if buf.Len() != 0 {
			t.Fatalf("delimiter %q: expected no output, got %q", d, buf.String())
		}
	}
}

func TestContract_StreamInvalidDelimiter(t *testing.T) {
	for _, d := range []rune{'"', '\n', '\r', 0, utf8.RuneError, -1, 0xd800, utf8.MaxRune + 1} {
		var buf bytes.Buffer
		s := NewStream(&buf, WithDelimiter(d), WithHeader([]string{"A", "B"}))
		if err := s.Render([]table.Value{table.String("x"), table.String("y")}); !errors.Is(err, table.ErrDelimiter) {
			t.Fatalf("delimiter %q: expected table.ErrDelimiter, got %v", d, err)
		}
		if buf.Len() != 0 {
			t.Fatalf("delimiter %q: expected no output, got %q", d, buf.String())
		}
	}
}

func TestContract_StreamCloseInvalidDelimiter(t *testing.T) {
	tests := []struct {
		name string
		opts []Option
	}{
		{
			name: "no sections",
		},
		{
			name: "header",
			opts: []Option{WithHeader([]string{"A", "B"})},
		},
		{
			name: "footer",
			opts: []Option{WithFooter(func() [][]string {
				return [][]string{{"A", "B"}}
			})},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, d := range []rune{'"', '\n', '\r', 0, utf8.RuneError, -1, 0xd800, utf8.MaxRune + 1} {
				var buf bytes.Buffer
				opts := append([]Option{WithDelimiter(d)}, test.opts...)
				s := NewStream(&buf, opts...)
				first := s.Close()
				if !errors.Is(first, table.ErrDelimiter) {
					t.Fatalf("delimiter %q: expected table.ErrDelimiter, got %v", d, first)
				}
				if err := s.Close(); err != first {
					t.Fatalf("delimiter %q: second Close: expected %v, got %v", d, first, err)
				}
				if buf.Len() != 0 {
					t.Fatalf("delimiter %q: expected no output, got %q", d, buf.String())
				}
			}
		})
	}
}

func TestContract_ColumnOptionReplaces(t *testing.T) {
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithHeader([]string{"A"}),
		WithTransformer(Columns(0), func(table.Value) string {
			return "first"
		}),
		WithTransformer(Columns(0), func(table.Value) string {
			return "second"
		}),
	)
	if err := tb.Render([][]table.Value{{table.String("x")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertBytes(t, buf.Bytes(), []byte("A\nsecond\n"), "the later transformer")
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
	footer := []string{"sum", "extra"}
	var tableOutput bytes.Buffer
	if err := NewTable(&tableOutput,
		WithHeader([]string{"A"}),
		WithFooter(func() [][]string {
			return [][]string{footer}
		}),
	).Render([][]table.Value{{table.String("a")}}); !errors.Is(err, table.ErrColumnCount) {
		t.Fatalf("Table: expected table.ErrColumnCount, got %v", err)
	}
	if tableOutput.Len() != 0 {
		t.Fatalf("Table: expected no output, got:\n%s", tableOutput.String())
	}
	var streamOutput bytes.Buffer
	stream := NewStream(&streamOutput,
		WithHeader([]string{"A"}),
		WithFooter(func() [][]string {
			return [][]string{footer}
		}),
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
		WithFooter(func() [][]string {
			return [][]string{footer}
		}),
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
		{table.String("group-1"), table.String("value, quoted"), table.Int(100)},
		{table.String("group-1"), table.String("plain"), table.Int(99)},
		{table.String("group-2"), table.String("x"), table.Int(98)},
	}
	o := NewTable(io.Discard,
		WithHeader([]string{"Group", "Message", "Score"}),
		WithDelimiter(','),
		WithCRLF(),
		WithIndex(),
		WithFooter(func() [][]string {
			return [][]string{{"", "", "297"}}
		}),
	)
	a := acquireArena()
	footer := o.option.footer()
	config := a.newConfig(&o.option, footer, len(rows), len(rows[0]))
	config.prepare()
	compiler := a.newCompiler(config.output)
	compiler.prepare()
	compiler.compileHeader()
	compiler.compileBody(rows)
	compiler.compileFooter()
	painter := a.newPainter(compiler.output, o.w)
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

func TestContract_OptionIsShareable(t *testing.T) {
	opts := []Option{
		WithHeader([]string{"A", "B"}),
		WithTransformer(AllColumns(), func(table.Value) string {
			return "value"
		}),
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
	var buf bytes.Buffer
	tb := NewTable(&buf,
		WithIndex(),
		WithTransformer(AllColumns(), func(table.Value) string {
			return "all"
		}),
		WithTransformer(selector, func(table.Value) string {
			return "explicit"
		}),
		WithTransformer(Columns(math.MaxInt/2), func(table.Value) string {
			return "large"
		}),
		WithTransformer(Columns(math.MaxInt), func(table.Value) string {
			return "missing"
		}),
	)
	if err := tb.Render([][]table.Value{{table.String("a"), table.String("b"), table.String("c")}}); err != nil {
		t.Fatal(err)
	}
	testutil.AssertBytes(t, buf.Bytes(), []byte("1\tall\texplicit\tall\n"), "selectors")
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

func TestContract_TransformerValue(t *testing.T) {
	type args struct {
		input table.Value
		read  func(table.Value) any
	}
	type want struct {
		value any
		calls int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "string",
			args: args{
				input: table.String("text"),
				read: func(v table.Value) any {
					return v.AsString()
				},
			},
			want: want{
				value: "text",
				calls: 1,
			},
		},
		{
			name: "int",
			args: args{
				input: table.Int(1000),
				read: func(v table.Value) any {
					return v.AsInt()
				},
			},
			want: want{
				value: 1000,
				calls: 1,
			},
		},
		{
			name: "bytes capacity",
			args: args{
				input: table.Bytes(make([]byte, 2, 8)),
				read: func(v table.Value) any {
					return cap(v.AsBytes())
				},
			},
			want: want{
				value: 2,
				calls: 1,
			},
		},
		{
			name: "Any capacity",
			args: args{
				input: table.Any(make([]byte, 2, 8)),
				read: func(v table.Value) any {
					return cap(v.AsBytes())
				},
			},
			want: want{
				value: 8,
				calls: 1,
			},
		},
		{
			name: "typed nil",
			args: args{
				input: table.Any((*int)(nil)),
				read: func(v table.Value) any {
					return v.AsAny()
				},
			},
			want: want{
				value: (*int)(nil),
				calls: 1,
			},
		},
		{
			name: "zero",
			args: args{
				input: table.Value{},
				read: func(v table.Value) any {
					return v.AsAny()
				},
			},
			want: want{
				value: nil,
				calls: 1,
			},
		},
	}
	for _, test := range tests {
		t.Run("Table/"+test.name, func(t *testing.T) {
			got := want{}
			output := NewTable(io.Discard,
				WithHeader([]string{"VALUE"}),
				WithTransformer(Columns(0), func(v table.Value) string {
					got.value = test.args.read(v)
					got.calls++
					return "transformed"
				}),
			)
			err := output.Render([][]table.Value{{test.args.input}})
			testutil.AssertValue(t, err, nil, "Render")
			testutil.AssertValue(t, got, test.want, "transformer input")
		})
		t.Run("Stream/"+test.name, func(t *testing.T) {
			got := want{}
			output := NewStream(io.Discard,
				WithHeader([]string{"VALUE"}),
				WithTransformer(Columns(0), func(v table.Value) string {
					got.value = test.args.read(v)
					got.calls++
					return "transformed"
				}),
			)
			err := output.Render([]table.Value{test.args.input})
			closeErr := output.Close()
			testutil.AssertValue(t, err, nil, "Render")
			testutil.AssertValue(t, closeErr, nil, "Close")
			testutil.AssertValue(t, got, test.want, "transformer input")
		})
	}
}

func contractCases() []contractCase {
	return []contractCase{
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
			name:   "index",
			opts:   []Option{WithIndex()},
			header: []string{"Name", "Score"},
			rows:   [][]table.Value{{table.String("alice"), table.Int(100)}, {table.String("bob"), table.Int(200)}},
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
