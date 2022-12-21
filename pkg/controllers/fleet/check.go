package fleet

import (
	"context"
	goErr "errors"
	"reflect"

	"github.com/robolaunch/fleet-operator/internal/label"
	fleetv1alpha1 "github.com/robolaunch/fleet-operator/pkg/api/roboscale.io/v1alpha1"
	robotv1alpha1 "github.com/robolaunch/robot-operator/pkg/api/roboscale.io/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
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

func (r *FleetReconciler) reconcileCheckAttachedRobots(ctx context.Context, instance *fleetv1alpha1.Fleet) error {

	for k, obj := range instance.Status.AttachedRobots {
		robot := &robotv1alpha1.Robot{}
		err := r.Get(ctx, types.NamespacedName{Namespace: obj.Reference.Namespace, Name: obj.Reference.Name}, robot)
		if err != nil && errors.IsNotFound(err) {
			// TODO: Empty the reference fields
			return err
		} else if err != nil {
			return err
		} else {

			obj.FleetCompatibility.IsCompatible = true
			obj.FleetCompatibility.Reason = ""
			obj.Phase = robot.Status.Phase

			err := checkRobotDiscovery(*instance, *robot)
			if err != nil {
				obj.FleetCompatibility.IsCompatible = false
				obj.FleetCompatibility.Reason = err.Error()
			}

			err = checkTenancy(*instance, *robot)
			if err != nil {
				obj.FleetCompatibility.IsCompatible = false
				obj.FleetCompatibility.Reason = err.Error()
			}

		}

		instance.Status.AttachedRobots[k] = obj
	}

	return nil
}

func checkRobotDiscovery(fleet fleetv1alpha1.Fleet, robot robotv1alpha1.Robot) error {

	fleetDsConfig := fleet.Spec.DiscoveryServerTemplate
	robotDsConfig := robot.Spec.DiscoveryServerTemplate

	if robotDsConfig.Type == robotv1alpha1.DiscoveryServerInstanceTypeServer {
		return goErr.New("discovery server configuration is not compatible with fleet, wrong type")
	}

	if fleetDsConfig.Type == robotv1alpha1.DiscoveryServerInstanceTypeServer {

		if (fleetDsConfig.Cluster != robotDsConfig.Cluster) ||
			(fleetDsConfig.Hostname != robotDsConfig.Hostname) ||
			(fleetDsConfig.Subdomain != robotDsConfig.Subdomain) {
			return goErr.New("discovery server configuration is not compatible with fleet")
		}

		if (robotDsConfig.Reference.Name != fleet.GetDiscoveryServerMetadata().Name) ||
			(robotDsConfig.Reference.Namespace != fleet.GetDiscoveryServerMetadata().Namespace) {
			return goErr.New("discovery server configuration is not compatible with fleet, wrong reference")
		}

	} else if fleetDsConfig.Type == robotv1alpha1.DiscoveryServerInstanceTypeClient {

		if !reflect.DeepEqual(fleetDsConfig, robotDsConfig) {
			return goErr.New("discovery server configuration is not compatible with fleet")
		}

	} else {
		return goErr.New("discovery server configuration is not compatible with fleet")
	}

	return nil

}

func checkTenancy(fleet fleetv1alpha1.Fleet, robot robotv1alpha1.Robot) error {

	fleetTenancy := label.GetTenancy(&fleet)
	robotTenancy := label.GetTenancy(&robot)
	if !reflect.DeepEqual(fleetTenancy, robotTenancy) {
		return goErr.New("tenancy configuration is not compatible with fleet")
	}

	return nil
}
