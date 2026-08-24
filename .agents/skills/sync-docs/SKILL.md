---
name: sync-docs
description: Review completed code changes, update only the affected README and docs, verify that documentation matches the implementation, and request approval before committing. Use after changing code, tests, public APIs, behavior, architecture, development commands, or performance-sensitive paths, and before any commit that may require documentation changes.
---

# Sync Docs

Synchronize documentation with the implementation after code changes. Treat the code and tests as evidence, avoid repeating the same explanation across documents, and stop before committing.

## Workflow

1. Inspect staged, unstaged, and relevant untracked changes against `HEAD`.
2. Trace changed identifiers and behavior through callers, tests, examples, package comments, and related media packages. Do not infer the contract from filenames or the requested change alone.
3. Read `AGENTS.md` and use its reference table to select the canonical documents for the change. Read the document itself rather than relying on a previous summary.
4. Decide whether the code change affects documentation. Internal changes with no observable, structural, procedural, or performance-contract impact may require no documentation edit; state that conclusion explicitly.
5. Update every affected document, but keep each fact in its canonical document. Replace stale text instead of appending a second explanation. Keep `README.md` concise and link to detailed documents when detail is necessary.
6. Check the edited prose against the implementation and tests. Preserve the repository vocabulary and the media order `text`, `html`, `markdown`, `backlog`, `csv`.
7. Run validation appropriate to the edits, including `git diff --check` and the relevant commands in `docs/DEVELOPMENT.md`. If behavior or performance changed, verify the contracts in `docs/API.md` or the criteria and publication rules in `docs/BASELINE.md`.
8. Report the code-to-document mapping, edited files, omitted documents and reasons, and validation results.
9. Do not stage or commit. Ask the user to approve the complete code and documentation diff before committing, even if an earlier request authorized implementation.

## Guardrails

- Do not update documentation merely to mention an internal refactor.
- Do not copy architecture or design detail into `README.md`.
- Do not describe planned behavior as implemented behavior.
- Do not preserve historical explanations unless they are necessary to understand a current design constraint.
- Do not create a new document when an existing document owns the topic.
- Do not broaden the code change while synchronizing documentation.
