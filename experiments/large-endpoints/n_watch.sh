#!/bin/bash

N_WATCHES=5000

for i in {1..N_WATCHES}
do
   kubectl get --watch endpoints --all-namespaces > /dev/null &
done

echo "Started $N_WATCHES watches..."

wait
