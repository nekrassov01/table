package skills

import (
	"bytes"
	"context"
	"crypto/sha256"
	"go/ast"
	"go/token"
	"os"
	"path/filepath"
	"testing"

	"github.com/nekrassov01/table/internal/testutil"
)

func TestAudit_write(t *testing.T) {
	type fields struct {
		audit audit
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "writes summary and findings",
			fields: fields{
				audit: audit{
					medium: "text",
					tests: []goldenTest{
						{
							name: "common_value",
						},
						{},
					},
					files:          []string{"common_value", "orphan"},
					missingFiles:   []string{"table_missing"},
					orphanedFiles:  []string{"orphan"},
					badReferences:  []string{"invalid"},
					duplicates:     []duplicate{{names: []string{"a", "b"}}},
					options:        []string{"WithHeader"},
					uncoveredPairs: []string{"WithA + WithB"},
				},
			},
			want: "text: tests=2 calls=1 files=2 options=1 uncovered_pairs=1\n" +
				"  missing files: table_missing\n" +
				"  orphaned files: orphan\n" +
				"  invalid references: invalid\n" +
				"  identical bytes: a, b\n" +
				"  uncovered option pairs: WithA + WithB\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var w bytes.Buffer
			test.fields.audit.write(&w)
			testutil.AssertValue(t, w.String(), test.want, "write")
		})
	}
}

func TestAudit_reconcile(t *testing.T) {
	type fields struct {
		audit audit
	}
	type want struct {
		missingFiles  []string
		orphanedFiles []string
		badReferences []string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "reconciles files and references",
			fields: fields{
				audit: audit{
					tests: []goldenTest{
						{
							function: "TestGolden_TableCommon",
							name:     "common_value",
						},
						{
							function: "TestGolden_StreamCommon",
							name:     "common_value",
						},
						{
							function: "TestGolden_TableOnly",
							name:     "table_value",
						},
						{
							function: "TestGolden_StreamMissing",
							name:     "stream_missing",
						},
						{
							function: "TestGolden_CommonSingle",
							name:     "common_single",
						},
						{
							function: "TestGolden_Unknown",
							name:     "unknown_value",
						},
						{
							function: "TestGolden_Invalid",
						},
					},
					files: []string{"common_single", "common_value", "orphan", "table_value"},
				},
			},
			want: want{
				missingFiles:  []string{"stream_missing", "unknown_value"},
				orphanedFiles: []string{"orphan"},
				badReferences: []string{
					"TestGolden_Invalid: expected exactly one AssertGolden call",
					"common_single: references=1 want=2",
					"unknown_value: unknown prefix",
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := test.fields.audit
			o.reconcile()
			got := want{
				missingFiles:  o.missingFiles,
				orphanedFiles: o.orphanedFiles,
				badReferences: o.badReferences,
			}
			testutil.AssertValue(t, got, test.want, "reconcile")
		})
	}
}

func TestRunGoldens(t *testing.T) {
	type args struct {
		ctx   context.Context
		media func(*testing.T) []string
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
			name: "successful audit",
			args: args{
				ctx: context.Background(),
				media: func(t *testing.T) []string {
					return []string{newAuditMedium(t)}
				},
			},
			want: want{
				stdout: true,
			},
		},
		{
			name: "default media",
			args: args{
				ctx: context.Background(),
				media: func(t *testing.T) []string {
					previous := defaultMedia
					defaultMedia = []string{newAuditMedium(t)}
					t.Cleanup(func() {
						defaultMedia = previous
					})
					return nil
				},
			},
			want: want{
				stdout: true,
			},
		},
		{
			name: "audit findings",
			args: args{
				ctx: context.Background(),
				media: func(t *testing.T) []string {
					medium := newAuditMedium(t)
					if err := os.Remove(filepath.Join(medium, "testdata", "common_value.txt")); err != nil {
						t.Fatal(err)
					}
					return []string{medium}
				},
			},
			want: want{
				code:   1,
				stdout: true,
			},
		},
		{
			name: "failed audit",
			args: args{
				ctx: context.Background(),
				media: func(t *testing.T) []string {
					return []string{filepath.Join(t.TempDir(), "missing")}
				},
			},
			want: want{
				code:   1,
				stderr: true,
			},
		},
		{
			name: "canceled audit",
			args: args{
				ctx: canceledContext(),
				media: func(*testing.T) []string {
					return []string{"text"}
				},
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
			code := RunGoldens(test.args.ctx, test.args.media(t), &stdout, &stderr)
			got := want{
				code:   code,
				stdout: stdout.Len() > 0,
				stderr: stderr.Len() > 0,
			}
			testutil.AssertValue(t, got, test.want, "run")
		})
	}
}

