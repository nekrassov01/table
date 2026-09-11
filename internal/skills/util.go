package skills

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type command struct {
	name      string
	args      []string
	directory string
	stdout    io.Writer
	stderr    io.Writer
}

type commandExecutor func(context.Context, command) error

func newCommandExecutor() commandExecutor {
	return func(ctx context.Context, input command) error {
		// #nosec G204 -- command names are selected internally and no shell is invoked.
		cmd := exec.CommandContext(ctx, input.name, input.args...)
		cmd.Dir = input.directory
		cmd.Stdout = input.stdout
		cmd.Stderr = input.stderr
		return cmd.Run()
	}
}

func (o commandExecutor) resolveRepositoryRoot(ctx context.Context) (string, error) {
	return o.output(ctx, command{
		name: "git",
		args: []string{"rev-parse", "--show-toplevel"},
	})
}

func (o commandExecutor) output(ctx context.Context, input command) (string, error) {
	var contents bytes.Buffer
	input.stdout = &contents
	input.stderr = &contents
	if err := o(ctx, input); err != nil {
		args := append([]string{input.name}, input.args...)
		return "", fmt.Errorf("run %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(contents.String()))
	}
	return strings.TrimSpace(contents.String()), nil
}
