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
	goErr "errors"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/go-logr/logr"
	fleetErr "github.com/robolaunch/fleet-operator/internal/error"
	"github.com/robolaunch/fleet-operator/internal/label"
	fleetv1alpha1 "github.com/robolaunch/fleet-operator/pkg/api/roboscale.io/v1alpha1"
	robotv1alpha1 "github.com/robolaunch/robot-operator/pkg/api/roboscale.io/v1alpha1"
)

// FleetReconciler reconciles a Fleet object
type FleetReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	DynamicClient dynamic.Interface
}

//+kubebuilder:rbac:groups=fleet.roboscale.io,resources=fleets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=fleet.roboscale.io,resources=fleets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=fleet.roboscale.io,resources=fleets/finalizers,verbs=update

//+kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=robot.roboscale.io,resources=discoveryservers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=robot.roboscale.io,resources=robots,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=types.kubefed.io,resources=federatednamespaces,verbs=get;list;watch;create;update;patch;delete

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

	if instance.Spec.Hybrid && label.GetInstanceType(instance) == label.InstanceTypePhysicalInstance {
		err = r.reconcileCheckRemoteNamespace(ctx, instance)
		if err != nil {
			var e *fleetErr.NamespaceNotFoundError
			if goErr.As(err, &e) {
				logger.Info("STATUS: Searching for namespace.")
				return ctrl.Result{
					Requeue:      true,
					RequeueAfter: 3 * time.Second,
				}, nil
			}
			return ctrl.Result{}, nil
		}
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

	switch instance.Status.NamespaceStatus.Ready {
	case true:

		switch instance.Status.DiscoveryServerStatus.Created {
		case true:

			switch instance.Status.DiscoveryServerStatus.Phase {
			case string(robotv1alpha1.DiscoveryServerPhaseReady):

				instance.Status.Phase = fleetv1alpha1.FleetPhaseReady

				err := r.reconcileHandleAttachments(ctx, instance)
				if err != nil {
					return err
				}

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

		switch instance.Spec.Hybrid {
		case true:

			switch label.GetInstanceType(instance) {
			case label.InstanceTypeCloudInstance:

				switch instance.Status.NamespaceStatus.Created {
				case true:

					switch instance.Status.NamespaceStatus.Federated {
					case false:

						instance.Status.Phase = fleetv1alpha1.FleetPhaseCreatingNamespace
						err := r.createFederatedNamespace(ctx, instance, instance.GetNamespaceMetadata())
						if err != nil {
							return err
						}
						instance.Status.NamespaceStatus.Federated = true
						instance.Status.NamespaceStatus.Ready = true

					}

				case false:

					instance.Status.Phase = fleetv1alpha1.FleetPhaseCreatingNamespace
					err := r.createNamespace(ctx, instance, instance.GetNamespaceMetadata())
					if err != nil {
						return err
					}
					instance.Status.NamespaceStatus.Created = true

				}

			case label.InstanceTypePhysicalInstance:

				// do nothing

			}

		case false:

			instance.Status.Phase = fleetv1alpha1.FleetPhaseCreatingNamespace
			err := r.createNamespace(ctx, instance, instance.GetNamespaceMetadata())
			if err != nil {
				return err
			}
			instance.Status.NamespaceStatus.Created = true
			instance.Status.NamespaceStatus.Ready = true

		}

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

	err = r.reconcileCheckAttachedRobots(ctx, instance)
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
		Watches(
			&source.Kind{Type: &robotv1alpha1.Robot{}},
			handler.EnqueueRequestsFromMapFunc(r.watchRobots),
		).
		Complete(r)
}

func (r *FleetReconciler) watchRobots(o client.Object) []reconcile.Request {

	obj := o.(*robotv1alpha1.Robot)

	robot := &fleetv1alpha1.Fleet{}
	err := r.Get(context.TODO(), types.NamespacedName{
		Name:      label.GetTargetFleet(obj),
		Namespace: obj.Namespace,
	}, robot)
	if err != nil {
		return []reconcile.Request{}
	}

	return []reconcile.Request{
		{
			NamespacedName: types.NamespacedName{
				Name:      robot.Name,
				Namespace: robot.Namespace,
			},
		},
	}
}