func TestAuditMedium(t *testing.T) {
	type args struct {
		medium func(*testing.T) string
	}
	type want struct {
		tests          []goldenTest
		files          []string
		duplicates     []duplicate
		options        []string
		uncoveredPairs []string
		failed         bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "audits medium",
			args: args{
				medium: newAuditMedium,
			},
			want: want{
				tests: []goldenTest{
					{
						function: "TestGolden_TableValue",
						name:     "common_value",
						options:  []string{"WithAlpha", "WithBeta"},
					},
					{
						function: "TestGolden_StreamValue",
						name:     "common_value",
						options:  []string{"WithAlpha", "WithBeta"},
					},
				},
				files:          []string{"common_value"},
				duplicates:     []duplicate{},
				options:        []string{"WithAlpha", "WithBeta"},
				uncoveredPairs: []string{},
			},
		},
		{
			name: "invalid golden tests",
			args: args{
				medium: func(t *testing.T) string {
					medium := newAuditMedium(t)
					writeFile(t, filepath.Join(medium, "golden_test.go"), "package")
					return medium
				},
			},
			want: want{
				failed: true,
			},
		},
		{
			name: "invalid golden directory pattern",
			args: args{
				medium: func(t *testing.T) string {
					medium := filepath.Join(t.TempDir(), "[")
					if err := os.Mkdir(medium, 0o700); err != nil {
						t.Fatal(err)
					}
					writeAuditMedium(t, medium)
					return medium
				},
			},
			want: want{
				failed: true,
			},
		},
		{
			name: "invalid options",
			args: args{
				medium: func(t *testing.T) string {
					medium := newAuditMedium(t)
					writeFile(t, filepath.Join(medium, "option.go"), "package")
					return medium
				},
			},
			want: want{
				failed: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := auditMedium(test.args.medium(t))
			got := want{
				tests:          result.tests,
				files:          result.files,
				duplicates:     result.duplicates,
				options:        result.options,
				uncoveredPairs: result.uncoveredPairs,
				failed:         err != nil,
			}
			testutil.AssertValue(t, got, test.want, "auditMedium")
		})
	}
}

func TestParseGoldenTests(t *testing.T) {
	type args struct {
		contents string
	}
	type want struct {
		tests  []goldenTest
		failed bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "parses tests and options",
			args: args{
				contents: `package fixture

func TestGolden_Table(t *testing.T) {
	WithBeta()
	pkg.WithAlpha()
	testutil.AssertGolden(t, "table_value", nil)
}

func TestGolden_Invalid(t *testing.T) {
	testutil.AssertGolden(t, "table_first", nil)
	testutil.AssertGolden(t, "table_second", nil)
}

func TestGolden_NonString(t *testing.T) {
	testutil.AssertGolden(t, name, nil)
}

func helper() {}

type suite struct{}

func (suite) TestGolden_Method(t *testing.T) {}
`,
			},
			want: want{
				tests: []goldenTest{
					{
						function: "TestGolden_Table",
						name:     "table_value",
						options:  []string{"WithAlpha", "WithBeta"},
					},
					{
						function: "TestGolden_Invalid",
					},
					{
						function: "TestGolden_NonString",
					},
				},
			},
		},
		{
			name: "invalid source",
			args: args{
				contents: "package",
			},
			want: want{
				failed: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path := writeTestFile(t, "golden_test.go", test.args.contents)
			tests, err := parseGoldenTests(path)
			got := want{
				tests:  tests,
				failed: err != nil,
			}
			testutil.AssertValue(t, got, test.want, "parseGoldenTests")
		})
	}
}

