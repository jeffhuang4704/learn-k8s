## olly agent

### install mock API receiver

I developed a mock API receiver to handle requests from Olly agents and store them in a local database.

<details><summary>deployment yaml</summary>

```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: receiver-deployment
  labels:
    app: receiver
spec:
  replicas: 1
  selector:
    matchLabels:
      app: receiver
  template:
    metadata:
      labels:
        app: receiver
    spec:
      containers:
        - name: receiver
          image: chihjenhuang/receiver-image:v5
          ports:
            - containerPort: 8443
          volumeMounts:
            - name: receiver-volume
              mountPath: /output
      volumes:
        - name: receiver-volume
          hostPath:
            path: /output
            type: DirectoryOrCreate
---
apiVersion: v1
kind: Service
metadata:
  name: receiver-service
spec:
  selector:
    app: receiver
  ports:
    - protocol: TCP
      port: 443
      targetPort: 8443
  type: NodePort
```

</details>

### helm installation - Install Olly Agents

If you installed the mock API receiver in the previous step, you can extract the IP address and service port and apply them to the stackstate.url below.
The following settings will enable the debug log for the agent.

```
laborant@dev-machine:~$ kubectl get nodes -owide
NAME        STATUS   ROLES           AGE   VERSION   INTERNAL-IP   EXTERNAL-IP   OS-IMAGE             KERNEL-VERSION   CONTAINER-RUNTIME
cplane-01   Ready    control-plane   27m   v1.32.2   172.16.0.2   👈 <none>        Ubuntu 24.04.2 LTS   5.10.230         containerd://1.7.25
node-01     Ready    <none>          27m   v1.32.2   172.16.0.3    <none>        Ubuntu 24.04.2 LTS   5.10.230         containerd://1.7.25
node-02     Ready    <none>          27m   v1.32.2   172.16.0.4    <none>        Ubuntu 24.04.2 LTS   5.10.230         containerd://1.7.25

laborant@dev-machine:~$ kubectl get svc
NAME               TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)         AGE
kubernetes         ClusterIP   10.96.0.1        <none>        443/TCP         28m
receiver-service   NodePort    10.104.159.164   <none>        443:32319/TCP   26m   👈
laborant@dev-machine:~$

```

Get the node IP an Node service port, format the URL like `https://172.16.0.2:32319`

<details><summary>helm...</summary>

```

helm repo add suse-observability https://charts.rancher.com/server-charts/prime/suse-observability
helm repo update

helm upgrade --install \
    --namespace suse-observability \
    --create-namespace \
    --set-string 'stackstate.apiKey'='1111998d-973b-4998-9f0b-08f52c1ebb75' \
    --set-string 'stackstate.cluster.name'='test0313' \
    --set-string 'stackstate.url'='https://172.16.0.2:32319/receiver/stsAgent' \  👈
    --set 'nodeAgent.skipKubeletTLSVerify'=true \
    --set 'nodeAgent.containers.agent.logLevel=debug' \
    --set 'nodeAgent.logLevel=debug' \
    --set 'nodeAgent.containers.processAgent.logLevel=debug' \
    --set 'nodeAgent.logLevel=debug' \
    --set 'clusterAgent.logLevel=debug' \
    --set 'checksAgent.logLevel=debug' \
suse-observability-agent suse-observability/suse-observability-agent

```

</details>

### node-agent container in node-agent pod

Listing pods in `suse-observability` namespace

```
laborant@dev-machine:~$ kubectl get pods -n suse-observability
NAME                                                      READY   STATUS    RESTARTS   AGE
suse-observability-agent-checks-agent-6ff9f65f6c-mz5mx    1/1     Running   0          52s
suse-observability-agent-cluster-agent-64db9dc6fd-jrhgx   1/1     Running   0          52s
suse-observability-agent-logs-agent-4d7qd                 1/1     Running   0          52s
suse-observability-agent-logs-agent-ncwlg                 1/1     Running   0          52s
suse-observability-agent-logs-agent-tfgdn                 1/1     Running   0          52s
suse-observability-agent-node-agent-6dcqh                 2/2     Running   0          52s
suse-observability-agent-node-agent-pxrhv                 2/2     Running   0          52s
suse-observability-agent-node-agent-rt5fh                 2/2     Running   0          52s


```

Exec into the `node-agent` container. Make sure to specify the container name `node-agent`, as the pod contains two containers.

```
laborant@dev-machine:~$ kubectl exec -it suse-observability-agent-node-agent-6dcqh -c node-agent -n suse-observability -- bash
root@node-01:/# cd /opt/stackstate-agent/bin/agent
root@node-01:/opt/stackstate-agent/bin/agent# ./agent --help

The Datadog Agent faithfully collects events and metrics and brings them
to Datadog on your behalf so that you can do something useful with your
monitoring and performance data.

Usage:
  agent [command]

Available Commands:
  check                 Run the specified check
  completion            Generate the autocompletion script for the specified shell
  config                Print the runtime configuration of a running agent
  configcheck           Print all configurations loaded & resolved of a running agent
  diagnose              Validate Agent installation, configuration and environment
  ...........
```

View `node-agent` container log. Make sure to specify the container name `node-agent`, as the pod contains two containers.

```
laborant@dev-machine:~$ kubectl logs -f suse-observability-agent-node-agent-6dcqh -c node-agent  -n suse-observability

2025-03-19 06:18:56 UTC | CORE | DEBUG | (pkg/collector/worker/check_logger.go:44 in CheckStarted) | check:file_handle | Running check...
2025-03-19 06:18:56 UTC | CORE | DEBUG | (pkg/collector/worker/check_logger.go:61 in CheckFinished) | check:file_handle | Done running check
2025-03-19 06:18:58 UTC | CORE | DEBUG | (pkg/collector/worker/check_logger.go:44 in CheckStarted) | check:cri | Running check...
```

### process-agent container in node-agent pod

Exec into the `process-agent` container. Make sure to specify the container name `process-agent`, as the pod contains two containers.

```
laborant@dev-machine:~$ kubectl exec -it suse-observability-agent-node-agent-6dcqh -c process-agent -n suse-observability -- bash
root@node-01:/opt/stackstate-agent/bin/agent# ls -l
total 164064
-rwxrwxr-x 1 root root 168000312 Mar 17 01:16 process-agent
root@node-01:/opt/stackstate-agent/bin/agent# ./process-agent --help
Usage of ./process-agent:
  -check string
        Run a specific check and print the results. Choose from: process, connections, realtime
  -config string
        Path to stackstate.yaml config (default "/etc/stackstate-agent/stackstate.yaml")
  -info
        Show info about running process agent and exit
  -pid string
        Path to set pidfile for process
  -version
        Print the version and exit
root@node-01:/opt/stackstate-agent/bin/agent#

```

View `process-agent` container log. Make sure to specify the container name `process-agent`, as the pod contains two containers.

```
laborant@dev-machine:~$ kubectl logs -f suse-observability-agent-node-agent-6dcqh -c process-agent -n suse-observability
2025-03-19 06:14:06 DEBUG (config.go:770) - Filter process: /sbin/agetty -o -p -- \u --keep-baud 115200,57600,38400,9600 - vt220 based on blacklist: ^/sbin/
2025-03-19 06:14:06 DEBUG (config.go:770) - Filter process: sshd: laborant@pts/0 based on blacklist: ^sshd:
2025-03-19 06:14:06 DEBUG (config.go:770) - Filter process: /usr/lib/systemd/systemd --user based on blacklist: ^/usr/lib/systemd/
2
```
