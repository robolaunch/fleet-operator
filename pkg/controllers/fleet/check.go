package fleet

import (
	"context"
	goErr "errors"
	"reflect"

	fleetErr "github.com/robolaunch/fleet-operator/internal/error"
	"github.com/robolaunch/fleet-operator/internal/label"
	"github.com/robolaunch/fleet-operator/internal/reference"
	fleetv1alpha1 "github.com/robolaunch/fleet-operator/pkg/api/roboscale.io/v1alpha1"
	robotv1alpha1 "github.com/robolaunch/robot-operator/pkg/api/roboscale.io/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

func (r *FleetReconciler) reconcileCheckNamespace(ctx context.Context, instance *fleetv1alpha1.Fleet) error {

	switch instance.Spec.Hybrid {
	case true:

		switch label.GetInstanceType(instance) {
		case label.InstanceTypeCloudInstance:

			// check ns
			namespaceQuery := &corev1.Namespace{}
			err := r.Get(ctx, *instance.GetNamespaceMetadata(), namespaceQuery)
			if err != nil && errors.IsNotFound(err) {
				instance.Status.NamespaceStatus = fleetv1alpha1.OwnedNamespaceStatus{}
			} else if err != nil {
				return err
			} else {
				instance.Status.NamespaceStatus.Resource.Created = true

				// check federated ns
				resourceInterface := r.DynamicClient.Resource(schema.GroupVersionResource{
					Group:    "types.kubefed.io",
					Version:  "v1beta1",
					Resource: "federatednamespaces",
				})

				instance.Status.NamespaceStatus.Federated = true
				instance.Status.NamespaceStatus.Ready = true

				unstructuredFedNs, err := resourceInterface.Namespace(instance.GetNamespaceMetadata().Name).Get(ctx, instance.GetNamespaceMetadata().Name, metav1.GetOptions{})
				if err != nil {
					instance.Status.NamespaceStatus.Federated = false
					instance.Status.NamespaceStatus.Ready = false
					return nil
				}

				desiredInstancesMapSlice := []map[string]interface{}{}
				for _, i := range instance.Spec.Instances {
					desiredInstancesMapSlice = append(desiredInstancesMapSlice, map[string]interface{}{
						"name": i,
					})
				}

				actualInstancesMapSlice := []map[string]interface{}{}
				actualInstancesInterface := unstructuredFedNs.Object["spec"].(map[string]interface{})["placement"].(map[string]interface{})["clusters"].([]interface{})
				for _, v := range actualInstancesInterface {
					actualInstancesMapSlice = append(actualInstancesMapSlice, v.(map[string]interface{}))
				}

				if !reflect.DeepEqual(desiredInstancesMapSlice, actualInstancesMapSlice) {
					err := r.updateFederatedNamespace(ctx, instance, instance.GetNamespaceMetadata(), unstructuredFedNs)
					if err != nil {
						logger.Info("UPDATE: Updating FederatedNamespace.")
						return err
					}
				}

			}

		}

	case false:

		namespaceQuery := &corev1.Namespace{}
		err := r.Get(ctx, *instance.GetNamespaceMetadata(), namespaceQuery)
		if err != nil && errors.IsNotFound(err) {
			instance.Status.NamespaceStatus = fleetv1alpha1.OwnedNamespaceStatus{}
		} else if err != nil {
			return err
		} else {
			instance.Status.NamespaceStatus.Resource.Created = true
			instance.Status.NamespaceStatus.Ready = true
		}

	}

	return nil
}

func (r *FleetReconciler) updateFederatedNamespace(ctx context.Context, instance *fleetv1alpha1.Fleet, nsNamespacedName *types.NamespacedName, federatedNs *unstructured.Unstructured) error {

	resourceInterface := r.DynamicClient.Resource(schema.GroupVersionResource{
		Group:    "types.kubefed.io",
		Version:  "v1beta1",
		Resource: "federatednamespaces",
	})

	instancesMapSlice := []map[string]interface{}{}
	for _, i := range instance.Spec.Instances {
		instancesMapSlice = append(instancesMapSlice, map[string]interface{}{
			"name": i,
		})
	}

	desiredFederatedNamespace := unstructured.Unstructured{
		Object: map[string]interface{}{
			"spec": map[string]interface{}{
				"placement": map[string]interface{}{
					"clusters": instancesMapSlice,
				},
			},
		},
	}

	federatedNs.Object["spec"] = desiredFederatedNamespace.Object["spec"]

	_, err := resourceInterface.Namespace(nsNamespacedName.Name).Update(ctx, federatedNs, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (r *FleetReconciler) reconcileCheckRemoteNamespace(ctx context.Context, instance *fleetv1alpha1.Fleet) error {

	namespaceQuery := &corev1.Namespace{}
	err := r.Get(ctx, *instance.GetNamespaceMetadata(), namespaceQuery)
	if err != nil && errors.IsNotFound(err) {
		instance.Status.NamespaceStatus = fleetv1alpha1.OwnedNamespaceStatus{}
		instance.Status.Phase = fleetv1alpha1.FleetPhaseCheckingRemoteNamespace

		err := r.reconcileUpdateInstanceStatus(ctx, instance)
		if err != nil {
			return err
		}

		return &fleetErr.NamespaceNotFoundError{
			ResourceKind:      instance.Kind,
			ResourceName:      instance.Name,
			ResourceNamespace: instance.Namespace,
		}
	} else if err != nil {
		return err
	} else {
		instance.Status.NamespaceStatus.Ready = true
		reference.SetReference(&instance.Status.NamespaceStatus.Resource.Reference, namespaceQuery.TypeMeta, namespaceQuery.ObjectMeta)
	}

	return nil
}

func (r *FleetReconciler) reconcileCheckDiscoveryServer(ctx context.Context, instance *fleetv1alpha1.Fleet) error {

	discoveryServerQuery := &robotv1alpha1.DiscoveryServer{}
	err := r.Get(ctx, *instance.GetDiscoveryServerMetadata(), discoveryServerQuery)
	if err != nil && errors.IsNotFound(err) {
		instance.Status.DiscoveryServerStatus = robotv1alpha1.OwnedResourceStatus{}
	} else if err != nil {
		return err
	} else {
		instance.Status.DiscoveryServerStatus.Created = true
		instance.Status.DiscoveryServerStatus.Phase = string(discoveryServerQuery.Status.Phase)
		reference.SetReference(&instance.Status.DiscoveryServerStatus.Reference, discoveryServerQuery.TypeMeta, discoveryServerQuery.ObjectMeta)
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
