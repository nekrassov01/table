# Performance Requirements

This document defines the performance requirements and measurement process for maintainers.

## Table of contents

- [Performance Requirements](#performance-requirements)
  - [Table of contents](#table-of-contents)
  - [Evaluation criteria](#evaluation-criteria)
  - [Benchmark structure](#benchmark-structure)
  - [Comparing changes](#comparing-changes)
  - [Checking allocation counts](#checking-allocation-counts)
  - [Checking execution time](#checking-execution-time)
  - [Checking allocated bytes](#checking-allocated-bytes)
  - [Accepting performance changes](#accepting-performance-changes)

## Evaluation criteria

Use the following criteria in order. Absolute thresholds are not fixed because results depend on the environment. Measure the unchanged and changed implementations under the same conditions and use the unchanged result as the baseline. Unless stated otherwise, a result is the median of the repeated samples for that benchmark and metric. `B/op` can vary slightly with the initial growth of `sync.Pool` storage and benchmark execution order.

| Priority | Metric    | Acceptance condition                                                                                                                                                                                                       |
| -------- | --------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1        | allocs/op | The changed result must not exceed the baseline.<br/>A `Table` `Reuse` case whose caller-provided functions do not allocate must remain at 0.                                                                              |
| 2        | ns/op     | The changed result must not exceed the baseline.<br/>Run multiple measurements with the same environment and input.<br/>Increase the iteration and sample counts if results do not converge.                               |
| 3        | B/op      | If neither higher-priority metric regresses and at least one improves, obtain approval for an increase of 10% or more.<br/>If both higher-priority metrics are unchanged, the changed result must not exceed the baseline. |

An intentional regression required by a feature must be approved with its reason, affected paths, and considered alternatives. An improvement in a lower-priority metric does not offset a regression in a higher-priority metric.

## Benchmark structure

The `benchmarks` module uses shared sample data and the corresponding format-specific `Option` sets. Benchmarks follow these naming rules.

| Pattern                             | Measurement                                                               |
| ----------------------------------- | ------------------------------------------------------------------------- |
| `Benchmark<Format>Table<Case>Fresh` | Construct a new `Table` and call `Render` in every iteration.             |
| `Benchmark<Format>Table<Case>Reuse` | Repeatedly call `Render` on the same `Table`.                             |
| `Benchmark<Format>Table<Case>Cold`  | Render one `Table` without reusing a pooled workspace.                    |
| `Benchmark<Format>Stream<Case>`     | Call `NewStream`, row-at-a-time `Render`, and `Close` in every iteration. |
| `Benchmark<Format>Stream<Case>Cold` | Render one `Stream` without reusing a pooled workspace.                   |

- `<Format>` is one of `Text`, `HTML`, `Markdown`, `Backlog`, or `CSV`.
- `<Case>` identifies the input, such as `Simple`, `Rowspan`, `Colspan`, `Footer`, `Complex`, or `Transformer`.
- `Table` and `Stream` identify the exported API being measured.
- `Fresh` and `Reuse` indicate whether a `Table` is reconstructed for each iteration or reused between iterations.
- `Cold` runs garbage collection twice per iteration to discard workspace retained by `sync.Pool`.

## Comparing changes

Save before-and-after results using the same machine, Go version, `benchtime`, and `count`. Compare the median of the repeated samples for each benchmark and metric.

```sh
make bench target=text count=10 cpuprofile=before.cpu memprofile=before.mem > before.txt
make bench target=text count=10 cpuprofile=after.cpu memprofile=after.mem > after.txt
benchstat before.txt after.txt
```

The example selects `text`; replace it with the affected format. Use `benchstat` to inspect variance and the overall direction of the change. The baseline-checking script reports the explicit per-benchmark medians used for acceptance.

When publishing measurements in `README.md` or another document, paste the complete console output from a run with at least five samples. Do not remove, reorder, or annotate sample lines inside the console block. Derive any summary from the same output and report the median for each benchmark and metric.

## Checking allocation counts

When allocs/op changes, first verify that both benchmarks use identical input and call order.

Slices and buffers used during rendering live in an `arena`, which retains reusable backing arrays across resets. Common paths avoid temporary slices and buffers, string-to-byte copies, and unnecessary interface conversions. Use compiler escape analysis to identify the source of a regression.

```sh
go test -gcflags='-m=2' ./text
```

Inspect the following areas:

- Values captured by an `Option` closure escaping to the heap.
- Row, cell, or column-measurement slices allocated per render instead of stored in the `arena`.
- Interface conversions or `fmt` calls added to frequently executed paths.
- Copies introduced by conversions between `string` and `[]byte`.
- `Stream` resets that discard reusable backing arrays.

Interpret the results as follows:

- The difference between `Fresh` and `Reuse` separates allocations caused by constructing `Table` from allocations caused by `Render` itself.
- `Complex` and `Transformer` may include allocations from caller-provided closures and must not automatically be treated like a regression in `Simple`.

## Checking execution time

Compare ns/op after allocs/op. Background processes and CPU frequency affect execution time, so use the median from multiple measurements.

Each pipeline stage passes forward results produced while transforming values, escaping text, or measuring display width. Later stages do not repeat value conversion or cell scans merely to derive the same result.

- Keep the machine, Go version, input data, options, `benchtime`, and `count` identical before and after the change.
- Select the median ns/op for each benchmark independently. Do not combine values from different benchmarks.
- If the median regresses, increase `benchtime` and `count` and measure both versions again.
- Use a CPU profile to locate the additional work and connect it to the affected implementation path.

A change in only the mean or minimum is not sufficient evidence of a regression. If the median does not stabilize, discard the results and measure again in a controlled environment.

## Checking allocated bytes

Compare B/op after allocs/op and ns/op. If B/op increases while allocs/op is unchanged, use a memory profile to identify the allocation site and size. Repeat the measurement under identical conditions if the difference depends on pool state or benchmark order.

To reduce buffer growth, each stage passes sizes and maxima obtained during its work to later capacity estimates. The estimate does not need to predict the exact output size; it should avoid buffer growth on the common path.

An underestimated capacity remains safe because `append` may grow the buffer. Estimation must never cause a panic or truncated output. A change that recomputes sizes or maxima solely for capacity estimation must demonstrate both simpler structure and a benchmark benefit.

## Accepting performance changes

Verify all of the following:

- Unit, contract, and golden tests pass with the race detector enabled.
- allocs/op meets the evaluation criteria for every affected benchmark.
- Repeated before-and-after measurements show no regression in the median ns/op.
- Any B/op difference can be explained by arena growth or values captured by closures.
- A local optimization in one format does not unnecessarily break responsibility and vocabulary symmetry with other formats.

Meeting the numeric criteria is not sufficient when a change introduces a stage without a responsibility, duplicates state, or unnecessarily disrupts symmetry between formats.
