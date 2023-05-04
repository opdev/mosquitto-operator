/*
Copyright 2023 Brad P. Crochet.

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

// MosquittoSpec defines the desired state of Mosquitto
type MosquittoSpec struct {
	Persist bool          `json:"persist,omitempty"`
	Auth    MosquittoAuth `json:"auth,omitempty"`
}

// MosquittoAuth defines the desired state of Auth for Mosquitto
// By default, it is disabled
type MosquittoAuth struct {
	//+kubebuilder:default:=false
	Enabled bool   `json:"enabled"`
	Secret  string `json:"secret,omitempty"`
}

// MosquittoStatus defines the observed state of Mosquitto
type MosquittoStatus struct {
	MosquittoConfConfigMap string `json:"mosquittoConfConfigMap,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Mosquitto is the Schema for the mosquittoes API
type Mosquitto struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MosquittoSpec   `json:"spec,omitempty"`
	Status MosquittoStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MosquittoList contains a list of Mosquitto
type MosquittoList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Mosquitto `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Mosquitto{}, &MosquittoList{})
}
