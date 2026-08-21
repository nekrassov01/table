package main

import (
	"bytes"
	"testing"

	"github.com/nekrassov01/table/internal/testutil"
)

func Test_runner_run(t *testing.T) {
	type fields struct {
		target string
		mode   string
		data   string
	}
	type want struct {
		targets []string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "runs every example",
			fields: fields{
				target: targetAll,
			},
			want: want{
				targets: []string{
					targetText,
					targetHTML,
					targetMarkdown,
					targetBacklog,
					targetCSV,
				},
			},
		},
		{
			name: "runs all targets",
			fields: fields{
				target: targetAll,
				mode:   modeTable,
				data:   dataSimple,
			},
			want: want{
				targets: []string{
					targetText,
					targetHTML,
					targetMarkdown,
					targetBacklog,
					targetCSV,
				},
			},
		},
		{
			name: "skips targets without selected data",
			fields: fields{
				target: targetAll,
				mode:   modeTable,
				data:   dataFooter,
			},
			want: want{
				targets: []string{
					targetText,
					targetHTML,
					targetBacklog,
					targetCSV,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var output bytes.Buffer
			o := runner{
				w:      &output,
				target: test.fields.target,
				mode:   test.fields.mode,
				data:   test.fields.data,
			}
			if err := o.run(); err != nil {
				t.Fatal(err)
			}
			got := output.Bytes()
			var expected bytes.Buffer
			for index, target := range test.want.targets {
				if index != 0 {
					expected.WriteByte('\n')
				}
				if err := newRunner(&expected, target, test.fields.mode, test.fields.data).run(); err != nil {
					t.Fatal(err)
				}
			}
			want := expected.Bytes()
			testutil.AssertBytes(t, got, want, "run")
		})
	}
}
