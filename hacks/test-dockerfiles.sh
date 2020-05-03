#!/usr/bin/env bash
set -eo pipefail

allbctl generate dockerfile ansible --stdout | docker build -
