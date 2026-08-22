# Repository Instructions

This repository provides a table rendering library for Go. This file defines only the reading order and repository-specific working rules; it does not repeat the documentation itself.

## References

Select the documents that correspond to the task and read them before changing the implementation. Read every applicable document when a change spans multiple areas.

| Area                                                                             | Document                                       |
| -------------------------------------------------------------------------------- | ---------------------------------------------- |
| Public API, format differences, output contracts, and specification conformance  | [`docs/API.md`](docs/API.md)                   |
| Package structure, pipelines, state ownership, and lifetimes                     | [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) |
| Design principles, responsibility boundaries, constraints, and their rationale   | [`docs/DESIGN.md`](docs/DESIGN.md)             |
| Test, coverage, benchmark, and static analysis commands                          | [`docs/DEVELOPMENT.md`](docs/DEVELOPMENT.md)   |
| Performance metrics, before-and-after measurements, and acceptance criteria      | [`docs/BASELINE.md`](docs/BASELINE.md)         |
| User-facing overview, examples, comparisons, and links to detailed documentation | [`README.md`](README.md)                       |

Read `ARCHITECTURE.md` and `DESIGN.md` before changing internal structure, and read `API.md` before changing public behavior. For performance-sensitive work, read the evaluation order and measurement conditions in `BASELINE.md` before implementation. Use the entry points documented in `DEVELOPMENT.md` rather than recalling commands from memory.

After code changes and validation are complete, use the `sync-docs` skill before committing. Do not commit until the skill has updated the affected files in `docs/` and `README.md` and the user has approved the complete change.

## Code

- Use the established vocabulary for an existing responsibility. Do not rename an existing concept merely to introduce a new term. If a name does not fit the vocabulary, reconsider the responsibility boundary.
- Name receivers `o`. Do not use a `get*` prefix; use `resolve*` for derived values.
- Row means a logical row, Line means a physical line, segment means a cell fragment, span/rowspan/colspan mean cell spanning, element means an HTML element, and box means output geometry. Do not introduce merge, group, unit, or renderer into identifiers.
- Within a file, order declarations as constants, types with their constructors and receiver methods, and other package functions. Put each constructor immediately below the type it constructs, followed by that type's receiver methods. Order other package functions by call order, and declare local variables with individual `:=` statements.
- Preserve changes that the user is editing. Do not format, move, or restore out-of-scope differences.

## Tests

- Except in `contract_test.go` and `golden_test.go`, follow the table-driven template in `text/arena_test.go`. Use `fields` for receiver state, `args` for arguments, and `want` for observed results.
- Keep test functions contiguous. Put file-local helpers after all test functions, and put shared pure assertions and mocks in `internal/testutil`.
- Do not add assertion branches for individual cases. Prefer the existing assertions in `internal/testutil`.
- Do not generate golden files mechanically. Show the complete new or changed output to the user and obtain approval before writing it by hand. Apply the same process when regenerating an existing golden file.
- Name golden cases `common_*` when Table and Stream share a contract, and `table_*` or `stream_*` when they differ. Verify fixes with fault injection to demonstrate that the test detects the failure.

## Commits

- Use English Conventional Commit messages. Separate structural changes from behavior changes, and keep every commit independently verifiable.
- Do not commit until the implementation, tests, related documentation, and validation results have been presented and approved. Approval applies only to the diff presented at that point.
