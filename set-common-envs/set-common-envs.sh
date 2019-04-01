#!/bin/bash

# This file exports env variables and should be run with: source .../set-common-envs.sh
# Running it without source will span a new bash process that won't modify parent env variables.

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
readonly PRESET_NAME=${1:-preset-e2e-scalability-common}

eval $(go run $DIR/print-common-envs.go --preset-name=$PRESET_NAME)
