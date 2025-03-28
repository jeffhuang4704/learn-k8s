## Manipulate CRD and CR wihtout coding - Part 1

### 1️⃣ CRD (== database table schema), CR (== a row in database table)

<details><summary>...</summary>

In Kubernetes, a Custom **Resource Definition (CRD)** is like a **database schema**. It defines the structure of a custom resource, specifying the kind of data it holds and how it should be validated.

On the other hand, a **Custom Resource (CR)** is similar to **a row in a database table**. It represents an actual instance of the data defined by the CRD, with specific values filled in according to the schema.

</details>

### 2️⃣ Why, When and Who

<details><summary>...</summary>

**Why**
We use CRDs in Kubernetes to **extend** its functionality beyond the built-in resources like Pods, Services, or Deployments. CRDs allow us to define our own custom resources with specific fields and behavior, tailored to the needs of our application or infrastructure.

**When**
When a product is installed, it defines a CRD to introduce new resource types that Kubernetes understands.

Once the CRD is in place, various actors can manipulate the CRs.

**Who**

- a user might manually create or update a CR through kubectl
- a CI/CD pipeline could automatically update CRs
- UI
- Any actor with the proper permissions.

</details>

### 3️⃣ Create CRD (~ create database schema)

<details><summary>CRD sample</summary>

```
# use chatgpt to create a CRD
# prompt
in k8s, help me to generate a CRD yaml. I just need a message with string type.  Ask me more info along the way.
```

```
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: hellomessages.susesecurity.com
spec:
  group: susesecurity.com
  names:
    kind: HelloMessage
    listKind: HelloMessageList
    plural: hellomessages
    singular: hellomessage
  scope: Namespaced
  versions:
    - name: v1alpha1
      served: true
      storage: true
      schema:
        openAPIV3Schema:
          type: object
          properties:
            apiVersion:
              type: string
            kind:
              type: string
            metadata:
              type: object
            spec:
              type: object
              properties:
                message:
                  type: string

```

</details>

<details><summary>🍉 Once applied, several things happen..</summary>

1. Kubernetes API Server recognizes the new resource (TODO: use a kubectl get xxx -v 6)
2. New API endpoint is created
3. kubectl can now manage the resource (CRUD)
4. `etcd` stores data for the CR (Custom Resource).

- The Kubernetes API server stores instances of your custom resource in **etcd**.
- Even if no controller exists, the objects persist in etcd.

5. No Automatic Controller (Until You Implement One)

- Kubernetes does **not** automatically provide controllers for your CRD.
- You need to write a **custom controller** (e.g., using Kubebuilder) to manage the lifecycle of the resource.

6. RBAC Considerations

- By default, only cluster administrators can manage the CRD.
- You need to define **RBAC roles** if you want other users, services, or controllers to interact with the resource.

7. Validation and Defaulting (Optional)

- If you define a **schema** in your CRD (`spec.versions.schema.openAPIV3Schema`), Kubernetes will validate requests.
- You can also define **default values** and **conversion webhooks**.

</details>

### 4️⃣ How to generate a sample CR

<details><summary>CR sample</summary>

```
# prompt
give me a corresponding CR yaml.
```

```
apiVersion: susesecurity.com/v1alpha1
kind: HelloMessage
metadata:
  name: example-hellomessage
  namespace: default
spec:
  message: "Hello, Kubernetes!"

```

</details>

### 5️⃣ Observe the content in etcd

<details><summary>steps</summary>

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

👉 /registry/susesecurity.com/hellomessages/default/example-hellomessage

# Get content given a key
etcdctl get /registry/susesecurity.com/hellomessages/default/example-hellomessage

# notes
    /registry/pods/         - Stores pod information
    /registry/deployments/  - Stores deployments
    /registry/services/     - Stores services
    /registry/nodes/        - Stores node information
    /registry/secrets/      - Stores secrets (encrypted if encryption is enabled)
```

</details>

### 6️⃣ Use kubectl to do CRUD

```
1. If the CR format is incorrect—for example, due to a typo—the API server will reject it.
            spec:
              type: object
              properties:
                message:
                  type: string

2.
  NeuVector Controller =>  K8s api-server
  Consul  =>  etcd
```

```
kubectl apply  (add/update)
kubectl delete
```

### 7️⃣ Use curl to do CRUD

<details><summary>...</summary>

Extract API Server Endpoint and certs

```
export KUBE_API=$(kubectl config view --raw -o jsonpath='{.clusters[0].cluster.server}')

kubectl config view --raw -o jsonpath='{.users[0].user.client-certificate-data}' | base64 -d > ~/client.crt
kubectl config view --raw -o jsonpath='{.users[0].user.client-key-data}' | base64 -d > ~/client.key
kubectl config view --raw -o jsonpath='{.clusters[0].cluster.certificate-authority-data}' | base64 -d > ~/ca.crt
```

TODO: add curl examples..

```
# prompt

