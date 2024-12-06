# Adding a new linter to KAL

Linters in KAL should live in their own package within the `pkg/analysis` directory.

Each linter is based on the `analysis.Analyzer` interface from the [golang.org/x/tools/go/analysis][go-analysis] package.

[go-analysis]: https://pkg.go.dev/golang.org/x/tools/go/analysis#hdr-Analyzer

The core of the linter is the `run` function, implemented with the signature:

```go
func (a *analyzer) run(pass *analysis.Pass) (interface{}, error)
```

It is recommended to implement the linter as a struct, that can contain configuration, and have the `run` function as a method on the struct.

It is also recommended to use the `inspect.Analyzer` pattern, which allows filtering the parsed syntax tree down to the types of nodes that are relevant to the linter.
This automates a lot of pre-work, and can be seen across existing linters, e.g. `jsontags` or `commentstart`.

Once you are within the `inspect.Preorder`, you can then implement the business logic of the linter, focusing on details of structs, fields, comments, etc.

## Registry

The registry in the analysis package co-ordinates the initialization of all linters.
Where linters have configuration, or are enabled/disabled by higher level configuration, the registry takes on making sure the linters are initialized correctly.

To enable the registry, each linter package must create an `Initializer` function that returns an `analysis.AnalyzerInitializer` interface (from `pkg/analysis`).

It is expected that each linter package contain a file `initializer.go`, the content of this file should be as follows:

```go
// Initializer returns the AnalyzerInitializer for this
// Analyzer so that it can be added to the registry.
func Initializer() initializer {
	return initializer{}
}

// intializer implements the AnalyzerInitializer interface.
type initializer struct{}

// Name returns the name of the Analyzer.
func (initializer) Name() string {
	return name
}

// Init returns the intialized Analyzer.
func (initializer) Init(cfg config.LintersConfig) (*analysis.Analyzer, error) {
	return newAnalyzer(cfg.MyLinterConfig)
}

// Default determines whether this Analyzer is on by default, or not.
func (initializer) Default() bool {
	return true // or false
}
```

This pattern allows the linter to be registered with the KAL registry, and allows the linter to be initialized with configuration.

Once you have created the `initializer.go` file, you will need to add the linter to the `pkg/analysis/registry.go` file.

Add the initializer to the `NewRegistry` function, and it will then be included in the linter builds.

```go
func NewRegistry() []*analysis.Analyzer {
    return []*analysis.Analyzer{
        // Add the new linter here
        mynewlinter.Initializer(),
    }
}
```

## Configuration

If the linter requires configuration, the configuration should be added to the `config` package.

Add a new structure (or structures) to the `linters_config.go` file.
Include a new field for the top level configuration to the `LintersConfig` struct, using the name of the linter, in camel case for the `json` tag.

Any options for the linter, should also be validated.
Validation lives in the `validation` package.

Within the `ValidateLintersConfig` function, in the `linters_config.go` file, you will need to add
a line as below, to include any configuration validation for the new linter.
There are already examples of this in the file.

```go
fieldErrors = append(fieldErrors, validateMyNewLint(lc.MyNewLinter, fldPath.Child("myNewLinter"))...)
```

The validations should use the `field.Error` pattern to provide consistent error messages.

## Helpers

The helpers package contains `analysis.Analyzer` implementations that can be used to source the common functionality required by linters.
For example, extracting information about `json` tags, or extracting `// +` style markers from comments.

Any new, common functionality should also be added to the helpers package.

Importantly:
* Helpers should expose a public `Analyzer` variable, of type `*analysis.Analyzer`.
* Helpers may not depend on any linter that needs to be initialized with configuration.
* Helpers themselves should not report lint issues, but should provide information to the linters that do.

In general, helpers return interfaces that can expose useful information in a simple way.
Exposing structs or maps directly as the result type of the helper means common
functions for accessing data must be implemented in each linter.

To use a helper, the helper `Analyzer` should be included in the linter `Analyzer`'s `Require` property.

```go
    return &analysis.Analyzer{
    Name:     "linterName",
    Doc:      "linter description",
    Requires: []*analysis.Analyzer{helpers.Analyzer},
    Run:      l.run,
}
```

Within the `run` function, the result of the helper can be extracted from the `*analysis.Pass` object.

```go
func (l *linter) run(pass *analysis.Pass) (interface{}, error) {
    helperResult := pass.ResultOf[helpers.Analyzer].(*helpers.ResultType)
    ...
}
```

## Tests

Basic tests can be implemented with the Go analysis test framework.

Create a `testdata` directory in the linter package and create a structure underneath.
Individual test files must be placed under `src` and then a subdirectory for each package.

If your linter has different configurations, e.g options to pass to the linter
you will need one configuration per option. Use one package per configuration
input for the linter.

```
mylinter
-- mylinter.go
-- mylinter_test.go
-- testdata
  -- src
    -- a
      -- a.go
    -- b
      -- b.go
```

The test suite can then be written using the standard go test framework.

```go
func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, commentstart.Analyzer, "a")
}
```

Each file within the test package should contain Go code to test the linter.
Typically this would mean a combination of constants, type declarations and struct declarations.

Where the linter is expected to return an issue, a comment can be added to the file to indicate the expected issue.

```go
type Foo struct {
    // this comment should be flagged by the linter // Want 'comment should start with a capital letter'
    Bar string

    foo string // want 'field is not exported'
}
```

If the expected output of the test happens to contain a regex string, then the regex within the `want` comment should be escaped.
The `jsontags` linter has an example of this pattern, which can be referred to.

### With suggested fixes

Where a linter also implements suggested fixes, the test suite can be extended to include the suggested fixes.

Replace `analysistest.Run` with `analysistest.RunWithSuggestedFixes`.

```go
func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, commentstart.Analyzer, "a")
}
```

For each file in the test packages, a corresponding `.golden` file should be created with the expected output, once the fixes have been applied.

The `commentstart` linter has an example of this pattern, which can be referred to.

## Docs

Each linter package should contain a `doc.go` file, which has a package level comment explaining what the linter is,
what it checks for, and how it can be configured, if appropriate.

The package level documentation is helpful when running `godoc` or accessing `pkg.go.dev`.

### ReadMe

The root `README.md` file should be updated to include the new linter and any relevant information about the linter.

The format for which should look something like:

````markdown
## LinterNameInPascalCase

Include a description of what the linter is checking for.

### Fixes (via standalone binary)

Include a description of what the linter fixes, if applicable.

### Configuration

Include an example of the configuration that can be used to configure the linter, if applicable.

```yaml
lintersConfig:
  linterNameInCamelCase:
    option: value
```
````

Please add the new linter in the correct position alphabetically.
