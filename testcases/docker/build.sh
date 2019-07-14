#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

DIR="$( cd -P "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

GIT_HEAD="$(git rev-parse --short=7 HEAD)"
GIT_DATE=$(git log HEAD -n1 --pretty='format:%cd' --date=format:'%Y%m%d-%H%M')

PROJECT=auto-test
PROJECT_TAG="${PROJECT}:${GIT_HEAD}-${GIT_DATE}"

BUILD_CONTEXT="${DIR}/.."
DOCKER_FILE="${DIR}/Dockerfile"

# Build docker image
docker build -f ${DOCKER_FILE} -t ${PROJECT_TAG} ${BUILD_CONTEXT}
