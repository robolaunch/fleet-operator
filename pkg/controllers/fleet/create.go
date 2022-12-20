package fleet

import (
	"context"

	"github.com/robolaunch/fleet-operator/internal/resources"
	fleetv1alpha1 "github.com/robolaunch/fleet-operator/pkg/api/roboscale.io/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *FleetReconciler) createNamespace(ctx context.Context, instance *fleetv1alpha1.Fleet, nsNamespacedName *types.NamespacedName) error {

	ns := resources.GetNamespace(instance, nsNamespacedName)

	err := ctrl.SetControllerReference(instance, ns, r.Scheme)
	if err != nil {
		return err
	}

	err = r.Create(ctx, ns)
	if err != nil && errors.IsAlreadyExists(err) {
		return nil
	} else if err != nil {
		return err
	}

	logger.Info("STATUS: Namespace " + ns.Name + " is created.")
	return nil
}
