package reconcilers

import (
	"context"

	"github.com/imdario/mergo"
	"github.com/opdev/subreconciler"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

// Generic function to reconcile a Kubernetes resource
// `current` should be an empty resource (e.g. &appsv1.Deployment{}). It is
// populated by the actual current state of the resource in the initial GET
// request.
func reconcileResource(
	ctx context.Context,
	rclient client.Client,
	desired client.Object,
	current client.Object,
) error {
	log := ctrllog.FromContext(ctx)
	log.Info(
		"Reconciling child resource",
		"Kind", desired.GetObjectKind().GroupVersionKind().Kind,
		"Name", desired.GetName(),
		"Namespace", desired.GetNamespace(),
	)

	key := types.NamespacedName{Name: desired.GetName(), Namespace: desired.GetNamespace()}
	if err := rclient.Get(ctx, key, current); err != nil {
		if !k8serrors.IsNotFound(err) {
			log.Error(
				err,
				"Error reading child resource",
				"Kind", desired.GetObjectKind().GroupVersionKind().Kind,
				"Name", desired.GetName(),
				"Namespace", desired.GetNamespace(),
			)
			return err
		}

		log.Info(
			"Creating a new child resource",
			"Kind", desired.GetObjectKind().GroupVersionKind().Kind,
			"Name", desired.GetName(),
			"Namespace", desired.GetNamespace(),
		)

		err = rclient.Create(ctx, desired)
		if err != nil {
			log.Error(
				err,
				"Failed to create a new child resource",
				"Kind", desired.GetObjectKind().GroupVersionKind().Kind,
				"Name", desired.GetName(),
				"Namespace", desired.GetNamespace(),
			)
			return err
		}
	} else {
		log.Info(
			"Patching existing child resource",
			"Kind", desired.GetObjectKind().GroupVersionKind().Kind,
			"Name", desired.GetName(),
			"Namespace", desired.GetNamespace(),
		)

		// This ensures that the resource is patched only if there is a
		// difference between desired and current.
		patchDiff := client.MergeFrom(current.DeepCopyObject().(client.Object))
		if err := mergo.Merge(current, desired, mergo.WithOverride); err != nil {
			log.Error(err, "Error in merge")
			return err
		}

		if err := rclient.Patch(ctx, current, patchDiff); err != nil {
			log.Error(
				err,
				"Failed to patch child resource",
				"Kind", desired.GetObjectKind().GroupVersionKind().Kind,
				"Name", desired.GetName(),
				"Namespace", desired.GetNamespace(),
			)
			return err
		}
	}

	current = desired.DeepCopyObject().(client.Object)

	log.Info(
		"Finished reconciling resource",
		"Kind", current.GetObjectKind().GroupVersionKind().Kind,
		"Name", current.GetName(),
		"Namespace", current.GetNamespace(),
	)
	return nil
}

func GetResource(ctx context.Context, client client.Client, req ctrl.Request, resource client.Object) (*ctrl.Result, error) {
	log := ctrllog.FromContext(ctx)

	err := client.Get(ctx, req.NamespacedName, resource)
	if err != nil && k8serrors.IsNotFound(err) {
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue, and we can get them on deleted requests. We'll need to wait for
		// a new notification.
		log.Error(
			err,
			"cannot find resource - has it been deleted?",
			"Name", resource.GetName(),
			"Namespace", resource.GetNamespace(),
		)
		return subreconciler.DoNotRequeue()
	}
	if err != nil && !k8serrors.IsNotFound(err) {
		log.Error(err, "error fetching resource")
		return subreconciler.RequeueWithError(err)
	}

	return subreconciler.ContinueReconciling()
}
