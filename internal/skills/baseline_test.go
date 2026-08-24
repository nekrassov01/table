package skills

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/nekrassov01/table/internal/testutil"
)

func TestNewExecution(t *testing.T) {
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
			output, err := newExecution().output(t.Context(), test.args.command)
			got := want{
				output: strings.HasPrefix(output, "go version"),
				failed: err != nil,
			}
			testutil.AssertValue(t, got, test.want, "newExecution")
		})
	}
}

func TestExecution_output(t *testing.T) {
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
			o := execution{
				execute: test.fields.execute,
			}
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

func TestRunBaseline(t *testing.T) {
	type args struct {
		arguments []string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "invalid arguments",
			args: args{
				arguments: []string{"-count", "0"},
			},
			want: 2,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := RunBaseline(t.Context(), test.args.arguments, &bytes.Buffer{}, &bytes.Buffer{})
			testutil.AssertValue(t, got, test.want, "RunBaseline")
		})
	}
}

func TestRunBaselineWithExecution(t *testing.T) {
	type args struct {
		arguments []string
		failure   string
	}
	type want struct {
		code   int
		stdout bool
		stderr bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "successful comparison",
			args: args{
				arguments: []string{"-bench", "^BenchmarkRender$", "-benchtime", "100x", "-count", "3"},
			},
			want: want{
				stdout: true,
			},
		},
		{
			name: "invalid arguments",
			args: args{
				arguments: []string{"-count", "0"},
			},
			want: want{
				code:   2,
				stderr: true,
			},
		},
		{
			name: "comparison failure",
			args: args{
				arguments: []string{"-bench", "^BenchmarkRender$"},
				failure:   "lookpath",
			},
			want: want{
				code:   1,
				stderr: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer
			state := &executionState{}
			o := newTestExecution(t, state, test.args.failure)
			code := runBaseline(t.Context(), o, test.args.arguments, &stdout, &stderr)
			got := want{
				code:   code,
				stdout: stdout.Len() > 0,
				stderr: stderr.Len() > 0,
			}
			testutil.AssertValue(t, got, test.want, "run")
		})
	}
}

func TestParseBaselineOptions(t *testing.T) {
	type args struct {
		arguments []string
	}
	type want struct {
		opts   options
		failed bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "defaults",
			want: want{
				opts: options{
					base:      "HEAD",
					target:    "all",
					benchtime: "10000x",
					count:     5,
				},
			},
		},
		{
			name: "values",
			args: args{
				arguments: []string{"-base", "main", "-target", "text", "-bench", "Render", "-benchtime", "20x", "-count", "7", "-keep"},
			},
			want: want{
				opts: options{
					base:      "main",
					target:    "text",
					bench:     "Render",
					benchtime: "20x",
					count:     7,
					keep:      true,
				},
			},
		},
		{
			name: "unexpected argument",
			args: args{
				arguments: []string{"argument"},
			},
			want: want{
				failed: true,
			},
		},
		{
			name: "invalid count",
			args: args{
				arguments: []string{"-count", "0"},
			},
			want: want{
				failed: true,
			},
		},
		{
			name: "invalid flag",
			args: args{
				arguments: []string{"-unknown"},
			},
			want: want{
				failed: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			opts, err := parseBaselineOptions(test.args.arguments)
			got := want{
				opts:   opts,
				failed: err != nil,
			}
			testutil.AssertValue(t, got, test.want, "parseOptions")
		})
	}
}

