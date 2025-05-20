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
package optionalfields

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	kalerrors "sigs.k8s.io/kube-api-linter/pkg/analysis/errors"
	"sigs.k8s.io/kube-api-linter/pkg/analysis/helpers/extractjsontags"
	"sigs.k8s.io/kube-api-linter/pkg/analysis/helpers/inspector"
	markershelper "sigs.k8s.io/kube-api-linter/pkg/analysis/helpers/markers"
	"sigs.k8s.io/kube-api-linter/pkg/analysis/utils"
	"sigs.k8s.io/kube-api-linter/pkg/config"
	"sigs.k8s.io/kube-api-linter/pkg/markers"

	"k8s.io/utils/ptr"
)

const (
	name = "optionalfields"

	optionalMarker            = markers.OptionalMarker
	requiredMarker            = markers.RequiredMarker
	kubebuilderOptionalMarker = markers.KubebuilderOptionalMarker
	kubebuilderRequiredMarker = markers.KubebuilderRequiredMarker

	minItemsMarker      = markers.KubebuilderMinItemsMarker
	minLengthMarker     = markers.KubebuilderMinLengthMarker
	minPropertiesMarker = markers.KubebuilderMinPropertiesMarker

	minimumMarker = markers.KubebuilderMinimumMarker
	maximumMarker = markers.KubebuilderMaximumMarker
)

func init() {
	markershelper.DefaultRegistry().Register(
		optionalMarker,
		requiredMarker,
		kubebuilderOptionalMarker,
		kubebuilderRequiredMarker,
		minItemsMarker,
		minLengthMarker,
		minPropertiesMarker,
		minimumMarker,
		maximumMarker,
	)
}

var (
	errMarkerMissingValue = errors.New("marker does not have a value")
)

type analyzer struct {
	pointerPolicy     config.OptionalFieldsPointerPolicy
	pointerPreference config.OptionalFieldsPointerPreference
	omitEmptyPolicy   config.OptionalFieldsOmitEmptyPolicy
}

// newAnalyzer creates a new analyzer.
func newAnalyzer(cfg config.OptionalFieldsConfig) *analysis.Analyzer {
	defaultConfig(&cfg)

	a := &analyzer{
		pointerPolicy:     cfg.Pointers.Policy,
		pointerPreference: cfg.Pointers.Preference,
		omitEmptyPolicy:   cfg.OmitEmpty.Policy,
	}

	return &analysis.Analyzer{
		Name: name,
		Doc: `Checks all optional fields comply with the configured policy.
		Depending on the configuration, this may include checking for the presence of the omitempty tag or
		whether the field is a pointer.
		For structs, this includes checking that if the field is marked as optional, it should be a pointer when it has omitempty.
		Where structs include required fields, they must be a pointer when they themselves are optional.
		`,
		Run:      a.run,
		Requires: []*analysis.Analyzer{inspector.Analyzer},
	}
}

func (a *analyzer) run(pass *analysis.Pass) (any, error) {
	inspect, ok := pass.ResultOf[inspector.Analyzer].(inspector.Inspector)
	if !ok {
		return nil, kalerrors.ErrCouldNotGetInspector
	}

	inspect.InspectFields(func(field *ast.Field, stack []ast.Node, jsonTagInfo extractjsontags.FieldTagInfo, markersAccess markershelper.Markers) {
		a.checkField(pass, field, markersAccess, jsonTagInfo)
	})

	return nil, nil //nolint:nilnil
}

func (a *analyzer) checkField(pass *analysis.Pass, field *ast.Field, markersAccess markershelper.Markers, jsonTags extractjsontags.FieldTagInfo) {
	if field == nil || len(field.Names) == 0 {
		return
	}

	fieldMarkers := markersAccess.FieldMarkers(field)

	fieldName := field.Names[0].Name

	if !isFieldOptional(fieldMarkers) {
		// The field is not marked optional, so we don't need to check it.
		return
	}

	a.checkFieldOmitEmpty(pass, field, fieldName, jsonTags)

	if field.Type == nil {
		// The field has no type? We can't check if it's a pointer.
		return
	}

	a.checkFieldPointers(pass, field, fieldName, markersAccess, jsonTags)
}

