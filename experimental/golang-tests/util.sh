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
