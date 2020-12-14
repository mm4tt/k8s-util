#!/bin/bash

set -euo pipefail

if [ $# -lt 1 ]
  then
    echo "Usage: ${0} <image> [gcp project]"
    exit 1
fi

image=${1?}
project=${2:-"mmatejczyk-gke-dev"}

echo "Pushing $image to gcr.io/$project/${image}"

docker pull "${image}"
docker tag "${image}" "gcr.io/$project/${image}"
docker push "gcr.io/$project/${image}"
