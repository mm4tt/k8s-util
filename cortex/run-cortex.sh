#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

docker stop cortex 2>/dev/null || true
docker rm cortex 2>/dev/null || true

docker run -d \
  --name=cortex \
  -v $DIR/cortex.yaml:/etc/cortex.yaml:ro \
  -v /mnt/disks/cortex-db:/cortex-db \
  quay.io/cortexproject/cortex:master-88541de7 -config.file=/etc/cortex.yaml
