/*
Copyright 2022.

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

package fleet

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/go-logr/logr"
	fleetv1alpha1 "github.com/robolaunch/fleet-operator/pkg/api/roboscale.io/v1alpha1"
	robotv1alpha1 "github.com/robolaunch/robot-operator/pkg/api/roboscale.io/v1alpha1"
)

// FleetReconciler reconciles a Fleet object
type FleetReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=fleet.roboscale.io,resources=fleets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=fleet.roboscale.io,resources=fleets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=fleet.roboscale.io,resources=fleets/finalizers,verbs=update

//+kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=robot.roboscale.io,resources=discoveryservers,verbs=get;list;watch;create;update;patch;delete

var logger logr.Logger

func (r *FleetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger = log.FromContext(ctx)

	instance, err := r.reconcileGetInstance(ctx, req.NamespacedName)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	err = r.reconcileCheckStatus(ctx, instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = r.reconcileUpdateInstanceStatus(ctx, instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = r.reconcileCheckResources(ctx, instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = r.reconcileUpdateInstanceStatus(ctx, instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *FleetReconciler) reconcileCheckStatus(ctx context.Context, instance *fleetv1alpha1.Fleet) error {

	switch instance.Status.NamespaceStatus.Created {
	case true:

		switch instance.Status.DiscoveryServerStatus.Created {
		case true:

			switch instance.Status.DiscoveryServerStatus.Status.Phase {
			case robotv1alpha1.DiscoveryServerPhaseReady:

				instance.Status.Phase = fleetv1alpha1.FleetPhaseReady

			}

		case false:

			instance.Status.Phase = fleetv1alpha1.FleetPhaseCreatingDiscoveryServer
			err := r.createDiscoveryServer(ctx, instance, instance.GetDiscoveryServerMetadata())
			if err != nil {
				return err
			}
			instance.Status.DiscoveryServerStatus.Created = true

		}

	case false:

		instance.Status.Phase = fleetv1alpha1.FleetPhaseCreatingNamespace
		err := r.createNamespace(ctx, instance, instance.GetNamespaceMetadata())
		if err != nil {
			return err
		}
		instance.Status.NamespaceStatus.Created = true

	}

	return nil
}

func (r *FleetReconciler) reconcileCheckResources(ctx context.Context, instance *fleetv1alpha1.Fleet) error {

	err := r.reconcileCheckNamespace(ctx, instance)
	if err != nil {
		return err
	}

	err = r.reconcileCheckDiscoveryServer(ctx, instance)
	if err != nil {
		return err
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *FleetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&fleetv1alpha1.Fleet{}).
		Owns(&corev1.Namespace{}).
		Owns(&robotv1alpha1.DiscoveryServer{}).
		Complete(r)
}
