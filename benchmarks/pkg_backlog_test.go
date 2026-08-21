package benchmarks

import (
	"bytes"
	"testing"

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
