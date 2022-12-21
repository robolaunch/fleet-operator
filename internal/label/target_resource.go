package label

import (
	"github.com/robolaunch/fleet-operator/internal"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetTargetFleet(obj metav1.Object) string {
	labels := obj.GetLabels()

	if targetFleet, ok := labels[internal.FLEET_LABEL_KEY]; ok {
		return targetFleet
	}

	return ""
}
