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
	"fmt"
	"go/ast"
	"go/token"
	"strconv"

	"golang.org/x/tools/go/analysis"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/kube-api-linter/pkg/analysis/helpers/markers"
	"sigs.k8s.io/kube-api-linter/pkg/config"
)

// isStarExpr checks if the expression is a pointer type.
// If it is, it returns the expression inside the pointer.
func isStarExpr(expr ast.Expr) (bool, ast.Expr) {
	if ptrType, ok := expr.(*ast.StarExpr); ok {
		return true, ptrType.X
	}

	return false, expr
}

// isPointerType checks if the expression is a pointer type.
// This is for types that are always implemented as pointers and therefore should
// not be the underlying type of a star expr.
func isPointerType(expr ast.Expr) bool {
	switch expr.(type) {
	case *ast.StarExpr, *ast.MapType, *ast.ArrayType:
		return true
	default:
		return false
	}
}

// structContainsRequiredFields checks whether the struct has any required fields.
// Having a required field means that `{}` is not a valid entity for the struct,
// and therefore the struct should be a pointer when optional.
func structContainsRequiredFields(structType *ast.StructType, markersAccess markers.Markers) bool {
	if structType == nil {
		return false
	}

	for _, field := range structType.Fields.List {
		fieldMarkers := markersAccess.FieldMarkers(field)

		if isFieldRequired(fieldMarkers) {
			return true
		}
	}

	return false
}

// isFieldRequired checks if a field has a required marker.
func isFieldRequired(fieldMarkers markers.MarkerSet) bool {
	return fieldMarkers.Has(requiredMarker) || fieldMarkers.Has(kubebuilderRequiredMarker)
}

// isFieldOptional checks if a field has an optional marker.
func isFieldOptional(fieldMarkers markers.MarkerSet) bool {
	return fieldMarkers.Has(optionalMarker) || fieldMarkers.Has(kubebuilderOptionalMarker)
}

// reportShouldAddPointer adds an analysis diagnostic that explains that a pointer should be added.
// Where the pointer policy is suggest fix, it also adds a suggested fix to add the pointer.
func reportShouldAddPointer(pass *analysis.Pass, field *ast.Field, pointerPolicy config.OptionalFieldsPointerPolicy, fieldName, messageFmt string) {
	switch pointerPolicy {
	case config.OptionalFieldsPointerPolicySuggestFix:
		pass.Report(analysis.Diagnostic{
			Pos:     field.Pos(),
			Message: fmt.Sprintf(messageFmt, fieldName),
			SuggestedFixes: []analysis.SuggestedFix{
				{
					Message: "should make the field a pointer",
					TextEdits: []analysis.TextEdit{
						{
							Pos:     field.Pos() + token.Pos(len(fieldName)+1),
							NewText: []byte("*"),
						},
					},
				},
			},
		})
	case config.OptionalFieldsPointerPolicyWarn:
		pass.Reportf(field.Pos(), messageFmt, fieldName)
	default:
		panic(fmt.Sprintf("unknown pointer policy: %s", pointerPolicy))
	}
}

// reportShouldRemovePointer adds an analysis diagnostic that explains that a pointer should be removed.
// Where the pointer policy is suggest fix, it also adds a suggested fix to remove the pointer.
func reportShouldRemovePointer(pass *analysis.Pass, field *ast.Field, pointerPolicy config.OptionalFieldsPointerPolicy, fieldName, messageFmt string) {
	switch pointerPolicy {
	case config.OptionalFieldsPointerPolicySuggestFix:
		pass.Report(analysis.Diagnostic{
			Pos:     field.Pos(),
			Message: fmt.Sprintf(messageFmt, fieldName),
			SuggestedFixes: []analysis.SuggestedFix{
				{
					Message: "should remove the pointer",
					TextEdits: []analysis.TextEdit{
						{
							Pos: field.Pos() + token.Pos(len(fieldName)+1),
							End: field.Pos() + token.Pos(len(fieldName)+2),
						},
					},
				},
			},
		})
	case config.OptionalFieldsPointerPolicyWarn:
		pass.Reportf(field.Pos(), messageFmt, fieldName)
	default:
		panic(fmt.Sprintf("unknown pointer policy: %s", pointerPolicy))
	}
}

