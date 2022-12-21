package fleet

import (
	"context"
	"sort"

	"github.com/robolaunch/fleet-operator/internal"
	fleetv1alpha1 "github.com/robolaunch/fleet-operator/pkg/api/roboscale.io/v1alpha1"
	robotv1alpha1 "github.com/robolaunch/robot-operator/pkg/api/roboscale.io/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *FleetReconciler) reconcileHandleAttachments(ctx context.Context, instance *fleetv1alpha1.Fleet) error {

	// select attached robots
	err := r.reconcileAttachRobots(ctx, instance)
	if err != nil {
		return err
	}

	return nil
}

func (r *FleetReconciler) reconcileAttachRobots(ctx context.Context, instance *fleetv1alpha1.Fleet) error {

	// Get attached robots for this robot
	requirements := []labels.Requirement{}
	newReq, err := labels.NewRequirement(internal.FLEET_LABEL_KEY, selection.In, []string{instance.Name})
	if err != nil {
		return err
	}
	requirements = append(requirements, *newReq)

	fleetSelector := labels.NewSelector().Add(requirements...)

	robotList := robotv1alpha1.RobotList{}
	err = r.List(ctx, &robotList, &client.ListOptions{Namespace: instance.Namespace, LabelSelector: fleetSelector})
	if err != nil {
		return err
	}

	if len(robotList.Items) == 0 {
		instance.Status.AttachedRobots = []fleetv1alpha1.AttachedRobot{}
		return nil
	}

	// Sort attached robots for this robot according to their creation timestamps
	sort.SliceStable(robotList.Items[:], func(i, j int) bool {
		return robotList.Items[i].CreationTimestamp.String() < robotList.Items[j].CreationTimestamp.String()
	})

	instance.Status.AttachedRobots = []fleetv1alpha1.AttachedRobot{}

	for _, robot := range robotList.Items {

		fleetCompatibility := fleetv1alpha1.FleetCompatibilityStatus{}
		fleetCompatibility.IsCompatible = true

		err := checkRobotDiscovery(*instance, robot)
		if err != nil {
			fleetCompatibility.IsCompatible = false
			fleetCompatibility.Reason = err.Error()
		}

		err = checkTenancy(*instance, robot)
		if err != nil {
			fleetCompatibility.IsCompatible = false
			fleetCompatibility.Reason = err.Error()
		}

		instance.Status.AttachedRobots = append(instance.Status.AttachedRobots, fleetv1alpha1.AttachedRobot{
			Reference: corev1.ObjectReference{
				Kind:            robot.Kind,
				Namespace:       robot.Namespace,
				Name:            robot.Name,
				UID:             robot.UID,
				APIVersion:      robot.APIVersion,
				ResourceVersion: robot.ResourceVersion,
			},
			Phase:              robot.Status.Phase,
			FleetCompatibility: fleetCompatibility,
		})

	}

	return nil
}
