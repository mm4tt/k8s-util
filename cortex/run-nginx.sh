#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

docker stop nginx 2>/dev/null || true
docker rm nginx 2>/dev/null || true
docker network create mynet 2>/dev/null || true

docker run -d \
  --name=nginx \
  --network=cortex-network \
  -p 9090:9090 \
  -v $DIR/nginx.conf:/etc/nginx/nginx.conf:ro \
  nginx:1.16.0