// checkFieldOmitEmpty checks if the field has the omitempty tag.
// Normally, all optional fields should have the omitempty tag.
// In certain cases, you may not want to have the omitempty tag
// (such as to be able to marshal an empty object with all possible options shown)
// and when this is desired, there is an option to ignore the omitempty tag.
func (a *analyzer) checkFieldOmitEmpty(pass *analysis.Pass, field *ast.Field, fieldName string, jsonTags extractjsontags.FieldTagInfo) {
	switch {
	case jsonTags.OmitEmpty, a.omitEmptyPolicy == config.OptionalFieldsOmitEmptyPolicyIgnore:
		// Nothing to do, either we have omitempty, or we are ignoring it.
	case a.omitEmptyPolicy == config.OptionalFieldsOmitEmptyPolicySuggestFix:
		pass.Report(analysis.Diagnostic{
			Pos:     field.Pos(),
			Message: fmt.Sprintf("field %s is optional and should be omitempty", fieldName),
			SuggestedFixes: []analysis.SuggestedFix{
				{
					Message: "add the omitempty tag",
					TextEdits: []analysis.TextEdit{
						{
							Pos:     jsonTags.Pos + token.Pos(len(jsonTags.Name)),
							NewText: []byte(",omitempty"),
						},
					},
				},
			},
		})
	case a.omitEmptyPolicy == config.OptionalFieldsOmitEmptyPolicyWarn:
		pass.Reportf(field.Pos(), "field %s is optional and should be omitempty", fieldName)
	}
}

// checkFieldPointers is used to determine if a field should be a pointer, and advise on the correct action.
func (a *analyzer) checkFieldPointers(pass *analysis.Pass, field *ast.Field, fieldName string, markersAccess markershelper.Markers, jsonTags extractjsontags.FieldTagInfo) {
	isStarExpr, underlyingType := isStarExpr(field.Type)

	if isPointerType(underlyingType) {
		a.checkFieldPointersPointerTypes(pass, field, fieldName, isStarExpr, markersAccess, jsonTags)

		return
	}

	switch a.pointerPreference {
	case config.OptionalFieldsPointerPreferenceAlways:
		a.checkFieldPointersPreferenceAlways(pass, field, fieldName, isStarExpr)
	case config.OptionalFieldsPointerPreferenceWhenRequired:
		a.checkFieldPointersPreferenceWhenRequired(pass, field, fieldName, isStarExpr, underlyingType, markersAccess, jsonTags)
	}
}

// checkFieldPointersPointerTypes checks for pointer types (eg Maps and Arrays)
// and ensures that they aren't pointered again (e.g. no *[]string).
// It will also check for the presence of non-zero min-properties and min-items
// when omitempty is missing, as this will.
func (a *analyzer) checkFieldPointersPointerTypes(pass *analysis.Pass, field *ast.Field, fieldName string, isStarExpr bool, markersAccess markershelper.Markers, jsonTags extractjsontags.FieldTagInfo) {
	if a.omitEmptyPolicy == config.OptionalFieldsOmitEmptyPolicyIgnore && !jsonTags.OmitEmpty {
		a.checkFieldPointersPointerTypesWithoutOmitEmpty(pass, field, fieldName, markersAccess)
	}

	// Pointer types should not be pointered again.
	if !isStarExpr {
		return
	}

	reportShouldRemovePointer(pass, field, a.pointerPolicy, fieldName, "field %s is a pointer type and should not be a pointer")
}

