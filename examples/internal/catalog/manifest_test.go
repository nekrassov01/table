package catalog

import (
	"go/ast"
	"go/token"
	"os"
	"path/filepath"
	"testing"

	"github.com/nekrassov01/table/internal/testutil"
)

const manifestSource = manifestDataNamesSource + manifestExampleDataSource + manifestOptionSource

const manifestDataNamesSource = `package main

var targetText = "text"
var dataSimple = "simple"

func dataNames(target string) ([]string, bool) {
	switch target {
	case targetText:
		return []string{dataSimple}, true
	default:
		return nil, false
	}
}
`

const manifestExampleDataSource = `
func exampleData(name string) (any, bool) {
	switch name {
	case dataSimple:
		return examples.SimpleData, true
	default:
		return nil, false
	}
}
`

const manifestOptionSource = `
type runner struct{ data string }

func (o runner) runText() {
	switch o.data {
	case dataSimple:
		opts = examples.TextOptionSimple
	}
}
`

func Test_parseManifest(t *testing.T) {
	type args struct {
		source string
	}
	type want struct {
		manifest manifest
		err      bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "manifest",
			args: args{
				source: manifestSource,
			},
			want: want{
				manifest: manifest{
					scenarios: []scenario{
						{
							name: "simple",
							data: "SimpleData",
							targets: []target{
								{
									name:       "text",
									option:     "TextOptionSimple",
									optionFile: "examples/option_text.go",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "invalid source",
			args: args{
				source: "package main\nfunc",
			},
			want: want{
				err: true,
			},
		},
		{
			name: "missing target data function",
			args: args{
				source: "package main\nvar targetText = \"text\"",
			},
			want: want{
				err: true,
			},
		},
		{
			name: "missing example data function",
			args: args{
				source: manifestDataNamesSource,
			},
			want: want{
				err: true,
			},
		},
		{
			name: "missing option method",
			args: args{
				source: manifestDataNamesSource + manifestExampleDataSource,
			},
			want: want{
				err: true,
			},
		},
		{
			name: "unresolved data",
			args: args{
				source: manifestDataNamesSource + `
func exampleData(name string) (any, bool) {
	switch name {
	default:
		return nil, false
	}
}
` + manifestOptionSource,
			},
			want: want{
				err: true,
			},
		},
		{
			name: "unresolved option",
			args: args{
				source: manifestDataNamesSource + manifestExampleDataSource + `
type runner struct{ data string }

func (o runner) runText() {
	switch o.data {
	case dataSimple:
		return
	}
}
`,
			},
			want: want{
				err: true,
			},
		},
	}
	original := targetSources
	targetSources = []targetSource{
		{name: "text", method: "runText", optionFile: "examples/option_text.go"},
	}
	t.Cleanup(func() {
		targetSources = original
	})
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := writeManifestSource(t, test.args.source)
			result, err := parseManifest(root)
			got := want{
				manifest: result,
				err:      err != nil,
			}
			testutil.AssertValue(t, got, test.want, "parseManifest")
		})
	}
}

func Test_parseSource(t *testing.T) {
	type args struct {
		filename func(*testing.T) string
	}
	type want struct {
		name string
		err  bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "source",
			args: args{
				filename: func(t *testing.T) string {
					return writeSource(t, "package fixture")
				},
			},
			want: want{
				name: "fixture",
			},
		},
		{
			name: "missing file",
			args: args{
				filename: func(t *testing.T) string {
					return filepath.Join(t.TempDir(), "missing.go")
				},
			},
			want: want{
				err: true,
			},
		},
		{
			name: "invalid source",
			args: args{
				filename: func(t *testing.T) string {
					return writeSource(t, "package fixture\nfunc")
				},
			},
			want: want{
				err: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			source, err := parseSource(test.args.filename(t))
			got := want{
				err: err != nil,
			}
			if source.file != nil {
				got.name = source.file.Name.Name
			}
			testutil.AssertValue(t, got, test.want, "parseSource")
		})
	}
}

func Test_sourceFile_stringValues(t *testing.T) {
	type args struct {
		source func(*testing.T) sourceFile
	}
	type want struct {
		values map[string]string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "string variables",
			args: args{
				source: func(t *testing.T) sourceFile {
					return parseFixture(t, `package fixture

type value string

var (
	first = "one"
	second, third = "two", "three"
	missing string
	computed = value("four")
)

func ignored() {}
`)
				},
			},
			want: want{
				values: map[string]string{
					"first":  "one",
					"second": "two",
					"third":  "three",
				},
			},
		},
		{
			name: "invalid literal",
			args: args{
				source: func(t *testing.T) sourceFile {
					source := parseFixture(t, `package fixture

var value = "valid"
`)
					declaration := source.file.Decls[0].(*ast.GenDecl)
					value := declaration.Specs[0].(*ast.ValueSpec)
					value.Values[0].(*ast.BasicLit).Value = `"\x"`
					return source
				},
			},
			want: want{
				values: map[string]string{},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			values := test.args.source(t).stringValues()
			got := want{
				values: values,
			}
			testutil.AssertValue(t, got, test.want, "stringValues")
		})
	}
}

