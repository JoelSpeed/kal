package markers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractMarkerIdAndExpressions(t *testing.T) {
	type testcase struct {
		marker              string
		expectedId          string
		expectedExpressions map[string]string
	}

	testcases := []testcase{
		{
			marker:     "kubebuilder:object:root=true",
			expectedId: "kubebuilder:object:root",
			expectedExpressions: map[string]string{
				"": "true",
			},
		},
		{
			marker:     "kubebuilder:object:root:=true",
			expectedId: "kubebuilder:object:root",
			expectedExpressions: map[string]string{
				"": "true",
			},
		},
		{
			marker:              "required",
			expectedId:          "required",
			expectedExpressions: map[string]string{},
		},
		{
			marker:     "kubebuilder:validation:XValidation:rule='has(self.field)',message='must have field!'",
			expectedId: "kubebuilder:validation:XValidation",
			expectedExpressions: map[string]string{
				"rule":    "'has(self.field)'",
				"message": "'must have field!'",
			},
		},
	}

	for _, tc := range testcases {
		id, expressions := extractMarkerIdAndExpressions(tc.marker)
		assert.Equal(t, tc.expectedId, id, "marker", tc.marker)
		assert.Equal(t, tc.expectedExpressions, expressions, "marker", tc.marker)
	}
}
