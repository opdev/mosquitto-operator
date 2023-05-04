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

func ReconcileService(mosquitto *v1alpha1.Mosquitto, client client.Client, fs fs.ReadFileFS) subreconciler.Fn {
	return func(ctx context.Context) (*ctrl.Result, error) {
		svc, err := templates.ResourceFromTemplate[v1alpha1.Mosquitto, corev1.Service](mosquitto, "service", fs)
		if err != nil {
			return subreconciler.RequeueWithError(err)
		}

		if err := ctrl.SetControllerReference(mosquitto, svc, client.Scheme()); err != nil {
			return subreconciler.RequeueWithError(err)
		}

		if err := reconcileResource(ctx, client, svc, &corev1.Service{}); err != nil {
			return subreconciler.RequeueWithError(err)
		}

		return subreconciler.ContinueReconciling()
	}
}

func FinalizeService() {

}
