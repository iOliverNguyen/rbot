#!/bin/bash
set -eo pipefail
source "$(dirname "$0")/_init.sh"

USAGE="Usage: be-build.sh [COMMAND]

Commands:
  build     Build backend (rbot-server)
  docker    Use docker for build (useful when you need to build Linux version on Mac)
  help      Display this instruction
"
TARGET_DIR="/ws/rbot"

build_docker() {
    if docker ps -a | grep 'project_golang$' | grep Exited ; then
        docker start project_golang
    fi
    if ! docker ps | grep 'project_golang$' ; then
        docker run -d --name project_golang \
            -v "$PWD":"$TARGET_DIR" \
            -w "$TARGET_DIR" olvrng/golang-toolbox \
            sleep 3600
    fi

    # ENV_FILE: environment variables to pass into _build_inner.sh
    # it should be relative path from the project dir
    if [[ -n $ENV_FILE ]]; then _env_file="-e=ENV_FILE=$ENV_FILE" ; fi

    docker exec -it -e COMMIT="$COMMIT" $_env_file \
        project_golang scripts/_build_inner.sh "$TARGET_DIR"
}

case "$1" in
""|build)
    if [[ -n $ENV_FILE ]]; then source "$ENV_FILE" ; fi
    scripts/_build_inner.sh
    ;;
docker)
    build_docker
    ;;
*)
    echo "$USAGE"
    exit 2
esac
