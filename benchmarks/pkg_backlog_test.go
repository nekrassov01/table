package benchmarks

import (
	"bytes"
	"io"
	"strconv"
	"testing"

	"github.com/nekrassov01/table"
	"github.com/nekrassov01/table/backlog"
	"github.com/nekrassov01/table/examples"
)

func BenchmarkBacklogTableSimpleFresh(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		t := backlog.NewTable(w, examples.BacklogOptionSimple...)
		if err := t.Render(examples.SimpleData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBacklogTableSimpleReuse(b *testing.B) {
	w := &bytes.Buffer{}
	t := backlog.NewTable(w, examples.BacklogOptionSimple...)
	for b.Loop() {
		w.Reset()
		if err := t.Render(examples.SimpleData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBacklogStreamSimple(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		s := backlog.NewStream(w, examples.BacklogOptionSimple...)
		for _, row := range examples.SimpleData.Body {
			if err := s.Render(row); err != nil {
				b.Fatal(err)
			}
		}
		if err := s.Close(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBacklogTableRowspanFresh(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		t := backlog.NewTable(w, examples.BacklogOptionRowspan...)
		if err := t.Render(examples.RowspanData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBacklogTableRowspanReuse(b *testing.B) {
	w := &bytes.Buffer{}
	t := backlog.NewTable(w, examples.BacklogOptionRowspan...)
	for b.Loop() {
		w.Reset()
		if err := t.Render(examples.RowspanData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBacklogStreamRowspan(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		s := backlog.NewStream(w, examples.BacklogOptionRowspan...)
		for _, row := range examples.RowspanData.Body {
			if err := s.Render(row); err != nil {
				b.Fatal(err)
			}
		}
		if err := s.Close(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBacklogTableColspanFresh(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		t := backlog.NewTable(w, examples.BacklogOptionColspan...)
		if err := t.Render(examples.ColspanData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBacklogTableColspanReuse(b *testing.B) {
	w := &bytes.Buffer{}
	t := backlog.NewTable(w, examples.BacklogOptionColspan...)
	for b.Loop() {
		w.Reset()
		if err := t.Render(examples.ColspanData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBacklogStreamColspan(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		s := backlog.NewStream(w, examples.BacklogOptionColspan...)
		for _, row := range examples.ColspanData.Body {
			if err := s.Render(row); err != nil {
				b.Fatal(err)
			}
		}
		if err := s.Close(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBacklogTableFooterFresh(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		t := backlog.NewTable(w, examples.BacklogOptionFooter...)
		if err := t.Render(examples.FooterData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBacklogTableFooterReuse(b *testing.B) {
	w := &bytes.Buffer{}
	t := backlog.NewTable(w, examples.BacklogOptionFooter...)
	for b.Loop() {
		w.Reset()
		if err := t.Render(examples.FooterData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBacklogStreamFooter(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		s := backlog.NewStream(w, examples.BacklogOptionFooter...)
		for _, row := range examples.FooterData.Body {
			if err := s.Render(row); err != nil {
				b.Fatal(err)
			}
		}
		if err := s.Close(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBacklogTableTransformerFresh(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		t := backlog.NewTable(w, examples.BacklogOptionTransformer...)
		if err := t.Render(examples.FooterData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBacklogTableTransformerReuse(b *testing.B) {
	w := &bytes.Buffer{}
	t := backlog.NewTable(w, examples.BacklogOptionTransformer...)
	for b.Loop() {
		w.Reset()
		if err := t.Render(examples.FooterData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBacklogStreamTransformer(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		s := backlog.NewStream(w, examples.BacklogOptionTransformer...)
		for _, row := range examples.FooterData.Body {
			if err := s.Render(row); err != nil {
				b.Fatal(err)
			}
		}
		if err := s.Close(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBacklogTableComplexFresh(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		t := backlog.NewTable(w, examples.BacklogOptionComplex...)
		if err := t.Render(examples.ComplexData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBacklogTableComplexReuse(b *testing.B) {
	w := &bytes.Buffer{}
	t := backlog.NewTable(w, examples.BacklogOptionComplex...)
	for b.Loop() {
		w.Reset()
		if err := t.Render(examples.ComplexData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBacklogStreamComplex(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		s := backlog.NewStream(w, examples.BacklogOptionComplex...)
		for _, row := range examples.ComplexData.Body {
			if err := s.Render(row); err != nil {
				b.Fatal(err)
			}
		}
		if err := s.Close(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBacklogTableStackedHeaderFresh(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		t := backlog.NewTable(w, examples.BacklogOptionStackedHeader...)
		if err := t.Render(examples.StackedHeaderData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBacklogTableStackedHeaderReuse(b *testing.B) {
	w := &bytes.Buffer{}
	t := backlog.NewTable(w, examples.BacklogOptionStackedHeader...)
	for b.Loop() {
		w.Reset()
		if err := t.Render(examples.StackedHeaderData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBacklogStreamStackedHeader(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		s := backlog.NewStream(w, examples.BacklogOptionStackedHeader...)
		for _, row := range examples.StackedHeaderData.Body {
			if err := s.Render(row); err != nil {
				b.Fatal(err)
			}
		}
		if err := s.Close(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBacklogTableValueInputMixed(b *testing.B) {
	records := make([]struct {
		name   string
		number int
		ratio  float64
		data   []byte
	}, 1000)
	for i := range records {
		records[i].name = "resource-" + strconv.Itoa(i)
		records[i].number = 1000 + i
		records[i].ratio = float64(i) + 0.25
		records[i].data = []byte("payload")
	}
	header := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	t := backlog.NewTable(io.Discard, backlog.WithHeader(header))
	b.ReportAllocs()
	for b.Loop() {
		rows := make([][]table.Value, len(records))
		cells := make([]table.Value, len(records)*8)
		for i, v := range records {
			rows[i] = append(cells[i*8:i*8:i*8+8], table.String(v.name), table.String(v.name), table.String(v.name), table.Int(v.number), table.Int64(int64(v.number)), table.Float64(v.ratio), table.Bool(true), table.Bytes(v.data))
		}
		if err := t.Render(rows); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBacklogStreamValueInputMixed(b *testing.B) {
	records := make([]struct {
		name   string
		number int
		ratio  float64
		data   []byte
	}, 1000)
	for i := range records {
		records[i].name = "resource-" + strconv.Itoa(i)
		records[i].number = 1000 + i
		records[i].ratio = float64(i) + 0.25
		records[i].data = []byte("payload")
	}
	header := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	row := make([]table.Value, 0, 8)
	b.ReportAllocs()
	for b.Loop() {
		s := backlog.NewStream(io.Discard, backlog.WithHeader(header))
		for _, v := range records {
			row = append(row[:0], table.String(v.name), table.String(v.name), table.String(v.name), table.Int(v.number), table.Int64(int64(v.number)), table.Float64(v.ratio), table.Bool(true), table.Bytes(v.data))
			if err := s.Render(row); err != nil {
				b.Fatal(err)
			}
		}
		if err := s.Close(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBacklogTableValueInputSmallInts(b *testing.B) {
	records := make([]int, 1000)
	for i := range records {
		records[i] = 1000 + i
	}
	header := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	t := backlog.NewTable(io.Discard, backlog.WithHeader(header))
	b.ReportAllocs()
	for b.Loop() {
		rows := make([][]table.Value, len(records))
		cells := make([]table.Value, len(records)*8)
		for i, v := range records {
			n := v % 128
			rows[i] = append(cells[i*8:i*8:i*8+8], table.Int(n), table.Int(n), table.Int(n), table.Int(n), table.Int(n), table.Int(n), table.Int(n), table.Int(n))
		}
		if err := t.Render(rows); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBacklogStreamValueInputSmallInts(b *testing.B) {
	records := make([]int, 1000)
	for i := range records {
		records[i] = 1000 + i
	}
	header := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	row := make([]table.Value, 0, 8)
	b.ReportAllocs()
	for b.Loop() {
		s := backlog.NewStream(io.Discard, backlog.WithHeader(header))
		for _, v := range records {
			n := v % 128
			row = append(row[:0], table.Int(n), table.Int(n), table.Int(n), table.Int(n), table.Int(n), table.Int(n), table.Int(n), table.Int(n))
			if err := s.Render(row); err != nil {
				b.Fatal(err)
			}
		}
		if err := s.Close(); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkBacklogTableTransformerValue isolates callback input costs.
// Callbacks return existing strings to exclude formatting allocations.
func BenchmarkBacklogTableTransformerValue(b *testing.B) {
	message := "request completed"
	rows := make([][]table.Value, 1000)
	cells := make([]table.Value, len(rows)*2)
	for i := range rows {
		cells[i*2] = table.String(message)
		cells[i*2+1] = table.Int(1000 + i)
		rows[i] = cells[i*2 : i*2+2]
	}
	messageText := func(v table.Value) string {
		return v.AsString()
	}
	levelText := func(v table.Value) string {
		if v.AsInt() >= 1500 {
			return "warn"
		}
		return "info"
	}
	t := backlog.NewTable(io.Discard,
		backlog.WithHeader([]string{"MESSAGE", "LEVEL"}),
		backlog.WithTransformer(backlog.Columns(0), func(v table.Value) (string, *backlog.Color, *backlog.Decoration) {
			return messageText(v), nil, nil
		}),
		backlog.WithTransformer(backlog.Columns(1), func(v table.Value) (string, *backlog.Color, *backlog.Decoration) {
			return levelText(v), nil, nil
		}),
	)
	b.ReportAllocs()
	for b.Loop() {
		if err := t.Render(rows); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkBacklogStreamTransformerValue isolates callback input costs.
// Callbacks return existing strings to exclude formatting allocations.
func BenchmarkBacklogStreamTransformerValue(b *testing.B) {
	message := "request completed"
	const rowCount = 1000
	messageText := func(v table.Value) string {
		return v.AsString()
	}
	levelText := func(v table.Value) string {
		if v.AsInt() >= 1500 {
			return "warn"
		}
		return "info"
	}
	row := make([]table.Value, 2)
	b.ReportAllocs()
	for b.Loop() {
		s := backlog.NewStream(io.Discard,
			backlog.WithHeader([]string{"MESSAGE", "LEVEL"}),
			backlog.WithTransformer(backlog.Columns(0), func(v table.Value) (string, *backlog.Color, *backlog.Decoration) {
				return messageText(v), nil, nil
			}),
			backlog.WithTransformer(backlog.Columns(1), func(v table.Value) (string, *backlog.Color, *backlog.Decoration) {
				return levelText(v), nil, nil
			}),
		)
		for i := range rowCount {
			row[0] = table.String(message)
			row[1] = table.Int(1000 + i)
			if err := s.Render(row); err != nil {
				b.Fatal(err)
			}
		}
		if err := s.Close(); err != nil {
			b.Fatal(err)
		}
	}
}
