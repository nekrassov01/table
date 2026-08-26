// Package catalog generates the examples catalog from its implementation.
package catalog

import (
	"bytes"
	"fmt"
	"path/filepath"
)

type catalogView struct {
	Scenarios []scenarioView
}

type scenarioView struct {
	Name    string
	Data    string
	Targets []targetView
}

type targetView struct {
	Name         string
	Option       string
	TableOutput  string
	StreamOutput string
	SameOutput   bool
}

// Generate returns the examples catalog using the data, options, and output
// combinations declared beneath root.
func Generate(root string) ([]byte, error) {
	manifest, err := parseManifest(root)
	if err != nil {
		return nil, err
	}
	command, err := newExampleCommand(root)
	if err != nil {
		return nil, err
	}
	defer command.close()
	view := catalogView{
		Scenarios: make([]scenarioView, 0, len(manifest.scenarios)),
	}
	for _, scenario := range manifest.scenarios {
		input, err := inputSource(root, scenario.data)
		if err != nil {
			return nil, err
		}
		scenarioView := scenarioView{
			Name:    scenario.name,
			Data:    input,
			Targets: make([]targetView, 0, len(scenario.targets)),
		}
		for _, target := range scenario.targets {
			option, err := declarationSource(filepath.Join(root, target.optionFile), target.option)
			if err != nil {
				return nil, err
			}
			tableOutput, err := command.run(target.name, "table", scenario.name)
			if err != nil {
				return nil, err
			}
			streamOutput, err := command.run(target.name, "stream", scenario.name)
			if err != nil {
				return nil, err
			}
			scenarioView.Targets = append(scenarioView.Targets, targetView{
				Name:         target.name,
				Option:       option,
				TableOutput:  tableOutput,
				StreamOutput: streamOutput,
				SameOutput:   tableOutput == streamOutput,
			})
		}
		view.Scenarios = append(view.Scenarios, scenarioView)
	}
	var output bytes.Buffer
	if err := catalogTemplate.Execute(&output, view); err != nil {
		return nil, fmt.Errorf("execute examples template: %w", err)
	}
	return append(bytes.TrimRight(output.Bytes(), "\n"), '\n'), nil
}

func inputSource(root, name string) (string, error) {
	filename := filepath.Join(root, "examples", "examples.go")
	source, called, err := declarationAndCall(filename, name)
	if err != nil {
		return "", err
	}
	if called == "" {
		return source, nil
	}
	helper, err := declarationSource(filename, called)
	if err != nil {
		return "", err
	}
	return source + "\n\n" + helper, nil
}
