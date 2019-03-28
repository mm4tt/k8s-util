# How Tos?

### How to update addon managed by [addon manager](https://github.com/kubernetes/kubernetes/tree/master/cluster/addons/addon-manager)?

1. Ssh into master
2. Edit the proper addon template in `/etc/kubernetes/addons`

### How to create your own [CoreDNS](https://github.com/coredns/corednskubernetes/kubernetes/tree/master/cluster/addons/addon-manager) release?

1. Make sure you use the same go version as in 'go.mod', e.g. if it specifies 1.12 run
   ```
   gvm use go1.12
   ```
1. `make -f Makefile.release  DOCKER=gcr.io/mmatejczyk-gke-dev release`
1. `make -f Makefile.release  DOCKER=gcr.io/mmatejczyk-gke-dev docker`
1. `gcloud auth print-access-token | docker login -u oauth2accesstoken --password-stdin https://gcr.io`
1. `make -f Makefile.release  DOCKER=gcr.io/mmatejczyk-gke-dev docker-push`
