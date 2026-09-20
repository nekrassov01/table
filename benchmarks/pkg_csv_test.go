package benchmarks

import (
	"bytes"
	"io"
	"strconv"
	"testing"

	"github.com/nekrassov01/table"
	"github.com/nekrassov01/table/csv"
	"github.com/nekrassov01/table/examples"
)

func BenchmarkCSVTableSimpleFresh(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		t := csv.NewTable(w, examples.CSVOptionSimple...)
		if err := t.Render(examples.SimpleData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCSVTableSimpleReuse(b *testing.B) {
	w := &bytes.Buffer{}
	t := csv.NewTable(w, examples.CSVOptionSimple...)
	for b.Loop() {
		w.Reset()
		if err := t.Render(examples.SimpleData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCSVStreamSimple(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		s := csv.NewStream(w, examples.CSVOptionSimple...)
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

func BenchmarkCSVTableFooterFresh(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		t := csv.NewTable(w, examples.CSVOptionFooter...)
		if err := t.Render(examples.FooterData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCSVTableFooterReuse(b *testing.B) {
	w := &bytes.Buffer{}
	t := csv.NewTable(w, examples.CSVOptionFooter...)
	for b.Loop() {
		w.Reset()
		if err := t.Render(examples.FooterData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCSVStreamFooter(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		s := csv.NewStream(w, examples.CSVOptionFooter...)
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

func BenchmarkCSVTableTransformerFresh(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		t := csv.NewTable(w, examples.CSVOptionTransformer...)
		if err := t.Render(examples.FooterData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCSVTableTransformerReuse(b *testing.B) {
	w := &bytes.Buffer{}
	t := csv.NewTable(w, examples.CSVOptionTransformer...)
	for b.Loop() {
		w.Reset()
		if err := t.Render(examples.FooterData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCSVStreamTransformer(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		s := csv.NewStream(w, examples.CSVOptionTransformer...)
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

func BenchmarkCSVTableComplexFresh(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		t := csv.NewTable(w, examples.CSVOptionComplex...)
		if err := t.Render(examples.ComplexData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCSVTableComplexReuse(b *testing.B) {
	w := &bytes.Buffer{}
	t := csv.NewTable(w, examples.CSVOptionComplex...)
	for b.Loop() {
		w.Reset()
		if err := t.Render(examples.ComplexData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCSVStreamComplex(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		s := csv.NewStream(w, examples.CSVOptionComplex...)
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

func BenchmarkCSVTableCommaIncludedFresh(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		t := csv.NewTable(w, examples.CSVOptionCommaIncluded...)
		if err := t.Render(examples.CommaIncludedData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCSVTableCommaIncludedReuse(b *testing.B) {
	w := &bytes.Buffer{}
	t := csv.NewTable(w, examples.CSVOptionCommaIncluded...)
	for b.Loop() {
		w.Reset()
		if err := t.Render(examples.CommaIncludedData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCSVStreamCommaIncluded(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		s := csv.NewStream(w, examples.CSVOptionCommaIncluded...)
		for _, row := range examples.CommaIncludedData.Body {
			if err := s.Render(row); err != nil {
				b.Fatal(err)
			}
		}
		if err := s.Close(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCSVTableValueInputMixed(b *testing.B) {
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
	t := csv.NewTable(io.Discard, csv.WithHeader(header))
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

func BenchmarkCSVStreamValueInputMixed(b *testing.B) {
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
		s := csv.NewStream(io.Discard, csv.WithHeader(header))
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

func BenchmarkCSVTableValueInputSmallInts(b *testing.B) {
	records := make([]int, 1000)
	for i := range records {
		records[i] = 1000 + i
	}
	header := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	t := csv.NewTable(io.Discard, csv.WithHeader(header))
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

func BenchmarkCSVStreamValueInputSmallInts(b *testing.B) {
	records := make([]int, 1000)
	for i := range records {
		records[i] = 1000 + i
	}
	header := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	row := make([]table.Value, 0, 8)
	b.ReportAllocs()
	for b.Loop() {
		s := csv.NewStream(io.Discard, csv.WithHeader(header))
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

// BenchmarkCSVTableTransformerValue isolates callback input costs.
// Callbacks return existing strings to exclude formatting allocations.
func BenchmarkCSVTableTransformerValue(b *testing.B) {
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
	t := csv.NewTable(io.Discard,
		csv.WithHeader([]string{"MESSAGE", "LEVEL"}),
		csv.WithTransformer(csv.Columns(0), func(v table.Value) string {
			return messageText(v)
		}),
		csv.WithTransformer(csv.Columns(1), func(v table.Value) string {
			return levelText(v)
		}),
	)
	b.ReportAllocs()
	for b.Loop() {
		if err := t.Render(rows); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkCSVStreamTransformerValue isolates callback input costs.
// Callbacks return existing strings to exclude formatting allocations.
func BenchmarkCSVStreamTransformerValue(b *testing.B) {
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
		s := csv.NewStream(io.Discard,
			csv.WithHeader([]string{"MESSAGE", "LEVEL"}),
			csv.WithTransformer(csv.Columns(0), func(v table.Value) string {
				return messageText(v)
			}),
			csv.WithTransformer(csv.Columns(1), func(v table.Value) string {
				return levelText(v)
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
