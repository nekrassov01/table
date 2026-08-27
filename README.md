<p align="center">
  <img src="./assets/logo.png" alt="table logo" width="120">
</p>
<h1 align="center">TABLE</h1>

<p align="center">A high-performance table rendering library for Go, with streaming APIs for every output format.</p>
<p align="center">
  <a href="https://github.com/nekrassov01/table/actions/workflows/ci.yml"><img src="https://github.com/nekrassov01/table/actions/workflows/ci.yml/badge.svg?branch=main" alt="CI"></a>
  <a href="https://pkg.go.dev/github.com/nekrassov01/table"><img src="https://pkg.go.dev/badge/github.com/nekrassov01/table.svg" alt="Go Reference"></a>
  <img src="https://img.shields.io/github/license/nekrassov01/table" alt="License">
</p>

## Table of contents

- [Table of contents](#table-of-contents)
- [Overview](#overview)
- [Motivation](#motivation)
- [Output formats](#output-formats)
- [Text output gallery](#text-output-gallery)
- [Installation](#installation)
- [Examples](#examples)
  - [Table](#table)
  - [Stream](#stream)
- [Table library comparison](#table-library-comparison)
  - [Formats](#formats)
  - [Features](#features)
- [Performance](#performance)
- [Documentation](#documentation)
- [Author](#author)
- [License](#license)

## Overview

`nekrassov01/table` renders Go data as terminal tables, markup tables, or CSV records. Each output package provides `Table` for complete data sets and `Stream` for row-at-a-time output while retaining the selected format's own structure and escaping rules.

Browse the [examples catalog](./docs/EXAMPLES.md) for a broad range of ways to use the library, with input, options, and exact output for each example.

- `table` uses functional options for clear, reusable configuration.
- In the bundled comparisons, `table` runs 5.9 to 7.0 times as fast as the next-fastest alternative; see [Performance](#performance).
- `table` reuses internal buffers to minimize steady-state allocations.
- `TableOf` and `StreamOf` adapt typed slices and error-returning iterators.
- `text` measures Unicode by terminal display width, including ambiguous character widths in CJK locales.
- Format-specific options add headers, calculated footers, indexes, placeholders, transformations, alignment, decoration, and cell spans.

## Motivation

The project was created for four reasons:

- To provide a table renderer that is fast, efficient, and easy to use.
- To provide row-at-a-time output across every format, which no comparable library offered at the time.
- To manage table-oriented output formats for applications such as CLIs in one module.
- To support the less common Backlog table notation, which only the author's earlier [`mintab`](https://github.com/nekrassov01/mintab) project covered at the time.

## Output formats

Choose an output package for the destination. The root [`table`](.) package provides shared interfaces, typed row adapters, and errors; it does not select an output format.

| Package                  | Output                                                                                                                                                                                                                 | Use it for                     |
| ------------------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------ |
| [`text`](./text)         | Unicode or ASCII bordered tables                                                                                                                                                                                       | CLIs, terminals, and logs      |
| [`html`](./html)         | Semantic HTML tables                                                                                                                                                                                                   | Web pages and reports          |
| [`markdown`](./markdown) | GFM tables                                                                                                                                                                                                             | READMEs and GitHub             |
| [`backlog`](./backlog)   | [Backlog table notation](https://help-center.backlog.com/%E3%83%86%E3%82%AD%E3%82%B9%E3%83%88%E6%95%B4%E5%BD%A2%E3%83%AB%E3%83%BC%E3%83%AB%EF%BC%88Backlog%E8%A8%98%E6%B3%95%EF%BC%89-6a1d4d7f3abb3ada78c5658b#index7) | Backlog issues and Wiki pages  |
| [`csv`](./csv)           | Configurable delimiter-separated records                                                                                                                                                                               | TSV, CSV, and data interchange |

## Text output gallery

The examples below show the `text` package in several configurations.

ASCII

![ASCII table](./assets/examples/text/ascii.png)

Simple

![Simple table](./assets/examples/text/simple.png)

Compact layout with horizontal lines omitted and a colored rounded border style

![Compact table](./assets/examples/text/compact.png)

Row spans with a colored heavy border style

![Table with row spans](./assets/examples/text/rowspan.png)

Column spans with a colored light border style

![Table with column spans](./assets/examples/text/colspan.png)

Stacked header with a colored light border style

![Table with a stacked header](./assets/examples/text/stacked-header.png)

Calculated footer with CJK text and a colored light border style

![Table with a calculated footer](./assets/examples/text/footer.png)

Calculated footer with CJK text, value transformations, and a colored light border style

![Table with transformed values](./assets/examples/text/transformer.png)

Complex values with a colored double border style

![Table containing complex values](./assets/examples/text/complex.png)

## Installation

Install with:

```sh
go get github.com/nekrassov01/table
```

## Examples

The runnable examples use `target`, `mode`, and `data` to select the output package, `Table` or `Stream`, and a data set. Omitting `target` runs every output package, omitting `mode` runs both APIs, and omitting `data` runs every data set available for each selected package. When selecting one data set, also specify `mode`. This command renders every text `Table` example:

```sh
make example target=text mode=table
```

### Table

Use `TableOf` to keep application data typed until the row boundary. This example renders deployment status to a terminal.

```go
package main

import (
    "log"
    "os"

    "github.com/nekrassov01/table"
    "github.com/nekrassov01/table/text"
)

type Deployment struct {
    Service string
    Desired int
    Ready   int
    Status  string
}

func deploymentRow(deployment Deployment) []any {
    return []any{
        deployment.Service,
        deployment.Desired,
        deployment.Ready,
        deployment.Status,
    }
}

func main() {
    deployments := []Deployment{
        {Service: "payments", Desired: 4, Ready: 4, Status: "healthy"},
        {Service: "search", Desired: 3, Ready: 2, Status: "degraded"},
        {Service: "worker", Desired: 8, Ready: 8, Status: "healthy"},
    }

    output := text.NewTable(os.Stdout,
        text.WithHeader([]string{"SERVICE", "DESIRED", "READY", "STATUS"}),
        text.WithAlign(text.ScopeBody, text.Columns(1, 2), text.AlignRight),
        text.WithCompact(),
    )
    if err := output.Render(table.TableOf(deployments, deploymentRow)); err != nil {
        log.Fatal(err)
    }
}
```

This program produces the following table:

```text
┌──────────┬─────────┬───────┬──────────┐
│ SERVICE  │ DESIRED │ READY │  STATUS  │
╞══════════╪═════════╪═══════╪══════════╡
│ payments │       4 │     4 │ healthy  │
│ search   │       3 │     2 │ degraded │
│ worker   │       8 │     8 │ healthy  │
└──────────┴─────────┴───────┴──────────┘
```

### Stream

Use `StreamOf` to adapt an `iter.Seq2[T, error]` to streaming table rows. Each successful value becomes one output row, and the first source error is forwarded. Additionally, call `Stream.Close` for every stream so it can write deferred output, such as a footer or closing border, and report late errors.

```go
package report

import (
    "errors"
    "io"
    "iter"
    "time"

    "github.com/nekrassov01/table"
    "github.com/nekrassov01/table/text"
)

type AuditEvent struct {
    Time     time.Time
    Actor    string
    Action   string
    Resource string
}

func WriteAuditEvents(w io.Writer, events iter.Seq2[AuditEvent, error]) (err error) {
    output := text.NewStream(w,
        text.WithHeader([]string{"TIME", "ACTOR", "ACTION", "RESOURCE"}),
        text.WithCompact(),
    )
    defer func() {
        err = errors.Join(err, output.Close())
    }()

    rows := table.StreamOf(events, func(event AuditEvent) []any {
        return []any{
            event.Time.Format(time.RFC3339),
            event.Actor,
            event.Action,
            event.Resource,
        }
    })
    for row, sourceErr := range rows {
        if sourceErr != nil {
            return sourceErr
        }
        if renderErr := output.Render(row); renderErr != nil {
            return renderErr
        }
    }
    return nil
}
```

Given an iterator of audit events, the function produces output like this:

```text
┌───────────────────────────┬────────────┬───────────┬──────────────┐
│           TIME            │   ACTOR    │  ACTION   │   RESOURCE   │
╞═══════════════════════════╪════════════╪═══════════╪══════════════╡
│ 2026-08-21T09:00:00+09:00 │ deploy-bot │ reconcile │ payments-api │
│ 2026-08-21T09:03:00+09:00 │ alice      │ scale     │ worker       │
│ 2026-08-21T09:08:00+09:00 │ bob        │ rollback  │ search-api   │
└───────────────────────────┴────────────┴───────────┴──────────────┘
```

## Table library comparison

The following tables compare the public APIs in the versions pinned by the [benchmark module](./benchmarks/go.mod).

### Formats

This table records the output implementations documented by each library. `✓` means the library provides a dedicated output mode for the format, and `-` means it does not. `table` targets the GFM table extension; the other Markdown entries indicate generic Markdown table output.

| Output format    | `table` | [`go-pretty` v6.8.3](https://github.com/jedib0t/go-pretty/tree/v6.8.3) | [`tablewriter` v1.1.4](https://github.com/olekukonko/tablewriter/tree/v1.1.4) | [`simpletable` v1.0.0](https://github.com/alexeyco/simpletable/tree/v1.0.0) |
| ---------------- | ------- | ---------------------------------------------------------------------- | ----------------------------------------------------------------------------- | --------------------------------------------------------------------------- |
| Text             | ✓       | ✓                                                                      | ✓                                                                             | ✓                                                                           |
| HTML             | ✓       | ✓                                                                      | ✓                                                                             | -                                                                           |
| Markdown         | ✓       | ✓                                                                      | ✓                                                                             | ✓                                                                           |
| Backlog notation | ✓       | -                                                                      | -                                                                             | -                                                                           |
| CSV or TSV       | ✓       | ✓                                                                      | -                                                                             | -                                                                           |
| SVG              | -       | -                                                                      | ✓                                                                             | -                                                                           |

### Features

This table records whether each library exposes a direct public API for a capability in at least one output implementation. It does not imply that every format can express the capability.

| Feature                         | `table` | [`go-pretty` v6.8.3](https://github.com/jedib0t/go-pretty/tree/v6.8.3) | [`tablewriter` v1.1.4](https://github.com/olekukonko/tablewriter/tree/v1.1.4) | [`simpletable` v1.0.0](https://github.com/alexeyco/simpletable/tree/v1.0.0) |
| ------------------------------- | ------- | ---------------------------------------------------------------------- | ----------------------------------------------------------------------------- | --------------------------------------------------------------------------- |
| Streaming API                   | ✓       | -                                                                      | ✓                                                                             | -                                                                           |
| Header                          | ✓       | ✓                                                                      | ✓                                                                             | ✓                                                                           |
| Footer                          | ✓       | ✓                                                                      | ✓                                                                             | ✓                                                                           |
| Placeholder                     | ✓       | ✓ (HTML)                                                               | -                                                                             | -                                                                           |
| Index column                    | ✓       | ✓                                                                      | -                                                                             | -                                                                           |
| Vertical merge                  | ✓       | ✓                                                                      | ✓                                                                             | -                                                                           |
| Horizontal merge                | ✓       | ✓                                                                      | ✓                                                                             | ✓                                                                           |
| Caller-defined row adapter      | ✓       | -                                                                      | ✓                                                                             | -                                                                           |
| Reflection-based struct input   | -       | -                                                                      | ✓                                                                             | -                                                                           |
| CSV input                       | -       | -                                                                      | ✓                                                                             | -                                                                           |
| Per-column transformation       | ✓       | ✓                                                                      | ✓                                                                             | -                                                                           |
| Built-in sorting and filtering  | -       | ✓                                                                      | -                                                                             | -                                                                           |
| Pagination                      | -       | ✓                                                                      | -                                                                             | -                                                                           |
| Column hiding                   | -       | ✓                                                                      | ✓                                                                             | -                                                                           |
| Width, wrapping, and truncation | ✓       | ✓                                                                      | ✓                                                                             | -                                                                           |
| Automatic terminal fit          | ✓       | -                                                                      | -                                                                             | -                                                                           |
| Title or caption                | ✓       | ✓                                                                      | ✓                                                                             | -                                                                           |
| Pluggable output implementation | -       | -                                                                      | ✓                                                                             | -                                                                           |

For `table`, the feature matrix has the following qualifications:

- Footer callbacks derive values such as totals and averages from captured state.
- Column hiding is intentionally left to input adaptation, so `TableOf` and `StreamOf` can omit fields before rows reach the output package.
- Merge behavior depends on the selected output format and is documented in the [Public API guide](./docs/API.md).

The `go-pretty` placeholder entry refers to its HTML `EmptyColumn` setting.

## Performance

> [!NOTE]
> `table` is designed to minimize allocations and reuses internal buffers via `sync.Pool`. In steady-state benchmarks, common workloads reach one allocation per render; cold runs require additional allocations to initialize pooled state.
>
> Run `make bench target=comparison benchtime=1x count=1` to reduce steady-state amortization and expose one-iteration setup costs. For explicit pool-drained measurements of `table`, run `make bench target=cold`.

Run the comparison on your machine with `make bench target=comparison`. The following output records all five samples on an Apple M2 with Go 1.27.0:

```console
$ make bench target=comparison
go test -benchmem -count 5 -benchtime 10000x -cpuprofile cpu.prof -memprofile mem.prof . -bench '^BenchmarkComparison'
goos: darwin
goarch: arm64
pkg: benchmarks
cpu: Apple M2
BenchmarkComparisonTableSimple-8           10000        2436 ns/op       225 B/op         1 allocs/op
BenchmarkComparisonTableSimple-8           10000        1847 ns/op       224 B/op         1 allocs/op
BenchmarkComparisonTableSimple-8           10000        1739 ns/op       224 B/op         1 allocs/op
BenchmarkComparisonTableSimple-8           10000        1804 ns/op       224 B/op         1 allocs/op
BenchmarkComparisonTableSimple-8           10000        1773 ns/op       224 B/op         1 allocs/op
BenchmarkComparisonGoPrettySimple-8        10000       10580 ns/op      8155 B/op       110 allocs/op
BenchmarkComparisonGoPrettySimple-8        10000       10572 ns/op      8153 B/op       110 allocs/op
BenchmarkComparisonGoPrettySimple-8        10000       11115 ns/op      8153 B/op       110 allocs/op
BenchmarkComparisonGoPrettySimple-8        10000       11826 ns/op      8153 B/op       110 allocs/op
BenchmarkComparisonGoPrettySimple-8        10000       10688 ns/op      8153 B/op       110 allocs/op
BenchmarkComparisonTableWriterSimple-8     10000       81974 ns/op    486948 B/op       973 allocs/op
BenchmarkComparisonTableWriterSimple-8     10000       79278 ns/op    486948 B/op       973 allocs/op
BenchmarkComparisonTableWriterSimple-8     10000       81598 ns/op    486948 B/op       973 allocs/op
BenchmarkComparisonTableWriterSimple-8     10000       87737 ns/op    486949 B/op       973 allocs/op
BenchmarkComparisonTableWriterSimple-8     10000       80748 ns/op    486948 B/op       973 allocs/op
BenchmarkComparisonSimpleTableSimple-8     10000       23490 ns/op     13109 B/op       425 allocs/op
BenchmarkComparisonSimpleTableSimple-8     10000       23955 ns/op     13098 B/op       425 allocs/op
BenchmarkComparisonSimpleTableSimple-8     10000       22876 ns/op     13121 B/op       425 allocs/op
BenchmarkComparisonSimpleTableSimple-8     10000       22583 ns/op     13116 B/op       425 allocs/op
BenchmarkComparisonSimpleTableSimple-8     10000       22650 ns/op     13098 B/op       425 allocs/op
BenchmarkComparisonTableComplex-8          10000        9192 ns/op      1150 B/op        35 allocs/op
BenchmarkComparisonTableComplex-8          10000        9158 ns/op      1150 B/op        35 allocs/op
BenchmarkComparisonTableComplex-8          10000        9176 ns/op      1149 B/op        35 allocs/op
BenchmarkComparisonTableComplex-8          10000        9128 ns/op      1149 B/op        35 allocs/op
BenchmarkComparisonTableComplex-8          10000        9106 ns/op      1150 B/op        35 allocs/op
BenchmarkComparisonGoPrettyComplex-8       10000       63998 ns/op     49260 B/op       317 allocs/op
BenchmarkComparisonGoPrettyComplex-8       10000       63791 ns/op     49260 B/op       317 allocs/op
BenchmarkComparisonGoPrettyComplex-8       10000       61935 ns/op     49263 B/op       317 allocs/op
BenchmarkComparisonGoPrettyComplex-8       10000       64352 ns/op     49261 B/op       317 allocs/op
BenchmarkComparisonGoPrettyComplex-8       10000       62202 ns/op     49261 B/op       317 allocs/op
BenchmarkComparisonTableWriterComplex-8    10000      276792 ns/op    720005 B/op      4749 allocs/op
BenchmarkComparisonTableWriterComplex-8    10000      280501 ns/op    719999 B/op      4749 allocs/op
BenchmarkComparisonTableWriterComplex-8    10000      277848 ns/op    720000 B/op      4749 allocs/op
BenchmarkComparisonTableWriterComplex-8    10000      278027 ns/op    719994 B/op      4749 allocs/op
BenchmarkComparisonTableWriterComplex-8    10000      279195 ns/op    719997 B/op      4749 allocs/op
PASS
ok      benchmarks      24.222s
```

The table summarizes those five samples. Each cell shows `allocs/op · ns/op`; both values are medians.

| Scenario       | `table`        | `go-pretty`  | `tablewriter`   | `simpletable` |
| -------------- | -------------- | ------------ | --------------- | ------------- |
| Simple         | **1 · 1,804**  | 110 · 10,688 | 973 · 81,598    | 425 · 22,876  |
| Complex values | **35 · 9,158** | 317 · 63,791 | 4,749 · 278,027 | -             |

`-` indicates that a library cannot express the scenario with the benchmark input.

The comparison benchmark uses the shared Simple and Complex data sets. Static data is converted to each library's required row type before timing begins. Each timed iteration constructs a table, processes the rows, and writes the result to a reused buffer. Complex compares native value handling rather than equivalent rendered bytes.

The benchmark preserves each library's configuration model: `table` uses only functional options, `go-pretty` accumulates settings through setters, `tablewriter` combines constructor options with methods, and `simpletable` receives prebuilt cells through exposed table sections. These native construction paths remain inside each timed iteration. Only the settings needed to align table structure and preserve header text are applied; border characters and value formatting retain each library's defaults.

## Documentation

Use these references to select the API, understand its design, and work on the module itself.

| Resource                                                        | Contents                                                        |
| --------------------------------------------------------------- | --------------------------------------------------------------- |
| [Go Reference](https://pkg.go.dev/github.com/nekrassov01/table) | Exact declarations and symbol documentation                     |
| [Public API guide](./docs/API.md)                               | Options, output behavior, defaults, and format capabilities     |
| [Architecture](./docs/ARCHITECTURE.md)                          | Package structure, data flow, and state ownership               |
| [Design specification](./docs/DESIGN.md)                        | Design decisions, invariants, tradeoffs, and non-goals          |
| [Development guide](./docs/DEVELOPMENT.md)                      | Test, benchmark, coverage, and analysis commands                |
| [Performance baseline](./docs/BASELINE.md)                      | Benchmark procedure and performance acceptance criteria         |

## Author

[nekrassov01](https://github.com/nekrassov01)

## License

[MIT](https://github.com/nekrassov01/table/blob/main/LICENSE)
