## md to pdf controller

### kubebuilder

```
kubebuilder init --domain example.com --repo example.com/pdfdocument
kubebuilder create api --group tools --version v1 --kind PdfDocument

```

### code

// api/v1/pdfdocument_types.go

```


        DocumentName string `json:"documentName,omitempty"`
        Text         string `json:"text,omitempty"`
```

// internal/controller/pdfdocument_controller.go

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
	// pdfdocumentv1 "susesecurity.com/pdfdocument/api/v1"
)
```

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

```
make run
```

```
# CRD file
# /home/jeff/myprojects/pdfdocument/config/crd/bases/pdfdocument.susesecurity.com_pdfdocuments.yaml
```

CR

```
laborant@dev-machine:~/projects/pdfdocument/config/crd/bases$ cat cr.yaml
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
# step - 1
laborant@dev-machine:~$  kubectl cp sample-document-job-nhqcj:/data/my-document.pdf ${PWD}/my-document.pdf
Defaulted container "main" out of: main, store-to-md (init), convert-to-pdf (init)
tar: removing leading '/' from member names

laborant@dev-machine:~$ ls -l *.pdf
-rw-rw-r-- 1 laborant laborant 46532 Mar  8 06:26 my-document.pdf
laborant@dev-machine:~$

# step - 2
# exit to WSL
jeff@SUSE-387793:/mnt/c/demo ()$ labctl cp 67cbbe726a18929a7ce141ec:~/my-document.pdf .
Warning: Permanently added '[127.0.0.1]:45386' (ED25519) to the list of known hosts.
Done!

```