func TestCompareBaseline(t *testing.T) {
	type args struct {
		opts    options
		failure string
	}
	type want struct {
		failed          bool
		benchmark       bool
		target          bool
		medians         bool
		artifacts       bool
		benchmarks      int
		comparisons     int
		worktreeAdds    int
		worktreeRemoves int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "focused benchmark",
			args: args{
				opts: options{
					base:      "HEAD",
					bench:     "^BenchmarkRender$",
					benchtime: "100x",
					count:     3,
				},
			},
			want: want{
				benchmark:       true,
				medians:         true,
				benchmarks:      2,
				comparisons:     1,
				worktreeAdds:    1,
				worktreeRemoves: 1,
			},
		},
		{
			name: "make target",
			args: args{
				opts: options{
					base:      "HEAD",
					target:    "text",
					benchtime: "100x",
					count:     3,
				},
			},
			want: want{
				target:          true,
				medians:         true,
				benchmarks:      2,
				comparisons:     1,
				worktreeAdds:    1,
				worktreeRemoves: 1,
			},
		},
		{
			name: "missing benchstat",
			args: args{
				opts: options{
					base: "HEAD",
				},
				failure: "lookpath",
			},
			want: want{
				failed: true,
			},
		},
		{
			name: "repository lookup error",
			args: args{
				opts: options{
					base: "HEAD",
				},
				failure: "root",
			},
			want: want{
				failed: true,
			},
		},
		{
			name: "base revision error",
			args: args{
				opts: options{
					base: "missing",
				},
				failure: "base",
			},
			want: want{
				failed: true,
			},
		},
		{
			name: "temporary directory error",
			args: args{
				opts: options{
					base: "HEAD",
				},
				failure: "temporary",
			},
			want: want{
				failed: true,
			},
		},
		{
			name: "results directory error",
			args: args{
				opts: options{
					base: "HEAD",
				},
				failure: "mkdir",
			},
			want: want{
				failed: true,
			},
		},
		{
			name: "worktree error",
			args: args{
				opts: options{
					base: "HEAD",
				},
				failure: "worktree",
			},
			want: want{
				failed:       true,
				worktreeAdds: 1,
			},
		},
		{
			name: "go version error cleans worktree",
			args: args{
				opts: options{
					base: "HEAD",
				},
				failure: "version",
			},
			want: want{
				failed:          true,
				worktreeAdds:    1,
				worktreeRemoves: 1,
			},
		},
		{
			name: "platform error cleans worktree",
			args: args{
				opts: options{
					base: "HEAD",
				},
				failure: "platform",
			},
			want: want{
				failed:          true,
				worktreeAdds:    1,
				worktreeRemoves: 1,
			},
		},
		{
			name: "benchmark error cleans worktree",
			args: args{
				opts: options{
					base:      "HEAD",
					bench:     "^BenchmarkRender$",
					benchtime: "100x",
					count:     3,
				},
				failure: "benchmark",
			},
			want: want{
				failed:          true,
				benchmark:       true,
				benchmarks:      1,
				worktreeAdds:    1,
				worktreeRemoves: 1,
			},
		},
		{
			name: "second benchmark error cleans worktree",
			args: args{
				opts: options{
					base:      "HEAD",
					bench:     "^BenchmarkRender$",
					benchtime: "100x",
					count:     3,
				},
				failure: "after",
			},
			want: want{
				failed:          true,
				benchmark:       true,
				benchmarks:      2,
				worktreeAdds:    1,
				worktreeRemoves: 1,
			},
		},
		{
			name: "comparison error cleans worktree",
			args: args{
				opts: options{
					base:      "HEAD",
					bench:     "^BenchmarkRender$",
					benchtime: "100x",
					count:     3,
				},
				failure: "comparison",
			},
			want: want{
				failed:          true,
				benchmark:       true,
				benchmarks:      2,
				comparisons:     1,
				worktreeAdds:    1,
				worktreeRemoves: 1,
			},
		},
		{
			name: "invalid benchmark output cleans worktree",
			args: args{
				opts: options{
					base:      "HEAD",
					bench:     "^BenchmarkRender$",
					benchtime: "100x",
					count:     3,
				},
				failure: "result",
			},
			want: want{
				failed:          true,
				benchmark:       true,
				benchmarks:      2,
				comparisons:     1,
				worktreeAdds:    1,
				worktreeRemoves: 1,
			},
		},
		{
			name: "cleanup error",
			args: args{
				opts: options{
					base:      "HEAD",
					bench:     "^BenchmarkRender$",
					benchtime: "100x",
					count:     3,
				},
				failure: "cleanup",
			},
			want: want{
				failed:          true,
				benchmark:       true,
				medians:         true,
				benchmarks:      2,
				comparisons:     1,
				worktreeAdds:    1,
				worktreeRemoves: 1,
			},
		},
		{
			name: "temporary file cleanup error",
			args: args{
				opts: options{
					base:      "HEAD",
					bench:     "^BenchmarkRender$",
					benchtime: "100x",
					count:     3,
				},
				failure: "remove",
			},
			want: want{
				failed:          true,
				benchmark:       true,
				medians:         true,
				benchmarks:      2,
				comparisons:     1,
				worktreeAdds:    1,
				worktreeRemoves: 1,
			},
		},
		{
			name: "keeps artifacts",
			args: args{
				opts: options{
					base:      "HEAD",
					bench:     "^BenchmarkRender$",
					benchtime: "100x",
					count:     3,
					keep:      true,
				},
			},
			want: want{
				benchmark:       true,
				medians:         true,
				artifacts:       true,
				benchmarks:      2,
				comparisons:     1,
				worktreeAdds:    1,
				worktreeRemoves: 1,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer
			state := &executionState{}
			o := newTestExecution(t, state, test.args.failure)
			err := compareBaseline(t.Context(), o, test.args.opts, &stdout, &stderr)
			got := want{
				failed:          err != nil,
				benchmark:       strings.Contains(stdout.String(), "benchmark: "),
				target:          strings.Contains(stdout.String(), "target: "),
				medians:         strings.Contains(stdout.String(), "Medians"),
				artifacts:       strings.Contains(stdout.String(), "artifacts: "),
				benchmarks:      state.benchmarks,
				comparisons:     state.comparisons,
				worktreeAdds:    state.worktreeAdds,
				worktreeRemoves: state.worktreeRemoves,
			}
			testutil.AssertValue(t, got, test.want, "compare")
		})
	}
}

