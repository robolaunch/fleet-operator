package resources

import (
	fleetv1alpha1 "github.com/robolaunch/fleet-operator/pkg/api/roboscale.io/v1alpha1"
	robotv1alpha1 "github.com/robolaunch/robot-operator/pkg/api/roboscale.io/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func GetNamespace(fleet *fleetv1alpha1.Fleet, nsNamespacedName *types.NamespacedName) *corev1.Namespace {

	namespace := corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: nsNamespacedName.Name,
		},
	}

	return &namespace
}

func GetDiscoveryServer(fleet *fleetv1alpha1.Fleet, dsNamespacedName *types.NamespacedName) *robotv1alpha1.DiscoveryServer {

	discoveryServer := robotv1alpha1.DiscoveryServer{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dsNamespacedName.Name,
			Namespace: dsNamespacedName.Namespace,
			Labels:    fleet.Labels,
		},
		Spec: fleet.Spec.DiscoveryServerTemplate,
	}

	return &discoveryServer
}
