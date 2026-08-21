package benchmarks

import (
	"bytes"
	"testing"

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
