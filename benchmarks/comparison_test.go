// The comparison benchmarks render the examples datasets with the options
// below chosen for parity, not with raw defaults. Content (header,
// rows, footer) is passed wherever the library can express it; mintab
// has no footer and rejects nested values, so the Footer scenario
// omits its footer and the Complex scenario skips it. Every library
// draws an outer border and a header separator without per-row
// separators, header auto-casing is turned off in go-pretty and
// tablewriter so all four render the input verbatim, and merging and
// colors stay off. Border charsets keep each library's default look:
// Unicode for table and tablewriter, ASCII for mintab and go-pretty.

package benchmarks

import (
	"bytes"
	"io"
	"testing"

	goprettytable "github.com/jedib0t/go-pretty/v6/table"
	goprettytext "github.com/jedib0t/go-pretty/v6/text"
	"github.com/nekrassov01/mintab"
	"github.com/nekrassov01/table/examples"
	"github.com/nekrassov01/table/text"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
)

var tablewriterOpts = []tablewriter.Option{
	tablewriter.WithHeaderAutoFormat(tw.Off),
	tablewriter.WithFooterAutoFormat(tw.Off),
}

func BenchmarkComparisonTableSimple(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		t := text.NewTable(w,
			text.WithHeader(examples.SimpleData.Header...),
			text.WithStyle(text.StyleLight),
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
		t := tablewriter.NewTable(w, tablewriterOpts...)
		t.Header(examples.SimpleData.Header[0])
		if err := t.Bulk(examples.SimpleData.Body); err != nil {
			b.Fatal(err)
		}
		t.Render()
	}
}

func BenchmarkComparisonTableFooter(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		t := text.NewTable(w,
			text.WithHeader(examples.FooterData.Header...),
			text.WithFooter(examples.FooterData.Footer),
			text.WithStyle(text.StyleLight),
			text.WithCompact(),
		)
		if err := t.Render(examples.FooterData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkComparisonMintabFooter(b *testing.B) {
	w := &bytes.Buffer{}
	data := mintab.Input{
		Header: examples.FooterData.Header[0],
		Data:   examples.FooterData.Body,
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

func BenchmarkComparisonGoPrettyFooter(b *testing.B) {
	w := &bytes.Buffer{}
	h := goprettyRow(examples.FooterData.Header[0])
	f := goprettyRow(examples.FooterData.Footer()[0])
	r := goprettyRows(examples.FooterData.Body)
	for b.Loop() {
		w.Reset()
		t := goprettyWriter(w)
		t.AppendHeader(h)
		t.AppendRows(r)
		t.AppendFooter(f)
		t.Render()
	}
}

func BenchmarkComparisonTableWriterFooter(b *testing.B) {
	w := &bytes.Buffer{}
	f := examples.FooterData.Footer()[0]
	for b.Loop() {
		w.Reset()
		t := tablewriter.NewTable(w, tablewriterOpts...)
		t.Header(examples.FooterData.Header[0])
		if err := t.Bulk(examples.FooterData.Body); err != nil {
			b.Fatal(err)
		}
		t.Footer(f)
		t.Render()
	}
}

func BenchmarkComparisonTableCompact(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		t := text.NewTable(w,
			text.WithHeader(examples.CompactData.Header...),
			text.WithStyle(text.StyleLight),
			text.WithCompact(),
		)
		if err := t.Render(examples.CompactData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkComparisonMintabCompact(b *testing.B) {
	w := &bytes.Buffer{}
	data := mintab.Input{
		Header: examples.CompactData.Header[0],
		Data:   examples.CompactData.Body,
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

func BenchmarkComparisonGoPrettyCompact(b *testing.B) {
	w := &bytes.Buffer{}
	h := goprettyRow(examples.CompactData.Header[0])
	r := goprettyRows(examples.CompactData.Body)
	for b.Loop() {
		w.Reset()
		t := goprettyWriter(w)
		t.AppendHeader(h)
		t.AppendRows(r)
		t.Render()
	}
}

func BenchmarkComparisonTableWriterCompact(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		t := tablewriter.NewTable(w, tablewriterOpts...)
		t.Header(examples.CompactData.Header[0])
		if err := t.Bulk(examples.CompactData.Body); err != nil {
			b.Fatal(err)
		}
		t.Render()
	}
}

func BenchmarkComparisonTableRowspan(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		t := text.NewTable(w,
			text.WithHeader(examples.RowspanData.Header...),
			text.WithStyle(text.StyleLight),
			text.WithCompact(),
		)
		if err := t.Render(examples.RowspanData.Body); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkComparisonMintabRowspan(b *testing.B) {
	w := &bytes.Buffer{}
	data := mintab.Input{
		Header: examples.RowspanData.Header[0],
		Data:   examples.RowspanData.Body,
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

func BenchmarkComparisonGoPrettyRowspan(b *testing.B) {
	w := &bytes.Buffer{}
	h := goprettyRow(examples.RowspanData.Header[0])
	r := goprettyRows(examples.RowspanData.Body)
	for b.Loop() {
		w.Reset()
		t := goprettyWriter(w)
		t.AppendHeader(h)
		t.AppendRows(r)
		t.Render()
	}
}

func BenchmarkComparisonTableWriterRowspan(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		t := tablewriter.NewTable(w, tablewriterOpts...)
		t.Header(examples.RowspanData.Header[0])
		if err := t.Bulk(examples.RowspanData.Body); err != nil {
			b.Fatal(err)
		}
		t.Render()
	}
}

func BenchmarkComparisonTableComplex(b *testing.B) {
	w := &bytes.Buffer{}
	for b.Loop() {
		w.Reset()
		t := text.NewTable(w,
			text.WithHeader(examples.ComplexData.Header...),
			text.WithStyle(text.StyleLight),
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
		t := tablewriter.NewTable(w, tablewriterOpts...)
		t.Header(examples.ComplexData.Header[0])
		if err := t.Bulk(examples.ComplexData.Body); err != nil {
			b.Fatal(err)
		}
		t.Render()
	}
}

// goprettyRow converts a header or footer to a go-pretty row.
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

// goprettyWriter returns a writer mirrored to w with header and footer
// auto-casing off, so the content matches the input verbatim as the
// other libraries render it.
func goprettyWriter(w io.Writer) goprettytable.Writer {
	t := goprettytable.NewWriter()
	t.SetOutputMirror(w)
	t.Style().Format.Header = goprettytext.FormatDefault
	t.Style().Format.Footer = goprettytext.FormatDefault
	return t
}
