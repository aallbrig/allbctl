#!/usr/bin/env bash
# Watches one GitHub Actions job until its latest run completes, printing a
# live step checklist plus a glanceable "[sha] job: ✓/✗" summary line.
#
# Invoked once per trigger by watchexec panes in .tmuxinator/allbctl.yml,
# which watch one of the .ci-watch/*.trigger sentinel files touched by
# githooks/pre-push (or, for deploy-docs.yml, chained from this script's own
# "Create GitHub Release" invocation below). Also runnable manually via
# `task watch:ci-job WORKFLOW=... JOB=...` or directly with positional args.
#
# WORKFLOW/JOB come via env vars (tmuxinator panes set them as a shell prefix
# before `watchexec`) or positional args — env vars matter here because
# watchexec rejoins its post-`--` argv into a shell string without
# re-quoting, so a JOB value with spaces/parens (e.g. "test (ubuntu-latest)")
# breaks if passed as a CLI arg after `--`; env vars sidestep that since
# they're never re-parsed as shell text.
#
# Usage: WORKFLOW=<file> JOB=<name> watch-ci-job.sh
#    or: watch-ci-job.sh <workflow-file> <job-name>
set -uo pipefail

WORKFLOW="${1:-${WORKFLOW:?usage: WORKFLOW=<file> JOB=<name> watch-ci-job.sh (or positional args)}}"
JOB="${2:-${JOB:?usage: WORKFLOW=<file> JOB=<name> watch-ci-job.sh (or positional args)}}"

STATE_DIR=".ci-watch/state"
POLL_INTERVAL=10
FIND_RETRIES=10
FIND_DELAY=3

mkdir -p "$STATE_DIR"

slug=$(printf '%s' "${WORKFLOW}__${JOB}" | tr -c 'A-Za-z0-9._-' '-')
state_file="$STATE_DIR/${slug}.last-sha"
last_sha=""
[ -f "$state_file" ] && last_sha=$(cat "$state_file")

echo "==> watching '$JOB' in $WORKFLOW"

# Find the latest run for this workflow. Retries briefly in case this fired
# right after a push, before GitHub has registered the run yet.
run_line=""
i=0
while [ "$i" -lt "$FIND_RETRIES" ]; do
  run_line=$(gh run list --workflow="$WORKFLOW" --limit 1 \
    --json databaseId,headSha \
    -q '.[0] | "\(.databaseId)\t\(.headSha)"' 2>/dev/null || true)
  [ -n "$run_line" ] && [ "$run_line" != "null" ] && break
  sleep "$FIND_DELAY"
  i=$((i + 1))
done

if [ -z "$run_line" ] || [ "$run_line" = "null" ]; then
  echo "no $WORKFLOW runs found yet"
  exit 0
fi

run_id=${run_line%%$'\t'*}
head_sha=${run_line#*$'\t'}
short_sha=${head_sha:0:7}

if [ "$head_sha" = "$last_sha" ]; then
  echo "[$short_sha] $JOB: already reported, nothing new"
  exit 0
fi

# Find this job's ID within the run. Retries briefly in case the job list
# hasn't populated yet for a run that was just created.
job_id=""
i=0
while [ "$i" -lt "$FIND_RETRIES" ]; do
  job_id=$(gh run view "$run_id" --json jobs \
    -q ".jobs[] | select(.name==\"$JOB\") | .databaseId" 2>/dev/null || true)
  [ -n "$job_id" ] && break
  sleep "$FIND_DELAY"
  i=$((i + 1))
done

if [ -z "$job_id" ]; then
  echo "[$short_sha] $JOB: job not found in run $run_id"
  exit 0
fi

echo "[$short_sha] $JOB: run $run_id started, watching..."

status="unknown"
while [ "$status" != "completed" ]; do
  status=$(gh run view "$run_id" --json jobs \
    -q ".jobs[] | select(.databaseId==$job_id) | .status" 2>/dev/null || true)
  [ "$status" = "completed" ] && break

  clear
  gh run view "$run_id" --job "$job_id" 2>/dev/null || true
  sleep "$POLL_INTERVAL"
done

job_view=$(gh run view "$run_id" --job "$job_id" 2>/dev/null || true)
conclusion=$(gh run view "$run_id" --json jobs \
  -q ".jobs[] | select(.databaseId==$job_id) | .conclusion" 2>/dev/null || true)

clear
printf '%s\n' "$job_view"

header=$(printf '%s\n' "$job_view" | grep "(ID $job_id)")
echo
echo "[$short_sha] $header"

if [ "$conclusion" != "success" ]; then
  echo
  echo "── failed step log ──"
  gh run view "$run_id" --job "$job_id" --log-failed 2>/dev/null || true
fi

echo "$head_sha" > "$state_file"

# deploy-docs.yml fires on a GitHub-side "release published" event with no
# local git action of its own — chain off this job instead of polling for it.
if [ "$JOB" = "Create GitHub Release" ] && [ "$conclusion" = "success" ]; then
  mkdir -p .ci-watch
  touch .ci-watch/deploy-docs.trigger
fi
