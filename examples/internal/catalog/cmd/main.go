package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nekrassov01/table/examples/internal/catalog"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run ./internal/catalog/cmd <repository-root>")
		os.Exit(1)
	}
	root, err := filepath.Abs(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	output, err := catalog.Generate(root)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	filename := filepath.Join(root, "docs", "EXAMPLES.md")
	// #nosec G306 -- The generated documentation is a repository source file.
	if err := os.WriteFile(filename, output, 0o644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
