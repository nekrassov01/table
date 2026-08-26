package catalog

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nekrassov01/table/internal/testutil"
)

func Test_newExampleCommand(t *testing.T) {
	type args struct {
		root    string
		tempDir func(*testing.T) string
	}
	type want struct {
		command bool
		file    bool
		err     bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "command",
			args: args{
				root: filepath.Clean("../../.."),
			},
			want: want{
				command: true,
				file:    true,
			},
		},
		{
			name: "build error",
			args: args{
				root: t.TempDir(),
			},
			want: want{
				err: true,
			},
		},
		{
			name: "temporary directory error",
			args: args{
				root: filepath.Clean("../../.."),
				tempDir: func(t *testing.T) string {
					filename := filepath.Join(t.TempDir(), "file")
					if err := os.WriteFile(filename, nil, 0o600); err != nil {
						t.Fatal(err)
					}
					return filename
				},
			},
			want: want{
				err: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.args.tempDir != nil {
				dir := test.args.tempDir(t)
				t.Setenv("TMPDIR", dir)
				t.Setenv("TMP", dir)
				t.Setenv("TEMP", dir)
			}
			command, err := newExampleCommand(test.args.root)
			got := want{
				command: command != nil,
				err:     err != nil,
			}
			if command != nil {
				_, statErr := os.Stat(command.path)
				got.file = statErr == nil
				command.close()
			}
			testutil.AssertValue(t, got, test.want, "newExampleCommand")
		})
	}
}

func Test_exampleCommand_run(t *testing.T) {
	type args struct {
		target string
		mode   string
		data   string
	}
	type want struct {
		output bool
		err    bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "example",
			args: args{
				target: "text",
				mode:   "table",
				data:   "ascii",
			},
			want: want{
				output: true,
			},
		},
		{
			name: "command error",
			args: args{
				target: "missing",
				mode:   "table",
				data:   "ascii",
			},
			want: want{
				err: true,
			},
		},
	}
	command, err := newExampleCommand(filepath.Clean("../../.."))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(command.close)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output, err := command.run(test.args.target, test.args.mode, test.args.data)
			got := want{
				output: strings.Contains(output, "INSTANCE ID"),
				err:    err != nil,
			}
			testutil.AssertValue(t, got, test.want, "run")
		})
	}
}

func Test_exampleCommand_close(t *testing.T) {
	type fields struct {
		dir string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "removes directory",
			fields: fields{
				dir: t.TempDir(),
			},
			want: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &exampleCommand{dir: test.fields.dir}
			o.close()
			_, err := os.Stat(test.fields.dir)
			got := os.IsNotExist(err)
			testutil.AssertValue(t, got, test.want, "close")
		})
	}
}

func Test_exampleCommandName(t *testing.T) {
	type args struct {
		goos string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "windows",
			args: args{
				goos: "windows",
			},
			want: "examples.exe",
		},
		{
			name: "other system",
			args: args{
				goos: "darwin",
			},
			want: "examples",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := exampleCommandName(test.args.goos)
			testutil.AssertValue(t, got, test.want, "exampleCommandName")
		})
	}
}