func TestRunBenchmark(t *testing.T) {
	type args struct {
		opts    options
		failure string
	}
	type want struct {
		name      string
		bench     bool
		target    bool
		directory string
		contents  string
		failed    bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "focused benchmark",
			args: args{
				opts: options{
					bench:     "Render",
					benchtime: "10x",
					count:     2,
				},
			},
			want: want{
				name:     "go",
				bench:    true,
				contents: benchmarkOutput,
			},
		},
		{
			name: "make target",
			args: args{
				opts: options{
					target:    "text",
					benchtime: "10x",
					count:     2,
				},
			},
			want: want{
				name:     "make",
				target:   true,
				contents: benchmarkOutput,
			},
		},
		{
			name: "create output error",
			args: args{
				opts: options{
					bench: "Render",
				},
				failure: "create",
			},
			want: want{
				failed: true,
			},
		},
		{
			name: "benchmark error",
			args: args{
				opts: options{
					bench: "Render",
				},
				failure: "execute",
			},
			want: want{
				name:     "go",
				bench:    true,
				contents: benchmarkOutput,
				failed:   true,
			},
		},
		{
			name: "close output error",
			args: args{
				opts: options{
					bench: "Render",
				},
				failure: "close",
			},
			want: want{
				name:   "go",
				bench:  true,
				failed: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			directory := t.TempDir()
			results := t.TempDir()
			var observed command
			o := execution{
				create: func(path string) (*os.File, error) {
					if test.args.failure == "create" {
						return nil, testutil.NewError()
					}
					// #nosec G304 -- path is built from this test's temporary result directory.
					file, err := os.Create(path)
					if err == nil && test.args.failure == "close" {
						_ = file.Close()
					}
					return file, err
				},
				execute: func(_ context.Context, input command) error {
					observed = input
					_, _ = fmt.Fprint(input.stdout, benchmarkOutput)
					if test.args.failure == "execute" {
						return testutil.NewError()
					}
					return nil
				},
			}
			err := runBenchmark(t.Context(), o, test.args.opts, directory, results, "before", &bytes.Buffer{})
			// #nosec G304 -- the path is inside this test's temporary result directory.
			contents, _ := os.ReadFile(filepath.Join(results, "before.txt"))
			got := want{
				name:      observed.name,
				bench:     slices.Contains(observed.args, "-bench"),
				target:    slices.Contains(observed.args, "target=text"),
				directory: observed.directory,
				contents:  string(contents),
				failed:    err != nil,
			}
			if test.want.name != "" {
				test.want.directory = directory
				if test.args.opts.bench != "" {
					test.want.directory = filepath.Join(directory, "benchmarks")
				}
			}
			testutil.AssertValue(t, got, test.want, "runBenchmark")
		})
	}
}