// checkFieldPointersPointerTypesWithoutOmitEmpty handles the case where the field is a pointer type (array or map)
// and does not have the omitempty tag.
// In both cases, the field should not have a minimum number of items or properties greater than 0.
// Without omitempty, empty json objects and arrays would be rendered ('{}' and '[]') and would breach the minimum properties/items.
func (a *analyzer) checkFieldPointersPointerTypesWithoutOmitEmpty(pass *analysis.Pass, field *ast.Field, fieldName string, markersAccess markershelper.Markers) {
	switch field.Type.(type) {
	case *ast.MapType:
		reportShouldRemoveAllInstancesOfIntegerMarker(pass, field, markersAccess, minPropertiesMarker, fieldName, "field %s has a greater than zero minimum number of properties without omitempty. The minimum number of properties should be removed.")
	case *ast.ArrayType:
		reportShouldRemoveAllInstancesOfIntegerMarker(pass, field, markersAccess, minItemsMarker, fieldName, "field %s has a greater than zero minimum number of items without omitempty. The minimum number of items should be removed.")
	}
}

// checkFieldPointersPreferenceAlways checks if the field is a pointer and if not, suggests that it should be.
func (a *analyzer) checkFieldPointersPreferenceAlways(pass *analysis.Pass, field *ast.Field, fieldName string, isStarExpr bool) {
	if isStarExpr {
		return // The field is already a pointer, so we don't need to do anything.
	}

	reportShouldAddPointer(pass, field, a.pointerPolicy, fieldName, "field %s is optional and should be a pointer")
}

// checkFieldPointersPreferenceWhenRequired checks if the field needs to be a pointer.
// This means, is the zero value valid as a user desired value.
// For example, if it's an integer range that includes 0, then it should be a pointer.
// If the range doesn't include zero, and the field has omitempty, then the field doesn't need to be a pointer
// As `0` should never be committed to the API, and round trips would not be affected by the JSON library omitting the zero value when marshalling.
func (a *analyzer) checkFieldPointersPreferenceWhenRequired(pass *analysis.Pass, field *ast.Field, fieldName string, isStarExpr bool, underlyingType ast.Expr, markersAccess markershelper.Markers, jsonTags extractjsontags.FieldTagInfo) {
	ident, ok := underlyingType.(*ast.Ident)
	if !ok {
		// All fields should be idents, not sure when this would happen?
		return
	}

	switch ident.Name {
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		a.checkFieldPointersPreferenceWhenRequiredInteger(pass, field, fieldName, isStarExpr, markersAccess, jsonTags)
		return
	case "string":
		a.checkFieldPointersPreferenceWhenRequiredString(pass, field, fieldName, isStarExpr, markersAccess, jsonTags)
		return
	case "bool":
		a.checkFieldPointersPreferenceWhenRequiredBool(pass, field, fieldName, isStarExpr, jsonTags)
		return
	case "float32", "float64":
		a.checkFieldPointersPreferenceWhenRequiredFloat(pass, field, fieldName, isStarExpr, markersAccess, jsonTags)
		return
	}

	// The field is not a simple type in our switch, so try looking up the type spec.
	typeSpec, ok := utils.LookupTypeSpec(pass, ident)
	if !ok {
		return
	}

	a.checkFieldPointersPreferenceWhenRequiredIdentObj(pass, field, fieldName, isStarExpr, typeSpec, markersAccess, jsonTags)
}

func (a *analyzer) checkFieldPointersPreferenceWhenRequiredIdentObj(pass *analysis.Pass, field *ast.Field, fieldName string, isStarExpr bool, decl *ast.TypeSpec, markersAccess markershelper.Markers, jsonTags extractjsontags.FieldTagInfo) {
	switch t := decl.Type.(type) {
	case *ast.StructType:
		a.checkFieldPointersPreferenceWhenRequiredStructType(pass, field, fieldName, isStarExpr, t, markersAccess, jsonTags)
	case *ast.Ident:
		// The field is using a type alias.
		a.checkFieldPointersPreferenceWhenRequired(pass, field, fieldName, isStarExpr, t, markersAccess, jsonTags)
	case *ast.ArrayType, *ast.MapType:
		a.checkFieldPointersPointerTypes(pass, field, fieldName, isStarExpr, markersAccess, jsonTags)
	}
}

