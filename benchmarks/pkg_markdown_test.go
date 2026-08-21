package benchmarks

import (
	"bytes"
	"testing"

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
