## Manipulate CRD and CR Part 1

### CRD

use databae as analogy, draw diagram

### CR

### use curl

```
# use chatgpt to create a CRD
# prompt
in k8s, help me to generate a CRD yaml. I just need a message with string type.  Ask me more info along the way.

# prompt
give me a corresponding CR yaml.

# prompt
I would like to use curl to do CRUD for this CRD. API server is stored in $KUBE_API.
Generate related curl commands.
```

### try to observe the content in etcd

```
# find etcd pod
kubectl get pod -n kube-system

# exec into it
kubectl exec -it etcd-cplane-01 -n kube-system -- sh

# set environment variables
export ETCDCTL_API=3
export ETCDCTL_CACERT=/etc/kubernetes/pki/etcd/ca.crt
export ETCDCTL_CERT=/etc/kubernetes/pki/etcd/server.crt
export ETCDCTL_KEY=/etc/kubernetes/pki/etcd/server.key
export ETCDCTL_ENDPOINTS=https://127.0.0.1:2379

# List all keys stored in etcd
etcdctl get "" --prefix --keys-only

# Get content given a key
etcdctl get /registry/serviceaccounts/ns1/default

# notes
    /registry/pods/ - Stores pod information
    /registry/deployments/ - Stores deployments
    /registry/services/ - Stores services
    /registry/nodes/ - Stores node information
    /registry/secrets/ - Stores secrets (encrypted if encryption is enabled)

```