// checkFieldPointersPreferenceWhenRequiredStructType determines if the struct field needs to be a pointer.
// Any struct that has a minimum number of properties, or has required fields, should be a pointer.
// Without a pointer, the JSON library cannot omit the field, and will always render a `{}`.
// A rendered empty object would then violate the minimum number of properties/required field checks.
func (a *analyzer) checkFieldPointersPreferenceWhenRequiredStructType(pass *analysis.Pass, field *ast.Field, fieldName string, isStarExpr bool, typeExpr *ast.StructType, markersAccess markershelper.Markers, jsonTags extractjsontags.FieldTagInfo) {
	hasRequiredFields := structContainsRequiredFields(typeExpr, markersAccess)

	hasMinimumProperties, err := structHasGreaterThanZeroMinProperties(typeExpr, markersAccess.StructMarkers(typeExpr))
	if err != nil {
		pass.Reportf(field.Pos(), "error checking struct for min properties: %v", err)
		return
	}

	fieldHasMinimumProperties, err := structHasGreaterThanZeroMinProperties(typeExpr, markersAccess.FieldMarkers(field))
	if err != nil {
		pass.Reportf(field.Pos(), "error checking field for min properties: %v", err)
		return
	}

	if a.omitEmptyPolicy == config.OptionalFieldsOmitEmptyPolicyIgnore && !jsonTags.OmitEmpty {
		a.checkFieldPointersPreferenceWhenRequiredStructTypeWithoutOmitEmpty(pass, field, fieldName, isStarExpr, hasMinimumProperties, fieldHasMinimumProperties, markersAccess)
	} else {
		a.checkFieldPointersPreferenceWhenRequiredStructTypeWithOmitEmpty(pass, field, fieldName, isStarExpr, hasRequiredFields, hasMinimumProperties || fieldHasMinimumProperties)
	}
}

// checkFieldPointersPreferenceWhenRequiredStructTypeWithOmitEmpty recommends adding/removing pointers based on whether the struct
// has any required fields or minimum properties present.
func (a *analyzer) checkFieldPointersPreferenceWhenRequiredStructTypeWithOmitEmpty(pass *analysis.Pass, field *ast.Field, fieldName string, isStarExpr, hasRequiredFields, hasMinimumProperties bool) {
	switch {
	case hasRequiredFields && !isStarExpr:
		reportShouldAddPointer(pass, field, a.pointerPolicy, fieldName, "field %s is optional, but contains required field(s) and should be a pointer")
	case hasMinimumProperties && !isStarExpr:
		reportShouldAddPointer(pass, field, a.pointerPolicy, fieldName, "field %s has a greater than zero minimum number of properties and should be a pointer")
	case isStarExpr && !hasRequiredFields && !hasMinimumProperties:
		reportShouldRemovePointer(pass, field, a.pointerPolicy, fieldName, "field %s is optional, and contains no required field(s) and does not need to be a pointer")
	}
}

// checkFieldPointersPreferenceWhenRequiredStructTypeWithoutOmitEmpty recommends adding/removing pointers based on whether the struct
// has any required fields or minimum properties present, where the struct does not have omitempty as a tag.
func (a *analyzer) checkFieldPointersPreferenceWhenRequiredStructTypeWithoutOmitEmpty(pass *analysis.Pass, field *ast.Field, fieldName string, isStarExpr, hasMinimumProperties, fieldHasMinimumProperties bool, markersAccess markershelper.Markers) {
	switch {
	case hasMinimumProperties && isStarExpr, fieldHasMinimumProperties && isStarExpr:
		// The field is already a pointer and should be a pointer, so we don't need to do anything.
	case hasMinimumProperties && !isStarExpr:
		reportShouldAddPointer(pass, field, a.pointerPolicy, fieldName, "field %s has a greater than zero minimum number of properties and should be a pointer")
	case fieldHasMinimumProperties && !isStarExpr:
		reportShouldRemoveAllInstancesOfIntegerMarker(pass, field, markersAccess, minPropertiesMarker, fieldName, "field %s has a greater than zero minimum number of properties without omitempty. The minimum number of properties should be removed.")
	case isStarExpr:
		// The field is a pointer and should not be a pointer, so we need to remove the pointer.
		reportShouldRemovePointer(pass, field, a.pointerPolicy, fieldName, "field %s is an optional struct without omitempty. It should not be a pointer")
	}
}

