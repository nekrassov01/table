# Development

This document is the command reference for maintainers.

## Table of contents

- [Development](#development)
  - [Table of contents](#table-of-contents)
  - [Tests](#tests)
  - [Coverage](#coverage)
  - [Benchmarks](#benchmarks)
  - [Validation](#validation)

## Tests

Run tests from the repository root. Every test command enables the race detector and coverage collection and updates `coverage.out`.

- Unit tests follow the table-driven test template and exercise functions and methods directly.
- Contract tests verify the `Table` and `Stream` lifecycles and error handling.
- Golden tests verify the exact output bytes.

When behavior changes, update the relevant tests, user documentation, and package comments.

| Purpose                 | Command                                                      |
| ----------------------- | ------------------------------------------------------------ |
| Run all tests           | `make test target=all`                                       |
| Test the root package   | `make test target=table`                                     |
| Test internal packages  | `make test target=internal`                                  |
| Test `text`             | `make test target=text`                                      |
| Test `html`             | `make test target=html`                                      |
| Test `markdown`         | `make test target=markdown`                                  |
| Test `backlog`          | `make test target=backlog`                                   |
| Test `csv`              | `make test target=csv`                                       |
| Run only contract tests | `make test target=contract`                                  |
| Run only golden tests   | `make test target=golden`                                    |
| Run one test            | `make test target=markdown run='^TestGolden_TableAutolink$'` |

Running `make test` without `target` is equivalent to `target=all`. The `run` argument accepts the regular expression passed to `go test -run`.

## Coverage

Generate an HTML report from the `coverage.out` file produced by `make test`:

```sh
make cover
```

The report is written to `coverage.html`.

## Benchmarks

The benchmarks live in a separate Go module. The root `Makefile` changes into the `benchmarks` directory before running them.

| Purpose                           | Command                        |
| --------------------------------- | ------------------------------ |
| Run all benchmarks                | `make bench target=all`        |
| Benchmark `text`                  | `make bench target=text`       |
| Benchmark `html`                  | `make bench target=html`       |
| Benchmark `markdown`              | `make bench target=markdown`   |
| Benchmark `backlog`               | `make bench target=backlog`    |
| Benchmark `csv`                   | `make bench target=csv`        |
| Compare external packages         | `make bench target=comparison` |
| Measure without pooled workspaces | `make bench target=cold`       |

Running `make bench` without `target` is equivalent to `target=all`. Regular benchmarks use `benchtime=10000x` and `count=5`, and write CPU and memory profiles to the `benchmarks` directory. Both values can be overridden:

```sh
make bench target=html benchtime=20000x count=10
```

Each `Cold` benchmark runs garbage collection twice per iteration to discard the `sync.Pool` workspace, and uses `benchtime=100x` and `count=1`.

## Validation

Static analysis and vulnerability checks can be run separately or as part of the complete validation sequence.

| Purpose                                                                    | Command      |
| -------------------------------------------------------------------------- | ------------ |
| Run static analysis                                                        | `make lint`  |
| Check for known vulnerabilities                                            | `make vuln`  |
| Run tests, coverage, benchmarks, static analysis, and vulnerability checks | `make check` |

`make lint` and `make vuln` install their required tools when they are not already available.
