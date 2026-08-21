package main

import (
	"crypto/sha256"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
)

var defaultMedia = []string{"text", "html", "markdown", "backlog", "csv"}

type goldenTest struct {
	function string
	name     string
	options  []string
}

type duplicate struct {
	names []string
}

type audit struct {
	medium         string
	tests          []goldenTest
	files          []string
	missingFiles   []string
	orphanedFiles  []string
	badReferences  []string
	duplicates     []duplicate
	options        []string
	uncoveredPairs []string
}

func main() {
	media := os.Args[1:]
	if len(media) == 0 {
		media = defaultMedia
	}
	failed := false
	for _, medium := range media {
		result, err := auditMedium(medium)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", medium, err)
			failed = true
			continue
		}
		result.print()
		if len(result.missingFiles) > 0 || len(result.orphanedFiles) > 0 || len(result.badReferences) > 0 {
			failed = true
		}
	}
	if failed {
		os.Exit(1)
	}
}

func (o audit) print() {
	fmt.Printf("%s: tests=%d calls=%d files=%d options=%d uncovered_pairs=%d\n", o.medium, len(o.tests), countCalls(o.tests), len(o.files), len(o.options), len(o.uncoveredPairs))
	printItems("  missing files", o.missingFiles)
	printItems("  orphaned files", o.orphanedFiles)
	printItems("  invalid references", o.badReferences)
	for _, item := range o.duplicates {
		fmt.Printf("  identical bytes: %s\n", strings.Join(item.names, ", "))
	}
	printItems("  uncovered option pairs", o.uncoveredPairs)
}

func auditMedium(medium string) (audit, error) {
	tests, err := parseGoldenTests(filepath.Join(medium, "golden_test.go"))
	if err != nil {
		return audit{}, err
	}
	files, hashes, err := readGoldenFiles(filepath.Join(medium, "testdata"))
	if err != nil {
		return audit{}, err
	}
	options, err := parseOptions(filepath.Join(medium, "option.go"))
	if err != nil {
		return audit{}, err
	}
	result := audit{
		medium:         medium,
		tests:          tests,
		files:          files,
		duplicates:     resolveDuplicates(hashes),
		options:        options,
		uncoveredPairs: resolveUncoveredPairs(options, tests),
	}
	result.reconcile()
	return result, nil
}

func (o *audit) reconcile() {
	references := make(map[string]int, len(o.tests))
	for _, test := range o.tests {
		if test.name == "" {
			o.badReferences = append(o.badReferences, test.function+": expected exactly one AssertGolden call")
			continue
		}
		references[test.name]++
	}
	fileSet := make(map[string]struct{}, len(o.files))
	for _, name := range o.files {
		fileSet[name] = struct{}{}
		if references[name] == 0 {
			o.orphanedFiles = append(o.orphanedFiles, name)
		}
	}
	for name, count := range references {
		if _, ok := fileSet[name]; !ok {
			o.missingFiles = append(o.missingFiles, name)
		}
		want := 0
		switch {
		case strings.HasPrefix(name, "common_"):
			want = 2
		case strings.HasPrefix(name, "table_"), strings.HasPrefix(name, "stream_"):
			want = 1
		default:
			o.badReferences = append(o.badReferences, fmt.Sprintf("%s: unknown prefix", name))
		}
		if want != 0 && count != want {
			o.badReferences = append(o.badReferences, fmt.Sprintf("%s: references=%d want=%d", name, count, want))
		}
	}
	slices.Sort(o.missingFiles)
	slices.Sort(o.orphanedFiles)
	slices.Sort(o.badReferences)
}

