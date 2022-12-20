package fleet

import (
	"context"

	fleetv1alpha1 "github.com/robolaunch/fleet-operator/pkg/api/roboscale.io/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
)

func (r *FleetReconciler) reconcileCheckNamespace(ctx context.Context, instance *fleetv1alpha1.Fleet) error {

	namespaceQuery := &corev1.Namespace{}
	err := r.Get(ctx, *instance.GetNamespaceMetadata(), namespaceQuery)
	if err != nil && errors.IsNotFound(err) {
		instance.Status.NamespaceStatus = fleetv1alpha1.NamespaceStatus{}
	} else if err != nil {
		return err
	} else {
		instance.Status.NamespaceStatus.Created = true
	}

	return nil
}
