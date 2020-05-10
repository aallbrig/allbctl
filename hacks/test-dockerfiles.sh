#!/usr/bin/env bash
set -eo pipefail

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

"${SCRIPT_DIR}"/../bin/allbctl generate dockerfile --name=Ansible --stdout | docker build -
"${SCRIPT_DIR}"/../bin/allbctl generate dockerfile --name=Alpine --stdout | docker build -
