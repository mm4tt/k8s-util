#!/bin/bash

set -euo pipefail

if [ $# -ne 1 ]; then
  echo "Usage: ${0} <job_name>"
  exit 1
fi
job_name=${1?}

export KUBECONFIG=$HOME/kubeconfigs/scalability_prow
unset CLOUDSDK_API_ENDPOINT_OVERRIDES_CONTAINER 

gcloud container clusters --project=gke-scalability-prow get-credentials prow --region=us-central1

cd $GOPATH/src/k8s.io/test-infra
export GO111MODULE=on
go run prow/cmd/mkpj/main.go \
  --job="${job_name}" \
  --job-config-path=$GOPATH/src/k8s.io/test-infra/config/jobs/kubernetes/sig-scalability/ \
  --config-path=$GOPATH/src/gke-internal/test-infra/prow/gke-scalability-prow/config.yaml | tee /tmp/${job_name}.yaml

kubectl apply -f /tmp/${job_name}.yaml
cd -
