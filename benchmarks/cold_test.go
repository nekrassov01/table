//go:build cold

// The Cold benchmarks measure the memory one render actually demands
// on the 100x datasets. The steady-state benchmarks amortize to ~0 B/op
// because the pooled arenas retain their high-water capacity between
// iterations, so each iteration here drains the pools first. Batch
// rendering builds cells for every row at once, while streaming holds
// one row at a time; the B/op gap between the two is the point. The GC
// cycles make these heavy, hence the build tag: run them with
// "make bench target=cold".

package benchmarks

import (
	"bytes"
	"runtime"
	"testing"

	"github.com/nekrassov01/table/backlog"
	"github.com/nekrassov01/table/csv"
	"github.com/nekrassov01/table/examples"
	"github.com/nekrassov01/table/html"
	"github.com/nekrassov01/table/markdown"
	"github.com/nekrassov01/table/text"
)

func BenchmarkTextTableSimpleCold(b *testing.B) {
	w := &bytes.Buffer{}
	t := text.NewTable(w, examples.TextOptionSimple...)
	if err := t.Render(examples.SimpleDataLarge.Body); err != nil {
		b.Fatal(err)
	}
	for b.Loop() {
		drain(b, w)
		t := text.NewTable(w, examples.TextOptionSimple...)
		if err := t.Render(examples.SimpleDataLarge.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTextStreamSimpleCold(b *testing.B) {
	w := &bytes.Buffer{}
	s := text.NewStream(w, examples.TextOptionSimple...)
	for _, row := range examples.SimpleDataLarge.Body {
		if err := s.Render(row); err != nil {
			b.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		b.Fatal(err)
	}
	for b.Loop() {
		drain(b, w)
		s := text.NewStream(w, examples.TextOptionSimple...)
		for _, row := range examples.SimpleDataLarge.Body {
			if err := s.Render(row); err != nil {
				b.Fatal(err)
			}
		}
		if err := s.Close(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTextTableFooterCold(b *testing.B) {
	w := &bytes.Buffer{}
	t := text.NewTable(w, examples.TextOptionFooter...)
	if err := t.Render(examples.FooterDataLarge.Body); err != nil {
		b.Fatal(err)
	}
	for b.Loop() {
		drain(b, w)
		t := text.NewTable(w, examples.TextOptionFooter...)
		if err := t.Render(examples.FooterDataLarge.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTextStreamFooterCold(b *testing.B) {
	w := &bytes.Buffer{}
	s := text.NewStream(w, examples.TextOptionFooter...)
	for _, row := range examples.FooterDataLarge.Body {
		if err := s.Render(row); err != nil {
			b.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		b.Fatal(err)
	}
	for b.Loop() {
		drain(b, w)
		s := text.NewStream(w, examples.TextOptionFooter...)
		for _, row := range examples.FooterDataLarge.Body {
			if err := s.Render(row); err != nil {
				b.Fatal(err)
			}
		}
		if err := s.Close(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTextTableCompactCold(b *testing.B) {
	w := &bytes.Buffer{}
	t := text.NewTable(w, examples.TextOptionCompact...)
	if err := t.Render(examples.CompactDataLarge.Body); err != nil {
		b.Fatal(err)
	}
	for b.Loop() {
		drain(b, w)
		t := text.NewTable(w, examples.TextOptionCompact...)
		if err := t.Render(examples.CompactDataLarge.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTextStreamCompactCold(b *testing.B) {
	w := &bytes.Buffer{}
	s := text.NewStream(w, examples.TextOptionCompact...)
	for _, row := range examples.CompactDataLarge.Body {
		if err := s.Render(row); err != nil {
			b.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		b.Fatal(err)
	}
	for b.Loop() {
		drain(b, w)
		s := text.NewStream(w, examples.TextOptionCompact...)
		for _, row := range examples.CompactDataLarge.Body {
			if err := s.Render(row); err != nil {
				b.Fatal(err)
			}
		}
		if err := s.Close(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTextTableRowspanCold(b *testing.B) {
	w := &bytes.Buffer{}
	t := text.NewTable(w, examples.TextOptionRowspan...)
	if err := t.Render(examples.RowspanDataLarge.Body); err != nil {
		b.Fatal(err)
	}
	for b.Loop() {
		drain(b, w)
		t := text.NewTable(w, examples.TextOptionRowspan...)
		if err := t.Render(examples.RowspanDataLarge.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTextStreamRowspanCold(b *testing.B) {
	w := &bytes.Buffer{}
	s := text.NewStream(w, examples.TextOptionRowspan...)
	for _, row := range examples.RowspanDataLarge.Body {
		if err := s.Render(row); err != nil {
			b.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		b.Fatal(err)
	}
	for b.Loop() {
		drain(b, w)
		s := text.NewStream(w, examples.TextOptionRowspan...)
		for _, row := range examples.RowspanDataLarge.Body {
			if err := s.Render(row); err != nil {
				b.Fatal(err)
			}
		}
		if err := s.Close(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTextTableComplexCold(b *testing.B) {
	w := &bytes.Buffer{}
	t := text.NewTable(w, examples.TextOptionComplex...)
	if err := t.Render(examples.ComplexDataLarge.Body); err != nil {
		b.Fatal(err)
	}
	for b.Loop() {
		drain(b, w)
		t := text.NewTable(w, examples.TextOptionComplex...)
		if err := t.Render(examples.ComplexDataLarge.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTextStreamComplexCold(b *testing.B) {
	w := &bytes.Buffer{}
	s := text.NewStream(w, examples.TextOptionComplex...)
	for _, row := range examples.ComplexDataLarge.Body {
		if err := s.Render(row); err != nil {
			b.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		b.Fatal(err)
	}
	for b.Loop() {
		drain(b, w)
		s := text.NewStream(w, examples.TextOptionComplex...)
		for _, row := range examples.ComplexDataLarge.Body {
			if err := s.Render(row); err != nil {
				b.Fatal(err)
			}
		}
		if err := s.Close(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarkdownTableSimpleCold(b *testing.B) {
	w := &bytes.Buffer{}
	t := markdown.NewTable(w, examples.MarkdownOptionSimple...)
	if err := t.Render(examples.SimpleDataLarge.Body); err != nil {
		b.Fatal(err)
	}
	for b.Loop() {
		drain(b, w)
		t := markdown.NewTable(w, examples.MarkdownOptionSimple...)
		if err := t.Render(examples.SimpleDataLarge.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarkdownStreamSimpleCold(b *testing.B) {
	w := &bytes.Buffer{}
	s := markdown.NewStream(w, examples.MarkdownOptionSimple...)
	for _, row := range examples.SimpleDataLarge.Body {
		if err := s.Render(row); err != nil {
			b.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		b.Fatal(err)
	}
	for b.Loop() {
		drain(b, w)
		s := markdown.NewStream(w, examples.MarkdownOptionSimple...)
		for _, row := range examples.SimpleDataLarge.Body {
			if err := s.Render(row); err != nil {
				b.Fatal(err)
			}
		}
		if err := s.Close(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarkdownTableRowspanCold(b *testing.B) {
	w := &bytes.Buffer{}
	t := markdown.NewTable(w, examples.MarkdownOptionRowspan...)
	if err := t.Render(examples.RowspanDataLarge.Body); err != nil {
		b.Fatal(err)
	}
	for b.Loop() {
		drain(b, w)
		t := markdown.NewTable(w, examples.MarkdownOptionRowspan...)
		if err := t.Render(examples.RowspanDataLarge.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarkdownStreamRowspanCold(b *testing.B) {
	w := &bytes.Buffer{}
	s := markdown.NewStream(w, examples.MarkdownOptionRowspan...)
	for _, row := range examples.RowspanDataLarge.Body {
		if err := s.Render(row); err != nil {
			b.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		b.Fatal(err)
	}
	for b.Loop() {
		drain(b, w)
		s := markdown.NewStream(w, examples.MarkdownOptionRowspan...)
		for _, row := range examples.RowspanDataLarge.Body {
			if err := s.Render(row); err != nil {
				b.Fatal(err)
			}
		}
		if err := s.Close(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarkdownTableComplexCold(b *testing.B) {
	w := &bytes.Buffer{}
	t := markdown.NewTable(w, examples.MarkdownOptionComplex...)
	if err := t.Render(examples.ComplexDataLarge.Body); err != nil {
		b.Fatal(err)
	}
	for b.Loop() {
		drain(b, w)
		t := markdown.NewTable(w, examples.MarkdownOptionComplex...)
		if err := t.Render(examples.ComplexDataLarge.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarkdownStreamComplexCold(b *testing.B) {
	w := &bytes.Buffer{}
	s := markdown.NewStream(w, examples.MarkdownOptionComplex...)
	for _, row := range examples.ComplexDataLarge.Body {
		if err := s.Render(row); err != nil {
			b.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		b.Fatal(err)
	}
	for b.Loop() {
		drain(b, w)
		s := markdown.NewStream(w, examples.MarkdownOptionComplex...)
		for _, row := range examples.ComplexDataLarge.Body {
			if err := s.Render(row); err != nil {
				b.Fatal(err)
			}
		}
		if err := s.Close(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBacklogTableSimpleCold(b *testing.B) {
	w := &bytes.Buffer{}
	t := backlog.NewTable(w, examples.BacklogOptionSimple...)
	if err := t.Render(examples.SimpleDataLarge.Body); err != nil {
		b.Fatal(err)
	}
	for b.Loop() {
		drain(b, w)
		t := backlog.NewTable(w, examples.BacklogOptionSimple...)
		if err := t.Render(examples.SimpleDataLarge.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBacklogStreamSimpleCold(b *testing.B) {
	w := &bytes.Buffer{}
	s := backlog.NewStream(w, examples.BacklogOptionSimple...)
	for _, row := range examples.SimpleDataLarge.Body {
		if err := s.Render(row); err != nil {
			b.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		b.Fatal(err)
	}
	for b.Loop() {
		drain(b, w)
		s := backlog.NewStream(w, examples.BacklogOptionSimple...)
		for _, row := range examples.SimpleDataLarge.Body {
			if err := s.Render(row); err != nil {
				b.Fatal(err)
			}
		}
		if err := s.Close(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBacklogTableRowspanCold(b *testing.B) {
	w := &bytes.Buffer{}
	t := backlog.NewTable(w, examples.BacklogOptionRowspan...)
	if err := t.Render(examples.RowspanDataLarge.Body); err != nil {
		b.Fatal(err)
	}
	for b.Loop() {
		drain(b, w)
		t := backlog.NewTable(w, examples.BacklogOptionRowspan...)
		if err := t.Render(examples.RowspanDataLarge.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBacklogStreamRowspanCold(b *testing.B) {
	w := &bytes.Buffer{}
	s := backlog.NewStream(w, examples.BacklogOptionRowspan...)
	for _, row := range examples.RowspanDataLarge.Body {
		if err := s.Render(row); err != nil {
			b.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		b.Fatal(err)
	}
	for b.Loop() {
		drain(b, w)
		s := backlog.NewStream(w, examples.BacklogOptionRowspan...)
		for _, row := range examples.RowspanDataLarge.Body {
			if err := s.Render(row); err != nil {
				b.Fatal(err)
			}
		}
		if err := s.Close(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBacklogTableComplexCold(b *testing.B) {
	w := &bytes.Buffer{}
	t := backlog.NewTable(w, examples.BacklogOptionComplex...)
	if err := t.Render(examples.ComplexDataLarge.Body); err != nil {
		b.Fatal(err)
	}
	for b.Loop() {
		drain(b, w)
		t := backlog.NewTable(w, examples.BacklogOptionComplex...)
		if err := t.Render(examples.ComplexDataLarge.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBacklogStreamComplexCold(b *testing.B) {
	w := &bytes.Buffer{}
	s := backlog.NewStream(w, examples.BacklogOptionComplex...)
	for _, row := range examples.ComplexDataLarge.Body {
		if err := s.Render(row); err != nil {
			b.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		b.Fatal(err)
	}
	for b.Loop() {
		drain(b, w)
		s := backlog.NewStream(w, examples.BacklogOptionComplex...)
		for _, row := range examples.ComplexDataLarge.Body {
			if err := s.Render(row); err != nil {
				b.Fatal(err)
			}
		}
		if err := s.Close(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCSVTableSimpleCold(b *testing.B) {
	w := &bytes.Buffer{}
	t := csv.NewTable(w, examples.CSVOptionSimple...)
	if err := t.Render(examples.SimpleDataLarge.Body); err != nil {
		b.Fatal(err)
	}
	for b.Loop() {
		drain(b, w)
		t := csv.NewTable(w, examples.CSVOptionSimple...)
		if err := t.Render(examples.SimpleDataLarge.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCSVStreamSimpleCold(b *testing.B) {
	w := &bytes.Buffer{}
	s := csv.NewStream(w, examples.CSVOptionSimple...)
	for _, row := range examples.SimpleDataLarge.Body {
		if err := s.Render(row); err != nil {
			b.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		b.Fatal(err)
	}
	for b.Loop() {
		drain(b, w)
		s := csv.NewStream(w, examples.CSVOptionSimple...)
		for _, row := range examples.SimpleDataLarge.Body {
			if err := s.Render(row); err != nil {
				b.Fatal(err)
			}
		}
		if err := s.Close(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCSVTableComplexCold(b *testing.B) {
	w := &bytes.Buffer{}
	t := csv.NewTable(w, examples.CSVOptionComplex...)
	if err := t.Render(examples.ComplexDataLarge.Body); err != nil {
		b.Fatal(err)
	}
	for b.Loop() {
		drain(b, w)
		t := csv.NewTable(w, examples.CSVOptionComplex...)
		if err := t.Render(examples.ComplexDataLarge.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCSVStreamComplexCold(b *testing.B) {
	w := &bytes.Buffer{}
	s := csv.NewStream(w, examples.CSVOptionComplex...)
	for _, row := range examples.ComplexDataLarge.Body {
		if err := s.Render(row); err != nil {
			b.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		b.Fatal(err)
	}
	for b.Loop() {
		drain(b, w)
		s := csv.NewStream(w, examples.CSVOptionComplex...)
		for _, row := range examples.ComplexDataLarge.Body {
			if err := s.Render(row); err != nil {
				b.Fatal(err)
			}
		}
		if err := s.Close(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCSVTableCommaIncludedCold(b *testing.B) {
	w := &bytes.Buffer{}
	t := csv.NewTable(w, examples.CSVOptionCommaIncluded...)
	if err := t.Render(examples.CommaIncludedDataLarge.Body); err != nil {
		b.Fatal(err)
	}
	for b.Loop() {
		drain(b, w)
		t := csv.NewTable(w, examples.CSVOptionCommaIncluded...)
		if err := t.Render(examples.CommaIncludedDataLarge.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCSVStreamCommaIncludedCold(b *testing.B) {
	w := &bytes.Buffer{}
	s := csv.NewStream(w, examples.CSVOptionCommaIncluded...)
	for _, row := range examples.CommaIncludedDataLarge.Body {
		if err := s.Render(row); err != nil {
			b.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		b.Fatal(err)
	}
	for b.Loop() {
		drain(b, w)
		s := csv.NewStream(w, examples.CSVOptionCommaIncluded...)
		for _, row := range examples.CommaIncludedDataLarge.Body {
			if err := s.Render(row); err != nil {
				b.Fatal(err)
			}
		}
		if err := s.Close(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkHTMLTableSimpleCold(b *testing.B) {
	w := &bytes.Buffer{}
	t := html.NewTable(w, examples.HTMLOptionSimple...)
	if err := t.Render(examples.SimpleDataLarge.Body); err != nil {
		b.Fatal(err)
	}
	for b.Loop() {
		drain(b, w)
		t := html.NewTable(w, examples.HTMLOptionSimple...)
		if err := t.Render(examples.SimpleDataLarge.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkHTMLStreamSimpleCold(b *testing.B) {
	w := &bytes.Buffer{}
	s := html.NewStream(w, examples.HTMLOptionSimple...)
	for _, row := range examples.SimpleDataLarge.Body {
		if err := s.Render(row); err != nil {
			b.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		b.Fatal(err)
	}
	for b.Loop() {
		drain(b, w)
		s := html.NewStream(w, examples.HTMLOptionSimple...)
		for _, row := range examples.SimpleDataLarge.Body {
			if err := s.Render(row); err != nil {
				b.Fatal(err)
			}
		}
		if err := s.Close(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkHTMLTableRowspanCold(b *testing.B) {
	w := &bytes.Buffer{}
	t := html.NewTable(w, examples.HTMLOptionRowspan...)
	if err := t.Render(examples.RowspanDataLarge.Body); err != nil {
		b.Fatal(err)
	}
	for b.Loop() {
		drain(b, w)
		t := html.NewTable(w, examples.HTMLOptionRowspan...)
		if err := t.Render(examples.RowspanDataLarge.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkHTMLStreamRowspanCold(b *testing.B) {
	w := &bytes.Buffer{}
	s := html.NewStream(w, examples.HTMLOptionRowspan...)
	for _, row := range examples.RowspanDataLarge.Body {
		if err := s.Render(row); err != nil {
			b.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		b.Fatal(err)
	}
	for b.Loop() {
		drain(b, w)
		s := html.NewStream(w, examples.HTMLOptionRowspan...)
		for _, row := range examples.RowspanDataLarge.Body {
			if err := s.Render(row); err != nil {
				b.Fatal(err)
			}
		}
		if err := s.Close(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkHTMLTableColspanCold(b *testing.B) {
	w := &bytes.Buffer{}
	t := html.NewTable(w, examples.HTMLOptionColspan...)
	if err := t.Render(examples.ColspanDataLarge.Body); err != nil {
		b.Fatal(err)
	}
	for b.Loop() {
		drain(b, w)
		t := html.NewTable(w, examples.HTMLOptionColspan...)
		if err := t.Render(examples.ColspanDataLarge.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkHTMLStreamColspanCold(b *testing.B) {
	w := &bytes.Buffer{}
	s := html.NewStream(w, examples.HTMLOptionColspan...)
	for _, row := range examples.ColspanDataLarge.Body {
		if err := s.Render(row); err != nil {
			b.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		b.Fatal(err)
	}
	for b.Loop() {
		drain(b, w)
		s := html.NewStream(w, examples.HTMLOptionColspan...)
		for _, row := range examples.ColspanDataLarge.Body {
			if err := s.Render(row); err != nil {
				b.Fatal(err)
			}
		}
		if err := s.Close(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkHTMLTableComplexCold(b *testing.B) {
	w := &bytes.Buffer{}
	t := html.NewTable(w, examples.HTMLOptionComplex...)
	if err := t.Render(examples.ComplexDataLarge.Body); err != nil {
		b.Fatal(err)
	}
	for b.Loop() {
		drain(b, w)
		t := html.NewTable(w, examples.HTMLOptionComplex...)
		if err := t.Render(examples.ComplexDataLarge.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkHTMLStreamComplexCold(b *testing.B) {
	w := &bytes.Buffer{}
	s := html.NewStream(w, examples.HTMLOptionComplex...)
	for _, row := range examples.ComplexDataLarge.Body {
		if err := s.Render(row); err != nil {
			b.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		b.Fatal(err)
	}
	for b.Loop() {
		drain(b, w)
		s := html.NewStream(w, examples.HTMLOptionComplex...)
		for _, row := range examples.ComplexDataLarge.Body {
			if err := s.Render(row); err != nil {
				b.Fatal(err)
			}
		}
		if err := s.Close(); err != nil {
			b.Fatal(err)
		}
	}
}

// drain resets the writer and drops the pooled arenas while the timer
// is stopped, so allocations here stay out of the measurement;
// sync.Pool keeps a victim cache, hence two GC cycles.
func drain(b *testing.B, w *bytes.Buffer) {
	b.StopTimer()
	w.Reset()
	runtime.GC()
	runtime.GC()
	b.StartTimer()
}
