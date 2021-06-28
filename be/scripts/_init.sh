#!/usr/bin/env bash

# run this script only once
if [[ -z $INIT ]]; then
    export INIT=1

    # just to make sure we are running in the backend directory
    if ! grep 'olvrng/rbot' <go.mod >/dev/null; then
        echo
        echo 'ERROR: please run the script in the backend directory (review-bot/be)'
        exit 1
    fi

    dir=$(realpath "$(dirname "$0")"/..)
    export PROJECT_DIR="$dir"
fi
