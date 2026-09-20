package benchmarks

import (
	"bytes"
	"io"
	"strconv"
	"testing"

	"github.com/nekrassov01/table"
	"github.com/nekrassov01/table/examples"
	"github.com/nekrassov01/table/markdown"
)

func BenchmarkMarkdownTableSimpleFresh(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		t := markdown.NewTable(w, examples.MarkdownOptionSimple...)
		if err := t.Render(examples.SimpleData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarkdownTableSimpleReuse(b *testing.B) {
	w := &bytes.Buffer{}
	t := markdown.NewTable(w, examples.MarkdownOptionSimple...)
	for b.Loop() {
		w.Reset()
		if err := t.Render(examples.SimpleData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarkdownStreamSimple(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		s := markdown.NewStream(w, examples.MarkdownOptionSimple...)
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

func BenchmarkMarkdownTableRowspanFresh(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		t := markdown.NewTable(w, examples.MarkdownOptionRowspan...)
		if err := t.Render(examples.RowspanData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarkdownTableRowspanReuse(b *testing.B) {
	w := &bytes.Buffer{}
	t := markdown.NewTable(w, examples.MarkdownOptionRowspan...)
	for b.Loop() {
		w.Reset()
		if err := t.Render(examples.RowspanData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarkdownStreamRowspan(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		s := markdown.NewStream(w, examples.MarkdownOptionRowspan...)
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

func BenchmarkMarkdownTableColspanFresh(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		t := markdown.NewTable(w, examples.MarkdownOptionColspan...)
		if err := t.Render(examples.ColspanData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarkdownTableColspanReuse(b *testing.B) {
	w := &bytes.Buffer{}
	t := markdown.NewTable(w, examples.MarkdownOptionColspan...)
	for b.Loop() {
		w.Reset()
		if err := t.Render(examples.ColspanData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarkdownStreamColspan(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		s := markdown.NewStream(w, examples.MarkdownOptionColspan...)
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

func BenchmarkMarkdownTableTransformerFresh(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		t := markdown.NewTable(w, examples.MarkdownOptionTransformer...)
		if err := t.Render(examples.FooterData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarkdownTableTransformerReuse(b *testing.B) {
	w := &bytes.Buffer{}
	t := markdown.NewTable(w, examples.MarkdownOptionTransformer...)
	for b.Loop() {
		w.Reset()
		if err := t.Render(examples.FooterData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarkdownStreamTransformer(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		s := markdown.NewStream(w, examples.MarkdownOptionTransformer...)
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

func BenchmarkMarkdownTableComplexFresh(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		t := markdown.NewTable(w, examples.MarkdownOptionComplex...)
		if err := t.Render(examples.ComplexData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarkdownTableComplexReuse(b *testing.B) {
	w := &bytes.Buffer{}
	t := markdown.NewTable(w, examples.MarkdownOptionComplex...)
	for b.Loop() {
		w.Reset()
		if err := t.Render(examples.ComplexData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarkdownStreamComplex(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		s := markdown.NewStream(w, examples.MarkdownOptionComplex...)
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

func BenchmarkMarkdownTableValueInputMixed(b *testing.B) {
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
	t := markdown.NewTable(io.Discard, markdown.WithHeader(header))
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

func BenchmarkMarkdownStreamValueInputMixed(b *testing.B) {
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
		s := markdown.NewStream(io.Discard, markdown.WithHeader(header))
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

func BenchmarkMarkdownTableValueInputSmallInts(b *testing.B) {
	records := make([]int, 1000)
	for i := range records {
		records[i] = 1000 + i
	}
	header := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	t := markdown.NewTable(io.Discard, markdown.WithHeader(header))
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

func BenchmarkMarkdownStreamValueInputSmallInts(b *testing.B) {
	records := make([]int, 1000)
	for i := range records {
		records[i] = 1000 + i
	}
	header := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	row := make([]table.Value, 0, 8)
	b.ReportAllocs()
	for b.Loop() {
		s := markdown.NewStream(io.Discard, markdown.WithHeader(header))
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
