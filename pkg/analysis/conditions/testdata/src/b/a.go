package b

import (
	"go/ast"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ValidConditions struct {
	// conditions is an accurate representation of the desired state of a conditions object.
	// +listType=map
	// +listMapKey=type
	// +patchStrategy=merge
	// +patchMergeKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`

	// other fields
	OtherField string `json:"otherField,omitempty"`
}

type ConditionsNotFirst struct {
	// other fields
	OtherField string `json:"otherField,omitempty"`

	// conditions is an accurate representation of the desired state of a conditions object.
	// +listType=map
	// +listMapKey=type
	// +patchStrategy=merge
	// +patchMergeKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,2,rep,name=conditions"`
}

type ConditionsThird struct {
	// other fields
	OtherField string `json:"otherField,omitempty"`

	// another field
	AnotherField string `json:"anotherField,omitempty"`

	// conditions is an accurate representation of the desired state of a conditions object.
	// +listType=map
	// +listMapKey=type
	// +patchStrategy=merge
	// +patchMergeKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,3,rep,name=conditions"`
}

type ConditionsIncorrectType struct {
	// conditions has an incorrect type.
	Conditions map[string]metav1.Condition // want "Conditions field must be a slice of metav1.Condition"
}

type ConditionsIncorrectSliceElement struct {
	// conditions has an incorrect type.
	Conditions []string // want "Conditions field must be a slice of metav1.Condition"
}

type ConditionsIncorrectImportedSliceElement struct {
	// conditions has an incorrect type.
	Conditions []metav1.Time // want "Conditions field must be a slice of metav1.Condition"
}

type ConditionsIncorrectImportedPackage struct {
	// conditions has an incorrect type.
	Conditions []ast.Node // want "Conditions field must be a slice of metav1.Condition"
}

type MissingAllMarkers struct {
	// conditions is missing all markers.
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"` // want "Conditions field is missing the following markers: listType=map, listMapKey=type, patchStrategy=merge, patchMergeKey=type, optional"
}

type MissingListMarkers struct {
	// conditions is missing list markers.
	// +patchStrategy=merge
	// +patchMergeKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"` // want "Conditions field is missing the following markers: listType=map, listMapKey=type"
}

type MissingPatchMarkers struct {
	// conditions is missing patch markers.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"` // want "Conditions field is missing the following markers: patchStrategy=merge, patchMergeKey=type"
}

type MissingOptionalMarker struct {
	// conditions is missng the optional marker.
	// +listType=map
	// +listMapKey=type
	// +patchStrategy=merge
	// +patchMergeKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"` // want "Conditions field is missing the following markers: optional"
}

type MissingFieldTag struct {
	// conditions is missing the field tag.
	// +listType=map
	// +listMapKey=type
	// +patchStrategy=merge
	// +patchMergeKey=type
	// +optional
	Conditions []metav1.Condition // want "Conditions field is missing tags, should be: `json:\"conditions,omitempty\" patchStrategy:\"merge\" patchMergeKey:\"type\" protobuf:\"bytes,1,rep,name=conditions\"`"
}

type IncorrectFieldTag struct {
	// conditions has an incorrect field tag.
	// +listType=map
	// +listMapKey=type
	// +patchStrategy=merge
	// +patchMergeKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions"  patchMergeKey:"type" protobuf:"bytes,3,rep,name=conditions"` // want "Conditions field has incorrect tags, should be: `json:\"conditions,omitempty\" patchStrategy:\"merge\" patchMergeKey:\"type\" protobuf:\"bytes,1,rep,name=conditions\"`"
}
