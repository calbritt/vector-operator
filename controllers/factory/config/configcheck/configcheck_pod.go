package configcheck

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func createVectorConfigCheckPod(name, ns, image, hash string) *corev1.Pod {
	labels := labelsForVectorConfigCheck()

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getNameVectorConfigCheck(name, hash),
			Namespace: ns,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			ServiceAccountName: "vector-configcheck",
			Volumes:            generateVectorConfigCheckVolume(name, hash),
			SecurityContext:    &corev1.PodSecurityContext{},
			Containers: []corev1.Container{
				{
					Name:  "config-check",
					Image: image,
					Args:  []string{"validate", "/etc/vector/*.json"},
					Env:   generateVectorConfigCheckEnvs(),
					Ports: []corev1.ContainerPort{
						{
							Name:          "prom-exporter",
							ContainerPort: 9090,
							Protocol:      "TCP",
						},
					},
					VolumeMounts: generateVectorConfigCheckVolumeMounts(),
				},
			},
			RestartPolicy: "Never",
		},
	}

	return pod
}

func generateVectorConfigCheckVolume(name, hash string) []corev1.Volume {
	volume := []corev1.Volume{
		{
			Name: "config",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: getNameVectorConfigCheck(name, hash),
				},
			},
		},
		{
			Name: "data",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/var/lib/vector",
				},
			},
		},
		{
			Name: "var-log",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/var/log/",
				},
			},
		},
		{
			Name: "var-lib",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/var/lib/",
				},
			},
		},
		{
			Name: "procfs",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/proc",
				},
			},
		},
		{
			Name: "sysfs",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/sys",
				},
			},
		},
	}

	return volume
}

func generateVectorConfigCheckVolumeMounts() []corev1.VolumeMount {
	volumeMount := []corev1.VolumeMount{
		{
			Name:      "config",
			MountPath: "/etc/vector/",
		},
		{
			Name:      "data",
			MountPath: "/vector-data-dir",
		},
		{
			Name:      "var-log",
			MountPath: "/var/log/",
		},
		{
			Name:      "var-lib",
			MountPath: "/var/lib/",
		},
		{
			Name:      "procfs",
			MountPath: "/host/proc",
		},
		{
			Name:      "sysfs",
			MountPath: "/host/sys",
		},
	}

	return volumeMount
}

func generateVectorConfigCheckEnvs() []corev1.EnvVar {
	envs := []corev1.EnvVar{
		{
			Name: "VECTOR_SELF_NODE_NAME",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					APIVersion: "v1",
					FieldPath:  "spec.nodeName",
				},
			},
		},
		{
			Name: "VECTOR_SELF_POD_NAME",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					APIVersion: "v1",
					FieldPath:  "metadata.name",
				},
			},
		},
		{
			Name: "VECTOR_SELF_POD_NAMESPACE",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					APIVersion: "v1",
					FieldPath:  "metadata.namespace",
				},
			},
		},
		{
			Name:  "PROCFS_ROOT",
			Value: "/host/proc",
		},
		{
			Name:  "SYSFS_ROOT",
			Value: "/host/sys",
		},
	}

	return envs
}