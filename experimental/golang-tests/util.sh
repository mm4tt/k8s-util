#!/bin/bash


log_dir=~/log/${run_name}
mkdir -p ${log_dir}
log_file=${log_dir}/log_$(date +%Y%m%d_%H%M%S)

log() { echo $1 | ts | tee -a ${log_file}; }

apply_patch() {
 cl_id=${1?}
 revision=${2?}

 echo "Applying patch ${cl_id} at revision ${revision}"

 wget https://go-review.googlesource.com/changes/go~${cl_id}/revisions/${revision}/patch?zip -O patch.zip
 unzip patch.zip && rm patch.zip
 git apply --3way *.diff
 rm *.diff
 git add .
 git commit -a -m "Applied ${cl_id} revision ${revision}"
}

build_golang() {
  cd ~/golang/go/src
  ./make.bash
  cd -
}

build_k8s() {
  log "Building k8s"

  cd $GOPATH/src/k8s.io/kubernetes
  git checkout $k8s_branch

  cd build/build-image/cross/
  rm -rf go || true
  cp -R ~/golang/go/ go

  echo "$run_name" > VERSION

  git add .
  git commit -a -m "Update golang version for run ${run_name}" || true

  make build

  cd -
  make clean quick-release
}

verify_run_name() {
 if ! echo "$run_name" | grep -Po "^(?:[a-z](?:[-a-z0-9]{0,61}[a-z0-9])?)$" 1>/dev/null; then
   echo "Invalid run name: '$run_name', doesn't match ^(?:[a-z](?:[-a-z0-9]{0,61}[a-z0-9])?)$"
   exit 1
 fi
}

run_kubemark() {
  source $GOPATH/src/github.com/mm4tt/k8s-util/set-common-envs/set-common-envs.sh preset-e2e-kubemark-common ${test_infra_commit}
  source $GOPATH/src/github.com/mm4tt/k8s-util/set-common-envs/set-common-envs.sh preset-e2e-kubemark-gce-scale ${test_infra_commit}

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
      "${kubetest_extra_args}" \
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
}

run_full() {
  source $GOPATH/src/github.com/mm4tt/k8s-util/set-common-envs/set-common-envs.sh preset-e2e-scalability-common ${test_infra_commit}

  go run hack/e2e.go -- \
      --gcp-project=$PROJECT \
      --gcp-zone=$ZONE \
      --cluster=$CLUSTER \
      --gcp-node-size=n1-standard-1 \
      --gcp-nodes=5000 \
      --provider=gce \
      --check-version-skew=false \
      --up \
      "${kubetest_extra_args}" \
      --test=false \
      --test-cmd=$GOPATH/src/k8s.io/perf-tests/run-e2e.sh \
      --test-cmd-args=cluster-loader2 \
      --test-cmd-args=--enable-prometheus-server=true \
      --test-cmd-args=--experimental-gcp-snapshot-prometheus-disk=true \
      --test-cmd-args=--experimental-prometheus-disk-snapshot-name="${run_name}" \
      --test-cmd-args=--nodes=5000 \
      --test-cmd-args=--provider=gce \
      --test-cmd-args=--report-dir=/tmp/${run_name}/artifacts \
      --test-cmd-args=--tear-down-prometheus-server=true \
      --test-cmd-args=--testconfig=$GOPATH/src/k8s.io/perf-tests/clusterloader2/testing/density/config.yaml \
      --test-cmd-args=--testconfig=$GOPATH/src/k8s.io/perf-tests/clusterloader2/testing/load/config.yaml \
      --test-cmd-args=--testoverrides=./testing/density/5000_nodes/override.yaml \
      --test-cmd-name=ClusterLoaderV2 2>&1 | ts | tee -a ${log_file}
}