// checkFieldPointersPreferenceWhenRequiredString checks string fields for the minimum length marker.
// Where the minimum allowable length is 0, the field should be a pointer.
// Where the minimum length is greater than 0, the field should not be a pointer.
// When the field does not have omitempty, it should not be a pointer, and should not have a minimum length marker.
func (a *analyzer) checkFieldPointersPreferenceWhenRequiredString(pass *analysis.Pass, field *ast.Field, fieldName string, isStarExpr bool, markersAccess markershelper.Markers, jsonTags extractjsontags.FieldTagInfo) {
	if a.omitEmptyPolicy == config.OptionalFieldsOmitEmptyPolicyIgnore && !jsonTags.OmitEmpty {
		a.checkFieldPointersPreferenceWhenRequiredStringWithoutOmitEmpty(pass, field, fieldName, isStarExpr, markersAccess)
		return
	}

	fieldMarkers := markersAccess.FieldMarkers(field)
	if !fieldMarkers.Has(minLengthMarker) {
		if isStarExpr {
			pass.Reportf(field.Pos(), "field %s is an optional string and does not have a minimum length. Where the difference between omitted and the empty string is significant, set the minmum length to 0", fieldName)
		} else {
			pass.Reportf(field.Pos(), "field %s is an optional string and does not have a minimum length. Either set a minimum length or make %s a pointer where the difference between omitted and the empty string is significant", fieldName, fieldName)
		}

		return
	}

	if fieldMarkers.HasWithValue(minLengthMarker + "=0") {
		if isStarExpr {
			// With a minimum length of 0, the field should be a pointer.
		} else {
			reportShouldAddPointer(pass, field, a.pointerPolicy, fieldName, "field %s has a minimum length of 0. The empty string is a valid value and therefore the field should be a pointer")
		}

		return
	}

	// The field has a non-zero (assumed to be greater than zero) minimum length, so it doesn't need to be a pointer.
	if isStarExpr {
		reportShouldRemovePointer(pass, field, a.pointerPolicy, fieldName, "field %s has a greater than 0 length and does not need to be a pointer")
	}
}

// checkFieldPointersPreferenceWhenRequiredStringWithoutOmitEmpty checks string fields for minimum length markers
// and pointers and suggests that both are removed.
func (a *analyzer) checkFieldPointersPreferenceWhenRequiredStringWithoutOmitEmpty(pass *analysis.Pass, field *ast.Field, fieldName string, isStarExpr bool, markersAccess markershelper.Markers) {
	reportShouldRemoveAllInstancesOfIntegerMarker(pass, field, markersAccess, minLengthMarker, fieldName, "field %s has a greater than zero minimum length without omitempty. The minimum length should be removed.")

	// When non-omitempty, the string field should not be a pointer.
	// The empty string should be a valid/acceptable value.
	if isStarExpr {
		reportShouldRemovePointer(pass, field, a.pointerPolicy, fieldName, "field %s is an optional string without omitempty. It should not be a pointer.")
	}
}