func TestPrintMedians(t *testing.T) {
	type args struct {
		before string
		after  string
	}
	type want struct {
		output string
		failed bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "prints sorted medians",
			args: args{
				before: benchmarkLine("BenchmarkB-8", 20, 200, 2) + benchmarkLine("BenchmarkA-8", 10, 100, 1),
				after:  benchmarkLine("BenchmarkA-8", 11, 101, 2) + benchmarkLine("BenchmarkB-8", 21, 201, 3),
			},
			want: want{
				output: "\nMedians\n" +
					"benchmark\tbefore ns/op\tafter ns/op\tbefore B/op\tafter B/op\tbefore allocs/op\tafter allocs/op\n" +
					"BenchmarkA\t10\t11\t100\t101\t1\t2\n" +
					"BenchmarkB\t20\t21\t200\t201\t2\t3\n",
			},
		},
		{
			name: "benchmark absent from baseline",
			args: args{
				before: benchmarkLine("BenchmarkA-8", 10, 100, 1),
				after:  benchmarkLine("BenchmarkA-8", 11, 101, 1) + benchmarkLine("BenchmarkB-8", 20, 200, 2),
			},
			want: want{
				failed: true,
			},
		},
		{
			name: "invalid baseline output",
			args: args{
				before: "PASS\n",
				after:  benchmarkLine("BenchmarkA-8", 11, 101, 1),
			},
			want: want{
				failed: true,
			},
		},
		{
			name: "invalid changed output",
			args: args{
				before: benchmarkLine("BenchmarkA-8", 10, 100, 1),
				after:  "PASS\n",
			},
			want: want{
				failed: true,
			},
		},
		{
			name: "benchmark absent from change",
			args: args{
				before: benchmarkLine("BenchmarkA-8", 10, 100, 1) + benchmarkLine("BenchmarkB-8", 20, 200, 2),
				after:  benchmarkLine("BenchmarkA-8", 11, 101, 1),
			},
			want: want{
				output: "\nMedians\nbenchmark\tbefore ns/op\tafter ns/op\tbefore B/op\tafter B/op\tbefore allocs/op\tafter allocs/op\nBenchmarkA\t10\t11\t100\t101\t1\t1\n",
				failed: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			before := writeBenchmarkFile(t, test.args.before)
			after := writeBenchmarkFile(t, test.args.after)
			var w bytes.Buffer
			err := printMedians(&w, before, after)
			got := want{
				output: w.String(),
				failed: err != nil,
			}
			testutil.AssertValue(t, got, test.want, "printMedians")
		})
	}
}

func TestReadMedians(t *testing.T) {
	type args struct {
		path func(*testing.T) string
	}
	type want struct {
		metrics map[string]metric
		failed  bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "returns median metrics",
			args: args{
				path: func(t *testing.T) string {
					return writeBenchmarkFile(t, benchmarkLine("BenchmarkRender-8", 30, 300, 3)+
						benchmarkLine("BenchmarkRender-8", 10, 100, 1)+
						benchmarkLine("BenchmarkRender-8", 20, 200, 2))
				},
			},
			want: want{
				metrics: map[string]metric{
					"BenchmarkRender": {
						ns:     20,
						bytes:  200,
						allocs: 2,
					},
				},
			},
		},
		{
			name: "no benchmark results",
			args: args{
				path: func(t *testing.T) string {
					return writeBenchmarkFile(t, "PASS\n")
				},
			},
			want: want{
				failed: true,
			},
		},
		{
			name: "invalid execution time",
			args: args{
				path: func(t *testing.T) string {
					return writeBenchmarkFile(t, "BenchmarkRender-8 100 invalid ns/op 100 B/op 1 allocs/op\n")
				},
			},
			want: want{
				failed: true,
			},
		},
		{
			name: "invalid bytes",
			args: args{
				path: func(t *testing.T) string {
					return writeBenchmarkFile(t, "BenchmarkRender-8 100 10 ns/op invalid B/op 1 allocs/op\n")
				},
			},
			want: want{
				failed: true,
			},
		},
		{
			name: "invalid allocations",
			args: args{
				path: func(t *testing.T) string {
					return writeBenchmarkFile(t, "BenchmarkRender-8 100 10 ns/op 100 B/op invalid allocs/op\n")
				},
			},
			want: want{
				failed: true,
			},
		},
		{
			name: "missing file",
			args: args{
				path: func(t *testing.T) string {
					return filepath.Join(t.TempDir(), "missing.txt")
				},
			},
			want: want{
				failed: true,
			},
		},
		{
			name: "scan error",
			args: args{
				path: func(t *testing.T) string {
					return writeBenchmarkFile(t, strings.Repeat("x", 70_000))
				},
			},
			want: want{
				failed: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			metrics, err := readMedians(test.args.path(t))
			got := want{
				metrics: metrics,
				failed:  err != nil,
			}
			testutil.AssertValue(t, got, test.want, "readMedians")
		})
	}
}

func TestMedian(t *testing.T) {
	type args struct {
		values []float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "odd sample count",
			args: args{
				values: []float64{3, 1, 2},
			},
			want: 2,
		},
		{
			name: "even sample count",
			args: args{
				values: []float64{4, 1, 3, 2},
			},
			want: 2.5,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := median(test.args.values)
			testutil.AssertValue(t, got, test.want, "median")
		})
	}
}

func TestTrimCPUSuffix(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "CPU suffix",
			args: args{
				name: "BenchmarkRender-8",
			},
			want: "BenchmarkRender",
		},
		{
			name: "no suffix",
			args: args{
				name: "BenchmarkRender",
			},
			want: "BenchmarkRender",
		},
		{
			name: "non-numeric suffix",
			args: args{
				name: "BenchmarkRender-arm64",
			},
			want: "BenchmarkRender-arm64",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := trimCPUSuffix(test.args.name)
			testutil.AssertValue(t, got, test.want, "trimCPUSuffix")
		})
	}
}

