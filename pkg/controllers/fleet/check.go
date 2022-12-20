package fleet

import (
	"context"

	fleetv1alpha1 "github.com/robolaunch/fleet-operator/pkg/api/roboscale.io/v1alpha1"
	robotv1alpha1 "github.com/robolaunch/robot-operator/pkg/api/roboscale.io/v1alpha1"
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

func (r *FleetReconciler) reconcileCheckDiscoveryServer(ctx context.Context, instance *fleetv1alpha1.Fleet) error {

	discoveryServerQuery := &robotv1alpha1.DiscoveryServer{}
	err := r.Get(ctx, *instance.GetDiscoveryServerMetadata(), discoveryServerQuery)
	if err != nil && errors.IsNotFound(err) {
		instance.Status.DiscoveryServerStatus = fleetv1alpha1.DiscoveryServerInstanceStatus{}
	} else if err != nil {
		return err
	} else {
		instance.Status.DiscoveryServerStatus.Created = true
		instance.Status.DiscoveryServerStatus.Status = discoveryServerQuery.Status
	}

	return nil
}
