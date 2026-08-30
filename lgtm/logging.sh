#!/usr/bin/env bash

# Custom logging helper for grafana/otel-lgtm.
# Overrides the default (which discards all service output to /dev/null) so
# each service writes to its own log file under /tmp.  The "sites" tmux window
# then tails each file in a dedicated pane.
#
# Mounted into the container at /otel-lgtm/logging.sh via docker-compose.yml.

set -euo pipefail

function run_with_logging() {
	local name=$1; shift
	local _envvar=$1; shift   # ignored — we always log to a per-service file

	# Derive a short, lowercase service slug from the first word of the name
	# e.g. "Grafana 13.0.1" -> grafana, "OpenTelemetry Collector 0.x" -> opentelemetry
	local slug
	slug=$(echo "$name" | awk '{print tolower($1)}')
	local logfile="/tmp/${slug}.log"

	echo "Running ${name} → ${logfile}"
	exec "$@" > "${logfile}" 2>&1
}