func (a *analyzer) checkFieldPointersPreferenceWhenRequiredBool(pass *analysis.Pass, field *ast.Field, fieldName string, isStarExpr bool, jsonTags extractjsontags.FieldTagInfo) {
	if a.omitEmptyPolicy == config.OptionalFieldsOmitEmptyPolicyIgnore && !jsonTags.OmitEmpty {
		if isStarExpr {
			reportShouldRemovePointer(pass, field, a.pointerPolicy, fieldName, "field %s is an optional boolean without omitempty. It should not be a pointer")
		}

		return
	}

	// Optional bools should always be pointers.
	// When the bool is not a pointer, setting the value to false won't round trip.
	// You wouldn't create a bool where the only valid value is true or omitted.
	// This could be confusing for users.
	if !isStarExpr {
		reportShouldAddPointer(pass, field, a.pointerPolicy, fieldName, "field %s is an optional boolean and should be a pointer")
	}
}

//nolint:dupl
func (a *analyzer) checkFieldPointersPreferenceWhenRequiredInteger(pass *analysis.Pass, field *ast.Field, fieldName string, isStarExpr bool, markersAccess markershelper.Markers, jsonTags extractjsontags.FieldTagInfo) {
	fieldMarkers := markersAccess.FieldMarkers(field)

	minValue, err := getMarkerIntegerValueByName(fieldMarkers, minimumMarker)
	if err != nil && !errors.Is(err, errMarkerMissingValue) {
		pass.Reportf(field.Pos(), "field %s has a minimum value of %s, but it is not an integer", fieldName, fieldMarkers.Get(minimumMarker)[0].Expressions[""])
		return
	}

	maxValue, err := getMarkerIntegerValueByName(fieldMarkers, maximumMarker)
	if err != nil && !errors.Is(err, errMarkerMissingValue) {
		pass.Reportf(field.Pos(), "field %s has a maximum value of %s, but it is not an integer", fieldName, fieldMarkers.Get(maximumMarker)[0].Expressions[""])
		return
	}

	if a.omitEmptyPolicy == config.OptionalFieldsOmitEmptyPolicyIgnore && !jsonTags.OmitEmpty {
		a.checkFieldPointersPreferenceWhenRequiredIntegerWithoutOmitEmpty(pass, field, fieldName, minValue, maxValue, isStarExpr, fieldMarkers)
		return
	}

	a.checkFieldPointersPreferenceWhenRequiredIntegerWithOmitEmpty(pass, field, fieldName, isStarExpr, minValue, maxValue)
}

// checkFieldPointersPreferenceWhenRequiredIntegerWithOmitEmpty checks integers based on their minimum and maximum values
// to determine whether or not 0 is a valid value for the integer.
// Where 0 is a valid value, the field should be a pointer.
// Where the limits are ambiguous or missing, the linter will suggest adding a minimum/maximum value to help it decide.
//
//nolint:cyclop
func (a *analyzer) checkFieldPointersPreferenceWhenRequiredIntegerWithOmitEmpty(pass *analysis.Pass, field *ast.Field, fieldName string, isStarExpr bool, minValue, maxValue *int) {
	switch {
	case ptr.Deref(minValue, 0) > 0:
		if isStarExpr {
			reportShouldRemovePointer(pass, field, a.pointerPolicy, fieldName, "field %s has a greater than 0 minimum value and does not need to be a pointer")
		}
	case ptr.Deref(maxValue, 0) < 0:
		if isStarExpr {
			reportShouldRemovePointer(pass, field, a.pointerPolicy, fieldName, "field %s has a negative maximum value and does not need to be a pointer")
		}
	case integerRangeIncludesZero(minValue, maxValue):
		if !isStarExpr {
			reportShouldAddPointer(pass, field, a.pointerPolicy, fieldName, "field %s has a range of values including 0. The difference between omitted and 0 is significant and therefore the field should be a pointer")
		}
	case ptr.Deref(minValue, 0) < 0 && maxValue == nil:
		pass.Reportf(field.Pos(), "field %s has a negative minimum value and does not have a maximum value. A maximum value should be set", fieldName)
	case ptr.Deref(maxValue, 0) > 0 && minValue == nil:
		pass.Reportf(field.Pos(), "field %s has a positive maximum value and does not have a minimum value. A minimum value should be set", fieldName)
	case minValue == nil || maxValue == nil:
		if isStarExpr {
			pass.Reportf(field.Pos(), "field %s is an optional integer and does not have a minimum/maximum value. Where the difference between omitted and 0 is significant, set the minimum/maximum value to a range including 0", fieldName)
		} else {
			pass.Reportf(field.Pos(), "field %s is an optional integer and does not have a minimum/maximum value. Either set a minimum/maximum value or make %s a pointer where the difference between omitted and 0 is significant", fieldName, fieldName)
		}
	}
}

