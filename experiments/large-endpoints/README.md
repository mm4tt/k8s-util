# Large Endpoints in Large Clusters

Experiment details can be found in https://github.com/kubernetes/kubernetes/issues/75294

## How to run?

0. `cd <this_dir>`
1. Set-up cluster: `./set-up-cluster.sh`
1. Create test namespace: `kubectl create namespace test`
1. Monitor api-server on the master
    1. Create ssh tunnel to api-server 
        ```
        gcloud compute --project $PROJECT ssh --zone $ZONE "e2e-test-mmatejczyk-master" -- -L8081:localhost:8080
        ```
    2. Monitor via top (on master, you can use the above ssh session): 
       ```
       top | grep -e kube-apiserver -e etcd -e Cpu
       ```
    3. Pprof (localhost): `echo "exit" | go tool pprof -seconds 30 http://localhost:8081/debug/pprof/profile` 
1. Start watches:  `./start-watches.sh` 
 
  **TODO(mm4tt)**: Starting watches this way seems not to be working :( Right now a real big cluster is needed :(
1. Create large endpoints object: `kubectl apply -f large-endpoints.yaml`
1. Analyze pprof results with 
```
pprof -http=:8080 pprof.kube-apiserver.samples.cpu.005.pb
``` 
 

## Results

I saw CPU consumption grom from few cores up to 30 cores in kube-apiserver. It's still less than half of available resources
```
%Cpu(s): 39.2 us, 12.3 sy,  0.0 ni, 44.7 id,  0.2 wa,  0.0 hi,  3.5 si,  0.0 st
   3578 root      20   0   27.7g  27.5g  71512 S  3060  11.6 332:01.70 kube-apiserver                                                                                                                                                                                                     
   4124 root      20   0 7215652   1.7g 867736 S 308.3   0.7  43:06.67 etcd                                                                                                                                                                                                               
   4096 root      20   0 5523792 345976 132492 S   1.0   0.1   4:22.62 etcd  
``` 
