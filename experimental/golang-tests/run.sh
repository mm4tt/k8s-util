#!/bin/bash

set -euo pipefail

cd ~/golang/go
commit=$(git log --pretty=format:'%h' -n 1)

run_name=${1:-golang-$commit}

echo "Run name is: $run_name"

config=${2:-$GOPATH/src/github.com/mm4tt/k8s-util/experimental/golang-tests/config.sh}
echo "Loading config: $config"
source $config
source $GOPATH/src/github.com/mm4tt/k8s-util/experimental/golang-tests/util.sh

verify_run_name

log "Running the ${run_name} test with ${num_nodes} nodes"

if $build_k8s; then
  build_golang 2>&1 | ts | tee -a ${log_file}
  build_k8s 2>&1 | ts | tee -a ${log_file}
fi

log "k8s.io/perf-tests branch is: $perf_test_branch"
log "k8s.io/test-infra commit is: $test_infra_commit"

go install k8s.io/test-infra/kubetest

cd ~/go/src/k8s.io/perf-tests && git checkout ${perf_test_branch}
cd $GOPATH/src/k8s.io/kubernetes


export HEAPSTER_MACHINE_TYPE=n1-standard-32
export KUBE_DNS_MEMORY_LIMIT=300Mi

export CLUSTER=${run_name}
export KUBE_GCE_NETWORK=${CLUSTER}
export INSTANCE_PREFIX=${CLUSTER}
export KUBE_GCE_INSTANCE_PREFIX=${CLUSTER}

GODEBUG=gctrace=0

cd $GOPATH/src/k8s.io/kubernetes

log "Go version is: $(go version)"

log "Starting test..."
retval=0
if ! ($test_to_run); then
  retval=1
fi

$GOPATH/src/github.com/mm4tt/k8s-util/experimental/prometheus/add-snapshot.sh grafana $run_name || true

cd ~/golang/go

exit $retval