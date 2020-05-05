#!/bin/bash
set -eo pipefail

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

mkdir -p tmp
pushd "."
pushd "tmp"
trap "popd" EXIT
    "${SCRIPT_DIR}"/../bin/allbctl --help
    "${SCRIPT_DIR}"/../bin/allbctl
    "${SCRIPT_DIR}"/../bin/allbctl generate
    "${SCRIPT_DIR}"/../bin/allbctl generate ansible
    "${SCRIPT_DIR}"/../bin/allbctl generate ansible init
    "${SCRIPT_DIR}"/../bin/allbctl generate ansible config
    "${SCRIPT_DIR}"/../bin/allbctl generate ansible inventory
    "${SCRIPT_DIR}"/../bin/allbctl generate ansible hostVar
    "${SCRIPT_DIR}"/../bin/allbctl generate ansible groupVar
    "${SCRIPT_DIR}"/../bin/allbctl generate ansible role
    "${SCRIPT_DIR}"/../bin/allbctl generate dockerfile ansible
    "${SCRIPT_DIR}"/../bin/allbctl generate git
    "${SCRIPT_DIR}"/../bin/allbctl generate golang
    "${SCRIPT_DIR}"/../bin/allbctl generate java
    "${SCRIPT_DIR}"/../bin/allbctl generate kubernetes
    "${SCRIPT_DIR}"/../bin/allbctl generate node
    "${SCRIPT_DIR}"/../bin/allbctl generate python
    "${SCRIPT_DIR}"/../bin/allbctl generate ruby
    "${SCRIPT_DIR}"/../bin/allbctl generate scala
    "${SCRIPT_DIR}"/../bin/allbctl generate shell
    "${SCRIPT_DIR}"/../bin/allbctl youtube
    "${SCRIPT_DIR}"/../bin/allbctl youtube list
popd
rm -rf tmp