I want to use curl to watch the CR creation, API server is stored in $KUBE_API.
what's the curl command to use?

# prompt

I would like to use curl to do CRUD for this CRD. API server is stored in $KUBE_API.
Generate related curl commands.
```

</details>

### 8️⃣ Create more CRs

<details><summary>...</summary>

```
#!/bin/bash

for i in {1..20}
do
cat <<EOF > hellomessage-$i.yaml
apiVersion: susesecurity.com/v1alpha1
kind: HelloMessage
metadata:
  name: hellomessage-$i
  namespace: default
spec:
  message: "Hello, Kubernetes! This is message $i."
EOF
done
```

```
laborant@dev-machine:~/test$ ls -l
total 84
-rwxrwxr-x 1 laborant laborant 246 Mar 7 06:36 creaet_cr.sh
-rw-rw-r-- 1 laborant laborant 168 Mar 7 06:36 hellomessage-1.yaml
-rw-rw-r-- 1 laborant laborant 170 Mar 7 06:36 hellomessage-10.yaml
-rw-rw-r-- 1 laborant laborant 170 Mar 7 06:36 hellomessage-11.yaml
-rw-rw-r-- 1 laborant laborant 170 Mar 7 06:36 hellomessage-12.yaml
-rw-rw-r-- 1 laborant laborant 170 Mar 7 06:36 hellomessage-13.yaml
-rw-rw-r-- 1 laborant laborant 170 Mar 7 06:36 hellomessage-14.yaml
-rw-rw-r-- 1 laborant laborant 170 Mar 7 06:36 hellomessage-15.yaml
-rw-rw-r-- 1 laborant laborant 170 Mar 7 06:36 hellomessage-16.yaml
-rw-rw-r-- 1 laborant laborant 170 Mar 7 06:36 hellomessage-17.yaml
-rw-rw-r-- 1 laborant laborant 170 Mar 7 06:36 hellomessage-18.yaml
-rw-rw-r-- 1 laborant laborant 170 Mar 7 06:36 hellomessage-19.yaml
-rw-rw-r-- 1 laborant laborant 168 Mar 7 06:36 hellomessage-2.yaml
-rw-rw-r-- 1 laborant laborant 170 Mar 7 06:36 hellomessage-20.yaml
-rw-rw-r-- 1 laborant laborant 168 Mar 7 06:36 hellomessage-3.yaml
-rw-rw-r-- 1 laborant laborant 168 Mar 7 06:36 hellomessage-4.yaml
-rw-rw-r-- 1 laborant laborant 168 Mar 7 06:36 hellomessage-5.yaml
-rw-rw-r-- 1 laborant laborant 168 Mar 7 06:36 hellomessage-6.yaml
-rw-rw-r-- 1 laborant laborant 168 Mar 7 06:36 hellomessage-7.yaml
-rw-rw-r-- 1 laborant laborant 168 Mar 7 06:36 hellomessage-8.yaml
-rw-rw-r-- 1 laborant laborant 168 Mar 7 06:36 hellomessage-9.yaml
```

```
# apply
kubectl apply -f .

# get
laborant@dev-machine:~/test$ kubectl get hellomessages.susesecurity.com
NAME AGE
example-hellomessage 5h11m
hellomessage-1 31s
hellomessage-10 31s
hellomessage-11 31s
hellomessage-12 31s
hellomessage-13 31s
hellomessage-14 31s
hellomessage-15 31s
hellomessage-16 31s
hellomessage-17 31s
hellomessage-18 31s
hellomessage-19 31s
hellomessage-2 31s
hellomessage-20 31s
hellomessage-3 31s
hellomessage-4 31s
hellomessage-5 31s
hellomessage-6 31s
hellomessage-7 31s
hellomessage-8 31s
hellomessage-9 31s
laborant@dev-machine:~/test$

```

</details>

### 9️⃣ Custom controller

<details><summary>implementing a Controller...</summary>

If you want Kubernetes to take action when a CR is created, you must:

1. **Write a Controller** (using Kubebuilder, Operator SDK, or client-go).
2. **Watch for CR Events** (Create, Update, Delete).
3. **Reconcile Desired State** (Apply logic to manage resources based on the CR).
4. **Deploy the Controller as a Pod** inside the cluster.

The custom controller is responsible for **reconciling** the CR and ensuring the desired state is achieved

Follow Kubernetes' rules and best practices by understanding how it works. This includes key concepts like:

- Reconciliation loops: Continuously checking and fixing resources to match the desired state.
- Retries: Automatically trying again when an operation fails.
- Conflict management: Handling situations where multiple changes happen at the same time.

**Note:**

- A custom controller can be designed to process many kinds of Custom Resource Definitions (CRDs).
- A custom controller can process existing types of resources, not just custom ones (CRDs).

</details>
