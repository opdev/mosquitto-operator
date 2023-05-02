package reconcilers

import (
	"context"
	"io/fs"

	"github.com/opdev/subreconciler"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ReconcileAuthSecret(client client.Client, fs fs.ReadFileFS) subreconciler.FnWithRequest {
	return func(ctx context.Context, r controllerruntime.Request) (*controllerruntime.Result, error) {
		return subreconciler.ContinueReconciling()
	}
}
