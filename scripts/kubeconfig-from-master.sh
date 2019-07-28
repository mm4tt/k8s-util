#!/bin/bash

set -euo pipefail

if [ $# -ne 3 ]
  then
    echo "Usage: ${0} <project> <zone> <master_name>"
    exit 1
fi

project=${1?}
zone=${2?}
master_name=${3?}

gcloud beta compute --project "$project" ssh --zone "$zone" "$master_name"  --command "curl -s -H 'Metadata-Flavor: Google' http://metadata.google.internal/computeMetadata/v1/instance/attributes/kubeconfig > /tmp/kubeconfig"
gcloud beta compute --project "$project" scp --zone "$zone" "$master_name":/tmp/kubeconfig /tmp/kubeconfig_$master_name

echo "Run:"
echo "export KUBECONFIG=/tmp/kubeconfig_$master_name"
