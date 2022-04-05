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
	api "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DemoDeploymentSpec defines the desired state of DemoDeployment
type DemoDeploymentSpec struct {
	// +kubebuilder:validation:Maximum=2
	Replicas int32 `json:"replicas,omitempty"`

	Selector *metav1.LabelSelector `json:"selector,omitempty"`

	Template api.PodTemplateSpec `json:"template,omitempty"`
}

// DemoDeploymentStatus defines the observed state of DemoDeployment
type DemoDeploymentStatus struct {
	Replicas int32 `json:"replicas,omitempty"`

	UpdatedReplicas int32 `json:"updatedReplicas,omitempty"`

	ReadyReplicas int32 `json:"readyReplicas,omitempty"`

	AvailableReplicas int32 `json:"availablereplicas,omitempty"`

	UnavailableReplicas int32 `json:"unavailableReplicas,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// DemoDeployment is the Schema for the demodeployments API
type DemoDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DemoDeploymentSpec   `json:"spec,omitempty"`
	Status DemoDeploymentStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DemoDeploymentList contains a list of DemoDeployment
type DemoDeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DemoDeployment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DemoDeployment{}, &DemoDeploymentList{})
}
