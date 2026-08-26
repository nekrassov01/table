package examples

import (
	"bytes"
	"io"
	"testing"

	"github.com/nekrassov01/table/internal/testutil"
)

func TestRun(t *testing.T) {
	type args struct {
		values []string
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
				values: []string{targetText, modeTable, dataASCII},
			},
			want: want{
				output: true,
			},
		},
		{
			name: "missing arguments",
			want: want{
				err: true,
			},
		},
		{
			name: "too many arguments",
			args: args{
				values: []string{targetText, modeTable, dataASCII, "extra"},
			},
			want: want{
				err: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var output bytes.Buffer
			err := Run(&output, test.args.values...)
			got := want{
				output: output.Len() > 0,
				err:    err != nil,
			}
			testutil.AssertValue(t, got, test.want, "Run")
		})
	}
}

func Test_newRunner(t *testing.T) {
	type args struct {
		values []string
	}
	type want struct {
		target string
		mode   string
		data   string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "target",
			args: args{
				values: []string{targetText},
			},
			want: want{
				target: targetText,
			},
		},
		{
			name: "mode",
			args: args{
				values: []string{targetText, modeTable},
			},
			want: want{
				target: targetText,
				mode:   modeTable,
			},
		},
		{
			name: "data",
			args: args{
				values: []string{targetText, modeTable, dataSimple},
			},
			want: want{
				target: targetText,
				mode:   modeTable,
				data:   dataSimple,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := newRunner(&bytes.Buffer{}, test.args.values...)
			got := want{
				target: o.target,
				mode:   o.mode,
				data:   o.data,
			}
			testutil.AssertValue(t, got, test.want, "newRunner")
		})
	}
}

func Test_runner_run(t *testing.T) {
	type fields struct {
		w      func() io.Writer
		target string
		mode   string
		data   string
	}
	type want struct {
		output bool
		err    bool
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "all examples",
			fields: fields{
				target: targetAll,
			},
			want: want{
				output: true,
			},
		},
		{
			name: "selected example",
			fields: fields{
				target: targetText,
				mode:   modeStream,
				data:   dataASCII,
			},
			want: want{
				output: true,
			},
		},
		{
			name: "all targets skip unsupported data",
			fields: fields{
				target: targetAll,
				mode:   modeTable,
				data:   dataASCII,
			},
			want: want{
				output: true,
			},
		},
		{
			name: "unknown mode",
			fields: fields{
				target: targetText,
				mode:   "missing",
			},
			want: want{
				err: true,
			},
		},
		{
			name: "unknown target",
			fields: fields{
				target: "missing",
			},
			want: want{
				err: true,
			},
		},
		{
			name: "unknown data",
			fields: fields{
				target: targetText,
				data:   "missing",
			},
			want: want{
				err: true,
			},
		},
		{
			name: "unavailable data",
			fields: fields{
				target: targetCSV,
				data:   dataRowspan,
			},
			want: want{
				err: true,
			},
		},
		{
			name: "render error",
			fields: fields{
				w: func() io.Writer {
					return &testutil.ErrorWriter{Err: testutil.NewError()}
				},
				target: targetText,
				mode:   modeTable,
				data:   dataASCII,
			},
			want: want{
				err: true,
			},
		},
		{
			name: "separator error",
			fields: fields{
				w: func() io.Writer {
					return &testutil.MatchErrorWriter{Value: "\n", Err: testutil.NewError()}
				},
				target: targetText,
				mode:   modeTable,
			},
			want: want{
				err: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var output bytes.Buffer
			writer := io.Writer(&output)
			if test.fields.w != nil {
				writer = test.fields.w()
			}
			o := runner{
				w:      writer,
				target: test.fields.target,
				mode:   test.fields.mode,
				data:   test.fields.data,
			}
			err := o.run()
			got := want{
				output: output.Len() > 0,
				err:    err != nil,
			}
			testutil.AssertValue(t, got, test.want, "run")
		})
	}
}

