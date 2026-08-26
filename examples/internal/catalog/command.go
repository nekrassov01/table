package catalog

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

type exampleCommand struct {
	root string
	dir  string
	path string
}

func newExampleCommand(root string) (*exampleCommand, error) {
	dir, err := os.MkdirTemp("", "table-examples-catalog.*")
	if err != nil {
		return nil, fmt.Errorf("create example command directory: %w", err)
	}
	command := &exampleCommand{
		root: root,
		dir:  dir,
		path: filepath.Join(dir, exampleCommandName(runtime.GOOS)),
	}
	// #nosec G204 -- The executable and arguments are fixed by the generator.
	build := exec.Command("go", "build", "-o", command.path, "./examples/cmd")
	build.Dir = root
	var stderr bytes.Buffer
	build.Stderr = &stderr
	if err := build.Run(); err != nil {
		command.close()
		return nil, fmt.Errorf("build example command: %w: %s", err, stderr.String())
	}
	return command, nil
}

func (o *exampleCommand) run(target, mode, data string) (string, error) {
	// #nosec G204 -- The path is created by this package and arguments come from parsed declarations.
	command := exec.Command(o.path, target, mode, data)
	command.Dir = o.root
	command.Env = append(os.Environ(), "TZ=UTC")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	if err := command.Run(); err != nil {
		return "", fmt.Errorf("run example %s/%s/%s: %w: %s", target, mode, data, err, stderr.String())
	}
	return stdout.String(), nil
}

func (o *exampleCommand) close() {
	_ = os.RemoveAll(o.dir)
}

func exampleCommandName(goos string) string {
	if goos == "windows" {
		return "examples.exe"
	}
	return "examples"
}
