package nomaps_test

import (
	"testing"

	"github.com/JoelSpeed/kal/pkg/analysis/nomaps"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, nomaps.Analyzer, "a")
}
