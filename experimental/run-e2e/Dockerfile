FROM golang:1.12.4

ARG kubernetes_clone_cmd="git clone https://github.com/kubernetes/kubernetes.git --single-branch --branch master"
ARG perftests_clone_cmd="git clone https://github.com/kubernetes/perf-tests.git --single-branch --branch master"

RUN apt-get update
RUN apt-get install moreutils -y
RUN go get github.com/mm4tt/k8s-util/set-common-envs
RUN mkdir /workspace

WORKDIR /go/src/k8s.io
RUN $kubernetes_clone_cmd
RUN $perftests_clone_cmd

WORKDIR /go/src/k8s.io/kubernetes
RUN make quick release

# TODO(mm4tt): -e doesn't work, fix it
ENV project mmatejczyk-gke-dev
ENV zone us-east1-b
ENV num_nodes 100
ENV cluster mmat-test

RUN touch run.sh
RUN chmod +x run.sh

RUN echo "#!/bin/bash\n" >> run.sh

RUN echo "source $GOPATH/src/github.com/mm4tt/k8s-util/set-common-envs/set-common-envs.sh\n" >> run.sh

RUN echo "go run hack/e2e.go -- \n\
        --gcp-project=$project \n\
        --gcp-zone=$zone \n\
        --cluster=$cluster \n\
        --gcp-nodes=$num_nodes \n\
        --provider=gce \n\
        --up \n\
        --test=false \n\
        --test-cmd=$GOPATH/src/k8s.io/perf-tests/run-e2e.sh \n\
        --test-cmd-args=cluster-loader2 \n\
        --test-cmd-args=--nodes=$num_nodes \n\
        --test-cmd-args=--provider=gce \n\
        --test-cmd-args=--enable-prometheus-server=true \n\
        --test-cmd-args=--tear-down-prometheus-server=false \n\
        --test-cmd-args=--report-dir=/workspace/artifacts \n\
        --test-cmd-args=--testconfig=$GOPATH/src/k8s.io/perf-tests/clusterloader2/testing/density/config.yaml \n\
        --test-cmd-name=ClusterLoaderV2  2>&1 | ts | tee /workspace/test_log" >> run.sh

ENTRYPOINT ./run.sh



