#!/bin/bash
set -eo pipefail

allbctl --help
allbctl
allbctl generate
allbctl generate ansible
allbctl generate ansible init
allbctl generate ansible config
allbctl generate ansible inventory
allbctl generate ansible hostVar
allbctl generate ansible groupVar
allbctl generate ansible role
allbctl generate dockerfile ansible
allbctl generate git
allbctl generate golang
allbctl generate java
allbctl generate kubernetes
allbctl generate node
allbctl generate python
allbctl generate ruby
allbctl generate scala
allbctl generate shell