// checkFieldPointersPreferenceWhenRequiredIntegerWithoutOmitEmpty checks integers based on their minimum and maximum values
// to determine whether or not 0 is a valid value for the integer.
// Where 0 is not a valid value, the field should either add omitempty, or remove the limits.
// We assume since there's no omitempty that the API author wants the zero value to be marshalled and as such, suggest to remove the limits.
func (a *analyzer) checkFieldPointersPreferenceWhenRequiredIntegerWithoutOmitEmpty(pass *analysis.Pass, field *ast.Field, fieldName string, minValue, maxValue *int, isStarExpr bool, fieldMarkers markershelper.MarkerSet) {
	switch {
	case minValue != nil && *minValue > 0:
		reportShouldRemoveMarker(pass, field, fieldMarkers.Get(minimumMarker)[0], fieldName, "field %s has a greater than zero minimum value without omitempty. The minimum value should be removed.")
	case maxValue != nil && *maxValue < 0:
		reportShouldRemoveMarker(pass, field, fieldMarkers.Get(maximumMarker)[0], fieldName, "field %s has a less than zero maximum value without omitempty. The maximum value should be removed.")
	}

	if isStarExpr {
		reportShouldRemovePointer(pass, field, a.pointerPolicy, fieldName, "field %s is an optional integer without omitempty. It should not be a pointer.")
	}
}

//nolint:dupl
func (a *analyzer) checkFieldPointersPreferenceWhenRequiredFloat(pass *analysis.Pass, field *ast.Field, fieldName string, isStarExpr bool, markersAccess markershelper.Markers, jsonTags extractjsontags.FieldTagInfo) {
	fieldMarkers := markersAccess.FieldMarkers(field)

	minValue, err := getMarkerFloatValueByName(fieldMarkers, minimumMarker)
	if err != nil && !errors.Is(err, errMarkerMissingValue) {
		pass.Reportf(field.Pos(), "field %s has a minimum value of %s, but it is not a float", fieldName, fieldMarkers.Get(minimumMarker)[0].Expressions[""])
		return
	}

	maxValue, err := getMarkerFloatValueByName(fieldMarkers, maximumMarker)
	if err != nil && !errors.Is(err, errMarkerMissingValue) {
		pass.Reportf(field.Pos(), "field %s has a maximum value of %s, but it is not a float", fieldName, fieldMarkers.Get(maximumMarker)[0].Expressions[""])
		return
	}

	if a.omitEmptyPolicy == config.OptionalFieldsOmitEmptyPolicyIgnore && !jsonTags.OmitEmpty {
		a.checkFieldPointersPreferenceWhenRequiredFloatWithoutOmitEmpty(pass, field, fieldName, minValue, maxValue, isStarExpr, fieldMarkers)
		return
	}

	a.checkFieldPointersPreferenceWhenRequiredFloatWithOmitEmpty(pass, field, fieldName, isStarExpr, minValue, maxValue)
}

