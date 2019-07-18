#!/bin/bash

set -e

if [ $# -ne 2 ]
  then
    echo "Usage: ${0} <run-name> <num-node>"
    exit 1
fi

run_ name=${1?}
num_nodes=${2?}

perf_test_branch=etcd_tests
test_infra_commit=0bf601772

log_dir=~/log/etcd-test-kubemark/${run_ name}
mkdir -p ${log_dir}
log_file=${log_dir}/log_${num_nodes}

log() { echo $1 | ts | tee -a $log_file; }

log "Running the etcd kubemark test with ${num_nodes} nodes"
log "k8s.io/perf-tests branch is: $perf_test_branch"
log "k8s.io/test-infra commit is: $test_infra_commit"


go install k8s.io/test-infra/kubetest

cd ~/go/src/k8s.io/perf-tests && git checkout ${perf_test_branch} && cd -

source $GOPATH/src/github.com/mm4tt/k8s-util/set-common-envs/set-common-envs.sh preset-e2e-kubemark-common ${test_infra_commit}
source $GOPATH/src/github.com/mm4tt/k8s-util/set-common-envs/set-common-envs.sh preset-e2e-kubemark-gce-scale ${test_infra_commit}

export PROJECT=mmatejczyk-gke-dev
export ZONE=us-east1-b

export HEAPSTER_MACHINE_TYPE=n1-standard-32
export KUBE_DNS_MEMORY_LIMIT=300Mi

export CLUSTER=${run_ name}
export KUBE_GCE_NETWORK=${CLUSTER}
export INSTANCE_PREFIX=${CLUSTER}
export KUBE_GCE_INSTANCE_PREFIX=${CLUSTER}


go run hack/e2e.go -- \
    --gcp-project=$PROJECT \
    --gcp-zone=$ZONE \
    --cluster=$CLUSTER \
    --gcp-node-size=n1-standard-8 \
    --gcp-nodes=83 \
    --provider=gce \
    --kubemark \
    --kubemark-nodes=$NUM_NODES \
    --check-version-skew=false \
    --up \
    --down \
    --test=false \
    --test-cmd=$GOPATH/src/k8s.io/perf-tests/run-e2e.sh \
    --test-cmd-args=cluster-loader2 \
    --test-cmd-args=--enable-prometheus-server=true \
    --test-cmd-args=--experimental-gcp-snapshot-prometheus-disk=true \
    --test-cmd-args=--experimental-prometheus-disk-snapshot-name="${run_ name}_${num_nodes}" \
    --test-cmd-args=--nodes=$NUM_NODES \
    --test-cmd-args=--provider=kubemark \
    --test-cmd-args=--report-dir=~/log/golang1.13/artifacts \
    --test-cmd-args=--tear-down-prometheus-server=true \
    --test-cmd-args=--testconfig=$GOPATH/src/k8s.io/perf-tests/clusterloader2/testing/load/config.yaml \
    --test-cmd-args=--testoverrides=./testing/load/kubemark/throughput_override.yaml \
    --test-cmd-name=ClusterLoaderV2  2>&1 | ts | tee -a ${log_file}