func Test_sourceFile_scenariosByTarget(t *testing.T) {
	type args struct {
		source string
		values map[string]string
	}
	type want struct {
		values map[string][]string
		err    bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "target scenarios",
			args: args{
				source: `package fixture

func dataNames(target string) ([]string, bool) {
	switch target {
	case targetText:
		return []string{dataSimple, "literal", dataUnknown}, true
	case targetText, targetHTML:
		return []string{dataSimple}, true
	case "literal":
		return []string{dataSimple}, true
	case targetCSV:
		return nil, true
	case targetBacklog:
		break
	case targetNoResult:
		return
	case targetUnknown:
		return []string{dataSimple}, true
	}
	return nil, false
}
`,
				values: map[string]string{
					"targetText":     "text",
					"targetCSV":      "csv",
					"targetBacklog":  "backlog",
					"targetNoResult": "no-result",
					"dataSimple":     "simple",
				},
			},
			want: want{
				values: map[string][]string{
					"text": {"simple"},
				},
			},
		},
		{
			name: "missing function",
			args: args{
				source: "package fixture",
			},
			want: want{
				err: true,
			},
		},
		{
			name: "missing switch",
			args: args{
				source: "package fixture\nfunc dataNames() {}",
			},
			want: want{
				err: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			values, err := parseFixture(t, test.args.source).scenariosByTarget(test.args.values)
			got := want{
				values: values,
				err:    err != nil,
			}
			testutil.AssertValue(t, got, test.want, "scenariosByTarget")
		})
	}
}

func Test_sourceFile_dataByScenario(t *testing.T) {
	type args struct {
		source string
		values map[string]string
	}
	type want struct {
		values map[string]string
		err    bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "scenario data",
			args: args{
				source: `package fixture

func exampleData(name string) (any, bool) {
	switch name {
	case dataSimple, "literal", dataUnknown:
		return examples.SimpleData, true
	case dataScalar:
		return value, true
	case dataMissing:
		break
	}
	return nil, false
}
`,
				values: map[string]string{
					"dataSimple":  "simple",
					"dataScalar":  "scalar",
					"dataMissing": "missing",
				},
			},
			want: want{
				values: map[string]string{
					"simple": "SimpleData",
				},
			},
		},
		{
			name: "missing function",
			args: args{
				source: "package fixture",
			},
			want: want{
				err: true,
			},
		},
		{
			name: "missing switch",
			args: args{
				source: "package fixture\nfunc exampleData() {}",
			},
			want: want{
				err: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			values, err := parseFixture(t, test.args.source).dataByScenario(test.args.values)
			got := want{
				values: values,
				err:    err != nil,
			}
			testutil.AssertValue(t, got, test.want, "dataByScenario")
		})
	}
}

func Test_sourceFile_options(t *testing.T) {
	type args struct {
		source string
		method string
		values map[string]string
	}
	type want struct {
		values map[string]string
		err    bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "scenario options",
			args: args{
				source: `package fixture

type runner struct{ data string }

func (o runner) runText() {
	switch o.data {
	case dataSimple, "literal", dataUnknown:
		opts = examples.TextOptionSimple
	case dataMissing:
		return
	case dataScalar:
		opts = value
	}
}
`,
				method: "runText",
				values: map[string]string{
					"dataSimple":  "simple",
					"dataMissing": "missing",
					"dataScalar":  "scalar",
				},
			},
			want: want{
				values: map[string]string{
					"simple": "TextOptionSimple",
				},
			},
		},
		{
			name: "missing method",
			args: args{
				source: "package fixture",
				method: "runText",
			},
			want: want{
				err: true,
			},
		},
		{
			name: "missing switch",
			args: args{
				source: "package fixture\nfunc runText() {}",
				method: "runText",
			},
			want: want{
				err: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			values, err := parseFixture(t, test.args.source).options(test.args.method, test.args.values)
			got := want{
				values: values,
				err:    err != nil,
			}
			testutil.AssertValue(t, got, test.want, "options")
		})
	}
}

