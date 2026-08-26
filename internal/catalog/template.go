package catalog

import (
	"fmt"
	"strings"
	"text/template"
)

const catalogTemplateText = `# Examples

This file is generated from the example data, options, and command mappings. Run ` + "`make generate`" + ` after changing them; do not edit the catalog directly.

This catalog shows the input, configuration, and output of every runnable example.

## Running examples

Use ` + "`target`" + `, ` + "`mode`" + `, and ` + "`data`" + ` to select the output package, API, and scenario. All three variables are optional, but selecting ` + "`data`" + ` also requires ` + "`mode`" + `.

Omitting ` + "`target`" + ` runs every output package, omitting ` + "`mode`" + ` runs both APIs, and omitting ` + "`data`" + ` runs every scenario available for the selected package. The catalog below shows the same combinations without requiring the command to be run.

Each scenario contains its shared input declaration. Each output-package section then shows the exact Option declaration, commands for its ` + "`Table`" + ` and ` + "`Stream`" + ` examples, and the bytes they produce. Identical results are shown once.

## Catalog

{{range .Scenarios}}- [{{.Name}}](#{{.Name}})
{{end}}
{{range .Scenarios}}{{$scenario := .Name}}### {{.Name}}

Input:

{{code "go" .Data}}

{{range .Targets}}#### {{.Name}}

Configuration:

{{code "go" .Option}}

Example commands:

{{commands .Name $scenario}}

{{if .SameOutput}}` + "`Table`" + ` and ` + "`Stream`" + ` output:

{{output .Name .TableOutput}}
{{else}}` + "`Table`" + ` output:

{{output .Name .TableOutput}}

` + "`Stream`" + ` output:

{{output .Name .StreamOutput}}
{{end}}
{{end}}{{end}}`

var catalogTemplate = template.Must(template.New("examples").Funcs(template.FuncMap{
	"code":     codeBlock,
	"commands": exampleCommands,
	"output":   outputBlock,
}).Parse(catalogTemplateText))

func codeBlock(language, value string) string {
	fenceLength := 4
	runLength := 0
	for _, char := range value {
		if char != '`' {
			runLength = 0
			continue
		}
		runLength++
		if runLength >= fenceLength {
			fenceLength = runLength + 1
		}
	}
	fence := strings.Repeat("`", fenceLength)
	return fmt.Sprintf("%s%s\n%s\n%s", fence, language, strings.TrimSuffix(value, "\n"), fence)
}

func exampleCommands(target, data string) string {
	commands := fmt.Sprintf(
		"make example target=%s mode=table data=%s\nmake example target=%s mode=stream data=%s",
		target,
		data,
		target,
		data,
	)
	return codeBlock("sh", commands)
}

func outputBlock(target, value string) string {
	language := "text"
	if target == "html" {
		language = "html"
	}
	if target == "markdown" {
		language = "markdown"
	}
	return codeBlock(language, value)
}
