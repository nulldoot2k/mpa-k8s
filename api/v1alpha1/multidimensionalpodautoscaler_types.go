package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// ScaleTargetRef defines the target workload
type ScaleTargetRef struct {
	// +kubebuilder:validation:Required
	APIVersion string `json:"apiVersion"`

	// +kubebuilder:validation:Required
	Kind string `json:"kind"`

	// +kubebuilder:validation:Required
	Name string `json:"name"`
}

// MultidimensionalPodAutoscalerSpec defines desired state
type MultidimensionalPodAutoscalerSpec struct {
	// +kubebuilder:validation:Required
	ScaleTargetRef ScaleTargetRef `json:"scaleTargetRef"`

	// +kubebuilder:default=1
	// +kubebuilder:validation:Minimum=1
	MinReplicas *int32 `json:"minReplicas,omitempty"`

	// +kubebuilder:default=3
	// +kubebuilder:validation:Minimum=1
	MaxReplicas *int32 `json:"maxReplicas,omitempty"`
}

// MultidimensionalPodAutoscalerStatus defines observed state
type MultidimensionalPodAutoscalerStatus struct {
	LastAction    string       `json:"lastAction,omitempty"`
	LastScaleTime *metav1.Time `json:"lastScaleTime,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type MultidimensionalPodAutoscaler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MultidimensionalPodAutoscalerSpec   `json:"spec,omitempty"`
	Status MultidimensionalPodAutoscalerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type MultidimensionalPodAutoscalerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MultidimensionalPodAutoscaler `json:"items"`
}

func init() {
	SchemeBuilder.Register(
		&MultidimensionalPodAutoscaler{},
		&MultidimensionalPodAutoscalerList{},
	)
}
