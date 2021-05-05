package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	// Rules the rules to apply
	Rules []Rule `json:"rules,omitempty"`
}

// Rule the rules to apply
type Rule struct {
	// Resources the kubernetes resource dir to look for resources to verify
	Resources *Source `json:"resources,omitempty"`

	// Charts the charts to evaluate
	Charts *Charts `json:"charts,omitempty"`

	// Tests the tests to perform
	Tests Tests `json:"tests,omitempty"`
}

// Source the location of kubernetes resources to validate
type Source struct {
	// Dir the directory containing the kubernetes resources
	Dir string `json:"dir,omitempty"`
}

// Charts the charts to template and validate
type Charts struct {
	// Dir the directory containing a helm chart or the source to recurse through if recursive is enabled
	Dir string `json:"dir,omitempty"`

	// Recurse if enabled recurse through the directory to find any Chart.yaml files
	Recurse bool `json:"recurse,omitempty"`
}

// Test a kind of test
type Test struct {
	// Version optional override of the version to use
	Version string `json:"version,omitempty"`
	// Args optional additional comand line arguments to pass to the test
	Args []string `json:"args,omitempty"`
}

// Tests the tests to run on the resources
type Tests struct {
	// Conftest enables conftest tests
	Conftest *Test `json:"conftest,omitempty"`

	// Kubescore enables kube-score based tests
	Kubescore *Test `json:"kubescore,omitempty"`

	// Kubeval enables kubeval tests
	Kubeval *Test `json:"kubeval,omitempty"`

	// Polaris enables polaris tests
	Polaris *Test `json:"polaris,omitempty"`
}
