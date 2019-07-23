#!/bin/bash

###
# Things to change for every run
###

# This name will be used in grafana as datasource name (and many other places).
# KEEP THE RUN NAMES UNIQUE!
run_name=golang-xxx

golang_commit=master
# Golang patches to apply in format: cl_id:revision separated by commas,
# e.g. golang_patches=186598:3,186599:2

golang_patches=""

###
# Things that usually shouldn' be changed (unless k8s doesn't build)
###

k8s_branch=golang_kubemark_932487c7440b05_no_patches
# Some newer golang commits require some patches to build k8s, if k8s stops building, uncomment the line below.
#k8s_branch=golang_kubemark_932487c7440b05

# Golang commits to revert, seprated by commas.
golang_revert_commits=""
# If you're testing golang commit after the one below we need to revert
# f1a8ca30fcaa91803c353999448f6f3a292f1db1 as it breaks k8s build.
# So if golang_kubemark_932487c7440b05 branch didn't help, uncomment the line below.
#golang_revert_commits="f1a8ca30fcaa91803c353999448f6f3a292f1db1"

###
# Things that shouldn't be changed
###
num_nodes=2500
perf_test_branch=golang1.13
test_infra_commit=63eb09459