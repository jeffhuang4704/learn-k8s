## markdown to pdf controller (hello world)

Extend the Kubernetes API server functionality by implementing a custom controller that allows users to submit a Markdown file and generate a corresponding PDF. Users can submit the task via kubectl by creating a custom resource.

### 1️⃣ create project using kubebuilder

<details><summary>...</summary>

```
mkdir -p ~/projects/pdfdocument
cd ~/projects/pdfdocument

kubebuilder init --domain example.com --repo example.com/pdfdocument
kubebuilder create api --group tools --version v1 --kind PdfDocument
```

</details>

### 2️⃣ code

<details><summary>...</summary>

🅰️ modify the data structure (schema)

`// api/v1/pdfdocument_types.go`

```
        DocumentName string `json:"documentName,omitempty"`
        Text         string `json:"text,omitempty"`
```

`// internal/controller/pdfdocument_controller.go`

```
import (
	"context"
	"encoding/base64"
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"                   // For PodTemplateSpec and PodSpec
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1" // For ObjectMeta
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	// toolsv1 "example.com/pdfdocument/api/v1"
)
```

🅱️ Reconcile

```
func (r *PdfDocumentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var pdfDoc toolsv1.PdfDocument
	if err := r.Get(ctx, req.NamespacedName, &pdfDoc); err != nil {
		log.Error(err, "unable to fetch PdfDocument")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	jobspec, err := r.createJob(pdfDoc)
	if err != nil {
		log.Error(err, "unable to create Job")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if err := r.Create(ctx, &jobspec); err != nil {
		log.Error(err, "unable to create Job")
	}

	return ctrl.Result{}, nil
}
```

```
func (r *PdfDocumentReconciler) createJob(pdfDoc toolsv1.PdfDocument) (batchv1.Job, error) {
	image := "knsit/pandoc"
	base64text := base64.StdEncoding.EncodeToString([]byte(pdfDoc.Spec.Text))

	// Create a new Job
	job := batchv1.Job{
		TypeMeta: metav1.TypeMeta{APIVersion: batchv1.SchemeGroupVersion.String(), Kind: "Job"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      pdfDoc.Name + "-job",
			Namespace: pdfDoc.Namespace,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					InitContainers: []corev1.Container{
						{
							Name:    "store-to-md",
							Image:   "alpine",
							Command: []string{"/bin/sh"},
							Args:    []string{"-c", fmt.Sprintf("echo %s | base64 -d >> /data/text.md", base64text)},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "data-volume",
									MountPath: "/data",
								},
							},
						},
						{
							Name:    "convert-to-pdf",
							Image:   image,
							Command: []string{"sh", "-c"},
							Args:    []string{"-c", fmt.Sprintf("pandoc -s -o /data/%s.pdf /data/text.md", pdfDoc.Spec.DocumentName)},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "data-volume",
									MountPath: "/data",
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:    "main",
							Image:   "alpine",
							Command: []string{"sh", "-c", "sleep 3600"},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "data-volume",
									MountPath: "/data",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "data-volume",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
		},
	}

	return job, nil
}
```

</details>

### 3️⃣ Run

<details><summary>...</summary>

```
laborant@dev-machine:~/projects/pdfdocument$ make run
/home/laborant/projects/pdfdocument/bin/controller-gen rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases
/home/laborant/projects/pdfdocument/bin/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
go fmt ./...
go vet ./...
go run ./cmd/main.go

```

Create CRD (just like product installation)

```
# CRD yaml
laborant@dev-machine:~/projects/pdfdocument$ vim ./config/crd/bases/tools.example.com_pdfdocuments.yaml

# create CRD
laborant@dev-machine:~/projects/pdfdocument$ kubectl apply -f ./config/crd/bases/tools.example.com_pdfdocuments.yaml
customresourcedefinition.apiextensions.k8s.io/pdfdocuments.tools.example.com created

# view CRD
laborant@dev-machine:~/projects/pdfdocument$ kubectl get crds
NAME                                         CREATED AT
pdfdocuments.tools.example.com               2025-03-08T07:05:12Z

# view CRD definition
laborant@dev-machine:~/projects/pdfdocument$ kubectl get crd pdfdocuments.tools.example.com -oyaml

# view CR
laborant@dev-machine:~/projects/pdfdocument$ kubectl get pdfdocuments.tools.example.com
No resources found in default namespace.
```

Create CR (just like user submit request)

```
// cr.yaml
apiVersion: tools.example.com/v1
kind: PdfDocument
metadata:
  name: sample-document
  namespace: default
spec:
  documentName: my-document
  text: |
    ### My document
    Hello **world** !
```

Copy PDF out

```
# check the pod
laborant@dev-machine:~/projects/pdfdocument$ kubectl get pods --watch
NAME                        READY   STATUS    RESTARTS   AGE
sample-document-job-nhqcj   1/1     Running   0          28m

# step - 1
laborant@dev-machine:~$  kubectl cp sample-document-job-nhqcj:/data/my-document.pdf ${PWD}/my-document.pdf
Defaulted container "main" out of: main, store-to-md (init), convert-to-pdf (init)
tar: removing leading '/' from member names

laborant@dev-machine:~$ ls -l *.pdf
-rw-rw-r-- 1 laborant laborant 46532 Mar  8 06:26 my-document.pdf
laborant@dev-machine:~$

# step - 2
# exit to WSL
cd /mnt/c/demo

# get playground id
labctl playground list

jeff@SUSE-387793:/mnt/c/demo ()$ labctl cp 67cbbe726a18929a7ce141ec:~/my-document.pdf .
Warning: Permanently added '[127.0.0.1]:45386' (ED25519) to the list of known hosts.
Done!

```

<p align="center">
  <img src="./materils/md2pdf2.png" width="70%">
</p>

</details>

### debug

```
// start via delve debugger
laborant@dev-machine:~/projects/pdfdocument/bin$ dlv exec manager
Type 'help' for list of commands.
(dlv)


(dlv) funcs Reconcile
....
example.com/pdfdocument/internal/controller.(*PdfDocumentReconciler).Reconcile

// set breakpoint on Reconcile()
(dlv)  b example.com/pdfdocument/internal/controller.(*PdfDocumentReconciler).Reconcile

// continue
(dlv) c
```

Change buld option if necessary, add `-gcflags=all="-N -l"`

```
// Makefile
.PHONY: build
  1 build: manifests generate fmt vet ## Build manager binary.
    go build -o bin/manager cmd/main.go

	go build -gcflags=all="-N -l" -o bin/manager cmd/main.go
```

### Reference

<details><summary>...</summary>

[Writing Kubernetes Controllers](https://www.youtube.com/watch?v=q7b23612pSc)

<p align="center">
  <img src="./materils/md2pdf1.png" width="70%">
</p>

</details>