func TestFormatMetric(t *testing.T) {
	type args struct {
		value float64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "integer",
			args: args{
				value: 10,
			},
			want: "10",
		},
		{
			name: "fraction",
			args: args{
				value: 2.5,
			},
			want: "2.5",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := formatMetric(test.args.value)
			testutil.AssertValue(t, got, test.want, "formatMetric")
		})
	}
}

const benchmarkOutput = "BenchmarkRender-8 100 10 ns/op 100 B/op 1 allocs/op\n"

type executionState struct {
	root            string
	failure         string
	benchmarks      int
	comparisons     int
	worktreeAdds    int
	worktreeRemoves int
}

func (o *executionState) execute(_ context.Context, input command) error {
	arguments := strings.Join(input.args, " ")
	switch input.name {
	case "git":
		return o.executeGit(arguments, input)
	case "go":
		if strings.HasPrefix(arguments, "test ") {
			return o.executeBenchmark(input)
		}
		return o.executeGo(arguments, input)
	case "make":
		return o.executeBenchmark(input)
	case "benchstat-test":
		return o.executeComparison(input)
	default:
		return nil
	}
}

func (o *executionState) executeGit(arguments string, input command) error {
	switch {
	case arguments == "rev-parse --show-toplevel":
		if o.failure == "root" {
			return testutil.NewError()
		}
		_, _ = fmt.Fprintln(input.stdout, o.root)
	case strings.Contains(arguments, "rev-parse --verify"):
		if o.failure == "base" {
			return testutil.NewError()
		}
		_, _ = fmt.Fprintln(input.stdout, "base")
	case strings.Contains(arguments, "worktree add"):
		o.worktreeAdds++
		if o.failure == "worktree" {
			return testutil.NewError()
		}
	case strings.Contains(arguments, "worktree remove"):
		o.worktreeRemoves++
		if o.failure == "cleanup" {
			return testutil.NewError()
		}
	}
	return nil
}

func (o *executionState) executeGo(arguments string, input command) error {
	switch arguments {
	case "version":
		if o.failure == "version" {
			return testutil.NewError()
		}
		_, _ = fmt.Fprintln(input.stdout, "go version go1.test test/arch")
	case "env GOOS GOARCH":
		if o.failure == "platform" {
			return testutil.NewError()
		}
		_, _ = fmt.Fprintln(input.stdout, "test")
		_, _ = fmt.Fprintln(input.stdout, "arch")
	}
	return nil
}

func (o *executionState) executeBenchmark(input command) error {
	o.benchmarks++
	if o.failure == "benchmark" || o.failure == "after" && o.benchmarks == 2 {
		return testutil.NewError()
	}
	if o.failure == "result" {
		_, _ = fmt.Fprintln(input.stdout, "PASS")
		return nil
	}
	_, _ = fmt.Fprint(input.stdout, benchmarkOutput)
	return nil
}

func (o *executionState) executeComparison(input command) error {
	o.comparisons++
	if o.failure == "comparison" {
		return testutil.NewError()
	}
	_, _ = fmt.Fprintln(input.stdout, "comparison")
	return nil
}

func newTestExecution(t *testing.T, state *executionState, failure string) execution {
	t.Helper()
	state.root = t.TempDir()
	state.failure = failure
	return execution{
		mkdirTemp: func(directory, pattern string) (string, error) {
			if state.failure == "temporary" {
				return "", testutil.NewError()
			}
			path, err := os.MkdirTemp(directory, pattern)
			if err == nil {
				t.Cleanup(func() {
					_ = os.RemoveAll(path)
				})
			}
			return path, err
		},
		mkdir: func(path string, permissions os.FileMode) error {
			if state.failure == "mkdir" {
				return testutil.NewError()
			}
			return os.Mkdir(path, permissions)
		},
		removeAll: func(path string) error {
			if state.failure == "remove" {
				_ = os.RemoveAll(path)
				return testutil.NewError()
			}
			return os.RemoveAll(path)
		},
		create: os.Create,
		lookPath: func(string) (string, error) {
			if state.failure == "lookpath" {
				return "", testutil.NewError()
			}
			return "benchstat-test", nil
		},
		execute: state.execute,
	}
}

func benchmarkLine(name string, ns, bytes, allocs int) string {
	return fmt.Sprintf("%s 100 %d ns/op %d B/op %d allocs/op\n", name, ns, bytes, allocs)
}

func writeBenchmarkFile(t *testing.T, contents string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "benchmark.txt")
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}
