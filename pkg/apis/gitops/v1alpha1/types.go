package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// KubeTestFileName default name of the configuration file
	KubeTestFileName = "kube-test.yaml"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KubeTest represents the configuration of kube test
//
// +k8s:openapi-gen=true
type KubeTest struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata"`

	// Spec holds the desired state of the KubeTest from the client
	// +optional
	Spec KubeTestSpec `json:"spec"`
}

// KubeTestList contains a list of KubeTest
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type KubeTestList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KubeTest `json:"items"`
}

// KubeTestSpec defines the configuration of kube test
type KubeTestSpec struct {
}
