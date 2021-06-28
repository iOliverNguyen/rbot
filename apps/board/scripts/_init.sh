#!/bin/bash

# run this script only once
if [[ -z $INIT ]]; then
    export INIT=1

    # just to make sure we are running in the backend directory
    if ! grep 'rbot-board' <package.json >/dev/null; then
        echo
        echo 'ERROR: please run the script in the apps/* directory (rbot/apps/*)'
        exit 1
    fi

    dir=$(realpath "$(dirname "$0")"/..)
    export PROJECT_DIR="$dir"
fi
