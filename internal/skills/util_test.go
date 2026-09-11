package skills

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/nekrassov01/table/internal/testutil"
)

func Test_newCommandExecutor(t *testing.T) {
	type args struct {
		command command
	}
	type want struct {
		output bool
		failed bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "executes command",
			args: args{
				command: command{
					name: "go",
					args: []string{"version"},
				},
			},
			want: want{
				output: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output, err := newCommandExecutor().output(t.Context(), test.args.command)
			got := want{
				output: strings.HasPrefix(output, "go version"),
				failed: err != nil,
			}
			testutil.AssertValue(t, got, test.want, "newCommandExecutor")
		})
	}
}

func TestCommandExecutor_resolveRepositoryRoot(t *testing.T) {
	type fields struct {
		execute func(context.Context, command) error
	}
	type want struct {
		root   string
		failed bool
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "repository root",
			fields: fields{
				execute: func(_ context.Context, input command) error {
					_, _ = fmt.Fprintln(input.stdout, "/repository")
					return nil
				},
			},
			want: want{
				root: "/repository",
			},
		},
		{
			name: "command error",
			fields: fields{
				execute: func(context.Context, command) error {
					return testutil.NewError()
				},
			},
			want: want{
				failed: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := commandExecutor(test.fields.execute)
			root, err := o.resolveRepositoryRoot(t.Context())
			got := want{
				root:   root,
				failed: err != nil,
			}
			testutil.AssertValue(t, got, test.want, "resolveRepositoryRoot")
		})
	}
}

func TestCommandExecutor_output(t *testing.T) {
	type fields struct {
		execute func(context.Context, command) error
	}
	type want struct {
		output string
		failed bool
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "command output",
			fields: fields{
				execute: func(_ context.Context, input command) error {
					_, _ = fmt.Fprintln(input.stdout, " output ")
					return nil
				},
			},
			want: want{
				output: "output",
			},
		},
		{
			name: "command error",
			fields: fields{
				execute: func(_ context.Context, input command) error {
					_, _ = fmt.Fprint(input.stderr, "reason")
					return testutil.NewError()
				},
			},
			want: want{
				failed: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := commandExecutor(test.fields.execute)
			output, err := o.output(t.Context(), command{
				name: "command",
				args: []string{"argument"},
			})
			got := want{
				output: output,
				failed: err != nil,
			}
			testutil.AssertValue(t, got, test.want, "output")
		})
	}
}
