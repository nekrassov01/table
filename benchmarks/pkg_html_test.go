package benchmarks

import (
	"bytes"
	"io"
	"strconv"
	"testing"

	"github.com/nekrassov01/table"
	"github.com/nekrassov01/table/examples"
	"github.com/nekrassov01/table/html"
)

func BenchmarkHTMLTableSimpleFresh(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		t := html.NewTable(w, examples.HTMLOptionSimple...)
		if err := t.Render(examples.SimpleData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkHTMLTableSimpleReuse(b *testing.B) {
	w := &bytes.Buffer{}
	t := html.NewTable(w, examples.HTMLOptionSimple...)
	for b.Loop() {
		w.Reset()
		if err := t.Render(examples.SimpleData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkHTMLStreamSimple(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		s := html.NewStream(w, examples.HTMLOptionSimple...)
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

func BenchmarkHTMLTableRowspanFresh(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		t := html.NewTable(w, examples.HTMLOptionRowspan...)
		if err := t.Render(examples.RowspanData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkHTMLTableRowspanReuse(b *testing.B) {
	w := &bytes.Buffer{}
	t := html.NewTable(w, examples.HTMLOptionRowspan...)
	for b.Loop() {
		w.Reset()
		if err := t.Render(examples.RowspanData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkHTMLStreamRowspan(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		s := html.NewStream(w, examples.HTMLOptionRowspan...)
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

func BenchmarkHTMLTableColspanFresh(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		t := html.NewTable(w, examples.HTMLOptionColspan...)
		if err := t.Render(examples.ColspanData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkHTMLTableColspanReuse(b *testing.B) {
	w := &bytes.Buffer{}
	t := html.NewTable(w, examples.HTMLOptionColspan...)
	for b.Loop() {
		w.Reset()
		if err := t.Render(examples.ColspanData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkHTMLStreamColspan(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		s := html.NewStream(w, examples.HTMLOptionColspan...)
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

func BenchmarkHTMLTableFooterFresh(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		t := html.NewTable(w, examples.HTMLOptionFooter...)
		if err := t.Render(examples.FooterData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkHTMLTableFooterReuse(b *testing.B) {
	w := &bytes.Buffer{}
	t := html.NewTable(w, examples.HTMLOptionFooter...)
	for b.Loop() {
		w.Reset()
		if err := t.Render(examples.FooterData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkHTMLStreamFooter(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		s := html.NewStream(w, examples.HTMLOptionFooter...)
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

func BenchmarkHTMLTableTransformerFresh(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		t := html.NewTable(w, examples.HTMLOptionTransformer...)
		if err := t.Render(examples.FooterData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkHTMLTableTransformerReuse(b *testing.B) {
	w := &bytes.Buffer{}
	t := html.NewTable(w, examples.HTMLOptionTransformer...)
	for b.Loop() {
		w.Reset()
		if err := t.Render(examples.FooterData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkHTMLStreamTransformer(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		s := html.NewStream(w, examples.HTMLOptionTransformer...)
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

func BenchmarkHTMLTableComplexFresh(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		t := html.NewTable(w, examples.HTMLOptionComplex...)
		if err := t.Render(examples.ComplexData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkHTMLTableComplexReuse(b *testing.B) {
	w := &bytes.Buffer{}
	t := html.NewTable(w, examples.HTMLOptionComplex...)
	for b.Loop() {
		w.Reset()
		if err := t.Render(examples.ComplexData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkHTMLStreamComplex(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		s := html.NewStream(w, examples.HTMLOptionComplex...)
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

func BenchmarkHTMLTableStackedHeaderFresh(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		t := html.NewTable(w, examples.HTMLOptionStackedHeader...)
		if err := t.Render(examples.StackedHeaderData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkHTMLTableStackedHeaderReuse(b *testing.B) {
	w := &bytes.Buffer{}
	t := html.NewTable(w, examples.HTMLOptionStackedHeader...)
	for b.Loop() {
		w.Reset()
		if err := t.Render(examples.StackedHeaderData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkHTMLStreamStackedHeader(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		s := html.NewStream(w, examples.HTMLOptionStackedHeader...)
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

func BenchmarkHTMLTableValueInputMixed(b *testing.B) {
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
	t := html.NewTable(io.Discard, html.WithHeader(header))
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

func BenchmarkHTMLStreamValueInputMixed(b *testing.B) {
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
		s := html.NewStream(io.Discard, html.WithHeader(header))
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

func BenchmarkHTMLTableValueInputSmallInts(b *testing.B) {
	records := make([]int, 1000)
	for i := range records {
		records[i] = 1000 + i
	}
	header := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	t := html.NewTable(io.Discard, html.WithHeader(header))
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

func BenchmarkHTMLStreamValueInputSmallInts(b *testing.B) {
	records := make([]int, 1000)
	for i := range records {
		records[i] = 1000 + i
	}
	header := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	row := make([]table.Value, 0, 8)
	b.ReportAllocs()
	for b.Loop() {
		s := html.NewStream(io.Discard, html.WithHeader(header))
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
