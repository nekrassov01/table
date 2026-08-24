// Package skills implements maintenance commands invoked by repository skills.
package skills

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"
)

type options struct {
	base      string
	target    string
	bench     string
	benchtime string
	count     int
	keep      bool
}

type metric struct {
	ns     float64
	bytes  float64
	allocs float64
}

type command struct {
	name      string
	args      []string
	directory string
	stdout    io.Writer
	stderr    io.Writer
}

type execution struct {
	lookPath  func(string) (string, error)
	execute   func(context.Context, command) error
	mkdirTemp func(string, string) (string, error)
	mkdir     func(string, os.FileMode) error
	removeAll func(string) error
	create    func(string) (*os.File, error)
}

func newExecution() execution {
	return execution{
		lookPath:  exec.LookPath,
		mkdirTemp: os.MkdirTemp,
		mkdir:     os.Mkdir,
		removeAll: os.RemoveAll,
		create:    os.Create,
		execute: func(ctx context.Context, input command) error {
			// #nosec G204 -- command names are selected internally and no shell is invoked.
			cmd := exec.CommandContext(ctx, input.name, input.args...)
			cmd.Dir = input.directory
			cmd.Stdout = input.stdout
			cmd.Stderr = input.stderr
			return cmd.Run()
		},
	}
}

