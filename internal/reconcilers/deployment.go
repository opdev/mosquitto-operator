package reconcilers

import (
	"context"
	"io/fs"

	"github.com/opdev/mosquitto-operator/api/v1alpha1"
	"github.com/opdev/mosquitto-operator/internal/templates"

	"github.com/opdev/subreconciler"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func ReconcileDeployment(client client.Client, scheme *runtime.Scheme, fs fs.ReadFileFS) subreconciler.FnWithRequest {
	return func(ctx context.Context, req reconcile.Request) (*ctrl.Result, error) {
		var mosquitto v1alpha1.Mosquitto
		if r, err := getResource(ctx, client, req, &mosquitto); subreconciler.ShouldHaltOrRequeue(r, err) {
			return r, err
		}

		deployment, err := templates.ResourceFromTemplate[v1alpha1.Mosquitto, appsv1.Deployment](&mosquitto, "deployment", fs)
		if err != nil {
			return subreconciler.RequeueWithError(err)
		}

		if err := ctrl.SetControllerReference(&mosquitto, deployment, scheme); err != nil {
			return subreconciler.RequeueWithError(err)
		}

		if err := reconcileResource(ctx, client, deployment, &appsv1.Deployment{}); err != nil {
			return subreconciler.RequeueWithError(err)
		}

		return subreconciler.ContinueReconciling()
	}
}
