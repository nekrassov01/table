package catalog

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
)

type manifest struct {
	scenarios []scenario
}

type scenario struct {
	name    string
	data    string
	targets []target
}

type target struct {
	name       string
	option     string
	optionFile string
}

type targetSource struct {
	name       string
	method     string
	optionFile string
}

var targetSources = []targetSource{
	{name: "text", method: "runText", optionFile: "examples/option_text.go"},
	{name: "html", method: "runHTML", optionFile: "examples/option_html.go"},
	{name: "markdown", method: "runMarkdown", optionFile: "examples/option_markdown.go"},
	{name: "backlog", method: "runBacklog", optionFile: "examples/option_backlog.go"},
	{name: "csv", method: "runCSV", optionFile: "examples/option_csv.go"},
}

func parseManifest(root string) (manifest, error) {
	filename := filepath.Join(root, "examples", "cmd", "main.go")
	source, err := parseSource(filename)
	if err != nil {
		return manifest{}, err
	}
	values := source.stringValues()
	scenariosByTarget, err := source.scenariosByTarget(values)
	if err != nil {
		return manifest{}, err
	}
	dataByScenario, err := source.dataByScenario(values)
	if err != nil {
		return manifest{}, err
	}
	options := make(map[string]map[string]string, len(targetSources))
	for _, targetSource := range targetSources {
		options[targetSource.name], err = source.options(targetSource.method, values)
		if err != nil {
			return manifest{}, err
		}
	}
	var result manifest
	seen := make(map[string]bool)
	for _, targetSource := range targetSources {
		for _, name := range scenariosByTarget[targetSource.name] {
			if seen[name] {
				continue
			}
			seen[name] = true
			data, ok := dataByScenario[name]
			if !ok {
				return manifest{}, fmt.Errorf("resolve example data %q", name)
			}
			scenario := scenario{
				name: name,
				data: data,
			}
			for _, scenarioTarget := range targetSources {
				if !contains(scenariosByTarget[scenarioTarget.name], name) {
					continue
				}
				option, ok := options[scenarioTarget.name][name]
				if !ok {
					return manifest{}, fmt.Errorf("resolve %s option for %q", scenarioTarget.name, name)
				}
				scenario.targets = append(scenario.targets, target{
					name:       scenarioTarget.name,
					option:     option,
					optionFile: scenarioTarget.optionFile,
				})
			}
			result.scenarios = append(result.scenarios, scenario)
		}
	}
	return result, nil
}

type sourceFile struct {
	set    *token.FileSet
	file   *ast.File
	source []byte
}

func parseSource(filename string) (sourceFile, error) {
	// #nosec G304 -- Callers construct filenames beneath the repository root.
	source, err := os.ReadFile(filename)
	if err != nil {
		return sourceFile{}, fmt.Errorf("read %s: %w", filename, err)
	}
	set := token.NewFileSet()
	file, err := parser.ParseFile(set, filename, source, parser.ParseComments)
	if err != nil {
		return sourceFile{}, fmt.Errorf("parse %s: %w", filename, err)
	}
	return sourceFile{
		set:    set,
		file:   file,
		source: source,
	}, nil
}

func (o sourceFile) stringValues() map[string]string {
	values := make(map[string]string)
	for _, decl := range o.file.Decls {
		declaration, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range declaration.Specs {
			value, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for index, name := range value.Names {
				if index >= len(value.Values) {
					continue
				}
				literal, ok := value.Values[index].(*ast.BasicLit)
				if !ok || literal.Kind != token.STRING {
					continue
				}
				resolved, err := strconv.Unquote(literal.Value)
				if err != nil {
					continue
				}
				values[name.Name] = resolved
			}
		}
	}
	return values
}

func (o sourceFile) scenariosByTarget(values map[string]string) (map[string][]string, error) {
	function, err := o.function("dataNames")
	if err != nil {
		return nil, err
	}
	statement, err := switchStatement(function)
	if err != nil {
		return nil, err
	}
	result := make(map[string][]string)
	for _, clause := range statement.Body.List {
		caseClause := clause.(*ast.CaseClause)
		if len(caseClause.List) != 1 {
			continue
		}
		targetID, ok := caseClause.List[0].(*ast.Ident)
		if !ok {
			continue
		}
		targetName, ok := values[targetID.Name]
		if !ok {
			continue
		}
		returned := returnStatement(caseClause.Body)
		if returned == nil || len(returned.Results) == 0 {
			continue
		}
		list, ok := returned.Results[0].(*ast.CompositeLit)
		if !ok {
			continue
		}
		for _, element := range list.Elts {
			identifier, ok := element.(*ast.Ident)
			if !ok {
				continue
			}
			name, ok := values[identifier.Name]
			if ok {
				result[targetName] = append(result[targetName], name)
			}
		}
	}
	return result, nil
}

