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
	robotv1alpha1 "github.com/robolaunch/robot-operator/pkg/api/roboscale.io/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// FleetSpec defines the desired state of Fleet
type FleetSpec struct {
	DiscoveryServerTemplate robotv1alpha1.DiscoveryServerSpec `json:"discoveryServerTemplate,omitempty"`
}

type NamespaceStatus struct {
	Created bool `json:"created,omitempty"`
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

// FleetStatus defines the observed state of Fleet
type FleetStatus struct {
	Phase                 FleetPhase                    `json:"phase,omitempty"`
	NamespaceStatus       NamespaceStatus               `json:"namespaceStatus,omitempty"`
	DiscoveryServerStatus DiscoveryServerInstanceStatus `json:"discoveryServerStatus,omitempty"`
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