// reportShouldRemoveAllInstancesOfIntegerMarker adds an analysis diagnostic that explains that a marker should be removed.
// This function is used to find non-zero valued markers, and suggest that they are removed when the field is not a pointer.
// This is used for markers like MinLength and MinProperties.
func reportShouldRemoveAllInstancesOfIntegerMarker(pass *analysis.Pass, field *ast.Field, markersAccess markers.Markers, markerName, fieldName, messageFmt string) {
	fieldMarkers := markersAccess.FieldMarkers(field)

	for _, marker := range fieldMarkers.Get(markerName) {
		markerValue, err := getMarkerIntegerValue(marker)
		if err != nil {
			pass.Reportf(marker.Pos, "invalid value for %s marker: %v", markerName, err)
			return
		}

		if markerValue > 0 {
			reportShouldRemoveMarker(pass, field, marker, fieldName, messageFmt)
		}
	}
}

// reportShouldAddMarker adds an analysis diagnostic that explains that a marker should be removed.
// This is used where we see a marker that would conflict with a field that lacks omitempty.
func reportShouldRemoveMarker(pass *analysis.Pass, field *ast.Field, marker markers.Marker, fieldName, messageFmt string) {
	pass.Report(analysis.Diagnostic{
		Pos:     field.Pos(),
		Message: fmt.Sprintf(messageFmt, fieldName),
		SuggestedFixes: []analysis.SuggestedFix{
			{
				Message: fmt.Sprintf("should remove the marker: %s", marker.RawComment),
				TextEdits: []analysis.TextEdit{
					{
						Pos: marker.Pos,
						End: marker.End + 1,
					},
				},
			},
		},
	})
}

// getMarkerIntegerValueByName extracts the numeric value from the first
// instace of the marker with the given name.
// Works for markers like MaxLength, MinLength, etc.
func getMarkerIntegerValueByName(marker markers.MarkerSet, markerName string) (*int, error) {
	markerList := marker.Get(markerName)
	if len(markerList) == 0 {
		return nil, errMarkerMissingValue
	}

	markerValue, err := getMarkerIntegerValue(markerList[0])
	if err != nil {
		return nil, fmt.Errorf("error getting marker value: %w", err)
	}

	return &markerValue, nil
}

// getMarkerIntegerValue extracts a numeric value from the default
// value of a marker.
// Works for markers like MaxLength, MinLength, etc.
func getMarkerIntegerValue(marker markers.Marker) (int, error) {
	rawValue, ok := marker.Expressions[""]
	if !ok {
		return 0, errMarkerMissingValue
	}

	value, err := strconv.Atoi(rawValue)
	if err != nil {
		return 0, fmt.Errorf("error converting value to integer: %w", err)
	}

	return value, nil
}

// getMarkerFloatValueByName extracts the numeric value from the first
// instace of the marker with the given name.
// Works for markers like MaxLength, MinLength, etc.
func getMarkerFloatValueByName(marker markers.MarkerSet, markerName string) (*float64, error) {
	markerList := marker.Get(markerName)
	if len(markerList) == 0 {
		return nil, errMarkerMissingValue
	}

	markerValue, err := getMarkerFloatValue(markerList[0])
	if err != nil {
		return nil, fmt.Errorf("error getting marker value: %w", err)
	}

	return &markerValue, nil
}

// getMarkerFloatValue extracts a numeric value from the default
// value of a marker.
// Works for markers like MaxLength, MinLength, etc.
func getMarkerFloatValue(marker markers.Marker) (float64, error) {
	rawValue, ok := marker.Expressions[""]
	if !ok {
		return 0, errMarkerMissingValue
	}

	value, err := strconv.ParseFloat(rawValue, 64)
	if err != nil {
		return 0, fmt.Errorf("error converting value to float: %w", err)
	}

	return value, nil
}

// structHasGreaterThanZeroMaxProperties checks if the struct has a minProperties marker.
// It inspects the minProperties marker on the struct type itself.
func structHasGreaterThanZeroMinProperties(structType *ast.StructType, structMarkers markers.MarkerSet) (bool, error) {
	if structType == nil {
		return false, nil
	}

	for _, marker := range structMarkers.Get(minPropertiesMarker) {
		markerValue, err := getMarkerIntegerValue(marker)
		if err != nil {
			return false, fmt.Errorf("error getting marker value: %w", err)
		}

		if markerValue > 0 {
			return true, nil
		}
	}

	return false, nil
}

func integerRangeIncludesZero(minimum, maximum *int) bool {
	return ptr.Deref(minimum, -1) == 0 ||
		ptr.Deref(maximum, -1) == 0 ||
		ptr.Deref(minimum, 0) < 0 && ptr.Deref(maximum, 0) > 0
}

func floatRangeIncludesZero(minimum, maximum *float64) bool {
	return ptr.Deref(minimum, -1) == 0 ||
		ptr.Deref(maximum, -1) == 0 ||
		ptr.Deref(minimum, 0) < 0 && ptr.Deref(maximum, 0) > 0
}
