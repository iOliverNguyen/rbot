#!/bin/bash
set -e

binDir=bin
mountDir="$1"
if [[ -n $mountDir ]]; then
    buildDir=/project
    if [[ "$mountDir" == "$buildDir" ]]; then
        echo "invalid MOUNT_DIR"
        exit 1
    fi

    # copy the source to buildDir and start building in the new location
    # for faster build

    binDir="$mountDir/bin"
    mkdir -p "$binDir"

    cd /
    rm -rf "$buildDir" || true
    cp -r "$mountDir" "$buildDir"
    cd "$buildDir"

    echo "source copied, start building..."
fi

if [[ -n $ENV_FILE ]]; then source "$ENV_FILE" ; fi

# build
go version
go build -o "$binDir/rbot-server" ./cmd/rbot-server
