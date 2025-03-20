# learn-k8s

### 2025/03/06 - K8s operator pattern and kubebuilder, by Sam

```
# slide
https://docs.google.com/presentation/d/17RvRa15re2C4BGS00vQ2oJSfPvnBUryYdt7ZWU1ao5I/edit#slide=id.g33984add883_0_0

```

### 2025/03/06 - Kubernetes study group, by Sam

Sam introduce the `Kubernetes APIs` concept in the [tutorial](https://github.com/gianlucam76/kubernetes-controller-tutorial/blob/main/docs/custom-resources.md)

You can find the recording in the Gmail invitation for this session.

You can find additional tutorials from the series at:

```
Kubernetes Controller
https://github.com/gianlucam76/kubernetes-controller-tutorial/blob/main/docs/reconciler.md

Concurrent Reconciling
https://github.com/gianlucam76/kubernetes-controller-tutorial/blob/main/docs/concurrent_reconciling.md
```

### 2025/03/20, Jeff

I took some notes on the basic concepts of Kubernetes custom controllers.
You can consider these are prerequistic for the material introduce by Sam.

1️⃣ The notes start with how to [access the Kubernetes API using curl](https://github.com/jeffhuang4704/learn-k8s/blob/main/notes/access-api-server.md)

2️⃣ Next, how ot [manipulate CRDs and CRs without coding](https://github.com/jeffhuang4704/learn-k8s/blob/main/notes/crd-cr-part-1.md)

3️⃣ Finally, there's an example [demonstrating how to create a simple, hello-world-level custom controller](https://github.com/jeffhuang4704/learn-k8s/blob/main/notes/pdf-document.md). The controller accepts markdown text as input and generates a PDF from the provided markdown.
