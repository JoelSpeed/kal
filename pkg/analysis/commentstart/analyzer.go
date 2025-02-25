package commentstart

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/JoelSpeed/kal/pkg/analysis/helpers/extractjsontags"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const name = "commentstart"

var (
	errCouldNotGetInspector = errors.New("could not get inspector")
	errCouldNotGetJSONTags  = errors.New("could not get json tags")
)

// Analyzer is the analyzer for the commentstart package.
// It checks that all struct fields in an API have a godoc, and that the godoc starts with the serialised field name.
var Analyzer = &analysis.Analyzer{
	Name:     name,
	Doc:      "Check that all struct fields in an API have a godoc, and that the godoc starts with the serialised field name",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer, extractjsontags.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, errCouldNotGetInspector
	}

	jsonTags, ok := pass.ResultOf[extractjsontags.Analyzer].(extractjsontags.StructFieldTags)
	if !ok {
		return nil, errCouldNotGetJSONTags
	}

	// Filter to structs so that we can iterate over fields in a struct.
	nodeFilter := []ast.Node{
		(*ast.Field)(nil),
	}

	inspect.WithStack(nodeFilter, func(n ast.Node, push bool, stack []ast.Node) (proceed bool) {
		if !push {
			return false
		}

		if len(stack) < 2 {
			return true
		}

		// The 0th node in the stack is the *ast.File.
		// The 1st node in the stack is the *ast.GenDecl.
		decl, ok := stack[1].(*ast.GenDecl)
		if !ok {
			return false
		}

		if decl.Tok != token.TYPE {
			// Returning false here means we won't inspect non-type declarations (e.g. var, const, import).
			return false
		}

		field, ok := n.(*ast.Field)
		if !ok {
			return true
		}

		return checkField(pass, field, jsonTags)
	})

	return nil, nil //nolint:nilnil
}

func checkField(pass *analysis.Pass, field *ast.Field, jsonTags extractjsontags.StructFieldTags) (proceed bool) {
	if field == nil || len(field.Names) == 0 {
		// Returning false here means we don't inspect inline fields.
		// Types of inline fields will be inspected on its declaration.
		return false
	}

	tagInfo := jsonTags.FieldTags(field)
	if tagInfo.Ignored {
		// Returning false here means we won't inspect the children of an ignored field.
		return false
	}

	if tagInfo.Name == "" {
		return true
	}

	fieldName := field.Names[0].Name

	if field.Doc == nil {
		pass.Reportf(field.Pos(), "field %s is missing godoc comment", fieldName)
		return true
	}

	firstLine := field.Doc.List[0]
	if !strings.HasPrefix(firstLine.Text, "// "+tagInfo.Name+" ") {
		if strings.HasPrefix(strings.ToLower(firstLine.Text), strings.ToLower("// "+tagInfo.Name+" ")) {
			// The comment start is correct, apart from the casing, we can fix that.
			pass.Report(analysis.Diagnostic{
				Pos:     firstLine.Pos(),
				Message: fmt.Sprintf("godoc for field %s should start with '%s ...'", fieldName, tagInfo.Name),
				SuggestedFixes: []analysis.SuggestedFix{
					{
						Message: fmt.Sprintf("should replace first word with `%s`", tagInfo.Name),
						TextEdits: []analysis.TextEdit{
							{
								Pos:     firstLine.Pos(),
								End:     firstLine.Pos() + token.Pos(len(tagInfo.Name)) + token.Pos(4),
								NewText: []byte("// " + tagInfo.Name + " "),
							},
						},
					},
				},
			})
		} else {
			pass.Reportf(field.Doc.List[0].Pos(), "godoc for field %s should start with '%s ...'", fieldName, tagInfo.Name)
		}
	}

	return true
}
