package commentstart_test

import (
	"testing"

	"github.com/JoelSpeed/kal/pkg/analysis/commentstart"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, commentstart.Analyzer, "a/...")
}
