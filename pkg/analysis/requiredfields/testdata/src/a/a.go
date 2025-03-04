package a

type A struct {
	// optional field should not be picked up.
	// +optional
	OptionalField *string `json:"optionalField,omitempty"`

	// requiredCorrectField should not be picked up.
	// +required
	RequiredCorrectField string `json:"requiredCorrectField"`

	// requiredOmitEmptyField field should be picked up.
	// +required
	RequiredOmitEmptyField string `json:"requiredOmitEmptyField,omitempty"` // want "field RequiredOmitEmptyField is marked as required, but has the omitempty tag"

	// requiredPointerField should be picked up.
	// +required
	RequiredPointerField *string `json:"requiredPointerField"` // want "field RequiredPointerField is marked as required, should not be a pointer"

	// requiredPointerOmitEmptyField should be picked up.
	// +required
	RequiredPointerOmitEmptyField *string `json:"requiredPointerOmitEmptyField,omitempty"` // want "field RequiredPointerOmitEmptyField is marked as required, but has the omitempty tag" "field RequiredPointerOmitEmptyField is marked as required, should not be a pointer"

	// requiredKubebuilderMarkerField should not be picked up.
	// +kubebuilder:validation:Required
	RequiredKubebuilderMarkerField string `json:"requiredKubebuilderMarkerField"`

	// requiredKubebuilderMarkerOmitEmptyField should be picked up.
	// +kubebuilder:validation:Required
	RequiredKubebuilderMarkerOmitEmptyField string `json:"requiredKubebuilderMarkerOmitEmptyField,omitempty"` // want "field RequiredKubebuilderMarkerOmitEmptyField is marked as required, but has the omitempty tag"

	// requiredKubebuilderMarkerPointerField should be picked up.
	// +kubebuilder:validation:Required
	RequiredKubebuilderMarkerPointerField *string `json:"requiredKubebuilderMarkerPointerField"` // want "field RequiredKubebuilderMarkerPointerField is marked as required, should not be a pointer"

	// requiredKubebuilderMarkerPointerOmitEmptyField should be picked up.
	// +kubebuilder:validation:Required
	RequiredKubebuilderMarkerPointerOmitEmptyField *string `json:"requiredKubebuilderMarkerPointerOmitEmptyField,omitempty"` // want "field RequiredKubebuilderMarkerPointerOmitEmptyField is marked as required, but has the omitempty tag" "field RequiredKubebuilderMarkerPointerOmitEmptyField is marked as required, should not be a pointer"
}

// DoNothing is used to check that the analyser doesn't report on methods.
func (A) DoNothing() {}

type Interface interface {
	InaccessibleFunction() string
}
