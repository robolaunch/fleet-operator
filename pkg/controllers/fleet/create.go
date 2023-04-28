package fleet

import (
	"context"

	"github.com/robolaunch/fleet-operator/internal/resources"
	fleetv1alpha1 "github.com/robolaunch/fleet-operator/pkg/api/roboscale.io/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
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

func (r *FleetReconciler) createFederatedNamespace(ctx context.Context, instance *fleetv1alpha1.Fleet, nsNamespacedName *types.NamespacedName) error {

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

	federatedNamespace := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "types.kubefed.io/v1beta1",
			"kind":       "FederatedNamespace",
			"metadata": map[string]interface{}{
				"name":      nsNamespacedName.Name,
				"namespace": nsNamespacedName.Name,
				"ownerReferences": []map[string]interface{}{
					{
						"apiVersion":         instance.APIVersion,
						"blockOwnerDeletion": true,
						"controller":         true,
						"kind":               instance.Kind,
						"name":               instance.Name,
						"uid":                instance.UID,
					},
				},
			},
			"spec": map[string]interface{}{
				"placement": map[string]interface{}{
					"clusters": instancesMapSlice,
				},
			},
		},
	}

	_, err := resourceInterface.Namespace(nsNamespacedName.Name).Create(ctx, &federatedNamespace, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (r *FleetReconciler) createDiscoveryServer(ctx context.Context, instance *fleetv1alpha1.Fleet, dsNamespacedName *types.NamespacedName) error {

	discoveryServer := resources.GetDiscoveryServer(instance, dsNamespacedName)

	err := ctrl.SetControllerReference(instance, discoveryServer, r.Scheme)
	if err != nil {
		return err
	}

	err = r.Create(ctx, discoveryServer)
	if err != nil && errors.IsAlreadyExists(err) {
		return nil
	} else if err != nil {
		return err
	}

	logger.Info("STATUS: Discovery server " + discoveryServer.Name + " is created.")
	return nil
}