func parseGoldenTests(path string) ([]goldenTest, error) {
	set := token.NewFileSet()
	file, err := parser.ParseFile(set, path, nil, 0)
	if err != nil {
		return nil, err
	}
	tests := make([]goldenTest, 0)
	for _, declaration := range file.Decls {
		function, ok := declaration.(*ast.FuncDecl)
		if !ok || function.Recv != nil || !strings.HasPrefix(function.Name.Name, "TestGolden_") {
			continue
		}
		test := goldenTest{
			function: function.Name.Name,
		}
		calls := 0
		optionSet := make(map[string]struct{})
		ast.Inspect(function.Body, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}
			if name := callName(call.Fun); strings.HasPrefix(name, "With") {
				optionSet[name] = struct{}{}
			}
			selector, ok := call.Fun.(*ast.SelectorExpr)
			if !ok || selector.Sel.Name != "AssertGolden" || len(call.Args) < 2 {
				return true
			}
			literal, ok := call.Args[1].(*ast.BasicLit)
			if !ok || literal.Kind != token.STRING {
				return true
			}
			name, err := strconv.Unquote(literal.Value)
			if err != nil {
				return true
			}
			calls++
			test.name = name
			return true
		})
		if calls != 1 {
			test.name = ""
		}
		for option := range optionSet {
			test.options = append(test.options, option)
		}
		slices.Sort(test.options)
		tests = append(tests, test)
	}
	return tests, nil
}

func parseOptions(path string) ([]string, error) {
	set := token.NewFileSet()
	file, err := parser.ParseFile(set, path, nil, 0)
	if err != nil {
		return nil, err
	}
	options := make([]string, 0)
	for _, declaration := range file.Decls {
		function, ok := declaration.(*ast.FuncDecl)
		if ok && function.Recv == nil && strings.HasPrefix(function.Name.Name, "With") {
			options = append(options, function.Name.Name)
		}
	}
	slices.Sort(options)
	return options, nil
}

func readGoldenFiles(directory string) ([]string, map[[sha256.Size]byte][]string, error) {
	paths, err := filepath.Glob(filepath.Join(directory, "*.txt"))
	if err != nil {
		return nil, nil, err
	}
	names := make([]string, 0, len(paths))
	hashes := make(map[[sha256.Size]byte][]string)
	for _, path := range paths {
		contents, err := os.ReadFile(path)
		if err != nil {
			return nil, nil, err
		}
		name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		names = append(names, name)
		hash := sha256.Sum256(contents)
		hashes[hash] = append(hashes[hash], name)
	}
	slices.Sort(names)
	return names, hashes, nil
}

func resolveDuplicates(hashes map[[sha256.Size]byte][]string) []duplicate {
	result := make([]duplicate, 0)
	for _, names := range hashes {
		if len(names) < 2 {
			continue
		}
		slices.Sort(names)
		result = append(result, duplicate{names: names})
	}
	sort.Slice(result, func(i, j int) bool {
		return strings.Join(result[i].names, "\x00") < strings.Join(result[j].names, "\x00")
	})
	return result
}

func resolveUncoveredPairs(options []string, tests []goldenTest) []string {
	covered := make(map[string]struct{})
	for _, test := range tests {
		for left := 0; left < len(test.options); left++ {
			for right := left + 1; right < len(test.options); right++ {
				covered[test.options[left]+" + "+test.options[right]] = struct{}{}
			}
		}
	}
	uncovered := make([]string, 0)
	for left := 0; left < len(options); left++ {
		for right := left + 1; right < len(options); right++ {
			pair := options[left] + " + " + options[right]
			if _, ok := covered[pair]; !ok {
				uncovered = append(uncovered, pair)
			}
		}
	}
	return uncovered
}

func countCalls(tests []goldenTest) int {
	count := 0
	for _, test := range tests {
		if test.name != "" {
			count++
		}
	}
	return count
}

func callName(expression ast.Expr) string {
	switch value := expression.(type) {
	case *ast.Ident:
		return value.Name
	case *ast.SelectorExpr:
		return value.Sel.Name
	default:
		return ""
	}
}

func printItems(label string, items []string) {
	for _, item := range items {
		fmt.Printf("%s: %s\n", label, item)
	}
}