func (o execution) output(ctx context.Context, input command) (string, error) {
	var contents bytes.Buffer
	input.stdout = &contents
	input.stderr = &contents
	if err := o.execute(ctx, input); err != nil {
		args := append([]string{input.name}, input.args...)
		return "", fmt.Errorf("run %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(contents.String()))
	}
	return strings.TrimSpace(contents.String()), nil
}

// RunBaseline compares benchmark results with a baseline revision.
func RunBaseline(ctx context.Context, arguments []string, stdout, stderr io.Writer) int {
	return runBaseline(ctx, newExecution(), arguments, stdout, stderr)
}

func runBaseline(ctx context.Context, commands execution, arguments []string, stdout, stderr io.Writer) int {
	opts, err := parseBaselineOptions(arguments)
	if err != nil {
		_, _ = fmt.Fprintln(stderr, err)
		return 2
	}
	if err := compareBaseline(ctx, commands, opts, stdout, stderr); err != nil {
		_, _ = fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}

func parseBaselineOptions(arguments []string) (options, error) {
	var opts options
	flags := flag.NewFlagSet("compare", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.StringVar(&opts.base, "base", "HEAD", "revision immediately preceding the change")
	flags.StringVar(&opts.target, "target", "all", "Makefile benchmark target")
	flags.StringVar(&opts.bench, "bench", "", "focused go test benchmark regular expression")
	flags.StringVar(&opts.benchtime, "benchtime", "10000x", "Go benchmark benchtime")
	flags.IntVar(&opts.count, "count", 5, "number of benchmark samples")
	flags.BoolVar(&opts.keep, "keep", false, "keep output and profiles after a successful comparison")
	if err := flags.Parse(arguments); err != nil {
		return options{}, err
	}
	if flags.NArg() != 0 {
		return options{}, fmt.Errorf("unexpected arguments: %s", strings.Join(flags.Args(), " "))
	}
	if opts.count < 1 {
		return options{}, fmt.Errorf("count must be positive")
	}
	return opts, nil
}

func compareBaseline(ctx context.Context, commands execution, opts options, stdout, stderr io.Writer) (err error) {
	root, err := commands.output(ctx, command{
		name: "git",
		args: []string{"rev-parse", "--show-toplevel"},
	})
	if err != nil {
		return err
	}
	baseCommit, err := commands.output(ctx, command{
		name: "git",
		args: []string{"-C", root, "rev-parse", "--verify", opts.base + "^{commit}"},
	})
	if err != nil {
		return err
	}
	benchstat, err := commands.lookPath("benchstat")
	if err != nil {
		return fmt.Errorf("benchstat is unavailable; install it with go install golang.org/x/perf/cmd/benchstat@latest")
	}
	temporary, err := commands.mkdirTemp("", "table-check-baseline.")
	if err != nil {
		return err
	}
	baselineWorktree := filepath.Join(temporary, "baseline")
	results := filepath.Join(temporary, "results")
	worktreeAdded := false
	completed := false
	defer func() {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if worktreeAdded {
			cleanup := command{
				name:   "git",
				args:   []string{"-C", root, "worktree", "remove", "--force", baselineWorktree},
				stderr: stderr,
			}
			if cleanupErr := commands.execute(cleanupCtx, cleanup); cleanupErr != nil && err == nil {
				err = fmt.Errorf("remove baseline worktree: %w", cleanupErr)
			}
		}
		if !opts.keep || !completed {
			if cleanupErr := commands.removeAll(temporary); cleanupErr != nil && err == nil {
				err = fmt.Errorf("remove temporary files: %w", cleanupErr)
			}
		}
	}()
	if err := commands.mkdir(results, 0o755); err != nil {
		return err
	}

	worktree := command{
		name:   "git",
		args:   []string{"-C", root, "worktree", "add", "--detach", baselineWorktree, baseCommit},
		stdout: stderr,
		stderr: stderr,
	}
	if err := commands.execute(ctx, worktree); err != nil {
		return fmt.Errorf("add baseline worktree: %w", err)
	}
	worktreeAdded = true

	goVersion, err := commands.output(ctx, command{
		name: "go",
		args: []string{"version"},
	})
	if err != nil {
		return err
	}
	platform, err := commands.output(ctx, command{
		name: "go",
		args: []string{"env", "GOOS", "GOARCH"},
	})
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintf(stdout, "baseline: %s\n", baseCommit)
	_, _ = fmt.Fprintf(stdout, "go: %s\n", goVersion)
	_, _ = fmt.Fprintf(stdout, "platform: %s\n", strings.ReplaceAll(platform, "\n", "/"))
	if opts.bench != "" {
		_, _ = fmt.Fprintf(stdout, "benchmark: %s\n", opts.bench)
	} else {
		_, _ = fmt.Fprintf(stdout, "target: %s\n", opts.target)
	}
	_, _ = fmt.Fprintf(stdout, "benchtime: %s\n", opts.benchtime)
	_, _ = fmt.Fprintf(stdout, "count: %d\n\n", opts.count)

	before := filepath.Join(results, "before.txt")
	after := filepath.Join(results, "after.txt")
	if err := runBenchmark(ctx, commands, opts, baselineWorktree, results, "before", stderr); err != nil {
		return err
	}
	if err := runBenchmark(ctx, commands, opts, root, results, "after", stderr); err != nil {
		return err
	}

	comparison := command{
		name:      benchstat,
		args:      []string{"before.txt", "after.txt"},
		directory: results,
		stdout:    stdout,
		stderr:    stderr,
	}
	if err := commands.execute(ctx, comparison); err != nil {
		return fmt.Errorf("compare benchmarks: %w", err)
	}
	if err := printMedians(stdout, before, after); err != nil {
		return err
	}

	completed = true
	if opts.keep {
		_, _ = fmt.Fprintf(stdout, "\nartifacts: %s\n", temporary)
	}
	return nil
}

func runBenchmark(ctx context.Context, commands execution, opts options, directory, results, label string, stderr io.Writer) error {
	outputPath := filepath.Join(results, label+".txt")
	outputFile, err := commands.create(outputPath)
	if err != nil {
		return err
	}
	cpuProfile := filepath.Join(results, label+".cpu")
	memoryProfile := filepath.Join(results, label+".mem")
	var benchmark command
	if opts.bench != "" {
		benchmark = command{
			name: "go",
			args: []string{"test",
				"-benchmem",
				"-count", strconv.Itoa(opts.count),
				"-benchtime", opts.benchtime,
				"-cpuprofile", cpuProfile,
				"-memprofile", memoryProfile,
				".",
				"-bench", opts.bench,
			},
			directory: filepath.Join(directory, "benchmarks"),
		}
	} else {
		benchmark = command{
			name: "make",
			args: []string{"bench",
				"target=" + opts.target,
				"benchtime=" + opts.benchtime,
				"count=" + strconv.Itoa(opts.count),
				"cpuprofile=" + cpuProfile,
				"memprofile=" + memoryProfile,
			},
			directory: directory,
		}
	}
	benchmark.stdout = outputFile
	benchmark.stderr = stderr
	runErr := commands.execute(ctx, benchmark)
	closeErr := outputFile.Close()
	if runErr != nil {
		return fmt.Errorf("measure %s: %w", label, runErr)
	}
	if closeErr != nil {
		return fmt.Errorf("close %s benchmark output: %w", label, closeErr)
	}
	return nil
}

func printMedians(w io.Writer, beforePath, afterPath string) error {
	before, err := readMedians(beforePath)
	if err != nil {
		return err
	}
	after, err := readMedians(afterPath)
	if err != nil {
		return err
	}
	names := make([]string, 0, len(before))
	for name := range before {
		names = append(names, name)
	}
	for name := range after {
		if _, ok := before[name]; !ok {
			return fmt.Errorf("benchmark missing from baseline result: %s", name)
		}
	}
	slices.Sort(names)
	_, _ = fmt.Fprintln(w, "\nMedians")
	_, _ = fmt.Fprintln(w, "benchmark\tbefore ns/op\tafter ns/op\tbefore B/op\tafter B/op\tbefore allocs/op\tafter allocs/op")
	for _, name := range names {
		beforeMetric := before[name]
		afterMetric, ok := after[name]
		if !ok {
			return fmt.Errorf("benchmark missing from changed result: %s", name)
		}
		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			name,
			formatMetric(beforeMetric.ns),
			formatMetric(afterMetric.ns),
			formatMetric(beforeMetric.bytes),
			formatMetric(afterMetric.bytes),
			formatMetric(beforeMetric.allocs),
			formatMetric(afterMetric.allocs),
		)
	}
	return nil
}

func readMedians(path string) (map[string]metric, error) {
	// #nosec G304 -- path is an internally generated benchmark result path.
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()
	samples := make(map[string][]metric)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 8 || !strings.HasPrefix(fields[0], "Benchmark") || fields[3] != "ns/op" || fields[5] != "B/op" || fields[7] != "allocs/op" {
			continue
		}
		name := trimCPUSuffix(fields[0])
		ns, err := strconv.ParseFloat(fields[2], 64)
		if err != nil {
			return nil, fmt.Errorf("parse ns/op for %s: %w", name, err)
		}
		bytes, err := strconv.ParseFloat(fields[4], 64)
		if err != nil {
			return nil, fmt.Errorf("parse B/op for %s: %w", name, err)
		}
		allocs, err := strconv.ParseFloat(fields[6], 64)
		if err != nil {
			return nil, fmt.Errorf("parse allocs/op for %s: %w", name, err)
		}
		samples[name] = append(samples[name], metric{
			ns:     ns,
			bytes:  bytes,
			allocs: allocs,
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(samples) == 0 {
		return nil, fmt.Errorf("no benchmark results in %s", path)
	}
	result := make(map[string]metric, len(samples))
	for name, values := range samples {
		ns := make([]float64, len(values))
		bytes := make([]float64, len(values))
		allocs := make([]float64, len(values))
		for index, value := range values {
			ns[index] = value.ns
			bytes[index] = value.bytes
			allocs[index] = value.allocs
		}
		result[name] = metric{
			ns:     median(ns),
			bytes:  median(bytes),
			allocs: median(allocs),
		}
	}
	return result, nil
}

func median(values []float64) float64 {
	slices.Sort(values)
	middle := len(values) / 2
	if len(values)%2 != 0 {
		return values[middle]
	}
	return (values[middle-1] + values[middle]) / 2
}

func trimCPUSuffix(name string) string {
	index := strings.LastIndexByte(name, '-')
	if index < 0 {
		return name
	}
	if _, err := strconv.Atoi(name[index+1:]); err != nil {
		return name
	}
	return name[:index]
}

func formatMetric(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}