func (o sourceFile) dataByScenario(values map[string]string) (map[string]string, error) {
	function, err := o.function("exampleData")
	if err != nil {
		return nil, err
	}
	statement, err := switchStatement(function)
	if err != nil {
		return nil, err
	}
	result := make(map[string]string)
	for _, clause := range statement.Body.List {
		caseClause := clause.(*ast.CaseClause)
		returned := returnStatement(caseClause.Body)
		if returned == nil || len(returned.Results) == 0 {
			continue
		}
		selector, ok := returned.Results[0].(*ast.SelectorExpr)
		if !ok {
			continue
		}
		for _, expression := range caseClause.List {
			identifier, ok := expression.(*ast.Ident)
			if !ok {
				continue
			}
			name, ok := values[identifier.Name]
			if ok {
				result[name] = selector.Sel.Name
			}
		}
	}
	return result, nil
}

func (o sourceFile) options(method string, values map[string]string) (map[string]string, error) {
	function, err := o.function(method)
	if err != nil {
		return nil, err
	}
	statement, err := switchStatement(function)
	if err != nil {
		return nil, err
	}
	result := make(map[string]string)
	for _, clause := range statement.Body.List {
		caseClause := clause.(*ast.CaseClause)
		option := optionAssignment(caseClause.Body)
		if option == "" {
			continue
		}
		for _, expression := range caseClause.List {
			identifier, ok := expression.(*ast.Ident)
			if !ok {
				continue
			}
			name, ok := values[identifier.Name]
			if ok {
				result[name] = option
			}
		}
	}
	return result, nil
}

func (o sourceFile) function(name string) (*ast.FuncDecl, error) {
	for _, decl := range o.file.Decls {
		function, ok := decl.(*ast.FuncDecl)
		if ok && function.Name.Name == name {
			return function, nil
		}
	}
	return nil, fmt.Errorf("resolve function %s", name)
}

func (o sourceFile) slice(start, end token.Pos) string {
	return string(o.source[o.set.Position(start).Offset:o.set.Position(end).Offset])
}

func switchStatement(function *ast.FuncDecl) (*ast.SwitchStmt, error) {
	var result *ast.SwitchStmt
	ast.Inspect(function.Body, func(node ast.Node) bool {
		if result != nil {
			return false
		}
		statement, ok := node.(*ast.SwitchStmt)
		if ok {
			result = statement
		}
		return true
	})
	if result == nil {
		return nil, fmt.Errorf("resolve switch in %s", function.Name.Name)
	}
	return result, nil
}

func returnStatement(statements []ast.Stmt) *ast.ReturnStmt {
	for _, statement := range statements {
		returned, ok := statement.(*ast.ReturnStmt)
		if ok {
			return returned
		}
	}
	return nil
}

func optionAssignment(statements []ast.Stmt) string {
	for _, statement := range statements {
		assignment, ok := statement.(*ast.AssignStmt)
		if !ok || len(assignment.Lhs) != 1 || len(assignment.Rhs) != 1 {
			continue
		}
		name, ok := assignment.Lhs[0].(*ast.Ident)
		if !ok || name.Name != "opts" {
			continue
		}
		selector, ok := assignment.Rhs[0].(*ast.SelectorExpr)
		if ok {
			return selector.Sel.Name
		}
	}
	return ""
}

func contains(values []string, value string) bool {
	for _, candidate := range values {
		if candidate == value {
			return true
		}
	}
	return false
}

func declarationSource(filename, name string) (string, error) {
	source, _, err := declarationAndCall(filename, name)
	return source, err
}

func declarationAndCall(filename, name string) (string, string, error) {
	source, err := parseSource(filename)
	if err != nil {
		return "", "", err
	}
	for _, decl := range source.file.Decls {
		start := decl.Pos()
		end := decl.End()
		switch node := decl.(type) {
		case *ast.GenDecl:
			for _, spec := range node.Specs {
				value, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}
				for _, identifier := range value.Names {
					if identifier.Name != name {
						continue
					}
					if node.Doc != nil {
						start = node.Doc.Pos()
					}
					called := ""
					if len(value.Values) == 1 {
						if call, ok := value.Values[0].(*ast.CallExpr); ok {
							if function, ok := call.Fun.(*ast.Ident); ok {
								called = function.Name
							}
						}
					}
					return source.slice(start, end), called, nil
				}
			}
		case *ast.FuncDecl:
			if node.Name.Name == name {
				if node.Doc != nil {
					start = node.Doc.Pos()
				}
				return source.slice(start, end), "", nil
			}
		}
	}
	return "", "", fmt.Errorf("resolve declaration %s in %s", name, filename)
}
