package extractjsontags

import (
	"errors"
	"go/ast"
	"go/token"
	"reflect"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var (
	errCouldNotGetInspector          = errors.New("could not get inspector")
	errCouldNotCreateStructFieldTags = errors.New("could not create new structFieldTags")
)

// StructFieldTags is used to find information about
// json tags on fields within struct.
type StructFieldTags interface {
	FieldTags(*ast.StructType, *ast.Field) FieldTagInfo
}

type structFieldTags struct {
	structToFieldTags map[*ast.StructType]map[*ast.Field]FieldTagInfo
}

func newStructFieldTags() StructFieldTags {
	return &structFieldTags{
		structToFieldTags: make(map[*ast.StructType]map[*ast.Field]FieldTagInfo),
	}
}

func (s *structFieldTags) insertFieldTagInfo(styp *ast.StructType, field *ast.Field, tagInfo FieldTagInfo) {
	if s.structToFieldTags[styp] == nil {
		s.structToFieldTags[styp] = make(map[*ast.Field]FieldTagInfo)
	}

	s.structToFieldTags[styp][field] = tagInfo
}

// FieldTags find the tag information for the named field within the given struct.
func (s *structFieldTags) FieldTags(styp *ast.StructType, field *ast.Field) FieldTagInfo {
	structFields := s.structToFieldTags[styp]

	if structFields != nil {
		return structFields[field]
	}

	return FieldTagInfo{}
}

// Analyzer is the analyzer for the jsontags package.
// It checks that all struct fields in an API are tagged with json tags.
var Analyzer = &analysis.Analyzer{
	Name:       "extractjsontags",
	Doc:        "Iterates over all fields in structs and extracts their json tags.",
	Run:        run,
	Requires:   []*analysis.Analyzer{inspect.Analyzer},
	ResultType: reflect.TypeOf(newStructFieldTags()),
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, errCouldNotGetInspector
	}

	// Filter to structs so that we can iterate over fields in a struct.
	nodeFilter := []ast.Node{
		(*ast.StructType)(nil),
	}

	results, ok := newStructFieldTags().(*structFieldTags)
	if !ok {
		return nil, errCouldNotCreateStructFieldTags
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		sTyp, ok := n.(*ast.StructType)
		if !ok {
			return
		}

		if sTyp.Fields == nil {
			return
		}

		for i := 0; i < sTyp.Fields.NumFields(); i++ {
			field := sTyp.Fields.List[i]

			results.insertFieldTagInfo(sTyp, field, extractTagInfo(field.Tag))
		}
	})

	return results, nil
}

func extractTagInfo(tag *ast.BasicLit) FieldTagInfo {
	if tag == nil || tag.Value == "" {
		return FieldTagInfo{Missing: true}
	}

	rawTag, err := strconv.Unquote(tag.Value)
	if err != nil {
		// This means the way AST is treating tags has changed.
		panic(err)
	}

	tagValue, ok := reflect.StructTag(rawTag).Lookup("json")
	if !ok {
		return FieldTagInfo{Missing: true}
	}

	if tagValue == "" {
		return FieldTagInfo{}
	}

	pos := tag.Pos() + token.Pos(strings.Index(tag.Value, tagValue))
	end := pos + token.Pos(len(tagValue))

	tagValues := strings.Split(tagValue, ",")

	if len(tagValues) == 2 && tagValues[0] == "" && tagValues[1] == "inline" {
		return FieldTagInfo{
			Inline:   true,
			RawValue: tagValue,
			Pos:      pos,
			End:      end,
		}
	}

	tagName := tagValues[0]

	return FieldTagInfo{
		Name:      tagName,
		OmitEmpty: len(tagValues) == 2 && tagValues[1] == "omitempty",
		RawValue:  tagValue,
		Pos:       pos,
		End:       end,
	}
}

// FieldTagInfo contains information about a field's json tag.
// This is used to pass information about a field's json tag between analyzers.
type FieldTagInfo struct {
	// Name is the name of the field extracted from the json tag.
	Name string

	// OmitEmpty is true if the field has the omitempty option in the json tag.
	OmitEmpty bool

	// Inline is true if the field has the inline option in the json tag.
	Inline bool

	// Missing is true when the field had no json tag.
	Missing bool

	// RawValue is the raw value from the json tag.
	RawValue string

	// Pos marks the starting position of the json tag value.
	Pos token.Pos

	// End marks the end of the json tag value.
	End token.Pos
}
