package examples

import (
	"fmt"
	"io"
	"slices"

	"github.com/nekrassov01/table"
	"github.com/nekrassov01/table/backlog"
	"github.com/nekrassov01/table/csv"
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

// Run renders the examples selected by args to w.
func Run(w io.Writer, args ...string) error {
	if len(args) < 1 || len(args) > 3 {
		return fmt.Errorf("examples: expected <target> [mode] [data]")
	}
	return newRunner(w, args...).run()
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
		opts = TextOptionASCII
	case dataSimple:
		opts = TextOptionSimple
	case dataCompact:
		opts = TextOptionCompact
	case dataRowspan:
		opts = TextOptionRowspan
	case dataColspan:
		opts = TextOptionColspan
	case dataFooter:
		opts = TextOptionFooter
	case dataTransformer:
		opts = TextOptionTransformer
	case dataComplex:
		opts = TextOptionComplex
	case dataStackedHeader:
		opts = TextOptionStackedHeader
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
		opts = HTMLOptionSimple
	case dataRowspan:
		opts = HTMLOptionRowspan
	case dataColspan:
		opts = HTMLOptionColspan
	case dataFooter:
		opts = HTMLOptionFooter
	case dataTransformer:
		opts = HTMLOptionTransformer
	case dataComplex:
		opts = HTMLOptionComplex
	case dataStackedHeader:
		opts = HTMLOptionStackedHeader
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
		opts = MarkdownOptionSimple
	case dataRowspan:
		opts = MarkdownOptionRowspan
	case dataColspan:
		opts = MarkdownOptionColspan
	case dataTransformer:
		opts = MarkdownOptionTransformer
	case dataComplex:
		opts = MarkdownOptionComplex
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
		opts = BacklogOptionSimple
	case dataRowspan:
		opts = BacklogOptionRowspan
	case dataColspan:
		opts = BacklogOptionColspan
	case dataFooter:
		opts = BacklogOptionFooter
	case dataTransformer:
		opts = BacklogOptionTransformer
	case dataComplex:
		opts = BacklogOptionComplex
	case dataStackedHeader:
		opts = BacklogOptionStackedHeader
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
		opts = CSVOptionSimple
	case dataFooter:
		opts = CSVOptionFooter
	case dataTransformer:
		opts = CSVOptionTransformer
	case dataComplex:
		opts = CSVOptionComplex
	case dataCommaIncluded:
		opts = CSVOptionCommaIncluded
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

func exampleData(name string) (Data, bool) {
	switch name {
	case dataASCII, dataSimple:
		return SimpleData, true
	case dataCompact:
		return CompactData, true
	case dataRowspan:
		return RowspanData, true
	case dataColspan:
		return ColspanData, true
	case dataFooter, dataTransformer:
		return FooterData, true
	case dataComplex:
		return ComplexData, true
	case dataStackedHeader:
		return StackedHeaderData, true
	case dataCommaIncluded:
		return CommaIncludedData, true
	default:
		return Data{}, false
	}
}

func newError(name string, target string) error {
	return fmt.Errorf("examples: data %q is not available for target %q", name, target)
}
