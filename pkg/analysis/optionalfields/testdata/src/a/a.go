package a

type A struct {
	// required field should not be picked up.
	// +required
	RequiredField string `json:"requiredField,omitempty"`

	// pointerString is a pointer string field.
	// +optional
	PointerString *string `json:"pointerString,omitempty"`

	// pointerInt is a pointer int field.
	// +optional
	PointerInt *int `json:"pointerInt,omitempty"`

	// pointerStruct is a pointer struct field.
	// +optional
	PointerStruct *B `json:"pointerStruct,omitempty"`

	// string is a string field.
	// +optional
	String string `json:"string,omitempty"` // want "field String is optional and should be a pointer"

	// NonOmittedString is a string field without omitempty
	// +optional
	NonOmittedString string `json:"nonOmittedString"` // want "field NonOmittedString is optional and should be a pointer" "field NonOmittedString is optional and should be omitempty"

	// int is an int field.
	// +optional
	Int int `json:"int,omitempty"` // want "field Int is optional and should be a pointer"

	// nonOmittedInt is an int field without omitempty
	// +optional
	NonOmittedInt int `json:"nonOmittedInt"` // want "field NonOmittedInt is optional and should be a pointer" "field NonOmittedInt is optional and should be omitempty"

	// struct is a struct field.
	// +optional
	Struct B `json:"struct,omitempty"` // want "field Struct is optional and should be a pointer"

	// nonOmittedStruct is a struct field without omitempty.
	// +optional
	NonOmittedStruct B `json:"nonOmittedStruct"` // want "field NonOmittedStruct is optional and should be a pointer" "field NonOmittedStruct is optional and should be omitempty"

	// structWithMinProperties is a struct field with a minimum number of properties.
	// +kubebuilder:validation:MinProperties=1
	// +optional
	StructWithMinProperties B `json:"structWithMinProperties,omitempty"` // want "field StructWithMinProperties is optional and should be a pointer"

	// structWithMinPropertiesOnStruct is a struct field with a minimum number of properties on the struct.
	// +optional
	StructWithMinPropertiesOnStruct D `json:"structWithMinPropertiesOnStruct,omitempty"` // want "field StructWithMinPropertiesOnStruct is optional and should be a pointer"

	// slice is a slice field.
	// +optional
	Slice []string `json:"slice,omitempty"`

	// map is a map field.
	// +optional
	Map map[string]string `json:"map,omitempty"`

	// PointerSlice is a pointer slice field.
	// +optional
	PointerSlice *[]string `json:"pointerSlice,omitempty"` // want "field PointerSlice is a pointer type and should not be a pointer"

	// PointerMap is a pointer map field.
	// +optional
	PointerMap *map[string]string `json:"pointerMap,omitempty"` // want "field PointerMap is a pointer type and should not be a pointer"

	// PointerPointerString is a double pointer string field.
	// +optional
	DoublePointerString **string `json:"doublePointerString,omitempty"` // want "field DoublePointerString is a pointer type and should not be a pointer"
}

type B struct {
	// pointerString is a pointer string field.
	// +optional
	PointerString *string `json:"pointerString,omitempty"`
}

// +kubebuilder:validation:MinProperties=1
type D struct {
	// string is a string field.
	// +optional
	String *string `json:"string,omitempty"`

	// stringWithMinLength1 with minimum length is a string field.
	// +kubebuilder:validation:MinLength=1
	// +optional
	StringWithMinLength1 *string `json:"stringWithMinLength1,omitempty"`
}
