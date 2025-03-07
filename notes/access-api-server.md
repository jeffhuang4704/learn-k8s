##

## methods

-- curl
-- kubectl proxy
-- kubectl --raw

## curl

<details><summary>scripts</summary>

Extract API Server Endpoint

```
export KUBE_API=$(kubectl config view --raw -o jsonpath='{.clusters[0].cluster.server}')
```

Extract and Decode the Client Certificate / Client Key / CA Certificate

```
kubectl config view --raw -o jsonpath='{.users[0].user.client-certificate-data}' | base64 -d > ~/client.crt
kubectl config view --raw -o jsonpath='{.users[0].user.client-key-data}' | base64 -d > ~/client.key
kubectl config view --raw -o jsonpath='{.clusters[0].cluster.certificate-authority-data}' | base64 -d > ~/ca.crt
```

Use curl to Call the Kubernetes API

```
curl --cert ~/client.crt --key ~/client.key --cacert ~/ca.crt ${KUBE_API}/api
```

### create a deployment

HTTP POST

```
curl ${KUBE_API}/apis/apps/v1/namespaces/default/deployments \
  --cert ~/client.crt \
  --key ~/client.key \
  --cacert ~/ca.crt \
  -X POST \
  -H 'Content-Type: application/yaml' \
  -d '---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sleep
spec:
  replicas: 1
  selector:
    matchLabels:
      app: sleep
  template:
    metadata:
      labels:
        app: sleep
    spec:
      containers:
      - name: sleep
        image: curlimages/curl
        command: ["/bin/sleep", "365d"]
'

```

### get all objects in the default namespace

```
curl $KUBE_API/apis/apps/v1/namespaces/default/deployments \
  --cert ~/client.crt \
  --key ~/client.key \
  --cacert ~/ca.crt
```

### get an object by a name and a namespace

```
curl $KUBE_API/apis/apps/v1/namespaces/default/deployments/sleep \
  --cert ~/client.crt \
  --key ~/client.key \
  --cacert ~/ca.crt
```

### watch

```
curl ${KUBE_API}/apis/apps/v1/namespaces/default/deployments?watch=true \
  --cert ~/client.crt \
  --key ~/client.key \
  --cacert ~/ca.crt

```

</details>

### Reference

[How To Call Kubernetes API using Simple HTTP Client](https://iximiuz.com/en/posts/kubernetes-api-call-simple-http-client/)