func Test_runner_runExample(t *testing.T) {
	type fields struct {
		target string
		mode   string
		data   string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "text",
			fields: fields{
				target: targetText,
				mode:   modeTable,
				data:   dataASCII,
			},
		},
		{
			name: "html",
			fields: fields{
				target: targetHTML,
				mode:   modeTable,
				data:   dataSimple,
			},
		},
		{
			name: "markdown",
			fields: fields{
				target: targetMarkdown,
				mode:   modeTable,
				data:   dataSimple,
			},
		},
		{
			name: "backlog",
			fields: fields{
				target: targetBacklog,
				mode:   modeTable,
				data:   dataSimple,
			},
		},
		{
			name: "csv",
			fields: fields{
				target: targetCSV,
				mode:   modeTable,
				data:   dataSimple,
			},
		},
		{
			name: "unknown data",
			fields: fields{
				target: targetText,
				data:   "missing",
			},
			want: true,
		},
		{
			name: "unknown target",
			fields: fields{
				target: "missing",
				data:   dataSimple,
			},
			want: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := runner{
				w:      &bytes.Buffer{},
				target: test.fields.target,
				mode:   test.fields.mode,
				data:   test.fields.data,
			}
			got := o.runExample() != nil
			testutil.AssertValue(t, got, test.want, "runExample")
		})
	}
}

func Test_runner_runText(t *testing.T) {
	testRunnerTarget(t, targetText, dataASCII, func(o runner, rows [][]any) error {
		return o.runText(rows)
	})
}

func Test_runner_runHTML(t *testing.T) {
	testRunnerTarget(t, targetHTML, dataSimple, func(o runner, rows [][]any) error {
		return o.runHTML(rows)
	})
}

func Test_runner_runMarkdown(t *testing.T) {
	testRunnerTarget(t, targetMarkdown, dataSimple, func(o runner, rows [][]any) error {
		return o.runMarkdown(rows)
	})
}

func Test_runner_runBacklog(t *testing.T) {
	testRunnerTarget(t, targetBacklog, dataSimple, func(o runner, rows [][]any) error {
		return o.runBacklog(rows)
	})
}

func Test_runner_runCSV(t *testing.T) {
	testRunnerTarget(t, targetCSV, dataSimple, func(o runner, rows [][]any) error {
		return o.runCSV(rows)
	})
}

func Test_newExample(t *testing.T) {
	type args struct {
		rows [][]any
	}
	tests := []struct {
		name string
		args args
		want [][]any
	}{
		{
			name: "rows",
			args: args{
				rows: [][]any{{"value"}},
			},
			want: [][]any{{"value"}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := newExample(test.args.rows).rows
			testutil.AssertValue(t, got, test.want, "newExample")
		})
	}
}

func Test_example_run(t *testing.T) {
	testError := testutil.NewError()
	type fields struct {
		rows     [][]any
		tabular  *testutil.Tabular
		streamer *testutil.Streamer
	}
	type want struct {
		rows [][]any
		err  bool
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "table",
			fields: fields{
				rows:    [][]any{{"first"}, {"second"}},
				tabular: &testutil.Tabular{},
			},
			want: want{
				rows: [][]any{{"first"}, {"second"}},
			},
		},
		{
			name: "table error",
			fields: fields{
				rows: [][]any{{"value"}},
				tabular: &testutil.Tabular{
					Err: testError,
				},
			},
			want: want{
				rows: [][]any{{"value"}},
				err:  true,
			},
		},
		{
			name: "stream",
			fields: fields{
				rows:     [][]any{{"first"}, {"second"}},
				streamer: &testutil.Streamer{},
			},
			want: want{
				rows: [][]any{{"first"}, {"second"}},
			},
		},
		{
			name: "stream render error",
			fields: fields{
				rows: [][]any{{"value"}},
				streamer: &testutil.Streamer{
					RenderErr: testError,
				},
			},
			want: want{
				rows: [][]any{{"value"}},
				err:  true,
			},
		},
		{
			name: "stream close error",
			fields: fields{
				streamer: &testutil.Streamer{
					CloseErr: testError,
				},
			},
			want: want{
				err: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := example{
				rows: test.fields.rows,
			}
			if test.fields.tabular != nil {
				o.tabular = test.fields.tabular
			}
			if test.fields.streamer != nil {
				o.streamer = test.fields.streamer
			}
			err := o.run()
			var rows [][]any
			if test.fields.tabular != nil {
				rows = test.fields.tabular.Rows
			}
			if test.fields.streamer != nil {
				rows = test.fields.streamer.Rows
			}
			got := want{
				rows: rows,
				err:  err != nil,
			}
			testutil.AssertValue(t, got, test.want, "run")
		})
	}
}

