package nomaps

import (
	"errors"
	"fmt"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const name = "nomaps"

var (
	errCouldNotGetInspector = errors.New("could not get inspector")
)

// Analyzer is the analyzer for the nomaps package.
// It checks that the struct fields which is map type.
// if the type is map[string]string, it should be skipped since the type is used by Labels/Annotations.
var Analyzer = &analysis.Analyzer{
	Name:     name,
	Doc:      "no map type should be allowed and prefer list type",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, errCouldNotGetInspector
	}

	// Filter to fields so that we can iterate over fields in a struct.
	nodeFilter := []ast.Node{
		(*ast.StructType)(nil),
	}

	ignoreMapType := types.NewMap(types.Typ[types.String], types.Typ[types.String])

	// Preorder visits all the nodes of the AST in depth-first order. It calls
	// f(n) for each node n before it visits n's children.
	//
	// We use the filter defined above, ensuring we only look at struct fields.
	inspect.Preorder(nodeFilter, func(n ast.Node) {
		structType, ok := n.(*ast.StructType)
		if !ok {
			return
		}

		if structType.Fields == nil {
			return
		}

		for _, field := range structType.Fields.List {
			var m *types.Map
			underlyingType := pass.TypesInfo.TypeOf(field.Type).Underlying()
			if ptr, ok := underlyingType.(*types.Pointer); ok {
				underlyingType = ptr.Elem().Underlying()
			}
			m, ok := underlyingType.(*types.Map)
			if !ok {
				continue
			}
			if types.Identical(m, ignoreMapType) {
				continue
			}
			pass.Report(analysis.Diagnostic{
				Pos:     field.Pos(),
				Message: fmt.Sprintf("%s should not use a map type, use list type instead", field.Names[0].Name),
			})
		}
	})

	return nil, nil //nolint:nilnil
}
