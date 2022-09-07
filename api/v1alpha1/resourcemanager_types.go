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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ResourceManagerSpec defines the desired state of ResourceManager
type ResourceManagerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Disabled bool `json:"active,omitempty"`
	DryRun   bool `json:"dry-run,omitempty"`

	// TODO: add validation
	ResourceKind string                `json:"resourceKind"`
	Selector     *metav1.LabelSelector `json:"selector"`

	// TODO: add validation + enum
	Action string `json:"action"`

	Condition []ExpiryCondition `json:"condition"`
}

type Condition struct {
	Type string `json:"type"`
}

type ExpiryCondition struct {
	Condition `json:",inline"`
	After     string `json:"after"`
}

// type IntervalCondition struct {
// 	Condition
// }

// ResourceManagerStatus defines the observed state of ResourceManager
type ResourceManagerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ResourceManager is the Schema for the resourcemanagers API
type ResourceManager struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ResourceManagerSpec   `json:"spec,omitempty"`
	Status ResourceManagerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ResourceManagerList contains a list of ResourceManager
type ResourceManagerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ResourceManager `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ResourceManager{}, &ResourceManagerList{})
}
