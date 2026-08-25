---
name: check-goldens
description: Verify the golden corpus of each output package, measure which option pairs no case exercises, and extend the corpus only with the user's approval. Use after adding, changing, renaming, or reordering golden tests, after a fix that a golden should pin, and before any commit that touches testdata.
---

# Check Goldens

Keep the golden corpus honest. The corpus is the value of this library, and `go test` cannot defend it: a deleted test passes, an orphaned golden passes, and a case that proves nothing passes. Treat the rendered bytes as evidence and stop before committing.

## Workflow

1. Determine the output packages in scope from the changed files, preserving the order `text`, `html`, `markdown`, `backlog`, and `csv`. Default to all five when the scope cannot be derived. Report the scope before checking it.
2. Run `go run ./.agents/skills/check-goldens/scripts` with the packages in scope as arguments. Treat only `testdata/*.txt` as golden files; supporting assets such as `html/testdata/style.css` are outside the corpus.
3. Reconcile the audit with `golden_test.go`: every `TestGolden_*` function has exactly one `AssertGolden` call, every asserted name has a file, every golden file is asserted, `common_*` has two references, and `table_*` and `stream_*` have one reference.
4. Review the placement of changed test declarations against the surrounding file. Keep Table and Stream adjacent in the shared-behavior section, with Table first. Keep Stream-only and Table-only cases in their existing sections and order each section by case name. Determine the section from the inputs and contract; a `table_*` or `stream_*` filename alone does not prove that a case is API-specific.
5. Review byte-identical files reported by the audit as candidates, not failures. Inspect their inputs and options before deciding whether they intentionally pin equivalent behavior, such as an explicit default, or duplicate a case without adding a contract.
6. Use the audit's option inventory and uncovered pairs as the mechanical input. Investigate uncovered pairs in this order: geometry against geometry, such as width, padding, truncation, spans, the index column, and terminal fitting; then geometry against markup; then the remaining pairs. Do not equate an uncovered pair with a defect.
7. Propose cases only for uncovered pairs whose options can interact and whose result would pin a meaningful contract. Render each proposed case and show its complete output. Do not write the golden file.
8. When a case pins a fix, reproduce the candidate test and the inverse code change in a disposable copy or worktree and show that the test fails there. Never inject a fault into the user's active worktree.
9. Report per package: the reconciled counts, every failed invariant, byte-identical candidates and their explanations, uncovered pairs worth testing, proposed cases, and what the user accepted or declined.

## Guardrails

- Do not write or regenerate a golden file. Show the complete output and let the user place it by hand.
- Do not judge a golden from a truncated view, a diff summary, or a description of the output.
- Do not adjust test data so that Table and Stream agree. A `common_*` case that stops agreeing becomes a `table_*` and `stream_*` pair.
- Do not classify byte-identical files without inspecting the inputs and options that produce them.
- Do not report a fix as pinned without isolated fault injection that shows the golden failing.
- Do not add a test file that performs these checks. This skill is the check.
- Do not infer declaration sections from golden-name prefixes. When reordering is necessary, operate on the syntax tree, preserve Table and Stream pairs, and abort when the declaration count changes.
- Do not propose a case for every uncovered pair. Two options that cannot interact cost a golden without pinning anything.
- Do not leave disposable copies, worktrees, or temporary output behind.
- Do not stage or commit.
