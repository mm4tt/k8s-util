#!/bin/bash

docker stop grafana 2>/dev/null || true
docker rm grafana 2>/dev/null || true

docker run --name=grafana \
  -d \
  -p 3000:3000 \
  -e "GF_AUTH_ANONYMOUS_ENABLED=true" \
  -e "GF_AUTH_ANONYMOUS_ORG_ROLE=Admin" \
  grafana/grafana:6.0.0
