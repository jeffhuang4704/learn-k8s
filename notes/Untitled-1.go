

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