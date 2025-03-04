package a

import "a/pkg"

type CommentStartTestStruct struct {
	NoJSONTag     string
	EmptyJSONTag  string `json:""`
	InlineJSONTag string `json:",inline"`
	NoComment     string `json:"noComment"` // want "field NoComment is missing godoc comment"
	Ignored       string `json:"-"`
	Hyphen        string `json:"-,"` // want "field Hyphen is missing godoc comment"

	AnonymousStruct struct { // want "field AnonymousStruct is missing godoc comment"
		NoComment string `json:"noComment"` // want "field NoComment is missing godoc comment"
	} `json:"anonymousStruct"`

	AnonymousStructInlineJSONTag struct {
		NoComment string `json:"noComment"` // want "field NoComment is missing godoc comment"
	} `json:",inline"`

	IgnoredAnonymousStruct struct {
		NoComment string `json:"noComment"`
	} `json:"-"`

	StructForInlineField `json:",inline"`

	A `json:"a"` // want "field A is missing godoc comment"

	PkgA pkg.A `json:"pkgA"` // want "field PkgA is missing godoc comment"

	pkg.Embedded `json:"embedded"` // want "field pkg.Embedded is missing godoc comment"

	*pkg.EmbeddedPointer `json:"embeddedPointer"` // want "field \\*pkg.EmbeddedPointer is missing godoc comment"

	// IncorrectStartComment is a field with an incorrect start to the comment. // want "godoc for field IncorrectStartComment should start with 'incorrectStartComment ...'"
	IncorrectStartComment string `json:"incorrectStartComment"`

	// IncorrectStartOptionalComment is a field with an incorrect start to the comment. // want "godoc for field IncorrectStartOptionalComment should start with 'incorrectStartOptionalComment ...'"
	IncorrectStartOptionalComment string `json:"incorrectStartOptionalComment"`

	// correctStartComment is a field with a correct start to the comment.
	CorrectStartComment string `json:"correctStartComment"`

	// correctStartOptionalComment is a field with a correct start to the comment.
	CorrectStartOptionalComment string `json:"correctStartOptionalComment,omitempty"`

	// IncorrectMultiLineComment is a field with an incorrect start to the comment. // want "godoc for field IncorrectMultiLineComment should start with 'incorrectMultiLineComment ...'"
	// Except this time there are multiple lines to the comment.
	IncorrectMultiLineComment string `json:"incorrectMultiLineComment"`

	// correctMultiLineComment is a field with a correct start to the comment.
	// Except this time there are multiple lines to the comment.
	CorrectMultiLineComment string `json:"correctMultiLineComment"`

	// This comment just isn't correct at all, doesn't even start with anything resembling the field names. // want "godoc for field IncorrectComment should start with 'incorrectComment ...'"
	IncorrectComment string `json:"incorrectComment"`
}

// DoNothing is used to check that the analyser doesn't report on methods.
func (CommentStartTestStruct) DoNothing() {}

type StructForInlineField struct {
	NoComment string `json:"noComment"` // want "field NoComment is missing godoc comment"
}

type A struct {
	NoComment string `json:"noComment"` // want "field NoComment is missing godoc comment"
}

type unexportedStruct struct {
	NoComment string `json:"noComment"` // want "field NoComment is missing godoc comment"
}

type (
	MultipleTypeDeclaration1 struct {
		NoComment string `json:"noComment"` // want "field NoComment is missing godoc comment"
	}
	MultipleTypeDeclaration2 struct {
		NoComment string `json:"noComment"` // want "field NoComment is missing godoc comment"
	}
)

func FunctionWithStructs() {
	type InaccessibleStruct struct {
		NoComment string `json:"noComment"`
	}
}

type Interface interface {
	InaccessibleFunction() string
}
