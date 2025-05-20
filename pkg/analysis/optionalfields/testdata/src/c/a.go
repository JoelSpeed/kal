package a

type A struct {
	// required field should not be picked up.
	// +required
	RequiredField string `json:"requiredField,omitempty"`

	// pointerString is a pointer string field.
	// +optional
	PointerString *string `json:"pointerString,omitempty"` // want "field PointerString is an optional string and does not have a minimum length. Where the difference between omitted and the empty string is significant, set the minmum length to 0"

	// pointerStringWithoutOmitEmpty is a pointer string field without omitempty/
	// +optional
	PointerStringWithoutOmitEmpty *string `json:"pointerStringWithoutOmitEmpty"` // want "field PointerStringWithoutOmitEmpty is an optional string without omitempty. It should not be a pointer."

	// pointerStringWithMinLength1 with minimum length is a pointer string field.
	// +kubebuilder:validation:MinLength=1
	// +optional
	PointerStringWithMinLength1 *string `json:"pointerStringWithMinLength1,omitempty"` // want "field PointerStringWithMinLength1 has a greater than 0 length and does not need to be a pointer"

	// pointerStringWithMinLength1WithoutOmitEmpty with minimum length is a pointer string field without omitempty.
	// +kubebuilder:validation:MinLength=1
	// +optional
	PointerStringWithMinLength1WithoutOmitEmpty *string `json:"pointerStringWithMinLength1WithoutOmitEmpty"` // want "field PointerStringWithMinLength1WithoutOmitEmpty is an optional string without omitempty. It should not be a pointer." "field PointerStringWithMinLength1WithoutOmitEmpty has a greater than zero minimum length without omitempty. The minimum length should be removed."

	// pointerStringWithMinLength0 with minimum length is a pointer string field.
	// +kubebuilder:validation:MinLength=0
	// +optional
	PointerStringWithMinLength0 *string `json:"pointerStringWithMinLength0,omitempty"`

	// pointerStringWithMinLength0WithoutOmitempty with minimum length is a pointer string field without omitempty.
	// +kubebuilder:validation:MinLength=0
	// +optional
	PointerStringWithMinLength0WithoutOmitempty *string `json:"pointerStringWithMinLength0WithoutOmitempty"` // want "field PointerStringWithMinLength0WithoutOmitempty is an optional string without omitempty. It should not be a pointer."

	// pointerInt is a pointer int field.
	// +optional
	PointerInt *int `json:"pointerInt,omitempty"` // want "field PointerInt is an optional integer and does not have a minimum/maximum value. Where the difference between omitted and 0 is significant, set the minimum/maximum value to a range including 0"

	// pointerIntWithoutOmitEmpty is a pointer int field with omitempty.
	// +optional
	PointerIntWithoutOmitEmpty *int `json:"pointerIntWithoutOmitEmpty"` // want "field PointerIntWithoutOmitEmpty is an optional integer without omitempty. It should not be a pointer."

	// pointerIntWithMinValue1 with minimum value is a pointer int field.
	// +kubebuilder:validation:Minimum=1
	// +optional
	PointerIntWithMinValue1 *int `json:"pointerIntWithMinValue1,omitempty"` // want "field PointerIntWithMinValue1 has a greater than 0 minimum value and does not need to be a pointer"

	// pointerIntWithMinValue1WithoutOmitEmpty with minimum value is a pointer int field without omitempty.
	// +kubebuilder:validation:Minimum=1
	// +optional
	PointerIntWithMinValue1WithoutOmitEmpty *int `json:"pointerIntWithMinValue1WithoutOmitEmpty"` // want "field PointerIntWithMinValue1WithoutOmitEmpty is an optional integer without omitempty. It should not be a pointer." "field PointerIntWithMinValue1WithoutOmitEmpty has a greater than zero minimum value without omitempty. The minimum value should be removed."

	// pointerIntWithMinValue0 with minimum value is a pointer int field.
	// +kubebuilder:validation:Minimum=0
	// +optional
	PointerIntWithMinValue0 *int `json:"pointerIntWithMinValue0,omitempty"`

	// pointerIntWithMinValue0WithoutOmitEmpty with minimum value is a pointer int field without omitempty.
	// +kubebuilder:validation:Minimum=0
	// +optional
	PointerIntWithMinValue0WithoutOmitEmpty *int `json:"pointerIntWithMinValue0WithoutOmitEmpty"` // want "field PointerIntWithMinValue0WithoutOmitEmpty is an optional integer without omitempty. It should not be a pointer."

	// pointerIntWithNegativeMaximumValue with negative maximum value is a pointer int field.
	// +kubebuilder:validation:Maximum=-1
	// +optional
	PointerIntWithNegativeMaximumValue *int `json:"pointerIntWithNegativeMaximumValue,omitempty"` // want "field PointerIntWithNegativeMaximumValue has a negative maximum value and does not need to be a pointer"

	// pointerIntWithNegativeMaximumValueWithoutOmitEmpty with negative maximum value is a pointer int field without omitempty.
	// +kubebuilder:validation:Maximum=-1
	// +optional
	PointerIntWithNegativeMaximumValueWithoutOmitEmpty *int `json:"pointerIntWithNegativeMaximumValueWithoutOmitEmpty"` // want "field PointerIntWithNegativeMaximumValueWithoutOmitEmpty is an optional integer without omitempty. It should not be a pointer." "field PointerIntWithNegativeMaximumValueWithoutOmitEmpty has a less than zero maximum value without omitempty. The maximum value should be removed."

	// pointerIntWithNegativeMinimumValue with negative minimum value is a pointer int field.
	// +kubebuilder:validation:Minimum=-1
	// +optional
	PointerIntWithNegativeMinimumValue *int `json:"pointerIntWithNegativeMinimumValue,omitempty"` // want "field PointerIntWithNegativeMinimumValue has a negative minimum value and does not have a maximum value. A maximum value should be set"

	// pointerIntWithPositiveMaximumValue with positive maximum value is a pointer int field.
	// +kubebuilder:validation:Maximum=1
	// +optional
	PointerIntWithPositiveMaximumValue *int `json:"pointerIntWithPositiveMaximumValue,omitempty"` // want "field PointerIntWithPositiveMaximumValue has a positive maximum value and does not have a minimum value. A minimum value should be set"

	// pointerIntWithPositiveMaximumValueWithoutOmitEmpty with positive maximum value is a pointer int field without omitempty.
	// +kubebuilder:validation:Maximum=1
	// +optional
	PointerIntWithPositiveMaximumValueWithoutOmitEmpty *int `json:"pointerIntWithPositiveMaximumValueWithoutOmitEmpty"` // want "field PointerIntWithPositiveMaximumValueWithoutOmitEmpty is an optional integer without omitempty. It should not be a pointer."

	// pointerIntWithRange is a pointer int field with a range of values including 0.
	// +kubebuilder:validation:Minimum=-10
	// +kubebuilder:validation:Maximum=10
	// +optional
	PointerIntWithRange *int `json:"pointerIntWithRange,omitempty"`

	// pointerIntWithRangeWithoutOmitEmpty is a pointer int field with a range of values including 0 wihtout omitempty.
	// +kubebuilder:validation:Minimum=-10
	// +kubebuilder:validation:Maximum=10
	// +optional
	PointerIntWithRangeWithoutOmitEmpty *int `json:"pointerIntWithRangeWithoutOmitEmpty"` // want "field PointerIntWithRangeWithoutOmitEmpty is an optional integer without omitempty. It should not be a pointer."

	// pointerStruct is a pointer struct field.
	// +optional
	PointerStruct *B `json:"pointerStruct,omitempty"` // want "field PointerStruct is optional, and contains no required field\\(s\\) and does not need to be a pointer"

	// pointerStructWithoutOmitEmpty is a pointer struct field without omitempty.
	// +optional
	PointerStructWithoutOmitEmpty *B `json:"pointerStructWithoutOmitEmpty"` // want "field PointerStructWithoutOmitEmpty is an optional struct without omitempty. It should not be a pointer"

	// string is a string field.
	// +optional
	String string `json:"string,omitempty"` // want "field String is an optional string and does not have a minimum length. Either set a minimum length or make String a pointer where the difference between omitted and the empty string is significant"

	// stringWithoutOmitEmpty is a string field without omitempty.
	// +optional
	StringWithoutOmitEmpty string `json:"stringWithoutOmitEmpty"`

	// stringWithMinLength1 with minimum length is a string field.
	// +kubebuilder:validation:MinLength=1
	// +optional
	StringWithMinLength1 string `json:"stringWithMinLength1,omitempty"`

	// stringWithMinLength1WithoutOmitEmpty with minimum length is a string field without omitempty.
	// +kubebuilder:validation:MinLength=1
	// +optional
	StringWithMinLength1WithoutOmitEmpty string `json:"stringWithMinLength1WithoutOmitEmpty"` // want "field StringWithMinLength1WithoutOmitEmpty has a greater than zero minimum length without omitempty. The minimum length should be removed."

	// stringWithMinLength0 with minimum length is a string field.
	// +kubebuilder:validation:MinLength=0
	// +optional
	StringWithMinLength0 string `json:"stringWithMinLength0,omitempty"` // want "field StringWithMinLength0 has a minimum length of 0. The empty string is a valid value and therefore the field should be a pointer"

	// stringWithMinLength0WithoutOmitEmpty with minimum length is a string field without omitempty.
	// +kubebuilder:validation:MinLength=0
	// +optional
	StringWithMinLength0WithoutOmitEmpty string `json:"stringWithMinLength0WithoutOmitEmpty"`

	// int is an int field.
	// +optional
	Int int `json:"int,omitempty"` // want "field Int is an optional integer and does not have a minimum/maximum value. Either set a minimum/maximum value or make Int a pointer where the difference between omitted and 0 is significant"

	// intWithMinValue1 with minimum value is an int field.
	// +kubebuilder:validation:Minimum=1
	// +optional
	IntWithMinValue1 int `json:"intWithMinValue1,omitempty"`

	// intWithMinValue1WithoutOmitEmpty with minimum value is an int field without omitempty.
	// +kubebuilder:validation:Minimum=1
	// +optional
	IntWithMinValue1WithoutOmitEmpty int `json:"intWithMinValue1WithoutOmitEmpty"` // want "field IntWithMinValue1WithoutOmitEmpty has a greater than zero minimum value without omitempty. The minimum value should be removed."

	// intWithMinValue0 with minimum value is an int field.
	// +kubebuilder:validation:Minimum=0
	// +optional
	IntWithMinValue0 int `json:"intWithMinValue0,omitempty"` // want "field IntWithMinValue0 has a range of values including 0. The difference between omitted and 0 is significant and therefore the field should be a pointer"

	// intWithMinValue0WithoutOmitEmpty with minimum value is an int field without omitempty.
	// +kubebuilder:validation:Minimum=0
	// +optional
	IntWithMinValue0WithoutOmitEmpty int `json:"intWithMinValue0WithoutOmitEmpty"`

	// intWithNegativeMaximumValue with negative maximum value is an int field.
	// +kubebuilder:validation:Maximum=-1
	// +optional
	IntWithNegativeMaximumValue int `json:"intWithNegativeMaximumValue,omitempty"`

	// intWithNegativeMaximumValueWithoutOmitEmpty with negative maximum value is an int field without omitempty.
	// +kubebuilder:validation:Maximum=-1
	// +optional
	IntWithNegativeMaximumValueWithoutOmitEmpty int `json:"intWithNegativeMaximumValueWithoutOmitEmpty"` // want "field IntWithNegativeMaximumValueWithoutOmitEmpty has a less than zero maximum value without omitempty. The maximum value should be removed."

	// intWithNegativeMinimumValue with negative minimum value is an int field.
	// +kubebuilder:validation:Minimum=-1
	// +optional
	IntWithNegativeMinimumValue int `json:"intWithNegativeMinimumValue,omitempty"` // want "field IntWithNegativeMinimumValue has a negative minimum value and does not have a maximum value. A maximum value should be set"

	// intWithNegativeMinimumValueWithoutOmitEmpty with negative minimum value is an int field without omitempty.
	// +kubebuilder:validation:Minimum=-1
	// +optional
	IntWithNegativeMinimumValueWithoutOmitEmpty int `json:"intWithNegativeMinimumValueWithoutOmitEmpty"`

	// intWithPositiveMaximumValue with positive maximum value is an int field.
	// +kubebuilder:validation:Maximum=1
	// +optional
	IntWithPositiveMaximumValue int `json:"intWithPositiveMaximumValue,omitempty"` // want "field IntWithPositiveMaximumValue has a positive maximum value and does not have a minimum value. A minimum value should be set"

	// intWithPositiveMaximumValueWithoutOmitEmpty with positive maximum value is an int field without omitempty.
	// +kubebuilder:validation:Maximum=1
	// +optional
	IntWithPositiveMaximumValueWithoutOmitEmpty int `json:"intWithPositiveMaximumValueWithoutOmitEmpty"`

	// intWithRange is an int field with a range of values including 0.
	// +kubebuilder:validation:Minimum=-10
	// +kubebuilder:validation:Maximum=10
	// +optional
	IntWithRange int `json:"intWithRange,omitempty"` // want "field IntWithRange has a range of values including 0. The difference between omitted and 0 is significant and therefore the field should be a pointer"

	// intWithRangeWithoutOmitEmpty is an int field with a range of values including 0 without omitempty.
	// +kubebuilder:validation:Minimum=-10
	// +kubebuilder:validation:Maximum=10
	// +optional
	IntWithRangeWithoutOmitEmpty int `json:"intWithRangeWithoutOmitEmpty"`

	// intWithInvalidMinimumValue with invalid minimum value is an int field.
	// +kubebuilder:validation:Minimum=foo
	// +optional
	IntWithInvalidMinimumValue int `json:"intWithInvalidMinimumValue,omitempty"` // want "field IntWithInvalidMinimumValue has a minimum value of foo, but it is not an integer"

	// intWithInvalidMaximumValue with invalid maximum value is an int field.
	// +kubebuilder:validation:Maximum=foo
	// +optional
	IntWithInvalidMaximumValue int `json:"intWithInvalidMaximumValue,omitempty"` // want "field IntWithInvalidMaximumValue has a maximum value of foo, but it is not an integer"

	// float is a float field.
	// +optional
	Float float64 `json:"float,omitempty"` // want "field Float is an optional float and does not have a minimum/maximum value. Either set a minimum/maximum value or make Float a pointer where the difference between omitted and 0 is significant"

	// floatWithoutOmitEmpty is a float field without omitempty.
	// +optional
	FloatWithoutOmitEmpty float64 `json:"floatWithoutOmitEmpty"`

	// floatWithMinValue1 with minimum value is a float field.
	// +kubebuilder:validation:Minimum=1.0
	// +optional
	FloatWithMinValue1 float64 `json:"floatWithMinValue1,omitempty"`

	// floatWithMinValue1WithoutOmitEmpty with minimum value is a float field without omitempty.
	// +kubebuilder:validation:Minimum=1.0
	// +optional
	FloatWithMinValue1WithoutOmitEmpty float64 `json:"floatWithMinValue1WithoutOmitEmpty"` // want "field FloatWithMinValue1WithoutOmitEmpty has a greater than zero minimum value without omitempty. The minimum value should be removed."

	// floatWithMinValue0 with minimum value is a float field.
	// +kubebuilder:validation:Minimum=0.0
	// +optional
	FloatWithMinValue0 float64 `json:"floatWithMinValue0,omitempty"` // want "field FloatWithMinValue0 has a range of values including 0. The difference between omitted and 0 is significant and therefore the field should be a pointer"

	// floatWithMinValue0WithoutOmitEmpty with minimum value is a float field without omitempty.
	// +kubebuilder:validation:Minimum=0.0
	// +optional
	FloatWithMinValue0WithoutOmitEmpty float64 `json:"floatWithMinValue0WithoutOmitEmpty"`

	// floatWithNegativeMaximumValue with negative maximum value is a float field.
	// +kubebuilder:validation:Maximum=-1.0
	// +optional
	FloatWithNegativeMaximumValue float64 `json:"floatWithNegativeMaximumValue,omitempty"`

	// floatWithNegativeMaximumValueWithoutOmitEmpty with negative maximum value is a float field without omitempty.
	// +kubebuilder:validation:Maximum=-1.0
	// +optional
	FloatWithNegativeMaximumValueWithoutOmitEmpty float64 `json:"floatWithNegativeMaximumValueWithoutOmitEmpty"` // want "field FloatWithNegativeMaximumValueWithoutOmitEmpty has a less than zero maximum value without omitempty. The maximum value should be removed."

	// floatWithNegativeMinimumValue with negative minimum value is a float field.
	// +kubebuilder:validation:Minimum=-1.0
	// +optional
	FloatWithNegativeMinimumValue float64 `json:"floatWithNegativeMinimumValue,omitempty"` // want "field FloatWithNegativeMinimumValue has a negative minimum value and does not have a maximum value. A maximum value should be set"

	// floatWithNegativeMinimumValueWithoutOmitEmpty with negative minimum value is a float field without omitempty.
	// +kubebuilder:validation:Minimum=-1.0
	// +optional
	FloatWithNegativeMinimumValueWithoutOmitEmpty float64 `json:"floatWithNegativeMinimumValueWithoutOmitEmpty"`

	// floatWithPositiveMaximumValue with positive maximum value is a float field.
	// +kubebuilder:validation:Maximum=1.0
	// +optional
	FloatWithPositiveMaximumValue float64 `json:"floatWithPositiveMaximumValue,omitempty"` // want "field FloatWithPositiveMaximumValue has a positive maximum value and does not have a minimum value. A minimum value should be set"

	// floatWithPositiveMaximumValueWithoutOmitEmpty with positive maximum value is a float field without omitempty.
	// +kubebuilder:validation:Maximum=1.0
	// +optional
	FloatWithPositiveMaximumValueWithoutOmitEmpty float64 `json:"floatWithPositiveMaximumValueWithoutOmitEmpty"`

	// floatWithRange is a float field with a range of values including 0.
	// +kubebuilder:validation:Minimum=-10.0
	// +kubebuilder:validation:Maximum=10.0
	// +optional
	FloatWithRange float64 `json:"floatWithRange,omitempty"` // want "field FloatWithRange has a range of values including 0. The difference between omitted and 0 is significant and therefore the field should be a pointer"

	// floatWithRangeWithoutOmitEmpty is a float field with a range of values including 0 without omitempty.
	// +kubebuilder:validation:Minimum=-10.0
	// +kubebuilder:validation:Maximum=10.0
	// +optional
	FloatWithRangeWithoutOmitEmpty float64 `json:"floatWithRangeWithoutOmitEmpty"`

	// floatWithInvalidMinimumValue with invalid minimum value is a float field.
	// +kubebuilder:validation:Minimum=foo
	// +optional
	FloatWithInvalidMinimumValue float64 `json:"floatWithInvalidMinimumValue,omitempty"` // want "field FloatWithInvalidMinimumValue has a minimum value of foo, but it is not a float"

	// floatWithInvalidMaximumValue with invalid maximum value is a float field.
	// +kubebuilder:validation:Maximum=foo
	// +optional
	FloatWithInvalidMaximumValue float64 `json:"floatWithInvalidMaximumValue,omitempty"` // want "field FloatWithInvalidMaximumValue has a maximum value of foo, but it is not a float"

	// structWithOptionalFields is a struct field.
	// +optional
	StructWithOptionalFields B `json:"structWithOptionalFields,omitempty"`

	// structWithOptionalFieldsWithoutOmitEmpty is a struct field without omitempty.
	// +optional
	StructWithOptionalFieldsWithoutOmitEmpty B `json:"structWithOptionalFieldsWithoutOmitEmpty"`

	// structWithMinProperties is a struct field with a minimum number of properties.
	// +kubebuilder:validation:MinProperties=1
	// +optional
	StructWithMinProperties B `json:"structWithMinProperties,omitempty"` // want "field StructWithMinProperties has a greater than zero minimum number of properties and should be a pointer"

	// structWithMinPropertiesWithoutOmitEmpty is a struct field with a minimum number of properties without omitempty.
	// +kubebuilder:validation:MinProperties=1
	// +optional
	StructWithMinPropertiesWithoutOmitEmpty B `json:"structWithMinPropertiesWithoutOmitEmpty"` // want "field StructWithMinPropertiesWithoutOmitEmpty has a greater than zero minimum number of properties without omitempty. The minimum number of properties should be removed."

	// structWithMinPropertiesOnStruct is a struct field with a minimum number of properties on the struct.
	// +optional
	StructWithMinPropertiesOnStruct D `json:"structWithMinPropertiesOnStruct,omitempty"` // want "field StructWithMinPropertiesOnStruct has a greater than zero minimum number of properties and should be a pointer"

	// structWithMinPropertiesOnStructWithoutOmitEmpty is a struct field with a minimum number of properties on the struct without omitempty.
	// +optional
	StructWithMinPropertiesOnStructWithoutOmitEmpty D `json:"structWithMinPropertiesOnStructWithoutOmitEmpty"` // want "field StructWithMinPropertiesOnStructWithoutOmitEmpty has a greater than zero minimum number of properties and should be a pointer"

	// structWithRequiredFields is a struct field.
	// +optional
	StructWithRequiredFields C `json:"structWithRequiredFields,omitempty"` // want "field StructWithRequiredFields is optional, but contains required field\\(s\\) and should be a pointer"

	// structWithRequiredFieldsWithoutOmitEmpty is a struct field without omitempty.
	// +optional
	StructWithRequiredFieldsWithoutOmitEmpty C `json:"structWithRequiredFieldsWithoutOmitEmpty"`

	// pointerStructWithOptionalFields is a pointer struct field.
	// +optional
	PointerStructWithOptionalFields *B `json:"pointerStructWithOptionalFields,omitempty"` // want "field PointerStructWithOptionalFields is optional, and contains no required field\\(s\\) and does not need to be a pointer"

	// pointerStructWithOptionalFieldsWithoutOmitEmpty is a pointer struct field without omitempty.
	// +optional
	PointerStructWithOptionalFieldsWithoutOmitEmpty *B `json:"pointerStructWithOptionalFieldsWithoutOmitEmpty"` // want "field PointerStructWithOptionalFieldsWithoutOmitEmpty is an optional struct without omitempty. It should not be a pointer"

	// pointerStructWithRequiredFields is a pointer struct field.
	// +optional
	PointerStructWithRequiredFields *C `json:"pointerStructWithRequiredFields,omitempty"`

	// pointerStructWithRequiredFieldsWithoutOmitEmpty is a pointer struct field without omitempty.
	// +optional
	PointerStructWithRequiredFieldsWithoutOmitEmpty *C `json:"pointerStructWithRequiredFieldsWithoutOmitEmpty"` // want "field PointerStructWithRequiredFieldsWithoutOmitEmpty is an optional struct without omitempty. It should not be a pointer"

	// bool is a boolean field.
	// +optional
	Bool bool `json:"bool,omitempty"` // want "field Bool is an optional boolean and should be a pointer"

	// boolWithoutOmitEmpty is a boolean field without omitempty.
	// +optional
	BoolWithoutOmitEmpty bool `json:"boolWithoutOmitEmpty"`

	// boolPointer is a pointer boolean field.
	// +optional
	BoolPointer *bool `json:"boolPointer,omitempty"`

	// boolPointerWithoutOmitEmpty is a pointer boolean field without omitempty.
	// +optional
	BoolPointerWithoutOmitEmpty *bool `json:"boolPointerWithoutOmitEmpty"` // want "field BoolPointerWithoutOmitEmpty is an optional boolean without omitempty. It should not be a pointer"

	// slice is a slice field.
	// +optional
	Slice []string `json:"slice,omitempty"`

	// sliceWithMinItems is a slice field with a minimum number of items.
	// +kubebuilder:validation:MinItems=1
	// +optional
	SliceWithMinItems []string `json:"sliceWithMinItems,omitempty"`

	// sliceWithMinItemsWithoutOmitEmpty is a slice field with a minimum number of items without omitempty.
	// +kubebuilder:validation:MinItems=1
	// +optional
	SliceWithMinItemsWithoutOmitEmpty []string `json:"sliceWithMinItemsWithoutOmitEmpty"` // want "field SliceWithMinItemsWithoutOmitEmpty has a greater than zero minimum number of items without omitempty. The minimum number of items should be removed."

	// map is a map field.
	// +optional
	Map map[string]string `json:"map,omitempty"`

	// mapWithMinProperties is a map field with a minimum number of properties.
	// +kubebuilder:validation:MinProperties=1
	// +optional
	MapWithMinProperties map[string]string `json:"mapWithMinProperties,omitempty"`

	// mapWithMinPropertiesWithoutOmitEmpty is a map field with a minimum number of properties without omitempty.
	// +kubebuilder:validation:MinProperties=1
	// +optional
	MapWithMinPropertiesWithoutOmitEmpty map[string]string `json:"mapWithMinPropertiesWithoutOmitEmpty"` // want "field MapWithMinPropertiesWithoutOmitEmpty has a greater than zero minimum number of properties without omitempty. The minimum number of properties should be removed."

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
	// +kubebuilder:validation:MinLength=0
	// +optional
	PointerString *string `json:"pointerString,omitempty"`
}

type C struct {
	// string is a string field.
	// +required
	String string `json:"string"`
}

// +kubebuilder:validation:MinProperties=1
type D struct {
	// string is a string field.
	// +optional
	String string `json:"string"`

	// stringWithMinLength1 with minimum length is a string field.
	// +kubebuilder:validation:MinLength=1
	// +optional
	StringWithMinLength1 string `json:"stringWithMinLength1,omitempty"`
}
