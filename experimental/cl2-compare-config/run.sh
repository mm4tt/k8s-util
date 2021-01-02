#!/bin/bash

set -euo pipefail

if [ $# -ne 1 ]
  then
    echo "Usage: ${0} <pr-number>"
    echo "The tool compares flattened load test config between master and the PR branch"
    exit 1
fi

pr_num=${1?}

tmp_dir=$(mktemp -d -t cl2-compare-XXXXXXXXXX)
cd "$tmp_dir"

echo "Cloning perf-test master and $pr_num PR"
mkdir config && cd config
git clone https://github.com/kubernetes/perf-tests.git
cd perf-tests
git fetch origin "pull/$pr_num/head:pr$pr_num"

echo "Cloning mm4tt's cl2 copy for dumping test config"
cd "$tmp_dir"
mkdir cl2 && cd cl2
git clone --single-branch --branch dump_cl2_steps https://github.com/mm4tt/perf-tests.git
cd perf-tests/clusterloader2

echo "Dumping config steps for master"
cd "$tmp_dir/config/perf-tests"
git checkout master
cd "$tmp_dir/cl2/perf-tests/clusterloader2"
go run cmd/clusterloader.go --testconfig testing/load/config.yaml --report-dir="$tmp_dir" --provider=gce --kubeconfig=/dev/null
cd "$tmp_dir" && mv steps.yaml master-steps.yaml

echo "Dumping config steps for PR $pr_num"
cd "$tmp_dir/config/perf-tests"
git checkout "pr$pr_num"
cd "$tmp_dir/cl2/perf-tests/clusterloader2"
go run cmd/clusterloader.go --testconfig testing/load/config.yaml --report-dir="$tmp_dir" --provider=gce --kubeconfig=/dev/null
cd "$tmp_dir" && mv steps.yaml "pr$pr_num-steps.yaml"

echo "Master config dumped in $tmp_dir/master-steps.yaml"
echo "PR $pr_num config dumped in $tmp_dir/pr$pr_num-steps.yaml"

cd "$tmp_dir"
diff  master-steps.yaml "pr$pr_num-steps.yaml" > diff

echo "Diff stored in $tmp_dir/diff"
echo "Diff is:"
cat diff
