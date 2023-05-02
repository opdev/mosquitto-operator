package reconcilers

import (
	"context"
	"io/fs"

	"github.com/opdev/mosquitto-operator/api/v1alpha1"
	"github.com/opdev/mosquitto-operator/internal/templates"
	"github.com/opdev/subreconciler"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func ReconcileConfigMap(client client.Client, scheme *runtime.Scheme, fs fs.ReadFileFS) subreconciler.FnWithRequest {
	return func(ctx context.Context, req reconcile.Request) (*ctrl.Result, error) {
		var mosquitto v1alpha1.Mosquitto
		if r, err := getResource(ctx, client, req, &mosquitto); subreconciler.ShouldHaltOrRequeue(r, err) {
			return r, err
		}

		configmap, err := templates.ResourceFromTemplate[v1alpha1.Mosquitto, corev1.ConfigMap](&mosquitto, "configmap", fs)
		if err != nil {
			return subreconciler.RequeueWithError(err)
		}

		if err := ctrl.SetControllerReference(&mosquitto, configmap, scheme); err != nil {
			return subreconciler.RequeueWithError(err)
		}

		if err := reconcileResource(ctx, client, configmap, &corev1.ConfigMap{}); err != nil {
			return subreconciler.RequeueWithError(err)
		}

		return subreconciler.ContinueReconciling()
	}
}
