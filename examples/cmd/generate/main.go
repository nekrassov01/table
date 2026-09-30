package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nekrassov01/table/internal/catalog"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run ./examples/cmd/generate <repository-root>")
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
	// The CLI caller chooses the trusted repository root; the output suffix is fixed.
	// #nosec G306 G703 -- Trusted output path; 0644 is for repository documentation.
	if err := os.WriteFile(filename, output, 0o644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
