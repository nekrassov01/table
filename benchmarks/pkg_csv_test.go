package benchmarks

import (
	"bytes"
	"testing"

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
