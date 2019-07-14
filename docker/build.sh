#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

DIR="$( cd -P "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

${DIR}/../app/docker/build.sh
${DIR}/../testcases/docker/build.sh
