#!/usr/bin/env bash

# run this script only once
if [[ -z $INIT ]]; then
    export INIT=1

    # just to make sure we are running in the root directory
    if ! grep rbot <README.md >/dev/null; then
        echo
        echo ERROR: please run the script in the project root directory
        exit 1
    fi

    dir=$(realpath "$(dirname "$0")"/..)
    export PROJECT_DIR="$dir"
fi
