package a

// It must be ignored since it is not a type
// +kubebuilder:validation:Enum=foo;bar;baz
// +kubebuilder:validation:Enum=foo;bar;baz
var Variable string

// +kubebuilder:validation:Enum=foo;bar;baz
// +kubebuilder:validation:Enum=foo;bar;baz
type Enum string // want "Enum has duplicated markers kubebuilder:validation:Enum"

// +kubebuilder:validation:MaxLength=10
// +kubebuilder:validation:MaxLength=11
type MaxLength int

// +kubebuilder:validation:MaxLength=10
// +kubebuilder:validation:MaxLength=10
type DuplicatedMaxLength int // want "DuplicatedMaxLength has duplicated markers kubebuilder:validation:MaxLength=10"

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:object:root=true
type DuplicateMarkerSpec struct { // want "DuplicateMarkerSpec has duplicated markers kubebuilder:object:root"
	// +kubebuilder:validation:Required
	// should be ignored since it only has single marker
	Required string `json:"required"`

	// +listType=map
	// +listMapKey=primaryKey
	// +listMapKey=secondaryKey
	// +required
	// should be ignored since listMapKey is allowed to have different values
	Map Map `json:"map"`

	// +optional
	// +kubebuilder:validation:XValidation:rule="self >= 1 && self <= 3",message="must be 1 to 5"
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="replicas must be immutable"
	// should be ignored since XValidation is allowed to have different values
	Replicas *int `json:"replicas"`

	// +kubebuilder:validation:MaxLength=11
	// +kubebuilder:validation:MaxLength=10
	// should be ignored since MaxLength is allowed to have different values
	Maxlength int `json:"maxlength"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Required
	DuplicatedRequired string `json:"duplicatedRequired"` // want "DuplicatedRequired has duplicated markers kubebuilder:validation:Required"

	// +kubebuilder:validation:Enum=foo;bar;baz
	// +kubebuilder:validation:Enum=foo;bar;baz
	DuplicatedEnum string `json:"duplicatedEnum"` // want "DuplicatedEnum has duplicated markers kubebuilder:validation:Enum"

	// +kubebuilder:validation:MaxLength=10
	// +kubebuilder:validation:MaxLength=10
	DuplicatedMaxLength int `json:"duplicatedMaxLength"` // want "DuplicatedMaxLength has duplicated markers kubebuilder:validation:MaxLength=10"

	// +kubebuilder:validation:MaxLength=10
	DuplicatedMaxLengthIncludingTypeMarker MaxLength `json:"duplicatedMaxLengthIncludingTypeMarker"` // want "DuplicatedMaxLengthIncludingTypeMarker has duplicated markers kubebuilder:validation:MaxLength=10"

	// +listType=map
	// +listMapKey=primaryKey
	// +listMapKey=secondaryKey
	// +listType=map
	// +required
	DuplicatedListTypeMap Map `json:"duplicatedListTypeMap"` // want "DuplicatedListTypeMap has duplicated markers listType=map"

	// +optional
	// +kubebuilder:validation:XValidation:rule="self >= 1 && self <= 3",message="must be 1 to 5"
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="replicas must be immutable"
	// +kubebuilder:validation:XValidation:rule="self >= 1 && self <= 3",message="must be 1 to 5"
	DuplicatedReplicas *int `json:"duplicatedReplicas"` // want "DuplicatedReplicas has duplicated markers kubebuilder:validation:XValidation:rule=\"self >= 1 && self <= 3\",message=\"must be 1 to 5\""

	// +optional
	// +kubebuilder:validation:XValidation:rule="self >= 1 && self <= 3",message="must be 1 to 5"
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="replicas must be immutable"
	// +kubebuilder:validation:XValidation:message="must be 1 to 5",rule="self >= 1 && self <= 3"
	DuplicatedUnorderedValidationReplicas *int `json:"duplicatedUnorderedValidationReplicas"` // want "DuplicatedUnorderedValidationReplicas has duplicated markers kubebuilder:validation:XValidation:message=\"must be 1 to 5\",rule=\"self >= 1 && self <= 3\""

	StringFromAnotherFile StringFromAnotherFile `json:"stringFromAnotherFile"`

	// +kubebuilder:validation:MaxLength=10
	StringFromAnotherFileWithMaxLength StringFromAnotherFile `json:"stringFromAnotherFileWithMaxLength"` // want "StringFromAnotherFileWithMaxLength has duplicated markers kubebuilder:validation:MaxLength=10"
}

type Map struct {
	// +required
	PrimaryKey string `json:"primaryKey"`
	// +required
	SecondaryKey string `json:"secondaryKey"`
	// +required
	Value string `json:"value"`
}
