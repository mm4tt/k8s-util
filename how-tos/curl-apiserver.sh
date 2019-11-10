#!/bin/bash

set -euo pipefail

SERVICE_ACCOUNT=api-explorer

cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ServiceAccount
metadata:
  name: ${SERVICE_ACCOUNT}
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: ${SERVICE_ACCOUNT}-view
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: view
subjects:
  - kind: ServiceAccount
    name: ${SERVICE_ACCOUNT}
    namespace: default
EOF

SECRET=$(kubectl get serviceaccount ${SERVICE_ACCOUNT} -o json | jq -Mr '.secrets[].name | select(contains("token"))')
TOKEN=$(kubectl get secret ${SECRET} -o json | jq -Mr '.data.token' | base64 -D)
kubectl get secret ${SECRET} -o json | jq -Mr '.data["ca.crt"]' | base64 -D > /tmp/ca.crt
APISERVER=https://$(kubectl -n default get endpoints kubernetes --no-headers | awk '{ print $2 }')




# curl calls go here

curl -s  --cacert /tmp/ca.crt \
  --header "Authorization: Bearer $TOKEN" \
  $APISERVER/api/v1/namespaces/default
echo


# Get size of response
system_pods_bytes=$(curl -s  --cacert /tmp/ca.crt \
  --header "Authorization: Bearer $TOKEN" \
  --header "Accept: application/vnd.kubernetes.protobuf" \
  -so /dev/null -w '%{size_download}' \
  $APISERVER/api/v1/namespaces/kube-system/pods)
echo "System pods bytes: $system_pods_bytes"