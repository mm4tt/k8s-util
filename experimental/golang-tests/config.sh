#!/bin/bash

export PROJECT=k8s-scale-testing-2
export ZONE=us-central1-b
test_to_run=run_full

#export PROJECT=mmatejczyk-gke-dev
#export ZONE=us-east1-b
#test_to_run=run_kubemark

build_k8s=true

###
# Things that usually shouldn' be changed (unless k8s doesn't build)
###

k8s_branch=golang_kubemark_932487c7440b05_no_patches
# Some newer golang commits require some patches to build k8s, if k8s stops building, uncomment the line below.
#k8s_branch=golang_kubemark_932487c7440b05

###
# Things that shouldn't be changed
###
num_nodes=2500
perf_test_branch=golang1.13
test_infra_commit=63eb09459
