
## Documentation

When creating documentation files:

- Place them in the `docs/` directory or a subdirectory. Typical subdirectories include `docs/design/` and `docs/research/`.
- Prepend a UTC timestamp to the `.md` filename in the format `YYYYMMDD_HHMMSS_UTC_` so files sort chronologically (e.g., `docs/design/20260329_150800_UTC_architecture_overview.md`).

---

## Issue Identification & Labeling

allbctl is a developer-facing CLI. Its issue tracker is **GitHub Issues** on this repo, not TODO comments or scratch docs. When asked to audit the codebase, review a subsystem, or "find issues," use the workflow below: identify gaps, draft candidate issues, get them labeled with the taxonomy in this file, and confirm with Andrew before filing. Do not fix the underlying problems in the same pass unless he explicitly asks you to — this workflow is about surfacing and triaging work, not doing it.

### Where to look

The project's own quality gates are the map of what "good" looks like — a gap is anywhere a gate is missing, weak, or would fail. Check each dimension against the actual tooling, not in the abstract:

| Dimension | Look at | Local command |
|---|---|---|
| Unit test gaps | Packages under `pkg/` and `cmd/` with thin or missing `_test.go` coverage | `go test -cover ./...` |
| CI issues | `.github/workflows/cicd.yaml`, `release.yml`, `deploy-docs.yml` — missing gates, flaky steps, drift from `Taskfile.yml` | `task ci` |
| Lint issues | `.golangci.yml` coverage, `golangci-lint` findings, shellcheck findings on `scripts/*.sh` | `task lint-go`, `task lint-shell` |
| Formatting issues | Unformatted Go files | `task fmt-check` |
| Command logic issues | Behavior in `cmd/*.go` that's wrong, inconsistent across subcommands, or diverges from README/Hugo docs | manual read + `go run main.go <cmd>` |
| Parallelization / streaming opportunities | Long-running subcommands (e.g. `status`, `bootstrap`) that serialize independent work (network calls, package manager queries, per-runtime checks) instead of fanning out or streaming partial results as they arrive | read `pkg/status`, `cmd/status.go`, `cmd/bootstrap.go` |
| Acceptance testing gaps | End-to-end CLI behavior not covered by `cmd/integration_test.go` / `cmd/commands_test.go` | `go test -run TestCLI ./cmd/...` |
| Performance testing gaps | Subcommands with no benchmark/perf test analogous to `cmd/status_perf_test.go` | `go test -bench . ./...` |
| External CI / release-pipeline issues | Failures that only show up in the real GitHub Actions run, not in `task ci` — macOS/Windows matrix legs, `docker` build+push, `publish-chocolatey`, `publish-apt`, `update-homebrew` in `release.yml`, `deploy-docs.yml` | `gh run list`, `gh run view <id> --log-failed` — see "After a push or release" below |

Reference points for how a human works this repo day-to-day: `Taskfile.yml` (`task ci` is the composite gate list) and `.tmuxinator/allbctl.yml` (the `watchers` window runs fmt/vet/build/test/acceptance/lint/secrets live, per-gate, as one pane each) — a missing pane or a gate absent from `task ci` is itself a signal of a coverage gap.

### After a push or release

`task ci` runs everything that can be checked locally, but it isn't the whole picture: `cicd.yaml` also runs the test matrix on `macos-latest` and `windows-latest` (not just Linux), and `release.yml` runs steps `task ci` can't touch at all — building all six platform binaries, pushing a Docker image to GHCR, creating the GitHub Release, and publishing to Chocolatey, the APT repo, and the Homebrew tap. Local gates being green does not mean the pipeline is green. After pushing to `main` or pushing a release tag (including via `task release:major` / `task release:minor` / `task release:patch`), check the triggered run instead of assuming it passed:

- `gh run list --branch main --limit 3` (after a push to main) or `gh run list --workflow=release.yml --limit 3` (after a tag push) to find the run.
- `gh run watch <run-id>` to follow it live, or `gh run view <run-id>` once it's done for a per-job summary.
- `gh run view <run-id> --log-failed` to pull logs for just the failed job(s) — cheaper than tailing the whole run.

A failure here is a real, workflow-verified issue even if `task ci` was clean locally — file it the same way as any other finding (`area:ci`, or `area:command-logic` if the root cause is in `cmd/*.go` rather than the workflow YAML), and note in the issue body which job/step failed and that it's a GH-Actions-only reproduction (matrix OS, Docker build, or a publish step), since that's not obvious from the title alone.

### Label taxonomy

Every issue this workflow files gets exactly one `area:*` label and exactly one `priority:*` label. Don't invent new label names — extend this table (and confirm with Andrew) before using anything outside it.

| Label | Color | Use for |
|---|---|---|
| `area:testing` | `1d76db` | Unit test coverage gaps |
| `area:ci` | `0e8a16` | CI/workflow config issues, gates missing from `task ci` |
| `area:lint` | `fbca04` | Linting issues or `.golangci.yml` config gaps |
| `area:format` | `fef2c0` | Formatting issues (`gofmt` etc.) |
| `area:command-logic` | `e99695` | Bugs or inconsistencies in command behavior |
| `area:performance` | `5319e7` | Perf issues, incl. parallelization/streaming opportunities |
| `area:acceptance-testing` | `c5def5` | Missing end-to-end/CLI-behavior test coverage |
| `area:performance-testing` | `bfd4f2` | Missing benchmark/perf tests |
| `priority:high` | `b60205` | Should be addressed soon |
| `priority:medium` | `d93f0b` | Normal priority |
| `priority:low` | `c2e0c6` | Nice-to-have |

The existing GitHub defaults (`bug`, `enhancement`, `documentation`, `good first issue`, `help wanted`, etc.) stay in use alongside these where they fit — the `area:*`/`priority:*` pair adds the categorization this workflow needs, it doesn't replace the defaults.

### Workflow

1. `gh label list` — confirm which labels from the table above already exist in the repo.
2. `gh label create <name> --color <hex> --description "<use>"` for any missing from the table. This is a repo-visible action; it's pre-approved as part of standing up this taxonomy, but if the table above changes later, confirm before creating new labels.
3. `gh issue list --label area:<x> --state open` (and a plain title search) to check a candidate issue isn't already filed before drafting a new one.
4. Draft the issue (title, body describing the gap and why it matters, proposed `area:*` + `priority:*` labels) and show it to Andrew. Wait for explicit go-ahead before running `gh issue create` — filing issues is visible, shared-state action, not a local edit.
5. Once approved, `gh issue create --title "..." --body "..." --label area:x --label priority:y`.

### Out of scope for this workflow

Vagrant-based macOS coverage is a known gap (Vagrant can't officially virtualize macOS guests under Apple's EULA) but is intentionally **not** part of this identify-and-label pass — don't file an issue for it unless separately asked.

