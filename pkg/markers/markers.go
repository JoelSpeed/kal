/*
Copyright 2025 The Kubernetes Authors.

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
package markers

const (
	// OptionalMarker is the marker that indicates that a field is optional.
	OptionalMarker = "optional"

	// RequiredMarker is the marker that indicates that a field is required.
	RequiredMarker = "required"
)

const (
	// KubebuilderRootMarker is the marker that indicates that a struct is the object root for code and CRD generation.
	KubebuilderRootMarker = "kubebuilder:object:root"

	// KubebuilderStatusSubresourceMarker is the marker that indicates that the CRD generated for a struct should include the /status subresource.
	KubebuilderStatusSubresourceMarker = "kubebuilder:subresource:status"

	// KubebuilderEnumMarker is the marker that indicates that a field has an enum in kubebuilder.
	KubebuilderEnumMarker = "kubebuilder:validation:Enum"

	// KubebuilderFormatMarker is the marker that indicates that a field has a format in kubebuilder.
	KubebuilderFormatMarker = "kubebuilder:validation:Format"

	// KubebuilderMaximumMarker is the marker that indicates that a field has a maximum value in kubebuilder.
	KubebuilderMaximumMarker = "kubebuilder:validation:Maximum"

	// KubebuilderMaxItemsMarker is the marker that indicates that a field has a maximum number of items in kubebuilder.
	KubebuilderMaxItemsMarker = "kubebuilder:validation:MaxItems"

	// KubebuilderMaxLengthMarker is the marker that indicates that a field has a maximum length in kubebuilder.
	KubebuilderMaxLengthMarker = "kubebuilder:validation:MaxLength"

	// KubebuilderMinimumMarker is the marker that indicates that a field has a minimum value in kubebuilder.
	KubebuilderMinimumMarker = "kubebuilder:validation:Minimum"

	// KubebuilderMinItemsMarker is the marker that indicates that a field has a minimum number of items in kubebuilder.
	KubebuilderMinItemsMarker = "kubebuilder:validation:MinItems"

	// KubebuilderMinLengthMarker is the marker that indicates that a field has a minimum length in kubebuilder.
	KubebuilderMinLengthMarker = "kubebuilder:validation:MinLength"

	// KubebuilderMinPropertiesMarker is the marker that indicates that a field has a minimum number of properties in kubebuilder.
	KubebuilderMinPropertiesMarker = "kubebuilder:validation:MinProperties"

	// KubebuilderOptionalMarker is the marker that indicates that a field is optional in kubebuilder.
	KubebuilderOptionalMarker = "kubebuilder:validation:Optional"

	// KubebuilderRequiredMarker is the marker that indicates that a field is required in kubebuilder.
	KubebuilderRequiredMarker = "kubebuilder:validation:Required"

	// KubebuilderItemsMaxLengthMarker is the marker that indicates that a field has a maximum length in kubebuilder.
	KubebuilderItemsMaxLengthMarker = "kubebuilder:validation:items:MaxLength"

	// KubebuilderItemsEnumMarker is the marker that indicates that a field has an enum in kubebuilder.
	KubebuilderItemsEnumMarker = "kubebuilder:validation:items:Enum"

	// KubebuilderItemsFormatMarker is the marker that indicates that a field has a format in kubebuilder.
	KubebuilderItemsFormatMarker = "kubebuilder:validation:items:Format"
)

const (
	// K8sOptionalMarker is the marker that indicates that a field is optional in k8s declarative validation.
	K8sOptionalMarker = "k8s:optional"

	// K8sRequiredMarker is the marker that indicates that a field is required in k8s declarative validation.
	K8sRequiredMarker = "k8s:required"
)
