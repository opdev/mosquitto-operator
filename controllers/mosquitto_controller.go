/*
Copyright 2023 Brad P. Crochet.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	"github.com/opdev/mosquitto-operator/api/v1alpha1"
	"github.com/opdev/mosquitto-operator/internal/reconcilers"
	"github.com/opdev/mosquitto-operator/internal/templates"

	"github.com/opdev/subreconciler"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// MosquittoReconciler reconciles a Mosquitto object
type MosquittoReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=messaging.eclipse.org,resources=mosquittoes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=messaging.eclipse.org,resources=mosquittoes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=messaging.eclipse.org,resources=mosquittoes/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *MosquittoReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	var mosquitto v1alpha1.Mosquitto
	if r, err := reconcilers.GetResource(ctx, r.Client, req, &mosquitto); subreconciler.ShouldHaltOrRequeue(r, err) {
		return *r, err
	}

	fs := templates.Templates

	subreconcilers := []subreconciler.Fn{
		reconcilers.ReconcileConfigMap(&mosquitto, r.Client, fs),
		reconcilers.ReconcileDeployment(&mosquitto, r.Client, fs),
		reconcilers.ReconcileService(&mosquitto, r.Client, fs),
	}

	for _, r := range subreconcilers {
		res, err := r(ctx)
		if subreconciler.ShouldHaltOrRequeue(res, err) {
			return subreconciler.Evaluate(res, err)
		}
	}

	return subreconciler.Evaluate(subreconciler.DoNotRequeue())
}

// SetupWithManager sets up the controller with the Manager.
func (r *MosquittoReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Mosquitto{}).
		Complete(r)
}
