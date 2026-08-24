---
name: check-baseline
description: Measure a change against the immediately preceding implementation on the same machine and apply the evaluation criteria in docs/BASELINE.md. Use after touching a rendering path, pooled storage, capacity estimation, or anything else that can move allocs/op, ns/op, or B/op, and before any commit that claims a performance result.
---

# Check Baseline

Compare against a measurement, never against memory. `docs/BASELINE.md` owns the criteria, benchmark structure, and investigation process. This skill executes that process, explains every difference it reports, and stops before committing.

## Workflow

1. Identify the changed paths and the benchmarks that cover them, and report the scope. When no benchmark reaches a changed path, say so and propose one rather than measuring something adjacent.
2. Read the evaluation criteria and measurement process in `docs/BASELINE.md` before measuring. That document decides what counts as a regression; do not restate or reinterpret its criteria here.
3. Resolve the revision immediately preceding the measured change. Use `HEAD` for staged or unstaged changes. For an already committed change or range, identify its parent or explicit base from the requested scope and verify it with `git diff`; never substitute the branch merge base without evidence that the entire branch is the measured change.
4. Confirm that `benchstat` is available. If it is missing, stop and report `go install golang.org/x/perf/cmd/benchstat@latest`; do not install it without approval.
5. Run `go run ./.agents/skills/check-baseline/scripts` with the resolved revision, benchmark target, `benchtime`, and `count`. Use `--bench` instead of `--target` for a focused rerun. The script records the environment, creates a detached worktree, measures both revisions, runs `benchstat`, extracts the median of each metric, and removes its temporary worktree, output, and profiles. Add `--keep` only when profiles are needed to investigate a confirmed regression, then remove the reported artifact directory after the investigation.
6. Inspect the `benchstat` comparison for variance and direction. Use the script's per-benchmark medians as the acceptance values.
7. Apply the criteria in priority order: allocs/op first, then ns/op, then B/op. Confirm that every affected `Table` `Reuse` benchmark is still at 0 allocs/op.
8. Explain each difference before reporting it. Attribute an allocation difference with escape-analysis evidence. Attribute a byte difference to arena growth, pool state, benchmark order, or a caller-provided closure. Increase `benchtime` and `count` for both revisions when results do not converge; never select only a favorable run.
9. Run the tests with the race detector and confirm the structural acceptance conditions in `docs/BASELINE.md` that measurements cannot show.
10. Report a per-benchmark table of baseline against changed for all three metrics, the verdict for each criterion, the resolved baseline revision, the recorded environment, and anything still unexplained. Do not add machine-dependent measurements to `docs/BASELINE.md`.
11. Confirm that the script removed its temporary worktree and files, including after a failed measurement. Leave the active worktree and its existing changes untouched.

## Guardrails

- Do not infer the baseline revision from the branch name or merge base. Verify that it contains the implementation immediately preceding the measured change.
- Do not compare against remembered numbers, a previous session's output, or a run whose conditions were not recorded.
- Do not select a minimum or another favorable sample as evidence for ns/op.
- Do not offset a regression in a higher-priority metric with an improvement in a lower-priority one.
- Do not report a B/op difference as acceptable while its cause is unknown.
- Do not accept a numeric win that adds a stage without a responsibility, duplicates state, or breaks vocabulary and responsibility symmetry between the media.
- Do not edit `docs/BASELINE.md` merely to record a machine-dependent result.
- Do not leave worktrees, benchmark output, profiles, or temporary source changes behind.
- Do not stage or commit.
