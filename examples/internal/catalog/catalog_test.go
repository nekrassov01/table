package catalog

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"

	"github.com/nekrassov01/table/internal/testutil"
)

const generatorCommandSource = `package main

import (
	"fmt"
	"os"

	"fixture/examples"
)

var targetText = "text"
var dataSimple = "simple"

type runner struct{ data string }

func main() {
	if len(os.Args) > 2 && os.Args[2] == "{{FAILURE}}" {
		os.Exit(1)
	}
	fmt.Fprintln(os.Stdout, os.Args[2])
}

func dataNames(target string) ([]string, bool) {
	switch target {
	case targetText:
		return []string{dataSimple}, true
	default:
		return nil, false
	}
}

func exampleData(name string) (any, bool) {
	switch name {
	case dataSimple:
		return examples.{{DATA}}, true
	default:
		return nil, false
	}
}

func (o runner) runText() {
	var opts any
	switch o.data {
	case dataSimple:
		opts = examples.TextOptionSimple
	}
	_ = opts
}
`

func TestGenerate(t *testing.T) {
	root := filepath.Clean("../../..")
	document, err := os.ReadFile(filepath.Join(root, "docs", "EXAMPLES.md"))
	if err != nil {
		t.Fatal(err)
	}
	commandSource, err := os.ReadFile(filepath.Join(root, "examples", "cmd", "main.go"))
	if err != nil {
		t.Fatal(err)
	}
	buildErrorRoot := writeManifestSource(t, string(commandSource))
	inputErrorRoot := writeGeneratorModule(t, "MissingData", "")
	tableErrorRoot := writeGeneratorModule(t, "SimpleData", "table")
	streamErrorRoot := writeGeneratorModule(t, "SimpleData", "stream")
	type args struct {
		root            string
		missingOption   bool
		invalidTemplate bool
		singleTarget    bool
	}
	type want struct {
		output []byte
		err    bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "catalog matches documentation",
			args: args{
				root: root,
			},
			want: want{
				output: document,
			},
		},
		{
			name: "missing repository",
			args: args{
				root: t.TempDir(),
			},
			want: want{
				err: true,
			},
		},
		{
			name: "example command build error",
			args: args{
				root: buildErrorRoot,
			},
			want: want{
				err: true,
			},
		},
		{
			name: "missing option declaration",
			args: args{
				root:          root,
				missingOption: true,
			},
			want: want{
				err: true,
			},
		},
		{
			name: "missing input declaration",
			args: args{
				root:         inputErrorRoot,
				singleTarget: true,
			},
			want: want{
				err: true,
			},
		},
		{
			name: "table command error",
			args: args{
				root:         tableErrorRoot,
				singleTarget: true,
			},
			want: want{
				err: true,
			},
		},
		{
			name: "stream command error",
			args: args{
				root:         streamErrorRoot,
				singleTarget: true,
			},
			want: want{
				err: true,
			},
		},
		{
			name: "template error",
			args: args{
				root:            root,
				invalidTemplate: true,
			},
			want: want{
				err: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			originalSources := targetSources
			originalTemplate := catalogTemplate
			t.Cleanup(func() {
				targetSources = originalSources
				catalogTemplate = originalTemplate
			})
			if test.args.missingOption {
				targetSources = append([]targetSource(nil), targetSources...)
				targetSources[0].optionFile = "examples/missing.go"
			}
			if test.args.singleTarget {
				targetSources = []targetSource{
					{name: "text", method: "runText", optionFile: "examples/option_text.go"},
				}
			}
			if test.args.invalidTemplate {
				catalogTemplate = template.Must(template.New("invalid").Parse("{{call .Scenarios}}"))
			}
			output, err := Generate(test.args.root)
			got := want{
				output: output,
				err:    err != nil,
			}
			testutil.AssertValue(t, got, test.want, "Generate")
		})
	}
}

func Test_inputSource(t *testing.T) {
	type args struct {
		root func(*testing.T) string
		name string
	}
	type want struct {
		declaration bool
		helper      bool
		err         bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "direct declaration",
			args: args{
				root: func(*testing.T) string {
					return filepath.Clean("../../..")
				},
				name: "SimpleData",
			},
			want: want{
				declaration: true,
			},
		},
		{
			name: "declaration with helper",
			args: args{
				root: func(*testing.T) string {
					return filepath.Clean("../../..")
				},
				name: "FooterData",
			},
			want: want{
				declaration: true,
				helper:      true,
			},
		},
		{
			name: "missing declaration",
			args: args{
				root: func(*testing.T) string {
					return filepath.Clean("../../..")
				},
				name: "MissingData",
			},
			want: want{
				err: true,
			},
		},
		{
			name: "missing helper",
			args: args{
				root: func(t *testing.T) string {
					return writeExamplesSource(t, "package examples\nvar Value = missing()")
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
			source, err := inputSource(test.args.root(t), test.args.name)
			got := want{
				declaration: strings.Contains(source, "var "+test.args.name),
				helper:      strings.Contains(source, "func newFooterData"),
				err:         err != nil,
			}
			testutil.AssertValue(t, got, test.want, "inputSource")
		})
	}
}

func writeExamplesSource(t *testing.T, source string) string {
	t.Helper()
	root := t.TempDir()
	dir := filepath.Join(root, "examples")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "examples.go"), []byte(source), 0o600); err != nil {
		t.Fatal(err)
	}
	return root
}

func writeGeneratorModule(t *testing.T, dataName, failure string) string {
	t.Helper()
	root := t.TempDir()
	files := map[string]string{
		"go.mod": "module fixture\n\ngo 1.26\n",
		"examples/examples.go": `package examples

var SimpleData = struct{}{}
`,
		"examples/support.go": `package examples

var MissingData = struct{}{}
`,
		"examples/option_text.go": `package examples

var TextOptionSimple = []int{1}
`,
		"examples/cmd/main.go": strings.NewReplacer(
			"{{DATA}}", dataName,
			"{{FAILURE}}", failure,
		).Replace(generatorCommandSource),
	}
	for name, source := range files {
		filename := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(filename), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filename, []byte(source), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	return root
}
