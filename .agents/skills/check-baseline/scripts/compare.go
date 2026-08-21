package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"syscall"
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

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	if err := run(ctx, parseOptions()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, opts options) (err error) {
	root, err := output(ctx, exec.CommandContext(ctx, "git", "rev-parse", "--show-toplevel"))
	if err != nil {
		return err
	}
	baseCommit, err := output(ctx, exec.CommandContext(ctx, "git", "-C", root, "rev-parse", "--verify", opts.base+"^{commit}"))
	if err != nil {
		return err
	}
	benchstat, err := exec.LookPath("benchstat")
	if err != nil {
		return fmt.Errorf("benchstat is unavailable; install it with go install golang.org/x/perf/cmd/benchstat@latest")
	}
	temporary, err := os.MkdirTemp("", "table-check-baseline.")
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
			command := exec.CommandContext(cleanupCtx, "git", "-C", root, "worktree", "remove", "--force", baselineWorktree)
			if cleanupErr := command.Run(); cleanupErr != nil && err == nil {
				err = fmt.Errorf("remove baseline worktree: %w", cleanupErr)
			}
		}
		if !opts.keep || !completed {
			if cleanupErr := os.RemoveAll(temporary); cleanupErr != nil && err == nil {
				err = fmt.Errorf("remove temporary files: %w", cleanupErr)
			}
		}
	}()
	if err := os.Mkdir(results, 0o755); err != nil {
		return err
	}

	worktree := exec.CommandContext(ctx, "git", "-C", root, "worktree", "add", "--detach", baselineWorktree, baseCommit)
	worktree.Stdout = os.Stderr
	worktree.Stderr = os.Stderr
	if err := worktree.Run(); err != nil {
		return fmt.Errorf("add baseline worktree: %w", err)
	}
	worktreeAdded = true

	goVersion, err := output(ctx, exec.CommandContext(ctx, "go", "version"))
	if err != nil {
		return err
	}
	platform, err := output(ctx, exec.CommandContext(ctx, "go", "env", "GOOS", "GOARCH"))
	if err != nil {
		return err
	}
	fmt.Printf("baseline: %s\n", baseCommit)
	fmt.Printf("go: %s\n", goVersion)
	fmt.Printf("platform: %s\n", strings.ReplaceAll(platform, "\n", "/"))
	if opts.bench != "" {
		fmt.Printf("benchmark: %s\n", opts.bench)
	} else {
		fmt.Printf("target: %s\n", opts.target)
	}
	fmt.Printf("benchtime: %s\n", opts.benchtime)
	fmt.Printf("count: %d\n\n", opts.count)

	before := filepath.Join(results, "before.txt")
	after := filepath.Join(results, "after.txt")
	if err := runBenchmark(ctx, opts, baselineWorktree, results, "before"); err != nil {
		return err
	}
	if err := runBenchmark(ctx, opts, root, results, "after"); err != nil {
		return err
	}

	comparison := exec.CommandContext(ctx, benchstat, "before.txt", "after.txt")
	comparison.Dir = results
	comparison.Stdout = os.Stdout
	comparison.Stderr = os.Stderr
	if err := comparison.Run(); err != nil {
		return fmt.Errorf("compare benchmarks: %w", err)
	}
	if err := printMinima(before, after); err != nil {
		return err
	}

	completed = true
	if opts.keep {
		fmt.Printf("\nartifacts: %s\n", temporary)
	}
	return nil
}

func parseOptions() options {
	var opts options
	flag.StringVar(&opts.base, "base", "HEAD", "revision immediately preceding the change")
	flag.StringVar(&opts.target, "target", "all", "Makefile benchmark target")
	flag.StringVar(&opts.bench, "bench", "", "focused go test benchmark regular expression")
	flag.StringVar(&opts.benchtime, "benchtime", "10000x", "Go benchmark benchtime")
	flag.IntVar(&opts.count, "count", 5, "number of benchmark samples")
	flag.BoolVar(&opts.keep, "keep", false, "keep output and profiles after a successful comparison")
	flag.Parse()
	if flag.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "unexpected arguments: %s\n", strings.Join(flag.Args(), " "))
		os.Exit(2)
	}
	if opts.count < 1 {
		fmt.Fprintln(os.Stderr, "count must be positive")
		os.Exit(2)
	}
	return opts
}

func runBenchmark(ctx context.Context, opts options, directory, results, label string) error {
	outputPath := filepath.Join(results, label+".txt")
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	cpuProfile := filepath.Join(results, label+".cpu")
	memoryProfile := filepath.Join(results, label+".mem")
	var command *exec.Cmd
	if opts.bench != "" {
		command = exec.CommandContext(ctx, "go", "test",
			"-benchmem",
			"-count", strconv.Itoa(opts.count),
			"-benchtime", opts.benchtime,
			"-cpuprofile", cpuProfile,
			"-memprofile", memoryProfile,
			".",
			"-bench", opts.bench,
		)
		command.Dir = filepath.Join(directory, "benchmarks")
	} else {
		command = exec.CommandContext(ctx, "make", "bench",
			"target="+opts.target,
			"benchtime="+opts.benchtime,
			"count="+strconv.Itoa(opts.count),
			"cpuprofile="+cpuProfile,
			"memprofile="+memoryProfile,
		)
		command.Dir = directory
	}
	command.Stdout = outputFile
	command.Stderr = os.Stderr
	runErr := command.Run()
	closeErr := outputFile.Close()
	if runErr != nil {
		return fmt.Errorf("measure %s: %w", label, runErr)
	}
	if closeErr != nil {
		return fmt.Errorf("close %s benchmark output: %w", label, closeErr)
	}
	return nil
}

func printMinima(beforePath, afterPath string) error {
	before, err := readMinima(beforePath)
	if err != nil {
		return err
	}
	after, err := readMinima(afterPath)
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
	fmt.Println("\nRaw minima")
	fmt.Println("benchmark\tbefore ns/op\tafter ns/op\tbefore B/op\tafter B/op\tbefore allocs/op\tafter allocs/op")
	for _, name := range names {
		beforeMetric := before[name]
		afterMetric, ok := after[name]
		if !ok {
			return fmt.Errorf("benchmark missing from changed result: %s", name)
		}
		fmt.Printf("%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
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

func readMinima(path string) (map[string]metric, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	result := make(map[string]metric)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 8 || !strings.HasPrefix(fields[0], "Benchmark") || fields[3] != "ns/op" || fields[5] != "B/op" || fields[7] != "allocs/op" {
			continue
		}
		name := trimCPUSuffix(fields[0])
		current, ok := result[name]
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
		if !ok || ns < current.ns {
			current.ns = ns
		}
		if !ok || bytes < current.bytes {
			current.bytes = bytes
		}
		if !ok || allocs < current.allocs {
			current.allocs = allocs
		}
		result[name] = current
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("no benchmark results in %s", path)
	}
	return result, nil
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

func output(ctx context.Context, command *exec.Cmd) (string, error) {
	contents, err := command.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("run %s: %w: %s", command.String(), err, strings.TrimSpace(string(contents)))
	}
	return strings.TrimSpace(string(contents)), nil
}
