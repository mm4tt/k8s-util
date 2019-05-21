#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

docker rm nginx 2>/dev/null || true

docker run \
  --name=nginx \
  -p 80:80 \
  -v $DIR/reverse-proxy.conf:/etc/nginx/nginx.conf:ro \
  nginx:1.16.0 nginx-debug -g 'daemon off;'
