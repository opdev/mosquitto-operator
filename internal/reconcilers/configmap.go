package reconcilers

import (
	"context"
	"io/fs"

	"github.com/opdev/mosquitto-operator/api/v1alpha1"
	"github.com/opdev/mosquitto-operator/internal/templates"
	"github.com/opdev/subreconciler"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ReconcileConfigMap(mosquitto *v1alpha1.Mosquitto, client client.Client, fs fs.ReadFileFS) subreconciler.Fn {
	return func(ctx context.Context) (*ctrl.Result, error) {
		configmap, err := templates.ResourceFromTemplate[v1alpha1.Mosquitto, corev1.ConfigMap](mosquitto, "configmap", fs)
		if err != nil {
			return subreconciler.RequeueWithError(err)
		}

		if mosquitto.Status.MosquittoConfConfigMap != "" {
			configmap.Name = mosquitto.Status.MosquittoConfConfigMap
		}

		if err := ctrl.SetControllerReference(mosquitto, configmap, client.Scheme()); err != nil {
			return subreconciler.RequeueWithError(err)
		}

		if err := reconcileResource(ctx, client, configmap, &corev1.ConfigMap{}); err != nil {
			return subreconciler.RequeueWithError(err)
		}

		mosquitto.Status.MosquittoConfConfigMap = configmap.Name
		if err := client.Status().Update(ctx, mosquitto); err != nil {
			return subreconciler.DoNotRequeue()
		}

		return subreconciler.ContinueReconciling()
	}
}
