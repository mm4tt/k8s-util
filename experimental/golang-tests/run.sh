#!/bin/bash

set -euo pipefail

if [ $# -ne 1 ]
  then
    echo "Usage: ${0} <run_name>"
    echo ""
    echo "RUN NAMES SHOULD BE UNIQUE!"
    exit 1
fi

run_name=${1?}

config=${2:-$GOPATH/src/github.com/mm4tt/k8s-util/experimental/golang-tests/config.sh}
echo "Loading config: $config"
source $config
source $GOPATH/src/github.com/mm4tt/k8s-util/experimental/golang-tests/util.sh

verify_run_name

log "Running the ${run_name} test with ${num_nodes} nodes"

build_golang 2>&1 | ts | tee -a ${log_file}
build_k8s 2>&1 | ts | tee -a ${log_file}

log "k8s.io/perf-tests branch is: $perf_test_branch"
log "k8s.io/test-infra commit is: $test_infra_commit"

go install k8s.io/test-infra/kubetest

cd ~/go/src/k8s.io/perf-tests && git checkout ${perf_test_branch}
cd $GOPATH/src/k8s.io/kubernetes

source $GOPATH/src/github.com/mm4tt/k8s-util/set-common-envs/set-common-envs.sh preset-e2e-kubemark-common ${test_infra_commit}
source $GOPATH/src/github.com/mm4tt/k8s-util/set-common-envs/set-common-envs.sh preset-e2e-kubemark-gce-scale ${test_infra_commit}

export PROJECT=mmatejczyk-gke-dev
export ZONE=us-east1-b

export HEAPSTER_MACHINE_TYPE=n1-standard-32
export KUBE_DNS_MEMORY_LIMIT=300Mi

export CLUSTER=${run_name}
export KUBE_GCE_NETWORK=${CLUSTER}
export INSTANCE_PREFIX=${CLUSTER}
export KUBE_GCE_INSTANCE_PREFIX=${CLUSTER}

go run hack/e2e.go -- \
    --gcp-project=$PROJECT \
    --gcp-zone=$ZONE \
    --cluster=$CLUSTER \
    --gcp-node-size=n1-standard-8 \
    --gcp-nodes=50 \
    --provider=gce \
    --kubemark \
    --kubemark-nodes=$num_nodes \
    --check-version-skew=false \
    --up \
    --test=false \
    --test-cmd=$GOPATH/src/k8s.io/perf-tests/run-e2e.sh \
    --test-cmd-args=cluster-loader2 \
    --test-cmd-args=--enable-prometheus-server=true \
    --test-cmd-args=--experimental-gcp-snapshot-prometheus-disk=true \
    --test-cmd-args=--experimental-prometheus-disk-snapshot-name="${run_name}" \
    --test-cmd-args=--nodes=$num_nodes \
    --test-cmd-args=--provider=kubemark \
    --test-cmd-args=--report-dir=/tmp/${run_name}/artifacts \
    --test-cmd-args=--tear-down-prometheus-server=true \
    --test-cmd-args=--testconfig=$GOPATH/src/k8s.io/perf-tests/clusterloader2/testing/load/config.yaml \
    --test-cmd-args=--testoverrides=./testing/load/kubemark/throughput_override.yaml \
    --test-cmd-name=ClusterLoaderV2 2>&1 | ts | tee -a ${log_file}


$GOPATH/src/github.com/mm4tt/k8s-util/experimental/prometheus/add-snapshot.sh grafana $run_name