func TestCallName(t *testing.T) {
	type args struct {
		expression ast.Expr
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "identifier",
			args: args{
				expression: &ast.Ident{Name: "WithHeader"},
			},
			want: "WithHeader",
		},
		{
			name: "selector",
			args: args{
				expression: &ast.SelectorExpr{
					Sel: &ast.Ident{Name: "WithFooter"},
				},
			},
			want: "WithFooter",
		},
		{
			name: "other expression",
			args: args{
				expression: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"value"`,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := callName(test.args.expression)
			testutil.AssertValue(t, got, test.want, "callName")
		})
	}
}

func TestGoldenName(t *testing.T) {
	type args struct {
		expression ast.Expr
	}
	type want struct {
		name string
		ok   bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "string literal",
			args: args{
				expression: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"value"`,
				},
			},
			want: want{
				name: "value",
				ok:   true,
			},
		},
		{
			name: "other expression",
			args: args{
				expression: &ast.Ident{Name: "name"},
			},
		},
		{
			name: "invalid string literal",
			args: args{
				expression: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"\x"`,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			name, ok := goldenName(test.args.expression)
			got := want{
				name: name,
				ok:   ok,
			}
			testutil.AssertValue(t, got, test.want, "goldenName")
		})
	}
}

func TestReadGoldenFiles(t *testing.T) {
	type args struct {
		directory func(*testing.T) string
	}
	type want struct {
		names      []string
		duplicates []duplicate
		failed     bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "reads text files",
			args: args{
				directory: func(t *testing.T) string {
					directory := t.TempDir()
					files := map[string]string{
						"b.txt":     "same",
						"a.txt":     "same",
						"style.css": "ignored",
					}
					for name, contents := range files {
						writeFile(t, filepath.Join(directory, name), contents)
					}
					return directory
				},
			},
			want: want{
				names:      []string{"a", "b"},
				duplicates: []duplicate{{names: []string{"a", "b"}}},
			},
		},
		{
			name: "invalid pattern",
			args: args{
				directory: func(t *testing.T) string {
					return filepath.Join(t.TempDir(), "[")
				},
			},
			want: want{
				duplicates: []duplicate{},
				failed:     true,
			},
		},
		{
			name: "unreadable text file",
			args: args{
				directory: func(t *testing.T) string {
					directory := t.TempDir()
					if err := os.Mkdir(filepath.Join(directory, "value.txt"), 0o700); err != nil {
						t.Fatal(err)
					}
					return directory
				},
			},
			want: want{
				duplicates: []duplicate{},
				failed:     true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			names, hashes, err := readGoldenFiles(test.args.directory(t))
			got := want{
				names:      names,
				duplicates: resolveDuplicates(hashes),
				failed:     err != nil,
			}
			testutil.AssertValue(t, got, test.want, "readGoldenFiles")
		})
	}
}

func TestParseGoldenOptions(t *testing.T) {
	type args struct {
		contents string
	}
	type want struct {
		options []string
		failed  bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "parses option functions",
			args: args{
				contents: `package fixture

func WithBeta() {}
func NewTable() {}
func WithAlpha() {}

type option struct{}

func (option) WithMethod() {}
`,
			},
			want: want{
				options: []string{"WithAlpha", "WithBeta"},
			},
		},
		{
			name: "invalid source",
			args: args{
				contents: "package",
			},
			want: want{
				failed: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path := writeTestFile(t, "option.go", test.args.contents)
			options, err := parseGoldenOptions(path)
			got := want{
				options: options,
				failed:  err != nil,
			}
			testutil.AssertValue(t, got, test.want, "parseOptions")
		})
	}
}

func TestResolveDuplicates(t *testing.T) {
	type args struct {
		hashes func() map[[sha256.Size]byte][]string
	}
	tests := []struct {
		name string
		args args
		want []duplicate
	}{
		{
			name: "sorts duplicate groups",
			args: args{
				hashes: func() map[[sha256.Size]byte][]string {
					first := [sha256.Size]byte{1}
					second := [sha256.Size]byte{2}
					third := [sha256.Size]byte{3}
					return map[[sha256.Size]byte][]string{
						first:  {"d", "c"},
						second: {"b", "a"},
						third:  {"single"},
					}
				},
			},
			want: []duplicate{
				{names: []string{"a", "b"}},
				{names: []string{"c", "d"}},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := resolveDuplicates(test.args.hashes())
			testutil.AssertValue(t, got, test.want, "resolveDuplicates")
		})
	}
}

func TestResolveUncoveredPairs(t *testing.T) {
	type args struct {
		options []string
		tests   []goldenTest
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "returns pairs absent from tests",
			args: args{
				options: []string{"WithA", "WithB", "WithC"},
				tests: []goldenTest{
					{
						options: []string{"WithA", "WithC"},
					},
				},
			},
			want: []string{"WithA + WithB", "WithB + WithC"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := resolveUncoveredPairs(test.args.options, test.args.tests)
			testutil.AssertValue(t, got, test.want, "resolveUncoveredPairs")
		})
	}
}

func TestCountCalls(t *testing.T) {
	type args struct {
		tests []goldenTest
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "counts resolved names",
			args: args{
				tests: []goldenTest{
					{
						name: "common_value",
					},
					{},
				},
			},
			want: 1,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := countCalls(test.args.tests)
			testutil.AssertValue(t, got, test.want, "countCalls")
		})
	}
}

func TestWriteItems(t *testing.T) {
	type args struct {
		label string
		items []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "writes each item",
			args: args{
				label: "items",
				items: []string{"a", "b"},
			},
			want: "items: a\nitems: b\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var w bytes.Buffer
			writeItems(&w, test.args.label, test.args.items)
			testutil.AssertValue(t, w.String(), test.want, "writeItems")
		})
	}
}

func canceledContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}

func newAuditMedium(t *testing.T) string {
	t.Helper()
	directory := t.TempDir()
	writeAuditMedium(t, directory)
	return directory
}

func writeAuditMedium(t *testing.T, directory string) {
	t.Helper()
	writeFile(t, filepath.Join(directory, "golden_test.go"), `package fixture

func TestGolden_TableValue(t *testing.T) {
	WithBeta()
	WithAlpha()
	testutil.AssertGolden(t, "common_value", nil)
}

func TestGolden_StreamValue(t *testing.T) {
	WithAlpha()
	WithBeta()
	testutil.AssertGolden(t, "common_value", nil)
}
`)
	writeFile(t, filepath.Join(directory, "option.go"), `package fixture

func WithBeta() {}
func WithAlpha() {}
`)
	testdata := filepath.Join(directory, "testdata")
	if err := os.Mkdir(testdata, 0o700); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(testdata, "common_value.txt"), "value")
}

func writeTestFile(t *testing.T, name, contents string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	writeFile(t, path, contents)
	return path
}

func writeFile(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}
}
