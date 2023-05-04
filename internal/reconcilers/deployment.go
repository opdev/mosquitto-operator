package reconcilers

import (
	"context"

	"github.com/opdev/mosquitto-operator/api/v1alpha1"
	"github.com/opdev/mosquitto-operator/internal/templates"

	"github.com/opdev/subreconciler"
	appsv1 "k8s.io/api/apps/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ReconcileDeployment(mosquitto *v1alpha1.Mosquitto, client client.Client) subreconciler.Fn {
	return func(ctx context.Context) (*ctrl.Result, error) {
		deployment, err := templates.ResourceFromTemplate[v1alpha1.Mosquitto, appsv1.Deployment](mosquitto, "deployment")
		if err != nil {
			return subreconciler.RequeueWithError(err)
		}

		if err := ctrl.SetControllerReference(mosquitto, deployment, client.Scheme()); err != nil {
			return subreconciler.RequeueWithError(err)
		}

		if err := reconcileResource(ctx, client, deployment, &appsv1.Deployment{}); err != nil {
			return subreconciler.RequeueWithError(err)
		}

		return subreconciler.ContinueReconciling()
	}
}
