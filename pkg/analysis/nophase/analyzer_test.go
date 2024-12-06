package nophase_test

import (
	"testing"

	"github.com/JoelSpeed/kal/pkg/analysis/nophase"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, nophase.Analyzer, "a")
}
