package conditions_test

import (
	"testing"

	"github.com/JoelSpeed/kal/pkg/analysis/conditions"
	"github.com/JoelSpeed/kal/pkg/config"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestDefaultConfiguration(t *testing.T) {
	testdata := analysistest.TestData()

	a, err := conditions.Initializer().Init(config.LintersConfig{})
	if err != nil {
		t.Fatal(err)
	}

	analysistest.RunWithSuggestedFixes(t, testdata, a, "a")
}

func TestNotFieldFirst(t *testing.T) {
	testdata := analysistest.TestData()

	a, err := conditions.Initializer().Init(config.LintersConfig{
		Conditions: config.ConditionsConfig{
			IsFirstField: config.ConditionsFirstFieldIgnore,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	analysistest.RunWithSuggestedFixes(t, testdata, a, "b")
}

func TestIgnoreProtobuf(t *testing.T) {
	testdata := analysistest.TestData()

	a, err := conditions.Initializer().Init(config.LintersConfig{
		Conditions: config.ConditionsConfig{
			UseProtobuf: config.ConditionsUseProtobufIgnore,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	analysistest.RunWithSuggestedFixes(t, testdata, a, "c")
}
