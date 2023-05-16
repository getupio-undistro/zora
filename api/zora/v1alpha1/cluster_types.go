// Copyright 2022 Undistro Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/undistro/zora/pkg/discovery"
	"github.com/undistro/zora/pkg/formats"
)

const (
	LabelEnvironment = "zora.undistro.io/environment"

	ClusterReady               = "Ready"
	ClusterDiscovered          = "Discovered"
	ClusterResourcesDiscovered = "ResourcesDiscovered"
)

// ClusterSpec defines the desired state of Cluster
type ClusterSpec struct {
	// KubeconfigRef is a reference to a secret in the same namespace that contains the kubeconfig data
	KubeconfigRef *corev1.LocalObjectReference `json:"kubeconfigRef,omitempty"`
}

// ClusterStatus defines the observed state of Cluster
type ClusterStatus struct {
	Status                `json:",inline"`
	discovery.ClusterInfo `json:",inline"`

	// KubernetesVersion is the server's kubernetes version (git version).
	KubernetesVersion string `json:"kubernetesVersion,omitempty"`

	// Usage and available resources
	Resources discovery.ClusterResources `json:"resources,omitempty"`

	// Usage of memory in quantity and percentage
	MemoryUsage string `json:"memoryUsage,omitempty"`

	// Quantity of memory available in Mi
	MemoryAvailable string `json:"memoryAvailable,omitempty"`

	// Usage of CPU in quantity and percentage
	CPUUsage string `json:"cpuUsage,omitempty"`

	// Quantity of CPU available
	CPUAvailable string `json:"cpuAvailable,omitempty"`

	// Timestamp representing the server time of the last reconciliation
	LastReconciliationTime metav1.Time `json:"lastReconciliationTime,omitempty"`
}

// SetResources format and fill temporary fields about resources
func (in *ClusterStatus) SetResources(res discovery.ClusterResources) {
	in.Resources = res.DeepCopy()
	if m, found := res[corev1.ResourceMemory]; found {
		in.MemoryAvailable = formats.Memory(m.Available)
		in.MemoryUsage = formats.MemoryUsage(m.Usage, m.UsagePercentage)
	}
	if c, found := res[corev1.ResourceCPU]; found {
		in.CPUAvailable = formats.CPU(c.Available)
		in.CPUUsage = formats.CPUUsage(c.Usage, c.UsagePercentage)
	}
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Version",type="string",priority=0,JSONPath=".status.kubernetesVersion"
//+kubebuilder:printcolumn:name="MEM Available",type="string",priority=0,JSONPath=".status.memoryAvailable"
//+kubebuilder:printcolumn:name="MEM Usage (%)",type="string",priority=0,JSONPath=".status.memoryUsage"
//+kubebuilder:printcolumn:name="CPU Available",type="string",priority=0,JSONPath=".status.cpuAvailable"
//+kubebuilder:printcolumn:name="CPU Usage (%)",type="string",priority=0,JSONPath=".status.cpuUsage"
//+kubebuilder:printcolumn:name="Nodes",type="integer",priority=0,JSONPath=".status.totalNodes"
//+kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type==\"Ready\")].status"
//+kubebuilder:printcolumn:name="Age",type="date",priority=0,JSONPath=".status.creationTimestamp"
//+kubebuilder:printcolumn:name="Provider",type="string",priority=1,JSONPath=".status.provider"
//+kubebuilder:printcolumn:name="Region",type="string",priority=1,JSONPath=".status.region"

// Cluster is the Schema for the clusters API
// +genclient
// +genclient:onlyVerbs=list,get
// +genclient:noStatus
type Cluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterSpec   `json:"spec,omitempty"`
	Status ClusterStatus `json:"status,omitempty"`
}

func (in *Cluster) KubeconfigRefKey() *types.NamespacedName {
	if in.Spec.KubeconfigRef == nil {
		return nil
	}
	return &types.NamespacedName{Name: in.Spec.KubeconfigRef.Name, Namespace: in.Namespace}
}

func (in *Cluster) SetStatus(statusType string, status bool, reason, msg string) {
	s := metav1.ConditionFalse
	if status {
		s = metav1.ConditionTrue
	}
	in.Status.SetCondition(metav1.Condition{
		Type:               statusType,
		Status:             s,
		ObservedGeneration: in.Generation,
		Reason:             reason,
		Message:            msg,
	})
}

//+kubebuilder:object:root=true

// ClusterList contains a list of Cluster
type ClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Cluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Cluster{}, &ClusterList{})
}
