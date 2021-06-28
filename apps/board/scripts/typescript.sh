#!/usr/bin/env bash
set -eo pipefail
source "$(dirname "$0")/_init.sh"

rm build-ts/**/*.js || true

yarn tsc --watch
