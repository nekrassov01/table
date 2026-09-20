package benchmarks

import (
	"bytes"
	"io"
	"strconv"
	"testing"

	"github.com/nekrassov01/table"
	"github.com/nekrassov01/table/examples"
	"github.com/nekrassov01/table/text"
)

func BenchmarkTextTableASCIIFresh(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		t := text.NewTable(w, examples.TextOptionASCII...)
		if err := t.Render(examples.SimpleData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTextTableASCIIReuse(b *testing.B) {
	w := &bytes.Buffer{}
	t := text.NewTable(w, examples.TextOptionASCII...)
	for b.Loop() {
		w.Reset()
		if err := t.Render(examples.SimpleData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTextStreamASCII(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		s := text.NewStream(w, examples.TextOptionASCII...)
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

func BenchmarkTextTableSimpleFresh(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		t := text.NewTable(w, examples.TextOptionSimple...)
		if err := t.Render(examples.SimpleData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTextTableSimpleReuse(b *testing.B) {
	w := &bytes.Buffer{}
	t := text.NewTable(w, examples.TextOptionSimple...)
	for b.Loop() {
		w.Reset()
		if err := t.Render(examples.SimpleData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTextStreamSimple(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		s := text.NewStream(w, examples.TextOptionSimple...)
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

func BenchmarkTextTableCompactFresh(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		t := text.NewTable(w, examples.TextOptionCompact...)
		if err := t.Render(examples.CompactData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTextTableCompactReuse(b *testing.B) {
	w := &bytes.Buffer{}
	t := text.NewTable(w, examples.TextOptionCompact...)
	for b.Loop() {
		w.Reset()
		if err := t.Render(examples.CompactData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTextStreamCompact(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		s := text.NewStream(w, examples.TextOptionCompact...)
		for _, row := range examples.CompactData.Body {
			if err := s.Render(row); err != nil {
				b.Fatal(err)
			}
		}
		if err := s.Close(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTextTableRowspanFresh(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		t := text.NewTable(w, examples.TextOptionRowspan...)
		if err := t.Render(examples.RowspanData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTextTableRowspanReuse(b *testing.B) {
	w := &bytes.Buffer{}
	t := text.NewTable(w, examples.TextOptionRowspan...)
	for b.Loop() {
		w.Reset()
		if err := t.Render(examples.RowspanData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTextStreamRowspan(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		s := text.NewStream(w, examples.TextOptionRowspan...)
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

func BenchmarkTextTableColspanFresh(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		t := text.NewTable(w, examples.TextOptionColspan...)
		if err := t.Render(examples.ColspanData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTextTableColspanReuse(b *testing.B) {
	w := &bytes.Buffer{}
	t := text.NewTable(w, examples.TextOptionColspan...)
	for b.Loop() {
		w.Reset()
		if err := t.Render(examples.ColspanData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTextStreamColspan(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		s := text.NewStream(w, examples.TextOptionColspan...)
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

func BenchmarkTextTableFooterFresh(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		t := text.NewTable(w, examples.TextOptionFooter...)
		if err := t.Render(examples.FooterData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTextTableFooterReuse(b *testing.B) {
	w := &bytes.Buffer{}
	t := text.NewTable(w, examples.TextOptionFooter...)
	for b.Loop() {
		w.Reset()
		if err := t.Render(examples.FooterData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTextStreamFooter(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		s := text.NewStream(w, examples.TextOptionFooter...)
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

func BenchmarkTextTableTransformerFresh(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		t := text.NewTable(w, examples.TextOptionTransformer...)
		if err := t.Render(examples.FooterData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTextTableTransformerReuse(b *testing.B) {
	w := &bytes.Buffer{}
	t := text.NewTable(w, examples.TextOptionTransformer...)
	for b.Loop() {
		w.Reset()
		if err := t.Render(examples.FooterData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTextStreamTransformer(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		s := text.NewStream(w, examples.TextOptionTransformer...)
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

func BenchmarkTextTableComplexFresh(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		t := text.NewTable(w, examples.TextOptionComplex...)
		if err := t.Render(examples.ComplexData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTextTableComplexReuse(b *testing.B) {
	w := &bytes.Buffer{}
	t := text.NewTable(w, examples.TextOptionComplex...)
	for b.Loop() {
		w.Reset()
		if err := t.Render(examples.ComplexData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTextStreamComplex(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		s := text.NewStream(w, examples.TextOptionComplex...)
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

func BenchmarkTextTableStackedHeaderFresh(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		t := text.NewTable(w, examples.TextOptionStackedHeader...)
		if err := t.Render(examples.StackedHeaderData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTextTableStackedHeaderReuse(b *testing.B) {
	w := &bytes.Buffer{}
	t := text.NewTable(w, examples.TextOptionStackedHeader...)
	for b.Loop() {
		w.Reset()
		if err := t.Render(examples.StackedHeaderData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTextStreamStackedHeader(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		s := text.NewStream(w, examples.TextOptionStackedHeader...)
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

func BenchmarkTextTableValueInputMixed(b *testing.B) {
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
	t := text.NewTable(io.Discard, text.WithHeader(header))
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

func BenchmarkTextStreamValueInputMixed(b *testing.B) {
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
		s := text.NewStream(io.Discard, text.WithHeader(header))
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

func BenchmarkTextTableValueInputSmallInts(b *testing.B) {
	records := make([]int, 1000)
	for i := range records {
		records[i] = 1000 + i
	}
	header := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	t := text.NewTable(io.Discard, text.WithHeader(header))
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

func BenchmarkTextStreamValueInputSmallInts(b *testing.B) {
	records := make([]int, 1000)
	for i := range records {
		records[i] = 1000 + i
	}
	header := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	row := make([]table.Value, 0, 8)
	b.ReportAllocs()
	for b.Loop() {
		s := text.NewStream(io.Discard, text.WithHeader(header))
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

// BenchmarkTextTableTransformerValue isolates callback input costs.
// Callbacks return existing strings to exclude formatting allocations.
func BenchmarkTextTableTransformerValue(b *testing.B) {
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
	t := text.NewTable(io.Discard,
		text.WithHeader([]string{"MESSAGE", "LEVEL"}),
		text.WithTransformer(text.Columns(0), func(v table.Value) (string, *text.Attr) {
			return messageText(v), nil
		}),
		text.WithTransformer(text.Columns(1), func(v table.Value) (string, *text.Attr) {
			return levelText(v), nil
		}),
	)
	b.ReportAllocs()
	for b.Loop() {
		if err := t.Render(rows); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkTextStreamTransformerValue isolates callback input costs.
// Callbacks return existing strings to exclude formatting allocations.
func BenchmarkTextStreamTransformerValue(b *testing.B) {
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
		s := text.NewStream(io.Discard,
			text.WithHeader([]string{"MESSAGE", "LEVEL"}),
			text.WithTransformer(text.Columns(0), func(v table.Value) (string, *text.Attr) {
				return messageText(v), nil
			}),
			text.WithTransformer(text.Columns(1), func(v table.Value) (string, *text.Attr) {
				return levelText(v), nil
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
