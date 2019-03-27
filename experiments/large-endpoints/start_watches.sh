#!/bin/bash

N_WATCHES=2000
for ((i=0; i < $N_WATCHES; i++))
do
   kubectl get --watch endpoints --all-namespaces > /dev/null &
done

echo "Started $N_WATCHES watches..."

wait