func Test_dataNames(t *testing.T) {
	type args struct {
		target string
	}
	type want struct {
		count int
		ok    bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{name: "text", args: args{target: targetText}, want: want{count: 9, ok: true}},
		{name: "html", args: args{target: targetHTML}, want: want{count: 7, ok: true}},
		{name: "markdown", args: args{target: targetMarkdown}, want: want{count: 5, ok: true}},
		{name: "backlog", args: args{target: targetBacklog}, want: want{count: 7, ok: true}},
		{name: "csv", args: args{target: targetCSV}, want: want{count: 5, ok: true}},
		{name: "unknown", args: args{target: "missing"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			names, ok := dataNames(test.args.target)
			got := want{
				count: len(names),
				ok:    ok,
			}
			testutil.AssertValue(t, got, test.want, "dataNames")
		})
	}
}

func Test_exampleData(t *testing.T) {
	type args struct {
		name string
	}
	type want struct {
		rows int
		ok   bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{name: "ascii", args: args{name: dataASCII}, want: want{rows: len(SimpleData.Body), ok: true}},
		{name: "simple", args: args{name: dataSimple}, want: want{rows: len(SimpleData.Body), ok: true}},
		{name: "compact", args: args{name: dataCompact}, want: want{rows: len(CompactData.Body), ok: true}},
		{name: "rowspan", args: args{name: dataRowspan}, want: want{rows: len(RowspanData.Body), ok: true}},
		{name: "colspan", args: args{name: dataColspan}, want: want{rows: len(ColspanData.Body), ok: true}},
		{name: "footer", args: args{name: dataFooter}, want: want{rows: len(FooterData.Body), ok: true}},
		{name: "transformer", args: args{name: dataTransformer}, want: want{rows: len(FooterData.Body), ok: true}},
		{name: "complex", args: args{name: dataComplex}, want: want{rows: len(ComplexData.Body), ok: true}},
		{name: "stacked header", args: args{name: dataStackedHeader}, want: want{rows: len(StackedHeaderData.Body), ok: true}},
		{name: "comma included", args: args{name: dataCommaIncluded}, want: want{rows: len(CommaIncludedData.Body), ok: true}},
		{name: "unknown", args: args{name: "missing"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data, ok := exampleData(test.args.name)
			got := want{
				rows: len(data.Body),
				ok:   ok,
			}
			testutil.AssertValue(t, got, test.want, "exampleData")
		})
	}
}

func Test_newError(t *testing.T) {
	type args struct {
		name   string
		target string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "unavailable data",
			args: args{
				name:   dataASCII,
				target: targetCSV,
			},
			want: `examples: data "ascii" is not available for target "csv"`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := newError(test.args.name, test.args.target).Error()
			testutil.AssertValue(t, got, test.want, "newError")
		})
	}
}

func testRunnerTarget(t *testing.T, target, data string, run func(runner, [][]any) error) {
	t.Helper()
	type fields struct {
		mode string
		data string
		rows [][]any
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "table",
			fields: fields{
				mode: modeTable,
				data: data,
				rows: SimpleData.Body,
			},
		},
		{
			name: "stream",
			fields: fields{
				mode: modeStream,
				data: data,
				rows: SimpleData.Body,
			},
		},
		{
			name: "transformer fallback",
			fields: fields{
				mode: modeTable,
				data: dataTransformer,
				rows: [][]any{make([]any, len(FooterData.Header[0]))},
			},
		},
		{
			name: "complex fallback",
			fields: fields{
				mode: modeTable,
				data: dataComplex,
				rows: [][]any{make([]any, len(ComplexData.Header[0]))},
			},
		},
		{
			name: "unknown data",
			fields: fields{
				mode: modeTable,
				data: "missing",
			},
			want: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := runner{
				w:      &bytes.Buffer{},
				target: target,
				mode:   test.fields.mode,
				data:   test.fields.data,
			}
			err := run(o, test.fields.rows)
			got := err != nil
			testutil.AssertValue(t, got, test.want, "run target")
		})
	}
}