func Test_sourceFile_function(t *testing.T) {
	type args struct {
		name string
	}
	type want struct {
		name string
		err  bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "function",
			args: args{
				name: "target",
			},
			want: want{
				name: "target",
			},
		},
		{
			name: "missing",
			args: args{
				name: "missing",
			},
			want: want{
				err: true,
			},
		},
	}
	source := parseFixture(t, "package fixture\nvar value string\nfunc target() {}")
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			function, err := source.function(test.args.name)
			got := want{
				err: err != nil,
			}
			if function != nil {
				got.name = function.Name.Name
			}
			testutil.AssertValue(t, got, test.want, "function")
		})
	}
}

func Test_sourceFile_slice(t *testing.T) {
	source := parseFixture(t, "package fixture\nvar Value = 1")
	declaration := source.file.Decls[0]
	type args struct {
		start token.Pos
		end   token.Pos
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "source range",
			args: args{
				start: declaration.Pos(),
				end:   declaration.End(),
			},
			want: "var Value = 1",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := source.slice(test.args.start, test.args.end)
			testutil.AssertValue(t, got, test.want, "slice")
		})
	}
}

func Test_switchStatement(t *testing.T) {
	type args struct {
		function func(*testing.T) *ast.FuncDecl
	}
	type want struct {
		found bool
		err   bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "nested switch",
			args: args{
				function: func(t *testing.T) *ast.FuncDecl {
					return resolveFixtureFunction(t, "func target() { if true { switch {} } }", "target")
				},
			},
			want: want{
				found: true,
			},
		},
		{
			name: "missing switch",
			args: args{
				function: func(t *testing.T) *ast.FuncDecl {
					return resolveFixtureFunction(t, "func target() {}", "target")
				},
			},
			want: want{
				err: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			statement, err := switchStatement(test.args.function(t))
			got := want{
				found: statement != nil,
				err:   err != nil,
			}
			testutil.AssertValue(t, got, test.want, "switchStatement")
		})
	}
}

func Test_returnStatement(t *testing.T) {
	type args struct {
		statements func(*testing.T) []ast.Stmt
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "return",
			args: args{
				statements: func(t *testing.T) []ast.Stmt {
					return resolveFixtureFunction(t, "func target() { value(); return }", "target").Body.List
				},
			},
			want: true,
		},
		{
			name: "missing",
			args: args{
				statements: func(t *testing.T) []ast.Stmt {
					return resolveFixtureFunction(t, "func target() { value() }", "target").Body.List
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := returnStatement(test.args.statements(t)) != nil
			testutil.AssertValue(t, got, test.want, "returnStatement")
		})
	}
}

func Test_optionAssignment(t *testing.T) {
	type args struct {
		statement func(*testing.T) ast.Stmt
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "option",
			args: args{
				statement: func(t *testing.T) ast.Stmt {
					return resolveFixtureStatement(t, "opts = examples.Option")
				},
			},
			want: "Option",
		},
		{
			name: "not assignment",
			args: args{
				statement: func(t *testing.T) ast.Stmt {
					return resolveFixtureStatement(t, "value()")
				},
			},
		},
		{
			name: "multiple values",
			args: args{
				statement: func(t *testing.T) ast.Stmt {
					return resolveFixtureStatement(t, "opts, other = first, second")
				},
			},
		},
		{
			name: "non identifier",
			args: args{
				statement: func(t *testing.T) ast.Stmt {
					return resolveFixtureStatement(t, "values[0] = examples.Option")
				},
			},
		},
		{
			name: "different variable",
			args: args{
				statement: func(t *testing.T) ast.Stmt {
					return resolveFixtureStatement(t, "value = examples.Option")
				},
			},
		},
		{
			name: "non selector",
			args: args{
				statement: func(t *testing.T) ast.Stmt {
					return resolveFixtureStatement(t, "opts = value")
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := optionAssignment([]ast.Stmt{test.args.statement(t)})
			testutil.AssertValue(t, got, test.want, "optionAssignment")
		})
	}
}

