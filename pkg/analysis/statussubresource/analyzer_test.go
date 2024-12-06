package statussubresource_test

import (
	"testing"

	"github.com/JoelSpeed/kal/pkg/analysis/statussubresource"
	"github.com/JoelSpeed/kal/pkg/config"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestStatusSubresourceAnalyzerWithDefaultConfig(t *testing.T) {
	testdata := analysistest.TestData()
	initializer := statussubresource.Initializer()

	analyzer, err := initializer.Init(config.LintersConfig{})
	if err != nil {
		t.Fatal(err)
	}

	analysistest.RunWithSuggestedFixes(t, testdata, analyzer, "a")
}
