#!/bin/bash
set -eo pipefail

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
function cleanup() {
    popd
    rm -rf tmp
}

trap cleanup EXIT

mkdir -p tmp
pushd "."
pushd "tmp"

"${SCRIPT_DIR}"/../bin/allbctl --help
"${SCRIPT_DIR}"/../bin/allbctl
"${SCRIPT_DIR}"/../bin/allbctl completion
"${SCRIPT_DIR}"/../bin/allbctl completion zsh
"${SCRIPT_DIR}"/../bin/allbctl generate
"${SCRIPT_DIR}"/../bin/allbctl generate ansible
"${SCRIPT_DIR}"/../bin/allbctl generate ansible init --interactive=false
"${SCRIPT_DIR}"/../bin/allbctl generate ansible config --interactive=false
"${SCRIPT_DIR}"/../bin/allbctl generate ansible inventory --interactive=false
"${SCRIPT_DIR}"/../bin/allbctl generate ansible hostVar --interactive=false
"${SCRIPT_DIR}"/../bin/allbctl generate ansible groupVar --interactive=false
"${SCRIPT_DIR}"/../bin/allbctl generate ansible role --interactive=false
"${SCRIPT_DIR}"/../bin/allbctl generate dockerfile --name=Ansible --interactive=false
"${SCRIPT_DIR}"/../bin/allbctl generate dockerfile --name=Alpine --interactive=false
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
"${SCRIPT_DIR}"/../bin/allbctl youtube playlists
"${SCRIPT_DIR}"/../bin/allbctl youtube videos

cleanup
