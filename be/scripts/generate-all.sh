#!/bin/bash
set -eo pipefail
source "$(dirname "$0")/_init.sh"

# generate all packages
go generate ./cmd/ggen ./...
