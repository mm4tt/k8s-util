#!/bin/bash

set -euo pipefail

cluster_name=mmat-cheaper-kubemark

num_fake_nodes=5000
node_size=8
num_nodes=63

echo "Node size: $node_size vCPU, num_nodes: $num_nodes"

export CLUSTER=$cluster_name
export KUBE_GCE_NETWORK=${CLUSTER}
export INSTANCE_PREFIX=${CLUSTER}
export KUBE_GCE_INSTANCE_PREFIX=${CLUSTER}

source $GOPATH/src/github.com/mm4tt/k8s-util/set-common-envs/set-common-envs.sh preset-e2e-kubemark-common
source $GOPATH/src/github.com/mm4tt/k8s-util/set-common-envs/set-common-envs.sh preset-e2e-kubemark-gce-scale

cd $GOPATH/src/k8s.io/kubernetes
go run hack/e2e.go -- \
    --gcp-project=$PROJECT \
    --gcp-zone=$ZONE \
    --gcp-node-size=n1-standard-$node_size \
    --gcp-nodes=$num_nodes \
    --provider=gce \
    --kubemark \
    --kubemark-nodes=$num_fake_nodes \
    --up \
    --test=false \
    --test-cmd=/bin/true

cd -
