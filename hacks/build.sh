#!/bin/bash
set -eo pipefail

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
function cleanup() {
    popd
}

trap cleanup EXIT

pushd "${SCRIPT_DIR}"
pushd "${SCRIPT_DIR}/.."

pkger
go build -v -o bin/allbctl

cleanup
