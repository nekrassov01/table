---
name: prepare-pr
description: Inspect the current Git branch against its intended base, propose a concise pull request title, fill the repository pull request template with evidence-backed changes and verification, and save the draft under docs/pr without committing or opening the pull request. Use when the user asks to draft, prepare, write, or update a pull request title or description for the current repository.
---

# Prepare PR

Create a reviewer-ready pull request draft from the complete branch diff. Treat the repository template and the implementation as the sources of truth, and never infer the pull request from the latest commit alone.

## Workflow

1. Find the repository root and read its `AGENTS.md` and applicable documentation before describing the change.
2. Resolve the base branch.
   - Honor a base named by the user.
   - Otherwise use the remote default branch when it is available, then fall back to `main` or `master`.
   - Verify the selected base and report it.
3. Inspect the complete change.
   - Run `git status --short` and distinguish committed changes from staged, unstaged, and untracked work.
   - Review `git log --reverse <base>..HEAD`, `git diff --stat <base>...HEAD`, and `git diff --name-status <base>...HEAD`.
   - Read the substantive source, tests, documentation, and generated assets needed to explain behavior and intent accurately.
   - Audit PR granularity without redefining it. If the history or diff contains concerns that can be reviewed and reverted independently, report the possible split before drafting and request direction. Treat implementation, tests, and documentation for one behavior as one concern.
   - Include every in-scope branch change. Do not make the newest commit stand in for the branch.
4. Read `.github/PULL_REQUEST_TEMPLATE.md` from the current repository. Preserve its section order and checklist structure; do not maintain a copied template in this skill.
5. Draft the pull request in English.
   - Put the proposed title in an H1 above the template body when the template has no title field.
   - Make the title concise and outcome-oriented. Follow the repository's title convention when one is evident.
   - Use `Summary` for what changed and why, `Changes` for reviewer-oriented implementation groups, `Verification` for commands actually run and their results, and `Notes` only for material constraints, deferred work, or review context.
   - Normalize verification commands to the repository's documented entry points. Omit environment-manager wrappers, cache variables, absolute temporary paths, and other author-specific setup unless they are required to reproduce the result.
   - Check only the applicable `Type` entries.
   - Do not invent validation. Run the relevant documented commands when authorized; otherwise state precisely what was not run.
6. Save the draft under `docs/pr/`.
   - Use `docs/pr/YYYY-MM-DD-<short-slug>.md` with the repository's local date and a concise slug derived from the PR title or scope. Never use a date-only filename.
   - If that path already contains a different draft, preserve it and choose a more specific slug.
   - Do not alter the repository pull request template.
7. Verify every claim against the diff and test output. Check the completed Markdown for placeholders, stale comments, unchecked applicable types, and formatting errors.
8. Report the proposed title, output path, base branch, compared commit range, validation results, and any uncommitted work excluded from the draft.

## Guardrails

- Do not stage, commit, push, open, or update a remote pull request.
- Do not overwrite an unrelated draft.
- Do not hide unrelated branch commits to produce a narrower story.
- Do not list files mechanically when a responsibility-oriented explanation is clearer.
- Do not claim compatibility, performance, coverage, or successful validation without direct evidence.
- Do not duplicate repository documentation in the pull request; link or summarize only what reviewers need.
