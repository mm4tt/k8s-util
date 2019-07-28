#!/bin/bash

if [ $# -ne 2 ]
  then
    echo "Usage: ${0} <cluster_name> <num_fake_nodes>"
    exit 1
fi

cluster_name=${1?}
num_fake_nodes=${2?}

node_size=4
if [[ $num_fake_nodes -gt 100 ]]; then
  node_size=8
fi

num_nodes=$((13*$num_fake_nodes/100/$node_size+2))

echo "Node size: $node_size, num_nodes: $num_nodes"

export CLUSTER=$cluster_name
export KUBE_GCE_NETWORK=${CLUSTER}
export INSTANCE_PREFIX=${CLUSTER}
export KUBE_GCE_INSTANCE_PREFIX=${CLUSTER}

source $GOPATH/src/github.com/mm4tt/k8s-util/set-common-envs/set-common-envs.sh preset-e2e-kubemark-common
if [[ $num_fake_nodes -ge 1000 ]]; then
  source $GOPATH/src/github.com/mm4tt/k8s-util/set-common-envs/set-common-envs.sh preset-e2e-kubemark-gce-scale
fi

go run $GOPATH/src/k8s.io/kubernetes/hack/e2e.go -- \
    --gcp-project=$PROJECT \
    --gcp-zone=$ZONE \
    --gcp-node-size=n1-standard-$node_size \
    --gcp-nodes=$num_nodes \
    --provider=gce \
    --kubemark \
    --kubemark-nodes=$num_fake_nodes \
    --up
