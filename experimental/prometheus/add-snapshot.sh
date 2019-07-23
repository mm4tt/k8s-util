#!/bin/bash

set -euo pipefail

prom_instance_name=${1?}
snapshot_name=${2?}
disk_name=${3:-$snapshot_name}

echo "Creating disk ${disk_name}"
gcloud compute disks create --project="${PROJECT}" --zone="${ZONE}" --source-snapshot="${snapshot_name}" --type=pd-ssd "${disk_name}"
echo "Created disk ${disk_name}"

echo "Attaching disk ${disk_name} to VM ${prom_instance_name}"
gcloud compute instances attach-disk ${prom_instance_name} --disk ${disk_name} --project=${PROJECT} --zone=${ZONE} --device-name=${disk_name}
echo "Attached disk ${disk_name}"

gcloud compute ssh --project="${PROJECT}" --zone="${ZONE}" "${prom_instance_name}" --command /bin/bash <<EOF
set -euo pipefail

nextPrometheusPort() {
  for i in {9091..10000} ; do
    if ! docker ps --no-trunc | grep "web.listen-address=:\$i" 1>/dev/null; then
      echo "\$i"
      return
    fi
  done
}

port=\$(nextPrometheusPort)
echo "Next port: \$port"

sudo mkdir /mnt/disks/${disk_name}
sudo mount /dev/disk/by-id/google-${disk_name} /mnt/disks/${disk_name}
sudo chmod -R 777 /mnt/disks/${disk_name}

docker run -d \
  --name=${disk_name} \
  --net=host \
  -v /mnt/disks/${disk_name}/prometheus-db:/prometheus-db \
  -v /tmp/prometheus.yml:/etc/prometheus/prometheus.yml \
  quay.io/prometheus/prometheus:v2.9.2 \
    --storage.tsdb.path=/prometheus-db \
    --storage.tsdb.retention=365d \
    --config.file=/etc/prometheus/prometheus.yml \
    --web.listen-address=:\${port}

echo "Sleeping 10s before adding datasource to grafana"
sleep 10

curl --header "Content-Type: application/json" \
  --request POST \
  --data "{\"name\":\"$disk_name\", \"type\":\"prometheus\", \"url\":\"http://localhost:\${port}\", \"access\":\"proxy\", \"basicAuth\":false}" \
  http://localhost:3000/api/datasources
EOF
