#!/bin/bash
set -eo pipefail
source "$(dirname "$0")/_init.sh"

yarn snowpack dev
