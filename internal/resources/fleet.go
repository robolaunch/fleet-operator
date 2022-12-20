package resources

import (
	fleetv1alpha1 "github.com/robolaunch/fleet-operator/pkg/api/roboscale.io/v1alpha1"
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