// checkFieldPointersPreferenceWhenRequiredFloatWithOmitEmpty checks floats based on their minimum and maximum values
// to determine whether or not 0 is a valid value for the float.
// Where 0 is a valid value, the field should be a pointer.
// Where the limits are ambiguous or missing, the linter will suggest adding a minimum/maximum value to help it decide.
//
//nolint:cyclop
func (a *analyzer) checkFieldPointersPreferenceWhenRequiredFloatWithOmitEmpty(pass *analysis.Pass, field *ast.Field, fieldName string, isStarExpr bool, minValue, maxValue *float64) {
	switch {
	case ptr.Deref(minValue, 0) > 0:
		if isStarExpr {
			reportShouldRemovePointer(pass, field, a.pointerPolicy, fieldName, "field %s has a greater than 0 minimum value and does not need to be a pointer")
		}
	case ptr.Deref(maxValue, 0) < 0:
		if isStarExpr {
			reportShouldRemovePointer(pass, field, a.pointerPolicy, fieldName, "field %s has a negative maximum value and does not need to be a pointer")
		}
	case floatRangeIncludesZero(minValue, maxValue):
		if !isStarExpr {
			reportShouldAddPointer(pass, field, a.pointerPolicy, fieldName, "field %s has a range of values including 0. The difference between omitted and 0 is significant and therefore the field should be a pointer")
		}
	case ptr.Deref(minValue, 0) < 0 && maxValue == nil:
		pass.Reportf(field.Pos(), "field %s has a negative minimum value and does not have a maximum value. A maximum value should be set", fieldName)
	case ptr.Deref(maxValue, 0) > 0 && minValue == nil:
		pass.Reportf(field.Pos(), "field %s has a positive maximum value and does not have a minimum value. A minimum value should be set", fieldName)
	case minValue == nil || maxValue == nil:
		if isStarExpr {
			pass.Reportf(field.Pos(), "field %s is an optional float and does not have a minimum/maximum value. Where the difference between omitted and 0 is significant, set the minimum/maximum value to a range including 0", fieldName)
		} else {
			pass.Reportf(field.Pos(), "field %s is an optional float and does not have a minimum/maximum value. Either set a minimum/maximum value or make %s a pointer where the difference between omitted and 0 is significant", fieldName, fieldName)
		}
	}
}

// checkFieldPointersPreferenceWhenRequiredFloatWithoutOmitEmpty checks floats based on their minimum and maximum values
// to determine whether or not 0 is a valid value for the float.
// Where 0 is not a valid value, the field should either add omitempty, or remove the limits.
// We assume since there's no omitempty that the API author wants the zero value to be marshalled and as such, suggest to remove the limits.
func (a *analyzer) checkFieldPointersPreferenceWhenRequiredFloatWithoutOmitEmpty(pass *analysis.Pass, field *ast.Field, fieldName string, minValue, maxValue *float64, isStarExpr bool, fieldMarkers markershelper.MarkerSet) {
	switch {
	case minValue != nil && *minValue > 0:
		reportShouldRemoveMarker(pass, field, fieldMarkers.Get(minimumMarker)[0], fieldName, "field %s has a greater than zero minimum value without omitempty. The minimum value should be removed.")
	case maxValue != nil && *maxValue < 0:
		reportShouldRemoveMarker(pass, field, fieldMarkers.Get(maximumMarker)[0], fieldName, "field %s has a less than zero maximum value without omitempty. The maximum value should be removed.")
	}

	if isStarExpr {
		reportShouldRemovePointer(pass, field, a.pointerPolicy, fieldName, "field %s is an optional float without omitempty. It should not be a pointer.")
	}
}

func defaultConfig(cfg *config.OptionalFieldsConfig) {
	if cfg.Pointers.Policy == "" {
		cfg.Pointers.Policy = config.OptionalFieldsPointerPolicySuggestFix
	}

	if cfg.Pointers.Preference == "" {
		cfg.Pointers.Preference = config.OptionalFieldsPointerPreferenceAlways
	}

	if cfg.OmitEmpty.Policy == "" {
		cfg.OmitEmpty.Policy = config.OptionalFieldsOmitEmptyPolicySuggestFix
	}
}
