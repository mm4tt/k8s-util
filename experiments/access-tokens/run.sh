#!/bin/bash

set -euo pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

$GOPATH/src/k8s.io/perf-tests/clusterloader2/run-e2e.sh \
  --prometheus-scrape-kube-proxy=false \
  --provider=gke \
  --enable-prometheus-server=true \
  --tear-down-prometheus-server=false \
  --report-dir=/tmp/access-tokens \
  --testconfig=$DIR/manifests/cl2.yaml 2>&1 | tee /tmp/log
