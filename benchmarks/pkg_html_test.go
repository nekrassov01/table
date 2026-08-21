package benchmarks

import (
	"bytes"
	"testing"

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
