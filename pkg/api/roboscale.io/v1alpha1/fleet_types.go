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

// FleetSpec defines the desired state of Fleet
type FleetSpec struct {
	DiscoveryServerTemplate robotv1alpha1.DiscoveryServerSpec `json:"discoveryServerTemplate,omitempty"`
	Hybrid                  bool                              `json:"hybrid,omitempty"`
}

type NamespaceStatus struct {
	Created bool `json:"created,omitempty"`
	Ready   bool `json:"ready,omitempty"`
}

type DiscoveryServerInstanceStatus struct {
	Created bool                                `json:"created,omitempty"`
	Status  robotv1alpha1.DiscoveryServerStatus `json:"status,omitempty"`
}

type FleetPhase string

const (
	FleetPhaseCreatingNamespace       FleetPhase = "CreatingNamespace"
	FleetPhaseCreatingDiscoveryServer FleetPhase = "CreatingDiscoveryServer"
	FleetPhaseReady                   FleetPhase = "Ready"
)

type FleetCompatibilityStatus struct {
	IsCompatible bool   `json:"isCompatible"`
	Reason       string `json:"reason,omitempty"`
}

type AttachedRobot struct {
	Reference          corev1.ObjectReference   `json:"reference,omitempty"`
	Phase              robotv1alpha1.RobotPhase `json:"phase,omitempty"`
	FleetCompatibility FleetCompatibilityStatus `json:"fleetCompatibility,omitempty"`
}

// FleetStatus defines the observed state of Fleet
type FleetStatus struct {
	Phase                 FleetPhase                    `json:"phase,omitempty"`
	NamespaceStatus       NamespaceStatus               `json:"namespaceStatus,omitempty"`
	DiscoveryServerStatus DiscoveryServerInstanceStatus `json:"discoveryServerStatus,omitempty"`
	// Attached launch object information
	AttachedRobots []AttachedRobot `json:"attachedRobots,omitempty"`
}

//+kubebuilder:resource:scope=Cluster
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Fleet is the Schema for the fleets API
type Fleet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FleetSpec   `json:"spec,omitempty"`
	Status FleetStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// FleetList contains a list of Fleet
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
