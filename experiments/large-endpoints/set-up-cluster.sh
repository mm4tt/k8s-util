#!/bin/bash


make quick-release && \
MASTER_SIZE=n1-standard-64 \
HEAPSTER_MACHINE_TYPE=n1-standard-4 \
NUM_NODES=50 \
go run hack/e2e.go -- \
  --gcp-zone=$ZONE \
  --gcp-project=$PROJECT \
  --cluster=$CLUSTER \
  --provider=$PROVIDER \
  --up
