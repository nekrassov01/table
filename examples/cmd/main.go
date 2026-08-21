// Command cmd renders a selected table example.
package main

import (
	"fmt"
	"io"
	"os"
	"slices"

	"github.com/nekrassov01/table"
	"github.com/nekrassov01/table/backlog"
	"github.com/nekrassov01/table/csv"
	"github.com/nekrassov01/table/examples"
	"github.com/nekrassov01/table/html"
	"github.com/nekrassov01/table/markdown"
	"github.com/nekrassov01/table/text"
)

var (
	targetAll      = "all"
	targetText     = "text"
	targetHTML     = "html"
	targetMarkdown = "markdown"
	targetBacklog  = "backlog"
	targetCSV      = "csv"
)

var (
	modeTable  = "table"
	modeStream = "stream"
)

var (
	dataASCII         = "ascii"
	dataSimple        = "simple"
	dataCompact       = "compact"
	dataRowspan       = "rowspan"
	dataColspan       = "colspan"
	dataFooter        = "footer"
	dataTransformer   = "transformer"
	dataComplex       = "complex"
	dataStackedHeader = "stacked-header"
	dataCommaIncluded = "comma-included"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 || len(args) > 3 {
		fmt.Fprintln(os.Stderr, "usage: go run ./examples/cmd <target> [mode] [data]")
		os.Exit(1)
	}
	if err := newRunner(os.Stdout, args...).run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type runner struct {
	w      io.Writer
	target string
	mode   string
	data   string
}

func newRunner(w io.Writer, args ...string) runner {
	o := runner{
		w:      w,
		target: args[0],
	}
	if len(args) > 1 {
		o.mode = args[1]
	}
	if len(args) > 2 {
		o.data = args[2]
	}
	return o
}

func (o runner) run() error {
	modes := []string{modeTable, modeStream}
	if o.mode != "" {
		if !slices.Contains(modes, o.mode) {
			return fmt.Errorf("examples: unknown mode %q", o.mode)
		}
		modes = []string{o.mode}
	}
	allTargets := o.target == targetAll
	targets := []string{o.target}
	if allTargets {
		targets = []string{
			targetText,
			targetHTML,
			targetMarkdown,
			targetBacklog,
			targetCSV,
		}
	} else if _, ok := dataNames(o.target); !ok {
		return fmt.Errorf("examples: unknown target %q", o.target)
	}
	data := o.data
	if data != "" {
		if _, ok := exampleData(data); !ok {
			return fmt.Errorf("examples: unknown data %q", data)
		}
	}
	first := true
	for _, target := range targets {
		names, _ := dataNames(target)
		if data != "" {
			if !slices.Contains(names, data) {
				if allTargets {
					continue
				}
				return newError(data, target)
			}
			names = []string{data}
		}
		for _, mode := range modes {
			for _, name := range names {
				if !first {
					if _, err := io.WriteString(o.w, "\n"); err != nil {
						return err
					}
				}
				o.target = target
				o.mode = mode
				o.data = name
				if err := o.runExample(); err != nil {
					return err
				}
				first = false
			}
		}
	}
	return nil
}

func (o runner) runExample() error {
	data, ok := exampleData(o.data)
	if !ok {
		return fmt.Errorf("examples: unknown data %q", o.data)
	}
	switch o.target {
	case targetText:
		return o.runText(data.Body)
	case targetHTML:
		return o.runHTML(data.Body)
	case targetMarkdown:
		return o.runMarkdown(data.Body)
	case targetBacklog:
		return o.runBacklog(data.Body)
	case targetCSV:
		return o.runCSV(data.Body)
	default:
		return fmt.Errorf("examples: unknown target %q", o.target)
	}
}

func (o runner) runText(rows [][]any) error {
	var opts []text.Option
	switch o.data {
	case dataASCII:
		opts = examples.TextOptionASCII
	case dataSimple:
		opts = examples.TextOptionSimple
	case dataCompact:
		opts = examples.TextOptionCompact
	case dataRowspan:
		opts = examples.TextOptionRowspan
	case dataColspan:
		opts = examples.TextOptionColspan
	case dataFooter:
		opts = examples.TextOptionFooter
	case dataTransformer:
		opts = examples.TextOptionTransformer
	case dataComplex:
		opts = examples.TextOptionComplex
	case dataStackedHeader:
		opts = examples.TextOptionStackedHeader
	default:
		return newError(o.data, o.target)
	}
	example := newExample(rows)
	if o.mode == modeTable {
		example.tabular = text.NewTable(o.w, opts...)
	} else {
		example.streamer = text.NewStream(o.w, opts...)
	}
	return example.run()
}

func (o runner) runHTML(rows [][]any) error {
	var opts []html.Option
	switch o.data {
	case dataSimple:
		opts = examples.HTMLOptionSimple
	case dataRowspan:
		opts = examples.HTMLOptionRowspan
	case dataColspan:
		opts = examples.HTMLOptionColspan
	case dataFooter:
		opts = examples.HTMLOptionFooter
	case dataTransformer:
		opts = examples.HTMLOptionTransformer
	case dataComplex:
		opts = examples.HTMLOptionComplex
	case dataStackedHeader:
		opts = examples.HTMLOptionStackedHeader
	default:
		return newError(o.data, o.target)
	}
	example := newExample(rows)
	if o.mode == modeTable {
		example.tabular = html.NewTable(o.w, opts...)
	} else {
		example.streamer = html.NewStream(o.w, opts...)
	}
	return example.run()
}

func (o runner) runMarkdown(rows [][]any) error {
	var opts []markdown.Option
	switch o.data {
	case dataSimple:
		opts = examples.MarkdownOptionSimple
	case dataRowspan:
		opts = examples.MarkdownOptionRowspan
	case dataColspan:
		opts = examples.MarkdownOptionColspan
	case dataTransformer:
		opts = examples.MarkdownOptionTransformer
	case dataComplex:
		opts = examples.MarkdownOptionComplex
	default:
		return newError(o.data, o.target)
	}
	example := newExample(rows)
	if o.mode == modeTable {
		example.tabular = markdown.NewTable(o.w, opts...)
	} else {
		example.streamer = markdown.NewStream(o.w, opts...)
	}
	return example.run()
}

func (o runner) runBacklog(rows [][]any) error {
	var opts []backlog.Option
	switch o.data {
	case dataSimple:
		opts = examples.BacklogOptionSimple
	case dataRowspan:
		opts = examples.BacklogOptionRowspan
	case dataColspan:
		opts = examples.BacklogOptionColspan
	case dataFooter:
		opts = examples.BacklogOptionFooter
	case dataTransformer:
		opts = examples.BacklogOptionTransformer
	case dataComplex:
		opts = examples.BacklogOptionComplex
	case dataStackedHeader:
		opts = examples.BacklogOptionStackedHeader
	default:
		return newError(o.data, o.target)
	}
	example := newExample(rows)
	if o.mode == modeTable {
		example.tabular = backlog.NewTable(o.w, opts...)
	} else {
		example.streamer = backlog.NewStream(o.w, opts...)
	}
	return example.run()
}

func (o runner) runCSV(rows [][]any) error {
	var opts []csv.Option
	switch o.data {
	case dataSimple:
		opts = examples.CSVOptionSimple
	case dataFooter:
		opts = examples.CSVOptionFooter
	case dataTransformer:
		opts = examples.CSVOptionTransformer
	case dataComplex:
		opts = examples.CSVOptionComplex
	case dataCommaIncluded:
		opts = examples.CSVOptionCommaIncluded
	default:
		return newError(o.data, o.target)
	}
	example := newExample(rows)
	if o.mode == modeTable {
		example.tabular = csv.NewTable(o.w, opts...)
	} else {
		example.streamer = csv.NewStream(o.w, opts...)
	}
	return example.run()
}

type example struct {
	rows     [][]any
	tabular  table.Tabular
	streamer table.Streamer
}

func newExample(rows [][]any) example {
	return example{rows: rows}
}

func (o example) run() error {
	if o.tabular != nil {
		return o.tabular.Render(o.rows)
	}
	for _, row := range o.rows {
		if err := o.streamer.Render(row); err != nil {
			return err
		}
	}
	return o.streamer.Close()
}

func dataNames(target string) ([]string, bool) {
	switch target {
	case targetText:
		return []string{
			dataASCII,
			dataSimple,
			dataCompact,
			dataRowspan,
			dataColspan,
			dataFooter,
			dataTransformer,
			dataComplex,
			dataStackedHeader,
		}, true
	case targetHTML:
		return []string{
			dataSimple,
			dataRowspan,
			dataColspan,
			dataFooter,
			dataTransformer,
			dataComplex,
			dataStackedHeader,
		}, true
	case targetMarkdown:
		return []string{
			dataSimple,
			dataRowspan,
			dataColspan,
			dataTransformer,
			dataComplex,
		}, true
	case targetBacklog:
		return []string{
			dataSimple,
			dataRowspan,
			dataColspan,
			dataFooter,
			dataTransformer,
			dataComplex,
			dataStackedHeader,
		}, true
	case targetCSV:
		return []string{
			dataSimple,
			dataFooter,
			dataTransformer,
			dataComplex,
			dataCommaIncluded,
		}, true
	default:
		return nil, false
	}
}

func exampleData(name string) (examples.Data, bool) {
	switch name {
	case dataASCII, dataSimple:
		return examples.SimpleData, true
	case dataCompact:
		return examples.CompactData, true
	case dataRowspan:
		return examples.RowspanData, true
	case dataColspan:
		return examples.ColspanData, true
	case dataFooter, dataTransformer:
		return examples.FooterData, true
	case dataComplex:
		return examples.ComplexData, true
	case dataStackedHeader:
		return examples.StackedHeaderData, true
	case dataCommaIncluded:
		return examples.CommaIncludedData, true
	default:
		return examples.Data{}, false
	}
}

func newError(name string, target string) error {
	return fmt.Errorf("examples: data %q is not available for target %q", name, target)
}
