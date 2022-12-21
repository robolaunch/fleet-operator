package fleet

import (
	"context"

	fleetv1alpha1 "github.com/robolaunch/fleet-operator/pkg/api/roboscale.io/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"
)

func (r *FleetReconciler) reconcileGetInstance(ctx context.Context, meta types.NamespacedName) (*fleetv1alpha1.Fleet, error) {
	instance := &fleetv1alpha1.Fleet{}
	err := r.Get(ctx, meta, instance)
	if err != nil {
		return &fleetv1alpha1.Fleet{}, err
	}

	return instance, nil
}

func (r *FleetReconciler) reconcileUpdateInstanceStatus(ctx context.Context, instance *fleetv1alpha1.Fleet) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		instanceLV := &fleetv1alpha1.Fleet{}
		err := r.Get(ctx, types.NamespacedName{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		}, instanceLV)

		if err == nil {
			instance.ResourceVersion = instanceLV.ResourceVersion
		}

		err1 := r.Status().Update(ctx, instance)
		return err1
	})
}
