package catalog

import (
	"testing"

	"github.com/nekrassov01/table/internal/testutil"
)

func Test_codeBlock(t *testing.T) {
	type args struct {
		language string
		value    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "source",
			args: args{
				language: "go",
				value:    "var value = 1\n",
			},
			want: "````go\nvar value = 1\n````",
		},
		{
			name: "nested fence",
			args: args{
				language: "markdown",
				value:    "````text\nvalue\n````\n",
			},
			want: "`````markdown\n````text\nvalue\n````\n`````",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := codeBlock(test.args.language, test.args.value)
			testutil.AssertValue(t, got, test.want, "codeBlock")
		})
	}
}

func Test_exampleCommands(t *testing.T) {
	type args struct {
		target string
		data   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "table and stream",
			args: args{
				target: "text",
				data:   "simple",
			},
			want: "````sh\nmake example target=text mode=table data=simple\nmake example target=text mode=stream data=simple\n````",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := exampleCommands(test.args.target, test.args.data)
			testutil.AssertValue(t, got, test.want, "exampleCommands")
		})
	}
}

func Test_outputBlock(t *testing.T) {
	type args struct {
		target string
		value  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "text",
			args: args{
				target: "text",
				value:  "value\n",
			},
			want: "````text\nvalue\n````",
		},
		{
			name: "html",
			args: args{
				target: "html",
				value:  "<table></table>\n",
			},
			want: "````html\n<table></table>\n````",
		},
		{
			name: "markdown",
			args: args{
				target: "markdown",
				value:  "| value |\n",
			},
			want: "````markdown\n| value |\n````",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := outputBlock(test.args.target, test.args.value)
			testutil.AssertValue(t, got, test.want, "outputBlock")
		})
	}
}