func Test_contains(t *testing.T) {
	type args struct {
		values []string
		value  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "found",
			args: args{
				values: []string{"first", "second"},
				value:  "second",
			},
			want: true,
		},
		{
			name: "missing",
			args: args{
				values: []string{"first", "second"},
				value:  "third",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := contains(test.args.values, test.args.value)
			testutil.AssertValue(t, got, test.want, "contains")
		})
	}
}

func Test_declarationSource(t *testing.T) {
	type args struct {
		filename func(*testing.T) string
		name     string
	}
	type want struct {
		source string
		err    bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "declaration",
			args: args{
				filename: func(t *testing.T) string {
					return writeSource(t, "package fixture\n\n// Value documents value.\nvar Value = \"value\"\n")
				},
				name: "Value",
			},
			want: want{
				source: "// Value documents value.\nvar Value = \"value\"",
			},
		},
		{
			name: "missing declaration",
			args: args{
				filename: func(t *testing.T) string {
					return writeSource(t, "package fixture")
				},
				name: "Value",
			},
			want: want{
				err: true,
			},
		},
		{
			name: "missing file",
			args: args{
				filename: func(t *testing.T) string {
					return filepath.Join(t.TempDir(), "missing.go")
				},
				name: "Value",
			},
			want: want{
				err: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			source, err := declarationSource(test.args.filename(t), test.args.name)
			got := want{
				source: source,
				err:    err != nil,
			}
			testutil.AssertValue(t, got, test.want, "declarationSource")
		})
	}
}

func Test_declarationAndCall(t *testing.T) {
	type args struct {
		source string
		name   string
	}
	type want struct {
		source string
		called string
		err    bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "called function",
			args: args{
				source: "package fixture\n\n// Value documents value.\nvar Value = newValue()\n",
				name:   "Value",
			},
			want: want{
				source: "// Value documents value.\nvar Value = newValue()",
				called: "newValue",
			},
		},
		{
			name: "selector call",
			args: args{
				source: "package fixture\n\nvar Value = fixture.New()\n",
				name:   "Value",
			},
			want: want{
				source: "var Value = fixture.New()",
			},
		},
		{
			name: "function",
			args: args{
				source: "package fixture\n\n// newValue creates value.\nfunc newValue() {}\n",
				name:   "newValue",
			},
			want: want{
				source: "// newValue creates value.\nfunc newValue() {}",
			},
		},
		{
			name: "missing",
			args: args{
				source: "package fixture",
				name:   "Value",
			},
			want: want{
				err: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			source, called, err := declarationAndCall(writeSource(t, test.args.source), test.args.name)
			got := want{
				source: source,
				called: called,
				err:    err != nil,
			}
			testutil.AssertValue(t, got, test.want, "declarationAndCall")
		})
	}
}

func writeManifestSource(t *testing.T, source string) string {
	t.Helper()
	root := t.TempDir()
	dir := filepath.Join(root, "examples", "cmd")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	// #nosec G703 -- dir is created beneath testing.T.TempDir.
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(source), 0o600); err != nil {
		t.Fatal(err)
	}
	return root
}

func writeSource(t *testing.T, source string) string {
	t.Helper()
	filename := filepath.Join(t.TempDir(), "fixture.go")
	if err := os.WriteFile(filename, []byte(source), 0o600); err != nil {
		t.Fatal(err)
	}
	return filename
}

func parseFixture(t *testing.T, source string) sourceFile {
	t.Helper()
	parsed, err := parseSource(writeSource(t, source))
	if err != nil {
		t.Fatal(err)
	}
	return parsed
}

func resolveFixtureFunction(t *testing.T, declaration, name string) *ast.FuncDecl {
	t.Helper()
	source := parseFixture(t, "package fixture\n"+declaration)
	function, err := source.function(name)
	if err != nil {
		t.Fatal(err)
	}
	return function
}

func resolveFixtureStatement(t *testing.T, statement string) ast.Stmt {
	t.Helper()
	function := resolveFixtureFunction(t, "func target() { "+statement+" }", "target")
	return function.Body.List[0]
}
