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

package v1alpha1

import (
	"github.com/robolaunch/fleet-operator/internal"
	robotv1alpha1 "github.com/robolaunch/robot-operator/pkg/api/roboscale.io/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type FleetSpec struct {
	// Discovery server configuration of fleet. For detailed information, refer the document for the API group `robot.roboscale.io`.
	DiscoveryServerTemplate robotv1alpha1.DiscoveryServerSpec `json:"discoveryServerTemplate,omitempty"`
	// Determines if the fleet should be federated across clusters or not.
	Hybrid bool `json:"hybrid,omitempty"`
	// If `.spec.hybrid` is true, this field includes Kubernetes cluster names which the fleet will be federated to.
	Instances []string `json:"instances,omitempty"`
}

type OwnedNamespaceStatus struct {
	// Generic structure of the most recent status of an owned object. For detailed information, refer the document for the API group `robot.roboscale.io`.
	Resource robotv1alpha1.OwnedResourceStatus `json:"resource,omitempty"`
	// Sets to `true` if the owned namespace is federated.
	Federated bool `json:"federated,omitempty"`
	// Sets to `true` if the namespace is ready for the resources to be deployed such as robot.
	Ready bool `json:"ready,omitempty"`
}

type FleetCompatibilityStatus struct {
	// Indicates the robot's compatibility with fleet.
	IsCompatible bool `json:"isCompatible"`
	// Indicates the possible incompatibility reason of an attached robot.
	Reason string `json:"reason,omitempty"`
}

type AttachedRobot struct {
	// Resource reference for attached robot.
	Reference corev1.ObjectReference `json:"reference,omitempty"`
	// Attached robot phase. For detailed information, refer the document for the API group `robot.roboscale.io`.
	Phase robotv1alpha1.RobotPhase `json:"phase,omitempty"`
	// Compatibility status of attached robot with the fleet.
	FleetCompatibility FleetCompatibilityStatus `json:"fleetCompatibility,omitempty"`
}

type FleetPhase string

const (
	FleetPhaseCheckingRemoteNamespace FleetPhase = "CheckingRemoteNamespace"
	FleetPhaseCreatingNamespace       FleetPhase = "CreatingNamespace"
	FleetPhaseCreatingDiscoveryServer FleetPhase = "CreatingDiscoveryServer"
	FleetPhaseReady                   FleetPhase = "Ready"
)

type FleetStatus struct {
	// Fleet phase.
	Phase FleetPhase `json:"phase,omitempty"`
	// Namespace status. Fleet creates namespace if the `.spec.hybrid` is set to `true`. It creates `FederatedNamespace` if `false`.
	NamespaceStatus OwnedNamespaceStatus `json:"namespaceStatus,omitempty"`
	// Discovery server instance status. For detailed information, refer the document for the API group `robot.roboscale.io`.
	DiscoveryServerStatus robotv1alpha1.OwnedResourceStatus `json:"discoveryServerStatus,omitempty"`
	// Attached launch object information.
	AttachedRobots []AttachedRobot `json:"attachedRobots,omitempty"`
}

//+kubebuilder:resource:scope=Cluster
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Fleet manages lifecycle and configuration of multiple robots and robot's connectivity layer that contains DDS Discovery Server and ROS bridge services.
type Fleet struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// Specification of the desired behavior of the Fleet.
	Spec FleetSpec `json:"spec,omitempty"`
	// Most recently observed status of the Fleet.
	Status FleetStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// FleetList contains a list of Fleet.
type FleetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Fleet `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Fleet{}, &FleetList{})
}

func (fleet *Fleet) GetNamespaceMetadata() *types.NamespacedName {
	return &types.NamespacedName{
		Name: fleet.Name,
	}
}

func (fleet *Fleet) GetDiscoveryServerMetadata() *types.NamespacedName {
	return &types.NamespacedName{
		Name:      fleet.Name + internal.DISCOVERY_SERVER_FLEET_POSTFIX,
		Namespace: fleet.GetNamespaceMetadata().Name,
	}
}
