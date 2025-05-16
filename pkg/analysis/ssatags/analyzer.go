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
package ssatags

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"sigs.k8s.io/kube-api-linter/pkg/analysis/helpers/extractjsontags"
	"sigs.k8s.io/kube-api-linter/pkg/analysis/helpers/inspector"
	"sigs.k8s.io/kube-api-linter/pkg/analysis/helpers/markers"
)

const name = "ssatags"

// Analyzer is the analyzer for the ssatags package.
// It checks whether the array types in the API have the SSA tags and the usage of the tags is correct.
var Analyzer = &analysis.Analyzer{
	Name:     name,
	Doc:      "Check that all array types in the API have the SSA tags and the usage of the tags is correct",
	Run:      run,
	Requires: []*analysis.Analyzer{inspector.Analyzer},
}

func run(pass *analysis.Pass) (any, error) {
	inspect, ok := pass.ResultOf[inspector.Analyzer].(inspector.Inspector)
	if !ok {
		return nil, nil
	}

	inspect.InspectFields(func(field *ast.Field, stack []ast.Node, jsonTagInfo extractjsontags.FieldTagInfo, markersAccess markers.Markers) {
	})

	inspect.InspectTypeSpec(func(typeSpec *ast.TypeSpec, markersAccess markers.Markers) {
	})

	return nil, nil //nolint:nilnil
}
