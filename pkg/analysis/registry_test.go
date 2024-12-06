package analysis_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	goanalysis "golang.org/x/tools/go/analysis"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/JoelSpeed/kal/pkg/analysis"
	"github.com/JoelSpeed/kal/pkg/config"
)

var _ = Describe("Registry", func() {
	Context("DefaultLinters", func() {
		It("should return the default linters", func() {
			r := analysis.NewRegistry()
			Expect(r.DefaultLinters().UnsortedList()).To(ConsistOf(
				"conditions",
				"commentstart",
				"jsontags",
				"optionalorrequired",
				"requiredfields",
			))
		})
	})

	Context("AllLinters", func() {
		It("should return the all known linters", func() {
			r := analysis.NewRegistry()
			Expect(r.AllLinters().UnsortedList()).To(ConsistOf(
				"conditions",
				"commentstart",
				"jsontags",
				"optionalorrequired",
				"requiredfields",
			))
		})
	})

	Context("InitializeLinters", func() {
		type initLintersTableInput struct {
			config        config.Linters
			lintersConfig config.LintersConfig

			expectedLinters []string
		}

		DescribeTable("Initialize Linters", func(in initLintersTableInput) {
			r := analysis.NewRegistry()
			linters, err := r.InitializeLinters(in.config, in.lintersConfig)
			Expect(err).NotTo(HaveOccurred())

			toLinterNames := func(a []*goanalysis.Analyzer) []string {
				names := []string{}

				for _, linter := range a {
					names = append(names, linter.Name)
				}

				return names
			}

			Expect(linters).To(WithTransform(toLinterNames, ConsistOf(in.expectedLinters)))
		},
			Entry("Empty config", initLintersTableInput{
				config:          config.Linters{},
				lintersConfig:   config.LintersConfig{},
				expectedLinters: analysis.NewRegistry().DefaultLinters().UnsortedList(),
			}),
			Entry("With wildcard enabled linters", initLintersTableInput{
				config: config.Linters{
					Enable: []string{config.Wildcard},
				},
				lintersConfig:   config.LintersConfig{},
				expectedLinters: analysis.NewRegistry().AllLinters().UnsortedList(),
			}),
			Entry("With wildcard enabled linters and a disabled linter", initLintersTableInput{
				config: config.Linters{
					Enable:  []string{config.Wildcard},
					Disable: []string{"jsontags"},
				},
				lintersConfig:   config.LintersConfig{},
				expectedLinters: analysis.NewRegistry().AllLinters().Difference(sets.New("jsontags")).UnsortedList(),
			}),
			Entry("With wildcard disabled linters", initLintersTableInput{
				config: config.Linters{
					Disable: []string{config.Wildcard},
				},
				lintersConfig:   config.LintersConfig{},
				expectedLinters: []string{},
			}),
			Entry("With wildcard disabled linters and an enabled linter", initLintersTableInput{
				config: config.Linters{
					Disable: []string{config.Wildcard},
					Enable:  []string{"jsontags"},
				},
				lintersConfig:   config.LintersConfig{},
				expectedLinters: []string{"jsontags"},
			}),
		)
	})
})
