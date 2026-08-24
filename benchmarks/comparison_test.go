// The comparison benchmarks render the shared example data sets. Options are
// selected for parity rather than using each library's raw defaults.
//
// Every library receives the same header and body. Mintab rejects the nested
// values in Complex, so that scenario omits it. Each output has an outer border
// and a header separator without per-row separators. Header auto-casing remains
// off, and border characters retain each library's default style: Unicode for
// table and tablewriter, and ASCII for mintab, simpletable, and go-pretty. Value
// formatting otherwise retains each library's default behavior.
//
// Static data is converted to each library's required row type before timing.
// Construction, row ingestion, and output remain inside each loop.

package benchmarks

import (
	"bytes"
	"io"
	"testing"

	"github.com/alexeyco/simpletable"
	goprettytable "github.com/jedib0t/go-pretty/v6/table"
	goprettytext "github.com/jedib0t/go-pretty/v6/text"
	"github.com/nekrassov01/mintab"
	"github.com/nekrassov01/table/examples"
	"github.com/nekrassov01/table/text"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
)

func BenchmarkComparisonTableSimple(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		t := text.NewTable(w,
			text.WithHeader(examples.SimpleData.Header...),
			text.WithCompact(),
		)
		if err := t.Render(examples.SimpleData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkComparisonMintabSimple(b *testing.B) {
	w := &bytes.Buffer{}
	data := mintab.Input{
		Header: examples.SimpleData.Header[0],
		Data:   examples.SimpleData.Body,
	}
	for b.Loop() {
		w.Reset()
		t := mintab.New(w,
			mintab.WithFormat(mintab.CompressedTextFormat),
		)
		if err := t.Load(data); err != nil {
			b.Fatal(err)
		}
		t.Render()
	}
}

func BenchmarkComparisonSimpleTableSimple(b *testing.B) {
	w := &bytes.Buffer{}
	h := simpletableHeader(examples.SimpleData.Header[0])
	r := simpletableRows(examples.SimpleData.Body)
	for b.Loop() {
		w.Reset()
		t := simpletable.New()
		t.Header.Cells = h
		t.Body.Cells = r
		_, _ = w.WriteString(t.String())
		_ = w.WriteByte('\n')
	}
}

func BenchmarkComparisonGoPrettySimple(b *testing.B) {
	w := &bytes.Buffer{}
	h := goprettyRow(examples.SimpleData.Header[0])
	r := goprettyRows(examples.SimpleData.Body)
	for b.Loop() {
		w.Reset()
		t := goprettyWriter(w)
		t.AppendHeader(h)
		t.AppendRows(r)
		t.Render()
	}
}

func BenchmarkComparisonTableWriterSimple(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		t := tablewriter.NewTable(w, tablewriter.WithHeaderAutoFormat(tw.Off))
		t.Header(examples.SimpleData.Header[0])
		if err := t.Bulk(examples.SimpleData.Body); err != nil {
			b.Fatal(err)
		}
		if err := t.Render(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkComparisonTableComplex(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		t := text.NewTable(w,
			text.WithHeader(examples.ComplexData.Header...),
			text.WithCompact(),
		)
		if err := t.Render(examples.ComplexData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkComparisonGoPrettyComplex(b *testing.B) {
	w := &bytes.Buffer{}
	h := goprettyRow(examples.ComplexData.Header[0])
	r := goprettyRows(examples.ComplexData.Body)
	for b.Loop() {
		w.Reset()
		t := goprettyWriter(w)
		t.AppendHeader(h)
		t.AppendRows(r)
		t.Render()
	}
}

func BenchmarkComparisonTableWriterComplex(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		t := tablewriter.NewTable(w, tablewriter.WithHeaderAutoFormat(tw.Off))
		t.Header(examples.ComplexData.Header[0])
		if err := t.Bulk(examples.ComplexData.Body); err != nil {
			b.Fatal(err)
		}
		if err := t.Render(); err != nil {
			b.Fatal(err)
		}
	}
}

// simpletableHeader converts string cells to a simpletable header.
func simpletableHeader(cells []string) []*simpletable.Cell {
	h := make([]*simpletable.Cell, len(cells))
	for index, cell := range cells {
		h[index] = &simpletable.Cell{Text: cell}
	}
	return h
}

// simpletableRows converts string-valued data rows to simpletable rows.
func simpletableRows(rows [][]any) [][]*simpletable.Cell {
	r := make([][]*simpletable.Cell, len(rows))
	for rowIndex, row := range rows {
		cells := make([]*simpletable.Cell, len(row))
		for columnIndex, cell := range row {
			cells[columnIndex] = &simpletable.Cell{Text: cell.(string)}
		}
		r[rowIndex] = cells
	}
	return r
}

// goprettyRow converts string cells to a go-pretty row.
func goprettyRow(cells []string) goprettytable.Row {
	r := make(goprettytable.Row, len(cells))
	for i, cell := range cells {
		r[i] = cell
	}
	return r
}

// goprettyRows converts data rows to go-pretty rows.
func goprettyRows(rows [][]any) []goprettytable.Row {
	r := make([]goprettytable.Row, len(rows))
	for i, row := range rows {
		r[i] = goprettytable.Row(row)
	}
	return r
}

// goprettyWriter returns a writer mirrored to w with header auto-casing off,
// so the content matches the input verbatim as the other libraries render it.
func goprettyWriter(w io.Writer) goprettytable.Writer {
	t := goprettytable.NewWriter()
	t.SetOutputMirror(w)
	t.Style().Format.Header = goprettytext.FormatDefault
	return t
}
