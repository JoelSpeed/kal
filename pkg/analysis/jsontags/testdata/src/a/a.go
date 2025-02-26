package a

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type JSONTagTestStruct struct {
	NoJSONTag            string // want "field NoJSONTag is missing json tag"
	EmptyJSONTag         string `json:""`                        // want "field EmptyJSONTag has empty json tag"
	NonCamelCaseJSONTag  string `json:"non_camel_case_json_tag"` // want "field NonCamelCaseJSONTag json tag does not match pattern \"\\^\\[a-z\\]\\[a-z0-9\\]\\*\\(\\?:\\[A-Z\\]\\[a-z0-9\\]\\*\\)\\*\\$\": non_camel_case_json_tag"
	WithHyphensJSONTag   string `json:"with-hyphens-json-tag"`   // want "field WithHyphensJSONTag json tag does not match pattern \"\\^\\[a-z\\]\\[a-z0-9\\]\\*\\(\\?:\\[A-Z\\]\\[a-z0-9\\]\\*\\)\\*\\$\": with-hyphens-json-tag"
	PascalCaseJSONTag    string `json:"PascalCaseJSONTag"`       // want "field PascalCaseJSONTag json tag does not match pattern \"\\^\\[a-z\\]\\[a-z0-9\\]\\*\\(\\?:\\[A-Z\\]\\[a-z0-9\\]\\*\\)\\*\\$\": PascalCaseJSONTag"
	NonTerminatedJSONTag string `json:"nonTerminatedJSONTag`     // want "field NonTerminatedJSONTag is missing json tag"
	XMLTag               string `xml:"xmlTag"`                   // want "field XMLTag is missing json tag"
	InlineJSONTag        string `json:",inline"`
	ValidJSONTag         string `json:"validJsonTag"`
	ValidOptionalJSONTag string `json:"validOptionalJsonTag,omitempty"`
	JSONTagWithID        string `json:"jsonTagWithID"`
	JSONTagWithTTL       string `json:"jsonTagWithTTL"`
	JSONTagWithGiB       string `json:"jsonTagWithGiB"`
	Ignored              string `json:"-"`

	IgnoredAnonymousStruct struct {
		A string `json:""` // want "field A has empty json tag"
	} `json:"-"`

	A `json:",inline"`
	B `json:"bar,omitempty"`
	C             // want "embedded field C is missing json tag"
	D `json:""`   // want "embedded field D has empty json tag"
	E `json:"e-"` // want "embedded field E json tag does not match pattern \"\\^\\[a-z\\]\\[a-z0-9\\]\\*\\(\\?:\\[A-Z\\]\\[a-z0-9\\]\\*\\)\\*\\$\": e-"

	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
}

type A struct{}

func (A) DoNothing() {}

type B struct{}

type C struct{}

type D struct{}

type E struct{}
